// Package sync
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
package sync

import (
	"context"
	"time"

	"github.com/bilibili/HCP/app/interface/v1/configs"
	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	common "github.com/bilibili/HCP/common/models"
	"github.com/bilibili/HCP/utils/aes"
	"github.com/go-kratos/kratos/pkg/log"
)

// SyncCloudResource 同步云资源数据
func (c *Client) SyncCloudResource(ctx context.Context, req *common.SyncCloudResourceRequest) error {
	requests, err := c.SyncRequestToQueryRequest(ctx, req)
	if err != nil {
		return err
	}
	syncFunc, ok := c.ResourceSyncFunc[req.ResourceType]
	if !ok {
		return err
	}
	for idx, _ := range requests {
		err = syncFunc(ctx, &requests[idx])
		if err != nil {
			return err
		}
	}
	return nil
}

// CloudAccountDecrypt 解密云账户SK信息
func (c *Client) CloudAccountDecrypt(data []*biz.Account) ([]*biz.Account, error) {
	for idx, _ := range data {
		secretIdEnc := data[idx].SyncSecretID
		secretKeyEnc, err := aes.Decrypt(configs.Conf.CloudSecret.SecretAesKey, data[idx].SyncSecretKey)
		if err != nil {
			return nil, err
		}
		data[idx].SyncSecretID = secretIdEnc
		data[idx].SyncSecretKey = secretKeyEnc
	}
	return data, nil
}

// SyncRequestToQueryRequest 同步请求转查询请求，返回查询请求列表
func (c *Client) SyncRequestToQueryRequest(ctx context.Context, req *common.SyncCloudResourceRequest) ([]common.QueryCloudResourceRequest, error) {
	conditions := map[string]interface{}{
		"cid":  req.AccountCID,
		"name": req.AccountName,
	}
	accounts, err := c.Account.QueryAccount(&biz.AccountWhere{Conditions: conditions}, &biz.AccountOutput{})
	accounts, err = c.CloudAccountDecrypt(accounts)
	if err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, err
	}
	var requests []common.QueryCloudResourceRequest
	for _, acc := range accounts {
		requests = append(requests, common.QueryCloudResourceRequest{
			NewClientReq: common.NewCloudClientRequest{
				Region:     req.Region,
				SecretId:   acc.SyncSecretID,
				SecretKey:  acc.SyncSecretKey,
				SecretAk:   "",
				SecretSk:   "",
				PublicKey:  "",
				PrivateKey: "",
			},
			ResourceType: req.ResourceType,
			CloudID:      req.CloudID,
			AccountCID:   acc.CID,
			AccountName:  acc.Name,
			AccountAlias: acc.Alias,
			Region:       req.Region,
			ResourceID:   req.ResourceID,
			ResourceName: req.ResourceName,
		})
	}
	return requests, err
}

// BeforeSyncCloudResource 同步云资源前置操作
func (c *Client) BeforeSyncCloudResource(resourceType, cloudId string, acc *biz.Account) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	if resourceType == common.CloudRegion {
		req := &common.SyncCloudResourceRequest{
			Region:       "",
			CloudID:      cloudId,
			AccountCID:   acc.CID,
			AccountName:  acc.Name,
			AccountAlias: acc.Alias,
			ResourceType: resourceType,
			ResourceName: "",
			ResourceID:   "",
		}
		err := c.SyncCloudResource(ctx, req)
		if err != nil {
			log.Error("BeforeSyncCloudResource SyncCloudResource cloud id (%v) resource type (%v) error (%v)", cloudId, resourceType, err)
		}
		return
	}
	conditions := map[string]interface{}{
		"provider_id": acc.ProviderID,
	}
	resp, err := c.CloudRegion.QueryCloudRegion(&biz.CloudRegionWhere{Conditions: conditions}, &biz.CloudRegionOutput{})
	if err != nil {
		log.Error("BeforeSyncCloudResource QueryCloudRegion cloud id (%v) resource type (%v) account name (%v) error (%v)", cloudId, resourceType, acc.Name, err)
		return
	}
	for _, dt := range resp {
		ctx, _ = context.WithTimeout(context.Background(), time.Minute*10)
		req := &common.SyncCloudResourceRequest{
			Region:       dt.CID,
			CloudID:      cloudId,
			AccountCID:   acc.CID,
			AccountName:  acc.Name,
			AccountAlias: acc.Alias,
			ResourceType: resourceType,
			ResourceName: "",
			ResourceID:   "",
		}
		err = c.SyncCloudResource(ctx, req)
		if err != nil {
			log.Error("BeforeSyncCloudResource SyncCloudResource cloud id (%v) resource type (%v) region (%v) error (%v)", cloudId, resourceType, dt.CID, err)
		}
	}
}

