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
	"time"

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/app/interface/v1/internal/biz/enum"
	"github.com/bilibili/HCP/utils"
	"github.com/go-kratos/kratos/pkg/log"
)

// JobCreatorPipelineExecutable 创建Pipeline可执行结构体
type JobCreatorPipelineExecutable struct {
	// Job 任务结构体
	*biz.Job
	// DefaultExecutable 默认执行器
	*DefaultExecutable
}

// NewJobCreatorPipelineExecutable 创建JobCreator可执行对象
func NewJobCreatorPipelineExecutable(do *biz.Job) *JobCreatorPipelineExecutable {
	return &JobCreatorPipelineExecutable{Job: do}
}

// Execute
// Executable 已知job，默认 scheduler
// 目标：创建 job step、job object
func (jp *JobCreatorPipelineExecutable) Execute(ctx context.Context) error {
	var (
		jobId   = jp.Job.Id
		jobName = jp.Job.Name
	)
	template, err := jobUseCase.GetTemplate(jobName) //通过任务名称获取template
	if err != nil {
		return err
	}
	input, err := jobUseCase.GetTemplateContent(jobId) //获取表单提交的内容
	if err != nil {
		return err
	}
	formInputCopy := make(map[string]interface{})
	for k, v := range input {
		formInputCopy[k] = v
	}

	if err := jobUseCase.CreateJobSteps(jobId, template, input); err != nil { //创建job step
		log.Error("JobCreatorPipelineExecutable: CreateJobSteps (%v)", err)
		return err
	}

	if template.Workflow == "" {
		err = jobUseCase.BatchUpdateJobStatus([]int{jobId}, enum.ExecutableJobStateIdle)
		if err != nil {
			log.Error("JobCreatorPipelineExecutable: BatchUpdateJobStatus (%v)", err)
			return err
		}
	}

	// 任务处理
	return nil
}

// Commit 后置操作
func (jp *JobCreatorPipelineExecutable) Commit(ctx context.Context) error {
	return nil
}

// Cancel 取消操作
func (jp *JobCreatorPipelineExecutable) Cancel() {
	return
}

// ExecutableID 获取可执行id
func (jp *JobCreatorPipelineExecutable) ExecutableID() int {
	return jp.Id
}

// ExecutableStatus 获取可执行状态
func (jp *JobCreatorPipelineExecutable) ExecutableStatus() string {
	return jp.Status
}

// ExecutableType 获取可执行类型
func (jp *JobCreatorPipelineExecutable) ExecutableType() string {
	return enum.ExecutableTypeJobCreator
}

// ScheduleType 获取调度类型
func (jp *JobCreatorPipelineExecutable) ScheduleType() string {
	return enum.ScheduleTypePipeline
}

// TimeCost 获取耗时
func (jp *JobCreatorPipelineExecutable) TimeCost() float64 {
	return jp.EndTime.Sub(jp.StartTime).Minutes()
}

// JobCreatorPipelineInput 创建Pipeline Creator可执行结构体
type JobCreatorPipelineInput struct {
}

// NewJobCreatorPipelineInput 实例化JobCreatorPipelineInput可执行结构体
func NewJobCreatorPipelineInput() *JobCreatorPipelineInput {
	return &JobCreatorPipelineInput{}
}

// ExecutableType 获取可执行类型
func (jp *JobCreatorPipelineInput) ExecutableType() string {
	return enum.ExecutableTypeJobCreator
}

// Input 获取可执行列表
func (jp *JobCreatorPipelineInput) Input(ctx context.Context) ([]Executable, error) {
	list := make([]Executable, 0)
	doList, err := jobUseCase.ListCreateWaitJob(100)
	if err != nil {
		return list, err
	}
	if len(doList) == 0 {
		return list, nil
	}
	ids := make([]int, 0)
	now := time.Now()
	for _, do := range doList {
		exe := NewJobCreatorPipelineExecutable(do)
		exe.Status = enum.ExecutableJobStateRunning
		exe.Reviewers = ""
		exe.StartTime = now
		ids = append(ids, do.Id)
		list = append(list, exe)
		time.Sleep(10 * time.Millisecond)
	}
	if len(ids) == 0 {
		return list, nil
	}
	if err := jobUseCase.BatchUpdateJob(ids, enum.ExecutableStatusCreating, now, []string{}); err != nil {
		return nil, err
	}
	return list, nil
}

// InputCloser 输入是否关闭
func (jp *JobCreatorPipelineInput) InputCloser(ctx context.Context) bool {
	return utils.ContextWait(ctx)
}

// InputRecover 输入恢复
func (jp *JobCreatorPipelineInput) InputRecover(ctx context.Context) ([]Executable, error) {
	list, err := jobUseCase.ListCreatingJob()
	if err != nil {
		return []Executable{}, err
	}
	restore := make([]Executable, 0)
	for _, do := range list {
		exe := NewJobCreatorPipelineExecutable(do)
		restore = append(restore, exe)
	}
	return restore, nil
}
