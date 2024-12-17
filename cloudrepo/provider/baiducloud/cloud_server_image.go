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
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/bilibili/HCP/cloudrepo"
)

// CloudServerImage 云服务器镜像
type CloudServerImage struct {
	cloudrepo.CloudProductCommon

	OsName string
	Type   string
}

// GetCID 获取云镜像ID
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
	cli, err := c.clientBcc()
	if err != nil {
		return nil, err
	}
	maxItems := 1000
	var marker = new(string)
	*marker = ""
	args := &api.ListImageArgs{
		MaxKeys: maxItems,
	}
	var data []api.ImageModel
	for {
		args.Marker = *marker
		resp, err := cli.ListImage(args)
		if err != nil {
			return nil, err
		}
		data = append(data, resp.Images...)
		marker = &resp.NextMarker
		if !resp.IsTruncated {
			req.DisablePage = true
			break
		}
	}
	return list2DoServerImage(&api.ListImageResult{
		Images: data,
	})
}

func list2DoServerImage(resp *api.ListImageResult) (list []cloudrepo.CloudServerImage, err error) {
	for _, v := range resp.Images {
		_ = v
		list = append(list, &CloudServerImage{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.Id,
				Name:   fmt.Sprintf("%v%v", v.OsName, v.Name),
				Status: cloudrepo.GetCloudEnum(cloudrepo.CloudBaidu, "cloud_server_image", string(v.Status)),
			},
			OsName: v.OsName,
			Type:   cloudrepo.GetCloudEnum(cloudrepo.CloudBaidu, "type_image", fmt.Sprintf("%v", v.Type)),
		})
	}
	return
}
