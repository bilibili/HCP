// Package cloudrepo
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
package cloudrepo

const (
	CloudVMStatusPending      = "pending"       // 创建中
	CloudVMStatusRunning      = "running"       // 运行中
	CloudVMStatusStarting     = "starting"      // 启动中
	CloudVMStatusStopping     = "stopping"      // 停止中
	CloudVMStatusSuspend      = "suspend"       // 暂停中
	CloudVMStatusStopped      = "stopped"       // 已停止
	CloudVMStatusTerminated   = "terminated"    // 已销毁
	CloudVMStatusShuttingDown = "shutting-down" // 关机中
	CloudVMStatusChangeFlavor = "change_flavor" // 变更规格中
	CloudVMStatusDeploying    = "deploying"     // 部署中
	CloudVMStatueDeleted      = "deleted"       // 已删除
	CloudVMStatusOther        = "other"         // 其他状态

	CloudVmChargeTypePrePaid  = "prepaid"  // 预付费/包年包月
	CloudVmChargeTypePostPaid = "postpaid" // 后付费/按需付费
)
