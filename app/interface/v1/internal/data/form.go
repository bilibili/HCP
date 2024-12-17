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
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	common "github.com/bilibili/HCP/common/models"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// CreateJobSteps generate job
// 目标 根据jobId，创建 Steps 和 Objects
// 创建 job_step
func (jd *jobRepo) CreateJobSteps(jobId int, template *biz.FormTemplate, input map[string]interface{}) error {
	var list []*biz.Job
	err := jd.db.GormDB.Table("job").Where("id = ?", jobId).Model(&biz.Job{}).Find(&list).Error
	if err != nil {
		return errors.Wrap(err, "QueryJobStep")
	}
	if len(list) == 0 {
		return errors.New("job not found")
	}
	operate := list[0].Operate
	steps, err := NewStepCreator(jd).recursiveCreateStep(*template, input, nil, operate)
	if err != nil {
		return err
	}
	var total int
	for _, s := range steps {
		total += len(s.Objects)
	}
	// 遍历 Steps，按照200的容量创建其中的 objects
	for k1 := range steps {
		steps[k1].SequenceNumber = k1 + 1
		var (
			pageSize  = 100
			pageNum   = 1
			cnt       = len(steps[k1].Objects)
			objectIds []int
		)
		//创建step
		steps[k1].JobId = jobId //JobID
		err = jd.db.GormDB.Table("job_step").Create(steps[k1]).Error
		if err != nil {
			return err
		}
		//创建object
		for k2 := range steps[k1].Objects {
			steps[k1].Objects[k2].JobStepId = steps[k1].Id //批次ID
		}
		for pageNum = 1; pageNum <= (cnt-1)/pageSize+1; pageNum++ {
			p1 := (pageNum - 1) * pageSize
			p2 := pageNum * pageSize
			if p2 > cnt {
				p2 = cnt
			}
			err = jd.db.GormDB.Table("job_object").Create(steps[k1].Objects[p1:p2]).Error
			if err != nil {
				return err
			}
			for _, v := range steps[k1].Objects[p1:p2] {
				objectIds = append(objectIds, v.Id)
			}
		}
		//更新step中step_total_object数量
		err = jd.db.GormDB.Table("job_step").Where("id = ?", steps[k1].Id).Update("step_total_object", len(steps[k1].Objects)).Error
		if err != nil {
			return err
		}
	}
	var parallelStepEnable int
	if len(steps) > 1 {
		parallelStepEnable = 1
	}
	jobUpdateMap := map[string]interface{}{
		"total_step":           len(steps),
		"total_object":         total,
		"parallel_step_enable": parallelStepEnable,
	}
	err = jd.db.GormDB.Table("job").Where("id = ?", jobId).Updates(jobUpdateMap).Error
	if err != nil {
		return err
	}
	return err
}

// GetTemplate 获取模板
func (jd *jobRepo) GetTemplate(formName string) (template *biz.FormTemplate, err error) {
	conditions := map[string]interface{}{
		"name": formName,
	}
	var list []*biz.Job
	err = jd.db.GormDB.Table("job").Preload("FormTemplate").Where(conditions).Model(&biz.Job{}).Find(&list).Error
	if err != nil {
		log.Error("GetTemplate job table query data error (%v)", err)
		return nil, err
	}
	if len(list) != 1 {
		return template, fmt.Errorf("表单名称 %v 关联任务不唯一", formName)
	}
	data := list[0].FormTemplate.Data
	var formTemplate *biz.FormTemplate

	if err = json.Unmarshal([]byte(data), &formTemplate); err != nil {
		log.Error("GetTemplate json unmarshal error (%v)", err)
		return template, err
	}
	return formTemplate, nil
}

// GetTemplateContent 将Job数据填充进模板
func (jd *jobRepo) GetTemplateContent(jobId int) (input map[string]interface{}, err error) {
	conditions := map[string]interface{}{
		"id": jobId,
	}
	var list []*biz.Job
	err = jd.db.GormDB.Select("id, raw").Table("job").Where(conditions).Model(&biz.Job{}).Find(&list).Error
	if err != nil {
		log.Error("GetTemplateContent job table query data error (%v)", err)
		return nil, err
	}
	if len(list) != 1 {
		return input, errors.Errorf("GetTemplateContent: query job failed, jobId(%v)", jobId)
	}
	if err := json.Unmarshal([]byte(list[0].Raw), &input); err != nil {
		return input, errors.Wrap(err, "GetTemplateContent unmarshal body")
	}
	return input, nil
}

