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
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/bilibili/HCP/app/interface/v1/internal/biz"
	"github.com/bilibili/HCP/app/interface/v1/internal/biz/enum"
	"github.com/bilibili/HCP/common/ecode"
	common "github.com/bilibili/HCP/common/models"
	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/go-kratos/kratos/pkg/net/http/blademaster/binding"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

/*
	todo 处理前端提交的表单

示例： {"input":[{"project":{"id":"0x1124af1"},"cloudAccount":{"id":"0x43710"},"num":1,"configuration":"test","formName":""}]}
*/
// CreateJob 创建任务请求
type CreateJobReq struct {
	// 表单模版ID
	TemplateId int `json:"template_id"`
	// 任务Name
	Name string `json:"name"`
	// 任务名称
	Title string `json:"title"`
	// 提交数据
	Form json.RawMessage `json:"form"`
}

// CreateJobResp 创建任务响应
type CreateJobResp struct {
	// 任务ID
	JobId int `json:"job_id"`
}

// GetJobReq 获取任务请求
type GetJobReq struct {
	// 任务ID
	Id int `json:"id"`
	// 批次数
	StepNum int `json:"step_num"`
	// 页码
	PageNum int `json:"page_num"`
	// 条数
	PageSize int `json:"page_size"`
}

// GetJobResp 获取任务响应
type GetJobResp struct {
	// 任务ID
	JobId int `json:"job_id"`
	// 任务类型
	JobType string `json:"job_type"`
	// 任务状态
	JobStatus string `json:"job_status"`
	// 任务名称
	JobTitle string `json:"job_title"`
	// 任务创建者
	User string `json:"user"`
	// 审批人
	Reviewers []string `json:"reviewers"`
	// 创建时间
	CreateTime string `json:"create_time"`
	// 任务开始时间
	StartTime string `json:"start_time"`
	// 任务结束时间
	EndTime string `json:"end_time"`
	// 任务总对象数
	Total int `json:"total"`
	// 成功数量
	Success int `json:"success"`
	// 失败数量
	Failure int `json:"failure"`
	// 就绪数量
	Idle int `json:"idle"`
	// 批次
	Steps []*Step `json:"steps"`
}

// Object 任务对象
type Object struct {
	// 对象ID
	Id int `json:"id"`
	// 对象Code
	Code string `json:"code"`
	// 对象名称
	Name string `json:"name"`
	// 对象属性
	Attributes map[string]string `json:"attributes"`
	// 对象状态
	Status string `json:"status"`
	// 开始时间
	StartTime string `json:"start_time"`
	// 结束时间
	EndTime string `json:"end_time"`
}

// Step 任务批次
type Step struct {
	// 批次ID
	Id int `json:"id"`
	// 批次名称
	Title string `json:"title"`
	// 对象数
	Total int `json:"total"`
	// 批次状态
	Status string `json:"status"`
	// 批次开始时间
	SequenceNumber int `json:"sequence_number"`
	// 批次开始时间
	StartTime string `json:"start_time"`
	// 批次结束时间
	EndTime string `json:"end_time"`
	// 对象列表
	Objects []*Object `json:"objects"`
}

// CreateJob 创建任务
func CreateJob(c *bm.Context, w http.ResponseWriter, r *http.Request) {
	//Todo FillInput 填充模型字段, 参考 doValidAndFill
	req := &CreateJobReq{}
	if err := c.BindWith(req, binding.JSON); err != nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
		return
	}
	userInfo := r.Context().Value("userInfo").(biz.UserValues)
	if userInfo.ID == 0 {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, "用户id不能为空"))
		return
	}
	if req.TemplateId == 0 {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, "模板id不能为空"))
		return
	}
	if req.Form == nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, "表单不能为空"))
		return
	}
	bytes, _ := json.Marshal(req.Form)
	job := &biz.Job{
		UserId:         userInfo.ID,
		FormTemplateID: req.TemplateId,
		Raw:            string(bytes),
		Status:         enum.ExecutableStatusCreateWait,
	}
	resp, err := Svc.Job.CreateJob(job)
	if err != nil {
		c.JSON(nil, errors.Errorf("Job creation failed, err: %v", err))
		return
	}
	data := &CreateJobResp{
		JobId: resp,
	}
	c.JSON(data, err)
	return
}

