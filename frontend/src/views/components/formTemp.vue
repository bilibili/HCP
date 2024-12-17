<template>
  <div class='autoBox general_box'>
    <el-dialog :close-on-click-modal="false" :title="newTaskName" :visible.sync="reviewStatusNew" class="diaBox"
      :width="widthPro">
      <div>
        <el-form :model="infoReviewForm" ref="infoReviewForm" label-width="100px" class="demo-ruleForm">
          <el-table :data="infoReviewAttr" row-key="_id" class="infoTable diaBoxFormMain" align="left"
            style="padding-bottom:10px">
            <el-table-column label="属性">
              <template slot-scope="scope">
                <div v-if="scope.row.paramGroup">
                  <div>
                    <h4
                      style="margin: 10px 0px 15px; padding: 0px; color: rgb(51, 51, 51); font-size: 15px; font-weight: 600; text-indent: 20px; border-left: 4px solid rgb(35, 173, 229);line-height: 20px;">
                      {{ scope.row.paramGroup }} </h4>
                  </div>
                </div>
                <div style="display:flex;width:100%;" v-show="scope.row.thisShowIf == '1'">
                  <el-form-item style="flex-grow:1;" :label="scope.row.paramTitle"
                    :rules="scope.row.paramRequired ? [{ required: true, message: '请选择', trigger: 'change' }, { required: true, pattern: /\S/, message: '必填不能为空', trigger: 'change' }] : []">
                    <div v-if="scope.row.paramType == 'object'">
                      <template v-if="scope.row.paramStyle == 'table'">
                        <el-button
                          @click="clickServerBatch(scope, scope.row, infoReviewForm, infoReviewAttr)">选择列表</el-button>
                        <el-table border class="specialTable" v-if="scope.row.childAttrArr && scope.row.childAttrArr[0]"
                          :data="infoReviewForm[scope.row.paramName + 'ChildObj'] ? JSON.parse(JSON.stringify(infoReviewForm[scope.row.paramName + 'ChildObj'])) : []"
                          style="width: 100%;margin-top:10px" size="mini" :show-header="true" max-height="300"
                          min-height="100">
                          <el-table-column :key="tti" :show-overflow-tooltip="true" :label="ttitem.paramTitle"
                            :prop="ttitem.paramName" v-for="(ttitem, tti) in scope.row.childAttrArr[0]">
                            <template slot="header" slot-scope="scopeds">
                              {{ ttitem.paramTitle }}
                              <el-popover placement="right" width="330" trigger="click">
                                <div v-if="scopeds">
                                  <el-input style="max-height:400px;overflow-y:auto" type="textarea"
                                    :autosize="{ minRows: 3 }" resize="none" v-model="batchSearchInput"
                                    clearable></el-input>
                                  <div style="margin-top:10px;text-align:right;"><el-button
                                      @click="changebatchSearchInput(infoReviewForm[scope.row.paramName + 'ChildObj'], ttitem.paramName, batchSearchInput, scope.row.childAttrArr, tti)">批量添加</el-button>
                                  </div>
                                </div>
                                <span v-if="ttitem.paramType != 'object'"
                                  style="margin-right:10px;color:rgb(64, 158, 255);" class="el-icon-plus"
                                  @click="batchSearchInput = ''" slot="reference"></span>
                              </el-popover>
                            </template>
                            <template slot-scope="scopeds">
                              <template v-if="ttitem.paramType == 'object'">
                                <el-select v-if="ttitem.paramStyle == 'search' || ttitem.paramreference == 'User'"
                                  :popper-class="widthPro == '50%' ? 'specal_select_width' : 'specal_slect_classs'"
                                  filterable clearable style="flex-grow:1;width:100%;"
                                  @change="changeChildSelect($event, scope, tti, scope.row.paramName, infoReviewForm, infoReviewAttr, scopeds.$index)"
                                  :multiple="scope.row.childAttrArr[scopeds.$index].paramList"
                                  v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][scopeds.$index][ttitem.paramName]" 
                                  remote :remote-method="value => filterMethod(value, { row: scope.row.childAttrArr[scopeds.$index] }, infoReviewForm[scope.row.paramName + 'ChildObj'][scopeds.$index])"
                                  placeholder="请选择"
                                  :class="scope.row.childAttrArr[scopeds.$index][tti].colorRed === '5' ? 'colorRedss' : ''"
                                  disabled>
                                  <el-option v-for="item in scope.row.childAttrArr[scopeds.$index][tti].paramDataList"
                                    :key="item.id" :label="item.name" :value="item.value">
                                    <span style="display:block;font-size:15px;font-weight:blod">{{ item.name }}
                                    </span>
                                    <span style="color: #8492a6; font-size: 13px">{{ item.content }}</span>
                                  </el-option>
                                </el-select>
                                <el-select v-else
                                  :popper-class="widthPro == '50%' ? 'specal_select_width' : 'specal_slect_classs'"
                                  filterable clearable style="flex-grow:1;width:100%;"
                                  @change="changeChildSelect($event, scope, tti, scope.row.paramName, infoReviewForm, infoReviewAttr, scopeds.$index)"
                                  @focus="focusChildSelect(scope, tti, scope.row.paramName, infoReviewForm, scope.row.childAttrArr[scopeds.$index], scopeds.$index, ttitem.paramName)"
                                  :multiple="scope.row.childAttrArr[scopeds.$index].paramList"
                                  v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][scopeds.$index][ttitem.paramName]"
                                  placeholder="请选择"
                                  :class="scope.row.childAttrArr[scopeds.$index][tti].colorRed === '5' ? 'colorRedss' : ''"
                                  disabled>
                                  <el-option v-for="item in scope.row.childAttrArr[scopeds.$index][tti].paramDataList"
                                    :key="item.id" :label="item.name" :value="item.value">
                                    <span style="display:block;font-size:15px;font-weight:blod">{{ item.name }}
                                    </span>
                                    <span style="color: #8492a6; font-size: 13px">{{ item.content }}</span>
                                  </el-option>
                                </el-select>
                              </template>
                              <template v-else>
                                <div v-if="ttitem.paramEnum">
                                  <el-select filterable clearable style="width:100%;"
                                    :multiple="scope.row.childAttrArr[scopeds.$index].paramList"
                                    v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][scopeds.$index][ttitem.paramName]"
                                    placeholder="请选择"
                                    @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName + 'ChildObj'][scopeds.$index][ttitem.paramName], scope.row.childAttrArr[scopeds.$index], scopeds.$index, tti)"
                                    :class="scope.row.childAttrArr[scopeds.$index][tti].colorRed === '5' ? 'colorRedss' : ''"
                                    disabled>
                                    <el-option
                                      v-for="(item, index) in scope.row.childAttrArr[scopeds.$index][tti].paramDataList"
                                      :key="index" :label="item.name" :value="item.value">
                                    </el-option>
                                  </el-select>
                                </div>
                                <div v-else>
                                  <div v-if="ttitem.paramType == 'string'">
                                    <el-input v-if="ttitem.paramStyle == 'password'" style="width:100%;"
                                      placeholder="请输入" clearable
                                      v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][scopeds.$index][ttitem.paramName]"
                                      show-password type="text" :maxlength="ttitem.paramMaxLength"
                                      @change="validatingchild($event, scope, scopeds.$index, tti, ttitem.paramMinLength, ttitem.paramMaxLength, ttitem.paramPatterns, 'string')"
                                      :class="scope.row.childAttrArr[scopeds.$index][tti].colorRed === '2' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '3' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '4' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '5' ? 'colorReds' : ''"
                                      show-word-limit></el-input>
                                    <el-input v-if="ttitem.paramStyle == 'textarea' || ttitem.paramStyle == 'text'"
                                      clearable style="width:100%;margin-top:5px" placeholder="请输入"
                                      v-model="infoReviewForm[scope.row.paramName + 'ChildObj'][scopeds.$index][ttitem.paramName]"
                                      type="textarea" :maxlength="ttitem.paramMaxLength"
                                      @change="validatingchild($event, scope, scopeds.$index, tti, ttitem.paramMinLength, ttitem.paramMaxLength, ttitem.paramPatterns, 'string')"
                                      :class="scope.row.childAttrArr[scopeds.$index][tti].colorRed === '2' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '3' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '4' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '5' ? 'colorReds' : ''"
                                      :autosize="{ minRows: 3, maxRows: 6 }" show-word-limit></el-input>
                                    <el-input style="width:100%;border:none" disabled
                                      v-if="ttitem.paramStyle == 'description'" placeholder="请输入"
                                      v-model="ttitem.paramDefault" class="noneBorder"></el-input>
                                    <el-input v-if="!ttitem.paramStyle || ttitem.paramStyle == ''" style="width:100%;"
                                      placeholder="请输入"
                                      v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][scopeds.$index][ttitem.paramName]"
                                      type="text" :maxlength="ttitem.paramMaxLength"
                                      @change="validatingchild($event, scope, scopeds.$index, tti, ttitem.paramMinLength, ttitem.paramMaxLength, ttitem.paramPatterns, 'string')"
                                      :class="scope.row.childAttrArr[scopeds.$index][tti].colorRed === '2' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '3' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '4' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '5' ? 'colorReds' : ''"
                                      show-word-limit clearable></el-input>
                                  </div>
                                  <div v-else>
                                    <template v-if="ttitem.paramStyle == 'description'">
                                      <el-input style="width:100%;border:none" disabled
                                        v-if="ttitem.paramType == 'number' && !ttitem.paramMutipleof" clearable
                                        placeholder="请输入" v-model="ttitem.paramDefault" class="noneBorder"></el-input>
                                      <el-input-number style="width:100%;" clearable disabled placeholder="请输入" v-else
                                        v-model.number="infoReviewForm[scope.row.paramName + 'ChildObj'][scopeds.$index][ttitem.paramName]"
                                        type="number"
                                        @change="validatingchild($event, scope, scopeds.$index, tti, ttitem.paramMinmum, ttitem.paramMaxmum, ttitem.paramPatterns, 'number')"
                                        :class="scope.row.childAttrArr[scopeds.$index][tti].colorRed === '2' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '3' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '4' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '5' ? 'colorReds noneBorder' : 'noneBorder'"
                                        :step="ttitem.paramMutipleof" step-strictly></el-input-number>
                                    </template>
                                    <template v-else>
                                      <el-input style="width:100%;border:none" clearable
                                        v-if="ttitem.paramType == 'number' && !ttitem.paramMutipleof" placeholder="请输入"
                                        v-model.number="infoReviewForm[scope.row.paramName + 'ChildObj'][scopeds.$index][ttitem.paramName]"></el-input>
                                      <el-input-number style="width:100%;" clearable placeholder="请输入" v-else
                                        v-model.number="infoReviewForm[scope.row.paramName + 'ChildObj'][scopeds.$index][ttitem.paramName]"
                                        type="number"
                                        @change="validatingchild($event, scope, scopeds.$index, tti, ttitem.paramMinmum, ttitem.paramMaxmum, ttitem.paramPatterns, 'number')"
                                        :class="scope.row.childAttrArr[scopeds.$index][tti].colorRed === '2' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '3' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '4' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '5' ? 'colorReds' : ''"
                                        :step="ttitem.paramMutipleof" step-strictly></el-input-number>
                                    </template>
                                  </div>
                                </div>
                              </template>
                              <div
                                v-if="ttitem.paramType == 'string' && (ttitem.paramStyle == 'textarea' || ttitem.paramStyle == 'text')">
                                <div
                                  v-if="ttitem.paramDescription && scope.row.childAttrArr[scopeds.$index][tti].colorRed === '1'"
                                  style="position: relative;font-size: 12px;line-height: 10px;bottom: -6px;color: #c1c1c1;"
                                  v-html="ttitem.paramDescription"></div>
                                <div
                                  v-if="scope.row.childAttrArr[scopeds.$index][tti].colorRed === '2' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '3' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '4' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '5'"
                                  style="position: relative;font-size: 12px;line-height: 10px;bottom: -6px;color: red;">
                                  {{ ttitem.paramConstraintDescription ? ttitem.paramConstraintDescription : ttitem.customDescription }}
                                </div>
                              </div>
                              <div v-else>
                                <div
                                  v-if="ttitem.paramDescription && scope.row.childAttrArr[scopeds.$index][tti].colorRed === '1'"
                                  style="position: relative;font-size: 12px;line-height: 12px;bottom:-2px;color: #c1c1c1;"
                                  v-html="ttitem.paramDescription"></div>
                                <div
                                  v-if="scope.row.childAttrArr[scopeds.$index][tti].colorRed === '2' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '3' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '4' || scope.row.childAttrArr[scopeds.$index][tti].colorRed === '5'"
                                  style="position: relative;font-size: 12px;line-height: 12px;bottom:-2px;color: red;">
                                  {{ ttitem.paramConstraintDescription ? ttitem.paramConstraintDescription  : ttitem.customDescription }}
                                </div>
                              </div>
                            </template>
                          </el-table-column>
                          <el-table-column fixed="right" label="操作" width="50" align="center">
                            <template slot-scope="scopeds">
                              <i style="cursor: pointer;" class="el-icon-delete"
                                @click="deletCardAttr(scope, scope.row.paramName, scopeds.$index)"
                                v-if="scope.row.childAttrArr.length != 1"></i>
                            </template>
                          </el-table-column>
                        </el-table>
                      </template>
                      <template v-else>
                        <div v-if="scope.row.childAttrArr || scope.row.paramreference == 'FormTemplate'">
                          <el-button type='primary' v-if="nodeCheckName == 'FormInitBatchCloudServer'"
                            @click="batchIn = true">批量导入</el-button>
                          <el-card style="margin-top:13px;padding-right: 60px;position:relative"
                            v-for="(childItemAtt, iii) in scope.row.childAttrArr" :key="iii">
                            <span
                              style="position:absolute;right:15px;top:0px;font-size:15px;cursor: pointer;z-index:99;">
                              <i class="el-icon-document-copy" title="复制" style="margin-right: 10px;"
                                @click="copyCardAttr(scope, scope.row.paramName, iii)"></i>
                              <i class="el-icon-delete" title="删除"
                                @click="deletCardAttr(scope, scope.row.paramName, iii)"></i>
                            </span>
                            <template v-if="infoReviewForm[scope.row.paramName + 'ChildObj'][iii]">
                              <div style="font-size: 13px;font-weight: 900;cursor: pointer;">
                                <span @click="checkShow(scope.row.paramName, iii)"
                                  style="position: absolute;top: 0;left: 0px;display: block;width: 100%;padding-left: 15px;">
                                  <span
                                    :class="infoReviewForm[scope.row.paramName + 'ChildObj'][iii]['_isShowStatus'] ? 'el-icon-arrow-right' : 'el-icon-arrow-down'">
                                    <span style="font-size:13px;margin-left:10px;position: relative;"> 批次{{ iii + 1 }} -
                                      {{ childItemAtt[0].paramFromTitle }}
                                    </span>
                                    <span v-for="(ttitem, tti) in childItemAtt[0].showTitleArr" :key="tti">
                                      <span
                                        v-if="idChangeNameObj[infoReviewForm[scope.row.paramName + 'ChildObj'][iii][ttitem]]">
                                        {{ idChangeNameObj[infoReviewForm[scope.row.paramName +'ChildObj'][iii][ttitem]] }}
                                      </span>
                                      <span v-else>
                                        {{ infoReviewForm[scope.row.paramName + 'ChildObj'][iii][ttitem] }}
                                      </span>
                                    </span>
                                  </span>
                                </span>
                              </div>
                              <div style="margin-top:20px;"
                                v-show="lookType ? infoReviewForm[scope.row.paramName + 'ChildObj'][iii]['_isShowStatus'] : !infoReviewForm[scope.row.paramName + 'ChildObj'][iii]['_isShowStatus']">
                                <template v-for="(childItem, iiii) in childItemAtt">
                                  <div v-if="childItem.paramGroup" :key="'g' + iiii">
                                    <div v-if="iiii != 0 && childItem.thisShowIf == '1'">
                                      <h4
                                        style="margin: 10px 0px 15px; padding: 0px; color: rgb(51, 51, 51); font-size: 15px; font-weight: 600; text-indent: 20px; border-left: 4px solid rgb(35, 173, 229);line-height: 20px;">
                                        {{ childItem.paramGroup }} </h4>
                                    </div>
                                  </div>
                                  <el-form-item v-show="childItem.thisShowIf == '1'" :label="childItem.paramTitle"
                                    :rules="childItem.paramRequired ? [{ required: true, message: '请选择', trigger: 'change' }, { required: true, pattern: /\S/, message: '必填不能为空', trigger: 'change' }] : []"
                                    :key="iiii">
                                    <div v-if="childItem.paramType == 'object'">
                                      <template v-if="childItem.paramreference == 'FormTemplate'">
                                        <form-tem-child :key="'g' + iiii" v-if="childItem.childAttrArrThird"
                                          :childData='childItem.childAttrArrThird' :ifStyleTable="childItem.paramStyle"
                                          :childObject="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName + 'ChildObj3']"
                                          :infoReviewForm="infoReviewForm" :infoReviewAttr="infoReviewAttr"
                                          :waiIndex='waiIndex' :batchIndex="batchIndex"
                                          :aiii="{ erindex: iii, sani: iiii }" @emitChild='emitChild'
                                          @emitChildvalidating='emitChildvalidating' @emitChildcopy="emitChildcopy"
                                          @emitChilddelOne="emitChilddelOne" @emitChildSureSer="emitChildSureSer"
                                          @emitChildfocusChildGLArr='emitChildfocusChildGLArr'
                                          :smallName="childItem.paramName + 'ChildObj3'"
                                          :parentName="scope.row.paramName + 'ChildObj'"
                                          :erinfoChildattr='scope.row.childAttrArr'
                                          :erinfoClidObj="infoReviewForm[scope.row.paramName + 'ChildObj'][iii]"
                                          :oneinfoClidObj="infoReviewForm"></form-tem-child>
                                        <template
                                          v-if="childItem.paramDataList && !detailOrLook && childItem.paramStyle != 'table'">
                                          <el-dropdown v-if="childItem.paramDataList.length != 1" size="small"
                                            type="primary"
                                            @command="showChildAttrThird(scope, childItem.paramName, $event, null, iiii, null, scope.$index, iii)">
                                            <span class="el-dropdown-link" style="font-size:12px">
                                              添加批次<i class="el-icon-caret-bottom  el-icon--right"></i>
                                            </span>
                                            <el-dropdown-menu slot="dropdown">
                                              <el-dropdown-item v-for="(item, index) in childItem.paramDataList"
                                                :key="index" :command="item.value">{{ item.name }}</el-dropdown-item>
                                            </el-dropdown-menu>
                                          </el-dropdown>
                                          <template v-else>
                                            <span style="color:#409eff;cursor: pointer;margin-top:10px;font-size:12px"
                                              v-for="(item, index) in childItem.paramDataList" :key="index"
                                              @click="showChildAttrThird(scope, childItem.paramName, childItem.paramDataList[0].value, null, iiii, null, scope.$index, iii)">
                                              <i class="el-icon-plus" style="color:#409eff;"></i>添加批次
                                            </span>
                                          </template>
                                        </template>
                                        <template
                                          v-if="childItem.paramDataList && detailOrLook && users.indexOf(usernameAll) != -1 && childItem.paramStyle != 'table'">
                                          <el-dropdown v-if="childItem.paramDataList.length != 1" size="small"
                                            type="primary"
                                            @command="showChildAttrThird(scope, childItem.paramName, $event, null, iiii, null, scope.$index, iii)">
                                            <span class="el-dropdown-link" style="font-size:12px">
                                              添加批次<i class="el-icon-caret-bottom  el-icon--right"></i>
                                            </span>
                                            <el-dropdown-menu slot="dropdown">
                                              <el-dropdown-item v-for="(item, index) in childItem.paramDataList"
                                                :key="index" :command="item.value">{{ item.name }}</el-dropdown-item>
                                            </el-dropdown-menu>
                                          </el-dropdown>
                                          <template v-else>
                                            <span style="color:#409eff;cursor: pointer;margin-top:10px;font-size:12px"
                                              v-for="(item, index) in childItem.paramDataList" :key="index"
                                              @click="showChildAttrThird(scope, childItem.paramName, childItem.paramDataList[0].value, null, iiii, null, scope.$index, iii)">
                                              <i class="el-icon-plus" style="color:#409eff;"></i>添加批次
                                            </span>
                                          </template>
                                        </template>
                                      </template>
                                      <template v-else>
                                        <el-button
                                          v-if="(childItem.paramList && (childItem.paramreferenceShow.length > 1 || childItem.paramreferenceDisplay) && childItem.paramreference != 'FormTemplate' && childItem.paramStyle != 'search') || (childItem.paramList && childItem.paramreference != 'FormTemplate' && childItem.paramStyle == 'tableStyle')"
                                          @click="choseServerList(scope, iiii, scope.row.paramName, infoReviewForm, infoReviewAttr, childItem, childItem.paramName, iii, childItem.paramTitle, false)">选择{{childItem.paramTitle }}</el-button>
                                        <template v-else>
                                          <!-- 子表单radio -->
                                          <div style="display:flex;width:100%;">
                                            <div style="flex-grow:1;" v-if="childItem.paramStyle == 'radio'">
                                              <div class="getQLBox">
                                                <template v-if="!childItem.paramList">
                                                  <el-radio-group
                                                    v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                                    size="mini"
                                                    @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName], childItem, iii, iiii)"
                                                    :class="childItem.colorRed === '5' ? 'colorRedss' : ''">
                                                    <el-radio-button
                                                      v-for="(itemOne, eeeindex) in childItem.paramDataList"
                                                      :key="eeeindex" :label="itemOne.value">
                                                      <template>
                                                        <el-tooltip class="item" effect="dark" placement="top">
                                                          <div slot="content">
                                                            <div style="text-align:center;">{{ itemOne.name }}</div>
                                                            {{ itemOne.content }}
                                                          </div>
                                                          <div>
                                                            {{ itemOne.name }}
                                                          </div>
                                                        </el-tooltip>
                                                      </template>
                                                    </el-radio-button>
                                                  </el-radio-group>
                                                </template>
                                                <el-checkbox-group v-else
                                                  v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                                  size="mini"
                                                  @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName], childItem, iii, iiii)"
                                                  :class="childItem.colorRed === '5' ? 'colorRedss' : ''">
                                                  <el-checkbox-button
                                                    v-for="(itemOne, eeeindex) in childItem.paramDataList"
                                                    :key="eeeindex" :label="itemOne.value">
                                                    <template>
                                                      <el-tooltip class="item" effect="dark" placement="top">
                                                        <div slot="content">
                                                          <div style="text-align:center;">{{ itemOne.name }}</div>
                                                          {{ itemOne.content }}
                                                        </div>
                                                        <div>
                                                          {{ itemOne.name }}
                                                        </div>
                                                      </el-tooltip>
                                                    </template>
                                                  </el-checkbox-button>
                                                </el-checkbox-group>
                                              </div>
                                            </div>
                                            <template v-else>
                                              <el-select
                                                v-if="childItem.paramStyle == 'search' || childItem.paramreference == 'User'"
                                                :popper-class="widthPro == '50%' ? 'specal_select_width' : 'specal_slect_classs'"
                                                filterable clearable style="flex-grow:1;"
                                                @focus="getMetaSchemasListvalFocus('', { row: childItem }, infoReviewForm[scope.row.paramName + 'ChildObj'][iii])"
                                                @change="changeChildSelect($event, scope, iiii, scope.row.paramName, infoReviewForm, infoReviewAttr, iii)"
                                                :multiple="childItem.paramList"
                                                v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                                remote :remote-method="value => filterMethod(value, { row: childItem }, infoReviewForm[scope.row.paramName + 'ChildObj'][iii])"
                                                placeholder="请选择"
                                                :class="childItem.colorRed === '5' ? 'colorRedss' : ''">
                                                <el-option v-for="item in childItem.paramDataList" :key="item.id"
                                                  :label="item.name" :value="item.value">
                                                  <span
                                                    style="display:block;font-size:15px;font-weight:blod">{{item.name }}
                                                  </span>
                                                  <span style="color: #8492a6; font-size: 13px">{{ item.content
                                                    }}</span>
                                                </el-option>
                                              </el-select>
                                              <el-select v-else
                                                :popper-class="widthPro == '50%' ? 'specal_select_width' : 'specal_slect_classs'"
                                                filterable clearable style="flex-grow:1;"
                                                @change="changeChildSelect($event, scope, iiii, scope.row.paramName, infoReviewForm, infoReviewAttr, iii)"
                                                @focus="focusChildSelect(scope, iiii, scope.row.paramName, infoReviewForm, childItem, iii, childItem.paramName)"
                                                :multiple="childItem.paramList"
                                                v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                                placeholder="请选择"
                                                :class="childItem.colorRed === '5' ? 'colorRedss' : ''">
                                                <el-option v-for="item in childItem.paramDataList" :key="item.id"
                                                  :label="item.name" :value="item.value">
                                                  <span
                                                    style="display:block;font-size:15px;font-weight:blod">{{item.name }}
                                                  </span>
                                                  <span style="color: #8492a6; font-size: 13px">{{ item.content
                                                    }}</span>
                                                </el-option>
                                              </el-select>
                                            </template>
                                            <span v-if="childItem.paramStyle == 'radio'"
                                              style="width: 30px;height: 30px;line-height: 30px;font-size: 16px;cursor: pointer;position: absolute;right: -50px;top: 7px;"
                                              class="el-icon-refresh"
                                              @click="focusChildSelect(scope, iiii, scope.row.paramName, infoReviewForm, childItem, iii, childItem.paramName, true)"></span>
                                          </div>
                                        </template>
                                        <el-table border class="specialTable"
                                          v-if="childItem.paramList && (childItem.paramreferenceShow.length > 1 || childItem.paramreferenceDisplay) && childItem.paramreference != 'FormTemplate' && (theadArrShow[iii] && theadArrShow[iii].length)"
                                          :data="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                          style="width: 100%;margin-top:10px;max-width:700px;" size="mini"
                                          :show-header="true" max-height="300" min-height="100">
                                          <template v-for="(taItems, i) in theadArrShow[iii]">
                                            <el-table-column :key="i" :show-overflow-tooltip="true"
                                              :label="taItems.title" :prop="taItems.value"
                                              v-if="childItem.paramreferenceLink && (taItems.value == 'name' || taItems.title == '名称')">
                                              <template slot-scope="scoped">
                                                <el-link type="primary" style="font-size:12px;"
                                                  @click="linkJumpTo(childItem, scoped.row)">{{scoped.row[taItems.value] }}</el-link>
                                              </template>
                                            </el-table-column>
                                            <el-table-column :key="i" :show-overflow-tooltip="true"
                                              :label="taItems.title" :prop="taItems.value" v-else>
                                            </el-table-column>
                                          </template>
                                          <el-table-column fixed="right" v-if="theadArrShow[iii]" label="操作" width="60"
                                            align="center">
                                            <template slot-scope="scoped">
                                              <i style="cursor: pointer;" class="el-icon-delete"
                                                @click="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName].splice(scoped.$index, 1)"></i>
                                            </template>
                                          </el-table-column>
                                        </el-table>
                                        <span style="color: #666;font-size:13px;"
                                          v-if="childItem.paramList && (childItem.paramreferenceShow.length > 1 || childItem.paramreferenceDisplay) && infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName] && childItem.paramreference != 'FormTemplate' && (theadArrShow[iii] && theadArrShow[iii].length)">已选数量：{{infoReviewForm[scope.row.paramName +  'ChildObj'][iii][childItem.paramName].length }}</span>
                                      </template>
                                    </div>
                                    <div v-else>
                                      <div v-if="childItem.paramEnum">
                                        <div v-if="childItem.paramStyle == 'radio'">
                                          <template v-if="!childItem.paramList">
                                            <el-radio-group
                                              v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                              @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName], childItem, iii, iiii)"
                                              :class="childItem.colorRed === '5' ? 'colorRedss' : ''" size="mini">
                                              <el-radio-button v-for="(item, index) in childItem.paramDataList"
                                                :key="index" :label="item.value">{{ item.name }}</el-radio-button>
                                            </el-radio-group>
                                          </template>
                                          <el-checkbox-group v-else
                                            v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                            @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName], childItem, iii, iiii)"
                                            :class="childItem.colorRed === '5' ? 'colorRedss' : ''" size="mini">
                                            <el-checkbox-button v-for="(item, index) in childItem.paramDataList"
                                              :key="index" :label="item.value">{{item.name }}</el-checkbox-button>
                                          </el-checkbox-group>
                                        </div>
                                        <el-select v-else filterable clearable style="width:100%;"
                                          :multiple="childItem.paramList"
                                          v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                          placeholder="请选择"
                                          @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName], childItem, iii, iiii)"
                                          :class="childItem.colorRed === '5' ? 'colorRedss' : ''">
                                          <el-option v-for="(item, index) in childItem.paramDataList" :key="index"
                                            :label="item.name" :value="item.value">
                                          </el-option>
                                        </el-select>
                                      </div>
                                      <div v-else>
                                        <template v-if="childItem.paramprerequisite.length">
                                          <el-select style="width:100%;" filterable clearable
                                            @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName], childItem, iii, iiii)"
                                            @focus="focusChildGLArr(childItem, iii, scope.row.paramName, iiii)"
                                            :multiple="childItem.paramList"
                                            v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                            placeholder="请选择" :class="childItem.colorRed === '5' ? 'colorRedss' : ''">
                                            <el-option v-for="item in childItem.paramDataList" :key="item.name"
                                              :label="item.name" :value="item.value">
                                            </el-option>
                                          </el-select>
                                        </template>
                                        <template v-else>
                                          <template>
                                            <div v-if="childItem.paramList">
                                              <template
                                                v-if="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]">
                                                <div style="display: flex;align-items: center;margin-bottom:7px;"
                                                  v-for="(item, index) in infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                                  :key="index">
                                                  <el-input style="width:100%;" placeholder="请输入"
                                                    v-if="childItem.paramType == 'string'" clearable
                                                    v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName][index]"
                                                    type="input"
                                                    @change="validatingchild($event, scope, iii, iiii, childItem.paramMinLength, childItem.paramMaxLength, childItem.paramPatterns, 'string', infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName])"></el-input>
                                                  <el-input-number style="width:100%;" clearable placeholder="请输入"
                                                    v-else
                                                    v-model.number="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName][index]"
                                                    type="number" :maxlength="childItem.paramMaxmum"
                                                    :minlength="childItem.paramMinmum" :step="childItem.paramMutipleof"
                                                    step-strictly
                                                    @change="validatingchild($event, scope, iii, iiii, childItem.paramMinLength, childItem.paramMaxLength, childItem.paramPatterns, 'number', infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName])"></el-input-number>
                                                </div>
                                              </template>
                                              <span class="el-icon-plus"
                                                style="font-size:18px;margin:0 5px;cursor: pointer;"
                                                @click="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName].push(null)"></span>
                                              <span class="el-icon-minus"
                                                style="font-size:18px;margin:0 5px;cursor: pointer;color:#f60;"
                                                @click="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName].splice(index, 1)"></span>
                                            </div>
                                            <div v-else>
                                              <template
                                                v-if="childItem.paramType == 'string' && (childItem.paramStyle == 'datetime' || childItem.paramStyle == 'date')">
                                                <el-date-picker v-model="value1" type="datetime" placeholder="选择日期时间"
                                                  v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                                  style="width:100%"
                                                  @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName], childItem, iii, iiii)"
                                                  v-if="childItem.paramStyle == 'datetime'">
                                                </el-date-picker>
                                                <el-date-picker v-model="value1" type="date" placeholder="选择日期"
                                                  v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                                  style="width:100%"
                                                  @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName], childItem, iii, iiii)"
                                                  v-if="childItem.paramStyle == 'date'" value-format="yyyy-MM-dd">
                                                </el-date-picker>
                                              </template>
                                              <div v-else>
                                                <div v-if="childItem.paramType == 'string'">
                                                  <el-input v-if="childItem.paramStyle == 'password'"
                                                    style="width:100%;" placeholder="请输入" clearable
                                                    v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                                    show-password type="text" :maxlength="childItem.paramMaxLength"
                                                    @change="validatingchild($event, scope, iii, iiii, childItem.paramMinLength, childItem.paramMaxLength, childItem.paramPatterns, 'string')"
                                                    :class="childItem.colorRed === '2' || childItem.colorRed === '3' || childItem.colorRed === '4' || childItem.colorRed === '5' ? 'colorReds' : ''"
                                                    show-word-limit></el-input>
                                                  <el-input
                                                    v-if="childItem.paramStyle == 'textarea' || childItem.paramStyle == 'text'"
                                                    clearable style="width:100%;margin-top:5px" placeholder="请输入"
                                                    v-model="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                                    type="textarea" :maxlength="childItem.paramMaxLength"
                                                    @change="validatingchild($event, scope, iii, iiii, childItem.paramMinLength, childItem.paramMaxLength, childItem.paramPatterns, 'string')"
                                                    :class="childItem.colorRed === '2' || childItem.colorRed === '3' || childItem.colorRed === '4' || childItem.colorRed === '5' ? 'colorReds' : ''"
                                                    :autosize="{ minRows: 3, maxRows: 6 }" show-word-limit></el-input>
                                                  <el-input style="width:100%;border:none" disabled
                                                    v-if="childItem.paramStyle == 'description'" placeholder="请输入"
                                                    v-model="childItem.paramDefault" class="noneBorder"></el-input>
                                                  <el-input clearable
                                                    v-if="!childItem.paramStyle || childItem.paramStyle == ''"
                                                    style="width:100%;" placeholder="请输入"
                                                    v-model.trim="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                                    type="text" :maxlength="childItem.paramMaxLength"
                                                    @change="validatingchild($event, scope, iii, iiii, childItem.paramMinLength, childItem.paramMaxLength, childItem.paramPatterns, 'string')"
                                                    :class="childItem.colorRed === '2' || childItem.colorRed === '3' || childItem.colorRed === '4' || childItem.colorRed === '5' ? 'colorReds' : ''"
                                                    show-word-limit></el-input>
                                                </div>
                                                <div v-else>
                                                  <template v-if="childItem.paramStyle == 'description'">
                                                    <el-input style="width:100%;border:none" clearable disabled
                                                      v-if="childItem.paramType == 'number' && !childItem.paramMutipleof"
                                                      placeholder="请输入" v-model="childItem.paramDefault"
                                                      class="noneBorder"></el-input>
                                                    <el-input-number style="width:100%;" clearable disabled
                                                      placeholder="请输入" v-else
                                                      v-model.number="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                                      type="number"
                                                      @change="validatingchild($event, scope, iii, iiii, childItem.paramMinmum, childItem.paramMaxmum, childItem.paramPatterns, 'number')"
                                                      :class="childItem.colorRed === '2' || childItem.colorRed === '3' || childItem.colorRed === '4' || childItem.colorRed === '5' ? 'colorReds noneBorder' : 'noneBorder'"
                                                      :step="childItem.paramMutipleof" step-strictly></el-input-number>
                                                  </template>
                                                  <template v-else>
                                                    <el-input style="width:100%;border:none"
                                                      v-if="childItem.paramType == 'number' && !childItem.paramMutipleof"
                                                      placeholder="请输入" clearable
                                                      v-model.number="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                                      @change="validatingchild($event, scope, iii, iiii, childItem.paramMinmum, childItem.paramMaxmum, childItem.paramPatterns, 'number')"></el-input>
                                                    <el-input-number style="width:100%;" clearable placeholder="请输入"
                                                      v-else
                                                      v-model.number="infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName]"
                                                      type="number"
                                                      @change="validatingchild($event, scope, iii, iiii, childItem.paramMinmum, childItem.paramMaxmum, childItem.paramPatterns, 'number')"
                                                      :class="childItem.colorRed === '2' || childItem.colorRed === '3' || childItem.colorRed === '4' || childItem.colorRed === '5' ? 'colorReds' : ''"
                                                      :step="childItem.paramMutipleof" step-strictly></el-input-number>
                                                  </template>
                                                </div>
                                              </div>
                                            </div>
                                          </template>
                                        </template>
                                      </div>
                                    </div>
                                    <i v-if="childItem.paramStyle == 'radio'" class="el-icon-circle-close"
                                      style="width: 30px;height: 30px;line-height: 30px;font-size: 16px;cursor: pointer;position: absolute;right: -25px;top: 7px;"
                                      @click="clearableRadioBtn(childItem, scope, infoReviewForm[scope.row.paramName + 'ChildObj'][iii][childItem.paramName], childItem, iii, iiii, infoReviewForm[scope.row.paramName + 'ChildObj'][iii], childItem.paramName)"></i>
                                    <div
                                      v-if="childItem.paramType == 'string' && (childItem.paramStyle == 'textarea' || childItem.paramStyle == 'text')">
                                      <div v-if="childItem.paramDescription && childItem.colorRed === '1'"
                                        style="position: relative;font-size: 12px;line-height: 10px;bottom: -6px;color: #c1c1c1;"
                                        v-html="childItem.paramDescription"></div>
                                      <div
                                        v-if="childItem.colorRed === '2' || childItem.colorRed === '3' || childItem.colorRed === '4' || childItem.colorRed === '5'"
                                        style="position: relative;font-size: 12px;line-height: 10px;bottom: -6px;color: red;">
                                        {{ childItem.paramConstraintDescription ? childItem.paramConstraintDescription : childItem.customDescription }}
                                      </div>
                                    </div>
                                    <div v-else>
                                      <div v-if="childItem.paramDescription && childItem.colorRed === '1'"
                                        style="position: relative;font-size: 12px;line-height: 12px;bottom:-2px;color: #c1c1c1;"
                                        v-html="childItem.paramDescription"></div>
                                      <div
                                        v-if="childItem.colorRed === '2' || childItem.colorRed === '3' || childItem.colorRed === '4' || childItem.colorRed === '5'"
                                        style="position: relative;font-size: 12px;line-height: 12px;bottom:-2px;color: red;">
                                        {{ childItem.paramConstraintDescription ? childItem.paramConstraintDescription : childItem.customDescription }}
                                      </div>
                                    </div>
                                  </el-form-item>
                                </template>
                              </div>
                            </template>
                          </el-card>
                        </div>
                        <div v-else>
                          <div
                            v-if="scope.row.paramList && (scope.row.paramreferenceShow.length > 1 || scope.row.paramreferenceDisplay) && scope.row.paramStyle != 'search' && scope.row.paramreference != 'FormTemplate'">
                            <el-button
                              v-if="(scope.row.paramList && (scope.row.paramreferenceShow.length > 1 || scope.row.paramreferenceDisplay)) || (scope.row.paramList && scope.row.paramStyle == 'tableStyle' && scope.row.paramreference != 'FormTemplate')"
                              @click="choseList(scope, scope.row.paramName, infoReviewForm, infoReviewAttr, scope.row)">选择{{scope.row.paramTitle }}</el-button>
                            <div v-if="listObj[scope.row.paramName]">
                              <el-table border class="specialTable"
                                v-if="scope.row.paramList && (scope.row.paramreferenceShow.length > 1 || scope.row.paramreferenceDisplay) && infoReviewForm[scope.row.paramName]"
                                :data="infoReviewForm[scope.row.paramName]" style="width: 100%;margin-top:10px"
                                size="mini" :show-header="true" max-height="300" min-height="100">
                                <template v-for="(taItems, i) in listObj[scope.row.paramName]">
                                  <el-table-column :key="i" :show-overflow-tooltip="true" :label="taItems.title"
                                    :prop="taItems.value"
                                    v-if="scope.row.paramreferenceLink && (taItems.title == '名称' || taItems.value == 'name')">
                                    <template slot-scope="scoped">
                                      <el-link type="primary" style="font-size:12px;"
                                        @click="linkJumpTo(scope.row, scoped.row)">{{ scoped.row[taItems.value]
                                        }}</el-link>
                                    </template>
                                  </el-table-column>
                                  <el-table-column :key="i" :show-overflow-tooltip="true" :label="taItems.title"
                                    :prop="taItems.value" v-else>
                                  </el-table-column>
                                </template>
                                <el-table-column fixed="right" v-if="infoReviewForm[scope.row.paramName]" label="操作"
                                  width="100" align="center">
                                  <template slot-scope="scoped">
                                    <i style="cursor: pointer;" class="el-icon-delete"
                                      @click="infoReviewForm[scope.row.paramName].splice(scoped.$index, 1)"></i>
                                  </template>
                                </el-table-column>
                              </el-table>
                              <span style="color: #666;font-size:13px;"
                                v-if="scope.row.paramList && (scope.row.paramreferenceShow.length > 1 || scope.row.paramreferenceDisplay) && infoReviewForm[scope.row.paramName]">已选数量：{{  infoReviewForm[scope.row.paramName] ? infoReviewForm[scope.row.paramName].length : 0}}</span>
                            </div>
                          </div>
                          <div v-else>
                            <div>
                              <el-select style="width:100%;" filterable clearable
                                @change="changeSelectValue($event, scope, scope.row.paramName, infoReviewForm, infoReviewAttr)"
                                @focus="focusSelectValue(scope, scope.row.paramName, infoReviewForm)"
                                :disabled="scope.row.paramName == 'loadBalance' && nodeCheckName == 'FormUpdateCloudServerGroup' || (scope.row.paramName == 'formObjects' && nodeCheckName == 'FormCloudVpcBindZone') || (scope.row.paramName == 'formObject' && (newTaskName.indexOf('编辑') != -1 || newTaskName.indexOf('删除') != -1 ||newTaskName.indexOf('资源巡检添优化') != -1 || newTaskName.indexOf('资源巡检添加白') != -1 || newTaskName.indexOf('恢复') != -1))"
                                v-if="scope.row.paramreference == 'User' || scope.row.paramStyle == 'search'"
                                :multiple="scope.row.paramList" v-model.trim="infoReviewForm[scope.row.paramName]"
                                remote :remote-method="value => filterMethod(value, scope, infoReviewForm)" placeholder="请选择"
                                :popper-class="widthPro == '50%' ? 'specal_select_width' : 'specal_slect_classs'"
                                :class="scope.row.colorRed === '5' ? 'colorRedss' : ''">
                                <el-option v-for="(item, index) in scope.row.paramDataList" :key="index"
                                  :label="item.name" :value="item.value">
                                  <span style="display:block;font-size:15px;font-weight:blod">{{ item.name }}
                                  </span>
                                  <span style="color: #8492a6; font-size: 13px">{{ item.content }}</span>
                                </el-option>
                              </el-select>
                              <div v-else>
                                <template v-if="scope.row.paramStyle == 'radio'">
                                  <div class="getQLBox">
                                    <template v-if="!scope.row.paramList">
                                      <el-radio-group
                                        :disabled="(scope.row.paramName == 'formObject' && newTaskName.indexOf('编辑') != -1)"
                                        v-model.trim="infoReviewForm[scope.row.paramName]" size="mini"
                                        @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName])"
                                        :class="scope.row.colorRed === '5' ? 'colorRedss' : ''">
                                        <el-radio-button v-for="(item, index) in scope.row.paramDataList" :key="index"
                                          :label="item.value">
                                          <template>
                                            <el-tooltip class="item" effect="dark" placement="top">
                                              <div slot="content">
                                                <div style="text-align:center;">{{ item.name }}</div>
                                                {{ item.content }}
                                              </div>
                                              <div>
                                                {{ item.name }}
                                              </div>
                                            </el-tooltip>
                                          </template>
                                        </el-radio-button>
                                      </el-radio-group>
                                    </template>
                                    <el-checkbox-group v-else v-model.trim="infoReviewForm[scope.row.paramName]"
                                      size="mini"
                                      :disabled="(scope.row.paramName == 'formObject' && newTaskName.indexOf('编辑') != -1)"
                                      @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName])"
                                      :class="scope.row.colorRed === '5' ? 'colorRedss' : ''">
                                      <el-checkbox-button v-for="(item, index) in scope.row.paramDataList" :key="index"
                                        :label="item.value">
                                        <template>
                                          <el-tooltip class="item" effect="dark" placement="top">
                                            <div slot="content">
                                              <div style="text-align:center;">{{ item.name }}</div>
                                              {{ item.content }}
                                            </div>
                                            <div>
                                              {{ item.name }}
                                            </div>
                                          </el-tooltip>
                                        </template>
                                      </el-checkbox-button>
                                    </el-checkbox-group>
                                  </div>
                                </template>
                                <el-select v-else
                                  :disabled="(scope.row.paramName == 'formObject' && (newTaskName.indexOf('编辑') != -1|| newTaskName.indexOf('删除') != -1))"
                                  style="width:100%;" filterable clearable
                                  @change="changeSelectValue($event, scope, scope.row.paramName, infoReviewForm, infoReviewAttr)"
                                  @focus="focusSelectValue(scope, scope.row.paramName, infoReviewForm)"
                                  :multiple="scope.row.paramList" v-model.trim="infoReviewForm[scope.row.paramName]"
                                  placeholder="请选择"
                                  :popper-class="widthPro == '50%' ? 'specal_select_width' : 'specal_slect_classs'"
                                  :class="scope.row.colorRed === '5' ? 'colorRedss' : ''">
                                  <el-option v-for="(item, index) in scope.row.paramDataList" :key="index"
                                    :label="item.name" :value="item.value">
                                    <span style="display:block;font-size:15px;font-weight:blod">{{ item.name }}
                                    </span>
                                    <span style="color: #8492a6; font-size: 13px">{{ item.content }}</span>
                                  </el-option>
                                </el-select>
                              </div>
                            </div>
                          </div>
                        </div>
                        <div v-if="scope.row.paramreferenceZ.length > 0">
                          <template v-if="scope.row.paramDataList && !detailOrLook && scope.row.paramStyle != 'table'">
                            <el-dropdown v-if="scope.row.paramDataList.length != 1" size="small" type="primary"
                              @command="showChildAttr(scope, scope.row.paramName, $event)">
                              <span class="el-dropdown-link" style="font-size:12px">
                                添加批次<i class="el-icon-caret-bottom  el-icon--right"></i>
                              </span>
                              <el-dropdown-menu slot="dropdown">
                                <el-dropdown-item v-for="(item, index) in scope.row.paramDataList" :key="index"
                                  :command="item.value">{{ item.name }}</el-dropdown-item>
                              </el-dropdown-menu>
                            </el-dropdown>
                            <template v-else>
                              <span style="color:#409eff;cursor: pointer;margin-top:10px;font-size:12px"
                                v-for="(item, index) in scope.row.paramDataList" :key="index"
                                @click="showChildAttr(scope, scope.row.paramName, scope.row.paramDataList[0].value)">
                                <i class="el-icon-plus" style="color:#409eff;"></i>添加批次
                              </span>
                            </template>
                          </template>
                          <template
                            v-if="scope.row.paramDataList && detailOrLook && users.indexOf(usernameAll) != -1 && scope.row.paramStyle != 'table'">
                            <el-dropdown v-if="scope.row.paramDataList.length != 1" size="small" type="primary"
                              @command="showChildAttr(scope, scope.row.paramName, $event)">
                              <span class="el-dropdown-link" style="font-size:12px">
                                添加批次<i class="el-icon-caret-bottom  el-icon--right"></i>
                              </span>
                              <el-dropdown-menu slot="dropdown">
                                <el-dropdown-item v-for="(item, index) in scope.row.paramDataList" :key="index"
                                  :command="item.value">{{ item.name }}</el-dropdown-item>
                              </el-dropdown-menu>
                            </el-dropdown>
                            <template v-else>
                              <span style="color:#409eff;cursor: pointer;margin-top:10px;font-size:12px"
                                v-for="(item, index) in scope.row.paramDataList" :key="index"
                                @click="showChildAttr(scope, scope.row.paramName, scope.row.paramDataList[0].value)">
                                <i class="el-icon-plus" style="color:#409eff;"></i>添加批次
                              </span>
                            </template>
                          </template>
                        </div>
                      </template>
                    </div>
                    <div v-else>
                      <div v-if="scope.row.paramEnum">
                        <div v-if="scope.row.paramStyle == 'radio'">
                          <template v-if="!scope.row.paramList">
                            <el-radio-group v-model.trim="infoReviewForm[scope.row.paramName]" size="mini"
                              @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName])"
                              :class="scope.row.colorRed === '5' ? 'colorRedss' : ''">
                              <el-radio-button v-for="(item, index) in scope.row.paramDataList" :key="index"
                                :label="item.value">{{ item.name }}</el-radio-button>
                            </el-radio-group>
                          </template>
                          <el-checkbox-group v-else v-model.trim="infoReviewForm[scope.row.paramName]" size="mini"
                            @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName])"
                            :class="scope.row.colorRed === '5' ? 'colorRedss' : ''">
                            <el-checkbox-button v-for="(item, index) in scope.row.paramDataList" :key="index"
                              :label="item.value">{{ item.name }}</el-checkbox-button>
                          </el-checkbox-group>
                        </div>
                        <el-select v-else filterable clearable style="width:100%;" :multiple="scope.row.paramList"
                          v-model.trim="infoReviewForm[scope.row.paramName]" placeholder="请选择"
                          @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName])"
                          :class="scope.row.colorRed === '5' ? 'colorRedss' : ''">
                          <el-option v-for="(item, index) in scope.row.paramDataList" :key="index" :label="item.name"
                            :value="item.value">
                          </el-option>
                        </el-select>
                      </div>
                      <div v-else>
                        <div v-if="scope.row.paramList">
                          <template>
                            <div style="display: flex;align-items: center;margin-bottom:5px;"
                              v-for="(item, index) in infoReviewForm[scope.row.paramName]" :key="index">
                              <el-input :clearable='!scope.row.paramRequired' style="width:100%;"
                                @change="validating($event, scope, scope.row.paramMinmum, scope.row.paramMaxmum, scope.row.paramPatterns, 'string', infoReviewForm[scope.row.paramName])"
                                placeholder="请输入" v-if="scope.row.paramType == 'string'"
                                v-model.trim="infoReviewForm[scope.row.paramName][index]" type="input"></el-input>
                              <el-input-number :clearable='!scope.row.paramRequired' style="width:100%;"
                                @change="validating($event, scope, scope.row.paramMinmum, scope.row.paramMaxmum, scope.row.paramPatterns, 'number', infoReviewForm[scope.row.paramName])"
                                placeholder="请输入" v-if="scope.row.paramType == 'number'"
                                v-model="infoReviewForm[scope.row.paramName][index]" type="number"
                                :step="scope.row.paramMutipleof ? scope.row.paramMutipleof : 1"
                                step-strictly></el-input-number>
                              <span class="el-icon-minus"
                                style="font-size:18px;margin:0 5px;cursor: pointer;color:#f60;"
                                @click="infoReviewForm[scope.row.paramName].splice(index, 1)"></span>
                            </div>
                            <span class="el-icon-plus" style="font-size:18px;cursor: pointer;color:#409EFF;"
                              @click="infoReviewForm[scope.row.paramName].push(null)"></span>
                          </template>
                        </div>
                        <div v-else>
                          <template
                            v-if="scope.row.paramType == 'string' && (scope.row.paramStyle == 'datetime' || scope.row.paramStyle == 'date')">
                            <el-date-picker v-model="value1" type="datetime" placeholder="选择日期时间"
                              v-model.trim="infoReviewForm[scope.row.paramName]" style="width:100%"
                              @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName])"
                              v-if="scope.row.paramStyle == 'datetime'">
                            </el-date-picker>
                            <el-date-picker v-model="value1" type="date" placeholder="选择日期"
                              v-model.trim="infoReviewForm[scope.row.paramName]" style="width:100%"
                              @change="tureOrFalseCoant($event, scope, infoReviewForm[scope.row.paramName])"
                              v-if="scope.row.paramStyle == 'date'" value-format='yyyy-MM-dd'>
                            </el-date-picker>
                          </template>
                          <div v-else>
                            <div v-if="scope.row.paramType == 'string'">
                              <template v-if="scope.row.paramStyle == 'json'">
                                <CodeEditor v-model="infoReviewForm[scope.row.paramName]" :auto-format="true"
                                  :smart-indent="true" theme="dracula" :indent-unit="4" :line-wrap="false" ref="editor">
                                </CodeEditor>
                              </template>
                              <template v-else>
                                <!-- 如果不是list且是string则判断是什么框 -->
                                <el-input style="width:100%;" clearable placeholder="请输入"
                                  v-if="scope.row.paramStyle == 'password'" show-password
                                  :maxlength="scope.row.paramMaxLength"
                                  @change="validating($event, scope, scope.row.paramMinLength, scope.row.paramMaxLength, scope.row.paramPatterns, 'string')"
                                  v-model.trim="infoReviewForm[scope.row.paramName]" type="text"
                                  :class="scope.row.colorRed === '2' || scope.row.colorRed === '3' || scope.row.colorRed === '4' || scope.row.colorRed === '5' ? 'colorReds' : ''"
                                  show-word-limit></el-input>
                                <el-input style="width:100%;" clearable placeholder="请输入"
                                  v-if="scope.row.paramStyle == 'secret'" :maxlength="scope.row.paramMaxLength"
                                  @change="validating($event, scope, scope.row.paramMinLength, scope.row.paramMaxLength, scope.row.paramPatterns, 'string')"
                                  v-model.trim="infoReviewForm[scope.row.paramName]" type="text"
                                  :class="scope.row.colorRed === '2' || scope.row.colorRed === '3' || scope.row.colorRed === '4' || scope.row.colorRed === '5' ? 'colorReds' : ''"></el-input>
                                <el-input style="width:100%;margin-top:5px" clearable placeholder="请输入"
                                  v-if="scope.row.paramStyle == 'textarea' || scope.row.paramStyle == 'text'"
                                  :maxlength="scope.row.paramMaxLength"
                                  @change="validating($event, scope, scope.row.paramMinLength, scope.row.paramMaxLength, scope.row.paramPatterns, 'string')"
                                  v-model="infoReviewForm[scope.row.paramName]" type="textarea"
                                  :class="scope.row.colorRed === '2' || scope.row.colorRed === '3' || scope.row.colorRed === '4' || scope.row.colorRed === '5' ? 'colorReds' : ''"
                                  :autosize="{ minRows: 3, maxRows: 6 }" show-word-limit></el-input>
                                <el-input style="width:100%;border:none" disabled
                                  v-if="scope.row.paramStyle == 'description'" placeholder="请输入"
                                  v-model="scope.row.paramDefault" class="noneBorder"></el-input>
                                <el-input style="width:100%;" clearable placeholder="请输入"
                                  v-if="!scope.row.paramStyle || scope.row.paramStyle == ''"
                                  :maxlength="scope.row.paramMaxLength"
                                  @change="validating($event, scope, scope.row.paramMinLength, scope.row.paramMaxLength, scope.row.paramPatterns, 'string')"
                                  v-model.trim="infoReviewForm[scope.row.paramName]" type="text"
                                  :class="scope.row.colorRed === '2' || scope.row.colorRed === '3' || scope.row.colorRed === '4' || scope.row.colorRed === '5' ? 'colorReds' : ''"
                                  show-word-limit></el-input>
                              </template>
                            </div>
                            <div v-else>
                              <template v-if="scope.row.paramStyle == 'description'">
                                <el-input style="width:100%;border:none" clearable disabled
                                  v-if="scope.row.paramType == 'number' && !scope.row.paramMutipleof" placeholder="请输入"
                                  v-model="scope.row.paramDefault" class="noneBorder"></el-input>
                                <el-input-number style="width:100%;" clearable v-else disabled placeholder="请输入"
                                  v-model.number="infoReviewForm[scope.row.paramName]" type="number"
                                  @change="validating($event, scope, scope.row.paramMinmum, scope.row.paramMaxmum, scope.row.paramPatterns, 'number')"
                                  :class="scope.row.colorRed === '2' || scope.row.colorRed === '3' || scope.row.colorRed === '4' || scope.row.colorRed === '5' ? 'colorReds noneBorder' : 'noneBorder'"
                                  :step="scope.row.paramMutipleof" step-strictly></el-input-number>
                              </template>
                              <template v-else>
                                <el-input type="number" style="width:100%;border:none" clearable
                                  v-if="scope.row.paramType == 'number' && !scope.row.paramMutipleof" placeholder="请输入"
                                  v-model.number="infoReviewForm[scope.row.paramName]"></el-input>
                                <el-input-number style="width:100%;" v-else clearable placeholder="请输入"
                                  v-model.number="infoReviewForm[scope.row.paramName]" type="number"
                                  @change="validating($event, scope, scope.row.paramMinmum, scope.row.paramMaxmum, scope.row.paramPatterns, 'number')"
                                  :class="scope.row.colorRed === '2' || scope.row.colorRed === '3' || scope.row.colorRed === '4' || scope.row.colorRed === '5' ? 'colorReds' : ''"
                                  :step="scope.row.paramMutipleof" step-strictly></el-input-number>
                              </template>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                    <template>
                      <div
                        v-if="scope.row.paramType == 'string' && (scope.row.paramStyle == 'textarea' || scope.row.paramStyle == 'text')">
                        <div v-if="scope.row.paramDescription && scope.row.colorRed === '1'"
                          style="position: relative;font-size: 12px; line-height: 10px;bottom: -6px;color: #c1c1c1;"
                          v-html="scope.row.paramDescription"></div>
                        <div
                          v-if="scope.row.colorRed === '2' || scope.row.colorRed === '3' || scope.row.colorRed === '4' || scope.row.colorRed === '5'"
                          style="position: relative;font-size: 12px;line-height: 10px;bottom: -6px;color: red;">
                          {{ scope.row.paramConstraintDescription ? scope.row.paramConstraintDescription : scope.row.customDescription }}
                        </div>
                      </div>
                      <div v-else>
                        <div v-if="scope.row.paramDescription && scope.row.colorRed === '1'"
                          style="position: relative;font-size: 12px;line-height: 12px;bottom:-2px;color: #c1c1c1;"
                          v-html="scope.row.paramDescription"></div>
                        <div
                          v-if="scope.row.colorRed === '2' || scope.row.colorRed === '3' || scope.row.colorRed === '4' || scope.row.colorRed === '5'"
                          style="position: relative;font-size: 12px;line-height: 12px;bottom:-2px;color: red;">
                          {{ scope.row.paramConstraintDescription ? scope.row.paramConstraintDescription : scope.row.customDescription }}
                        </div>
                      </div>
                    </template>
                  </el-form-item>
                  <i v-if="scope.row.paramStyle == 'radio'" class="el-icon-circle-close"
                    style="width: 30px;height: 30px;line-height: 30px;font-size: 20px;cursor: pointer;margin-top: 4px;text-align: center;"
                    @click="clearableRadioBtn(scope.row, scope, infoReviewForm[scope.row.paramName], null, null, null, infoReviewForm, scope.row.paramName)"></i>
                  <span
                    v-if="scope.row.paramreference && scope.row.paramreference != 'FormTemplate' && scope.row.paramreference != 'User' && scope.row.paramStyle == 'radio'"
                    style="width: 30px;height: 30px;line-height: 30px;font-size: 20px;cursor: pointer;margin-top: 4px;text-align: center;"
                    class="el-icon-refresh"
                    @click="focusSelectValue(scope, scope.row.paramName, infoReviewForm)"></span>
                  <span v-else style="width: 30px;height: 30px;"></span>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-form>
        <div style="text-align:right;padding-top:10px;" v-if="newTaskName != '弹性资源模板查看'">
          <el-button plain type="primary" @click="subAddModel(false)"
            v-if="detailOrLook && users.indexOf(usernameAll) != -1 && (jobIdStatus && jobIdStatus == '审核中') && !copyOrupdateType"
            :loading="reviewload">更 新</el-button>
          <el-button plain type="primary" @click="subAddModel(false)"
            v-if="detailOrLook && (jobIdStatus && jobIdStatus != '进行中') && copyOrupdateType"
            :loading="reviewload">提 交</el-button>
          <el-button plain type="success" @click="subAddModel(true)" :loading="reviewload"
            v-if="detailOrLook && LeaderType && (jobIdStatus && jobIdStatus == '审核中') && !copyOrupdateType">流 转</el-button>
          <el-button plain type="danger" @click="subAddModelDeny"
            v-if="detailOrLook && LeaderType && (jobIdStatus && jobIdStatus == '审核中') && !copyOrupdateType"
            :loading="reviewload">拒 绝</el-button>
          <el-button plain @click="reviewStatusNew = false" v-if='detailOrLook && users.indexOf(usernameAll) != -1'
            :loading="reviewload">取 消</el-button>
          <el-button @click="reviewStatusNew = false" v-if='!detailOrLook' :loading="reviewload">取 消</el-button>
          <el-button type="primary" @click="subAddModel(false, true)" v-if="!detailOrLook"
            :loading="reviewload">{{(newTaskName  == '资源模板编辑' || newTaskName == '弹性资源模板编辑') ? '保 存' : '立即提交' }}</el-button>
        </div>
      </div>
    </el-dialog>
    <!-- 选择服务器弹框 -->
    <el-dialog :title="serverTitle" :visible.sync="choseServerLog" :close-on-click-modal="false" class="diaBoxFormMain"
      width="60%">
      <el-dialog :modal='false' :show-close='false' :close-on-press-escape='true' custom-class='advanced_searchFormTem'
        :visible.sync="poHide" style="float:left;margin-left:15px">
        <div style="background:#f8f9fa;padding:10px;border-radius:5px">
          <span v-for="(Onetem, leII) in advSearch" :key="leII">
            <advancedSearchOne :searchMoreArr='searchMoreArr' :tableHeaderArr='tableHeaderArr' :searchLevel='1'
              :searchLevelData='Onetem' :key="'aa' + leII" :ownIndex='leII' @addNewsearch="addNewsearch"
              @delNewsearch='delNewsearch' :oneIndex='leII' :advSearch='advSearch' :parentIndex='null'
              style="margin-top:10px" v-if="!Onetem.child || Onetem.child.length == 0"></advancedSearchOne>
            <template v-if="Onetem.child && Onetem.child.length != 0">
              <div v-if='advSearch[0].category'
                style="border:1px solid #e1e1e1;position:relative;margin:10px 0;font-size:10px">
                <span
                  style="position:absolute;top:-5px;left:-1px;background:#f8f9fa;font-weight: 600;">{{advSearch[0].category }}</span>
              </div>
              <el-card style="margin-top:15px" class="gjSearchBody">
                <div class='gjSearchBodyItem'>
                  <div v-for="(levItem, levI) in Onetem.child" :key="levI">
                    <advancedSearchOne :searchMoreArr='searchMoreArr' :tableHeaderArr='tableHeaderArr' :searchLevel='2'
                      :searchLevelData='levItem' :ownIndex='levI' :oneIndex='leII' @addNewsearch="addNewsearch"
                      @delNewsearch='delNewsearch' :parentIndex='leII' :advSearch='advSearch'
                      v-if="!levItem.child || levItem.child.length == 0" style="margin-top:10px"></advancedSearchOne>
                    <div v-if="levItem.child && levItem.child.length != 0">
                      <div v-if="advSearch[leII].child[0].category"
                        style="border:1px solid #e1e1e1;position:relative;margin:10px 0;font-size:10px">
                        <span
                          style="position:absolute;top:-5px;left:-1px;background:#f8f9fa;font-weight: 600">{{ advSearch[leII].child[0].category }}</span>
                      </div>
                      <el-card style="margin-top:15px" class="gjSearchBody">
                        <div class='gjSearchBodyItem'>
                          <div v-for="(levItem2, levI2) in levItem.child" :key="levI2">
                            <advancedSearchOne :searchMoreArr='searchMoreArr' :tableHeaderArr='tableHeaderArr' levI
                              levI2 :searchLevel='3' :searchLevelData='levItem2' :ownIndex='levI2' :parentIndex='levI'
                              @addNewsearch="addNewsearch" :oneIndex='leII' @delNewsearch='delNewsearch'
                              :advSearch='advSearch' style="margin-top:10px"></advancedSearchOne>
                          </div>
                        </div>
                      </el-card>
                    </div>
                  </div>
                </div>
              </el-card>
            </template>
          </span>
          <div style="text-align:right;margin-top:15px">
            <el-button @click="resetSearch">重 置</el-button>
            <el-button type="primary" plain @click="searchTableServer('search')">搜 索</el-button>
          </div>
        </div>
      </el-dialog>
      <div style="border:1px solid #e1e1e1;height:28px;border-radius:4px;width:650px">
        <el-button
          style="float:left;padding: 6px 10px !important;background: #f5f7fa;color: #409EFF;    border: none;position: relative;z-index: 1;border-radius: 4px 0 0 4px;height:26px;border-right:1px solid #e1e1e1;"
          icon="el-icon-plus" @click="searchShowStatus1">高级搜索</el-button>
        <el-input v-model="searchFormServe" @keydown.enter.native="searchTableServer('search')"
          style="width:560px;float: left;height:26px" class="batchSearchClass" placeholder="请输入关键字(长度大于2)"
          clearable></el-input>
      </div>
      <div v-if="searchObjTag" style="margin:10px 0 10px;cursor: pointer;">
        <div v-if='searchObjTag.length' style='position: relative;float: left;'>

          <div v-for='(tagsValitem, tagsValindex) in searchObjTag' :key='tagsValindex'
            style="margin-right:20px;background:#ecf5ff;padding: 5px;float:left;margin-bottom:10px">
            <span :key="index + 1" v-for="(item, index) in tagsValitem" @click="enterTagLast(tagsValindex)">
              <span v-if="!item.child || !item.child.length">
                <span v-if="(tagsValitem[0] && tagsValitem[0].category) && (index != 0)"
                  class="tagTypeClass">{{tagsValitem[0].category }}</span>
                <span class="speical_tag" slot="reference" :disable-transitions="false"
                  style="margin:0 3px;cursor: pointer;">
                  <span v-if="item.type == 'select'" class="widthwarp">
                    {{ item.name }}:{{ item.value }}
                  </span>
                  <span v-if="item.type == '_GT'" class="widthwarp">
                    {{ item.name }}>{{ item.value }}
                  </span>
                  <span v-if="item.type == '_LT'" class="widthwarp">
                    {{ item.name + "&lt;" + item.value }}
                  </span>
                  <span v-if="item.type == ''" class="widthwarp">
                    {{ item.name }}:{{ item.value }}
                  </span>
                  <span v-if="item.type == '_REGEX'" class="widthwarp">
                    {{ item.name }}:/{{ item.value }}/
                  </span>
                  <span v-if="item.type == '_HGT'" class="widthwarp">
                    {{ item.name }}>{{ new Date(item.value).getTime() | filterTimeShow }}
                  </span>
                  <span v-if="item.type == '_DGT'" class="widthwarp">
                    {{ item.name }}>{{ new Date(item.value).getTime() | filterTimeShow('YYYY-DD-MM') }}
                  </span>
                  <span v-if="item.type == '_HLT'" class="widthwarp">
                    {{ item.name }}&lt;{{ new Date(item.value).getTime() | filterTimeShow }}
                  </span>
                  <span v-if="item.type == '_DLT'" class="widthwarp">
                    {{ item.name }}&lt;{{ new Date(item.value).getTime() | filterTimeShow('YYYY-DD-MM') }}
                  </span>
                  <span v-if="item.type == 'arr'" style="top:-4px;position: relative;">
                    <span v-if="item.value && item.value.split('\n').length > 8">
                      {{ item.name + ':[' + item.value.split('\n')[0] + ',' +
      item.value.split('\n')[1] + '等' + item.value.split('\n').length + '个]' }}
                    </span>
                    <span v-else>
                      {{ item.name + ':[' + item.value + ']' }}
                    </span>
                  </span>
                  <span v-if="item.type == 'selectNull'" class="widthwarp">
                    {{ item.name }}
                    <span v-if="item.value != 'no'">!</span>
                  </span>

                </span>
              </span>

              <span v-else>
                <span class="tagTypeClass"><span style='margin-right:3px'>{{ tagsValitem[0].category }}</span>(</span>
                <span :key="index1 + 1" v-for="(itemTwo, index1) in item.child">
                  <span v-if="!itemTwo.child || !itemTwo.child.length">
                    <span v-if="item.child[0].category && index1 != 0" class="tagTypeClass">{{ item.child[0].category
                      }}</span>
                    <span class="speical_tag" slot="reference" :disable-transitions="false"
                      style="margin:0 3px;cursor: pointer;">
                      <span v-if="itemTwo.type == 'select'" class="widthwarp">
                        {{ itemTwo.name }}:{{ itemTwo.value }}
                      </span>
                      <span v-if="itemTwo.type == '_GT'" class="widthwarp">
                        {{ itemTwo.name }}>{{ itemTwo.value }}
                      </span>
                      <span v-if="itemTwo.type == '_LT'" class="widthwarp">
                        {{ itemTwo.name + "&lt;" + itemTwo.value }}
                      </span>
                      <span v-if="itemTwo.type == ''" class="widthwarp">
                        {{ itemTwo.name }}:{{ itemTwo.value }}
                      </span>
                      <span v-if="itemTwo.type == '_REGEX'" class="widthwarp">
                        {{ itemTwo.name }}:/{{ itemTwo.value }}/
                      </span>
                      <span v-if="itemTwo.type == '_HGT'" class="widthwarp">
                        {{ itemTwo.name }}>{{ new Date(itemTwo.value).getTime() | filterTimeShow }}
                      </span>
                      <span v-if="itemTwo.type == '_DGT'" class="widthwarp">
                        {{ itemTwo.name }}>{{ new Date(itemTwo.value).getTime() | filterTimeShow('YYYY-DD-MM') }}
                      </span>
                      <span v-if="itemTwo.type == '_HLT'" class="widthwarp">
                        {{ itemTwo.name }}&lt;{{ new Date(itemTwo.value).getTime() | filterTimeShow }}
                      </span>
                      <span v-if="itemTwo.type == '_DLT'" class="widthwarp">
                        {{ itemTwo.name }}&lt;{{ new Date(itemTwo.value).getTime() | filterTimeShow('YYYY-DD-MM') }}
                      </span>
                      <span v-if="itemTwo.type == 'arr'" style="top:-4px;position: relative;">
                        <span v-if="itemTwo.value && itemTwo.value.split('\n').length > 8">
                          {{ itemTwo.name + ':[' + itemTwo.value.split('\n')[0] + ',' +
      itemTwo.value.split('\n')[1] + '等' + itemTwo.value.split('\n').length + '个]' }}
                        </span>
                        <span v-else>
                          {{ itemTwo.name + ':[' + itemTwo.value + ']' }}
                        </span>
                      </span>
                      <span v-if="itemTwo.type == 'selectNull'" class="widthwarp">
                        {{ itemTwo.name }}
                        <span v-if="itemTwo.value != 'no'">!</span>
                      </span>
                    </span>
                  </span>
                  <span v-else>
                    <span class="tagTypeClass"><span style='margin-right:3px'>{{ item.child[0].category
                        }}</span>(</span>
                    <span :key="index3 + 1" v-for="(itemThere, index3) in itemTwo.child">
                      <span v-if="itemTwo.child[0].category && index3 != 0" class="tagTypeClass">{{
      itemTwo.child[0].category }}</span>
                      <span class="speical_tag" slot="reference" :disable-transitions="false"
                        style="margin:0 3px;cursor: pointer;" v-if="!itemThere.child || !itemThere.child.length">
                        <span v-if="itemThere.type == 'select'" class="widthwarp">
                          {{ itemThere.name }}:{{ itemThere.value }}
                        </span>
                        <span v-if="itemThere.type == '_GT'" class="widthwarp">
                          {{ itemThere.name }}>{{ itemThere.value }}
                        </span>
                        <span v-if="itemThere.type == '_LT'" class="widthwarp">
                          {{ itemThere.name + "&lt;" + itemThere.value }}
                        </span>
                        <span v-if="itemThere.type == ''" class="widthwarp">
                          {{ itemThere.name }}:{{ itemThere.value }}
                        </span>
                        <span v-if="itemThere.type == '_REGEX'" class="widthwarp">
                          {{ itemThere.name }}:/{{ itemThere.value }}/
                        </span>
                        <span v-if="itemThere.type == '_HGT'" class="widthwarp">
                          {{ itemThere.name }}>{{ new Date(itemThere.value).getTime() | filterTimeShow }}
                        </span>
                        <span v-if="itemThere.type == '_DGT'" class="widthwarp">
                          {{ itemThere.name }}>{{ new Date(itemThere.value).getTime() | filterTimeShow('YYYY-DD-MM') }}
                        </span>
                        <span v-if="itemThere.type == '_HLT'" class="widthwarp">
                          {{ itemThere.name }}&lt;{{ new Date(itemThere.value).getTime() | filterTimeShow }}
                        </span>
                        <span v-if="itemThere.type == '_DLT'" class="widthwarp">
                          {{ itemThere.name }}&lt;{{ new Date(itemThere.value).getTime() | filterTimeShow('YYYY-DD-MM')
                          }}
                        </span>
                        <span v-if="itemThere.type == 'arr'" style="top:-4px;position: relative;">
                          <span v-if="itemThere.value && itemThere.value.split('\n').length > 8">
                            {{ itemThere.name + ':[' + itemThere.value.split('\n')[0] + ',' +
      itemThere.value.split('\n')[1] + '等' + itemThere.value.split('\n').length + '个]' }}
                          </span>
                          <span v-else>
                            {{ itemThere.name + ':[' + itemThere.value + ']' }}
                          </span>
                        </span>
                        <span v-if="itemThere.type == 'selectNull'" class="widthwarp">
                          {{ itemThere.name }}
                          <span v-if="itemThere.value != 'no'">!</span>
                        </span>
                      </span>
                    </span>
                    <span class="tagTypeClass" style="margin-right:3px">)</span>
                  </span>
                </span>
                <span v-if="tagsValitem[0] && tagsValitem[0].category" class="tagTypeClass"
                  style="margin-right:3px">)</span>
              </span>
            </span>

            <i class='el-icon-circle-close' style="cursor: pointer; margin-left: 10px;font-size: 12px;"
              @click="refAdvancedSearch(tagsValindex)"></i>

          </div>
          <el-link v-if='searchObjTag.length' size="mini"
            style="cursor: pointer; margin-left: 10px;font-size: 12px;position: absolute;bottom: 15px;right: -15px;"
            type="primary" @click="refAdvancedSearchLast">重置</el-link>
        </div>
      </div>

      <el-card v-show="searchShowStatus" class="searchBox" style="margin:10px;">
        <el-form :model="searchObj" ref="searchObjForm" label-width="80px" class="demo-ruleForm">
          <el-row style="margin-bottom:10px">
            <el-form-item :label="serverTitle">
              <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 6 }" resize="none" v-model="batchSearchName"
                placeholder="支持名称批量搜索（换行输入）" clearable></el-input>
            </el-form-item>
          </el-row>
          <el-row v-if="showCidSearch" style="margin-bottom:10px">
            <el-form-item label="云ID">
              <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 6 }" resize="none" v-model="batchSearchName1"
                placeholder="支持云ID批量搜索（换行输入）" clearable></el-input>
            </el-form-item>
          </el-row>
        </el-form>
        <div style="text-align:right;">
          <el-button @click="resetVal">重 置</el-button>
          <el-button type="primary" plain @click="searchTableServer('search')">搜 索</el-button>
        </div>
      </el-card>

      <el-table border :data="searchHostTable" style="width: 100%;margin-top:10px" size="mini" class="RSHostTable"
        max-height="300" min-height="100" @selection-change="handleSelectionChangeSearth">
        <el-table-column type="selection" width="40"></el-table-column>
        <el-table-column v-for="(taItem, i) in referenceShowArrHeader" :key="i" :show-overflow-tooltip="true"
          :label="taItem.title" :prop="taItem.value"></el-table-column>
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
    <el-dialog :title="listTitle" :visible.sync="choseListStatus" :close-on-click-modal="false" class="diaBoxFormMain"
      :width="sreachWidth">
      <el-dialog :modal='false' :show-close='false' :close-on-press-escape='true' custom-class='advanced_searchFormTem1'
        :visible.sync="poHide" style="float:left;margin-left:15px">
        <div style="background:#f8f9fa;padding:10px;border-radius:5px">
          <span v-for="(Onetem, leII) in advSearch" :key="leII">
            <advancedSearchOne :searchMoreArr='searchMoreArr' :tableHeaderArr='tableHeaderArr' :searchLevel='1'
              :searchLevelData='Onetem' :key="'aa' + leII" :ownIndex='leII' @addNewsearch="addNewsearch"
              @delNewsearch='delNewsearch' :oneIndex='leII' :advSearch='advSearch' :parentIndex='null'
              style="margin-top:10px" v-if="!Onetem.child || Onetem.child.length == 0"></advancedSearchOne>
            <template v-if="Onetem.child && Onetem.child.length != 0">
              <div v-if='advSearch[0].category'
                style="border:1px solid #e1e1e1;position:relative;margin:10px 0;font-size:10px">
                <span style="position:absolute;top:-5px;left:-1px;background:#f8f9fa;font-weight: 600;">{{
      advSearch[0].category }}</span>
              </div>
              <el-card style="margin-top:15px" class="gjSearchBody">
                <div class='gjSearchBodyItem'>
                  <div v-for="(levItem, levI) in Onetem.child" :key="levI">
                    <advancedSearchOne :searchMoreArr='searchMoreArr' :tableHeaderArr='tableHeaderArr' :searchLevel='2'
                      :searchLevelData='levItem' :ownIndex='levI' :oneIndex='leII' @addNewsearch="addNewsearch"
                      @delNewsearch='delNewsearch' :parentIndex='leII' :advSearch='advSearch'
                      v-if="!levItem.child || levItem.child.length == 0" style="margin-top:10px"></advancedSearchOne>
                    <div v-if="levItem.child && levItem.child.length != 0">
                      <div v-if="advSearch[leII].child[0].category"
                        style="border:1px solid #e1e1e1;position:relative;margin:10px 0;font-size:10px">
                        <span style="position:absolute;top:-5px;left:-1px;background:#f8f9fa;font-weight: 600">{{
      advSearch[leII].child[0].category }}</span>
                      </div>
                      <el-card style="margin-top:15px" class="gjSearchBody">
                        <div class='gjSearchBodyItem'>
                          <div v-for="(levItem2, levI2) in levItem.child" :key="levI2">
                            <advancedSearchOne :searchMoreArr='searchMoreArr' :tableHeaderArr='tableHeaderArr' levI
                              levI2 :searchLevel='3' :searchLevelData='levItem2' :ownIndex='levI2' :parentIndex='levI'
                              @addNewsearch="addNewsearch" :oneIndex='leII' @delNewsearch='delNewsearch'
                              :advSearch='advSearch' style="margin-top:10px"></advancedSearchOne>
                          </div>
                        </div>
                      </el-card>
                    </div>
                  </div>
                </div>
              </el-card>
            </template>
          </span>

          <div style="text-align:right;margin-top:15px">
            <el-button @click="resetSearch">重 置</el-button>
            <el-button type="primary" plain @click="searchTableServer('search')">搜 索</el-button>
          </div>
        </div>
      </el-dialog>
      <div style="border:1px solid #e1e1e1;height:28px;border-radius:4px;width:650px">
        <el-button
          style="float:left;padding: 6px 10px !important;background: #f5f7fa;color: #409EFF;    border: none;position: relative;z-index: 1;border-radius: 4px 0 0 4px;height:26px;border-right:1px solid #e1e1e1;"
          icon="el-icon-plus" @click="searchShowStatus1">高级搜索</el-button>
        <el-input v-model="searchFormServe" @keydown.enter.native="searchTableServer('search')"
          style="width:560px;float: left;height:26px" class="batchSearchClass" placeholder="请输入关键字(长度大于2)"
          clearable></el-input>
      </div>
      <div v-if="searchObjTag" style="margin:10px 0 10px;cursor: pointer;">
        <div v-if='searchObjTag.length' style='position: relative;float: left;'>

          <div v-for='(tagsValitem, tagsValindex) in searchObjTag' :key='tagsValindex'
            style="margin-right:20px;background:#ecf5ff;padding: 5px;float:left;margin-bottom:10px">
            <span :key="index + 1" v-for="(item, index) in tagsValitem" @click="enterTagLast(tagsValindex)">
              <span v-if="!item.child || !item.child.length">
                <span v-if="(tagsValitem[0] && tagsValitem[0].category) && (index != 0)" class="tagTypeClass">{{
      tagsValitem[0].category }}</span>
                <span class="speical_tag" slot="reference" :disable-transitions="false"
                  style="margin:0 3px;cursor: pointer;">
                  <span v-if="item.type == 'select'" class="widthwarp">
                    {{ item.name }}:{{ item.value }}
                  </span>
                  <span v-if="item.type == '_GT'" class="widthwarp">
                    {{ item.name }}>{{ item.value }}
                  </span>
                  <span v-if="item.type == '_LT'" class="widthwarp">
                    {{ item.name + "&lt;" + item.value }}
                  </span>
                  <span v-if="item.type == ''" class="widthwarp">
                    {{ item.name }}:{{ item.value }}
                  </span>
                  <span v-if="item.type == '_REGEX'" class="widthwarp">
                    {{ item.name }}:/{{ item.value }}/
                  </span>
                  <span v-if="item.type == '_HGT'" class="widthwarp">
                    {{ item.name }}>{{ new Date(item.value).getTime() | filterTimeShow }}
                  </span>
                  <span v-if="item.type == '_DGT'" class="widthwarp">
                    {{ item.name }}>{{ new Date(item.value).getTime() | filterTimeShow('YYYY-DD-MM') }}
                  </span>
                  <span v-if="item.type == '_HLT'" class="widthwarp">
                    {{ item.name }}&lt;{{ new Date(item.value).getTime() | filterTimeShow }}
                  </span>
                  <span v-if="item.type == '_DLT'" class="widthwarp">
                    {{ item.name }}&lt;{{ new Date(item.value).getTime() | filterTimeShow('YYYY-DD-MM') }}
                  </span>
                  <span v-if="item.type == 'arr'" style="top:-4px;position: relative;">
                    <span v-if="item.value && item.value.split('\n').length > 8">
                      {{ item.name + ':[' + item.value.split('\n')[0] + ',' +
      item.value.split('\n')[1] + '等' + item.value.split('\n').length + '个]' }}
                    </span>
                    <span v-else>
                      {{ item.name + ':[' + item.value + ']' }}
                    </span>
                  </span>
                  <span v-if="item.type == 'selectNull'" class="widthwarp">
                    {{ item.name }}
                    <span v-if="item.value != 'no'">!</span>
                  </span>

                </span>
              </span>

              <span v-else>
                <span class="tagTypeClass"><span style='margin-right:3px'>{{ tagsValitem[0].category }}</span>(</span>
                <span :key="index1 + 1" v-for="(itemTwo, index1) in item.child">
                  <span v-if="!itemTwo.child || !itemTwo.child.length">
                    <span v-if="item.child[0].category && index1 != 0" class="tagTypeClass">{{ item.child[0].category
                      }}</span>
                    <span class="speical_tag" slot="reference" :disable-transitions="false"
                      style="margin:0 3px;cursor: pointer;">
                      <span v-if="itemTwo.type == 'select'" class="widthwarp">
                        {{ itemTwo.name }}:{{ itemTwo.value }}
                      </span>
                      <span v-if="itemTwo.type == '_GT'" class="widthwarp">
                        {{ itemTwo.name }}>{{ itemTwo.value }}
                      </span>
                      <span v-if="itemTwo.type == '_LT'" class="widthwarp">
                        {{ itemTwo.name + "&lt;" + itemTwo.value }}
                      </span>
                      <span v-if="itemTwo.type == ''" class="widthwarp">
                        {{ itemTwo.name }}:{{ itemTwo.value }}
                      </span>
                      <span v-if="itemTwo.type == '_REGEX'" class="widthwarp">
                        {{ itemTwo.name }}:/{{ itemTwo.value }}/
                      </span>
                      <span v-if="itemTwo.type == '_HGT'" class="widthwarp">
                        {{ itemTwo.name }}>{{ new Date(itemTwo.value).getTime() | filterTimeShow }}
                      </span>
                      <span v-if="itemTwo.type == '_DGT'" class="widthwarp">
                        {{ itemTwo.name }}>{{ new Date(itemTwo.value).getTime() | filterTimeShow('YYYY-DD-MM') }}
                      </span>
                      <span v-if="itemTwo.type == '_HLT'" class="widthwarp">
                        {{ itemTwo.name }}&lt;{{ new Date(itemTwo.value).getTime() | filterTimeShow }}
                      </span>
                      <span v-if="itemTwo.type == '_DLT'" class="widthwarp">
                        {{ itemTwo.name }}&lt;{{ new Date(itemTwo.value).getTime() | filterTimeShow('YYYY-DD-MM') }}
                      </span>
                      <span v-if="itemTwo.type == 'arr'" style="top:-4px;position: relative;">
                        <span v-if="itemTwo.value && itemTwo.value.split('\n').length > 8">
                          {{ itemTwo.name + ':[' + itemTwo.value.split('\n')[0] + ',' +
      itemTwo.value.split('\n')[1] + '等' + itemTwo.value.split('\n').length + '个]' }}
                        </span>
                        <span v-else>
                          {{ itemTwo.name + ':[' + itemTwo.value + ']' }}
                        </span>
                      </span>
                      <span v-if="itemTwo.type == 'selectNull'" class="widthwarp">
                        {{ itemTwo.name }}
                        <span v-if="itemTwo.value != 'no'">!</span>
                      </span>
                    </span>
                  </span>
                  <span v-else>
                    <span class="tagTypeClass"><span style='margin-right:3px'>{{ item.child[0].category
                        }}</span>(</span>
                    <span :key="index3 + 1" v-for="(itemThere, index3) in itemTwo.child">
                      <span v-if="itemTwo.child[0].category && index3 != 0" class="tagTypeClass">{{
      itemTwo.child[0].category }}</span>
                      <span class="speical_tag" slot="reference" :disable-transitions="false"
                        style="margin:0 3px;cursor: pointer;" v-if="!itemThere.child || !itemThere.child.length">
                        <span v-if="itemThere.type == 'select'" class="widthwarp">
                          {{ itemThere.name }}:{{ itemThere.value }}
                        </span>
                        <span v-if="itemThere.type == '_GT'" class="widthwarp">
                          {{ itemThere.name }}>{{ itemThere.value }}
                        </span>
                        <span v-if="itemThere.type == '_LT'" class="widthwarp">
                          {{ itemThere.name + "&lt;" + itemThere.value }}
                        </span>
                        <span v-if="itemThere.type == ''" class="widthwarp">
                          {{ itemThere.name }}:{{ itemThere.value }}
                        </span>
                        <span v-if="itemThere.type == '_REGEX'" class="widthwarp">
                          {{ itemThere.name }}:/{{ itemThere.value }}/
                        </span>
                        <span v-if="itemThere.type == '_HGT'" class="widthwarp">
                          {{ itemThere.name }}>{{ new Date(itemThere.value).getTime() | filterTimeShow }}
                        </span>
                        <span v-if="itemThere.type == '_DGT'" class="widthwarp">
                          {{ itemThere.name }}>{{ new Date(itemThere.value).getTime() | filterTimeShow('YYYY-DD-MM') }}
                        </span>
                        <span v-if="itemThere.type == '_HLT'" class="widthwarp">
                          {{ itemThere.name }}&lt;{{ new Date(itemThere.value).getTime() | filterTimeShow }}
                        </span>
                        <span v-if="itemThere.type == '_DLT'" class="widthwarp">
                          {{ itemThere.name }}&lt;{{ new Date(itemThere.value).getTime() | filterTimeShow('YYYY-DD-MM')
                          }}
                        </span>
                        <span v-if="itemThere.type == 'arr'" style="top:-4px;position: relative;">
                          <span v-if="itemThere.value && itemThere.value.split('\n').length > 8">
                            {{ itemThere.name + ':[' + itemThere.value.split('\n')[0] + ',' +
      itemThere.value.split('\n')[1] + '等' + itemThere.value.split('\n').length + '个]' }}
                          </span>
                          <span v-else>
                            {{ itemThere.name + ':[' + itemThere.value + ']' }}
                          </span>
                        </span>
                        <span v-if="itemThere.type == 'selectNull'" class="widthwarp">
                          {{ itemThere.name }}
                          <span v-if="itemThere.value != 'no'">!</span>
                        </span>
                      </span>
                    </span>
                    <span class="tagTypeClass" style="margin-right:3px">)</span>
                  </span>
                </span>
                <span v-if="tagsValitem[0] && tagsValitem[0].category" class="tagTypeClass"
                  style="margin-right:3px">)</span>
              </span>
            </span>
            <i class='el-icon-circle-close' style="cursor: pointer; margin-left: 10px;font-size: 12px;"
              @click="refAdvancedSearch(tagsValindex)"></i>
          </div>
          <el-link v-if='searchObjTag.length' size="mini"
            style="cursor: pointer; margin-left: 10px;font-size: 12px;position: absolute;bottom: 15px;right: -15px;"
            type="primary" @click="refAdvancedSearchLast">重置</el-link>
        </div>
      </div>
      <el-card v-show="searchShowList" class="searchBox" style="margin:10px;">
        <el-form :model="searchObj" ref="searchObjForm" label-width="80px" class="demo-ruleForm">
          <el-row style="margin-bottom:10px">
            <el-form-item :label="listTitle">
              <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 6 }" resize="none" v-model="batchSearchName"
                placeholder="支持名称批量搜索（换行输入）" clearable></el-input>
            </el-form-item>
          </el-row>
          <el-row v-if="showCidSearch" style="margin-bottom:10px">
            <el-form-item label="云ID">
              <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 6 }" resize="none" v-model="batchSearchName1"
                placeholder="支持云ID批量搜索（换行输入）" clearable></el-input>
            </el-form-item>
          </el-row>
        </el-form>
        <div style="text-align:right;">
          <el-button @click="resetVal">重 置</el-button>
          <el-button type="primary" plain @click="searchTableServer('search')">搜 索</el-button>
        </div>
      </el-card>
      <el-table border :data="searchHostTable" style="width: 100%;margin-top:10px" size="mini" class="RSHostTable"
        max-height="300" min-height="100" @selection-change="handleSelectionChangeSearth">
        <el-table-column type="selection" width="40"></el-table-column>
        <el-table-column v-for="(taItem, i) in referenceShowListArrHeader" :key="i" :show-overflow-tooltip="true"
          :label="taItem.title" :prop="taItem.value"></el-table-column>
      </el-table>
      <div style="text-align: center; margin-top: 15px; margin-bottom: 15px;">
        <el-pagination :pager-count="5" @size-change="handleSizeChange" @current-change="changehandleCurrent"
          :page-sizes="[10, 20, 30, 50, 100, 150, 200]" layout="total,sizes, prev, pager, next"
          :current-page="searchpageNum" :page-size="searchpagesize" :total="searchtotalLen">
        </el-pagination>
      </div>
      <div style="text-align:right;padding-top:10px;">
        <el-button @click="choseListStatus = false">取 消</el-button>
        <el-button type="primary" @click="choseListBtn">确 定</el-button>
      </div>
    </el-dialog>
    <!-- 批量导入弹框 -->
    <el-dialog title="批量导入" :visible.sync="batchIn" :close-on-click-modal="false" class="diaBoxFormMain" width="50%">
      <el-input type="textarea" :autosize="{ minRows: 4, maxRows: 6 }" placeholder="请输入内容" v-model="batchTextarea"
        clearable>
      </el-input>
      <div style="text-align:right;padding-top:10px;">
        <el-button @click="batchTextarea = ''">重 置</el-button>
        <el-button type="primary" @click="batchInBtn">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="newTaskName" :visible.sync="batchCearteStatus" :close-on-click-modal="false" class="diaBoxCearte"
      style="margin-top:10vh;" width="600px">
      <div>
        <el-form ref="formPost" :rules="rulesPost" :model="formPost" label-width="80px">
          <el-form-item label="名称" prop="name">
            <el-input v-model="formPost.name" placeholder="请输入" clearable></el-input>
          </el-form-item>
          <el-form-item label="项目" prop="project" v-if="false">
            <el-select style="width:100%;" filterable v-model="formPost.project"
              placeholder="请选择项目">
              <el-option v-for="(item, index) in projecArr" :label="item.name" :value="item.value"
                :key="index"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="配额" prop="cloudQuotaConfig"
            v-if="newTaskName == '弹性资源模板编辑' || newTaskName == '弹性资源模板新建' || this.newTaskName == '弹性资源模板改配'">
            <el-radio-group v-model.trim="formPost.cloudQuotaConfig" size="mini">
              <el-radio-button v-for="(itemOne, eeeindex) in cloudQuotaConfigArr" :key="eeeindex"
                :label="itemOne.value">
                <template>
                  <el-tooltip class="item" effect="dark" placement="top">
                    <div slot="content">
                      <div style="text-align:center;">{{ itemOne.name }}</div>
                      {{ itemOne.content }}
                    </div>
                    <div>
                      {{ itemOne.name }}
                    </div>
                  </el-tooltip>
                </template>
              </el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="描述" prop="remark">
            <el-input type="textarea" :autosize="{ minRows: 3, maxRows: 5 }" placeholder="请输入" resize="none"
              v-model="formPost.remark" clearable>
            </el-input>
          </el-form-item>
        </el-form>
        <div style="text-align:right;">
          <el-button @click="batchCearteStatus = false">取消</el-button>
          <el-button type="primary" :loading="reviewload">确认</el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>
