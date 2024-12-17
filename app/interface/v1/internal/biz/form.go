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

// MetaSchema 表单参数
type MetaSchema struct {
	DefaultAcl  string         `json:"defaultAcl,omitempty"`
	Description string         `json:"description,omitempty"`
	Display     []string       `json:"display,omitempty"`
	Id          string         `json:"id,omitempty"`
	Interfaces  []MetaSchema   `json:"interfaces,omitempty"`
	IsInterface bool           `json:"isInterface,omitempty"`
	MutRecord   bool           `json:"mutRecord,omitempty"`
	Name        string         `json:"name,omitempty"`
	Objects     []MetaSchema   `json:"objects,omitempty"`
	Operator    []FormTemplate `json:"operator,omitempty"`
	Title       string         `json:"title,omitempty"`
	Total       int64          `json:"total,omitempty"`
	Type        string         `json:"type,omitempty"`
}

// FormTemplate 表单模版
type FormTemplate struct {
	AbstractTemplate   string          `json:"abstractTemplate,omitempty"`
	Description        string          `json:"description,omitempty"`
	Id                 string          `json:"id,omitempty"`
	InitialStepStatus  string          `json:"initialStepStatus,omitempty"`
	JobMode            string          `json:"jobMode,omitempty"`
	JobSteps           []Step          `gorm:"-" json:"jobSteps,omitempty"`
	Jobs               []Job           `gorm:"-" json:"jobs,omitempty"`
	Linker             string          `json:"linker,omitempty"`
	MainOperand        *MetaSchema     `json:"mainOperand,omitempty"`
	MainOperandLock    string          `json:"mainOperandLock,omitempty"`
	Name               string          `json:"name,omitempty"`
	ParallelStep       string          `json:"parallelStep,omitempty"`
	ParallelStepEnable bool            `json:"parallelStepEnable,omitempty"`
	Parameters         []FormParameter `gorm:"-" json:"parameters,omitempty"`
	ParentParameters   []FormParameter `gorm:"-" json:"parentParameters,omitempty"`
	Prompt             string          `json:"prompt,omitempty"`
	SequenceNumber     float64         `json:"sequenceNumber,omitempty"`
	SubmittedCallback  string          `json:"submittedCallback,omitempty"`
	SubmittingCallback []string        `gorm:"-" json:"submittingCallback,omitempty"`
	Title              string          `json:"title,omitempty"`
	Workflow           string          `json:"workflow,omitempty"`
	Operate            string          `json:"operate,omitempty"`
}

// FormParameter 表单参数
type FormParameter struct {
	Annotation            string         `json:"annotation,omitempty"`
	AvailableCondition    string         `json:"availableCondition,omitempty"`
	ConstraintDescription string         `json:"constraintDescription,omitempty"`
	Default               string         `json:"default,omitempty"`
	DependentSchema       string         `json:"dependentSchema,omitempty"`
	Description           string         `json:"description,omitempty"`
	Format                string         `json:"format,omitempty"`
	Group                 string         `json:"group,omitempty"`
	Id                    string         `json:"id,omitempty"`
	JoinTitle             string         `json:"joinTitle,omitempty"`
	Link                  string         `json:"link,omitempty"`
	List                  bool           `json:"list,omitempty"`
	MaxLength             float64        `json:"maxLength,omitempty"`
	MaxNumber             int64          `json:"maxNumber,omitempty"`
	Maximum               float64        `json:"maximum,omitempty"`
	MinLength             float64        `json:"minLength,omitempty"`
	Minimum               float64        `json:"minimum,omitempty"`
	Mutipleof             float64        `json:"mutipleof,omitempty"`
	Name                  string         `json:"name,omitempty"`
	Nullable              bool           `json:"nullable,omitempty"`
	ParentTemplate        *FormTemplate  `json:"parentTemplate,omitempty"`
	Pattern               string         `json:"pattern,omitempty"`
	Prerequisite          []string       `json:"prerequisite,omitempty"`
	Reference             string         `json:"reference,omitempty"`
	ReferenceDisplay      string         `json:"referenceDisplay,omitempty"`
	ReferenceMutation     string         `json:"referenceMutation,omitempty"`
	ReferenceQuery        string         `json:"referenceQuery,omitempty"`
	ReferenceSelect       string         `json:"referenceSelect,omitempty"`
	ReferenceShow         []string       `json:"referenceShow,omitempty"`
	Required              bool           `json:"required,omitempty"`
	SchemaName            string         `json:"schemaName,omitempty"`
	SequenceNumber        int64          `json:"sequenceNumber,omitempty"`
	Style                 string         `json:"style,omitempty"`
	Templates             []FormTemplate `json:"templates,omitempty"`
	Title                 string         `json:"title,omitempty"`
	Type                  string         `json:"type,omitempty"`
	Unit                  string         `json:"unit,omitempty"`
}
