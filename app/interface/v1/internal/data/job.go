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
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/app/interface/v1/internal/biz/enum"
	common "github.com/bilibili/HCP/common/models"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// jobRepo 任务DB
type jobRepo struct {
	db *biz.DB
}

// NewJobRepo 初始化任务DB
func NewJobRepo(db *biz.DB) biz.JobRepo {
	return &jobRepo{
		db: db,
	}
}

// timeFormat 时间格式化
func timeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// timeParse 时间解析
func timeParse(str string) (t time.Time) {
	t, _ = time.Parse("2006-01-02T15:04:05-07:00", str)
	return
}

// ModelConversionTable 模型转换表
func modelConversionTable(model string) string {
	var tableMap = map[string]string{
		"CloudServer": common.CloudServer,
	}
	if v, ok := tableMap[model]; ok {
		return v
	}
	return ""
}

// ParseJob 任务查询条件处理
func (jd *jobRepo) ParseJob(where *biz.JobWhere) (conditions map[string]interface{}, err error) {
	jsonBytes, err := json.Marshal(where.Conditions)
	if err != nil {
		return nil, err
	}
	conditions, err = biz.ParseCloudData(jsonBytes)
	if err != nil {
		log.Error("job table parse data error(%v)", err)
		return nil, err
	}
	return conditions, nil
}

// CreateJob 创建任务
func (jd *jobRepo) CreateJob(job *biz.Job) (int, error) {
	job.CreateTime = time.Now()
	job.UpdateTime = time.Now()
	job.StartTime = timeParse("0001-01-01T00:00:00+08:00")
	job.EndTime = timeParse("0001-01-01T00:00:00+08:00")
	job.AuditTime = timeParse("0001-01-01T00:00:00+08:00")
	//查询模版
	var template []*biz.FormTemplateData
	err := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("form_template_data").Model(&biz.FormTemplateData{}).
		Where("id = ?", job.FormTemplateID).Find(&template).Error
	if err != nil {
		return 0, errors.Wrap(err, "data: CreateJob Query FormTemplate failed")
	}
	if len(template) == 0 {
		return 0, errors.New("data: CreateJob Query FormTemplate Data is empty")
	}
	job.Name = fmt.Sprintf("%v-%v", template[0].Name, time.Now().UnixNano())
	job.Title = fmt.Sprintf("【HCP】%v", template[0].Title)
	job.Operate = template[0].Operate
	err = jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job").Create(job).Error
	if err != nil {
		return 0, errors.Wrap(err, "data: CreateJob failed")
	}
	return job.Id, nil
}

// CountJob 任务查询数量
func (jd *jobRepo) CountJob(where *biz.JobWhere) (total int64, err error) {
	conditions, err := jd.ParseJob(where) //条件转换
	if err != nil {
		return 0, err
	}
	session := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job").Model(&biz.Job{})
	if where != nil {
		for key, value := range conditions {
			session.Where(key, value)
		}
	}
	err = session.Count(&total).Error
	if err != nil {
		log.Error("job table query data count error (%v)", err)
		return 0, err
	}
	return total, nil
}

