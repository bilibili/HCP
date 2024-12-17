// Package enum
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
package enum

import (
	"github.com/pkg/errors"
)

const (
	ScheduleTypePipeline     = "pipeline"
	ScheduleTypeTrigger      = "trigger"
	ExecutableTypeJob        = "job"
	ExecutableTypeAuthJob    = "auth_job"
	ExecutableTypeStep       = "step"
	ExecutableTypeObject     = "object"
	ExecutableTypeJobCreator = "job_creator"
	ExecutableTypeOverdue    = "overdue"

	ExecutableStateSuccess       = "success"
	ExecutableStateFailure       = "failure"
	ExecutableJobStateRunning    = "running"
	ExecutableJobStateConfirming = "confirming"
	ExecutableJobStateIdle       = "idle"
	ExecutableStateAuth          = "auth"
	ExecutableStateAuthReq       = "auth_req"
	ExecutableStateAuthWait      = "auth_wait"
	ExecutableStateAuthDeny      = "auth_deny"
	ExecutableStatusAuthSleep    = "auth_sleep"
	ExecutablesStateAuthWaitSync = "auth_wait_sync"
	ExecutableStateSuspend       = "suspend"
	ExecutableStateClose         = "close"
	ExecutableStatusCreateWait   = "create_wait"
	ExecutableStatusCreating     = "creating"

	ResourceTypeCloudServer      = "cloud_server"
	ResourceTypeCloudServerImage = "cloud_server_image"
)

var (
	ErrDataNotFound = errors.New("data not found")
)

var (
	ExecutableStateDisplayMapping = map[string]string{
		ExecutableStateSuccess:       "成功",
		ExecutableJobStateRunning:    "进行中",
		ExecutableJobStateConfirming: "待确认",
		ExecutableStateFailure:       "失败",
		ExecutableJobStateIdle:       "等待",
		ExecutableStateAuth:          "审核中",
		ExecutableStateAuthReq:       "审核中",
		ExecutableStateAuthWait:      "审核中",
		ExecutableStateAuthDeny:      "已拒绝",
		ExecutableStateSuspend:       "已暂停",
		ExecutableStateClose:         "关闭",
		ExecutableStatusCreateWait:   "创建中",
		ExecutableStatusCreating:     "创建中",
	}
)
