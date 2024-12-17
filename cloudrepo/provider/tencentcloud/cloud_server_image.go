// Package tencentcloud
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
package tencentcloud

import (
	"encoding/json"

	"github.com/bilibili/HCP/cloudrepo"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

// CloudServerImage 云镜像信息
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

// ListServerImage 获取服务器镜像列表
func (c *Client) ListServerImage(req *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudServerImage, error) {
	cli, err := c.clientEcs()
	if err != nil {
		return nil, err
	}
	var data []Image
	pageNum, pageSize := 1, 100
	request := cvm.NewDescribeImagesRequest()
	if req.ResourceID != "" {
		request.ImageIds = []*string{&req.ResourceID}
	}
	for {
		offset := uint64((pageNum - 1) * pageSize)
		limit := uint64(pageSize)
		request.Offset = &offset
		request.Limit = &limit
		resp, err := cli.DescribeImages(request)
		if err != nil {
			return nil, err
		}
		response := &struct {
			Response struct {
				// 一个关于镜像详细信息的结构体，主要包括镜像的主要状态与属性。
				ImageSet []Image `json:"ImageSet,omitempty" name:"ImageSet"`
				// 符合要求的镜像数量。
				TotalCount int64 `json:"TotalCount,omitempty" name:"TotalCount"`
				// 唯一请求 ID，每次请求都会返回。定位问题时需要提供该次请求的 RequestId。
				RequestId string `json:"RequestId,omitempty" name:"RequestId"`
			} `json:"Response"`
		}{}
		err = json.Unmarshal([]byte(resp.ToJsonString()), response)
		if err != nil {
			return nil, err
		}
		data = append(data, response.Response.ImageSet...)
		if len(data) >= int(response.Response.TotalCount) {
			break
		}
		pageNum += 1
	}
	return list2DoServerImage(data)
}

func list2DoServerImage(resp []Image) (list []cloudrepo.CloudServerImage, err error) {
	for _, v := range resp {
		_ = v
		list = append(list, &CloudServerImage{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:    v.ImageId,
				Name:   v.ImageName,
				Status: cloudrepo.GetCloudEnum(cloudrepo.CloudTencent, "cloud_server_image", v.ImageState),
			},
			OsName: v.OsName,
			Type:   cloudrepo.GetCloudEnum(cloudrepo.CloudTencent, "type_image", v.ImageType),
		})
	}
	return
}
