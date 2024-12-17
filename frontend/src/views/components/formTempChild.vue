<template>
  <div class="dashboard-container">
    <div v-if="childData" class="diabox33">
      <template v-if="ifStyleTable == 'table'">
        <el-button @click="clickServerBatch(childObject[0], childData)">选择服务器</el-button>
        <el-table border class="specialTable" v-if="childData[0]"
          :data="childObject ? JSON.parse(JSON.stringify(childObject)) : []" style="width: 100%;margin-top:10px"
          size="mini" :show-header="true" max-height="300" min-height="100">
          <el-table-column :show-overflow-tooltip="true" :label="ttitem.paramTitle" :prop="ttitem.paramName"
            v-for="(ttitem, tti) in childData[0]" :key="tti">
            <template slot="header" slot-scope="scopeds">
              {{ ttitem.paramTitle }}
              <el-popover placement="right" width="330" trigger="click">
                <div v-if="scopeds">
                  <el-input style="max-height:400px;overflow-y:auto" type="textarea" clearable
                    :autosize="{ minRows: 3 }" resize="none" v-model="batchSearchInput"></el-input>
                  <div style="margin-top:10px;text-align:right;"><el-button
                      @click="changebatchSearchInput(childObject, ttitem.paramName, batchSearchInput)">批量添加</el-button>
                  </div>
                </div>
                <span v-if="ttitem.paramType != 'object'" style="margin-right:10px;color:rgb(64, 158, 255);"
                  class="el-icon-plus" @click="batchSearchInput = ''" slot="reference"></span>
              </el-popover>
            </template>
            <template slot-scope="scopeds">
              <template v-if="ttitem.paramType == 'object'">
                <el-select v-if="ttitem.paramStyle == 'search' || ttitem.paramreference == 'User'"
                  :popper-class="widthPro == '40%' ? 'specal_select_width' : 'specal_slect_classs'" filterable clearable
                  style="flex-grow:1;"
                  @change="changeChildSelect($event, tti, childItem3.paramName, childObject, childData, scopeds.$index)"
                  :multiple="childData[scopeds.$index].paramList"
                  v-model.trim="childObject[scopeds.$index][ttitem.paramName]"
                  remote  :remote-method="value => filterMethod(value, ttitem, childObject[scopeds.$index], scopeds.$index, tti)"
                  placeholder="请选择" :class="childData[scopeds.$index][tti].colorRed === '5' ? 'colorRedss' : ''">
                  <el-option v-for="(item, index) in childData[scopeds.$index][tti].paramDataList"
                    :key="index + 'i' + item.value" :label="item.name" :value="item.value">
                    <span style="display:block;font-size:15px;font-weight:blod">{{ item.name }}
                    </span>
                    <span style="color: #8492a6; font-size: 13px">{{ item.content }}</span>
                  </el-option>
                </el-select>
                <el-select v-else :popper-class="widthPro == '40%' ? 'specal_select_width' : 'specal_slect_classs'"
                  filterable clearable style="flex-grow:1;"
                  @change="changeChildSelect($event, tti, childItem3.paramName, childObject, childData, scopeds.$index)"
                  @focus="focusChildSelect(scopeds.$index, tti, childObject, ttitem, ttitem.paramName)"
                  :multiple="childData[scopeds.$index].paramList"
                  v-model.trim="childObject[scopeds.$index][ttitem.paramName]" placeholder="请选择"
                  :class="childData[scopeds.$index][tti].colorRed === '5' ? 'colorRedss' : ''">
                  <el-option v-for="(item, index) in childData[scopeds.$index][tti].paramDataList"
                    :key="index + 'ii' + item.value" :label="item.name" :value="item.value">
                    <span style="display:block;font-size:15px;font-weight:blod">{{ item.name }}
                    </span>
                    <span style="color: #8492a6; font-size: 13px">{{ item.content }}</span>
                  </el-option>
                </el-select>
              </template>
              <template v-else>
                <div v-if="ttitem.paramEnum">
                  <el-select filterable clearable style="width:100%;" key-value="index"
                    :multiple="childData[scopeds.$index].paramList"
                    v-model.trim="childObject[scopeds.$index][ttitem.paramName]" placeholder="请选择"
                    @change="tureOrFalseCoant($event, scopeds.$index, childObject[scopeds.$index][ttitem.paramName], ttitem, tti)"
                    :class="childData[scopeds.$index][tti].colorRed === '5' ? 'colorRedss' : ''">
                    <el-option v-for="(item, index) in childData[scopeds.$index][tti].paramDataList"
                      :key="index + 'iii' + item.value" :label="item.name" :value="item.value">
                    </el-option>
                  </el-select>
                </div>
                <div v-else>
                  <div v-if="ttitem.paramType == 'string'">
                    <el-input v-if="ttitem.paramStyle == 'password'" clearable style="width:100%;" placeholder="请输入"
                      v-model.trim="childObject[scopeds.$index][ttitem.paramName]" show-password type="text"
                      :maxlength="ttitem.paramMaxLength"
                      @change="validatingchild($event, scopeds.$index, tti, ttitem.paramMinLength, ttitem.paramMaxLength, ttitem.paramPatterns, 'string')"
                      :class="childData[scopeds.$index][tti].colorRed === '2' || childData[scopeds.$index][tti].colorRed === '3' || childData[scopeds.$index][tti].colorRed === '4' || childData[scopeds.$index][tti].colorRed === '5' ? 'colorReds' : ''"
                      show-word-limit></el-input>
                    <el-input v-if="ttitem.paramStyle == 'textarea' || ttitem.paramStyle == 'text'"
                      style="width:100%;margin-top:5px" clearable placeholder="请输入"
                      v-model="childObject[scopeds.$index][ttitem.paramName]" type="textarea"
                      :maxlength="ttitem.paramMaxLength"
                      @change="validatingchild($event, scopeds.$index, tti, ttitem.paramMinLength, ttitem.paramMaxLength, ttitem.paramPatterns, 'string')"
                      :class="childData[scopeds.$index][tti].colorRed === '2' || childData[scopeds.$index][tti].colorRed === '3' || childData[scopeds.$index][tti].colorRed === '4' || childData[scopeds.$index][tti].colorRed === '5' ? 'colorReds' : ''"
                      :autosize="{ minRows: 3, maxRows: 6 }" show-word-limit></el-input>
                    <el-input style="width:100%;border:none" disabled v-if="ttitem.paramStyle == 'description'"
                      clearable placeholder="请输入" v-model="ttitem.paramDefault" class="noneBorder"></el-input>
                    <el-input v-if="!ttitem.paramStyle || ttitem.paramStyle == ''" style="width:100%;" placeholder="请输入"
                      v-model.trim="childObject[scopeds.$index][ttitem.paramName]" type="text"
                      :maxlength="ttitem.paramMaxLength" clearable
                      @change="validatingchild($event, scopeds.$index, tti, ttitem.paramMinLength, ttitem.paramMaxLength, ttitem.paramPatterns, 'string')"
                      :class="childData[scopeds.$index][tti].colorRed === '2' || childData[scopeds.$index][tti].colorRed === '3' || childData[scopeds.$index][tti].colorRed === '4' || childData[scopeds.$index][tti].colorRed === '5' ? 'colorReds' : ''"
                      show-word-limit></el-input>
                  </div>
                  <div v-else>
                    <template v-if="ttitem.paramStyle == 'description'">
                      <el-input style="width:100%;border:none" disabled
                        v-if="ttitem.paramType == 'number' && !ttitem.paramMutipleof" placeholder="请输入" clearable
                        v-model="ttitem.paramDefault" class="noneBorder"></el-input>
                      <el-input-number style="width:100%;" disabled clearable placeholder="请输入" v-else
                        v-model.number="childObject[scopeds.$index][ttitem.paramName]" type="number"
                        @change="validatingchild($event, scopeds.$index, tti, ttitem.paramMinmum, ttitem.paramMaxmum, ttitem.paramPatterns, 'number')"
                        :class="childData[scopeds.$index][tti].colorRed === '2' || childData[scopeds.$index][tti].colorRed === '3' || childData[scopeds.$index][tti].colorRed === '4' || childData[scopeds.$index][tti].colorRed === '5' ? 'colorReds noneBorder' : 'noneBorder'"
                        :step="ttitem.paramMutipleof" step-strictly></el-input-number>
                    </template>
                    <template v-else>
                      <el-input style="width:100%;border:none" clearable
                        v-if="ttitem.paramType == 'number' && !ttitem.paramMutipleof" placeholder="请输入"
                        v-model.number="childObject[scopeds.$index][ttitem.paramName]"></el-input>
                      <el-input-number style="width:100%;" clearable placeholder="请输入" v-else
                        v-model.number="childObject[scopeds.$index][ttitem.paramName]" type="number"
                        @change="validatingchild($event, scopeds.$index, tti, ttitem.paramMinmum, ttitem.paramMaxmum, ttitem.paramPatterns, 'number')"
                        :class="childData[scopeds.$index][tti].colorRed === '2' || childData[scopeds.$index][tti].colorRed === '3' || childData[scopeds.$index][tti].colorRed === '4' || childData[scopeds.$index][tti].colorRed === '5' ? 'colorReds' : ''"
                        :step="ttitem.paramMutipleof" step-strictly></el-input-number>
                    </template>
                  </div>
                </div>
              </template>
              <div
                v-if="ttitem.paramType == 'string' && (ttitem.paramStyle == 'textarea' || ttitem.paramStyle == 'text')">
                <div v-if="ttitem.paramDescription && childData[scopeds.$index][tti].colorRed === '1'"
                  style="position: relative;font-size: 12px;line-height: 10px;bottom: -6px;color: #c1c1c1;"
                  v-html="ttitem.paramDescription"></div>
                <div
                  v-if="childData[scopeds.$index][tti].colorRed === '2' || childData[scopeds.$index][tti].colorRed === '3' || childData[scopeds.$index][tti].colorRed === '4' || childData[scopeds.$index][tti].colorRed === '5'"
                  style="position: relative;font-size: 12px;line-height: 10px;bottom: -6px;color: red;">
                  {{ ttitem.paramConstraintDescription ? ttitem.paramConstraintDescription : ttitem.customDescription }}
                </div>
              </div>
              <div v-else>
                <div v-if="ttitem.paramDescription && childData[scopeds.$index][tti].colorRed === '1'"
                  style="position: relative;font-size: 12px;line-height: 12px;bottom:-2px;color: #c1c1c1;"
                  v-html="ttitem.paramDescription"></div>
                <div
                  v-if="childData[scopeds.$index][tti].colorRed === '2' || childData[scopeds.$index][tti].colorRed === '3' || childData[scopeds.$index][tti].colorRed === '4' || childData[scopeds.$index][tti].colorRed === '5'"
                  style="position: relative;font-size: 12px;line-height: 12px;bottom:-2px;color: red;">
                  {{ ttitem.paramConstraintDescription ? ttitem.paramConstraintDescription : ttitem.customDescription }}
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column fixed="right" label="操作" width="60" align="center">
            <template slot-scope="scopeds">
              <i style="cursor: pointer;" v-if="childData.length != 1" class="el-icon-delete"
                @click="deletCardAttr(childData[scopeds.$index], '', scopeds.$index)"></i>
            </template>
          </el-table-column>
        </el-table>
      </template>
      <template v-else>
        <el-card v-for="(childItem3, childii3) in childData" :key="'c1' + childii3"
          style="margin-top:10px;position:relative">
          <el-form v-if="childObject[childii3]" :model="childObject[childii3]" ref="childObject" label-width="100px"
            class="demo-ruleForm">
            <span style="position:absolute;right:15px;top:0px;font-size:15px;cursor: pointer;z-index:99;">
              <i class="el-icon-document-copy" title="复制" style="margin-right: 10px;"
                @click="copyCardAttr(childItem3, childItem3.paramName, childii3)"></i>
              <i class="el-icon-delete" title="删除"
                @click="deletCardAttr(childItem3, childItem3.paramName, childii3)"></i>
            </span>
            <div style="font-size: 13px;font-weight: 900;cursor: pointer;">
              <span @click="checkShow(childItem3, childii3)"
                style="position: absolute;top: 0;left: 0px;display: block;width: 100%;padding-left: 15px;">
                <span>
                  <span
                    :class="(childObject[childii3] && childObject[childii3]['_isShowStatus']) ? 'el-icon-arrow-right' : 'el-icon-arrow-down'">
                    <span style="font-size:13px;margin-left:10px;position: relative;">
                      批次{{ aiii.erindex + 1 }}.{{ childii3 + 1 }} - {{ childItem3[0].paramFromTitle }}
                    </span>
                    <span v-for="(ttitem, tti) in childItem3[0].showTitleArr" :key="tti">
                      <span v-if="childObject[childii3] && idChangeNameObj[childObject[childii3][ttitem]]">
                        {{ idChangeNameObj[childObject[childii3][ttitem]] }}
                      </span>
                      <span v-else>
                        {{ childObject[childii3][ttitem] }}
                      </span>
                    </span>
                  </span>
                </span>
              </span>
            </div>
            <div v-show="childObject[childii3] && childObject[childii3]['_isShowStatus']">
              <div v-for="(Item33, itemII3) in childItem3" :key="itemII3" style="margin-top:5px">
                <el-form-item v-show="Item33.thisShowIf == '1'" :label="Item33.paramTitle"
                  :rules="Item33.paramRequired ? [{ required: true, message: '请选择', trigger: 'change' }, { required: true, pattern: /\S/, message: '必填不能为空', trigger: 'change' }] : []">
                  <div v-if="Item33.paramType == 'object'">
                    <el-button
                      v-if="(Item33.paramList && (Item33.paramreferenceShow.length > 1 || Item33.paramreferenceDisplay) && Item33.paramreference != 'FormTemplate' && Item33.paramStyle != 'search') || (Item33.paramList && Item33.paramStyle == 'tableStyle' && Item33.paramreference != 'FormTemplate')"
                      @click="choseServerList(itemII3, childObject[childii3], childData, Item33, Item33.paramName, childii3, Item33.paramTitle)">选择{{
                        Item33.paramTitle }}</el-button>
                    <template v-else>
                      <div style="display:flex;width:100%;">
                        <div style="flex-grow:1;" v-if="Item33.paramStyle == 'radio'">
                          <div class="getQLBox">
                            <template v-if="!Item33.paramList">
                              <el-radio-group v-model.trim="childObject[childii3][Item33.paramName]" size="mini"
                                @change="tureOrFalseCoant($event, childii3, childObject[childii3][Item33.paramName], Item33, itemII3)"
                                :class="Item33.colorRed === '5' ? 'colorRedss' : ''">
                                <el-radio-button v-for="(eeeitem, eeeindex) in Item33.paramDataList" :key="eeeindex"
                                  :label="eeeitem.value">
                                  <template>
                                    <el-tooltip class="item" effect="dark" placement="top">
                                      <div slot="content">
                                        <div style="text-align:center;">
                                          {{ eeeitem.name }}</div>
                                        {{ eeeitem.content }}
                                      </div>
                                      <div>
                                        {{ eeeitem.name }}
                                      </div>
                                    </el-tooltip>
                                  </template>
                                </el-radio-button>
                              </el-radio-group>
                            </template>
                            <el-checkbox-group v-else v-model.trim="childObject[childii3][Item33.paramName]" size="mini"
                              @change="tureOrFalseCoant($event, childii3, childObject[childii3][Item33.paramName], Item33, itemII3)"
                              :class="Item33.colorRed === '5' ? 'colorRedss' : ''">
                              <el-checkbox-button v-for="(eeeitem, eeeindex) in Item33.paramDataList" :key="eeeindex"
                                :label="eeeitem.value">
                                <template>
                                  <el-tooltip class="item" effect="dark" placement="top">
                                    <div slot="content">
                                      <div style="text-align:center;">{{ eeeitem.name }}
                                      </div>
                                      {{ eeeitem.content }}
                                    </div>
                                    <div>
                                      {{ eeeitem.name }}
                                    </div>
                                  </el-tooltip>
                                </template>
                              </el-checkbox-button>
                            </el-checkbox-group>
                          </div>
                        </div>
                        <template v-else>
                          <el-select v-if="Item33.paramStyle == 'search' || Item33.paramreference == 'User'"
                            popper-class="specal_slect_classs" filterable clearable style="flex-grow:1;"
                            @focus="getMetaSchemasListvalFocus('', Item33, childObject[childii3], childii3, itemII3)"
                            @change="changeChildSelect($event, itemII3, childItem3.paramName, childObject, childData, childii3)"
                            :multiple="Item33.paramList" v-model.trim="childObject[childii3][Item33.paramName]"
                            remote  :remote-method="value => filterMethod(value, Item33, childObject[childii3], childii3, itemII3)"
                            placeholder="请选择" :class="Item33.colorRed === '5' ? 'colorRedss' : ''">
                            <el-option v-for="(item, index) in Item33.paramDataList" :key="index + '4i' + item.name"
                              :label="item.name" :value="item.value">
                              <span style="display:block;font-size:15px;font-weight:blod">{{
                                item.name }}
                              </span>
                              <span style="color: #8492a6; font-size: 13px">{{ item.content
                              }}</span>
                            </el-option>
                          </el-select>
                          <el-select v-if="Item33.paramStyle != 'search' && Item33.paramreference != 'User'"
                            popper-class="specal_slect_classs" filterable clearable style="flex-grow:1;"
                            @change="changeChildSelect($event, itemII3, childItem3.paramName, childObject, childData, childii3)"
                            @focus="focusChildSelect(childii3, itemII3, childObject, Item33, Item33.paramName)"
                            :multiple="Item33.paramList" v-model.trim="childObject[childii3][Item33.paramName]"
                            placeholder="请选择" :class="Item33.colorRed === '5' ? 'colorRedss' : ''">
                            <el-option v-for="(item, index) in Item33.paramDataList" :key="index + '5i' + item.value"
                              :label="item.name" :value="item.value">
                              <span style="display:block;font-size:15px;font-weight:blod">{{
                                item.name }}
                              </span>
                              <span style="color: #8492a6; font-size: 13px">{{ item.content
                              }}</span>
                            </el-option>
                          </el-select>
                        </template>
                        <span v-if="Item33.paramStyle == 'radio'"
                          style="width: 30px;height: 30px;line-height: 30px;font-size: 16px;cursor: pointer;position: absolute;right: 15px;top: 7px;"
                          class="el-icon-refresh"
                          @click="focusChildSelect(childii3, itemII3, childObject, Item33, Item33.paramName)"></span>
                      </div>
                    </template>
                    <el-table border class="specialTable"
                      v-if="Item33.paramList && (Item33.paramreferenceShow.length > 1 || Item33.paramreferenceDisplay) && Item33.paramreference != 'FormTemplate' && Item33.paramStyle != 'search' && (Item33.paramTableTitle && Item33.paramTableTitle.length)"
                      :data="childObject[childii3][Item33.paramName]"
                      style="width: 100%;margin-top:10px;max-width:700px;" size="mini" :show-header="true"
                      max-height="300" min-height="100">
                      <template v-for="(taItems, i) in Item33.paramTableTitle">
                        <el-table-column :key="i" :show-overflow-tooltip="true" :label="taItems.title"
                          :prop="taItems.value"
                          v-if="Item33.paramreferenceLink && (taItems.title == '名称' || taItems.value == 'name')">
                          <template slot-scope="scoped">
                            <el-link type="primary" style="font-size:12px;" @click="linkJumpTo(Item33, scoped.row)">{{
                              scoped.row[taItems.value] }}</el-link>
                          </template>
                        </el-table-column>
                        <el-table-column :key="i" :show-overflow-tooltip="true" :label="taItems.title"
                          :prop="taItems.value" v-else>
                        </el-table-column>
                      </template>
                      <el-table-column fixed="right" v-if="Item33.paramTableTitle" label="操作" width="60" align="center">
                        <template slot-scope="scoped">
                          <i style="cursor: pointer;" class="el-icon-delete"
                            @click="deloneline(childii3, Item33.paramName, itemII3, scoped)"></i>
                        </template>
                      </el-table-column>
                    </el-table>
                    <span style="color: #666;font-size:13px;"
                      v-if="Item33.paramList && (Item33.paramreferenceShow.length > 1 || Item33.paramreferenceDisplay) && childObject[childii3][Item33.paramName] && Item33.paramreference != 'FormTemplate' && Item33.paramStyle != 'search' && (Item33.paramTableTitle && Item33.paramTableTitle.length)">已选数量：{{
                        childObject[childii3][Item33.paramName].length }}</span>
                  </div>
                  <div v-else>
                    <div v-if="Item33.paramEnum">
                      <div v-if="Item33.paramStyle == 'radio'">
                        <template v-if="!Item33.paramList">
                          <el-radio-group v-model.trim="childObject[childii3][Item33.paramName]"
                            @change="tureOrFalseCoant($event, childii3, childObject[childii3][Item33.paramName], Item33, itemII3)"
                            :class="Item33.colorRed === '5' ? 'colorRedss' : ''" size="mini">
                            <el-radio-button v-for="(item, index) in Item33.paramDataList" :key="index"
                              :label="item.value">{{ item.name }}</el-radio-button>
                          </el-radio-group>
                        </template>
                        <el-checkbox-group v-else v-model.trim="childObject[childii3][Item33.paramName]"
                          @change="tureOrFalseCoant($event, childii3, childObject[childii3][Item33.paramName], Item33, itemII3)"
                          :class="Item33.colorRed === '5' ? 'colorRedss' : ''" size="mini">
                          <el-checkbox-button v-for="(item, index) in Item33.paramDataList" :key="index"
                            :label="item.value">{{ item.name }}</el-checkbox-button>
                        </el-checkbox-group>
                      </div>
                      <el-select v-else filterable clearable style="width:100%;" key-value="index"
                        :multiple="Item33.paramList" v-model.trim="childObject[childii3][Item33.paramName]"
                        placeholder="请选择"
                        @change="tureOrFalseCoant($event, childii3, childObject[childii3][Item33.paramName], Item33, itemII3)"
                        :class="Item33.colorRed === '5' ? 'colorRedss' : ''">
                        <el-option v-for="(item, index) in Item33.paramDataList" :key="index + '6i' + item.value"
                          :label="item.name" :value="item.value">
                        </el-option>
                      </el-select>
                    </div>
                    <div v-else>
                      <template v-if="Item33.paramprerequisite && Item33.paramprerequisite.length">
                        <el-select style="width:100%;" filterable clearable
                          @change="tureOrFalseCoant($event, childii3, childObject[childii3][Item33.paramName], Item33, itemII3)"
                          @focus="focusChildGLArr(Item33, childii3, itemII3, Item33.paramName)"
                          :multiple="Item33.paramList" key-value="value"
                          v-model.trim="childObject[childii3][Item33.paramName]" placeholder="请选择"
                          :class="Item33.colorRed === '5' ? 'colorRedss' : ''" size="mini">
                          <el-option v-for="(item, index) in Item33.paramDataList" :key="index + '7i' + item.value"
                            :label="item.name" :value="item.value">
                          </el-option>
                        </el-select>
                      </template>
                      <template v-else>
                        <div v-if="Item33.paramList">
                        </div>
                        <div v-else>
                          <template>
                            <template
                              v-if="Item33.paramType == 'string' && (Item33.paramStyle == 'datetime' || Item33.paramStyle == 'date')">
                              <el-date-picker v-model="value1" type="datetime" placeholder="选择日期时间"
                                v-model.trim="childObject[childii3][Item33.paramName]" style="width:100%"
                                @change="tureOrFalseCoant($event, childii3, childObject[childii3][Item33.paramName], Item33, itemII3)"
                                v-if="Item33.paramStyle == 'datetime'">
                              </el-date-picker>
                              <el-date-picker v-model="value1" type="date" placeholder="选择日期"
                                v-model.trim="childObject[childii3][Item33.paramName]" style="width:100%"
                                @change="tureOrFalseCoant($event, childii3, childObject[childii3][Item33.paramName], Item33, itemII3)"
                                v-if="Item33.paramStyle == 'date'" value-format="yyyy-MM-dd">
                              </el-date-picker>
                            </template>
                            <div v-else>
                              <div v-if="Item33.paramType == 'string'">
                                <el-input v-if="Item33.paramStyle == 'password'" clearable style="width:100%;"
                                  placeholder="请输入" v-model.trim="childObject[childii3][Item33.paramName]" show-password
                                  type="text" :maxlength="Item33.paramMaxLength"
                                  @change="validatingchild($event, childii3, itemII3, Item33.paramMinLength, Item33.paramMaxLength, Item33.paramPatterns, 'string')"
                                  :class="Item33.colorRed === '2' || Item33.colorRed === '3' || Item33.colorRed === '4' || Item33.colorRed === '5' ? 'colorReds' : ''"
                                  show-word-limit></el-input>
                                <el-input v-if="Item33.paramStyle == 'textarea' || Item33.paramStyle == 'text'"
                                  style="width:100%;margin-top:5px" clearable placeholder="请输入"
                                  v-model="childObject[childii3][Item33.paramName]" type="textarea"
                                  :maxlength="Item33.paramMaxLength"
                                  @change="validatingchild($event, childii3, itemII3, Item33.paramMinLength, Item33.paramMaxLength, Item33.paramPatterns, 'string')"
                                  :class="Item33.colorRed === '2' || Item33.colorRed === '3' || Item33.colorRed === '4' || Item33.colorRed === '5' ? 'colorReds' : ''"
                                  :autosize="{ minRows: 3, maxRows: 6 }" show-word-limit></el-input>
                                <el-input style="width:100%;border:none" disabled
                                  v-if="Item33.paramStyle == 'description'" clearable placeholder="请输入"
                                  v-model="Item33.paramDefault" class="noneBorder"></el-input>
                                <el-input v-if="!Item33.paramStyle || Item33.paramStyle == ''" style="width:100%;"
                                  placeholder="请输入" clearable v-model.trim="childObject[childii3][Item33.paramName]"
                                  type="text" :maxlength="Item33.paramMaxLength"
                                  @change="validatingchild($event, childii3, itemII3, Item33.paramMinLength, Item33.paramMaxLength, Item33.paramPatterns, 'string')"
                                  :class="Item33.colorRed === '2' || Item33.colorRed === '3' || Item33.colorRed === '4' || Item33.colorRed === '5' ? 'colorReds' : ''"
                                  show-word-limit></el-input>
                              </div>
                              <div v-else>
                                <template v-if="Item33.paramStyle == 'description'">
                                  <el-input style="width:100%;border:none" disabled
                                    v-if="Item33.paramType == 'number' && !Item33.paramMutipleof" clearable
                                    placeholder="请输入" v-model="Item33.paramDefault" class="noneBorder"></el-input>
                                  <el-input-number style="width:100%;" clearable disabled placeholder="请输入" v-else
                                    v-model.number="childObject[childii3][Item33.paramName]" type="number"
                                    @change="validatingchild($event, childii3, itemII3, Item33.paramMinmum, Item33.paramMaxmum, Item33.paramPatterns, 'number')"
                                    :class="Item33.colorRed === '2' || Item33.colorRed === '3' || Item33.colorRed === '4' || Item33.colorRed === '5' ? 'colorReds noneBorder' : 'noneBorder'"
                                    :step="Item33.paramMutipleof" step-strictly></el-input-number>
                                </template>
                                <template v-else>
                                  <el-input style="width:100%;border:none" clearable
                                    v-if="Item33.paramType == 'number' && !Item33.paramMutipleof" placeholder="请输入"
                                    v-model.number="childObject[childii3][Item33.paramName]"
                                    @change="validatingchild($event, childii3, itemII3, Item33.paramMinmum, Item33.paramMaxmum, Item33.paramPatterns, 'number')"
                                    :class="Item33.colorRed === '2' || Item33.colorRed === '3' || Item33.colorRed === '4' || Item33.colorRed === '5' ? 'colorReds' : ''"></el-input>
                                  <el-input-number style="width:100%;" placeholder="请输入" clearable v-else
                                    v-model.number="childObject[childii3][Item33.paramName]" type="number"
                                    @change="validatingchild($event, childii3, itemII3, Item33.paramMinmum, Item33.paramMaxmum, Item33.paramPatterns, 'number')"
                                    :class="Item33.colorRed === '2' || Item33.colorRed === '3' || Item33.colorRed === '4' || Item33.colorRed === '5' ? 'colorReds' : ''"
                                    :step="Item33.paramMutipleof" step-strictly></el-input-number>
                                </template>
                              </div>
                            </div>
                          </template>
                        </div>
                      </template>
                    </div>
                  </div>
                  <i v-if="Item33.paramStyle == 'radio'" class="el-icon-circle-close"
                    style="width: 30px;height: 30px;line-height: 30px;font-size: 16px;cursor: pointer;position: absolute;right: 45px;top: 7px;"
                    @click="clearableRadioBtn(Item33, childii3, childObject[childii3][Item33.paramName], Item33, itemII3, childObject[childii3], Item33.paramName)"></i>
                  <div
                    v-if="Item33.paramType == 'string' && (Item33.paramStyle == 'textarea' || Item33.paramStyle == 'text')">
                    <div v-if="Item33.paramDescription && Item33.colorRed === '1'"
                      style="position: relative;font-size: 12px;line-height: 10px;bottom: -6px;color: #c1c1c1;"
                      v-html="Item33.paramDescription"></div>
                    <div
                      v-if="Item33.colorRed === '2' || Item33.colorRed === '3' || Item33.colorRed === '4' || Item33.colorRed === '5'"
                      style="position: relative;font-size: 12px;line-height: 10px;bottom: -6px;color: red;">
                      {{ Item33.paramConstraintDescription ? Item33.paramConstraintDescription : Item33.customDescription
                      }}
                    </div>
                  </div>
                  <div v-else>
                    <div v-if="Item33.paramDescription && Item33.colorRed === '1'"
                      style="position: relative;font-size: 12px;line-height: 12px;bottom:-2px;color: #c1c1c1;"
                      v-html="Item33.paramDescription"></div>
                    <div
                      v-if="Item33.colorRed === '2' || Item33.colorRed === '3' || Item33.colorRed === '4' || Item33.colorRed === '5'"
                      style="position: relative;font-size: 12px;line-height: 12px;bottom:-2px;color: red;">
                      {{ Item33.paramConstraintDescription ? Item33.paramConstraintDescription : Item33.customDescription
                      }}
                    </div>
                  </div>
                </el-form-item>
              </div>
            </div>

          </el-form>
        </el-card>
      </template>
    </div>
    <!-- 选择服务器弹框 -->
    <el-dialog :title="serverTitle" :visible.sync="choseServerLog" :close-on-click-modal="false" append-to-body
      class="diaBox11" width="60%">
      <div style="text-align:right">
        <el-link type="primary" style="font-size:12px;margin-right:10px;"
          @click="searchShowStatus = !searchShowStatus">更多搜索<i class="el-icon-d-arrow-right"></i></el-link>
        <el-input v-model="searchFormServe" @keydown.enter.native="searchTableServer('search')"
          style="width:220px" placeholder="请输入关键字(长度大于2)" clearable></el-input>
      </div>

      <el-card v-show="searchShowStatus" class="searchBox" style="margin:10px;">
        <el-form :model="searchObj" ref="searchObjForm" label-width="80px" class="demo-ruleForm">
          <el-row style="margin-bottom:10px">
            <el-form-item :label="serverTitle">
              <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 6 }" resize="none" v-model="batchSearchName"
                placeholder="支持名称批量搜索（换行输入）"></el-input>
            </el-form-item>
          </el-row>
        </el-form>
        <div style="text-align:right;">
          <el-button @click="batchSearchName = ''">重 置</el-button>
          <el-button type="primary" plain @click="searchTableServer('search')">搜
            索</el-button>
        </div>
      </el-card>
      <el-table border :data="searchHostTable" style="width: 100%;margin-top:10px" size="mini" class="RSHostTable"
        max-height="300" min-height="100" @selection-change="handleSelectionChangeSearth">
        <el-table-column type="selection" width="40"></el-table-column>
        <el-table-column v-for="(taItem, i) in referenceShowArrHeader" :key="i" :show-overflow-tooltip="true"
          :label="taItem.title" :prop="taItem.value">
        </el-table-column>
      </el-table>
      <div style="text-align: center; margin-top: 15px; margin-bottom: 15px;">
        <el-pagination :pager-count="5" @size-change="handleSizeChange" @current-change="changehandleCurrent"
          :page-sizes="[10, 20, 30, 50, 100, 150, 200]" layout="total,sizes, prev, pager, next"
          :current-page="searchpageNum" :page-size="searchpagesize" :total="searchtotalLen">
        </el-pagination>
      </div>
      <div style="text-align:right;padding-top:10px;">
        <el-button @click="choseServerLog = false">取 消</el-button>
        <el-button type="primary" @click="choseServerBtn" v-if="!serverLogStatus">确 定</el-button>
        <el-button type="primary" @click="choseServerBtnBathjoin" v-if="serverLogStatus">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { formTempChild } from '@/components/mixin/formTempChild'
