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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CloudProductCommon struct {
	// ID
	ID int `gorm:"column:id" json:"id"`
	// 云ID
	CID string `gorm:"column:cid" json:"cid"`
	// 名称
	Name string `gorm:"column:name" json:"name"`
}

type OutPutCommon struct {
	// 条数
	PageSize int
	// 页码
	PageNum int
	// 排序
	Order string
}

type QueryDataResponse struct {
	// 数据
	Data string `json:"data"`
	// 总条数
	Total int64 `json:"total"`
}

// UserValues 用户信息
type UserValues struct {
	// 用户ID
	ID int
	// 用户名
	Name string
	// 角色
	Role int
	// 角色名称
	RoleName string
	// 关联的本地项目ID
	ProjectConfigID []int
}

// GetLastUnderscoreString 获取最后一个下划线前的字符串
func GetLastUnderscoreString(input string) string {
	// 查找最后一个下划线的位置
	underscoreIndex := strings.LastIndex(input, "_")
	if underscoreIndex == -1 {
		// 如果没有下划线，返回原字符串
		return input
	}
	// 返回最后一个下划线前的字符串
	return input[:underscoreIndex]
}

// GetStringToIntSlice 字符串转int切片
func GetStringToIntSlice(input string) []int {
	nums := strings.Split(input, ",")
	intNums := make([]int, len(nums))
	for i, s := range nums {
		n, _ := strconv.Atoi(s)
		intNums[i] = n
	}
	return intNums
}

// TimeParse 字符串转时间
func timeParse(str string) (t time.Time) {
	t, _ = time.Parse("2006-01-02T15:04:05-07:00", str)
	return
}

// ParseCloudData 解析查询条件
func ParseCloudData(where json.RawMessage) (conditions map[string]interface{}, err error) {
	if where == nil {
		return nil, nil
	}
	var conditionsMap map[string]interface{}
	err = json.Unmarshal(where, &conditionsMap)
	if err != nil {
		return nil, err
	}
	conditions = make(map[string]interface{})
	for k, v := range conditionsMap {
		if strings.Contains(strings.ToLower(k), strings.ToLower("IN")) {
			//IN查询
			//获取最后一个下划线前的字符串
			field := GetLastUnderscoreString(k)
			conditions[fmt.Sprintf("%s IN ?", field)] = v
		} else if strings.Contains(strings.ToLower(k), strings.ToLower("REGEX")) {
			//模糊查询
			//获取最后一个下划线前的字符串
			field := GetLastUnderscoreString(k)
			conditions[fmt.Sprintf("%s LIKE ?", field)] = fmt.Sprintf("%%%s%%", v)
		} else if strings.Contains(strings.ToLower(k), strings.ToLower("OR")) {
			//OR查询
			conditions["or"] = queryConditionDo(v)
		} else if strings.Contains(strings.ToLower(k), strings.ToLower("AND")) {
			//AND查询
			conditions["and"] = queryAndConditionDo(v)
		} else {
			conditions[fmt.Sprintf("%s = ?", k)] = v
		}
	}
	return conditions, nil
}

// queryAndConditionDo And查询条件处理
func queryAndConditionDo(con interface{}) []map[string]interface{} {
	conditionMaps, ok := con.([]interface{})
	if !ok {
		return nil
	}
	var conditions []map[string]interface{}
	for _, val := range conditionMaps {
		do := queryConditionDo(val)
		conditions = append(conditions, do)
	}
	return conditions
}

