// Package biz
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
package biz

import (
	"context"
	"encoding/json"
	"time"
)

// JobRepo 任务接口
type JobRepo interface {
	// CreateJob 创建任务
	CreateJob(job *Job) (int, error)
	// CountJob 任务计数
	CountJob(where *JobWhere) (total int64, err error)
	// QueryJob 查询任务
	QueryJob(where *JobWhere, output *JobOutput) ([]*Job, error)
	// GetTemplate 获取FormTemplate
	GetTemplate(formName string) (template *FormTemplate, err error)
	// GetTemplateContent 获取FormTemplate 内容
	GetTemplateContent(jobId int) (input map[string]interface{}, err error)
	//CreateJobSteps 填充Steps
	CreateJobSteps(jobId int, template *FormTemplate, input map[string]interface{}) error
	// UpdateJob 更新任务
	UpdateJob(job *Job, currentStep int) (*Job, error)
	// BatchUpdateJobStatus 批量更新任务状态
	BatchUpdateJobStatus(id []int, status string) error
	// BatchUpdateJob 批量更新任务
	BatchUpdateJob(id []int, status string, startAt time.Time, reviewers []string) error
	// ListRunningJob 查看正在运行任务
	ListRunningJob() ([]*Job, error)
	//ListCreatingJob 查询等待中任务列表
	ListCreatingJob() ([]*Job, error)
	//ListCreateWaitJob 查询等待中任务列表
	ListCreateWaitJob(limit int) ([]*Job, error)
	// ListIdleJob 查询等待中任务列表
	ListIdleJob(limit int) ([]*Job, error)
	// GetJobByJobId 查看任务
	GetJobByJobId(id int) (*Job, error)
	// CountObjectsStatus 获取具体对象状态信息
	CountObjectsStatus(jobId int) (success, failure, idle int, err error)
	// GetJobObjectCountByStepId  查看步骤的对象数
	GetJobObjectCountByStepId(id int) (cnt int, err error)
	// CountJobByJobStatus 状态任务的计数
	CountJobByJobStatus(ctx context.Context, status string) (int, error)
	// UpdateStep 更新步骤
	UpdateStep(step *Step) (*Step, error)
	// ListStepByJobId 步骤列表
	ListStepByJobId(id int) ([]*Step, error)
	// GetIndexStepByJobId 查看当前步骤
	GetIndexStepByJobId(index, id int) (*Step, error)
	// GetStepByStepId 根据 ID 查询步骤
	GetStepByStepId(id int) (*Step, error)
	// CountStepByJobIdAndStepStatus 状态步骤的计数
	CountStepByJobIdAndStepStatus(id int, status string) (int, error)
	// UpdateObject 更新对象
	UpdateObject(object *Object) (*Object, error)
	// UpdateObjectReservationTime 更新对象-保留时间
	UpdateObjectReservationTime(object *Object) (*Object, error)
	// BatchUpdateObjectStatus 批量更新对象状态
	BatchUpdateObjectStatus(stepId int, objectId []int, status string, startTime time.Time) (err error)
	// ListPageObjectByStepId  分页查看步骤对象列表
	ListPageObjectByStepId(id int, pageNum int, pageSize int) (objects []*Object, err error)
	// ListRunningObjectByStepId 查看步骤正在执行对象
	ListRunningObjectByStepId(id int) ([]*Object, error)
	// ListIdleObjectByStepId 查看步骤待执行对象
	ListIdleObjectByStepId(id int, limit, offset int) ([]*Object, error)
	// GetObjectByObjectId 查看对象
	GetObjectByObjectId(id int) (*Object, error)
	// CountObjectByStepIdAndObjectStatus 状态对象的计数
	CountObjectByStepIdAndObjectStatus(id int, status string) (int, error)
}

// JobWhere 任务查询条件
type JobWhere struct {
	// Query 查询语句
	Query string
	// Arg 查询参数
	Arg interface{}
	// Conditions 查询条件
	Conditions map[string]interface{}
}

type JobOutput struct {
	OutPutCommon
}

