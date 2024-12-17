// Package worker
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
package worker

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/app/interface/v1/internal/biz/enum"
	"github.com/bilibili/HCP/cloudrepo"
	common "github.com/bilibili/HCP/common/models"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/pkg/errors"
)

// ObjectTriggerExecutable 对象任务执行结构体
type ObjectTriggerExecutable struct {
	// Object 对象结构体
	*biz.Object
	// DefaultExecutable 默认可执行对象
	*DefaultExecutable

	// cancel 取消信号
	cancel chan struct{}
}

// NewObjectTriggerExecutable 创建对象任务执行器
func NewObjectTriggerExecutable(do *biz.Object) *ObjectTriggerExecutable {
	return &ObjectTriggerExecutable{Object: do, cancel: make(chan struct{})}
}

// ResourceOperate 资源参数结构体
type ResourceOperate struct {
	AccountID     int
	InstanceCID   string
	InstanceName  string
	RegionCID     string
	ProjectCID    string
	ServerSpecCID string
	ImageCID      string
}

// Terraform 模版参数结构体
type Terraform struct {
	JobID          int
	StepID         int
	ObjectID       int
	Operate        string
	ProviderID     int
	bindPublicIP   string
	ProjectID      int
	AccountID      int
	NamePrefix     string
	NameSuffix     string
	OperandCounter int
	StatePath      string
	OutputPath     string
	RdsArch        string
	RedisArch      string
	DependentKey   string
	DependentValue string
	Input          string
	SecretID       string
	SecretKey      string
	KeyPairID      string
	KeyPairName    string
}

