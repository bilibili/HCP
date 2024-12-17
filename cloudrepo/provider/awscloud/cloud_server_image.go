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
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/bilibili/HCP/cloudrepo"
)

// CloudServerImage 云服务器镜像
type CloudServerImage struct {
	cloudrepo.CloudProductCommon

	OsName string
	Type   string
}

// GetCID 获取镜像CID
func (c *CloudServerImage) GetCID() string {
	return c.CID
}

// GetName 获取镜像名称
func (c *CloudServerImage) GetName() string {
	return c.Name
}

// GetStatus 获取镜像状态
func (c *CloudServerImage) GetStatus() string {
	return c.Status
}

// GetOsName 获取镜像系统
func (c *CloudServerImage) GetOsName() string {
	return c.OsName
}

// GetType 获取镜像类型
func (c *CloudServerImage) GetType() string {
	return c.Type
}

// ListServerImage 获取镜像列表
func (c *Client) ListServerImage(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServerImage, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancelFn()
	cli := c.clientEc2()
	var data []types.Image
	pageSize := 1000
	maxResult := int32(pageSize)
	request := &ec2.DescribeImagesInput{MaxResults: &maxResult, Owners: []string{"amazon", "aws-marketplace", "self"}}
	for {
		if req.Cursor != "" {
			request.NextToken = &req.Cursor
		}
		resp, err := cli.DescribeImages(ctx, request)
		if err != nil {
			return nil, err
		}
		if resp.Images != nil && len(resp.Images) > 0 {
			data = append(data, resp.Images...)
		}
		if resp.NextToken != nil {
			req.Cursor = *resp.NextToken
		}
		if resp.NextToken == nil {
			break
		}
	}
	return list2DoServerImage(data)
}

func list2DoServerImage(data []types.Image) (list []cloudrepo.CloudServerImage, err error) {
	for _, v := range data {
		_ = v
		list = append(list, &CloudServerImage{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    *v.ImageId,
				Name:   *v.Name,
				Status: cloudrepo.GetCloudEnum(cloudrepo.CloudAws, "cloud_server_image", string(v.State)),
			},
			OsName: *v.PlatformDetails,
			Type:   cloudrepo.GetCloudEnum(cloudrepo.CloudAws, "type_image", fmt.Sprintf("%t", *v.Public)),
		})
	}
	return
}
