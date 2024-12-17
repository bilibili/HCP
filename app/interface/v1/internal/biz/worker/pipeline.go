// Package worker
/*
 * Copyright 2024-2025 Bilibili Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package worker

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/pkg/log"
)

/*
 1. Pipeline 可以设置是否有 Inputtable
 2. 有 Inputtable 的 Pipeline 将自动进行本地调度 Executable，并在本地执行，不支持分布式特性
 3. 无 Inputtable 的 Pipeline 将等待外部派发 Executable，并在本地执行，支持分布式特性
*/

type Pipe <-chan Executable

// Pipeline 管道接口
type Pipeline interface {
	Executables() (list []Executable)
	Run()
	Shutdown()
	Status() string
	Load(...Executable)
	Unload(...string)
}

// pipeline 参数结构体
type pipeline struct {
	mutex       *sync.Mutex
	name        string
	concurrency int
	src         Inputtable
	issue       chan Executable
	cancel      chan struct{}
	status      string
	executables *sync.Map
}

// NewPipeline 创建 Pipeline
func NewPipeline(name string, concurrency int, src Inputtable) Pipeline {
	return &pipeline{
		mutex:       &sync.Mutex{},
		name:        name,
		concurrency: concurrency,
		src:         src,
		status:      "idle",
		executables: &sync.Map{},
	}
}

// Executables 插入管道数据
func (p *pipeline) Executables() (list []Executable) {
	p.executables.Range(func(key, value interface{}) bool { //遍历executables的Map向接口Executable切片中插入数据
		list = append(list, value.(Executable))
		return true
	})
	return
}

// Run 执行
func (p *pipeline) Run() {
	p.cancel = make(chan struct{})
	p.status = "running"
	p.issue = make(chan Executable, 1000)

	var issuePipe Pipe //数据执行的管道
	if p.src != nil {
		issuePipe = p.issuer() //获取管道数据
		log.Warn("pipelining %v", p.src.ExecutableType())
	} else {
		issuePipe = p.issue
	}
	time.Sleep(1 * time.Second)
	execPipes := make([]Pipe, 0)
	for i := 0; i < p.concurrency; i += 1 {
		execPipes = append(execPipes, p.executor(i, issuePipe))
	}
	time.Sleep(1 * time.Second)
	for exe := range p.committer(execPipes...) {
		log.Warn("executable[%v] accomplished, total time: %v", exe.ExecutableID(), exe.TimeCost())
	}
	if p.status != "running" {
		return
	}
	p.status = "idle"
}

// Shutdown 关闭
func (p *pipeline) Shutdown() {
	p.status = "shutdown"
	close(p.cancel)
	wg := &sync.WaitGroup{}
	p.executables.Range(func(key, value interface{}) bool {
		wg.Add(1)
		exe, ok := p.executables.Load(key)
		go func() {
			exe.(Executable).Cancel()
			for ok {
				_, ok = p.executables.Load(key)
			}
			wg.Done()
		}()
		return true
	})
	wg.Wait()
}

// Status 状态
func (p *pipeline) Status() string {
	return p.status
}

// Load 加载
func (p *pipeline) Load(input ...Executable) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.status != "running" || p.src != nil {
		return
	}
	loaded := []string{}
	for _, exe := range input {
		_, ok := p.executables.LoadOrStore(exe.ExecutableID(), exe)
		if ok {
			continue
		}
		loaded = append(loaded, fmt.Sprintf("%v-%v", exe.ExecutableType(), exe.ExecutableID()))
		p.issue <- exe
	}
	log.Warn("%v loaded %v executables [%v]", p.name, len(loaded), strings.Join(loaded, ","))
}