// Execute 执行
func (oe *ObjectTriggerExecutable) Execute(ctx context.Context) error {
	log.Warnv(ctx, []log.D{
		log.KVString("log", "MARS Pipeline Story"),
		log.KVString("type", "object"),
		log.KVString("step_id", strconv.Itoa(oe.JobStepId)),
		log.KVString("object_id", strconv.Itoa(oe.Id)),
		log.KVString("stage", "executing"),
		log.KVString("at", time.Now().Format(time.RFC3339)),
	}...)
	for {
		select {
		case <-oe.cancel:
			return errors.Errorf("cancel by signal")
		default:
			switch oe.Operate {
			case "create_provider": //云厂商创建
				var provider *biz.Provider
				// 解析JSON数据到结构体切片
				err := json.Unmarshal([]byte(oe.Raw), &provider)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "create_provider json unmarshal failed")
				}
				var list []*biz.Provider
				_, err = providerUseCase.CreateProvider(append(list, provider))
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "Job:provider table create data error")
				}
				oe.Status = enum.ExecutableStateSuccess
			case "update_provider": //云厂商编辑
				var provider *biz.Provider
				// 解析JSON数据到结构体切片
				err := json.Unmarshal([]byte(oe.Raw), &provider)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "update_provider json unmarshal failed")
				}
				err = providerUseCase.UpdateProvider(&biz.ProviderWhere{
					Query: "id = ?",
					Arg:   provider.FormObject.ID,
				}, provider)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return err
				}
				oe.Status = enum.ExecutableStateSuccess
			case "create_account": //云账号创建
				var account *biz.Account
				// 解析JSON数据到结构体切片
				err := json.Unmarshal([]byte(oe.Raw), &account)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "create_account json unmarshal failed")
				}
				var list []*biz.Account
				account.ProviderID = account.Provider.ID
				_, err = accountUseCase.CreateAccount(append(list, account))
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "Job:account table create data error")
				}
				oe.Status = enum.ExecutableStateSuccess
			case "update_account": //云账号编辑
				var account *biz.Account
				// 解析JSON数据到结构体切片
				err := json.Unmarshal([]byte(oe.Raw), &account)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "update_account json unmarshal failed")
				}
				account.ProviderID = account.Provider.ID
				err = accountUseCase.UpdateAccount(&biz.AccountWhere{
					Query: "id = ?",
					Arg:   account.FormObject.ID,
				}, account)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return err
				}
				oe.Status = enum.ExecutableStateSuccess
			case "delete_account": //云账号删除
				var (
					account *biz.AccountDelete
					ids     []int
				)
				// 解析JSON数据到结构体切片
				err := json.Unmarshal([]byte(oe.Raw), &account)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "delete_account json unmarshal failed")
				}
				for _, v := range account.FormObject {
					ids = append(ids, v.ID)
				}
				if len(ids) > 0 {
					err = accountUseCase.DeleteAccount(0, ids)
					if err != nil {
						oe.Status = enum.ExecutableStateFailure
						return err
					}
				}
				oe.Status = enum.ExecutableStateSuccess
			case "create_user": //用户创建
				var user *biz.User
				// 解析JSON数据到结构体切片
				err := json.Unmarshal([]byte(oe.Raw), &user)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "create_user json unmarshal failed")
				}
				var list []*biz.User
				_, err = userUseCase.CreateUser(append(list, user))
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "Job:user table create data error")
				}
				oe.Status = enum.ExecutableStateSuccess
			case "update_user": //用户编辑
				var user *biz.User
				// 解析JSON数据到结构体切片
				err := json.Unmarshal([]byte(oe.Raw), &user)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "update_user json unmarshal failed")
				}
				err = userUseCase.UpdateUser(&biz.UserWhere{
					Query: "id = ?",
					Arg:   user.FormObject.ID,
				}, user)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return err
				}
				oe.Status = enum.ExecutableStateSuccess
			case "delete_user": //用户删除
				var (
					user *biz.UserDelete
					ids  []int
				)
				// 解析JSON数据到结构体切片
				err := json.Unmarshal([]byte(oe.Raw), &user)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "delete_user json unmarshal failed")
				}
				for _, v := range user.FormObject {
					ids = append(ids, v.ID)
				}
				if len(ids) > 0 {
					err = userUseCase.DeleteUser(0, ids)
					if err != nil {
						oe.Status = enum.ExecutableStateFailure
						return err
					}
				}
				oe.Status = enum.ExecutableStateSuccess
			case "create_project_config": //本地项目创建
				var projectConfig *biz.ProjectConfig
				// 解析JSON数据到结构体切片
				err := json.Unmarshal([]byte(oe.Raw), &projectConfig)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "create_project_config json unmarshal failed")
				}
				var list []*biz.ProjectConfig
				projectConfigList, err := projectConfigUseCase.CreateProjectConfig(append(list, projectConfig))
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "Job:project_config table create data error")
				}
				if len(projectConfigList) > 0 {
					projectConfigID := projectConfigList[0].ID
					//创建研发负责人和项目的关联关系
					if rdLeader, isRdlMap := projectConfig.RdLeader.([]interface{}); isRdlMap {
						if len(rdLeader) > 0 {
							var rdLeaderCreate []*biz.ProjectUserConfig
							for _, v := range rdLeader {
								if rdLeaders, ok := v.(map[string]interface{}); ok {
									rdLeaderCreate = append(rdLeaderCreate, &biz.ProjectUserConfig{
										UserID:          int(rdLeaders["id"].(float64)),
										ProjectConfigID: projectConfigID,
										Role:            1,
									})
								}
							}
							if len(rdLeaderCreate) > 0 {
								_, err = projectUserConfigUseCase.CreateProjectUserConfig(rdLeaderCreate)
								if err != nil {
									oe.Status = enum.ExecutableStateFailure
									return errors.Wrap(err, "project_user_config table create rd_leader data error")
								}
							}
						}
					}

					//创建研发人员和项目的关联关系
					if rdMember, isRdmMap := projectConfig.RdMember.([]interface{}); isRdmMap {
						if len(rdMember) > 0 {
							var rdMemberCreate []*biz.ProjectUserConfig
							for _, v := range rdMember {
								if rdMembers, ok := v.(map[string]interface{}); ok {
									rdMemberCreate = append(rdMemberCreate, &biz.ProjectUserConfig{
										UserID:          int(rdMembers["id"].(float64)),
										ProjectConfigID: projectConfigID,
										Role:            2,
									})
								}
							}
							if len(rdMemberCreate) > 0 {
								_, err = projectUserConfigUseCase.CreateProjectUserConfig(rdMemberCreate)
								if err != nil {
									oe.Status = enum.ExecutableStateFailure
									return errors.Wrap(err, "project_user_config table create rd_member data error")
								}
							}
						}
					}

					//创建运维负责人和项目的关联关系
					if opLeader, isOplMap := projectConfig.OpLeader.([]interface{}); isOplMap {
						if len(opLeader) > 0 {
							var opLeaderCreate []*biz.ProjectUserConfig
							for _, v := range opLeader {
								if opLeaders, ok := v.(map[string]interface{}); ok {
									opLeaderCreate = append(opLeaderCreate, &biz.ProjectUserConfig{
										UserID:          int(opLeaders["id"].(float64)),
										ProjectConfigID: projectConfigID,
										Role:            3,
									})
								}
							}
							if len(opLeaderCreate) > 0 {
								_, err = projectUserConfigUseCase.CreateProjectUserConfig(opLeaderCreate)
								if err != nil {
									oe.Status = enum.ExecutableStateFailure
									return errors.Wrap(err, "project_user_config table create op_leader data error")
								}
							}
						}
					}

					//创建运维人员和项目的关联关系
					if opMember, isOpmMap := projectConfig.OpMember.([]interface{}); isOpmMap {
						if len(opMember) > 0 {
							var opMemberCreate []*biz.ProjectUserConfig
							for _, v := range opMember {
								if opMembers, ok := v.(map[string]interface{}); ok {
									opMemberCreate = append(opMemberCreate, &biz.ProjectUserConfig{
										UserID:          int(opMembers["id"].(float64)),
										ProjectConfigID: projectConfigID,
										Role:            4,
									})
								}
							}
							if len(opMemberCreate) > 0 {
								_, err = projectUserConfigUseCase.CreateProjectUserConfig(opMemberCreate)
								if err != nil {
									oe.Status = enum.ExecutableStateFailure
									return errors.Wrap(err, "project_user_config table create op_member data error")
								}
							}
						}
					}

				}
				oe.Status = enum.ExecutableStateSuccess
			case "update_project_config": //本地项目编辑
				var projectConfig *biz.ProjectConfig
				// 解析JSON数据到结构体切片
				err := json.Unmarshal([]byte(oe.Raw), &projectConfig)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "project_config json unmarshal failed")
				}
				err = projectConfigUseCase.UpdateProjectConfig(&biz.ProjectConfigWhere{
					Query: "id = ?",
					Arg:   projectConfig.FormObject.ID,
				}, projectConfig)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "project_config update data failed")
				}
				//查询数据
				var (
					rdLeaderExistIds []int
					rdLeaderIds      []int

					rdMemberExistIds []int
					rdMemberIds      []int

					opLeaderExistIds []int
					opLeaderIds      []int

					opMemberExistIds []int
					opMemberIds      []int
				)
				conditions := map[string]interface{}{
					"id": projectConfig.FormObject.ID,
				}
				projectConfigExists, err := projectConfigUseCase.QueryProjectConfig(&biz.ProjectConfigWhere{Conditions: conditions}, &biz.ProjectConfigOutput{})
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "project_config query data failed")
				}
				if projectConfigExists[0].RdLeader != nil {
					for _, v := range projectConfigExists[0].RdLeader.([]*biz.User) {
						rdLeaderExistIds = append(rdLeaderExistIds, v.ID)
					}
				}
				if projectConfigExists[0].RdMember != nil {
					for _, v := range projectConfigExists[0].RdMember.([]*biz.User) {
						rdMemberExistIds = append(rdMemberExistIds, v.ID)
					}
				}
				if projectConfigExists[0].OpLeader != nil {
					for _, v := range projectConfigExists[0].OpLeader.([]*biz.User) {
						opLeaderExistIds = append(opLeaderExistIds, v.ID)
					}
				}
				if projectConfigExists[0].OpMember != nil {
					for _, v := range projectConfigExists[0].OpMember.([]*biz.User) {
						opMemberExistIds = append(opMemberExistIds, v.ID)
					}
				}

				if projectConfig.FormObject.ID != 0 {
					projectConfigID := projectConfig.FormObject.ID
					//创建研发负责人和项目的关联关系
					if rdLeader, isRdlMap := projectConfig.RdLeader.([]interface{}); isRdlMap {
						if len(rdLeader) > 0 {
							var rdLeaderCreate []*biz.ProjectUserConfig
							for _, v := range rdLeader {
								if rdLeaders, ok := v.(map[string]interface{}); ok {
									rdLeaderIds = append(rdLeaderIds, int(rdLeaders["id"].(float64)))
									rdLeaderCreate = append(rdLeaderCreate, &biz.ProjectUserConfig{
										UserID:          int(rdLeaders["id"].(float64)),
										ProjectConfigID: projectConfigID,
										Role:            1,
									})
								}
							}
							if !biz.SlicesEqual(rdLeaderExistIds, rdLeaderIds) { //判断研发负责人是否变更
								where := map[string]interface{}{
									"project_config_id": projectConfigID,
									"role":              1,
								}
								err = projectUserConfigUseCase.DeleteProjectUserConfigByWhere(&biz.ProjectUserConfigWhere{
									Conditions: where,
								})
								if err != nil {
									oe.Status = enum.ExecutableStateFailure
									return errors.Wrap(err, "project_user_config delete data failed")
								}
								if len(rdLeaderCreate) > 0 {
									_, err = projectUserConfigUseCase.CreateProjectUserConfig(rdLeaderCreate)
									if err != nil {
										oe.Status = enum.ExecutableStateFailure
										return errors.Wrap(err, "project_user_config table create rd_leader data error")
									}
								}
							}
						}
					}

					//创建研发人员和项目的关联关系
					if rdMember, isRdmMap := projectConfig.RdMember.([]interface{}); isRdmMap {
						if len(rdMember) > 0 {
							var rdMemberCreate []*biz.ProjectUserConfig
							for _, v := range rdMember {
								if rdMembers, ok := v.(map[string]interface{}); ok {
									rdMemberIds = append(rdMemberIds, int(rdMembers["id"].(float64)))
									rdMemberCreate = append(rdMemberCreate, &biz.ProjectUserConfig{
										UserID:          int(rdMembers["id"].(float64)),
										ProjectConfigID: projectConfigID,
										Role:            2,
									})
								}
							}
							if !biz.SlicesEqual(rdMemberExistIds, rdMemberIds) { //判断研发人员是否变更
								where := map[string]interface{}{
									"project_config_id": projectConfigID,
									"role":              2,
								}
								err = projectUserConfigUseCase.DeleteProjectUserConfigByWhere(&biz.ProjectUserConfigWhere{
									Conditions: where,
								})
								if err != nil {
									oe.Status = enum.ExecutableStateFailure
									return errors.Wrap(err, "project_user_config delete data failed")
								}
								if len(rdMemberCreate) > 0 {
									_, err = projectUserConfigUseCase.CreateProjectUserConfig(rdMemberCreate)
									if err != nil {
										oe.Status = enum.ExecutableStateFailure
										return errors.Wrap(err, "project_user_config table create rd_member data error")
									}
								}
							}
						}
					}

					//创建运维负责人和项目的关联关系
					if opLeader, isOplMap := projectConfig.OpLeader.([]interface{}); isOplMap {
						if len(opLeader) > 0 {
							var opLeaderCreate []*biz.ProjectUserConfig
							for _, v := range opLeader {
								if opLeaders, ok := v.(map[string]interface{}); ok {
									opLeaderIds = append(opLeaderIds, int(opLeaders["id"].(float64)))
									opLeaderCreate = append(opLeaderCreate, &biz.ProjectUserConfig{
										UserID:          int(opLeaders["id"].(float64)),
										ProjectConfigID: projectConfigID,
										Role:            3,
									})
								}
							}
							if !biz.SlicesEqual(opLeaderExistIds, opLeaderIds) { //判断运维负责人是否变更
								where := map[string]interface{}{
									"project_config_id": projectConfigID,
									"role":              3,
								}
								err = projectUserConfigUseCase.DeleteProjectUserConfigByWhere(&biz.ProjectUserConfigWhere{
									Conditions: where,
								})
								if err != nil {
									oe.Status = enum.ExecutableStateFailure
									return errors.Wrap(err, "project_user_config delete data failed")
								}
								if len(opLeaderCreate) > 0 {
									_, err = projectUserConfigUseCase.CreateProjectUserConfig(opLeaderCreate)
									if err != nil {
										oe.Status = enum.ExecutableStateFailure
										return errors.Wrap(err, "project_user_config table create op_leader data error")
									}
								}
							}
						}
					}

					//创建运维人员和项目的关联关系
					if opMember, isOpmMap := projectConfig.OpMember.([]interface{}); isOpmMap {
						if len(opMember) > 0 {
							var opMemberCreate []*biz.ProjectUserConfig
							for _, v := range opMember {
								if opMembers, ok := v.(map[string]interface{}); ok {
									opMemberIds = append(opMemberIds, int(opMembers["id"].(float64)))
									opMemberCreate = append(opMemberCreate, &biz.ProjectUserConfig{
										UserID:          int(opMembers["id"].(float64)),
										ProjectConfigID: projectConfigID,
										Role:            4,
									})
								}
							}
							if !biz.SlicesEqual(opLeaderExistIds, opLeaderIds) { //判断运维人员是否变更
								where := map[string]interface{}{
									"project_config_id": projectConfigID,
									"role":              4,
								}
								err = projectUserConfigUseCase.DeleteProjectUserConfigByWhere(&biz.ProjectUserConfigWhere{
									Conditions: where,
								})
								if err != nil {
									oe.Status = enum.ExecutableStateFailure
									return errors.Wrap(err, "project_user_config delete data failed")
								}
								if len(opMemberCreate) > 0 {
									_, err = projectUserConfigUseCase.CreateProjectUserConfig(opMemberCreate)
									if err != nil {
										oe.Status = enum.ExecutableStateFailure
										return errors.Wrap(err, "project_user_config table create op_member data error")
									}
								}
							}
						}
					}

				}
				oe.Status = enum.ExecutableStateSuccess
			case "delete_project_config": //本地项目删除
				var (
					projectConfig *biz.ProjectConfigDelete
					ids           []int
				)
				// 解析JSON数据到结构体切片
				err := json.Unmarshal([]byte(oe.Raw), &projectConfig)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "delete_project_config json unmarshal failed")
				}
				for _, v := range projectConfig.FormObject {
					ids = append(ids, v.ID)
				}
				if len(ids) > 0 {
					err = projectConfigUseCase.DeleteProjectConfig(0, ids)
					if err != nil {
						oe.Status = enum.ExecutableStateFailure
						return err
					}
					for _, v := range ids {
						where := map[string]interface{}{
							"project_config_id": v,
						}
						err = projectUserConfigUseCase.DeleteProjectUserConfigByWhere(&biz.ProjectUserConfigWhere{
							Conditions: where,
						})
						if err != nil {
							oe.Status = enum.ExecutableStateFailure
							return errors.Wrap(err, "project_user_config delete data failed")
						}
					}
				}
				oe.Status = enum.ExecutableStateSuccess
			case "join_project_config": //用户加入项目
				var projectUserConfig *biz.ProjectUserConfig
				// 解析JSON数据到结构体切片
				err := json.Unmarshal([]byte(oe.Raw), &projectUserConfig)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "join_project_config json unmarshal failed")
				}
				roleMaps := map[string]int{
					"rdLeader": 1,
					"rdMember": 2,
					"opLeader": 3,
					"opMember": 4,
				}
				var role int
				if roleM, ok := roleMaps[projectUserConfig.ProjectRole]; ok {
					role = roleM
				} else {
					oe.Status = enum.ExecutableStateFailure
					return errors.Errorf("unknown role")
				}
				if len(projectUserConfig.TargetProject) > 0 {
					userID := projectUserConfig.FormObject.ID
					for _, v := range projectUserConfig.TargetProject {
						err = projectUserConfigUseCase.UpsertProjectUserConfig(&biz.ProjectUserConfigUpsert{
							ProjectUserConfig: biz.ProjectUserConfig{
								UserID:          userID,
								ProjectConfigID: v.ID,
								Role:            role,
							},
						})
						if err != nil {
							oe.Status = enum.ExecutableStateFailure
							return errors.Wrap(err, "Job:join_project_config project_user_config table upsert data error")
						}
					}
				}
				oe.Status = enum.ExecutableStateSuccess
			case "leave_project_config": //用户退出项目
				var projectUserConfig *biz.ProjectUserConfig
				// 解析JSON数据到结构体切片
				err := json.Unmarshal([]byte(oe.Raw), &projectUserConfig)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "leave_project_config json unmarshal failed")
				}
				if len(projectUserConfig.TargetProject) > 0 {
					userID := projectUserConfig.FormObject.ID
					for _, v := range projectUserConfig.TargetProject {
						err = projectUserConfigUseCase.DeleteProjectUserConfigByWhere(&biz.ProjectUserConfigWhere{
							Conditions: map[string]interface{}{
								"user_id":           userID,
								"project_config_id": v.ID,
							},
						})
						if err != nil {
							oe.Status = enum.ExecutableStateFailure
							return errors.Wrap(err, "Job:join_project_config project_user_config table delete data error")
						}
					}
				}
				oe.Status = enum.ExecutableStateSuccess
			case "append_project_account_config": //项目账号添加
				var projectAccountConfig *biz.ProjectAccountConfig
				// 解析JSON数据到结构体切片
				err := json.Unmarshal([]byte(oe.Raw), &projectAccountConfig)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "create_project_account_config json unmarshal failed")
				}
				var list []*biz.ProjectAccountConfig
				if projectAccountConfig.Account != nil {
					projectAccountConfig.AccountID = projectAccountConfig.Account.ID
				}
				if projectAccountConfig.ProjectConfig != nil {
					projectAccountConfig.ProjectConfigID = projectAccountConfig.ProjectConfig.ID
				}
				if projectAccountConfig.ConnectCloudProject == "true" {
					if projectAccountConfig.Project != nil {
						projectAccountConfig.ProjectID = projectAccountConfig.Project.ID
					}
				} else {
					projectAccountConfig.ProjectID = 0
				}
				_, err = projectAccountConfigUseCase.CreateProjectAccountConfig(append(list, projectAccountConfig))
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "Job:project_account_config table create data error")
				}
				oe.Status = enum.ExecutableStateSuccess
			case "update_project_account_config": //本地项目账号编辑
				var projectAccountConfig *biz.ProjectAccountConfig
				// 解析JSON数据到结构体切片
				err := json.Unmarshal([]byte(oe.Raw), &projectAccountConfig)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return errors.Wrap(err, "project_account_config json unmarshal failed")
				}
				projectAccountConfig.ProjectID = projectAccountConfig.Project.ID
				err = projectAccountConfigUseCase.UpdateProjectAccountConfig(&biz.ProjectAccountConfigWhere{
					Query: "id = ?",
					Arg:   projectAccountConfig.FormObject.ID,
				}, projectAccountConfig)
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return err
				}
				oe.Status = enum.ExecutableStateSuccess
			default: //云端资源操作
				err := oe.CloudResourceOperate()
				if err != nil {
					oe.Status = enum.ExecutableStateFailure
					return err
				}
				oe.Status = enum.ExecutableStateSuccess
			}
		}
		return nil
	}
}