// StepCreator 生成步骤
type StepCreator struct {
	// Job
	repo *jobRepo
	// 步骤数
	count int
	// 模版
	template string
	// 名称前缀
	namePrefix string
	// 已占用数
	occupied float64
	// 已占用map
	occupiedMap map[string]float64
	// 标题名称
	titleToName map[string]string
}

// NewStepCreator 创建步骤
func NewStepCreator(repo *jobRepo) *StepCreator {
	return &StepCreator{
		repo:        repo,
		count:       1,
		occupiedMap: make(map[string]float64, 0),
		titleToName: make(map[string]string, 0),
	}
}

// generateJobStep 生成步骤
func (s *StepCreator) recursiveCreateStep(template biz.FormTemplate, input map[string]interface{}, parent map[string]interface{}, operate string) ([]*biz.Step, error) {
	var create []*biz.Step
	if template.MainOperand != nil {
		step := s.generateJobStep(template, input, parent, operate)
		create = append(create, step)
	}
	for _, property := range template.Parameters {
		if property.Type != "object" {
			continue
		}
		if property.Reference != "FormTemplate" {
			continue
		}
		_, ok := input[property.Name]
		if !ok {
			continue
		}
		// 联合模型
		jointInputs := input[property.Name].([]interface{})
		delete(input, property.Name)
		createCnt := len(create)
		flag := true
		for _, jointInput := range jointInputs {
			var subName string
			for _, t := range property.Templates {
				if _, ok := jointInput.(map[string]interface{})["joint"+t.Name]; !ok {
					continue
				}
				subName = t.Name
			}
			if subName == "" {
				return create, errors.Errorf("the joint does not contains any subform template.")
			}
			subInput := jointInput.(map[string]interface{})["joint"+subName]
			subTmpl, err := s.repo.getFormTemplate(subName)
			if err != nil {
				return create, err
			}
			if subTmpl.MainOperand == nil {
				continue
			}
			// 获取子表单operate
			templateData, err := s.repo.getFormTemplateData(subTmpl.Name)
			if err != nil {
				return create, errors.Errorf("The sub formtemplate model name is incorrect.")
			}
			steps, err := s.recursiveCreateStep(subTmpl, subInput.(map[string]interface{}), input, templateData.Operate)
			if property.Style == "table" {
				if flag {
					create = append(create, steps...)
					flag = false
					continue
				}
				create[createCnt].Objects = append(create[createCnt].Objects, steps[0].Objects...)
				continue
			}
			create = append(create, steps...)
		}
	}
	return create, nil
}

// SetOccupied 设置标题名称
func (s *StepCreator) SetOccupied(template, namePrefix string, occupied float64) {
	s.occupiedMap[fmt.Sprintf("%v:%v", template, namePrefix)] = occupied
}

// GetOccupied 获取标题名称
func (s *StepCreator) GetOccupied(template, namePrefix string) float64 {
	return s.occupiedMap[fmt.Sprintf("%v:%v", template, namePrefix)]
}

