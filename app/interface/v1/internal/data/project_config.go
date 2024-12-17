// Package data
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
package data

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/go-kratos/kratos/pkg/log"
	"gorm.io/gorm"
)

// ProjectConfigRepo 本地项目DB
type ProjectConfigRepo struct {
	db *biz.DB
}

// NewProjectConfigRepo 初始化本地项目DB
func NewProjectConfigRepo(db *biz.DB) biz.ProjectConfigRepo {
	return &ProjectConfigRepo{db: db}
}

// ParseProjectConfig 本地项目查询条件处理
func (c *ProjectConfigRepo) ParseProjectConfig(where *biz.ProjectConfigWhere) (conditions map[string]interface{}, err error) {
	conditionMaps, conditionSources := make(map[string]interface{}), make(map[string]interface{})
	if rdLeader, ok := where.Conditions["rd_leader"]; ok && rdLeader != nil {
		jsonBytes, err := json.Marshal(rdLeader)
		if err != nil {
			return nil, err
		}
		rdLeaderConditions, err := biz.ParseCloudData(jsonBytes)
		if err != nil {
			log.Error("project_config rd_leader user table parse data error(%v)", err)
			return nil, err
		}
		//查询用户ID
		userSession := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("user").Model(&biz.User{})
		if where != nil {
			for key, value := range rdLeaderConditions {
				userSession.Where(key, value)
			}
		}
		var rdLeaderList []*biz.User                               //研发负责人列表
		var rdLeaderUserIds []int                                  //研发负责人ID集
		var rdLeaderProjectUserConfigList []*biz.ProjectUserConfig //研发负责人关联项目集
		var rdLeaderProjectConfigIds []int                         //研发负责人关联项目的ID集
		err = userSession.Find(&rdLeaderList).Error
		if err != nil {
			log.Error("project_config user table query rd_leader data error(%v)", err)
			return nil, err
		}
		for _, v := range rdLeaderList {
			rdLeaderUserIds = append(rdLeaderUserIds, v.ID)
		}
		if len(rdLeaderUserIds) > 0 {
			//查询用户和项目的关联关系
			err = c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("project_user_config").Model(&biz.ProjectUserConfig{}).
				Where("user_id IN ?", rdLeaderUserIds).Find(&rdLeaderProjectUserConfigList).Error
			if err != nil {
				log.Error("project_config project_user_config table query rd_leader data error(%v)", err)
				return nil, err
			}
			for _, v := range rdLeaderProjectUserConfigList {
				rdLeaderProjectConfigIds = append(rdLeaderProjectConfigIds, v.ProjectConfigID)
			}
		}
		if len(rdLeaderProjectConfigIds) > 0 {
			conditionMaps["id IN ?"] = rdLeaderProjectConfigIds
		}
		conditionSources["rd_leader"] = true
	}
	if rdMember, ok := where.Conditions["rd_member"]; ok && rdMember != nil {
		jsonBytes, err := json.Marshal(rdMember)
		if err != nil {
			return nil, err
		}
		rdMemberConditions, err := biz.ParseCloudData(jsonBytes)
		if err != nil {
			log.Error("project_config rd_member user table parse data error(%v)", err)
			return nil, err
		}
		//查询用户ID
		userSession := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("user").Model(&biz.User{})
		if where != nil {
			for key, value := range rdMemberConditions {
				userSession.Where(key, value)
			}
		}
		var rdMemberList []*biz.User                               //研发人员列表
		var rdMemberUserIds []int                                  //研发人员ID集
		var rdMemberProjectUserConfigList []*biz.ProjectUserConfig //研发人员关联项目集
		var rdMemberProjectConfigIds []int                         //研发人员关联项目的ID集
		err = userSession.Find(&rdMemberList).Error
		if err != nil {
			log.Error("project_config user table query rd_member data error(%v)", err)
			return nil, err
		}
		for _, v := range rdMemberList {
			rdMemberUserIds = append(rdMemberUserIds, v.ID)
		}
		if len(rdMemberUserIds) > 0 {
			//查询用户和项目的关联关系
			err = c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("project_user_config").Model(&biz.ProjectUserConfig{}).
				Where("user_id IN ?", rdMemberUserIds).Find(&rdMemberProjectUserConfigList).Error
			if err != nil {
				log.Error("project_config project_user_config table query rd_member data error(%v)", err)
				return nil, err
			}
			for _, v := range rdMemberProjectUserConfigList {
				rdMemberProjectConfigIds = append(rdMemberProjectConfigIds, v.ProjectConfigID)
			}
		}
		if len(rdMemberProjectConfigIds) > 0 {
			conditionMaps["id IN ?"] = rdMemberProjectConfigIds
		}
		conditionSources["rd_member"] = true
	}
	if opLeader, ok := where.Conditions["op_leader"]; ok && opLeader != nil {
		jsonBytes, err := json.Marshal(opLeader)
		if err != nil {
			return nil, err
		}
		opLeaderConditions, err := biz.ParseCloudData(jsonBytes)
		if err != nil {
			log.Error("project_config op_leader user table parse data error(%v)", err)
			return nil, err
		}
		//查询用户ID
		userSession := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("user").Model(&biz.User{})
		if where != nil {
			for key, value := range opLeaderConditions {
				userSession.Where(key, value)
			}
		}
		var opLeaderList []*biz.User                               //运维负责人列表
		var opLeaderUserIds []int                                  //运维负责人ID集
		var opLeaderProjectUserConfigList []*biz.ProjectUserConfig //运维负责人关联项目集
		var opLeaderProjectConfigIds []int                         //运维负责人关联项目的ID集
		err = userSession.Find(&opLeaderList).Error
		if err != nil {
			log.Error("project_config user table query op_Leader data error(%v)", err)
			return nil, err
		}
		for _, v := range opLeaderList {
			opLeaderUserIds = append(opLeaderUserIds, v.ID)
		}
		if len(opLeaderUserIds) > 0 {
			//查询用户和项目的关联关系
			err = c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("project_user_config").Model(&biz.ProjectUserConfig{}).
				Where("user_id IN ?", opLeaderUserIds).Find(&opLeaderProjectUserConfigList).Error
			if err != nil {
				log.Error("project_config project_user_config table query op_leader data error(%v)", err)
				return nil, err
			}
			for _, v := range opLeaderProjectUserConfigList {
				opLeaderProjectConfigIds = append(opLeaderProjectConfigIds, v.ProjectConfigID)
			}
		}
		if len(opLeaderProjectConfigIds) > 0 {
			conditionMaps["id IN ?"] = opLeaderProjectConfigIds
		}

		conditionSources["op_leader"] = true
	}
	if opMember, ok := where.Conditions["op_member"]; ok && opMember != nil {
		jsonBytes, err := json.Marshal(opMember)
		if err != nil {
			return nil, err
		}
		opMemberConditions, err := biz.ParseCloudData(jsonBytes)
		if err != nil {
			log.Error("project_config op_member user table parse data error(%v)", err)
			return nil, err
		}
		//查询用户ID
		userSession := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("user").Model(&biz.User{})
		if where != nil {
			for key, value := range opMemberConditions {
				userSession.Where(key, value)
			}
		}
		var opMemberList []*biz.User                               //运维人员列表
		var opMemberUserIds []int                                  //运维人员ID集
		var opMemberProjectUserConfigList []*biz.ProjectUserConfig //运维人员关联项目集
		var opMemberProjectConfigIds []int                         //运维人员关联项目的ID集
		err = userSession.Find(&opMemberList).Error
		if err != nil {
			log.Error("project_config user table query op_member data error(%v)", err)
			return nil, err
		}
		for _, v := range opMemberList {
			opMemberUserIds = append(opMemberUserIds, v.ID)
		}
		if len(opMemberUserIds) > 0 {
			//查询用户和项目的关联关系
			err = c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("project_user_config").Model(&biz.ProjectUserConfig{}).
				Where("user_id IN ?", opMemberUserIds).Find(&opMemberProjectUserConfigList).Error
			if err != nil {
				log.Error("project_config project_user_config table query op_member data error(%v)", err)
				return nil, err
			}
			for _, v := range opMemberProjectUserConfigList {
				opMemberProjectConfigIds = append(opMemberProjectConfigIds, v.ProjectConfigID)
			}
		}
		if len(opMemberProjectConfigIds) > 0 {
			conditionMaps["id IN ?"] = opMemberProjectConfigIds
		}
		conditionSources["rd_member"] = true
	}
	jsonBytes, err := json.Marshal(conditionMaps)
	if err != nil {
		return nil, err
	}
	conditionMaps, err = biz.ParseCloudData(jsonBytes)
	if err != nil {
		log.Error("project_config table parse data error(%v)", err)
		return nil, err
	}
	conditionDiff := biz.DifferenceMap(where.Conditions, conditionSources)
	if conditionDiff != nil {
		jsonByteDiffs, err := json.Marshal(conditionDiff)
		if err != nil {
			return nil, err
		}
		conditionDos, err := biz.ParseCloudData(jsonByteDiffs)
		if err != nil {
			log.Error("project_config table parse data error(%v)", err)
			return nil, err
		}
		conditions = c.ParesQueryProjectConfigConditions(conditionDos)
		for k, v := range conditionMaps {
			conditions[k] = v
		}
	}
	return conditions, nil
}