<script>
import { mapGetters } from "vuex";
import { advSearchHandle } from '@/components/mixin/advSearchHandle'
import { userMixin } from '@/components/mixin/user'
import { formTemp } from '@/components/mixin/formTemp'
import formTemChild from "@/views/components/formTempChild"
import Cookies from 'js-cookie'
import { CodeEditor } from 'bin-code-editor'
import advancedSearchOne from '@/views/components/advancedSearchOne'
//例如：import 《组件名称》 from '《组件路径》';
export default {
  mixins: [advSearchHandle, userMixin,  formTemp],
  components: {
    formTemChild, CodeEditor,
    advancedSearchOne
  },
  props: {
    //查看单行数据，任务详情，任务信息中的任务表单查看
    historyTableList: {
      type: Object,
      default: {}
    },
    baseConfigData: {
      type: Object,
      default: {}
    },
    //弹框状态
    reviewStatus: {
      type: Boolean,
      default: false
    },
    //false是新建任务  true是编辑
    detailOrLook: {
      type: Boolean,
      default: false,
    },
    //false任务  true是项目/项目标签
    detailOrsigle: {
      type: String,
    },
    //弹框中文名称
    newTaskName: {
      type: String,
      default: ''
    },
    //任务名称
    nodeCheckName: {
      type: String,
      default: ''
    },
    childNodename: {
      type: String,
      default: ''
    },
    jobIdStatus: {
      type: String,
      default: ''
    },
    FormTemplatjobMode: {
      type: String,
      default: ''
    },
    checkList: {
      type: Array,
      default: []
    },
  },
  data() {
    //这里存放数据
    return {
      checkListLast: [],
      serverLogStatus: false,
      idChangeNameObj: {},
      LeaderType: false,
      copyOrupdateType: false,
      mainOperandName: '',
      mainOperandvsrefrename: '',
      paramreferenceShowNew1: [],
      paramreferenceShowNew2: [],
      createChildARRName3: '',
      typeCheckArr3: [],
      gaipeiArr3: [],
      waiIndex: '',
      secondName: '',
      nodeCheckNameThird: '',
      toJumpDetailID: '',
      lowerName: '',
      widthPro: '',
      sreachWidth: '',
      usernameAll: '',
      batchTextarea: '',
      batchIndex: '',
      batchIn: false,
      typeCheckArr: [],
      serverTitle: '',
      serverTypeIndex: "",
      childRowTitle1: "",
      childRowTitle: "",
      serverTypeParentIndex: "",
      detailOrLookThis: this.detailOrLook,
      infoReviewForm: {},
      infoReviewAttr: [],
      reviewStatusNew: this.reviewStatus,
      searchFormServe: '',
      searchpageNum: 1,
      searchpagesize: 20,
      searchtotalLen: 0,
      choseServerLog: false,
      reviewload: false,
      searchHostTable: [],
      referenceShowArrHeader: [],
      choseListServer: [],
      gaipeiArr: [],
      searchShowStatus: false,
      searchShowList: false,
      searchObj: {},
      batchSearchName: "",
      batchSearchName1: "",
      getstr: "",
      theadArrShow: [],
      thead1ArrShow: [],
      listTitle: "",
      listName: "",
      listObj: {},
      choseListStatus: false,
      showCidSearch: false,
      referenceShowListArrHeader: [],
      historyparamreferenceShow1: [],
      lookType: false,
      batchCearteStatus: false,
      formPost: {
        name: "",
        project: "",
        remark: ''
      },
      potsCreateObj: {},
      rulesPost: {
        name: [
          { required: true, message: '请输入名称', trigger: 'blur' },
        ],
        project: [
          { required: true, message: '请选择项目', trigger: 'change' }
        ],
      },
      projecArr: [],
      cloudQuotaConfigArr: [],
      batchSearchInput: "",
      template_id:""
    };
  },
  watch: {
    reviewStatusNew() {
      this.$emit('addClose', false)
    },
    listObj: {
      handler: () => {
      },
      deep: true
    }
  },
  computed: {
    ...mapGetters(["loginUserName", "loginUserID"]),
  },
  created() {
    if (this.$route.path.indexOf('/general/job') == -1) {
      this.widthPro = '50%'
      this.sreachWidth = '50%'
    } else {
      this.widthPro = '80%'
      this.sreachWidth = '60%'
    }
  },
  //生命周期 - 挂载完成（可以访问DOM元素）
  mounted() {
    this.usernameAll = sessionStorage.getItem("username") ? sessionStorage.getItem("username") : Cookies.get("username")
    if (this.historyTableList.dealer && this.historyTableList.dealer.length && this.historyTableList.dealer.indexOf(this.usernameAll) != -1) {
      this.LeaderType = true
    }
    if (this.historyTableList._copyOrupdate) {
      this.copyOrupdateType = true
    }
    this.checkListLast = JSON.parse(JSON.stringify(this.checkList))
    if (Object.keys(this.historyTableList).length != 0) {
      this.getNodeInfo('look')
    } else {
      this.getNodeInfo()
    }
  }
}
</script>
<style scoped lang='scss'>
.speical_tag .el-tag__close {
  position: absolute !important;
  top: 6px !important;
  right: 0 !important;
}