type Job struct {
	Id                 int               `gorm:"primary_key;column:id" json:"id"`
	Name               string            `gorm:"column:name" json:"name"`
	Title              string            `gorm:"column:title" json:"title"`
	Raw                string            `gorm:"column:raw" json:"raw"`
	Operate            string            `gorm:"column:operate" json:"operate"`
	StartTime          time.Time         `gorm:"column:start_time" json:"startTime"`
	EndTime            time.Time         `gorm:"column:end_time" json:"endTime"`
	AuditTime          time.Time         `gorm:"column:audit_time" json:"auditTime"`
	Status             string            `gorm:"column:status" json:"status"`
	UserId             int               `gorm:"column:user_id" json:"userId"`
	User               *User             `json:"user"`
	Reviewers          string            `gorm:"column:reviewers" json:"reviewers"`
	IsDelete           int               `gorm:"column:is_delete" json:"is_delete"`
	CreateTime         time.Time         `gorm:"column:create_time" json:"create_time"`
	UpdateTime         time.Time         `gorm:"column:update_time" json:"update_time"`
	CurrentStep        int               `gorm:"column:current_step" json:"current_step"`
	TotalStep          int               `gorm:"column:total_step" json:"total_step"`
	ParallelStepEnable bool              `gorm:"column:parallel_step_enable" json:"parallel_step_enable"`
	FormTemplateID     int               `gorm:"column:form_template_id" json:"form_template_id"`
	FormTemplate       *FormTemplateData `json:"form_template"`
	TotalObject        int               `gorm:"column:total_object" json:"total_object"`
}

func (j *Job) JobTime() float64 {
	return j.EndTime.Sub(j.StartTime).Minutes()
}

type Step struct {
	Id              int
	Name            string
	JobId           int
	Job             *Job
	SequenceNumber  int
	StepTotalObject int
	Status          string
	CreateTime      time.Time
	UpdateTime      time.Time
	StartTime       time.Time
	EndTime         time.Time
	FormTemplateID  int
	FormTemplate    *FormTemplateData
	Objects         []*Object `gorm:"-"`
}

func (s *Step) StepTime() float64 {
	return s.EndTime.Sub(s.StartTime).Minutes()
}

type Object struct {
	Id                int
	Name              string
	JobStepId         int
	Operate           string
	Code              string
	IsDelete          int
	Status            string
	Raw               string
	ReservationPeriod int
	OperandCounter    int
	StartTime         time.Time
	EndTime           time.Time
	CreateTime        time.Time
	UpdateTime        time.Time
}

func (o *Object) ObjectTime() float64 {
	return o.EndTime.Sub(o.StartTime).Minutes()
}

type TreeTemplate struct {
	Name         string
	Title        string
	FormTemplate string
	Content      string
}

type TreeStory struct {
	ProcedureName  string
	ProcedureTitle string
	ObjectID       string
	ObjectStatus   string
	Operation      string
	CurrentID      string
	CurrentName    string
	CurrentTitle   string
	CurrentStatus  string
	Arguments      string
	Logs           string
	StoryTime      string
}

type Tree struct {
	Id             string
	JobId          string
	StepId         string
	ObjectId       string
	ProcedureCode  string
	ProcedureTitle string
	OperateCode    string
	JumpUid        uint64
	TreeTemplate   string
	TreeSerial     string
	TreeBlackboard map[string]interface{}
	TreeRoot       string
	LastLeaf       string
	Errors         []string
	Logs           string
}

func (t Tree) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

type TreeOption struct {
	WriteBack bool
}

type JobUseCase struct {
	repo JobRepo
}

func NewJobUseCase(repo JobRepo) *JobUseCase {
	return &JobUseCase{repo: repo}
}

func (j *JobUseCase) CreateJob(job *Job) (int, error) {
	jobId, err := j.repo.CreateJob(job)
	return jobId, err
}

func (c *JobUseCase) CountJob(where *JobWhere) (total int64, err error) {
	return c.repo.CountJob(where)
}

func (c *JobUseCase) QueryJob(where *JobWhere, output *JobOutput) ([]*Job, error) {
	return c.repo.QueryJob(where, output)
}

func (j *JobUseCase) UpdateJob(job *Job, currentStep int) (err error) {
	_, err = j.repo.UpdateJob(job, currentStep)
	return
}

func (j *JobUseCase) BatchUpdateJobStatus(id []int, status string) (err error) {
	err = j.repo.BatchUpdateJobStatus(id, status)
	return
}