// generateJobStep 生成步骤
func (s *StepCreator) generateJobStep(template biz.FormTemplate, formInput map[string]interface{}, shared map[string]interface{}, operate string) *biz.Step {
	formName := formInput["formName"].(string)
	formTitle := formName
	if formInput["formTitle"] != nil {
		formTitle = formInput["formTitle"].(string)
	}
	s.titleToName[formTitle] = formName

	step := &biz.Step{
		Name:           fmt.Sprintf("step%v-%v", s.count, template.Name),
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
		SequenceNumber: s.count,
		StartTime:      timeParse("0001-01-01T00:00:00+08:00"),
		EndTime:        timeParse("0001-01-01T00:00:00+08:00"),
		Status:         "idle",
	}
	s.count += 1

	// 非自动生成的FormTemplate关联
	conditions := map[string]interface{}{
		"name": template.Name,
	}
	var list []*biz.FormTemplateData
	err := s.repo.db.GormDB.Table("form_template_data").Where(conditions).Model(biz.FormTemplateData{}).Find(&list).Error
	if err != nil {
		log.Error("generateJobStep form_template_data table query data error (%v)", err)
		return step
	}
	if len(list) > 0 {
		step.FormTemplateID = list[0].ID
	}

	// 配置步骤默认状态
	if template.InitialStepStatus != "" {
		step.Status = template.InitialStepStatus
	}
	// 判断主对象存在
	operandParameters := getOperandParameter(template, formInput)

	// 未找到引用模型为主操作对象字段输入，判断为实例创建
	if len(operandParameters) == 0 {
		// 如果有数量输入，有则多实例创建，无则单实例创建
		namePrefix, ok := formInput["name_prefix"]
		if !ok {
			namePrefix = "object"
		}
		nameSuffix, ok := formInput["name_suffix"]
		if !ok {
			nameSuffix = "none"
		}
		occupied := s.repo.getMaximumOccupied(namePrefix.(string), template, s.GetOccupied(template.Name, namePrefix.(string)))
		count, ok := formInput["count"]
		if !ok || count == 0 {
			step.Objects = append(step.Objects, newJobObject(namePrefix.(string), nameSuffix.(string), operate, 0, occupied, template, formInput))
			s.SetOccupied(template.Name, namePrefix.(string), occupied+1)
			return step
		}
		for i := float64(0); i < count.(float64); i += 1 {
			step.Objects = append(step.Objects, newJobObject(namePrefix.(string), nameSuffix.(string), operate, i, occupied, template, formInput))
		}
		s.SetOccupied(template.Name, namePrefix.(string), occupied+count.(float64))
		return step
	}
	objectIndex := 1
	// 遍历所有主操作对象字段输入，判断为任务对象
	for _, p := range operandParameters {
		objectInput, ok := formInput[p.Name]
		if !ok {
			continue
		}
		// 单实例输入
		if !p.List {
			step.Objects = append(step.Objects, generateJobObject(objectIndex, template, objectInput.(map[string]interface{}), operate, formInput))
			objectIndex += 1
			continue
		}
		// 多实例输入
		for _, o := range objectInput.([]interface{}) {
			step.Objects = append(step.Objects, generateJobObject(objectIndex, template, o.(map[string]interface{}), operate, formInput))
			objectIndex += 1
		}
	}
	return step
}

// getOperandParameter 获取主操作对象字段输入
func getOperandParameter(template biz.FormTemplate, formInput map[string]interface{}) []biz.FormParameter {
	operandParameters := make([]biz.FormParameter, 0)

	// FormAdd 表单可能包含多个引用自己的模型， 如 Organization
	if template.MainOperand != nil && !strings.HasPrefix(template.Name, "FormAdd") {
		for i, sub := range template.Parameters {
			annotation := newAnnotation(sub.Annotation)
			ignore, ok := annotation["ignoreMainOperand"]
			if ok && ignore == true {
				continue
			}
			if sub.Reference != "FormTemplate" && sub.Reference == template.MainOperand.Name {
				operandParameters = append(operandParameters, template.Parameters[i])
			}
		}
	} else if template.MainOperand == nil {
		template.MainOperand = &biz.MetaSchema{Name: "Other", Title: "其他"}
	}
	return operandParameters
}

// newAnnotation 解析注解
func newAnnotation(str string) map[string]interface{} {
	a := make(map[string]interface{}, 0)
	_ = json.Unmarshal([]byte(str), &a)
	return a
}

// newJobObject 生成任务对象
func newJobObject(namePrefix, nameSuffix, operate string, index, occupied float64, template biz.FormTemplate, formInput map[string]interface{}) *biz.Object {
	code := fmt.Sprintf("%v-object%v", template.MainOperand.Name, index)
	name := code
	operandCounter := occupied + 1 + index
	switch nameSuffix {
	case "none":
		name = namePrefix
	case "standard-01":
		if occupied+1+index < 10 {
			name = fmt.Sprintf("%v-0%v", namePrefix, occupied+1+index)
		} else {
			name = fmt.Sprintf("%v-%v", namePrefix, occupied+1+index)
		}
	case "standard-1":
		name = fmt.Sprintf("%v-%v", namePrefix, occupied+1+index)
	}
	bytes, _ := json.Marshal(formInput)
	object := &biz.Object{
		Code:           code,
		Name:           name,
		Operate:        operate,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
		StartTime:      timeParse("0001-01-01T00:00:00+08:00"),
		EndTime:        timeParse("0001-01-01T00:00:00+08:00"),
		Status:         "idle",
		Raw:            string(bytes),
		OperandCounter: int(operandCounter),
	}
	return object
}

