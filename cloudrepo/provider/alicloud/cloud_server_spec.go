// Package alicloud
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
package alicloud

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
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

// ListServerSpec 获取云服务器规格
func (c *Client) ListServerSpec(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServerSpec, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	var data []ecs.InstanceType
	var nextToken = new(string)
	*nextToken = ""
	request := ecs.CreateDescribeInstanceTypesRequest()
	request.SetConnectTimeout(time.Second * 2)
	request.SetReadTimeout(time.Second * 3)
	for {
		request.NextToken = *nextToken
		resp, err := cli.DescribeInstanceTypes(request)
		if err != nil {
			return nil, err
		}
		data = append(data, resp.InstanceTypes.InstanceType...)
		nextToken = &resp.NextToken
		if *nextToken == "" {
			break
		}
	}
	return c.list2DoServerSpec(data)
}

func (c *Client) list2DoServerSpec(resp []ecs.InstanceType) (list []cloudrepo.CloudServerSpec, err error) {
	for _, v := range resp {
		_ = v
		bandwidth := v.InstanceBandwidthRx
		if bandwidth < v.InstanceBandwidthTx {
			bandwidth = v.InstanceBandwidthTx
		}
		instanceTypeFamilyArr := strings.Split(v.InstanceTypeFamily, ".")
		InstanceTypeFamily := ""
		if len(instanceTypeFamilyArr) > 1 {
			InstanceTypeFamily = instanceTypeFamilyArr[1]
		}
		name := strconv.Itoa(v.CpuCoreCount) + "C." + strconv.FormatFloat(v.MemorySize, 'f', -1, 32) + "GB " + InstanceTypeFamily + "规格族"
		params := strings.Split(v.InstanceTypeFamily, ".")
		if len(params) < 2 {
			continue
		}
		re := regexp.MustCompile(`([A-Za-z]{1,8})\d`)
		finds := re.FindStringSubmatch(params[1])
		category := "unknown"
		if len(finds) > 1 {
			category = cloudrepo.GetCloudEnum(cloudrepo.CloudAli, "cloud_server_type", finds[1])
		}
		GpuModel := cloudrepo.GetGpuModel(v.GPUSpec)
		if GpuModel == "" {
			GpuModel = strings.Split(v.GPUSpec, "*")[0]
			GpuModel = strings.Split(GpuModel, "/")[0]
			if len(GpuModel) < 5 {
				GpuModel = ""
			}
		}
		pps := v.InstancePpsRx
		if pps < v.InstancePpsTx {
			pps = v.InstancePpsTx
		}
		list = append(list, &CloudServerSpec{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.InstanceTypeId,
				Name:   name,
				Status: "available",
			},
			BandWidth: float64(bandwidth) / float64(1024000),
			Category:  category,
			CPU:       v.CpuCoreCount,
			GPU:       v.GPUAmount,
			GPUModel:  GpuModel,
			Family:    InstanceTypeFamily,
			Memory:    int(v.MemorySize) * 1024,
			PPS:       float64(pps) / float64(10000),
		})
	}
	return
}
