// Package tencentcloud
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
package tencentcloud

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/bilibili/HCP/cloudrepo"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

// CloudServerSpec 云服务器规格
type CloudServerSpec struct {
	cloudrepo.CloudProductCommon

	BandWidth float64
	Category  string
	CPU       int
	GPU       int
	GPUModel  string
	Family    string
	Memory    int
	PPS       float64
}

// GetCID 获取云服务器规格CID
func (c *CloudServerSpec) GetCID() string {
	return c.CID
}

// GetName 获取云服务器规格名称
func (c *CloudServerSpec) GetName() string {
	return c.Name
}

// GetStatus 获取云服务器规格状态
func (c *CloudServerSpec) GetStatus() string {
	return c.Status
}

// GetBandWidth 获取云服务器规格带宽
func (c *CloudServerSpec) GetBandWidth() float64 {
	return c.BandWidth
}

// GetCategory 获取云服务器规格类别
func (c *CloudServerSpec) GetCategory() string {
	return c.Category
}

// GetCPU 获取云服务器规格CPU
func (c *CloudServerSpec) GetCPU() int {
	return c.CPU
}

// GetGPU 获取云服务器规格GPU
func (c *CloudServerSpec) GetGPU() int {
	return c.GPU
}

// GetGPUModel 获取云服务器规格GPU型号
func (c *CloudServerSpec) GetGPUModel() string {
	return c.GPUModel
}

// GetFamily 获取云服务器规格系列
func (c *CloudServerSpec) GetFamily() string {
	return c.Family
}

// GetMemory 获取云服务器规格内存
func (c *CloudServerSpec) GetMemory() int {
	return c.Memory
}

// GetPPS 获取云服务器PPS
func (c *CloudServerSpec) GetPPS() float64 {
	return c.PPS
}

// ListServerSpec 获取服务器规格列表
func (c *Client) ListServerSpec(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServerSpec, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	request := cvm.NewDescribeZoneInstanceConfigInfosRequest()
	resp, err := cli.DescribeZoneInstanceConfigInfos(request)
	if err != nil {
		return nil, err
	}
	response := &struct {
		Response struct {
			// 可用区机型配置列表。
			InstanceTypeQuotaSet []InstanceTypeQuotaItem `json:"InstanceTypeQuotaSet,omitempty" name:"InstanceTypeQuotaSet"`
			// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
			RequestId string `json:"RequestId,omitempty" name:"RequestId"`
		} `json:"Response"`
	}{}
	err = json.Unmarshal([]byte(resp.ToJsonString()), response)
	if err != nil {
		return nil, err
	}
	return list2DoServerSpec(response.Response.InstanceTypeQuotaSet)
}

func list2DoServerSpec(resp []InstanceTypeQuotaItem) (list []cloudrepo.CloudServerSpec, err error) {
	for _, v := range resp {
		_ = v
		name := strconv.Itoa(int(v.Cpu)) + "C." + strconv.Itoa(int(v.Memory)) + "GB " + v.InstanceFamily + "规格族"
		category := cloudrepo.GetCloudEnum(cloudrepo.CloudTencent, "cloud_server_type", string(([]byte(v.InstanceFamily))[0]))
		gpuMode := ""
		if v.Gpu > 0 {
			gpuMode = cloudrepo.GetGpuModel(v.Remark)
			if gpuMode == "" {
				gpuMode = getQcloudGpuMode(v.Remark)
			}
		}
		list = append(list, &CloudServerSpec{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.InstanceType,
				Name:   name,
				Status: "available",
			},
			BandWidth: v.InstanceBandwidth,
			Category:  category,
			CPU:       int(v.Cpu),
			GPU:       int(v.Gpu),
			GPUModel:  gpuMode,
			Memory:    int(v.Memory) * 1024,
			Family:    v.InstanceFamily,
			PPS:       float64(v.InstancePps),
		})
	}
	return
}

func getQcloudGpuMode(raw string) string {
	raws := strings.Split(raw, " ")
	if len(raws) <= 2 {
		return raw
	}
	res := ""
	for _, s := range raws[2:] {
		if res == "" {
			res = s
			continue
		}
		res = res + " " + s
	}
	return res
}
