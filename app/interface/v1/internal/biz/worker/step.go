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

// StepPipelineExecutable 批次任务执行结构体
type StepPipelineExecutable struct {
	// Step 批次
	*biz.Step
	// DefaultExecutable 默认可执行对象
	*DefaultExecutable

	// pipeline 执行管道
	pipeline Pipeline
}

// NewStepPipelineExecutable 创建批次任务执行对象
func NewStepPipelineExecutable(do *biz.Step) *StepPipelineExecutable {
	return &StepPipelineExecutable{Step: do}
}

// Execute 执行批次任务
func (se *StepPipelineExecutable) Execute(ctx context.Context) error {
	if se.Step.Status == enum.ExecutableStateSuccess {
		return nil
	}
	se.pipeline = NewPipeline(fmt.Sprintf("step-%v object execution pipeline", se.Id), se.StepTotalObject, NewObjectTriggerInput(se.Step))
	se.pipeline.Run()
	if se.pipeline.Status() == "shutdown" {
		return errors.Errorf("cancel by signal.")
	}
	se.EndTime = time.Now()
	se.Status = enum.ExecutableStateSuccess
	return nil
}

// Commit 批次任务执行结果
func (se *StepPipelineExecutable) Commit(ctx context.Context) error {
	if se.pipeline.Status() == "shutdown" {
		return errors.Errorf("cancel by signal.")
	}
	if se.Status == enum.ExecutableStateFailure {
		se.Status = enum.ExecutableJobStateIdle
	}
	if err := jobUseCase.UpdateStep(se.Step); err != nil {
		return err
	}
	job, err := jobUseCase.GetJobByJobId(se.JobId)
	if err != nil {
		return err
	}
	job.CurrentStep += 1
	if err := jobUseCase.UpdateJob(job, job.CurrentStep); err != nil {
		return err
	}
	return nil
}

// Cancel 取消批次任务
func (se *StepPipelineExecutable) Cancel() {
	se.pipeline.Shutdown()
}

// ExecutableID 获取批次任务执行id
func (se *StepPipelineExecutable) ExecutableID() int {
	return se.Id
}

// TimeCost 获取批次任务执行时间
func (se *StepPipelineExecutable) TimeCost() float64 {
	return se.EndTime.Sub(se.StartTime).Minutes()
}

// ExecutableStatus 获取批次任务执行状态
func (se *StepPipelineExecutable) ExecutableStatus() string {
	return se.Status
}

// ExecutableType 获取批次任务执行类型
func (se *StepPipelineExecutable) ExecutableType() string {
	return enum.ExecutableTypeStep
}

// ScheduleType 获取批次任务执行调度类型
func (se *StepPipelineExecutable) ScheduleType() string {
	return enum.ScheduleTypePipeline
}

// StepPipelineInput 批次任务执行输入
type StepPipelineInput struct {
	job *biz.Job
}

// NewStepSource 创建批次任务执行输入
func NewStepSource(job *biz.Job) *StepPipelineInput {
	return &StepPipelineInput{job: job}
}

// ExecutableType 获取批次任务执行类型
func (ss *StepPipelineInput) ExecutableType() string {
	return enum.ExecutableTypeStep
}

// Input TODO 获取 steps 中忽略 objects
func (ss *StepPipelineInput) Input(ctx context.Context) ([]Executable, error) {
	list := make([]Executable, 0)
	count, err := jobUseCase.CountStepByJobIdAndStepStatus(ss.job.Id, enum.ExecutableJobStateIdle)
	if err != nil || count == 0 {
		return list, err
	}
	if ss.job.ParallelStepEnable {
		steps, err := jobUseCase.ListStepByJobId(ss.job.Id)
		if err != nil {
			return list, err
		}
		for _, s := range steps {
			if s.Status != enum.ExecutableJobStateIdle {
				continue
			}
			s.Status = enum.ExecutableJobStateRunning
			s.StartTime = time.Now()
			if err := jobUseCase.UpdateStep(s); err != nil {
				return list, err
			}
			list = append(list, NewStepPipelineExecutable(s))
		}
		return list, nil
	}
	job, err := jobUseCase.GetJobByJobId(ss.job.Id)
	if err != nil {
		return nil, err
	}
	step, err := jobUseCase.GetIndexStepByJobId(job.CurrentStep, job.Id)
	if err != nil {
		if errors.Is(err, enum.ErrDataNotFound) {
			return nil, nil
		}
		return nil, err
	}
	// 步骤执行中或已完成，则不重复派发，待执行或执行失败，则重新派发
	if step.Status != enum.ExecutableJobStateIdle {
		return nil, nil
	}
	// 步骤执行中或已完成，则不重复派发，待执行或执行失败，则重新派发
	step.Status = enum.ExecutableJobStateRunning
	step.StartTime = time.Now()
	if err := jobUseCase.UpdateStep(step); err != nil {
		return nil, err
	}
	list = append(list, NewStepPipelineExecutable(step))
	return list, nil
}

// InputRecover 获取批次任务执行输入
func (ss *StepPipelineInput) InputRecover(ctx context.Context) ([]Executable, error) {
	if ss.job.CurrentStep == ss.job.TotalStep {
		return []Executable{}, nil
	}
	list := make([]Executable, 0)
	steps, err := jobUseCase.ListStepByJobId(ss.job.Id)
	if err != nil {
		return list, err
	}
	for _, s := range steps {
		if s.Status != enum.ExecutableJobStateRunning {
			continue
		}
		list = append(list, NewStepPipelineExecutable(s))
	}
	return list, nil
}

// InputCloser 获取批次任务执行输入是否完成
func (ss *StepPipelineInput) InputCloser(ctx context.Context) bool {
	if ss.job.ParallelStepEnable {
		count, err := jobUseCase.CountStepByJobIdAndStepStatus(ss.job.Id, enum.ExecutableStateSuccess)
		if err != nil {
			return false
		}
		return ss.job.TotalStep == count
	}
	job, err := jobUseCase.GetJobByJobId(ss.job.Id)
	if err != nil {
		return false
	}
	if job.TotalStep > job.CurrentStep {
		return false
	}
	return true
}