// ParesQueryProjectConfigConditions 解析查询条件（根据用户筛选）
func (c *ProjectConfigRepo) ParesQueryProjectConfigConditions(conditions map[string]interface{}) map[string]interface{} {
	var ids []int
	conditionsMaps, conditionReturns := make(map[string]interface{}), make(map[string]interface{})
	and, isAnd := conditions["and"]
	or, isOr := conditions["or"]
	if !isAnd && !isOr {
		return conditions
	}
	if or != nil && isOr {
		if mapValues, isMap := or.(map[string]interface{}); isMap {
			for k, v := range mapValues {
				var roleStr string
				field := k
				parts := strings.Split(k, ".") //分割字符串
				if len(parts) > 1 {            //用户搜索
					field = parts[1]
					roleStr = parts[0] //获取权限
					if roleStr != "" { //判断权限
						roleMaps := map[string]int{
							"RdLeader": 1,
							"RdMember": 2,
							"OpLeader": 3,
							"OpMember": 4,
						}
						if roleM, ok := roleMaps[roleStr]; ok {
							conditionsMaps["role = ?"] = roleM
						}
					}
					conditionsMaps[field] = v
					resp, err := c.QueryByRoleProjectUserConfigRelation(conditionsMaps)
					if err != nil {
						return nil
					}
					ids = append(ids, resp...)
					if len(ids) > 0 {
						conditionReturns["id IN ?"] = ids
					}
				} else { //普通搜索
					conditionReturns[k] = v
				}
			}
			delete(conditions, "or")
		}
	}
	if and != nil && isAnd {
		if mapValues, isMap := and.([]map[string]interface{}); isMap {
			for _, v := range mapValues {
				for k1, v1 := range v {
					var role int
					field := k1
					parts := strings.Split(k1, ".") //分割字符串
					if len(parts) > 1 {             //用户搜索
						roleStr := parts[0] //获取权限
						if roleStr != "" {  //判断权限
							roleMaps := map[string]int{
								"RdLeader": 1,
								"RdMember": 2,
								"OpLeader": 3,
								"OpMember": 4,
							}
							if roleM, ok := roleMaps[roleStr]; ok {
								role = roleM
							}
						}
						field = parts[1]
						conditionsMaps[field] = v1
						resp, err := c.QueryProjectUserConfigRelation(conditionsMaps, role)
						if err != nil {
							return nil
						}
						ids = append(ids, resp...)
						if len(ids) > 0 {
							conditionReturns["id IN ?"] = ids
						}
					} else { //普通搜索
						conditionReturns[k1] = v1
					}
				}
			}
			delete(conditions, "and")
		}
	}
	return conditionReturns
}

