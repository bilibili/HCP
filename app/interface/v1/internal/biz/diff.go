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
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// DiffNew 判断数据是否存在差异
func DiffNew(before interface{}, after interface{}) (bool, error) {
	// 将切片数据转换为JSON
	beforeJson, err := json.Marshal(before)
	if err != nil {
		return false, err
	}
	afterJson, err := json.Marshal(after)
	if err != nil {
		return false, err
	}
	// 计算JSON数据的SHA256哈希值
	beforeHash := sha256.Sum256(beforeJson)
	afterHash := sha256.Sum256(afterJson)
	if fmt.Sprintf("%x", beforeHash) != fmt.Sprintf("%x", afterHash) {
		return true, nil
	}
	return false, nil
}