// getMaximumOccupiedData 获取最大占用数据
func (jd *jobRepo) getMaximumOccupiedData(namePrefix string, template biz.FormTemplate) string {
	tableName := modelConversionTable(template.MainOperand.Name)
	if tableName == "" {
		return ""
	}
	conditions := map[string]interface{}{
		"name LIKE ?": fmt.Sprintf("%%%s%%", namePrefix),
	}
	switch tableName {
	case common.CloudServer:
		var list []*biz.CloudServer
		serverSession := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table(tableName).Model(&biz.CloudServer{})
		for key, value := range conditions {
			serverSession.Where(key, value)
		}
		err := serverSession.Find(&list).Error
		if err != nil {
			log.Error("getMaximumOccupiedData cloud_server table parse data error(%v)", err)
			return ""
		}
		if len(list) < 1 {
			return ""
		}
		return list[0].Name
	}
	return ""
}

// getMaximumOccupied 获取最大占用
func (jd *jobRepo) getMaximumOccupied(namePrefix string, template biz.FormTemplate, occupied float64) float64 {
	name := jd.getMaximumOccupiedData(namePrefix, template)
	if name == "" {
		return math.Max(occupied, 0)
	}
	list := strings.Split(name, "-")
	suffix := list[len(list)-1]
	num, _ := strconv.Atoi(suffix)
	occupied = math.Max(occupied, float64(num))
	return occupied
}

// generateJobObject 生成任务对象
func generateJobObject(index int, template biz.FormTemplate, operandInput map[string]interface{}, operate string, formInput map[string]interface{}) *biz.Object {
	operandId := operandInput["id"]
	operandName, ok := operandInput["name"]
	if !ok {
		operandName = fmt.Sprintf("object%v", index)
	}
	delete(operandInput, "oid")
	formInput["formObjects"] = []map[string]interface{}{
		{"id": operandId},
	}
	bytes, _ := json.Marshal(formInput)
	create := &biz.Object{
		Code:       fmt.Sprintf("%v-%v", template.MainOperand.Name, operandId),
		Name:       fmt.Sprintf("%v", operandName),
		Operate:    operate,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		StartTime:  timeParse("0001-01-01T00:00:00+08:00"),
		EndTime:    timeParse("0001-01-01T00:00:00+08:00"),
		Status:     "idle",
		Raw:        string(bytes),
	}
	return create
}

// getFormTemplateData 获取表单模板数据
func (jd *jobRepo) getFormTemplateData(name string) (biz.FormTemplateData, error) {
	conditions := map[string]interface{}{
		"name": name,
	}
	var list []*biz.FormTemplateData
	err := jd.db.GormDB.Table("form_template_data").Where(conditions).Model(&biz.FormTemplateData{}).Find(&list).Error
	if err != nil {
		return biz.FormTemplateData{}, errors.Wrap(err, "mysql: handleForm form_template_data failed")
	}
	if len(list) == 0 {
		return biz.FormTemplateData{}, errors.New("mysql: handleForm failed, form_template_data not found")
	}
	return *list[0], nil
}

// getFormTemplate 获取表单模板
func (jd *jobRepo) getFormTemplate(name string) (biz.FormTemplate, error) {
	conditions := map[string]interface{}{
		"name": name,
	}
	var list []*biz.FormTemplateData
	err := jd.db.GormDB.Table("form_template_data").Where(conditions).Model(&biz.FormTemplateData{}).Find(&list).Error
	if err != nil {
		return biz.FormTemplate{}, errors.Wrap(err, "mysql: handleForm failed")
	}
	if len(list) == 0 {
		return biz.FormTemplate{}, errors.New("mysql: handleForm failed, template not found")
	}
	data := list[0].Data
	var formTemplate biz.FormTemplate
	if err := json.Unmarshal([]byte(data), &formTemplate); err != nil {
		return biz.FormTemplate{}, errors.Wrap(err, "getFormTemplate failed, json unmarshal error")
	}
	return formTemplate, nil
}
