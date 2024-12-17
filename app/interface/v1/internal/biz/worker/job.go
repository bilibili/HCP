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
	"time"

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/app/interface/v1/internal/biz/enum"
	"github.com/pkg/errors"
)

/*
JobPipelineExecutable: JobPipelineInput 目前仅使用Input、InputRecovers、InputWait方法 传入Job 数据
DistributedJobPipelineExecutable 具体执行Job的方法
*/

// JobPipelineExecutable Job可执行结构体
type JobPipelineExecutable struct {
	// Job 任务结构体
	*biz.Job
	// DefaultExecutable 默认可执行对象
	*DefaultExecutable

	// pipeline 任务执行管道
	pipeline Pipeline
}

// NewJobPipelineExecutable 创建Job可执行对象
func NewJobPipelineExecutable(do *biz.Job) *JobPipelineExecutable {
	return &JobPipelineExecutable{Job: do}
}

// Execute 执行Job
func (jp *JobPipelineExecutable) Execute(ctx context.Context) error {
	maxRunningStep := 1
	if jp.ParallelStepEnable {
		maxRunningStep = jp.TotalStep
	}
	jp.pipeline = NewPipeline(fmt.Sprintf("job-%v step execution pipeline", jp.Id), maxRunningStep, NewStepSource(jp.Job))
	jp.pipeline.Run()
	if jp.pipeline.Status() == "shutdown" {
		return errors.Errorf("cancel by signal.")
	}
	jp.Status = enum.ExecutableStateSuccess
	jp.EndTime = time.Now()
	return nil
}

// Commit 提交Job
func (jp *JobPipelineExecutable) Commit(ctx context.Context) error {
	if jp.pipeline.Status() == "shutdown" {
		return errors.Errorf("cancel by signal.")
	}
	err := jobUseCase.UpdateJob(jp.Job, -1)
	return err
}

// Cancel 取消Job
func (jp *JobPipelineExecutable) Cancel() {
	jp.pipeline.Shutdown()
}

// ExecutableID 获取JobID
func (jp *JobPipelineExecutable) ExecutableID() int {
	return jp.Id
}

// TimeCost 获取Job执行时间
func (jp *JobPipelineExecutable) TimeCost() float64 {
	return jp.EndTime.Sub(jp.StartTime).Minutes()
}

// ExecutableStatus 获取Job状态
func (jp *JobPipelineExecutable) ExecutableStatus() string {
	return jp.Status
}

// ExecutableType 获取Job类型
func (jp *JobPipelineExecutable) ExecutableType() string {
	return enum.ExecutableTypeJob
}

// ScheduleType 获取Job调度类型
func (jp *JobPipelineExecutable) ScheduleType() string {
	return enum.ScheduleTypePipeline
}

// JobPipelineInput 获取Job
type JobPipelineInput struct{}

// NewJobPipelineInput 创建Job可执行对象
func NewJobPipelineInput() *JobPipelineInput {
	return &JobPipelineInput{}
}

// ExecutableType 获取Job类型
func (jp *JobPipelineInput) ExecutableType() string {
	return enum.ExecutableTypeJob
}

// Input 获取Job
func (jp *JobPipelineInput) Input(ctx context.Context) ([]Executable, error) {
	list := make([]Executable, 0)
	count, err := jobUseCase.CountJobByJobStatus(ctx, enum.ExecutableJobStateIdle)
	if err != nil || count == 0 {
		return list, err
	}
	doList, err := jobUseCase.ListIdleJob(100)
	if err != nil {
		return list, err
	}
	idleId := make([]int, 0)
	now := time.Now()
	for _, do := range doList {
		exe := NewJobPipelineExecutable(do)
		exe.Status = enum.ExecutableJobStateRunning
		exe.Reviewers = ""
		exe.StartTime = now
		idleId = append(idleId, do.Id)
		list = append(list, exe)
		time.Sleep(10 * time.Millisecond)
	}
	if len(idleId) == 0 {
		return list, nil
	}
	if err := jobUseCase.BatchUpdateJob(idleId, enum.ExecutableJobStateRunning, now, []string{}); err != nil {
		return nil, err
	}
	return list, nil
}

// InputCloser 输入是否关闭
func (jp *JobPipelineInput) InputCloser(ctx context.Context) bool {
	return false
}

// InputRecover 输入恢复
func (jp *JobPipelineInput) InputRecover(ctx context.Context) ([]Executable, error) {
	list, err := jobUseCase.ListRunningJob()
	if err != nil {
		return []Executable{}, err
	}
	restore := make([]Executable, 0)
	for _, do := range list {
		exe := NewJobPipelineExecutable(do)
		restore = append(restore, exe)
	}
	return restore, nil
}