// CloudResourceOperate 云端资源操作
func (oe *ObjectTriggerExecutable) CloudResourceOperate() error {
	switch oe.Operate {
	case "create_server": // 云服务器创建
		var cloudServer *biz.CreateCloudServer
		// 解析JSON数据到结构体切片
		err := json.Unmarshal([]byte(oe.Raw), &cloudServer)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server create server json unmarshal failed")
		}
		// 获取任务ID
		step, err := jobUseCase.GetStepByStepId(oe.JobStepId)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server create server get step by id failed")
		}
		// 获取账号AK SK
		accountInfo, err := oe.GetAccountAkSk(cloudServer.Account.ID)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server create server pre operate failed")
		}
		// 创建工作环境
		statePath, outputPath, err := oe.CreateTerraformJobEnvironment(&Terraform{
			JobID:     step.JobId,
			StepID:    step.Id,
			ObjectID:  oe.Id,
			Operate:   oe.Operate,
			AccountID: accountInfo.ID,
			ProjectID: cloudServer.ProjectConfigID,
		})
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server create server create terraform job environment failed")
		}
		// 生成Terraform参数依赖
		cloudServerRebuild, err := oe.CreateTerraformServerParameterDependent(cloudServer)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server create server create terraform parameter dependent failed")
		}
		// 创建Terraform模板
		jsonData, err := json.Marshal(cloudServerRebuild)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server create server json marshal failed")
		}
		err = oe.CreateTerraformTemplates(&Terraform{
			Input:          string(jsonData),
			KeyPairName:    accountInfo.KeyPairName,
			SecretID:       accountInfo.OperateSecretID,
			SecretKey:      accountInfo.OperateSecretKey,
			ProviderID:     accountInfo.ProviderID,
			bindPublicIP:   cloudServer.BindPublicIP,
			NamePrefix:     cloudServer.NamePrefix,
			NameSuffix:     cloudServer.NameSuffix,
			OperandCounter: cloudServer.Count,
			StatePath:      statePath,
			RdsArch:        cloudServer.RdsArch,
			RedisArch:      cloudServer.RedisArch,
			Operate:        oe.Operate,
			AccountID:      accountInfo.ID,
			JobID:          step.JobId,
		}, accountInfo.ProviderID)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server create server create terraform templates failed")
		}
		// 执行Terraform Init
		err = oe.ExecTerraformInit(&Terraform{
			StatePath:  statePath,
			OutputPath: outputPath,
		})
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server create server terraform init failed")
		}
		// 执行Terraform Plan
		err = oe.ExecTerraformPlan(&Terraform{
			StatePath:  statePath,
			OutputPath: outputPath,
		})
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server create server terraform plan failed")
		}
		// 执行Terraform Apply
		err = oe.ExecTerraformApply(&Terraform{
			StatePath:  statePath,
			OutputPath: outputPath,
		})
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server create server terraform apply failed")
		}
		// 同步云数据
		request := &common.SyncCloudResourceRequest{
			Region:       cloudServer.Region.CID,
			CloudID:      accountInfo.Provider.Alias,
			AccountCID:   accountInfo.CID,
			AccountName:  accountInfo.Name,
			AccountAlias: accountInfo.Alias,
			ResourceType: "cloud_server",
			ResourceName: cloudServer.Name,
			ResourceID:   "",
		}
		err = oe.SyncCloudResource(request)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server create server sync cloud resource failed")
		}
		oe.Status = enum.ExecutableStateSuccess
	case "config_server": // 云服务器改配
		var cloudServer *biz.CloudServer
		// 解析JSON数据到结构体切片
		err := json.Unmarshal([]byte(oe.Raw), &cloudServer)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server config server json unmarshal failed")
		}
		if len(cloudServer.FormObjects) == 0 || cloudServer.FormObjects[0].ID == 0 {
			oe.Status = enum.ExecutableStateFailure
			return errors.New("cloud_server config server form object is nil")
		}
		// 前置操作
		accountInfo, resource, provider, err := oe.PreOperate(cloudServer.Account.ID, cloudServer.FormObjects[0].ID)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server config server pre operate failed")
		}
		// 执行关机操作
		err = oe.StopCloudServer(resource, provider)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server stop server failed")
		}
		// 等待关机完成
		err = oe.StopCloudServerWait(resource, provider)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server stop server wait failed")
		}
		// 执行改配操作
		spec, err := cloudServerSpecUseCase.QueryCloudServerSpec(&biz.CloudServerSpecWhere{
			Conditions: map[string]interface{}{
				"id": cloudServer.ServerSpec.ID,
			},
		}, &biz.CloudServerSpecOutput{})
		if err != nil || len(spec) == 0 {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server get server spec failed")
		}
		err = oe.ConfigCloudServer(resource, provider, spec[0].CID)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server config server failed")
		}
		time.Sleep(time.Minute * 2) // 等待2分钟执行改配操作
		// 等待改配完成
		err = oe.ConfigCloudServerWait(resource, provider, spec[0].CID, accountInfo.Provider.Alias)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server config server wait failed")
		}
		// 执行开机操作
		err = oe.StartCloudServerWait(resource, provider)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server config server start server failed")
		}
		// 同步云数据
		request := &common.SyncCloudResourceRequest{
			Region:       resource.RegionCID,
			CloudID:      accountInfo.Provider.Alias,
			AccountCID:   accountInfo.CID,
			AccountName:  accountInfo.Name,
			AccountAlias: accountInfo.Alias,
			ResourceType: "cloud_server",
			ResourceName: resource.InstanceName,
			ResourceID:   resource.InstanceCID,
		}
		err = oe.SyncCloudResource(request)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server config server sync cloud resource failed")
		}
		oe.Status = enum.ExecutableStateSuccess
	case "reboot_server": //云服务器重启
		var cloudServer *biz.CloudServer
		// 解析JSON数据到结构体切片
		err := json.Unmarshal([]byte(oe.Raw), &cloudServer)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server reboot server json unmarshal failed")
		}
		if len(cloudServer.FormObjects) == 0 || cloudServer.FormObjects[0].ID == 0 {
			oe.Status = enum.ExecutableStateFailure
			return errors.New("cloud_server reboot server form object is nil")
		}
		// 前置操作
		accountInfo, resource, provider, err := oe.PreOperate(cloudServer.Account.ID, cloudServer.FormObjects[0].ID)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server reboot server pre operate failed")
		}
		// 执行重启操作
		err = oe.RebootCloudServer(resource, provider)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server reboot server failed")
		}
		time.Sleep(time.Second * 30) // 等待30秒执行重启操作
		// 等待重启完成
		err = oe.RebootCloudServerWait(resource, provider, accountInfo.Provider.Alias)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server reboot server wait failed")
		}
		oe.Status = enum.ExecutableStateSuccess
	case "rename_server": //云服务器改名
		var cloudServer *biz.CloudServer
		// 解析JSON数据到结构体切片
		err := json.Unmarshal([]byte(oe.Raw), &cloudServer)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server rename server json unmarshal failed")
		}
		if cloudServer.FormObject.ID == 0 {
			oe.Status = enum.ExecutableStateFailure
			return errors.New("cloud_server rename server form object is nil")
		}
		// 前置操作
		accountInfo, resource, provider, err := oe.PreOperate(0, cloudServer.FormObject.ID)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server rename server pre operate failed")
		}
		// 执行改名操作
		resource.InstanceName = cloudServer.NewName
		err = oe.RenameCloudServer(resource, provider)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server rename server failed")
		}
		// 等待改名完成
		err = oe.RenameCloudServerWait(resource, provider)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server rename server wait failed")
		}
		// 同步云数据
		request := &common.SyncCloudResourceRequest{
			Region:       resource.RegionCID,
			CloudID:      accountInfo.Provider.Alias,
			AccountCID:   accountInfo.CID,
			AccountName:  accountInfo.Name,
			AccountAlias: accountInfo.Alias,
			ResourceType: "cloud_server",
			ResourceName: resource.InstanceName,
			ResourceID:   resource.InstanceCID,
		}
		err = oe.SyncCloudResource(request)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server rename server sync cloud resource failed")
		}
		oe.Status = enum.ExecutableStateSuccess
	case "reinstall_server": //云服务器重装
		var cloudServer *biz.CloudServer
		// 解析JSON数据到结构体切片
		err := json.Unmarshal([]byte(oe.Raw), &cloudServer)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server reinstall server json unmarshal failed")
		}
		if len(cloudServer.FormObjects) == 0 || cloudServer.FormObjects[0].ID == 0 {
			oe.Status = enum.ExecutableStateFailure
			return errors.New("cloud_server reinstall server form object is nil")
		}
		if cloudServer.ServerImage.ID == 0 {
			oe.Status = enum.ExecutableStateFailure
			return errors.New("cloud_server reinstall server form image id is nil")
		}
		// 前置操作
		accountInfo, resource, provider, err := oe.PreOperate(cloudServer.Account.ID, cloudServer.FormObjects[0].ID)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server reinstall server pre operate failed")
		}
		// 执行关机操作
		err = oe.StopCloudServer(resource, provider)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server stop server failed")
		}
		// 等待关机完成
		err = oe.StopCloudServerWait(resource, provider)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server stop server wait failed")
		}
		// 执行重装操作
		serverImage, err := oe.GetResource(cloudServer.ServerImage.ID, enum.ResourceTypeCloudServerImage)
		if err != nil || serverImage.ImageCID == "" {
			return errors.Wrap(err, "cloud_server get image failed")
		}
		err = oe.ReinstallCloudServer(resource, provider, serverImage.ImageCID, accountInfo.CID, accountInfo.KeyPairID, accountInfo.KeyPairName)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server reinstall server failed")
		}
		time.Sleep(time.Minute * 2) // 等待2分钟执行重装操作
		// 等待重装完成
		err = oe.ReinstallCloudServerWait(resource, provider, serverImage.ImageCID, accountInfo.Provider.Alias)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server reinstall server wait failed")
		}
		// 执行开机操作
		err = oe.StartCloudServerWait(resource, provider)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server reinstall server start server failed")
		}
		// 同步云数据
		request := &common.SyncCloudResourceRequest{
			Region:       resource.RegionCID,
			CloudID:      accountInfo.Provider.Alias,
			AccountCID:   accountInfo.CID,
			AccountName:  accountInfo.Name,
			AccountAlias: accountInfo.Alias,
			ResourceType: "cloud_server",
			ResourceName: resource.InstanceName,
			ResourceID:   resource.InstanceCID,
		}
		err = oe.SyncCloudResource(request)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server reinstall server sync cloud resource failed")
		}
		oe.Status = enum.ExecutableStateSuccess
	case "delete_server": //云服务器删除
		var cloudServer *biz.CloudServer
		// 解析JSON数据到结构体切片
		err := json.Unmarshal([]byte(oe.Raw), &cloudServer)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server delete server json unmarshal failed")
		}
		if len(cloudServer.FormObjects) == 0 || cloudServer.FormObjects[0].ID == 0 {
			oe.Status = enum.ExecutableStateFailure
			return errors.New("cloud_server delete server form object is nil")
		}
		// 前置操作
		accountInfo, resource, provider, err := oe.PreOperate(cloudServer.Account.ID, cloudServer.FormObjects[0].ID)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server delete server pre operate failed")
		}
		// 执行关机操作
		err = oe.StopCloudServer(resource, provider)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server stop server failed")
		}
		// 等待关机完成
		err = oe.StopCloudServerWait(resource, provider)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server stop server wait failed")
		}
		// 设置保留时间
		if cloudServer.ReservationPeriod != "" {
			period, err := strconv.Atoi(cloudServer.ReservationPeriod)
			if err != nil {
				oe.Status = enum.ExecutableStateFailure
				return errors.Wrap(err, "cloud_server set reservation time period failed")
			}
			err = oe.SetReservationTime(period)
			if err != nil {
				oe.Status = enum.ExecutableStateFailure
				return errors.Wrap(err, "cloud_server set reservation time failed")
			}
			// 等待保留时间完成
			err = oe.SetReservationTimeWait()
			if err != nil {
				oe.Status = enum.ExecutableStateFailure
				return errors.Wrap(err, "cloud_server set reservation time wait failed")
			}
		}
		// 变更计费方式
		reply, err := provider.DescribeServer(&cloudrepo.DescribeCloudServerReq{
			InstanceID: resource.InstanceCID,
			ProjectID:  resource.ProjectCID,
		})
		if err != nil {
			return errors.Wrap(err, "cloud_server delete server describe server failed")
		}
		if accountInfo.Provider.Alias == cloudrepo.CloudAli && reply.ChargeType != cloudrepo.CloudVmChargeTypePostPaid { // 阿里云支持删除或者释放一台按量付费实例或者到期的包年包月实例
			// 变更计费方式为按量付费
			err = oe.ChangeCloudServerChargeType(resource, provider)
			if err != nil {
				oe.Status = enum.ExecutableStateFailure
				return errors.Wrap(err, "cloud_server delete server change charge type failed")
			}
			// 等待变更计费方式完成
			err = oe.ChangeCloudServerChargeTypeWait(resource, provider)
			if err != nil {
				oe.Status = enum.ExecutableStateFailure
				return errors.Wrap(err, "cloud_server delete server change charge type wait failed")
			}
		}
		// 执行删除操作
		err = oe.DeleteCloudServer(resource, provider)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server delete server failed")
		}
		// 等待删除完成
		err = oe.DeleteCloudServerWait(resource, provider)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server delete server wait failed")
		}
		time.Sleep(time.Minute)
		// 同步云数据
		request := &common.SyncCloudResourceRequest{
			Region:       resource.RegionCID,
			CloudID:      accountInfo.Provider.Alias,
			AccountCID:   accountInfo.CID,
			AccountName:  accountInfo.Name,
			AccountAlias: accountInfo.Alias,
			ResourceType: "cloud_server",
		}
		err = oe.SyncCloudResource(request)
		if err != nil {
			oe.Status = enum.ExecutableStateFailure
			return errors.Wrap(err, "cloud_server reinstall delete sync cloud resource failed")
		}
		oe.Status = enum.ExecutableStateSuccess
	default:
		oe.Status = enum.ExecutableStateFailure
		return errors.New("CloudResourceOperate unknown operate")
	}
	return nil
}

