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

var (
	statusQ = map[string]map[string]string{
		CloudRegion: {
			"available": "available",
			"soldOut":   "unavailable",
		},
		CloudZone: {
			"AVAILABLE":   "available",
			"UNAVAILABLE": "unavailable",
		},
		CloudProject: {},
		CloudServerSpec: {
			"S": CloudServerSpecGeneral,
			"M": CloudServerSpecMemory,
			"C": CloudServerSpecCompute,
		},
		CloudServerImage: {
			"CREATING":     "creating",
			"NORMAL":       "available",
			"CREATEFAILED": "failed",
		},
		TypeImage: {
			"PUBLIC_IMAGE":  "public",
			"PRIVATE_IMAGE": "private",
			"SHARED_IMAGE":  "shared",
		},
		CloudVpc: {
			"": "available",
		},
		CloudSubnet: {
			"": "available",
		},
		CloudSecurityGroup: {},
		CloudServer: {
			"PENDING":       "pending",
			"LAUNCH_FAILED": "error",
			"RUNNING":       "running",
			"STOPPED":       "stopped",
			"STOPPING":      "stopping",
			"SHUTDOWN":      "stopped",
		},
		CloudPayResourceRenewStatus: {
			"AutoRenewal":   "auto",
			"ManualRenewal": "manual",
			"NotRenewal":    "no",
		},
		PayType: {
			"0":                CloudPayTypeByPeriod,
			"1":                CloudPayTypeByFlow,
			"PREPAID":          CloudPayTypeByPeriod,
			"POSTPAID_BY_HOUR": CloudPayTypeByFlow,
		},
	}
)