// QueryJob 任务查询
func (jd *jobRepo) QueryJob(where *biz.JobWhere, output *biz.JobOutput) (list []*biz.Job, err error) {
	conditions, err := jd.ParseJob(where) //条件转换
	if err != nil {
		return nil, err
	}
	session := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job").Model(&biz.Job{})
	if where != nil {
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
	err = session.Preload("User").Preload("FormTemplate").Find(&list).Error
	if err != nil {
		log.Error("job table query data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// UpdateJob 任务更新
func (jd *jobRepo) UpdateJob(do *biz.Job, currentStep int) (*biz.Job, error) {
	update := map[string]interface{}{
		//"reviewers":  do.Reviewers,
		"audit_time": timeFormat(do.AuditTime),
		"end_time":   timeFormat(do.EndTime),
		"start_time": timeFormat(do.StartTime),
		"status":     do.Status,
	}
	if currentStep != -1 { // 更新当前步骤
		update["current_step"] = currentStep
	}
	// 定义重试逻辑
	for attempt := 0; attempt < 3; attempt++ {
		err := jd.db.GormDB.Table("job").Where("id = ?", do.Id).Updates(update).Error
		if err == nil {
			return do, nil
		}
		if err != nil {
			updateBytes, _ := json.Marshal(update)
			log.Error("job table update data, where: id:%v, update: %v, error(%v)", do.Id, string(updateBytes), err)
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			time.Sleep(time.Duration(attempt) * 500 * time.Millisecond) // 等待一段时间后重试
		}
	}
	return do, nil
}

// BatchUpdateJobStatus 批量更新任务状态
func (jd *jobRepo) BatchUpdateJobStatus(id []int, status string) error {
	update := map[string]interface{}{
		"status": status,
	}
	err := jd.db.GormDB.Table("job").Where("id IN ?", id).Updates(update).Error
	if err != nil {
		updateBytes, _ := json.Marshal(update)
		log.Error("job table update data, where: id:%v, update: %v, error(%v)", id, string(updateBytes), err)
		return errors.Wrap(err, "data: BatchUpdateJobStatus failed")
	}
	return nil
}

// BatchUpdateJob 批量更新任务
func (jd *jobRepo) BatchUpdateJob(id []int, status string, startAt time.Time, reviewers []string) error {
	// 更新任务状态-数据库
	reStr := strings.Join(reviewers, ",")
	update := map[string]interface{}{
		"start_time": timeFormat(startAt),
		"status":     status,
		"reviewers":  reStr,
	}
	err := jd.db.GormDB.Table("job").Where("id IN ?", id).Updates(update).Error
	if err != nil {
		return errors.Errorf("data: BatchUpdateJob failed, err: %v", err)
	}
	return nil
}

// ListCreatingJob 获取创建中的任务
func (jd *jobRepo) ListCreatingJob() ([]*biz.Job, error) {
	conditions := map[string]interface{}{
		"status": enum.ExecutableStatusCreating,
	}
	session := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job").Where(conditions).Model(&biz.Job{})
	var list []*biz.Job
	err := session.Find(&list).Error
	if err != nil {
		log.Error("job table query create_wait data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// ListCreateWaitJob 获取创建等待中的任务
func (jd *jobRepo) ListCreateWaitJob(limit int) ([]*biz.Job, error) {
	conditions := map[string]interface{}{
		"status": enum.ExecutableStatusCreateWait,
	}
	session := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job").Where(conditions).Model(&biz.Job{})
	if limit != 0 {
		session.Limit(limit).Offset(0)
	}
	var list []*biz.Job
	err := session.Find(&list).Error
	if err != nil {
		log.Error("job table query create_wait data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// ListIdleJob 获取就绪中的任务
func (jd *jobRepo) ListIdleJob(limit int) ([]*biz.Job, error) {
	conditions := map[string]interface{}{
		"status": enum.ExecutableJobStateIdle,
	}
	session := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job").Where(conditions).Model(&biz.Job{})
	if limit != 0 {
		session.Limit(limit).Offset(0)
	}
	var list []*biz.Job
	err := session.Find(&list).Error
	if err != nil {
		log.Error("job table query idle data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// ListRunningJob 获取运行中的任务
func (jd *jobRepo) ListRunningJob() ([]*biz.Job, error) {
	conditions := map[string]interface{}{
		"status": enum.ExecutableJobStateRunning,
	}
	session := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job").Where(conditions).Model(&biz.Job{})
	var list []*biz.Job
	err := session.Find(&list).Error
	if err != nil {
		log.Error("job table query status running data error (%v)", err)
		return nil, err
	}
	return list, nil
}

// GetJobByJobId 通过任务ID获取任务
func (jd *jobRepo) GetJobByJobId(id int) (*biz.Job, error) {
	conditions := map[string]interface{}{
		"job.id": id,
	}
	session := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job").Where(conditions).Model(&biz.Job{})
	var list []*biz.Job
	err := session.Joins("User", jd.db.GormDB.Select("id", "name")).Joins("FormTemplate", jd.db.GormDB.Select("id", "name", "title")).Find(&list).Error
	if err != nil {
		log.Error("job table query data error (%v)", err)
		return nil, err
	}
	if len(list) < 1 {
		return nil, errors.New("job data not found")
	}
	return list[0], nil
}

// CountObjectsStatus 通过任务ID获取各对象状态数量
func (jd *jobRepo) CountObjectsStatus(jobId int) (success, failure, idle int, err error) {
	stepConditions := map[string]interface{}{
		"job_id": jobId,
	}
	session := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job_step").Where(stepConditions).Model(&biz.Step{})
	var (
		steps   []*biz.Step
		objects []*biz.Object
		stepIds []int
	)
	err = session.Select("id").Find(&steps).Error
	if err != nil {
		log.Error("CountObjectsStatus job_step table query data error (%v)", err)
		return success, failure, idle, err
	}
	if len(steps) < 1 {
		return success, failure, idle, errors.New("job data not found")
	}
	for _, v := range steps {
		stepIds = append(stepIds, v.Id)
	}
	err = jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job_object").Where("job_step_id in ?", stepIds).
		Model(&biz.Object{}).Select("id, status").Find(&objects).Error
	if err != nil {
		log.Error("CountObjectsStatus job_object table query data error (%v)", err)
		return success, failure, idle, err
	}
	if len(objects) < 1 {
		return success, failure, idle, errors.New("job data not found")
	}
	for _, v := range objects {
		if v.Status == enum.ExecutableStateSuccess {
			success++
		} else if v.Status == enum.ExecutableJobStateIdle {
			idle++
		} else if v.Status == enum.ExecutableStateFailure {
			failure++
		}
	}
	return success, failure, idle, nil
}

// GetJobObjectCountByStepId 获取步骤中对象数量
func (jd *jobRepo) GetJobObjectCountByStepId(id int) (cnt int, err error) {
	conditions := map[string]interface{}{
		"id": id,
	}
	var list []*biz.Object
	session := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job_object").Where(conditions).Model(&biz.Object{})
	err = session.Select("id").Find(&list).Error
	if err != nil {
		log.Error("CountObjectsStatus object table query data error (%v)", err)
		return cnt, err
	}
	return len(list), nil
}

// CountJobByJobStatus 通过状态获取任务数量
func (jd *jobRepo) CountJobByJobStatus(ctx context.Context, status string) (int, error) {
	return jd.countJob(&biz.JobWhere{Conditions: map[string]interface{}{"status": status}})
}

// countJob 获取任务数量
func (jd *jobRepo) countJob(where *biz.JobWhere) (int, error) {
	session := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job").Where(where.Conditions).Model(&biz.Job{})
	var total int64
	err := session.Count(&total).Error
	if err != nil {
		log.Error("job table query data count error (%v)", err)
		return 0, err
	}
	return int(total), nil
}

// UpdateStep 更新步骤
func (jd *jobRepo) UpdateStep(do *biz.Step) (*biz.Step, error) {
	// 更新步骤状态-数据库
	update := map[string]interface{}{
		"end_time":   timeFormat(do.EndTime),
		"start_time": timeFormat(do.StartTime),
		"status":     do.Status,
	}
	err := jd.db.GormDB.Table("job_step").Where("id = ?", do.Id).Updates(update).Error
	if err != nil {
		return &biz.Step{}, errors.Errorf("data: UpdateStep failed, err: %v", err)
	}
	return &biz.Step{
		Id:        do.Id,
		Status:    do.Status,
		StartTime: timeParse(timeFormat(do.StartTime)),
		EndTime:   timeParse(timeFormat(do.EndTime)),
	}, nil
}

// ListStepByJobId 获取任务步骤
func (jd *jobRepo) ListStepByJobId(jobId int) ([]*biz.Step, error) {
	var list []*biz.Step
	conditions := map[string]interface{}{"job_id": jobId}
	session := jd.db.GormDB.Table("job_step").Where(conditions).Model(&biz.Step{})
	err := session.Find(&list).Error
	if err != nil {
		log.Error("job_step table query data error (%v)", err)
		return list, err
	}
	return list, nil
}

// GetIndexStepByJobId 查看当前步骤
func (jd *jobRepo) GetIndexStepByJobId(index, jobId int) (*biz.Step, error) {
	conditions := map[string]interface{}{
		"job_id": jobId,
	}
	session := jd.db.GormDB.Table("job_step").Where(conditions).Model(&biz.Step{})
	session.Limit(1).Offset(index)
	var list []*biz.Step
	err := session.Find(&list).Error
	if err != nil {
		log.Error("job_step table query data error (%v)", err)
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("data: GetIndexStepByJobId failed")
	}
	return list[0], nil
}

// GetStepByStepId 根据 ID 查询步骤
func (jd *jobRepo) GetStepByStepId(stepId int) (*biz.Step, error) {
	conditions := map[string]interface{}{
		"id": stepId,
	}
	session := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job_step").Where(conditions).Model(&biz.Step{})
	var list []*biz.Step
	err := session.Find(&list).Error
	if err != nil {
		log.Error("job_step table query data error (%v)", err)
		return nil, err
	}
	if len(list) < 1 {
		return nil, errors.New("job_step data not found")
	}
	return list[0], nil
}

// CountStepByJobIdAndStepStatus 获取步骤数量
func (jd *jobRepo) CountStepByJobIdAndStepStatus(id int, status string) (int, error) {
	return jd.countStep(&biz.JobWhere{Conditions: map[string]interface{}{
		"status": status,
		"job_id": id,
	}})
}

// countStep 获取步骤数量
func (jd *jobRepo) countStep(where *biz.JobWhere) (int, error) {
	var total int64
	err := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job_step").Model(&biz.Step{}).Where(where.Conditions).Count(&total).Error
	if err != nil {
		log.Error("job_step table query data count error (%v)", err)
		return 0, err
	}
	return int(total), nil
}

// UpdateObject 更新对象
func (jd *jobRepo) UpdateObject(do *biz.Object) (*biz.Object, error) {
	// 更新对象状态-图数据库
	update := map[string]interface{}{
		"end_time":   timeFormat(do.EndTime),
		"start_time": timeFormat(do.StartTime),
		"status":     do.Status,
	}
	err := jd.db.GormDB.Table("job_object").Where("id = ?", do.Id).Updates(update).Error
	if err != nil {
		return &biz.Object{}, errors.Errorf("data: UpdateObject failed, err: %v", err)
	}
	return &biz.Object{
		Id:        do.Id,
		Code:      do.Code,
		Status:    do.Status,
		StartTime: timeParse(timeFormat(do.StartTime)),
		EndTime:   timeParse(timeFormat(do.EndTime)),
	}, nil
}

// UpdateObjectReservationTime 更新对象保留时间
func (jd *jobRepo) UpdateObjectReservationTime(do *biz.Object) (*biz.Object, error) {
	// 更新对象状态-图数据库
	update := map[string]interface{}{
		"reservation_period": do.ReservationPeriod,
	}
	err := jd.db.GormDB.Table("job_object").Where("id = ?", do.Id).Updates(update).Error
	if err != nil {
		return &biz.Object{}, errors.Errorf("data: UpdateObjectReservationTime failed, err: %v", err)
	}
	return &biz.Object{
		Id:                do.Id,
		Code:              do.Code,
		ReservationPeriod: do.ReservationPeriod,
	}, nil
}

// BatchUpdateObjectStatus 批量更新对象状态
func (jd *jobRepo) BatchUpdateObjectStatus(stepId int, objectId []int, status string, startTime time.Time) error {
	// 更新对象状态-数据库
	update := map[string]interface{}{
		"start_time": timeFormat(startTime),
		"status":     status,
	}
	err := jd.db.GormDB.Table("job_object").Where("id = ?", objectId).Updates(update).Error
	if err != nil {
		return errors.Errorf("data: BatchUpdateObjectStatus failed, err: %v", err)
	}
	return nil
}

// ListPageObjectByStepId 获取步骤下的对象
func (jd *jobRepo) ListPageObjectByStepId(id int, pageNum int, pageSize int) (objects []*biz.Object, err error) {
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}
	conditions := map[string]interface{}{
		"job_step_id": id,
	}
	session := jd.db.GormDB.Table("job_object").Where(conditions).Model(&biz.Object{})
	session.Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("create_time asc")
	err = session.Find(&objects).Error
	if err != nil {
		log.Error("ListPageObjectByStepId job_object table query data error (%v)", err)
		return nil, errors.Errorf("ListPageObjectByStepId: get object failed, err(%v)", err)
	}
	return objects, nil
}

// ListRunningObjectByStepId 获取步骤下的运行中对象
func (jd *jobRepo) ListRunningObjectByStepId(id int) ([]*biz.Object, error) {
	conditions := map[string]interface{}{
		"job_step_id": id,
		"Status":      enum.ExecutableJobStateRunning,
	}
	session := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job_object").Where(conditions).Model(&biz.Object{})
	var list []*biz.Object
	err := session.Find(&list).Error
	if err != nil {
		log.Error("job_object table query status running data error(%v)", err)
		return nil, err
	}
	return list, nil
}

// ListIdleObjectByStepId 获取步骤下的就绪对象
func (jd *jobRepo) ListIdleObjectByStepId(id int, limit, offset int) ([]*biz.Object, error) {
	conditions := map[string]interface{}{
		"job_step_id": id,
		"Status":      enum.ExecutableJobStateIdle,
	}
	session := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job_object").Model(&biz.Object{})
	for key, value := range conditions {
		session.Where(key, value)
	}
	if limit != 0 {
		session.Limit(limit)
	}
	session.Offset(offset)
	var list []*biz.Object
	err := session.Find(&list).Error
	if err != nil {
		log.Error("job table query status running data error(%v)", err)
		return nil, err
	}
	return list, nil
}

// GetObjectByObjectId 通过对象ID获取对象
func (jd *jobRepo) GetObjectByObjectId(id int) (*biz.Object, error) {
	conditions := map[string]interface{}{
		"id": id,
	}
	session := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job_object").Model(&biz.Object{})
	for key, value := range conditions {
		session.Where(key, value)
	}
	var data *biz.Object
	err := session.Find(&data).Error
	if err != nil {
		log.Error("object table query data by object id error(%v)", err)
		return nil, err
	}
	return data, nil
}

// CountObjectByStepIdAndObjectStatus 统计步骤下的对象数量
func (jd *jobRepo) CountObjectByStepIdAndObjectStatus(id int, status string) (int, error) {
	return jd.countObject(&biz.JobWhere{Conditions: map[string]interface{}{
		"job_step_id": id,
		"Status":      status,
	}})
}

// CountObjectByStepId 获取对象数量
func (jd *jobRepo) countObject(where *biz.JobWhere) (int, error) {
	session := jd.db.GormDB.Session(&gorm.Session{NewDB: true}).Table("job_object").Where(where.Conditions).Model(&biz.Object{})
	var total int64
	err := session.Count(&total).Error
	if err != nil {
		log.Error("job_object table query data count error (%v)", err)
		return 0, err
	}
	return int(total), nil
}
