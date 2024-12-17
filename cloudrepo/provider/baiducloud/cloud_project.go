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
	"encoding/json"
	"io/ioutil"
	"net/http"

	bceHttp "github.com/baidubce/bce-sdk-go/http"
	"github.com/baidubce/bce-sdk-go/services/bcc"
	"github.com/baidubce/bce-sdk-go/util"
	"github.com/bilibili/HCP/cloudrepo"
)

// CloudProject 云项目
type CloudProject struct {
	cloudrepo.CloudProductCommon
	//client    *Client
	ProjectID string

	Groups []struct {
		ParentId string        `json:"parentId"`
		GroupId  string        `json:"groupId"`
		Name     string        `json:"name"`
		Extra    string        `json:"extra"`
		Children []interface{} `json:"children"`
	} `json:"groups"`
}

// GetCID 获取云项目ID
func (c *CloudProject) GetCID() string {
	return c.CID
}

// GetName 获取云项目名称
func (c *CloudProject) GetName() string {
	return c.Name
}

// GetStatus 获取云项目状态
func (c *CloudProject) GetStatus() string {
	return c.Status
}

// ListProject 获取项目列表
func (c *Client) ListProject(_ *cloudrepo.GetCloudProductReq) ([]cloudrepo.CloudProject, error) {
	cli, err := c.clientBcc()
	if err != nil {
		return nil, err
	}
	request := &bceHttp.Request{}
	request.SetEndpoint("http://resourcemanager.baidubce.com")
	request.SetUri("/v1/res/group")
	respBody, err := requestWithSigner(request, cli)
	if err != nil {
		return nil, err
	}
	var response = &CloudProject{}
	if err = json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}
	return list2Do(response)
}

func list2Do(resp *CloudProject) (list []cloudrepo.CloudProject, err error) {
	for _, v := range resp.Groups {
		_ = v
		list = append(list, &CloudProject{
			CloudProductCommon: cloudrepo.CloudProductCommon{
				CID:  v.GroupId,
				Name: v.Name,
			},
			//ProjectID: v.GroupId,
		})
	}
	return
}

// requestWithSigner 构造请求
func requestWithSigner(request *bceHttp.Request, client *bcc.Client) ([]byte, error) {
	request.SetMethod(bceHttp.GET)
	request.SetHeader(bceHttp.HOST, request.Host())
	request.SetHeader(bceHttp.USER_AGENT, client.Config.UserAgent)
	request.SetHeader(bceHttp.BCE_DATE, util.FormatISO8601Date(util.NowUTCSeconds()))
	if request.Header(bceHttp.CONTENT_TYPE) == "" {
		request.SetHeader(bceHttp.CONTENT_TYPE, "application/json;charset=utf-8")
	}
	client.Signer.Sign(request, client.Config.Credentials, client.Config.SignOption)
	newRequest, err := http.NewRequest("GET", "http://resourcemanager.baidubce.com/v1/res/group", nil)
	if err != nil {
		return nil, err
	}
	for k, v := range request.Headers() {
		newRequest.Header.Set(k, v)
	}
	h := http.Client{}
	resp, err := h.Do(newRequest)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
