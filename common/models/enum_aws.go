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

package common

type AwsDBInstanceSpec struct {
	VcpuCount  int
	VmemSizeMb int
}

var (
	statusAws = map[string]map[string]string{
		CloudServerImage: {
			"pending":      "pending",
			"available":    "available",
			"invalid":      "invalid",
			"deregistered": "deregistered",
			"transient":    "transient",
			"failed":       "failed",
			"error":        "error",
		},
		CloudRegionLocation: {
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
		CloudServer: {
			"running":       CloudDBStatusRunning,
			"pending":       CloudDBStatusPending,
			"terminated":    CloudDBStatusStopped,
			"stopping":      CloudDBStatusStopping,
			"stopped":       CloudDBStatusStopped,
			"shutting-down": CloudDBStatusStopped,
		},
		CloudSubnet: {
			"available": "available",
			"pending":   "unavailable",
		},
		CloudVpc: {
			"available": "available",
			"pending":   "unavailable",
		},
		CloudZone: {
			"available":   "available",
			"information": "information",
			"impaired":    "impaired",
			"unavailable": "unavailable",
		},
		PayType: {
			"": CloudPayTypeByFlow,
		},
	}
)