// GetJob 获取任务-用于表单回填
func GetJob(c *bm.Context) {
	req := &GetJobReq{}
	if err := c.BindWith(req, binding.JSON); err != nil {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, err.Error()))
		return
	}
	if req.Id == 0 {
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, "任务id不能为空"))
		return
	}
	if req.StepNum == 0 {
		req.StepNum = 1
	}
	do, err := Svc.Job.GetJobByJobId(req.Id)
	if err != nil {
		log.Error("GetJob failed, err: %v", err)
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, "获取任务失败"))
		return
	}
	reply := &GetJobResp{
		JobId:      do.Id,
		JobType:    do.FormTemplate.Title,
		JobTitle:   do.Title,
		JobStatus:  enum.ExecutableStateDisplayMapping[do.Status],
		User:       do.User.Name,
		CreateTime: do.CreateTime.Format("2006-01-02 15:04:05"),
		StartTime:  do.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:    do.EndTime.Format("2006-01-02 15:04:05"),
		Total:      do.TotalObject,
	}
	reply.Success, reply.Failure, reply.Idle, err = Svc.Job.CountObjectsStatus(req.Id)
	if err != nil {
		log.Error("GetJobObjectStateCount failed, err(%v)", err)
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, "获取对象失败"))
		return
	}
	if do.Status == enum.ExecutableJobStateRunning {
		reply.EndTime = do.EndTime.Format("2006-01-02 15:--:--")
	}
	steps, err := Svc.Job.ListStepByJobId(do.Id)
	if err != nil {
		log.Error("GetJob failed, err: %v", err)
		c.JSON(nil, ecode.NewECode(common.ErrServiceCode, "获取批次失败"))
		return
	}
	seg := &errgroup.Group{}
	sc := make(chan *Step, len(steps))
	for i := range steps {
		step := steps[i]
		seg.Go(func() error {
			cnt, err := Svc.Job.GetJobObjectCountByStepId(step.Id)
			if err != nil {
				log.Error("GetJobObjectContByStepId failed, err(%v)", err)
				return err
			}
			stepDto := &Step{
				Id:             step.Id,
				Title:          step.Name,
				Total:          cnt,
				Status:         enum.ExecutableStateDisplayMapping[step.Status],
				SequenceNumber: step.SequenceNumber,
				StartTime:      step.StartTime.Format("2006-01-02 15:04:05"),
				EndTime:        step.EndTime.Format("2006-01-02 15:04:05"),
			}
			if step.Status == enum.ExecutableJobStateRunning {
				stepDto.EndTime = step.EndTime.Format("2006-01-02 15:--:--")
			}
			if step.SequenceNumber != req.StepNum {
				sc <- stepDto
				return nil
			}
			objects, err := Svc.Job.ListPageObjectByStepId(step.Id, req.PageNum, req.PageSize)
			if err != nil {
				log.Error("GetJob:ListPageObjectByStepId failed, err: %v", err)
				c.JSON(nil, ecode.NewECode(common.ErrServiceCode, "获取批次对象列表失败"))
				return err
			}
			oeg := &errgroup.Group{}
			oc := make(chan *Object, len(objects))
			for i := range objects {
				object := objects[i]
				oeg.Go(func() error {
					attr := make(map[string]string)
					attr["object_status"] = object.Status
					objectDto := &Object{
						Id:         object.Id,
						Code:       object.Code,
						Name:       object.Name,
						Attributes: attr,
						Status:     enum.ExecutableStateDisplayMapping[object.Status],
						StartTime:  object.StartTime.Format("2006-01-02 15:04:05"),
						EndTime:    object.EndTime.Format("2006-01-02 15:04:05"),
					}
					if object.Status == enum.ExecutableJobStateRunning {
						objectDto.EndTime = time.Now().Format("2006-01-02 15:--:--")
					}
					oc <- objectDto
					return nil
				})
			}
			if err := oeg.Wait(); err != nil {
				close(oc)
				return err
			}
			close(oc)
			for objectDto := range oc {
				stepDto.Objects = append(stepDto.Objects, objectDto)
			}
			sort.Slice(stepDto.Objects, func(i, j int) bool {
				return stepDto.Objects[i].Id < stepDto.Objects[j].Id
			})
			sc <- stepDto
			return nil
		})
	}
	if err := seg.Wait(); err != nil {
		c.JSON(reply, err)
		return
	}
	close(sc)
	for stepDto := range sc {
		reply.Steps = append(reply.Steps, stepDto)
	}
	sort.Slice(reply.Steps, func(i, j int) bool {
		return reply.Steps[i].SequenceNumber < reply.Steps[j].SequenceNumber
	})
	c.JSON(reply, err)
	return
}
