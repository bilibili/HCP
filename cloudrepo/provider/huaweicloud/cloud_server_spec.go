// Package huaweicloud
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
package huaweicloud

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bilibili/HCP/cloudrepo"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	"github.com/pkg/errors"
)

var (
	specEndpointMap = map[string]string{
		"cn-south-2": "https://ecs.cn-south-2.myhuaweicloud.com",
	}

	HwSeverSpecs = map[string]string{
		"normal":            "通用计算型",
		"entry":             "通用入门型",
		"cpuv1":             "计算I型",
		"cpuv2":             "计算II型",
		"computingv3":       "通用计算增强型",
		"kunpeng_computing": "鲲鹏通用计算增强型",
		"kunpeng_highmem":   "鲲鹏内存优化型",
		"highmem":           "内存优化型",
		"saphana":           "大内存型",
		"diskintensive":     "磁盘增强型",
		"highio":            "超高I/O型",
		"ultracpu":          "超高性能计算型",
		"gpu":               "GPU加速型",
		"fpga":              "FPGA加速型",
		"ascend":            "AI加速型",
		"arm64":             "鲲鹏计算",
	}
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
	// 部分Region实例化client会报错，这里进行异常捕获
	defer func() {
		e := recover()
		if e != nil {
			err := errors.New(fmt.Sprintf("%v", e))
			if err != nil {
				return
			}
		}
	}()
	var cli *ecs.EcsClient
	var data []model.Flavor
	ep, ok := specEndpointMap[req.Region]
	if ok {
		projects, err := c.ListProjectIam(req)
		if err != nil {
			return nil, err
		}
		for _, p := range projects {
			if p.Name != req.Region {
				continue
			}
			cli = c.clientEcsWithEndPoint(p.Id, ep)
			servers, err := c.queryServerSpec(req, cli)
			if err != nil {
				return nil, err
			}
			for idx, _ := range servers {
				data = append(data, servers[idx])
			}
		}
		return list2DoServerSpec(data)
	} else {
		cli = c.clientEcs()
		resp, err := c.queryServerSpec(req, cli)
		if err != nil {
			return nil, err
		}
		return list2DoServerSpec(resp)
	}
}

func list2DoServerSpec(resp []model.Flavor) (list []cloudrepo.CloudServerSpec, err error) {
	for _, v := range resp {
		_ = v
		category := ""
		bandwidth, pps := 0.0, 0.0
		name := v.Vcpus + "C." + strconv.Itoa(int(v.Ram/1024)) + "GB"
		if v.OsExtraSpecs != nil {
			if v.OsExtraSpecs.Ecsperformancetype != nil {
				name += " " + HwSeverSpecs[*v.OsExtraSpecs.Ecsperformancetype]
				category = cloudrepo.GetCloudEnum(cloudrepo.CloudHuawei, "cloud_server_type", *v.OsExtraSpecs.Ecsperformancetype)
			}
			if v.OsExtraSpecs.EcsinstanceArchitecture != nil {
				name += " （" + HwSeverSpecs[*v.OsExtraSpecs.EcsinstanceArchitecture] + "）"
			}
			if v.OsExtraSpecs.QuotamaxRate != nil {
				if maxRate, err := strconv.Atoi(*v.OsExtraSpecs.QuotamaxRate); err == nil {
					bandwidth = float64(maxRate) / float64(1000)
				}
			}
			if v.OsExtraSpecs.QuotamaxPps != nil {
				if maxPps, err := strconv.Atoi(*v.OsExtraSpecs.QuotamaxPps); err == nil {
					pps = float64(maxPps) / float64(10000)
				}
			}
		}
		gpu := 0
		gpuModel := ""
		gpuStr := ""
		if v.OsExtraSpecs != nil && v.OsExtraSpecs.Infogpuname != nil {
			gpuStr = *v.OsExtraSpecs.Infogpuname
		}
		tup := strings.Split(gpuStr, "/")
		if len(tup) > 1 {
			tap := strings.Split(tup[0], "*")
			if len(tap) > 1 {
				gpu, _ = strconv.Atoi(strings.TrimSpace(tap[0]))
				gpuModel = strings.TrimSpace(tap[1])
			}
		}
		cpu, _ := strconv.Atoi(v.Vcpus)
		typeTmp := strings.Split(v.Id, ".")
		family := typeTmp[0]
		list = append(list, &CloudServerSpec{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.Id,
				Name:   name,
				Status: "available",
			},
			BandWidth: bandwidth,
			Category:  category,
			CPU:       cpu,
			GPU:       gpu,
			GPUModel:  gpuModel,
			Family:    family,
			Memory:    int(v.Ram),
			PPS:       pps,
		})
	}
	return
}

func (c *Client) queryServerSpec(req *cloudrepo.GetCloudProductReq, cli *ecs.EcsClient) (res []model.Flavor, err error) {
	req.DisablePage = true
	request := model.ListFlavorsRequest{}
	resp, err := cli.ListFlavors(&request)
	if err != nil {
		return nil, err
	}
	return *resp.Flavors, nil
}