import { mapGetters } from 'vuex'
export default {
  name: 'Dashboard',
  mixins: [formTempChild],
  props: {
    childData: {
      type: Array,
      default: []
    },
    childObject: {
      type: Array,
      default: []
    },
    aiii: {
      type: Object,
      default: {}
    },
    erinfoChildattr: {
      type: Array,
      default: []
    },
    erinfoClidObj: {
      type: Object,
      default: {}
    },
    oneinfoClidObj: {
      type: Object,
      default: {}
    },
    batchIndex: {
      type: Number
    },
    waiIndex: {
      type: String,
      default: ''
    },
    parentName: {
      type: String,
      default: ''
    },
    smallName: {
      type: String,
      default: ''
    },
    ifStyleTable: {
      type: String,
      default: ''
    },
  },
  data() {
    return {
      serverLogStatus: false,
      idChangeNameObj: {},
      mainOperandName: '',
      mainOperandvsrefrename: '',
      getstr: '',
      paramreferenceShowNew1: [],
      paramreferenceShowNew2: [],
      nodeCheckName: '',
      searchtotalLen: 0,
      serverTitle: '',
      choseServerLog: false,
      searchFormServe: '',
      batchSearchName: "",
      serverTypeIndex: "",
      childRowTitle: '',
      childRowTitle1: '',
      choseListServer: [],
      searchShowStatus: false,
      searchHostTable: [],
      referenceShowArrHeader: [],
      searchpageNum: 1,
      searchpagesize: 20,
      searchtotalLen: 0,
      searchObj: {},
      batchSearchInput: "",
      batchBG: [],
      ObjDataRowServer:{}
    }
  },
  computed: {
    ...mapGetters([
      "manageUser", "manageUserGroup"
    ])
  }
}
</script>

<style scoped>
.dashboard-container .RSHostTable .el-table__header {
  display: block !important;
}

.diaBox11 {
  z-index: 5555 !important;
}

.diabox33 /deep/ .getQLBox .el-radio-button--mini .el-radio-button__inner {
  margin: 2px 0;
  min-height: 28px;
  width: 150px;
  border-left: 1px solid #dcdfe6;
  margin-right: -1px;
  overflow: hidden;
}

.getQLBox .is-active .getQLBoxOne {
  color: #fff;
  margin-top: 3px;
  position: relative;
  z-index: 999;
}

.diabox33
  /deep/
  .getQLBox
  .el-radio-button--mini.is-active
  .el-radio-button__inner {
  border-left: none;
}
</style>