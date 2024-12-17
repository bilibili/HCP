// Package service
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
package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/common/ecode"
	common "github.com/bilibili/HCP/common/models"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
)

// QueryReq 查询请求参数
type QueryReq struct {
	// 表名
	Schema string `json:"schema"`
	// 条件
	Where map[string]interface{} `json:"where"`
	// 条数
	PageSize int64 `json:"page_size"`
	// 页码
	PageNum int64 `json:"page_num"`
	// 排序
	Order string `json:"order"`
}

// QueryResp 查询响应参数
type QueryResp struct {
	// 总数
	Total int64 `json:"total"`
	// 响应数据
	Data json.RawMessage `json:"data"`
}

// Query 查询数据
func Query(c *bm.Context, w http.ResponseWriter, r *http.Request) {
	req := &QueryReq{}
	if err := c.BindWith(req, binding.JSON); err != nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
		return
	}
	conditions := req.Where
	if conditions == nil {
		conditions = make(map[string]interface{})
	}
	userInfo := r.Context().Value("userInfo").(biz.UserValues)
	permissionConditions, err := QueryUserPermission(userInfo, req.Schema)
	if err != nil {
		c.JSON(&QueryResp{}, ecode.NewECode(common.ErrServiceCode, err.Error()))
		return
	}
	for k, v := range permissionConditions {
		conditions[k] = v
	}
	switch req.Schema {
	case "cloud_server":
		resp, err := Svc.CloudServer.QueryCloudServer(&biz.CloudServerWhere{Conditions: conditions}, &biz.CloudServerOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.CloudServer.CountCloudServer(&biz.CloudServerWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "provider":
		resp, err := Svc.Provider.QueryProvider(&biz.ProviderWhere{Conditions: conditions}, &biz.ProviderOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.Provider.CountProvider(&biz.ProviderWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "cloud_project":
		resp, err := Svc.CloudProject.QueryCloudProject(&biz.CloudProjectWhere{Conditions: conditions}, &biz.CloudProjectOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.CloudProject.CountCloudProject(&biz.CloudProjectWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "account":
		resp, err := Svc.Account.QueryAccount(&biz.AccountWhere{Conditions: conditions}, &biz.AccountOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.Account.CountAccount(&biz.AccountWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "cloud_region":
		resp, err := Svc.CloudRegion.QueryCloudRegion(&biz.CloudRegionWhere{Conditions: conditions}, &biz.CloudRegionOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.CloudRegion.CountCloudRegion(&biz.CloudRegionWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: int64(total),
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "cloud_security_group":
		resp, err := Svc.CloudSecurityGroup.QueryCloudSecurityGroup(&biz.CloudSecurityGroupWhere{Conditions: conditions}, &biz.CloudSecurityGroupOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.CloudSecurityGroup.CountCloudSecurityGroup(&biz.CloudSecurityGroupWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "cloud_server_image":
		resp, err := Svc.CloudServerImage.QueryCloudServerImage(&biz.CloudServerImageWhere{Conditions: conditions}, &biz.CloudServerImageOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.CloudServerImage.CountCloudServerImage(&biz.CloudServerImageWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "cloud_server_spec":
		resp, err := Svc.CloudServerSpec.QueryCloudServerSpec(&biz.CloudServerSpecWhere{Conditions: conditions}, &biz.CloudServerSpecOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.CloudServerSpec.CountCloudServerSpec(&biz.CloudServerSpecWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "cloud_subnet":
		resp, err := Svc.CloudSubnet.QueryCloudSubnet(&biz.CloudSubnetWhere{Conditions: conditions}, &biz.CloudSubnetOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.CloudSubnet.CountCloudSubnet(&biz.CloudSubnetWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "cloud_vpc":
		resp, err := Svc.CloudVpc.QueryCloudVpc(&biz.CloudVpcWhere{Conditions: conditions}, &biz.CloudVpcOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.CloudVpc.CountCloudVpc(&biz.CloudVpcWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "cloud_zone":
		resp, err := Svc.CloudZone.QueryCloudZone(&biz.CloudZoneWhere{Conditions: conditions}, &biz.CloudZoneOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.CloudZone.CountCloudZone(&biz.CloudZoneWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "project_config":
		resp, err := Svc.ProjectConfig.QueryProjectConfig(&biz.ProjectConfigWhere{Conditions: conditions}, &biz.ProjectConfigOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.ProjectConfig.CountProjectConfig(&biz.ProjectConfigWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "project_account_config":
		resp, err := Svc.ProjectAccountConfig.QueryProjectAccountConfig(&biz.ProjectAccountConfigWhere{Conditions: conditions}, &biz.ProjectAccountConfigOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.ProjectAccountConfig.CountProjectAccountConfig(&biz.ProjectAccountConfigWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "user":
		field := "id, name, role, create_time, update_time"
		resp, err := Svc.User.QueryUser(&biz.UserWhere{Conditions: conditions}, &biz.UserOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}}, field)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.User.CountUser(&biz.UserWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "job":
		resp, err := Svc.Job.QueryJob(&biz.JobWhere{Conditions: conditions}, &biz.JobOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.Job.CountJob(&biz.JobWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "charge_type":
		resp, err := Svc.ChargeType.QueryChargeType(&biz.ChargeTypeWhere{Conditions: conditions}, &biz.ChargeTypeOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.ChargeType.CountChargeType(&biz.ChargeTypeWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	case "disk_type":
		resp, err := Svc.DiskType.QueryDiskType(&biz.DiskTypeWhere{Conditions: conditions}, &biz.DiskTypeOutput{OutPutCommon: biz.OutPutCommon{
			PageSize: int(req.PageSize),
			PageNum:  int(req.PageNum),
			Order:    req.Order,
		}})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		jsonData, err := json.Marshal(resp)
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		total, err := Svc.DiskType.CountDiskType(&biz.DiskTypeWhere{Conditions: conditions})
		if err != nil {
			c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
			return
		}
		data := &QueryResp{
			Total: total,
			Data:  jsonData,
		}
		c.JSON(data, err)
		return
	default:
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, fmt.Sprintf("query operation unknown schema %s", req.Schema)))
		return
	}
}

// QueryFormTemplate 查询表单模板
func QueryFormTemplate(c *bm.Context) {
	req := &QueryReq{}
	if err := c.BindWith(req, binding.JSON); err != nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
		return
	}
	resp, err := Svc.FormTemplate.QueryFormTemplate(&biz.FormTemplateWhere{Conditions: req.Where})
	if err != nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
		return
	}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
		return
	}
	data := &QueryResp{
		Data: jsonData,
	}
	c.JSON(data, err)
	return
}
