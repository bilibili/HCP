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
	statusHw = map[string]map[string]string{
		CloudRegion: {
			"": "available",
		},
		CloudZone: {
			"": "available",
		},
		CloudProject: {},
		CloudServerSpec: {
			"normal":      CloudServerSpecGeneral,
			"computingv3": CloudServerSpecCompute,
			"highmem":     CloudServerSpecMemory,
		},
		CloudServerImage: {
			"deleted": "deleted",
			"active":  "available",
		},
		TypeImage: {
			"gold":    "public",
			"private": "private",
			"shared":  "shared",
		},
		CloudVpc: {
			"ACTIVE":  "available",
			"PENDING": "unavailable",
			"OK":      "available",
		},
		CloudSubnet: {
			"ACTIVE":  "available",
			"UNKNOWN": "unavailable",
			"ERROR":   "unavailable",
		},
		CloudSecurityGroup: {},
		CloudServer: {
			"REBOOT":  "rebooting",
			"ACTIVE":  "running",
			"SHUTOFF": "stopped",
			"ERROR":   "error",
			"DELETED": "deleted",
		},
		CloudPayResourceRenewStatus: {
			"0": "no",
			"1": "no",
			"2": "no",
			"3": "auto",
			"4": "no",
			"5": "no",
		},
		PayType: {
			"1": CloudPayTypeByPeriod,
			"0": CloudPayTypeByFlow,
		},
	}
)
