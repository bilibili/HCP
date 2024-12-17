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
	"strings"

	"github.com/bilibili/HCP/cloudrepo"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2/model"
	"github.com/pkg/errors"
)

// CloudServerImage 云镜像
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
	var data []model.ImageInfo
	pageSize := 100
	cli := c.clientImage()
	request := model.ListImagesRequest{}
	if req.ResourceID != "" {
		request.Id = &req.ResourceID
	}
	for {
		if req.Cursor != "" {
			request.Marker = &req.Cursor
		}
		limit := int32(pageSize)
		request.Limit = &limit
		resp, err := cli.ListImages(&request)
		if err != nil {
			return nil, err
		}
		if resp != nil && resp.Images != nil && len(*resp.Images) > 0 {
			data = append(data, *resp.Images...)
			last := (*resp.Images)[len(*resp.Images)-1]
			req.Cursor = last.Id
			if &last.Id == nil || &last.Id == &req.Cursor {
				break
			}
		} else {
			break
		}
	}
	return list2DoServerImage(data)
}

func list2DoServerImage(resp []model.ImageInfo) (list []cloudrepo.CloudServerImage, err error) {
	for _, v := range resp {
		_ = v
		imageOs := ""
		if v.OsVersion != nil {
			imageOs = *v.OsVersion
		}
		status, _ := v.Status.MarshalJSON()
		imageStatus := string(status)
		imageStatus = strings.Replace(imageStatus, "\n", "", -1)
		imageStatus = strings.Replace(imageStatus, `"`, "", -1)
		imageTypeJ, _ := v.Imagetype.MarshalJSON()
		imageType := string(imageTypeJ)
		imageType = strings.Replace(imageType, "\n", "", -1)
		imageType = strings.Replace(imageType, `"`, "", -1)
		list = append(list, &CloudServerImage{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.Id,
				Name:   v.Name,
				Status: cloudrepo.GetCloudEnum(cloudrepo.CloudHuawei, "cloud_server_image", imageStatus),
			},
			OsName: imageOs,
			Type:   cloudrepo.GetCloudEnum(cloudrepo.CloudHuawei, "type_image", imageType),
		})
	}
	return
}