// QueryProjectUserConfigRelation 查询用户和项目的关联关系
func (c *ProjectConfigRepo) QueryProjectUserConfigRelation(conditions map[string]interface{}, role int) ([]int, error) {
	//查询用户ID
	userSession := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("user").Model(&biz.User{})
	for key, value := range conditions {
		userSession.Where(key, value)
	}
	var userList []*biz.User                           //用户列表
	var UserIds []int                                  //用户ID集
	var projectUserConfigList []*biz.ProjectUserConfig //用户关联项目集
	var projectConfigIds []int                         //用户关联项目的ID集
	err := userSession.Find(&userList).Error
	if err != nil {
		return nil, err
	}
	for _, v := range userList {
		UserIds = append(UserIds, v.ID)
	}
	if len(UserIds) > 0 {
		//查询用户和项目的关联关系
		conditionMaps := make(map[string]interface{})
		conditionMaps["user_id IN ?"] = UserIds
		if role != 0 {
			conditionMaps["role = ?"] = role
		}
		userConfigSession := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("project_user_config").Model(&biz.ProjectUserConfig{})
		for key, value := range conditionMaps {
			userConfigSession.Where(key, value)
		}
		err = userConfigSession.Find(&projectUserConfigList).Error
		if err != nil {
			return nil, err
		}
		for _, v := range projectUserConfigList {
			projectConfigIds = append(projectConfigIds, v.ProjectConfigID)
		}
	}
	return projectConfigIds, nil
}