.speical_tag {
  color: #409eff;
  font-weight: 600;
  font-size: 12px;
  position: relative;
  top: 4px;
}

.general_box {
  width: 100%;
}
</style>
<style scoped>
.diaBox /deep/ .el-dialog__body .bin-json-editor {
  height: 800px !important;
}

.diaBox /deep/ .el-dialog__body .CodeMirror {
  height: 800px !important;
  line-height: 20px !important;
}

.batchSearchClass /deep/ .el-input__inner {
  border: none;
  height: 26px;
}

.general_box /deep/ .el-form-item {
  width: 100% !important;
}

.general_box /deep/ .infoTable .el-table__header {
  display: none;
}

.general_box /deep/ .infoTable td {
  border-bottom: none;
}

.general_box /deep/ .infoTable td td {
  border-bottom: 1px solid #ebeef5;
}

.general_box /deep/ .infohide .el-form-item__label {
  display: none;
}

.diaBox /deep/ .el-table td {
  padding: 0 !important;
}

.diaBoxFormMain /deep/ .el-table .cell,
.diaBox /deep/ .el-card .el-table .cell {
  min-height: 44px !important;
  line-height: 44px !important;
}

.diaBox /deep/ .el-form-item__label {
  font-size: 12px;
  overflow: hidden;
  white-space: nowrap;
}