// Commit 提交执行结果
func (oe *ObjectTriggerExecutable) Commit(ctx context.Context) error {
	select {
	case <-oe.cancel:
		return errors.Errorf("cancel by signal")
	default:
		oe.EndTime = time.Now()
		if err := jobUseCase.UpdateObject(oe.Object); err != nil {
			return err
		}
		return nil
	}
}

// Cancel 取消执行
func (oe *ObjectTriggerExecutable) Cancel() {
	close(oe.cancel)
}

// ExecutableID 获取执行id
func (oe *ObjectTriggerExecutable) ExecutableID() int {
	return oe.Id
}

// TimeCost 获取执行时长
func (oe *ObjectTriggerExecutable) TimeCost() float64 {
	return oe.EndTime.Sub(oe.StartTime).Minutes()
}

// ExecutableStatus 获取执行状态
func (oe *ObjectTriggerExecutable) ExecutableStatus() string {
	return oe.Status
}

// ExecutableType 获取执行类型
func (oe *ObjectTriggerExecutable) ExecutableType() string {
	return enum.ExecutableTypeObject
}

// ScheduleType 获取调度类型
func (oe *ObjectTriggerExecutable) ScheduleType() string {
	return enum.ScheduleTypeTrigger
}

// ObjectTriggerInput 获取对象执行输入
type ObjectTriggerInput struct {
	// step 批次
	step *biz.Step
}