// Unload 卸载
func (p *pipeline) Unload(ids ...string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	wg := &sync.WaitGroup{}
	unloaded := []string{}
	for i := range ids {
		id := ids[i]
		exe, ok := p.executables.Load(id)
		if !ok {
			continue
		}
		wg.Add(1)
		go func() {
			exe.(Executable).Cancel()
			for ok {
				_, ok = p.executables.Load(id)
			}
			wg.Done()
		}()
		unloaded = append(unloaded, fmt.Sprintf("%v-%v", exe.(Executable).ExecutableType(), exe.(Executable).ExecutableID()))
	}
	wg.Wait()
	log.Warn("%v unloaded %v executables [%v]", p.name, len(unloaded), strings.Join(unloaded, ","))
}

// issuer 获取管道数据
func (p *pipeline) issuer() Pipe {
	out := make(chan Executable, 1000) //分配内存并初始化实例
	go func() {
		defer close(out)
		log.Warn("%v's issuer is up", p.name)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		//获取需要恢复的数据
		restore, err := p.src.InputRecover(ctx)
		if err != nil {
			log.Error("%v's issuer: load err %v", p.name, err)
		}
		log.Warn("%v's issuer: restore %v %vs", p.name, len(restore), p.src.ExecutableType())
		for _, exe := range restore {
			out <- exe
			p.executables.Store(exe.ExecutableID(), exe) //存储到map中
			log.Warn("%v's issuer: executable %v-%v restored", p.name, exe.ExecutableType(), exe.ExecutableID())
		}
		for {
			select {
			case <-p.cancel: //如果p.cancel通道关闭，则返回
				log.Warn("%v' issuer canceled", p.name)
				return
			default:
				if p.src.InputCloser(ctx) {
					log.Warn("%v' issuer is down", p.name)
					return
				}
				//触发执行批次
				input, err := p.src.Input(ctx)
				if err != nil {
					log.Error("%v's issuer: issue err: %v", p.name, err)
					time.Sleep(5 * time.Second)
					continue
				}
				if len(input) == 0 {
					time.Sleep(2 * time.Second)
					continue
				}
				for _, exe := range input {
					out <- exe
					p.executables.Store(exe.ExecutableID(), exe)
					log.Warn("%v's issuer: executable %v-%v issued", p.name, exe.ExecutableType(), exe.ExecutableID())
				}
			}
		}
	}()
	return out
}

// executor 执行
func (p *pipeline) executor(id int, pipe Pipe) Pipe {
	out := make(chan Executable, 999)
	go func() {
		defer close(out)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		for exe := range pipe {
			log.Warn("%v's executor-%v: executable %v-%v started", p.name, id, exe.ExecutableType(), exe.ExecutableID())
			if err := exe.Execute(ctx); err != nil {
				log.Error("%v's executor-%v: executable %v-%v execute err: %v", p.name, id, exe.ExecutableType(), exe.ExecutableID(), err)
			}
			log.Warn("%v's executor-%v: executable %v-%v exited", p.name, id, exe.ExecutableType(), exe.ExecutableID())
			out <- exe
		}
		log.Warn("%v's executor-%v is down", p.name, id)
	}()
	return out
}

// committer 提交
func (p *pipeline) committer(pipes ...Pipe) Pipe {
	out := make(chan Executable, 999)
	wg := &sync.WaitGroup{}
	wg.Add(len(pipes))
	output := func(pipe Pipe) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		for exe := range pipe {
			log.Warn("%v's committer: executable %v-%v started", p.name, exe.ExecutableType(), exe.ExecutableID())
			p.executables.Delete(exe.ExecutableID())
			if err := exe.Commit(ctx); err != nil {
				log.Error("%v's committer: executable %v-%v commit err: %v", p.name, exe.ExecutableType(), exe.ExecutableID(), err)
			}
			log.Warn("%v's committer: executable %v-%v exited", p.name, exe.ExecutableType(), exe.ExecutableID())
			out <- exe
		}
		wg.Done()
	}
	for _, pipe := range pipes {
		go output(pipe)
	}
	go func() {
		log.Warn("%v's committer is up", p.name)
		wg.Wait()
		close(out)
		log.Warn("%v's committer is down", p.name)
	}()
	return out
}
