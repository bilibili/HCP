<template>
  <div class='app-container app_Main'>
    <div style="min-width:840px;" class="detailContent">
      <div style="display:flex;">
        <div style="flex-grow:1;">
          <div style="height:calc(100vh - 46px);overflow: auto;">
            <el-card :class="tableHight == 'up' ? 'tableHight111' : 'cardBox'">
              <div
                style="border-bottom:1px solid #e6e6e6;padding-bottom: 10px;margin-bottom:10px;margin-top: 8px;font-size: 14px;">
                <span style="font-weight:600;">
                  {{ $route.meta.topTitle }}
                  <i class="el-icon-arrow-right" style="font-size: 14px;font-weight:600;margin: 0 6px;"
                    v-if="$route.meta.topTitle"></i>
                  {{ $route.meta.pageTitle }}
                </span>
              </div>
              <div>
                <div style="text-align:right;margin-bottom:7px;" v-show="tableHeaderArr.length">
                  <div style="margin-right:35px;float: left;">

                    <div v-if="$route.name == 'cloudServer'" style="float:left;">
                      <el-button type="primary"
                      style="cursor: pointer;margin-right:15px;"
                      @click="rebootBtnNew('FormCreateCloudProduct', '云服务新建')">新 建</el-button>
                    </div>
                    
                    <div v-if="$route.name == 'projectManage'" style="cursor: pointer;float:left;margin-right:15px;">
                      <el-button v-if="users.indexOf(username) != -1" type="primary"
                        @click="rebootBtnNew('FormAdminCreateCloudProjectConfig','项目创建')">项目创建</el-button>
                      <el-button v-else type="primary"
                        @click="rebootBtnNew('FormCreateCloudProjectConfig','项目新增')">项目新增</el-button>
                      <el-button type="primary"
                        @click="rebootBtnNew('FormJoinCloudProjectConfig','加入项目')">加入项目</el-button>
                      <el-button type="primary"
                        @click="rebootBtnNew('FormLeaveCloudProjectConfig','退出项目')">退出项目</el-button>
                    </div>

                    <div v-if="$route.name == 'accountManage'" style="cursor: pointer;float:left;margin-right:15px;">
                     <el-button style="cursor: pointer;float:left;margin-right:15px;"
                      v-if="users.indexOf(username) != -1" type="primary"
                      @click="rebootBtnNew('FormAdminAppendCloudProjectAccountConfig','项目账号创建')">项目账号创建</el-button>
                    </div>

                    <div v-if="$route.name == 'Useraccount'" style="cursor: pointer;float:left;margin-right:15px;">
                     <el-button type="primary"
                      style="cursor: pointer;float:left;margin-right:15px;"
                      @click="rebootBtnNew('FormAdminCreateCloudAccount','主账号新建')">新 建</el-button>
                    </div>

                    <div v-if="$route.name == 'User'" style="cursor: pointer;float:left;margin-right:15px;">
                     <el-button type="primary"
                      style="cursor: pointer;float:left;margin-right:15px;"
                      @click="rebootBtnNew('FormAdminCreateUser','用户新建')">新 建</el-button>
                    </div>

                    <el-popover placement="bottom" popper-class="popperDiaBox" trigger="click">
                      <div>
                        <div v-if="$route.name == 'cloudServer'">
                          <!-- <div>
                            <el-link style="cursor: pointer;font-size:12px;border:none;margin-bottom:7px;"
                              @click="batchReasseP('FormUpdateCloudServerProject', '', '云服务器批量调整项目', true)">批量调整项目</el-link>
                          </div> -->
                          <div>
                            <el-link style="cursor: pointer;font-size:12px;border:none;margin-bottom:7px;"
                              @click="batchReassembly('FormRebootCloudServer', '', '云服务器重启')">批量重启</el-link>
                          </div>
                          <div>
                            <el-link style="cursor: pointer;font-size:12px;border:none;margin-bottom:7px;"
                              @click="batchReassembly('FormReinstallCloudServer', '', '云服务器重装')">批量重装</el-link>
                          </div>
                          <div>
                            <el-link style="cursor: pointer;font-size:12px;border:none;margin-bottom:7px;"
                              @click="batchReassembly('FormConfigCloudProduct', 'FormConfigCloudServer', '云服务器改配')">批量改配</el-link>
                          </div>
                          <div>
                            <el-link style="cursor: pointer;font-size:12px;border:none;margin-bottom:7px;"
                              @click="batchReassembly('FormDeleteCloudProduct', 'FormDeleteCloudServer', '云服务器清退')">批量清退</el-link>
                          </div>
                          <div>
                            <el-link style="cursor: pointer;font-size:12px;border:none;"
                              @click="batchReassembly('FormRenameBatchCloudServer', 'FormRenameCloudServer', '云服务器改名')">批量改名</el-link>
                          </div>
                        </div>
                      </div>
                      <el-button v-if="$route.name == 'cloudServer'" class="batchBtn"
                        style="cursor: pointer;float:left;" slot="reference">批量操作 <i style="margin-left:5px;"
                          class="el-icon-caret-bottom"></i></el-button>
                    </el-popover>
                  </div>
                  <el-dialog :modal='false' :show-close='false' :close-on-press-escape='true'
                    custom-class='advanced_search' :visible.sync="poHide" style="float:left;margin-left:15px">
                    <div style="background:#f8f9fa;padding:10px;border-radius:5px">
                      <span v-for="(Onetem, leII) in advSearch" :key="leII">
                        <advancedSearchOne :searchMoreArr='searchMoreArr' :tableHeaderArr='tableHeaderArr'
                          :searchLevel='1' :searchLevelData='Onetem' :key="'aa' + leII" :ownIndex='leII'
                          @addNewsearch="addNewsearch" @delNewsearch='delNewsearch' :oneIndex='leII'
                          :advSearch='advSearch' :parentIndex='null' style="margin-top:10px"
                          v-if="!Onetem.child || Onetem.child.length == 0"></advancedSearchOne>
                        <template v-if="Onetem.child && Onetem.child.length != 0">
                          <div v-if='advSearch[0].category'
                            style="border:1px solid #e1e1e1;position:relative;margin:10px 0;font-size:10px">
                            <span style="position:absolute;top:-5px;left:-1px;background:#f8f9fa;font-weight: 600;">{{
                              advSearch[0].category }}</span>
                          </div>
                          <el-card style="margin-top:15px" class="gjSearchBody">
                            <div class='gjSearchBodyItem'>
                              <div v-for="(levItem, levI) in Onetem.child" :key="levI">
                                <advancedSearchOne :searchMoreArr='searchMoreArr' :tableHeaderArr='tableHeaderArr'
                                  :searchLevel='2' :searchLevelData='levItem' :ownIndex='levI' :oneIndex='leII'
                                  @addNewsearch="addNewsearch" @delNewsearch='delNewsearch' :parentIndex='leII'
                                  :advSearch='advSearch' v-if="!levItem.child || levItem.child.length == 0"
                                  style="margin-top:10px"></advancedSearchOne>
                                <div v-if="levItem.child && levItem.child.length != 0">
                                  <div v-if="advSearch[leII].child[0].category"
                                    style="border:1px solid #e1e1e1;position:relative;margin:10px 0;font-size:10px">
                                    <span
                                      style="position:absolute;top:-5px;left:-1px;background:#f8f9fa;font-weight: 600">{{
                                        advSearch[leII].child[0].category }}</span>
                                  </div>
                                  <el-card style="margin-top:15px" class="gjSearchBody">
                                    <div class='gjSearchBodyItem'>
                                      <div v-for="(levItem2, levI2) in levItem.child" :key="levI2">
                                        <advancedSearchOne :searchMoreArr='searchMoreArr'
                                          :tableHeaderArr='tableHeaderArr' levI levI2 :searchLevel='3'
                                          :searchLevelData='levItem2' :ownIndex='levI2' :parentIndex='levI'
                                          @addNewsearch="addNewsearch" :oneIndex='leII' @delNewsearch='delNewsearch'
                                          :advSearch='advSearch' style="margin-top:10px">
                                        </advancedSearchOne>
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
                        <el-button type="primary" plain @click="searchTable(false)">搜 索</el-button>
                      </div>
                    </div>
                  </el-dialog>
                  <div style="border:1px solid #e1e1e1;float:left;height:28px;border-radius:4px">
                    <el-button
                      style="float:left;padding: 6px 10px !important;background: #f5f7fa;color: #409EFF;    border: none;position: relative;z-index: 1;border-radius: 4px 0 0 4px;height:26px;border-right:1px solid #e1e1e1;"
                      icon="el-icon-plus" @click="searchShowStatus">高级搜索</el-button>
                    <el-input style="width:560px;float: left;height:26px" class="batchSearchClass" placeholder="支持模糊搜索"
                      clearable v-model="batchSearchName" @clear="batchSearchNameTable"
                      @keydown.enter.native="batchSearchNameTable">
                    </el-input>
                  </div>
                  <template v-if="this.showlistObj.oneSearch">
                    <el-select @change="batchSearchNameTable"
                      :class="searchOneObj[itemO.name].length > 1 ? 'hideoneacSel' : 'oneacSel'"
                      v-for="(itemO, indexO) in this.showlistObj.oneSearch" collapse-tags
                      :filter-method="(query) => getselectList1(query, itemO)" multiple reserve-keyword
                      :placeholder="'请选择' + itemO.title" filterable clearable
                      style="width:140px;margin-left:15px;float: left;" v-model="searchOneObj[itemO.name]"
                      :key="indexO">
                      <el-option v-for="(oitem, oindex) in itemO.DataList"
                        :value="oitem[itemO.name.split('.')[itemO.name.split('.').length - 1]]"
                        :label="oitem[itemO.name.split('.')[itemO.name.split('.').length - 1]]"
                        :key="oindex"></el-option>
                    </el-select>
                  </template>
                  <i class="el-icon-refresh" style="margin-left:6px;cursor: pointer;color:rgb(64, 158, 255)"
                    @click="searchTable('res')"></i>
                  <template>
                    <el-popover placement="right-start" width="300" trigger="click">
                      <div>
                        <el-checkbox-group v-model="checkboxTH" @change="showTH" size="small">
                          <el-row :gutter="20">
                            <template v-for="(item, index) in tableHeaderArr">
                              <el-col :span="12" v-if="item.title != 'id' && !item.hideIn" :key="index">
                                <el-checkbox v-if="item.name == 'name'" disabled style="width:100%;margin-bottom:7px;"
                                  :label="item.title" border>
                                  {{ item.title + (item.unit ? '(' + item.unit + ')' : '') }}
                                </el-checkbox>
                                <el-checkbox v-else style="width:100%;margin-bottom:7px;" :label="item.title" border>
                                  {{ item.title + (item.unit ? '(' + item.unit + ')' : '') }}
                                </el-checkbox>
                              </el-col>
                            </template>
                          </el-row>
                        </el-checkbox-group>
                        <el-button icon="el-icon-refresh-right" size="mini"
                          style="float: right;cursor: pointer;margin-top:5px" @click="showTHRefresh">重置</el-button>
                      </div>
                      <i style="margin-left: 7px;color: #409EFF;cursor: pointer;font-size: 16px;" slot="reference"
                        class="el-icon-setting"></i>
                    </el-popover>
                    <i style="margin-left: 7px;color: #409EFF;cursor: pointer;font-size: 16px;" class="el-icon-download"
                      v-if="$route.name == '下载按钮'" @click="openDownload"></i>
                  </template>
                </div>
              </div>
              <div v-if="searchObjTag" style="margin:20px 0 10px;cursor: pointer;">
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
                        <span class="tagTypeClass"><span style='margin-right:3px'>{{ tagsValitem[0].category
                        }}</span>(</span>
                        <span :key="index1 + 1" v-for="(itemTwo, index1) in item.child">
                          <span v-if="!itemTwo.child || !itemTwo.child.length">
                            <span v-if="item.child[0].category && index1 != 0" class="tagTypeClass">{{
                              item.child[0].category }}</span>
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
                                {{ itemTwo.name }}&lt;{{ new Date(itemTwo.value).getTime() | filterTimeShow('YYYY-DD-MM')
                                }}
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
                                style="margin:0 3px;cursor: pointer;"
                                v-if="!itemThere.child || !itemThere.child.length">
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
                                  {{ itemThere.name }}>{{ new Date(itemThere.value).getTime() |
                                    filterTimeShow('YYYY-DD-MM') }}
                                </span>
                                <span v-if="itemThere.type == '_HLT'" class="widthwarp">
                                  {{ itemThere.name }}&lt;{{ new Date(itemThere.value).getTime() | filterTimeShow }}
                                </span>
                                <span v-if="itemThere.type == '_DLT'" class="widthwarp">
                                  {{ itemThere.name }}&lt;{{ new Date(itemThere.value).getTime() |
                                    filterTimeShow('YYYY-DD-MM') }}
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
                    type="primary" @click="refAdvancedSearchLast()">重置</el-link>
                </div>
              </div>
              <el-table :data="tableData"
                :style="isCollapseStatus ? 'width: calc(100vw - 80px);' : 'width: calc(100vw - 170px);'"
                @selection-change="handleSelectionChangeTol" v-loading="tableLoading"
                :row-class-name='tableRowClassName'>
                <el-table-column v-if="tableHeaderArr.length" type="selection" width="40"></el-table-column>
                <template v-for="(item, index) in tableHeaderArr">
                  <template v-if="(item.title != 'id' && !item.hideIn)">
                    <template v-if="checkboxTH.indexOf(item.title) != -1">
                      <template v-if="item.enum">
                        <el-table-column :key="index" :prop="item.name"
                          :label="item.title + (item.unit ? '(' + item.unit + ')' : '')" :show-overflow-tooltip="true">
                          <template slot-scope="scope">
                            <div>
                              <template v-if="item.customStyle">
                                <span v-for="(eeitem, eeindex) in item.customStyle" :key="eeindex">
                                  <span style="font-size:12px;"
                                    v-if="eeitem.name && scope.row[eeitem.name.replace(/\./g, '')]">
                                    <!-- 有link链接 -->
                                    <span v-if="eeitem.link || eeitem.urlss">
                                      <span style="font-size:12px;color:#409EFF;cursor:pointer" type="primary"
                                        @click="linkOpen(scope.row[eeitem.linkId.replace(/\./g, '')], eeitem.link)"
                                        :underline="false" v-if="eeitem.link">
                                        <span v-if="eeitem.prefix">
                                          {{ eeitem.prefix }}
                                        </span>
                                        <span v-if="eeitem.name == item.name">
                                          <span v-for="(objitem, objindex) in  item.enum" :key="objindex">
                                            <span v-if="objitem.en == scope.row[item.name.replace(/\./g, '')]">
                                              {{ objitem.zh }}
                                            </span>
                                          </span>
                                        </span>
                                        <span v-else>
                                          <span v-if="scope.row['chargeType'] != 'postpaid'">
                                            {{ scope.row[eeitem.name.replace(/\./g, '')] }}
                                          </span>
                                        </span>
                                      </span>
                                      <span v-if="eeitem.urlss">
                                        <span v-if="eeitem.name == item.name">
                                          <span v-for="(objitem, objindex) in  item.enum" :key="objindex">
                                            <span v-if="objitem.en == scope.row[item.name.replace(/\./g, '')]">
                                              <img v-if="objitem.zh == '百度云'"
                                                style="vertical-align: middle; position: relative; top: -1px;"
                                                src="@/assets/icon/cloud/BDC.png" alt="" width="20px" height="20px">
                                              <img v-if="objitem.zh == '腾讯云'"
                                                style="vertical-align: middle; position: relative; top: -1px;"
                                                src="@/assets/icon/cloud/TCC.svg" alt="" width="20px" height="20px">
                                              <img v-if="objitem.zh == '华为云'"
                                                style="vertical-align: middle; position: relative; top: -1px;"
                                                src="@/assets/icon/cloud/HWC.svg" alt="" width="20px" height="20px">
                                              <img v-if="objitem.zh == '阿里云'"
                                                style="vertical-align: middle; position: relative; top: -1px;"
                                                src="@/assets/icon/cloud/ABC.svg" alt="" width="20px" height="20px">
                                              <img v-if="objitem.zh == '亚马逊' || objitem.zh == 'AWS'"
                                                style="vertical-align: middle; position: relative; top: -1px;"
                                                src="@/assets/icon/cloud/AWS.png" alt="" width="18px">
                                              {{ objitem.zh }}
                                            </span>
                                          </span>
                                        </span>
                                      </span>
                                    </span>
                                    <span v-else>
                                      <span v-if="eeitem.prefix">
                                        {{ eeitem.prefix }}
                                      </span>
                                      <span v-if="eeitem.description" style="color:#ccc">
                                        {{ eeitem.description }}
                                      </span>
                                      <span v-if="eeitem.color" :style="{ 'color': eeitem.color }">
                                        <span v-if="eeitem.name == item.name">
                                          <span v-for="(objitem, objindex) in  item.enum" :key="objindex">
                                            <span v-if="objitem.en == scope.row[item.name.replace(/\./g, '')]">
                                              {{ objitem.zh }}
                                            </span>
                                          </span>
                                        </span>
                                        <span v-else>
                                          {{ scope.row[eeitem.name.replace(/\./g, '')] }}
                                        </span>
                                      </span>
                                      <span v-else>
                                        <span v-if="eeitem.name == item.name">
                                          <span v-for="(objitem, objindex) in  item.enum" :key="objindex">
                                            <span v-if="objitem.en == scope.row[item.name.replace(/\./g, '')]">
                                              <span v-if="item.name == 'chargeType'">
                                                <span v-if="scope.row['chargeType'] == 'postpaid'">
                                                  {{ objitem.zh }}
                                                </span>
                                                <span v-else>
                                                  <span v-if="scope.row['renewStatus']">
                                                    <span v-if="scope.row['renewStatus'] == 'auto'"
                                                      style="color:#67C23A">
                                                      {{ objitem.zh }}
                                                    </span>
                                                    <span v-if="scope.row['renewStatus'] == 'manual'"
                                                      style="color:rgb(249, 225, 3)">
                                                      {{ objitem.zh }}
                                                    </span>
                                                    <span v-if="scope.row['renewStatus'] == 'no'" style="color:red">
                                                      {{ objitem.zh }}
                                                    </span>
                                                  </span>
                                                  <span v-else style="color:red">
                                                    {{ objitem.zh }}
                                                  </span>
                                                </span>
                                              </span>
                                              <span v-else>
                                                {{ objitem.zh }}
                                              </span>
                                            </span>
                                          </span>
                                        </span>
                                        <span v-else>
                                          <span v-if="scope.row['chargeType'] != 'postpaid'">
                                            <span v-if="scope.row['renewStatus'] && scope.row['renewStatus'] == 'auto'">
                                              {{ scope.row[eeitem.name.replace(/\./g, '')] | filterTimeShow }}
                                            </span>
                                            <span v-else>
                                              <span
                                                v-if="parseInt((new Date(scope.row['expireAt']) - new Date()) / (1000 * 60 * 60 * 24)) < 7"
                                                style="color:red">
                                                {{ scope.row[eeitem.name.replace(/\./g, '')] | filterTimeShow }}
                                              </span>
                                              <span
                                                v-if="parseInt((new Date(scope.row['expireAt']) - new Date()) / (1000 * 60 * 60 * 24)) > 7 && parseInt((new Date(scope.row['expireAt']) - new Date()) / (1000 * 60 * 60 * 24)) < 30"
                                                style="color:rgb(249, 225, 3)">
                                                {{ scope.row[eeitem.name.replace(/\./g, '')] | filterTimeShow }}
                                              </span>
                                              <span
                                                v-if="parseInt((new Date(scope.row['expireAt']) - new Date()) / (1000 * 60 * 60 * 24)) > 30">
                                                {{ scope.row[eeitem.name.replace(/\./g, '')] | filterTimeShow }}
                                              </span>
                                            </span>
                                          </span>
                                        </span>
                                      </span>
                                    </span>
                                  </span>
                                  <span v-if="eeitem.subName">
                                    ({{ scope.row[eeitem.subName.replace(/\./g, '')] }})
                                  </span><br
                                    v-if="eeitem.name && scope.row[eeitem.name.replace(/\./g, '')] && scope.row['chargeType'] != 'postpaid'">
                                </span>
                              </template>
                              <template v-else v-for="(objitem, objindex) in  item.enum">
                                <span :key="objindex">
                                  <template v-if="objitem.style">
                                    <template v-if="objitem.style == 'info'">
                                      <el-link :key="index" v-if="objitem.en == scope.row[item.name]" :underline="false"
                                        class="el-icon-remove" style="font-size:12px;" :type="objitem.style"><span
                                          style="margin-left:5px;">{{ objitem.zh }}</span></el-link>
                                    </template>
                                    <template v-else>
                                      <el-link :key="index" v-if="objitem.en == scope.row[item.name]" :underline="false"
                                        :class="'el-icon-' + (objitem.style == 'danger' ? 'error' : objitem.style)"
                                        style="font-size:12px;" :type="objitem.style"><span style="margin-left:5px;">{{
                                          objitem.zh }}</span></el-link>
                                    </template>
                                  </template>
                                  <template v-else>
                                    <el-link :key="index" v-if="objitem.en == scope.row[item.name]" :underline="false"
                                      style="font-size:12px;" :type="objitem.style"><span style="margin-left:5px;">{{
                                        objitem.zh }}</span></el-link>
                                  </template>
                                </span>
                              </template>
                              <template>
                                <div v-if="item.enumArr">
                                  <span v-if="item.enumArr.indexOf(scope.row[item.name]) == -1">
                                    {{ scope.row[item.name] }}
                                  </span>
                                </div>
                              </template>
                            </div>
                          </template>
                        </el-table-column>
                      </template>
                      <template v-else>
                        <el-table-column
                          v-if="(item.type == 'object')" :key="index"
                          :label="item.title + (item.unit ? '(' + item.unit + ')' : '')" :show-overflow-tooltip="true">
                          <template slot-scope="scope">
                            <div v-if="scope.row[item.name.split('.')[0]]">
                              <span v-if="item.name.split('.').length == 2">
                                <template>
                                  <template v-if="item.list">
                                    <template v-if="item.more && scope.row[item.name.split('.')[0]].length > 1">
                                      <span v-for="ccindex in 1" :key="ccindex">
                                        <div>{{ scope.row[item.name.split('.')[0]][ccindex - 1][item.name.split('.')[1]]
                                        }}
                                        </div>
                                      </span>
                                      <el-link
                                        style="cursor: pointer;color:#1890ff;font-size:12px;border:none;margin:0 5px;position: relative;left: -4px;"
                                        @click="reviewDetailTo(scope.row, item.more)">更多({{
                                          scope.row[item.name.split('.')[0]].length }})</el-link>
                                    </template>
                                    <template v-else>
                                      <span v-for="(ccitem, ccindex) in scope.row[item.name.split('.')[0]]"
                                        :key="ccindex">
                                        {{ ccindex == 0 ? ccitem[item.name.split('.')[1]] : ',' +
                                          ccitem[item.name.split('.')[1]] }}
                                      </span>
                                    </template>

                                  </template>
                                  <template v-else>
                                    <span>
                                      {{ scope.row[item.name.replace(/\./g, "")] }}
                                    </span>
                                  </template>
                                </template>
                              </span>
                              <span v-else>
                                <span v-if="$route.name == 'serverImage'">
                                  <span>{{ scope.row[item.name.split('.')[0]][0].provider.name }}</span>
                                </span>
                                <span v-if="$route.name == 'serverType'">
                                  <span>{{ scope.row[item.name.split('.')[0]][0].provider.name }}</span>
                                </span>
                                <span v-if="$route.name != 'serverType' && $route.name != 'serverImage'">
                                  <span v-if="item.more">
                                    <span
                                      v-if="item.more && scope.row[item.name.replace(/\./g, '')].split(',').length > 1">
                                      <div v-for="iit in 1" :key="iit">
                                        {{ scope.row[item.name.replace(/\./g, '')].split(',')[iit] }}
                                      </div>
                                      <el-link
                                        style="cursor: pointer;color:#1890ff;font-size:12px;border:none;margin:0 5px;position: relative;left: -4px;"
                                        @click="reviewDetailTo(scope.row, item.more)">更多({{
                                          scope.row[item.name.replace(/\./g, '')].split(',').length }})</el-link>
                                    </span>
                                    <span v-else>
                                      <div v-for="iit in 1" :key="iit">
                                        {{ scope.row[item.name.replace(/\./g, '')].split(',')[iit] }}
                                      </div>
                                    </span>
                                  </span>
                                  <span v-else>
                                    {{ scope.row[item.name.replace(/\./g, '')] }}
                                  </span>
                                </span>
                              </span>
                            </div>
                          </template>
                        </el-table-column>
                        <template v-else>
                          <template v-if="item.customStyle">
                            <el-table-column v-if="item.style == 'lineFeedCopy'" :prop="item.name"
                              :label="item.title + (item.unit ? '(' + item.unit + ')' : '')" :key="index" width="250"
                              fixed :show-overflow-tooltip="true">
                              <template slot-scope="scope">
                                <span style="font-size:12px;color:#409EFF;cursor:pointer" type="primary"
                                  @click="reviewDetail(scope.row, $route.name, scope)" :underline="false">
                                  <span class="hoverSpan">
                                    {{ scope.row.name == '' || !scope.row.name ? '无名称' : scope.row.name }}
                                  </span>
                                </span>
                                <br>
                                <span style="color:#ccc;" v-for="(eeitem, eeindex) in item.customStyle" :key="eeindex">
                                  {{ scope.row[eeitem.name.replace(/\./g, '')] }}
                                  <br>
                                </span>
                                <i style="cursor:pointer;vertical-align: middle;position: absolute; top: 33%;z-index: 2;right: 0px;margin-top: -4px;"
                                  class="el-icon-copy-document" @click.stop="_copy(scope.row[item.name])"></i>
                              </template>
                            </el-table-column>
                            <el-table-column v-if="item.style == 'serverType'" :prop="item.name"
                              :label="item.title + (item.unit ? '(' + item.unit + ')' : '')" :key="index"
                              :show-overflow-tooltip="true">
                              <template slot-scope="scope">
                                <span>
                                  <span v-if="scope.row[item.customStyle[0].cpuName.replace(/\./g, '')]">
                                    {{ scope.row[item.customStyle[0].cpuName.replace(/\./g, '')] }}
                                    <span>C</span>
                                  </span>
                                  <span v-if="scope.row[item.customStyle[0].memoryName.replace(/\./g, '')]">
                                    {{ scope.row[item.customStyle[0].memoryName.replace(/\./g, '')] }}
                                    <span v-if="item.customStyle[0].memoryName.indexOf('memoryGB') != -1">GB</span>
                                    <span v-else>MB</span>
                                  </span>
                                  <span v-if="scope.row[item.customStyle[0].gpuName.replace(/\./g, '')]">
                                    GPU:{{ scope.row[item.customStyle[0].gpuName.replace(/\./g, '')] }}个
                                  </span>
                                  <span v-if="scope.row[item.customStyle[0].gpuModel.replace(/\./g, '')]">
                                    {{ scope.row[item.customStyle[0].gpuModel.replace(/\./g, '')] }}
                                  </span>
                                  <span v-if="scope.row[item.customStyle[0].diskSizeName.replace(/\./g, '')]">
                                    磁盘{{ scope.row[item.customStyle[0].diskSizeName.replace(/\./g, '')] }}G
                                  </span>
                                </span><br>
                                <span>
                                  {{ scope.row[item.customStyle[1].name.replace(/\./g, '')] }}
                                </span>
                              </template>
                            </el-table-column>
                            <el-table-column v-if="!item.style" :prop="item.name"
                              :label="item.title + (item.unit ? '(' + item.unit + ')' : '')" :key="index"
                              :show-overflow-tooltip="true">
                              <template slot-scope="scope">
                                <span v-for="(eeitem, eeindex) in item.customStyle" :key="eeindex">
                                  <span style="font-size:12px;"
                                    v-if="eeitem.name && scope.row[eeitem.name.replace(/\./g, '')]">
                                    <!-- 有link链接 -->
                                    <span v-if="eeitem.link">
                                      <span style="font-size:12px;color:#409EFF;cursor:pointer" type="primary"
                                        @click="linkOpen(scope.row[eeitem.linkId.replace(/\./g, '')], eeitem.link)"
                                        :underline="false">
                                        <span v-if="eeitem.prefix">
                                          {{ eeitem.prefix }}
                                        </span>
                                        {{ scope.row[eeitem.name.replace(/\./g, '')] }}</span>
                                    </span>
                                    <span v-else>
                                      <span v-if="eeitem.urls">
                                        <img v-if="scope.row[eeitem.name.replace(/\./g, '')] == '百度云'"
                                          style="vertical-align: middle; position: relative; top: -1px;"
                                          src="@/assets/icon/cloud/BDC.png" alt="" width="20px" height="20px">
                                        <img v-if="scope.row[eeitem.name.replace(/\./g, '')] == '腾讯云'"
                                          style="vertical-align: middle; position: relative; top: -1px;"
                                          src="@/assets/icon/cloud/TCC.svg" alt="" width="20px" height="20px">
                                        <img v-if="scope.row[eeitem.name.replace(/\./g, '')] == '华为云'"
                                          style="vertical-align: middle; position: relative; top: -1px;"
                                          src="@/assets/icon/cloud/HWC.svg" alt="" width="20px" height="20px">
                                        <img v-if="scope.row[eeitem.name.replace(/\./g, '')] == '阿里云'"
                                          style="vertical-align: middle; position: relative; top: -1px;"
                                          src="@/assets/icon/cloud/ABC.svg" alt="" width="20px" height="20px">
                                        <img v-if="scope.row[eeitem.name.replace(/\./g, '')] == '亚马逊' || scope.row[eeitem.name.replace(/\./g, '')] == 'AWS'"
                                          style="vertical-align: middle; position: relative; top: -1px;"
                                          src="@/assets/icon/cloud/AWS.png" alt="" width="18px">
                                      </span>
                                      <span v-if="eeitem.prefix">
                                        {{ eeitem.prefix }}
                                      </span>
                                      <span v-if="eeitem.description" style="color:#ccc">
                                        {{ eeitem.description }}
                                      </span>
                                      <span v-if="eeitem.color" :style="{ 'color': eeitem.color }">
                                        {{ scope.row[eeitem.name.replace(/\./g, '')] }}
                                      </span>
                                      <span v-else>
                                        {{ scope.row[eeitem.name.replace(/\./g, '')] }}
                                      </span>
                                      <span v-if="eeitem.unit">
                                        {{ eeitem.unit }}
                                      </span>
                                    </span>
                                  </span>
                                  <!-- 带括号得 -->
                                  <span v-if="eeitem.subName">
                                    ({{ scope.row[eeitem.subName.replace(/\./g, '')] }})
                                  </span>
                                  <br v-if="eeitem.name && scope.row[eeitem.name.replace(/\./g, '')]">
                                </span>
                              </template>
                            </el-table-column>
                          </template>
                          <template v-if="!item.customStyle">
                            <template v-if="item.tagsArr">
                              <el-table-column :key="index"
                                :label="item.title + (item.unit ? '(' + item.unit + ')' : '')" width="100"
                                align="center">
                                <template slot-scope="scope">
                                  <template v-if="scope.row[item.name] && scope.row[item.name].length != 0">
                                    <el-tooltip class="item" effect="dark" placement="top-start">
                                      <div slot="content" v-for="(tagItem, tagInx) in scope.row[item.name]"
                                        :key="tagInx">
                                        {{ tagItem }}
                                      </div>
                                      <span style="font-size:12px;cursor: pointer;" v-if="item.name == 'tags'">
                                        <img
                                          style="vertical-align: middle; position: relative; top: -1px;margin-right:3px"
                                          src="@/assets/icon/label.svg" alt="" width="20px" height="20px">{{
                                            scope.row[item.name].length }}
                                      </span>
                                      <span style="font-size:12px;cursor: pointer;" v-if="item.name == 'labels'">
                                        <img
                                          style="vertical-align: middle; position: relative; top: -1px;margin-right:3px"
                                          src="@/assets/icon/cloudlabel.png" alt="" width="20px" height="20px">{{
                                            scope.row[item.name].length }}
                                      </span>
                                    </el-tooltip>
                                  </template>
                                </template>
                              </el-table-column>
                            </template>
                            <template v-else>
                              <el-table-column :key="index" :prop="item.name"
                                :label="item.title + (item.unit ? '(' + item.unit + ')' : '')"
                                :show-overflow-tooltip="true">
                                <template slot-scope="scope">
                                  <span v-if="$route.name == 'serverType' && item.name == 'accounts.provider.name'">{{
                                    scope.row[item.name.split('.')[0]][0].provider.name }}</span>
                                  <span v-else>
                                    <span v-if="item.name.indexOf('.') != -1">
                                      {{ scope.row[item.name.replace(/\./g, "")] }}
                                    </span>
                                    <span v-else>{{ scope.row[item.name] }}</span>
                                  </span>
                                </template>
                              </el-table-column>
                            </template>
                          </template>
                        </template>
                      </template>
                    </template>
                  </template>
                </template>
                <el-table-column v-if="tableHeaderArr.length" fixed="right" label="操作" align="center" width="140">
                  <template slot-scope="scope">
                    <el-link
                      v-if="$route.name != 'resourceTemplate' && $route.name != 'elastic' && $route.name != 'inspectionRules'"
                      style="cursor: pointer;color:#1890ff;font-size:12px;border:none;margin:0 5px;"
                      @click="reviewDetail(scope.row, $route.name, scope)">详情</el-link>
                    <span v-if="$route.name == 'cloudServer'">
                      <el-link style="cursor: pointer;color:#1890ff;font-size:12px;border:none;margin:0 5px;"
                        @click="rebootBtn('FormRebootCloudServer', '', [scope.row], '云服务器重启')">重启</el-link>
                      <el-popover popper-class="popperDiaBox1" placement="bottom" trigger="click">
                        <div>
                          <div>
                            <el-link style="cursor: pointer;font-size:12px;border:none;margin:0 15px 7px 0px;"
                              @click="rebootBtn('FormReinstallCloudServer', '', [scope.row], '云服务器重装')">重装</el-link>
                          </div>
                          <div>
                            <el-link style="cursor: pointer;font-size:12px;border:none;margin:0 15px 7px 0px;"
                              @click="rebootBtn('FormConfigCloudProduct', 'FormConfigCloudServer', [scope.row], '云服务器改配')">改配</el-link>
                          </div>
                          <div>
                            <el-link style="cursor: pointer;font-size:12px;border:none;margin:0 15px 7px 0px;"
                              @click="rebootBtn('FormDeleteCloudProduct', 'FormDeleteCloudServer', [scope.row], '云服务器清退')">清退</el-link>
                          </div>
                          <div>
                            <el-link style="cursor: pointer;font-size:12px;border:none;margin:0 15px 7px 0px;"
                              @click="rebootBtn('FormRenameBatchCloudServer', 'FormRenameCloudServer', [scope.row], '云服务器改名')">改名</el-link>
                          </div>
                        </div>
                        <el-link type="primary" style="margin-left:5px;cursor: pointer;font-size:12px;margin:0 5px;"
                          slot="reference">更多<i class="el-icon-caret-bottom"></i></el-link>
                      </el-popover>
                    </span>

                    <span v-if="$route.name=='projectManage'">
                      <el-link type="primary"  style="margin:0px 5px;font-size:12px;" :underline="false"
                        @click="rebootBtnNoProject('FormUpdateCloudProjectConfig','',[scope.row],'项目编辑')">编辑
                      </el-link>
                      <el-link type="primary"  style="margin:0px 5px;font-size:12px;" :underline="false"
                        @click="rebootBtnNoProject('FormDeleteCloudProjectConfig','',[scope.row],'项目删除')">删除
                      </el-link>
                    </span>
                   
                    <span v-if="$route.name=='accountManage'">
                      <el-link type="primary"  style="margin:0px 5px;font-size:12px;" :underline="false"
                        @click="rebootBtnNoProject('FormUpdateCloudProjectAccountConfig','',[scope.row],'账号配置编辑')">编辑
                      </el-link>
                    </span>

                    <span v-if="$route.name ==  'CloudProvider'">
                      <el-link style="cursor: pointer;color:#1890ff;font-size:12px;border:none;margin:0 5px;"
                        @click="rebootBtnNoProject('FormUpdateCloudProvider','',[scope.row],'云厂商编辑')">编辑</el-link>
                    </span>

                    <span v-if="$route.name == 'Useraccount'">
                      <el-link style="cursor: pointer;color:#1890ff;font-size:12px;border:none;margin:0 5px;"
                        @click="rebootBtnNoProject('FormAdminUpdateCloudAccount','',[scope.row],'主账号编辑',)">编辑</el-link>
                      <el-link style="cursor: pointer;color:#1890ff;font-size:12px;border:none;margin:0 5px;"
                        @click="rebootBtnNoProject('FormAdminDeleteCloudAccount','',[scope.row],'主账号删除')">删除</el-link>
                    </span>
                    
                    <span v-if="$route.name == 'User'">
                      <el-link style="cursor: pointer;color:#1890ff;font-size:12px;border:none;margin:0 5px;"
                        @click="rebootBtnNoProject('FormAdminUpdateUser','',[scope.row],'用户编辑',)">编辑</el-link>
                        <el-link style="cursor: pointer;color:#1890ff;font-size:12px;border:none;margin:0 5px;"
                        @click="rebootBtnNoProject('FormAdminDeleteUser','',[scope.row],'用户删除')">删除</el-link>
                    </span>

                  </template>
                </el-table-column>
              </el-table>
              <div style="margin-top:10px;" id="tablePagina" v-show="tableHeaderArr.length">
                <el-pagination @current-change="changehandleCurrent" @size-change="handleSizeChange"
                  :current-page='pageNum' :page-size="pageSize" :page-sizes="[10, 20, 30, 100, 150, 200]"
                  layout="total, sizes, prev, pager, next, jumper" :total="total" v-if="!paginationSatatus">
                </el-pagination>
                <el-pagination @current-change="changehandleCurrent" @size-change="handleSizeChange" :total="total"
                  :page-size="pageSize" :current-page='pageNum' small v-else>
                </el-pagination>
              </div>
            </el-card>
          </div>
        </div>
      </div>
    </div>
    <div v-if="drawDetail">
      <drawer-Box @addCloseDraw="addCloseDraw" @tableHeight='tableHeight' :jobIds="jobIds"
        :tableDetailStatus='tableDetailStatus'></drawer-Box>
    </div>
    <!-- 新建任务弹框 infoAddAttr 表格行信息 infoAddForm表单提交信息 reviewStatus弹框状态-->
    <div v-if="reviewStatus222">
      <formtem-Dialog :nodeCheckName='nodeCheckName' :detailOrLook="detailOrLook" @addClose="addClose"
        :detailOrsigle='detailOrsigle' :reviewStatus='reviewStatus222' :newTaskName="newTaskName"
        :historyTableList='historyTableList' :historyparamreferenceShow='historyparamreferenceShow'
        :baseConfigData='baseConfigData' :childNodename='childNodename' :checkList='checkList'
        :FormTemplatjobMode='FormTemplatjobMode'></formtem-Dialog>
    </div>
    <el-dialog title="表格下载" :visible.sync="dialogDownload" :close-on-click-modal="false" width="500px"
      class="diaBox111">
      <div>
        <div style="padding-left: 10px;margin-bottom: 5px;padding-bottom: 5px;border-bottom: 1px solid #d9d9d9;">
          <el-checkbox :indeterminate="isIndeterminate" v-model="checkAll"
            @change="handleCheckAllChange">全选</el-checkbox>
        </div>
        <div style="max-height:350px;overflow: auto;">
          <div style="padding-right:10px;">
            <el-checkbox-group @change="handleCheckedCitiesChange" v-model="checkDownload" size="small">
              <el-row :gutter="20">
                <template v-for="(item, index) in excelTitleDown">
                  <el-col :span="12" :key="index">
                    <el-checkbox style="width:100%;margin-bottom:7px;" :label="item" border>
                      {{ item.title }}
                    </el-checkbox>
                  </el-col>
                </template>
              </el-row>
            </el-checkbox-group>
          </div>
        </div>
        <div style="padding-top: 5px;text-align: right;">
          <el-button @click="dialogDownload = false">取 消</el-button>
          <el-button type="primary" @click="subDownload" :loading="Downloading">确 定</el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>
<script>
import drawerBox from "@/views/components/detail";
import formtemDialog from "@/views/components/formTemp";
import { userMixin } from "@/components/mixin/user";
import { JSONConfig } from "@/components/mixin/JSONConfig";
import { productJSON } from "@/components/mixin/productJSON";
import { productLIst } from "@/components/mixin/productLIst";
import advancedSearchOne from "@/views/components/advancedSearchOne";
import { advSearchHandle } from "@/components/mixin/advSearchHandle";
import { mapGetters } from "vuex";
//例如：import 《组件名称》 from '《组件路径》';
export default {
  //import引入的组件需要注入到对象中才能使用
  props: [],
  mixins: [userMixin,JSONConfig, productJSON, advSearchHandle, productLIst],
  components: {
    drawerBox,
    formtemDialog,
    advancedSearchOne,
  },
  computed: {
    ...mapGetters(["loginUserName", "formPolicyDomain","loginUserInfo"]),
  },
  data() {
    //这里存放数据
    return {
      tablePagina: "",
      paginationSatatus: null,
      advSearch: [
        {
          name: "",
          type: "",
          val: "",
          category: "",
          index: [],
          child: [],
          objContent: {},
        },
      ],
      lastAdvSearch: [],
      advSearchIndex: null,
      enumAllArr: [],
      cascaderMorestr: "",
      searchMoreArrSecond: [],
      defaultData: "",
      cascaderMorestrSecond: null,
      tableHight: "up",
      tableRowIndex: null,
      tableDetailStatus: "",
      FormTemplatjobMode: "",
      arrEunmName: [],
      batcharrLength: [],
      tableLoading: false,
      baseConfigData: {},
      batchSearchName: "",
      activeName: "first",
      filterText: "",
      treedata: [],
      tableHeaderArr: [],
      tableData: [],
      valueData: { other: "", running: "", stopped: "", total: "" },
      nodeTitle: "",
      total: 0,
      pageNum: 1,
      pageSize: 10,
      searchStr: "",
      searchObj: {},
      searchOneObj: {},
      tableHeaderObj: {},
      thzh: [],
      thaazh: [],
      getTitleNode: "",
      A1Arr: [],
      A2Arr: [],
      infoReviewForm: {},
      infoAddBFForm: {},
      infoReviewAttr: [],
      gaipeiArr: [],
      dialogReview: false,
      nodeCheckName: "",
      nodeCheckProductName: "",
      dialogtitle: "",
      reviewStatus: false,
      reviewStatus222: false,
      detailOrLook: false,
      detailOrsigle: "",
      newTaskName: "",
      historyTableList: {},
      historyparamreferenceShow: [],
      checkboxTH: [],
      checkboxTHBF: [],
      batchHostarr: [],
      isCollapseStatus: false,
      checkProjectData: [],
      projectNoget: false,
      hintTitle: "",
      username: "",
      drawDetail: false,
      jobIds: "",
      childNodename: "",
      checkList: [],
      showlistObj: {},
      searchFormData: {},
      searchObjTag: [],
      classWidth: "",
      searchMoreArr: [],
      cascaderMoreArr: [],
      tagStatus: false,
      tagTitle: "",
      tagAllArr: [],
      poHide: false,

      dialogDownload: false,
      Downloading: false,
      checkDownload: [],
      isIndeterminate: false,
      checkAll: false,
      excelTitleDown: [],
    };
  },
  watch: {
    filterText(val) {
      this.$refs.tree.filter(val);
    },
    $route: "routeChange",
  },
  //生命周期 - 挂载完成（可以访问DOM元素）
  mounted() {
    this.username = sessionStorage.getItem("username");
    this.getTreeData();
    this.$EventBus.$on("isCollapse", (value) => {
      this.isCollapseStatus = value ? true : false;
    });
    this.$EventBus.$on("checkProjectArr", (value) => {
      this.checkProjectData = value;
      // this.getTable(this.treedata.schema);
    });
    this.tableHeaderArr.map((k, i) => {
      if (!k.hide && !k.hideIn) {
        this.excelTitleDown.push({
          name: k.name,
          title: k.title,
        });
      }
      if (k.enum) {
        k.enum.map((euitem, eni) => {
          let aaa = `${k.name}:::${euitem.zh}`;
          this.enumAllArr[aaa] = euitem.en;
        });
      }
    });
    this.tablePagina = document.getElementById("tablePagina");
    this.tablePagina.style.textAlign = "center";
    this.tableHight = "";
    if (this.$route.query.detail_ID) {
      this.jobIds = Number(this.$route.query.detail_ID);
      this.drawDetail = true;
      this.tableDetailStatus = sessionStorage.getItem("tableDetailStatusC") ? sessionStorage.getItem("tableDetailStatusC") : "halfScreen"; //注释
    }
  },
  created() {
    this.username = sessionStorage.getItem("username");
    this.isCollapseStatus = sessionStorage.getItem("isCollapse") == 1 ? true : false;
    this.checkProjectData = sessionStorage.getItem("checkProjectArr") ? sessionStorage.getItem("checkProjectArr").split(",") : [];
  },
  destroyed() {
    this.$EventBus.$off("isCollapse");
    this.$EventBus.$off("checkProjectArr");
  },
};
</script>
<style scoped>
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

.batchSearchClass /deep/ .el-input__inner {
  border: none;
  height: 26px;
}

.langClass /deep/ .el-tag.el-tag--info {
  position: relative;
}

.langClass /deep/ .el-select__tags-text {
  width: 143px;
  overflow: hidden;
  display: inline-block;
  text-overflow: ellipsis;
  white-space: nowrap;
  vertical-align: middle;
  padding: 0 3px;
}

.langClass /deep/ .el-select__tags span.el-tag i {
  position: absolute;
  right: 0;
  top: 5px;
}

.app-container /deep/ .el-tooltip__popper {
  display: none !important;
}

.detailContent /deep/ .el-table::before {
  background: none !important;
}

.activeColor {
  color: #1890ff;
  font-weight: 700;
}

.app-container /deep/ .el-form-item__label {
  color: #999;
}

.app-container /deep/ h6 {
  color: #606266;
  font-weight: 600;
  font-size: 15px;
  margin: 0px 0 15px 0;
  border-left: 3px solid #1890ff;
  text-indent: 10px;
  line-height: 24px;
  height: 24px;
}

.app-container /deep/ .infoTable .el-table__header {
  display: none;
}

.hoverSpan:hover {
  border-bottom: 1px solid #409eff;
}

.app-container /deep/ .demo-ruleForm .specialTable .el-table__header {
  display: block;
}

.app-container /deep/ .infoTable td {
  border-bottom: none;
}

.app-container /deep/ .infoTable td td {
  border-bottom: 1px solid #ebeef5;
}

.app-container /deep/ .infohide .el-form-item__label {
  display: none;
}

.app-container /deep/ .el-textarea .el-textarea__inner,
.app-container /deep/ .el-input .el-input__inner,
.app-container /deep/ .el-select .el-input.is-disabled .el-input__inner {
  background: #fff !important;
}

.app_Main {
  padding: 0px !important;
}

.cardBox {
  min-height: calc(100vh - 60px);
}

.cardBox /deep/ .el-card__body {
  padding: 10px !important;
}

.searchBox /deep/ .el-card__body {
  padding: 20px 30px !important;
}

.searchBox /deep/ .el-form-item__label {
  font-size: 13px;
}

.searchBox /deep/ .el-form-item {
  margin-bottom: 0px !important;
}

.searchBox /deep/ .el-input__inner {
  height: 28px !important;
  line-height: 28px !important;
}

.diaBox111 /deep/ .el-table th {
  padding: 0px !important;
}

.diaBox111 /deep/ .el-table .cell {
  padding: 0px !important;
}

.diaBox111 /deep/ .el-table .el-table .cell {
  padding: 7px 2px !important;
}

.diaBox111 /deep/ .el-table td {
  padding: 0 !important;
}

.diaBox111 /deep/ .specal_slect_classs .el-select-dropdown__item {
  width: 100%;
}

.diaBox111 /deep/ .el-dialog__body {
  padding: 0 20px 20px 20px !important;
}

.diaBox111 /deep/ .el-dialog__body .el-dialog__body {
  padding: 20px 20px 20px 20px !important;
}

.app_Main /deep/ .el-table tr td .cell {
  font-size: 12px !important;
}

.el-date-editor.el-input,
.el-date-editor.el-input__inner {
  width: 200px !important;
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

.detailContent
  /deep/
  .hideoneacSel
  .el-select__tags
  span
  span.el-tag.el-tag--info:first-child {
  max-width: 40px !important;
}

.detailContent
  /deep/
  .oneacSel
  .el-select__tags
  span
  span.el-tag.el-tag--info:first-child {
  max-width: 75px !important;
  overflow: hidden;
  white-space: nowrap !important;
  text-overflow: ellipsis !important;
  padding-right: 3px;
}

.detailContent .hideoneacSel /deep/ .el-select__tags span {
  max-width: 44px !important;
  overflow: hidden !important;
  white-space: nowrap !important;
  text-overflow: ellipsis !important;
  padding: 0 3px;
}

.detailContent .hideoneacSel /deep/ .el-select__tags span.el-tag,
.detailContent .oneacSel /deep/ .el-select__tags span.el-tag {
  position: relative;
}

.detailContent .hideoneacSel /deep/ .el-select__tags span.el-tag i,
.detailContent .oneacSel /deep/ .el-select__tags span.el-tag i {
  position: absolute;
  right: 0;
  top: 5px;
}
</style>
<style>
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

.popper__arrow {
  display: none !important;
}

.el-cascader-panel.is-bordered {
  border: none !important;
}

.el-cascader-menu {
  border: none !important;
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

.specal_slect_classs .el-select-dropdown__item.selected {
  background: #e1e1e1;
}

.colorReds input {
  border: 1px solid red;
}

.colorRedss .el-input__inner {
  border: 1px solid red;
}

.batchBtn.el-button {
  padding: 6px 10px !important;
}

.popperDiaBox /deep/ .popper__arrow {
  display: none !important;
}

.popperDiaBox.el-popover.el-popover {
  min-width: 110px !important;
}

.popperDiaBox {
  top: 125px !important;
}

.popperDiaBox /deep/ div {
  white-space: nowrap;
}

.popperDiaBox1 /deep/ .popper__arrow {
  display: none !important;
}

.widthwarp {
  display: inline-block;
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  height: 14px;
}

.app_Main /deep/ .el-table {
  position: relative;
}

.fixed_row {
  background: #e1efff !important;
}

.widthwarp {
  display: inline-block;
  max-width: 250px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  height: 14px;
}

.advanced_search {
  margin-top: 138px !important;
  position: relative;
  margin-left: 137px;
}

.advanced_search .el-dialog__header {
  padding: 0 !important;
}

.advanced_search .el-dialog__body {
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

.advanced_search /deep/ .el-dialog__body {
  padding: 10px 30px 20px !important;
  background: #ebeff3 !important;
}

.tableHight111 {
  height: calc(50vh - 60px);
  overflow: auto;
}

.tableHight111 /deep/ .el-card__body {
  padding: 10px !important;
}
</style>@/components/mixin/productJSON