// QueryByRoleProjectUserConfigRelation 根据已知权限查询用户和项目的关联关系
func (c *ProjectConfigRepo) QueryByRoleProjectUserConfigRelation(conditions map[string]interface{}) ([]int, error) {
	//查询用户和项目的关联关系
	var projectUserConfigList []*biz.ProjectUserConfig //用户关联项目集
	var projectConfigIds []int                         //用户关联项目的ID集
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("project_user_config").Model(&biz.ProjectUserConfig{})
	for key, value := range conditions {
		session.Where(key, value)
	}
	err := session.Find(&projectUserConfigList).Error
	if err != nil {
		return nil, err
	}
	for _, v := range projectUserConfigList {
		projectConfigIds = append(projectConfigIds, v.ProjectConfigID)
	}
	return projectConfigIds, nil
}

// CountProjectConfig 本地项目查询数量
func (c *ProjectConfigRepo) CountProjectConfig(where *biz.ProjectConfigWhere) (total int64, err error) {
	conditions, err := c.ParseProjectConfig(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("project_config").Model(&biz.ProjectConfig{})
	if where != nil {
		if or, ok := conditions["or"]; ok && or != nil { // or条件
			if mapValues, isMap := or.(map[string]interface{}); isMap {
				for k, v := range mapValues {
					session = session.Or(k, v)
				}
			}
			delete(conditions, "or")
		}
		if and, ok := conditions["and"]; ok && and != nil { // and条件
			if mapValues, isMap := and.([]map[string]interface{}); isMap {
				for _, v := range mapValues {
					for k1, v1 := range v {
						session = session.Where(k1, v1)
					}
				}
			}
			delete(conditions, "and")
		}
		for key, value := range conditions {
			session.Where(key, value)
		}
	}
	err = session.Count(&total).Error
	if err != nil {
		log.Error("project_config table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryProjectConfig 本地项目查询
func (c *ProjectConfigRepo) QueryProjectConfig(where *biz.ProjectConfigWhere, output *biz.ProjectConfigOutput) (list []*biz.ProjectConfig, err error) {
	conditions, err := c.ParseProjectConfig(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := c.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("project_config").Model(&biz.ProjectConfig{})
	if where != nil {
		if or, ok := conditions["or"]; ok && or != nil { // or条件
			if mapValues, isMap := or.(map[string]interface{}); isMap {
				for k, v := range mapValues {
					session = session.Or(k, v)
				}
			}
			delete(conditions, "or")
		}
		if and, ok := conditions["and"]; ok && and != nil { // and条件
			if mapValues, isMap := and.([]map[string]interface{}); isMap {
				for _, v := range mapValues {
					for k1, v1 := range v {
						session = session.Where(k1, v1)
					}
				}
			}
			delete(conditions, "and")
		}
		for key, value := range conditions {
			session.Where(key, value)
		}
	}
	if output.PageSize != 0 && output.PageNum != 0 {
		session.Limit(output.PageSize).Offset((output.PageNum - 1) * output.PageSize)
	}
	if output.Order != "" {
		session.Order(output.Order)
	}
	err = session.
		Find(&list).Error
	for _, v := range list {
		var (
			projectUserConfigList []*biz.ProjectUserConfig
			rdLeader              []*biz.User
			rdMember              []*biz.User
			opLeader              []*biz.User
			opMember              []*biz.User
		)
		//查询用户和项目的关联关系
		err = c.db.GormDB.Table("project_user_config").Where("project_user_config.project_config_id = ?", v.ID).
			Joins("User", c.db.GormDB.Select("id", "name", "role")).Find(&projectUserConfigList).Error
		if err != nil {
			log.Error("project_config project_user_config table query data error (%v)", err)
			return nil, err
		}
		if len(projectUserConfigList) > 0 {
			for _, v := range projectUserConfigList {
				if v.Role == 1 {
					rdLeader = append(rdLeader, v.User)
				} else if v.Role == 2 {
					rdMember = append(rdMember, v.User)
				} else if v.Role == 3 {
					opLeader = append(opLeader, v.User)
				} else if v.Role == 4 {
					opMember = append(opMember, v.User)
				}
			}
		}
		v.RdLeader = rdLeader
		v.RdMember = rdMember
		v.OpLeader = opLeader
		v.OpMember = opMember
	}
	if err != nil {
		log.Error("project_config table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// CreateProjectConfig 本地项目创建
func (c *ProjectConfigRepo) CreateProjectConfig(create []*biz.ProjectConfig) (list []*biz.ProjectConfig, err error) {
	for _, v := range create {
		v.CreateTime = time.Now()
		v.UpdateTime = time.Now()
		v.DeleteTime = timeParse("0001-01-01T00:00:00+08:00")
		err = c.db.GormDB.Table("project_config").Create(v).Error
		if err != nil {
			log.Error("project_config table create data error(%v)", err)
			return nil, err
		}
		list = append(list, &biz.ProjectConfig{
			ID: v.ID,
		})
	}
	return list, nil
}

// UpdateProjectConfig 本地项目更新
func (c *ProjectConfigRepo) UpdateProjectConfig(where *biz.ProjectConfigWhere, update *biz.ProjectConfig) error {
	updateMap := make(map[string]interface{})
	if update != nil {
		updateMap["name"] = update.Name
		updateMap["alias"] = update.Alias
		updateMap["description"] = update.Description
		updateMap["status"] = update.Status
		updateMap["is_delete"] = update.IsDelete
	}
	if isDelete, ok := where.Conditions["is_delete"]; ok && isDelete == 1 {
		updateMap["delete_time"] = time.Now()
	}
	if where != nil && len(updateMap) > 0 {
		err := c.db.GormDB.Table("project_config").Where(where.Query, where.Arg).Updates(updateMap).Error
		if err != nil {
			log.Error("project_config table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// DeleteProjectConfig 本地项目删除
func (c *ProjectConfigRepo) DeleteProjectConfig(softDelete int, deleteID []int) error {
	if softDelete == 1 { // 软删除
		now := time.Now()
		deleteTime := now.Format("2006-01-02 15:04:05")
		var updateMap map[string]interface{}
		updateMap["is_delete"] = 1
		updateMap["delete_time"] = deleteTime
		err := c.db.GormDB.Table("project_config").Where("id IN ?", deleteID).Updates(updateMap).Error
		if err != nil {
			log.Error("project_config table delete data error(%v)", err)
			return nil
		}
	} else {
		err := c.db.GormDB.Table("project_config").Where("id IN ?", deleteID).Delete(&biz.ProjectConfig{}).Error
		if err != nil {
			log.Error(" project_config table update data error(%v)", err)
			return err
		}
	}
	return nil
}

// UpsertProjectConfig 本地项目更新或创建
func (c *ProjectConfigRepo) UpsertProjectConfig(upsert []*biz.ProjectConfigUpsert) error {
	return nil
}
