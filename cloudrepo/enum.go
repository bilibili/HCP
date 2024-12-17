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

var (
	statusAli = map[string]map[string]string{
		"cloud_subnet": {
			"Available": "available",
			"Pending":   "unavailable",
		},
		"cloud_server_type": {
			"g":  "通用型",
			"c":  "计算型",
			"ic": "计算型",
			"r":  "内存型",
			"re": "内存型",
		},
		"cloud_server_image": {
			"Creating":  "creating",
			"Available": "available",
		},
		"type_image": {
			"true":  "public",
			"false": "private",
		},
		"cloud_server": {
			"Running":  "running",
			"Pending":  "pending",
			"Starting": "starting",
			"Stopping": "stopping",
			"Stopped":  "stopped",
		},
		"pay_type": {
			"Postpaid":    "postpaid",
			"Prepaid":     "prepaid",
			"PrePaid":     "prepaid",
			"PostPaid":    "prepaid",
			"PrePay":      "prepaid",
			"PayOnDemand": "postpaid",
		},
		"cloud_pay_resource_renew_status": {
			"AutoRenewal":   "auto",
			"ManualRenewal": "manual",
			"NotRenewal":    "no",
		},
		"renew_type": {
			"AutoRenewal":   "auto",
			"ManualRenewal": "manual",
			"NotRenewal":    "no",
		},
	}
	statusAws = map[string]map[string]string{
		"cloud_zone": {
			"available":   "available",
			"information": "information",
			"impaired":    "impaired",
			"unavailable": "unavailable",
		},
		"cloud_vpc": {
			"available": "available",
			"pending":   "unavailable",
		},
		"cloud_subnet": {
			"available": "available",
			"pending":   "unavailable",
		},
		"cloud_server": {
			"running":      "running",
			"pending":      "pending",
			"available":    "available",
			"invalid":      "invalid",
			"deregistered": "deregistered",
			"transient":    "transient",
			"failed":       "failed",
			"error":        "error",
		},
		"pay_type": {
			"": "postpaid",
		},
		"cloud_region_location": {
			"us-east-2":      "美国东部(俄亥俄州)",
			"us-east-1":      "美国东部(弗吉尼亚北部)",
			"us-west-1":      "美国西部(加利福尼亚北部)",
			"us-west-2":      "美国西部(俄勒冈)",
			"ap-east-1":      "亚太区域(香港)",
			"ap-south-1":     "亚太区域(孟买)",
			"ap-northeast-3": "亚太区域(大阪-本地)",
			"ap-northeast-2": "亚太区域(首尔)",
			"ap-southeast-1": "亚太区域(新加坡)",
			"ap-southeast-2": "亚太区域(悉尼)",
			"ap-southeast-3": "亚太地区(雅加达)",
			"ap-northeast-1": "亚太区域(东京)",
			"ca-central-1":   "加拿大(中部)",
			"cn-north-1":     "中国(北京)",
			"cn-northwest-1": "中国(宁夏)",
			"eu-central-1":   "欧洲(法兰克福)",
			"eu-west-1":      "欧洲(爱尔兰)",
			"eu-west-2":      "欧洲(伦敦)",
			"eu-south-1":     "欧洲(米兰)",
			"eu-west-3":      "欧洲(巴黎)",
			"eu-north-1":     "欧洲(斯德哥尔摩)",
			"me-south-1":     "中东(巴林)",
			"sa-east-1":      "南美洲(圣保罗)",
			"us-gov-west-1":  "AWS GovCloud(美国西部)",
			"us-gov-east-1":  "AWS GovCloud(美国东部)",

			"af-south-1": "非洲(开普敦)",
		},
		"cloud_server_image": {
			"failed":    "unavailable",
			"error":     "unavailable",
			"disabled":  "unavailable",
			"available": "available",
		},
		"type_image": {
			"true":  "public",
			"false": "private",
		},
	}
	statusTencent = map[string]map[string]string{
		"cloud_zone": {
			"AVAILABLE":   "available",
			"UNAVAILABLE": "unavailable",
		},
		"cloud_server_type": {
			"S": "通用型",
			"M": "内存型",
			"C": "计算型",
		},
		"cloud_server_image": {
			"CREATING":     "creating",
			"NORMAL":       "available",
			"CREATEFAILED": "failed",
		},
		"type_image": {
			"PUBLIC_IMAGE":  "public",
			"PRIVATE_IMAGE": "private",
			"SHARED_IMAGE":  "shared",
		},
		"cloud_server": {
			"PENDING":       "pending",
			"LAUNCH_FAILED": "error",
			"RUNNING":       "running",
			"STOPPED":       "stopped",
			"STOPPING":      "stopping",
			"SHUTDOWN":      "stopped",
		},
		"pay_type": {
			"0":                "prepaid",
			"1":                "postpaid",
			"PREPAID":          "prepaid",
			"POSTPAID_BY_HOUR": "postpaid",
		},
		"renew_type": {
			"NOTIFY_AND_MANUAL_RENEW":         "manual",
			"NOTIFY_AND_AUTO_RENEW":           "auto",
			"DISABLE_NOTIFY_AND_MANUAL_RENEW": "manual",
		},
	}
	statusHuawei = map[string]map[string]string{
		"cloud_server_type": {
			"normal":      "通用型",
			"computingv3": "计算型",
			"highmem":     "内存型",
		},
		"cloud_server_image": {
			"deleted": "deleted",
			"active":  "available",
		},
		"type_image": {
			"gold":    "public",
			"private": "private",
			"shared":  "shared",
		},
		"cloud_server": {
			"REBOOT":  "rebooting",
			"ACTIVE":  "running",
			"SHUTOFF": "stopped",
			"ERROR":   "error",
			"DELETED": "deleted",
		},
		"pay_type": {
			"1": "prepaid",
			"0": "postpaid",
		},
		"renew_type": {
			"0": "no",
			"1": "no",
			"2": "no",
			"3": "auto",
			"4": "no",
			"5": "no",
		},
		"cloud_subnet": {
			"ACTIVE":  "available",
			"UNKNOWN": "unavailable",
			"ERROR":   "unavailable",
		},
	}
	statusBaidu = map[string]map[string]string{
		"cloud_server_image": {
			"Available":    "available",
			"NotAvailable": "deleted",
			"Error":        "deleted",
		},
		"type_image": {
			"System":      "public",
			"Integration": "public",
			"Custom":      "private",
		},
		"cloud_server": {
			"Running":  "running",
			"Pending":  "pending",
			"Starting": "starting",
			"Stopping": "stopping",
			"Stopped":  "stopped",
		},
		"pay_type": {
			"Prepaid":  "prepaid",
			"Postpaid": "postpaid",
		},
		"renew_type": {
			"false": "manual",
			"true":  "auto",
		},
	}
)
