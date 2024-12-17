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
	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/app/interface/v1/internal/biz/sync"
	"github.com/bilibili/HCP/app/interface/v1/internal/biz/worker"
	"github.com/google/wire"
)

var (
	ProviderSet = wire.NewSet(New, new(*Service))
	Svc         *Service
)

// Service is the service
type Service struct {
	sync                   *sync.Client
	Account                *biz.AccountUseCase
	Provider               *biz.ProviderUseCase
	CloudProject           *biz.CloudProjectUseCase
	CloudRegion            *biz.CloudRegionUseCase
	CloudRegionAssociation *biz.CloudRegionAssociationUseCase
	CloudSecurityGroup     *biz.CloudSecurityGroupUseCase
	CloudServer            *biz.CloudServerUseCase
	CloudServerImage       *biz.CloudServerImageUseCase
	CloudServerSpec        *biz.CloudServerSpecUseCase
	CloudSubnet            *biz.CloudSubnetUseCase
	CloudVpc               *biz.CloudVpcUseCase
	CloudZone              *biz.CloudZoneUseCase
	ProjectConfig          *biz.ProjectConfigUseCase
	ProjectAccountConfig   *biz.ProjectAccountConfigUseCase
	ProjectUserConfig      *biz.ProjectUserConfigUseCase
	User                   *biz.UserUseCase
	FormTemplate           *biz.FormTemplateUseCase
	Job                    *biz.JobUseCase
	ChargeType             *biz.ChargeTypeUseCase
	DiskType               *biz.DiskTypeUseCase
	Worker                 *worker.Worker
}

// New new a service
func New(sync *sync.Client, account *biz.AccountUseCase, provider *biz.ProviderUseCase, cloudProject *biz.CloudProjectUseCase, cloudRegion *biz.CloudRegionUseCase, cloudRegionAssociation *biz.CloudRegionAssociationUseCase,
	cloudSecurityGroup *biz.CloudSecurityGroupUseCase, cloudServer *biz.CloudServerUseCase, cloudServerImage *biz.CloudServerImageUseCase, cloudServerSpec *biz.CloudServerSpecUseCase, cloudSubnet *biz.CloudSubnetUseCase,
	cloudVpc *biz.CloudVpcUseCase, cloudZone *biz.CloudZoneUseCase, formTemplate *biz.FormTemplateUseCase, job *biz.JobUseCase, projectConfig *biz.ProjectConfigUseCase, projectAccountConfig *biz.ProjectAccountConfigUseCase,
	user *biz.UserUseCase, projectUserConfig *biz.ProjectUserConfigUseCase, chargeType *biz.ChargeTypeUseCase, diskType *biz.DiskTypeUseCase, worker *worker.Worker) (s *Service, cf func(), err error) {
	s = &Service{
		sync:                   sync,
		Account:                account,
		Provider:               provider,
		CloudProject:           cloudProject,
		CloudRegion:            cloudRegion,
		CloudRegionAssociation: cloudRegionAssociation,
		CloudSecurityGroup:     cloudSecurityGroup,
		CloudServer:            cloudServer,
		CloudServerImage:       cloudServerImage,
		CloudServerSpec:        cloudServerSpec,
		CloudSubnet:            cloudSubnet,
		CloudVpc:               cloudVpc,
		CloudZone:              cloudZone,
		ProjectConfig:          projectConfig,
		ProjectAccountConfig:   projectAccountConfig,
		ProjectUserConfig:      projectUserConfig,
		User:                   user,
		FormTemplate:           formTemplate,
		Job:                    job,
		ChargeType:             chargeType,
		DiskType:               diskType,
		Worker:                 worker,
	}
	cf = s.Close
	s.Worker.Start()
	Svc = s
	return
}

// Close close the resource.
func (s *Service) Close() {
}
