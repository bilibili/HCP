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
	"os"
)

// Executable 任务执行器接口
type Executable interface {
	// Execute 执行任务
	Execute(ctx context.Context) error
	// Commit 提交任务
	Commit(ctx context.Context) error
	// Cancel 取消任务
	Cancel()
	// ExecutableID 获取任务id
	ExecutableID() int
	// TimeCost 获取任务执行时间
	TimeCost() float64
	// ExecutableStatus 获取任务状态
	ExecutableStatus() string
	// ExecutableType 获取任务类型
	ExecutableType() string
	// ScheduleType 获取任务调度类型
	ScheduleType() string
	// Container 获取任务所在容器
	Container() string
}

// Inputtable 任务输入器接口
type Inputtable interface {
	// ExecutableType 获取任务类型
	ExecutableType() string
	// Input 获取任务
	Input(ctx context.Context) ([]Executable, error)
	// InputCloser 任务输入是否关闭
	InputCloser(ctx context.Context) bool
	// InputRecover 任务输入恢复
	InputRecover(ctx context.Context) ([]Executable, error)
}

// DefaultExecutable 默认任务执行器嵌入Executable接口
type DefaultExecutable struct {
	Executable
}

// Container 获取任务所在容器
func (d *DefaultExecutable) Container() string {
	return os.Getenv("POD_IP")
}