func (j *JobUseCase) BatchUpdateJob(id []int, status string, startAt time.Time, reviewers []string) (err error) {
	err = j.repo.BatchUpdateJob(id, status, startAt, reviewers)
	return
}

func (j *JobUseCase) ListCreateWaitJob(limit int) (list []*Job, err error) {
	list, err = j.repo.ListCreateWaitJob(limit)
	return
}

func (j *JobUseCase) ListIdleJob(limit int) (list []*Job, err error) {
	list, err = j.repo.ListIdleJob(limit)
	return
}

func (j *JobUseCase) ListRunningJob() (list []*Job, err error) {
	list, err = j.repo.ListRunningJob()
	return
}

func (j *JobUseCase) ListCreatingJob() (list []*Job, err error) {
	list, err = j.repo.ListCreatingJob()
	return
}

func (j *JobUseCase) GetJobByJobId(id int) (job *Job, err error) {
	job, err = j.repo.GetJobByJobId(id)
	return
}

func (j *JobUseCase) CountObjectsStatus(jobId int) (success, failure, idle int, err error) {
	return j.repo.CountObjectsStatus(jobId)
}

func (j *JobUseCase) GetJobObjectCountByStepId(id int) (cnt int, err error) {
	return j.repo.GetJobObjectCountByStepId(id)
}

func (j *JobUseCase) UpdateStep(step *Step) (err error) {
	_, err = j.repo.UpdateStep(step)
	return
}

func (j *JobUseCase) ListStepByJobId(id int) (list []*Step, err error) {
	list, err = j.repo.ListStepByJobId(id)
	return
}

func (j *JobUseCase) GetIndexStepByJobId(index, id int) (step *Step, err error) {
	step, err = j.repo.GetIndexStepByJobId(index, id)
	return
}

func (j *JobUseCase) GetStepByStepId(id int) (step *Step, err error) {
	step, err = j.repo.GetStepByStepId(id)
	return
}

func (j *JobUseCase) UpdateObject(object *Object) (err error) {
	_, err = j.repo.UpdateObject(object)
	return
}

func (j *JobUseCase) UpdateObjectReservationTime(object *Object) (err error) {
	_, err = j.repo.UpdateObjectReservationTime(object)
	return
}

func (j *JobUseCase) BatchUpdateObjectStatus(stepId int, objectId []int, status string, startTime time.Time) (err error) {
	return j.repo.BatchUpdateObjectStatus(stepId, objectId, status, startTime)
}

func (j *JobUseCase) ListPageObjectByStepId(id int, pageNum int, pageSize int) (objects []*Object, err error) {
	objects, err = j.repo.ListPageObjectByStepId(id, pageNum, pageSize)
	return
}

func (j *JobUseCase) ListRunningObjectByStepId(id int) (list []*Object, err error) {
	list, err = j.repo.ListRunningObjectByStepId(id)
	return
}

func (j *JobUseCase) ListIdleObjectByStepId(id int, limit, offset int) (list []*Object, err error) {
	list, err = j.repo.ListIdleObjectByStepId(id, limit, offset)
	return
}

func (j *JobUseCase) GetObjectByObjectId(id int) (object *Object, err error) {
	object, err = j.repo.GetObjectByObjectId(id)
	return
}

func (j *JobUseCase) CountJobByJobStatus(ctx context.Context, status string) (count int, err error) {
	count, err = j.repo.CountJobByJobStatus(ctx, status)
	return
}

func (j *JobUseCase) CountStepByJobIdAndStepStatus(id int, status string) (count int, err error) {
	count, err = j.repo.CountStepByJobIdAndStepStatus(id, status)
	return
}

func (j *JobUseCase) CountObjectByStepIdAndObjectStatus(id int, status string) (count int, err error) {
	count, err = j.repo.CountObjectByStepIdAndObjectStatus(id, status)
	return
}

func (j *JobUseCase) GetTemplate(formName string) (template *FormTemplate, err error) {
	return j.repo.GetTemplate(formName)
}

func (j *JobUseCase) GetTemplateContent(jobId int) (input map[string]interface{}, err error) {
	return j.repo.GetTemplateContent(jobId)
}

func (j *JobUseCase) CreateJobSteps(jobId int, template *FormTemplate, input map[string]interface{}) error {
	return j.repo.CreateJobSteps(jobId, template, input)
}
