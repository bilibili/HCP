// Package awscloud
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
package awscloud

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
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

// GetCID 获取云服务器规格ID
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
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFn()
	cli := c.clientEc2()
	var data []types.InstanceTypeInfo
	pageSize := 100
	maxResult := int32(pageSize)
	request := &ec2.DescribeInstanceTypesInput{MaxResults: &maxResult}
	for {
		if req.Cursor != "" {
			request.NextToken = &req.Cursor
		}
		resp, err := cli.DescribeInstanceTypes(ctx, request)
		if err != nil {
			return nil, err
		}
		if resp.InstanceTypes != nil && len(resp.InstanceTypes) > 0 {
			data = append(data, resp.InstanceTypes...)
		}
		if resp.NextToken != nil {
			req.Cursor = *resp.NextToken
		}
		if resp.NextToken == nil {
			break
		}
	}

	return list2DoServerSpec(data)
}

func list2DoServerSpec(resp []types.InstanceTypeInfo) (list []cloudrepo.CloudServerSpec, err error) {
	for _, v := range resp {
		_ = v
		GpuCount := 0
		GpuModel := ""
		if v.GpuInfo != nil {
			for _, g := range v.GpuInfo.Gpus {
				GpuCount += int(*g.Count)
				GpuModel += *g.Manufacturer + " " + *g.Name
			}

		}
		list = append(list, &CloudServerSpec{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    string(v.InstanceType),
				Name:   string(v.InstanceType),
				Status: "available",
			},
			BandWidth: 0,
			Category:  "",
			CPU:       int(*v.VCpuInfo.DefaultVCpus),
			GPU:       GpuCount,
			GPUModel:  GpuModel,
			Family:    "",
			Memory:    int(*v.MemoryInfo.SizeInMiB),
			PPS:       float64(0),
		})
	}
	return
}