// NewObjectTriggerInput 创建对象执行输入
func NewObjectTriggerInput(step *biz.Step) *ObjectTriggerInput {
	return &ObjectTriggerInput{step: step}
}

// ExecutableType 获取执行类型
func (os *ObjectTriggerInput) ExecutableType() string {
	return enum.ExecutableTypeObject
}

// Input 获取执行对象
func (os *ObjectTriggerInput) Input(ctx context.Context) ([]Executable, error) {
	list := make([]Executable, 0)
	count, err := jobUseCase.CountObjectByStepIdAndObjectStatus(os.step.Id, enum.ExecutableJobStateIdle)
	if err != nil || count == 0 {
		return list, err
	}
	doList, err := jobUseCase.ListIdleObjectByStepId(os.step.Id, 1, 0)
	if err != nil {
		return list, err
	}
	if len(doList) == 0 {
		return list, nil
	}
	idleId := make([]int, 0)
	startAt := time.Now()
	for _, do := range doList {
		do.StartTime = startAt
		idleId = append(idleId, do.Id)
		list = append(list, NewObjectTriggerExecutable(do))
	}
	if err := jobUseCase.BatchUpdateObjectStatus(os.step.Id, idleId, enum.ExecutableJobStateRunning, startAt); err != nil {
		return list, err
	}
	return list, nil
}

// InputRecover 获取恢复执行对象
func (os *ObjectTriggerInput) InputRecover(ctx context.Context) ([]Executable, error) {
	list, err := jobUseCase.ListRunningObjectByStepId(os.step.Id)
	if err != nil {
		return []Executable{}, err
	}
	restore := make([]Executable, 0)
	for _, do := range list {
		restore = append(restore, NewObjectTriggerExecutable(do))
	}
	return restore, nil
}

// InputCloser 获取是否关闭输入
func (os *ObjectTriggerInput) InputCloser(ctx context.Context) bool {
	countSuccess, err := jobUseCase.CountObjectByStepIdAndObjectStatus(os.step.Id, enum.ExecutableStateSuccess)
	if err != nil {
		return false
	}
	countClose, err := jobUseCase.CountObjectByStepIdAndObjectStatus(os.step.Id, enum.ExecutableStateClose)
	if err != nil {
		return false
	}
	if os.step.StepTotalObject > countSuccess+countClose {
		return false
	}
	return true
}
