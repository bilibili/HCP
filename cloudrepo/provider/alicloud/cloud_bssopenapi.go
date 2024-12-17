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
	"errors"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
)

// QueryAvailableInstancesById 查询BssOpenAPI实例信息
func (c *Client) QueryAvailableInstancesById(region string, instanceIds []string, productCode string) (instances []bssopenapi.Instance, err error) {
	cli, err := c.clientBss(region)
	if err != nil {
		return nil, err
	}
	cli.SetConnectTimeout(time.Minute)
	cli.SetReadTimeout(time.Minute)
	if len(instanceIds) == 0 {
		return nil, err
	}
	for i := 0; i < len(instanceIds); {
		request := bssopenapi.CreateQueryAvailableInstancesRequest()

		if i+20 > len(instanceIds) {
			request.InstanceIDs = strings.Join(instanceIds[i:], ",")
		} else {
			request.InstanceIDs = strings.Join(instanceIds[i:i+20], ",")
		}
		i += 20
		request.ProductCode = productCode
		resp, err := cli.QueryAvailableInstances(request)
		if err != nil {
			return nil, err
		}
		if resp.Code != "Success" {
			return nil, errors.New(resp.Message)
		}
		instances = append(instances, resp.Data.InstanceList...)
	}
	return
}

// QueryInstanceInfo 查询实例信息
func (c *Client) QueryInstanceInfo(region string, instanceIds []string, productCode string) ([]bssopenapi.Instance, error) {
	if region == "cn-hongkong" {
		region = "ap-southeast-1"
	}
	if region == "ap-northeast-2" {
		region = "ap-northeast-1"
	}
	resp, err := c.QueryAvailableInstancesById(region, instanceIds, productCode)
	if err != nil {
		region = "cn-qingdao"
		resp, err := c.QueryAvailableInstancesById(region, instanceIds, productCode)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return resp, nil
}