.diaBox /deep/ .getQLBox .el-radio-button--mini .el-radio-button__inner {
  margin: 2px 0;
  min-height: 28px;
  width: 150px;
  border-left: 1px solid #dcdfe6;
  margin-right: -1px;
  overflow: hidden;
}

.diaBox /deep/.el-radio-button__inner .el-tooltip {
  overflow: hidden;
  text-overflow: ellipsis;
}

.getQLBoxOne {
  color: #999;
  margin-top: 3px;
  width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.getQLBox .is-active .getQLBoxOne {
  color: #fff;
  margin-top: 3px;
  position: relative;
  z-index: 999;
}

.diaBox
  /deep/
  .getQLBox
  .el-radio-button--mini.is-active
  .el-radio-button__inner {
  border-left: none;
}

.diaBox /deep/ .el-dialog__body {
  padding: 0 20px 30px 20px !important;
}

.diaBox /deep/ .el-input-number__decrease,
.diaBox /deep/.el-input-number__increase {
  top: 5px;
  height: 26px;
  line-height: 26px;
}

.diaBox /deep/ .el-input-number {
  width: 150px !important;
}

.general_box /deep/ .demo-ruleForm .specialTable .el-table__header {
  display: block;
}

.general_box /deep/ .el-dropdown-link {
  cursor: pointer;
  color: #409eff;
}

.general_box /deep/ .el-card__body {
  min-height: 40px !important;
  padding: 20px !important;
}

.diaBoxCearte /deep/ .el-dialog__body {
  padding: 0 20px 20px !important;
}

.el-tag {
  height: unset !important;
  line-height: 18px !important;
  padding: 4px 16px 4px 10px;
  max-width: 160px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  vertical-align: middle;
  position: relative;
}
</style>
<style>
.myClassMessage {
  min-width: 255px !important;
}

.demo-ruleForm .el-textarea .el-input__count {
  bottom: -8px !important;
  background: none !important;
}

.specal_select_width .el-select-dropdown__item {
  max-width: 380px;
}

.specal_slect_classs .el-select-dropdown__item {
  width: 800px;
}

.specal_slect_classs .el-select-dropdown__item {
  font-size: 14px;
  padding: 6px 20px;
  position: relative;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: #606266;
  height: auto;
  line-height: 20px !important;
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
  cursor: pointer;
}

.specal_select_width .el-select-dropdown__item {
  font-size: 14px;
  padding: 6px 20px;
  position: relative;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: #606266;
  height: auto;
  line-height: 20px !important;
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
  cursor: pointer;
}

.specal_slect_classs .el-select-dropdown__item.hover,
.el-select-dropdown__item:hover {
  background-color: #e1e1e1;
}

.specal_slect_classs {
  max-width: 421px !important;
}

.colorReds input {
  border: 1px solid red;
}

.colorRedss .el-input__inner {
  border: 1px solid red;
}

.noneBorder.is-disabled .el-input__inner {
  border: none;
}

.diaBox /deep/ .el-table--enable-row-hover .el-table__body tr:hover > td {
  background: #fff !important;
}

.speical_tag .el-tag__close {
  position: absolute !important;
  top: 6px !important;
  right: 0 !important;
}

.el-select-dropdown {
  max-width: 500px !important;
}

.smallwidth {
  width: 206px !important;
}

.bigwidth {
  width: 595px !important;
  right: 3% !important;
  left: unset !important;
}

.widthwarp {
  display: inline-block;
  max-width: 250px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  height: 14px;
}

.advanced_searchFormTem {
  position: relative;
  top: 90px;
  left: -4.5%;
}

.advanced_searchFormTem1 {
  position: relative;
  top: 90px;
  left: -14.5%;
}

.diaBoxFormMain .advanced_searchFormTem .el-dialog__header {
  padding: 0 !important;
}

.diaBoxFormMain .advanced_searchFormTem1 .el-dialog__header {
  padding: 0 !important;
}

.diaBoxFormMain .advanced_searchFormTem .el-dialog__body {
  padding: 10px !important;
}

.diaBoxFormMain .advanced_searchFormTem1 .el-dialog__body {
  padding: 10px !important;
}

.tagTypeClass {
  font-size: 11px;
  color: #666;
  vertical-align: middle;
}

.gjSearchBody .gjSearchBodyItem {
  padding: 5px !important;
  background: #f8f9fa !important;
}

.diaBoxFormMain .advanced_searchFormTem /deep/ .el-dialog__body {
  padding: 10px 30px 20px !important;
  background: #ebeff3 !important;
}

.diaBoxFormMain .advanced_searchFormTem1 /deep/ .el-dialog__body {
  padding: 10px 30px 20px !important;
  background: #ebeff3 !important;
}
</style>
<style scoped>
.file {
  position: relative;
  display: inline-block;
  background: #fff;
  border: 1px dashed #ccc;
  border-radius: 4px;
  padding: 10px 80px;
  overflow: hidden;
  color: #606266;
  text-decoration: none;
  text-indent: 0;
  line-height: 20px;
  cursor: pointer;
}

.file input {
  position: absolute;
  font-size: 100px;
  right: 0;
  top: 0;
  opacity: 0;
  cursor: pointer;
}
</style>