// queryConditionDo 查询条件处理
func queryConditionDo(con interface{}) map[string]interface{} {
	conditionMaps := con.(map[string]interface{})
	queryOrConditions := make(map[string]interface{})
	for key, val := range conditionMaps {
		if mapValues, ok := val.(map[string]interface{}); ok {
			for key1, val1 := range mapValues { //连表查询例如查询条件为server_image.name = "xxx"
				if mapValues1, ok1 := val1.(map[string]interface{}); ok1 {
					for key2, val2 := range mapValues1 {
						mode1 := GetFuzzyOrPrecise(key2) //获取模糊查询、精确查询、IN查询
						if mode1 == 0 {                  //拼接表名到字段
							queryOrConditions[fmt.Sprintf("%s_%s.%s LIKE ?", FirstStrUpper(key), FirstStrUpper(key1), GetLastUnderscoreString(key2))] = fmt.Sprintf("%%%s%%", val2)
						} else if mode1 == 1 { //拼接表名到字段
							queryOrConditions[fmt.Sprintf("%s_%s.%s = ?", FirstStrUpper(key), FirstStrUpper(key1), key2)] = val2
						} else if mode1 == 2 { //IN 查询
							queryOrConditions[fmt.Sprintf("%s_%s.%s IN ?", FirstStrUpper(key), FirstStrUpper(key1), GetLastUnderscoreString(key2))] = val2
						} else if mode1 == 3 { //GT 查询
							queryOrConditions[fmt.Sprintf("%s_%s.%s > ?", FirstStrUpper(key), FirstStrUpper(key1), GetLastUnderscoreString(key2))] = val2
						} else if mode1 == 4 { //LT 查询
							queryOrConditions[fmt.Sprintf("%s_%s.%s < ?", FirstStrUpper(key), FirstStrUpper(key1), GetLastUnderscoreString(key2))] = val2
						}
					}
				} else {
					mode1 := GetFuzzyOrPrecise(key1) //获取模糊查询、精确查询、IN查询
					field := ToCamelCase(key)        //将下划线命名变为驼峰命名
					if mode1 == 0 {                  //拼接表名到字段
						queryOrConditions[fmt.Sprintf("%s.%s LIKE ?", field, GetLastUnderscoreString(key1))] = fmt.Sprintf("%%%s%%", val1)
					} else if mode1 == 1 { //拼接表名到字段
						queryOrConditions[fmt.Sprintf("%s.%s = ?", field, key1)] = val1
					} else if mode1 == 2 {
						queryOrConditions[fmt.Sprintf("%s.%s IN ?", field, GetLastUnderscoreString(key1))] = val1
					} else if mode1 == 3 {
						queryOrConditions[fmt.Sprintf("%s.%s > ?", field, GetLastUnderscoreString(key1))] = val1
					} else if mode1 == 4 {
						queryOrConditions[fmt.Sprintf("%s.%s < ?", field, GetLastUnderscoreString(key1))] = val1
					}
				}
			}
			continue
		}
		//直接查询例如查询条件为name = "xxx"
		mode := GetFuzzyOrPrecise(key) //获取模糊查询还是精确查询
		if mode == 0 {
			queryOrConditions[fmt.Sprintf("%s LIKE ?", GetLastUnderscoreString(key))] = fmt.Sprintf("%%%s%%", val)
			continue
		} else if mode == 1 {
			queryOrConditions[fmt.Sprintf("%s = ?", key)] = val
			continue
		} else if mode == 2 {
			//获取最后一个下划线前的字符串 例如id_IN 转为 id IN ?
			queryOrConditions[fmt.Sprintf("%s IN ?", GetLastUnderscoreString(key))] = val
			continue
		} else if mode == 3 {
			queryOrConditions[fmt.Sprintf("%s > ?", GetLastUnderscoreString(key))] = val
		} else if mode == 4 {
			queryOrConditions[fmt.Sprintf("%s < ?", GetLastUnderscoreString(key))] = val
		}
	}
	return queryOrConditions
}

// GetHandleConditions 处理and和or查询条件
func GetHandleConditions(conditions map[string]interface{}, name string) map[string]interface{} {
	conditionMaps := make(map[string]interface{})
	or, isOr := conditions["or"]
	and, isAnd := conditions["and"]
	if !isOr && !isAnd {
		for k, v := range conditions {
			field := k
			if name != "" {
				field = fmt.Sprintf("%s.%s", name, k) //拼接新键值
			}
			conditionMaps[field] = v
			delete(conditions, k)
		}
	}
	if or != nil && isOr {
		if mapValues, isMap := or.(map[string]interface{}); isMap {
			for k, v := range mapValues {
				field := k
				if !strings.Contains(k, ".") { //连表查询拼接表名
					if name != "" {
						field = fmt.Sprintf("%s.%s", name, k) //拼接新键值
					}
					delete(mapValues, k) //删除旧键值
				}
				mapValues[field] = v
			}
			conditionMaps["or"] = mapValues
		}
	}
	if and != nil && isAnd {
		if mapValues, isMap := and.([]map[string]interface{}); isMap {
			for k, v := range mapValues {
				for k1, v1 := range v {
					field := k1
					if !strings.Contains(k1, ".") { //连表查询拼接表名
						if name != "" {
							field = fmt.Sprintf("%s.%s", name, k1) //拼接新键值
						}
						delete(v, k1) //删除旧键值
					}
					mapValues[k][field] = v1
				}
			}
			conditionMaps["and"] = mapValues
		}
	}
	for k, v := range conditionMaps {
		conditions[k] = v
	}
	return conditions
}

// GetFuzzyOrPrecise 判断是否是模糊查询或精确查询
// 返回0是模糊查询， 1:是精确查询，2:IN查询 3:GT查询 4:LT查询
func GetFuzzyOrPrecise(name string) int {
	if strings.Contains(strings.ToLower(name), strings.ToLower("_REGEX")) {
		return 0
	} else if strings.Contains(strings.ToLower(name), strings.ToLower("_IN")) {
		return 2
	} else if strings.Contains(strings.ToLower(name), strings.ToLower("_GT")) {
		return 3
	} else if strings.Contains(strings.ToLower(name), strings.ToLower("_LT")) {
		return 4
	} else {
		return 1
	}
}

// ToCamelCase 转换下划线命名到驼峰命名
func ToCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)   // 将下划线替换为空格
	s = strings.Title(s)                   // 将字符串转换为标题格式，即首字母大写
	return strings.Replace(s, " ", "", -1) // 移除空格
}

// FirstStrUpper 首字母大写
func FirstStrUpper(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// SlicesEqual 判断切片是否相同
func SlicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// DifferenceMap 获取两个map差集
func DifferenceMap(map1, map2 map[string]interface{}) map[string]interface{} {
	diff := make(map[string]interface{})
	for k, v := range map1 {
		if _, ok := map2[k]; !ok {
			diff[k] = v
		}
	}
	return diff
}
