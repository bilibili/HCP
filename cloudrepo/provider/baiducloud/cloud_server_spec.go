// Package baiducloud
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
package baiducloud

import (
	"strconv"

	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/bilibili/HCP/cloudrepo"
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

// GetPPS 获取云服务器规格PPS
func (c *CloudServerSpec) GetPPS() float64 {
	return c.PPS
}

// ListServerSpec 获取服务器规格列表
func (c *Client) ListServerSpec(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServerSpec, error) {
	cli, err := c.clientBcc()
	if err != nil {
		return nil, err
	}
	request := &api.ListFlavorSpecArgs{}
	listVPCResult, err := cli.ListFlavorSpec(request)
	if err != nil {
		return nil, err
	}
	return list2DoServerSpec(listVPCResult)
}

func list2DoServerSpec(resp *api.ListFlavorSpecResult) (list []cloudrepo.CloudServerSpec, err error) {
	for _, v := range resp.ZoneResources {
		_ = v
		for _, flavor := range v.BccResources.FlavorGroups {
			_ = flavor
			var group = ""
			switch flavor.GroupId {
			case "calculate":
				group = "计算型"
			case "common":
				group = "通用型"
			case "memory":
				group = "内存型"
			}
			for _, f := range flavor.Flavors {
				_ = f
				name := strconv.Itoa(f.CpuCount) + "C." + strconv.Itoa(f.MemoryCapacityInGB) + "GB " + f.SpecId + group
				bandwidth, _ := strconv.ParseFloat(f.NetworkBandwidth, 64)
				pps, _ := strconv.ParseFloat(f.NetworkPackage, 64)
				list = append(list, &CloudServerSpec{
					CloudProductCommon: cloudrepo.CloudProductCommon{
						CID:    f.Spec,
						Name:   name,
						Status: "available",
					},
					BandWidth: bandwidth,
					Category:  group,
					CPU:       f.CpuCount,
					GPU:       f.GpuCardCount,
					GPUModel:  f.GpuCardType,
					Family:    group,
					Memory:    f.MemoryCapacityInGB * 1024,
					PPS:       pps,
				})
			}
		}
	}
	return
}
