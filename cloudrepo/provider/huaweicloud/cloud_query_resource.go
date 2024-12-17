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
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/bss/v2/model"
)

// ListPayPerUseCustomerResources 查询订单资源
func (c *Client) ListPayPerUseCustomerResources(resourceIds []string) (instance []model.OrderInstanceV2, err error) {
	if len(resourceIds) == 0 {
		return nil, nil
	}
	cli := c.clientBssWithProject()
	maxSendId := 50
	limit := int32(500)
	for i := 0; i < len(resourceIds); {
		var ids []string
		if i+maxSendId > len(resourceIds) {
			ids = resourceIds[i:]
		} else {
			ids = resourceIds[i : i+maxSendId]
		}
		i += maxSendId
		request := model.ListPayPerUseCustomerResourcesRequest{Body: &model.QueryResourcesReq{ResourceIds: &ids, Limit: &limit}}
		resp, err := cli.ListPayPerUseCustomerResources(&request)
		if err != nil {
			return nil, err
		}
		instance = append(instance, *(resp.Data)...)
	}
	return
}