// LoadData 启动goroutine数据同步
func (c *Client) LoadData() {
	var pds []string
	pds = append(pds, common.SyncCloudProducts...)
	pds = append(pds, common.SyncCloudProductsQuick...)
	pds = append(pds, common.SyncCloudProductsSlow...)
	go c.Sync()
	go c.SyncQuick()
	go c.SyncSlow()
}

func (c *Client) GetWorkerCh() chan syncJob {
	return c.syncCloudWorkCh
}

func (c *Client) GetWorkerQuickCh() chan syncJob {
	return c.syncCloudWorkQuickCh
}

func (c *Client) GetWorkerSlowCh() chan syncJob {
	return c.syncCloudWorkSlowCh
}

// StartSync 启动数据同步
func (c *Client) StartSync() {
	ctx := context.Background()
	_ = ctx
	if configs.Conf.CloudSync.Sync {
		go c.LoadData()
	}
}

// StartWorker 启动数据同步worker
func (c *Client) StartWorker() {
	for i := 0; i <= configs.Conf.CloudSync.ConcurrencyAccount*3; i++ {
		go func() {
			for {
				select {
				case job := <-c.GetWorkerCh():
					c.doSyncJob(job)
				}
			}
		}()
	}
}

// StartQuickWorker 启动数据同步worker
func (c *Client) StartQuickWorker() {
	for i := 0; i < configs.Conf.CloudSync.ConcurrencyAccount; i++ {
		go func() {
			for {
				select {
				case job := <-c.GetWorkerQuickCh():
					c.doSyncJob(job)
				}
			}
		}()
	}
}

// StartSlowWorker 启动数据同步worker
func (c *Client) StartSlowWorker() {
	go func() {
		for {
			select {
			case job := <-c.GetWorkerSlowCh():
				c.doSyncJob(job)
			}
		}
	}()
}

// doSyncJob 执行数据同步
func (c *Client) doSyncJob(job syncJob) {
	accounts, err := c.Account.QueryAccount(&biz.AccountWhere{
		Conditions: map[string]interface{}{
			"is_delete": 0,
		},
	}, &biz.AccountOutput{})
	if err != nil {
		return
	}
	providers, err := c.Provider.QueryProvider(&biz.ProviderWhere{
		Conditions: map[string]interface{}{
			"is_delete": 0,
		},
	}, &biz.ProviderOutput{})
	if err != nil {
		return
	}
	var providersMap = make(map[int]string)
	for _, dt := range providers {
		providersMap[dt.ID] = dt.Alias
	}
	for _, dt := range accounts {
		if cloudID, ok := providersMap[dt.ProviderID]; ok {
			c.BeforeSyncCloudResource(job.ResourceType, cloudID, &biz.Account{
				CloudProductCommon: biz.CloudProductCommon{
					ID:   dt.ID,
					CID:  dt.CID,
					Name: dt.Name,
				},
				Alias:      dt.Alias,
				ProviderID: dt.ProviderID,
			})
		}
	}
}

// Sync 向syncCloudWorkCh推送数据
func (c *Client) Sync() {
	for _, product := range common.SyncCloudProducts {
		c.GetWorkerCh() <- syncJob{
			ResourceType: product,
		}
		time.Sleep(time.Minute * time.Duration(configs.Conf.CloudSync.SyncInterval))
	}
}

// SyncQuick 向syncCloudWorkQuickCh推送数据
func (c *Client) SyncQuick() {
	for {
		for _, product := range common.SyncCloudProductsQuick {
			c.GetWorkerQuickCh() <- syncJob{
				ResourceType: product,
			}
			time.Sleep(time.Minute * time.Duration(configs.Conf.CloudSync.SyncInterval))
		}
	}
}

// SyncSlow 向syncCloudWorkSlowCh推送数据
func (c *Client) SyncSlow() {
	for {
		for _, product := range common.SyncCloudProductsSlow {
			c.GetWorkerSlowCh() <- syncJob{
				ResourceType: product,
			}
			time.Sleep(time.Minute * time.Duration(configs.Conf.CloudSync.SyncInterval))
		}
	}
}
