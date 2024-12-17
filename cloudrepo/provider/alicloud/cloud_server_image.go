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
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/bilibili/HCP/cloudrepo"
)

// CloudServerImage 云服务器镜像
type CloudServerImage struct {
	cloudrepo.CloudProductCommon

	OsName string
	Type   string
}

// GetCID 获取云镜像CID
func (c *CloudServerImage) GetCID() string {
	return c.CID
}

// GetName 获取云镜像名称
func (c *CloudServerImage) GetName() string {
	return c.Name
}

// GetStatus 获取云镜像状态
func (c *CloudServerImage) GetStatus() string {
	return c.Status
}

// GetOsName 获取云镜像系统
func (c *CloudServerImage) GetOsName() string {
	return c.OsName
}

// GetType 获取云镜像类型
func (c *CloudServerImage) GetType() string {
	return c.Type
}

// ListServerImage 获取镜像列表
func (c *Client) ListServerImage(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServerImage, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	var data []ecs.Image
	pageNum, pageSize := 1, 100
	request := ecs.CreateDescribeImagesRequest()
	request.SetConnectTimeout(time.Second * 2)
	request.SetConnectTimeout(time.Second * 6)
	cli.SetConnectTimeout(time.Second * 2)
	cli.SetReadTimeout(time.Second * 6)
	request.RegionId = req.Region
	request.ImageId = req.ResourceID
	request.ImageName = req.ResourceName
	for {
		request.PageNumber = requests.Integer(fmt.Sprintf("%d", pageNum))
		request.PageSize = requests.Integer(fmt.Sprintf("%d", pageSize))
		resp, err := cli.DescribeImages(request)
		if err != nil {
			return nil, err
		}
		data = append(data, resp.Images.Image...)
		if len(data) >= resp.TotalCount {
			break
		}
		pageNum += 1
	}

	return list2DoServerImage(data)
}

func list2DoServerImage(resp []ecs.Image) (list []cloudrepo.CloudServerImage, err error) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudServerImage{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.ImageId,
				Name:   v.ImageName,
				Status: cloudrepo.GetCloudEnum(cloudrepo.CloudAli, "cloud_server_image", v.Status),
			},
			OsName: v.OSName,
			Type:   cloudrepo.GetCloudEnum(cloudrepo.CloudAli, "type_image", fmt.Sprintf("%v", v.IsPublic)),
		})
	}
	return
}
