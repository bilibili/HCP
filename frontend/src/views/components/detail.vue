<template>
  <div class="taskDesk_box">
    <el-drawer :visible.sync="drawerDetail" :direction="draw_direction" :modal="false" size="100%"
      :wrapperClosable="false" :modal-append-to-body="false"
      :style="tableDetailStatusC == 'halfScreen' ? 'height:50vh;top:50vh' : 'height:96%;'"
      :class="isCollapseStatus ? tableDetailStatusC == 'halfScreenBan' ? 'iscollClass3' : 'iscollClass1' : tableDetailStatusC == 'halfScreenBan' ? 'iscollClass4' : 'iscollClass2'"
      custom-class="spiceDrawLog" :before-close="handleClose">
      <div
        style="height: 100%;  cursor: pointer;  width: 15px;  background: rgb(247, 247, 247);  position: absolute;  z-index: 999;  display: flex;  justify-content: center;  align-items: center;  left: 0;  top: 0px;"
        @click="handleClose" v-if="tableDetailStatusC != 'halfScreen'">
        <i class="el-icon-arrow-right" style="color: #409eff"> </i>
      </div>
      <div
        style="width: 100%;  cursor: pointer;  height: 15px;  background: rgb(247, 247, 247);  position: absolute;  z-index: 999;  display: flex;  justify-content: center;  align-items: center;  left: 0;  top: 0px;"
        @click="handleClose" v-if="tableDetailStatusC == 'halfScreen'">
        <i class="el-icon-arrow-down" style="color: #409eff"></i>
      </div>
      <div style="padding: 10px; position: relative" v-if="Object.keys(baseConfigData).length">
        <div style="font-size: 15px;  font-weight: 600;  margin-bottom: 10px;  padding-left: 3px;  padding-top: 10px;">
          <el-breadcrumb v-if="breadcrumbArr.length" separator-class="el-icon-arrow-right">
            <el-breadcrumb-item v-for="(item, index) in breadcrumbArr" :key="index">
              {{ item }}
            </el-breadcrumb-item>
            <el-breadcrumb-item><span
                style="font-weight: 600; color: #000; font-size: 14px">详情</span></el-breadcrumb-item>
          </el-breadcrumb>
          <div style="position: absolute; right: 26px; top: 14px">
            <el-tooltip content="关闭" placement="bottom"
              style="color: red;  cursor: pointer;  font-size: 15px;  float: right;  margin-left: 10px;  margin-top: 0.5px;">
              <i class="el-icon-circle-close" @click="handleClose"></i>
            </el-tooltip>
            <el-tooltip content="右推" placement="bottom"
              style="float: right;  cursor: pointer;  margin-left: 8px;  margin-top: 2px;"
              v-if="tableDetailStatusC != 'halfScreenBan'">
              <img src="@/assets/img/right.png" width="12px" height="12px" alt="" @click="amplifyAllBan" />
            </el-tooltip>
            <el-tooltip content="下推" placement="bottom"
              style="float: right;  cursor: pointer;  margin-left: 8px;  margin-top: 2px;"
              v-if="tableDetailStatusC != 'halfScreen'">
              <img src="@/assets/img/up.png" width="12px" height="12px" alt="" @click="amplifyScreen" />
            </el-tooltip>
            <el-tooltip content="全屏" placement="bottom"
              style="float: right;  cursor: pointer;  margin-left: 8px;  margin-top: 2px;"
              v-if="tableDetailStatusC != 'halfAll'">
              <img src="@/assets/img/all.png" width="12px" height="12px" alt="" @click="amplifyAll" />
            </el-tooltip>
          </div>
        </div>
        <el-card>
          <div style="position: relative; top: 5px">
            <span style="font-size: 18px">{{ baseConfigData[showDetNameKey]?baseConfigData[showDetNameKey]:'(无名称)' }}</span>
            <span v-if="baseConfigData.status">
              <template v-for="(item, index) in statusObj">
                <template v-if="item.style">
                  <template v-if="item.style == 'info'">
                    <el-link :key="index" v-if="item.en == baseConfigData.status" :underline="false"
                      class="el-icon-remove" style="font-size: 12px; margin-left: 20px" :type="item.style">
                      <span style="margin-left: 5px">{{ item.zh }}</span>
                    </el-link>
                  </template>
                  <template v-else>
                    <el-link :key="index" v-if="item.en == baseConfigData.status" :underline="false"
                      :class="'el-icon-' + (item.style == 'danger' ? 'error' : item.style)"
                      style="font-size: 12px; margin-left: 20px" :type="item.style">
                      <span style="margin-left: 5px">{{ item.zh }}</span>
                    </el-link> </template>
                </template>
                <template v-else>
                  <el-link :key="index" v-if="item.en == baseConfigData.status" :underline="false"
                    style="font-size: 12px; margin-left: 20px" :type="item.style">
                    <span style="margin-left: 5px">{{ item.zh }}</span>
                  </el-link>
                </template>
              </template>
            </span>
            <span style="float: right" @click="objectMessageUp = !objectMessageUp">
              <i class="el-icon-arrow-down" style="cursor: pointer" v-if="!objectMessageUp"></i>
              <i class="el-icon-arrow-up" style="cursor: pointer" v-if="objectMessageUp"></i>
            </span>
          </div>
          <div style="margin-top: 15px" v-if="objectMessageUp">
            <div v-if="$route.name == 'cloudServer'">
              <el-button
                @click="rebootBtn('FormRebootCloudProduct', 'FormRebootCloudServer', [baseConfigData], '云服务器重启')">重启</el-button>
              <el-button
                @click="rebootBtn('FormReinstallCloudProduct', 'FormReinstallCloudServer', [baseConfigData], '云服务器重装')">重装</el-button>
              <el-button
                @click="rebootBtn('FormConfigCloudProduct', 'FormConfigCloudServer', [baseConfigData], '云服务器改配')">改配</el-button>
              <el-button
                @click="rebootBtn('FormDeleteCloudProduct', 'FormDeleteCloudServer', [baseConfigData], '云服务器清退')">清退</el-button>
              <el-button
                @click="rebootBtn('FormRenameBatchCloudServer', 'FormRenameCloudServer', [baseConfigData], '云服务器改名')">改名</el-button>
            </div>
          </div>
        </el-card>
        <el-card style="margin-top: 10px">
          <el-tabs v-model="activeName" @tab-click="handleClick">
            <el-tab-pane v-for="(item, index) in showConfigData.tabs" :key="index" :label="item.title"
              :name="(index + 1).toString()">
              <div v-if="item.title == '基本信息'">
                <div v-for="(citem, cindex) in item.groups" :key="cindex">
                  <h4
                    style="margin: 10px 0px 15px;  padding: 0px;  color: rgb(51, 51, 51);  font-size: 15px;  font-weight: 600;  text-indent: 20px;  border-left: 4px solid rgb(35, 173, 229);">
                    {{ citem.title }}
                  </h4>
                  <el-row :gutter="20">
                    <template v-for="(ccitem, ccindex) in citem.properties">
                      <el-col style="margin-bottom: 15px" v-if="!ccitem.hide" :span="12" :key="ccindex">
                        <div v-if="ccitem.enum" class="detBox">
                          <label>{{ ccitem.title + (ccitem.unit ? "(" + ccitem.unit + ")" : "") }} :</label>
                          <span style="display: flex">
                            <template v-for="(objitem, objindex) in ccitem.enum">
                              <span :key="objindex">
                                <span v-for="(  valueItem, valueindex) in ccitem.showValue" :key="valueindex">
                                  <template v-if="objitem.style">
                                    <template v-if="objitem.style == 'info'">
                                      <el-link :key="index" v-if="objitem.en == valueItem" :underline="false"
                                        class="el-icon-remove" style="font-size: 12px" :type="objitem.style"><span
                                          style="margin-left: 5px">{{ objitem.zh }}</span></el-link>
                                    </template>
                                    <template v-else>
                                      <el-link :key="index" v-if="objitem.en == valueItem" :underline="false"
                                        :class="'el-icon-' + (objitem.style == 'danger' ? 'error' : objitem.style)"
                                        style="font-size: 12px" :type="objitem.style">
                                        <span style="margin-left: 5px">{{ objitem.zh }}</span>
                                      </el-link>
                                    </template>
                                  </template>
                                  <template v-else>
                                    <span :key="index" v-if="objitem.en == valueItem" style="font-size: 12px"
                                      :type="objitem.style"><span style="margin-left: 5px">{{ objitem.zh
                                      }}</span></span>
                                  </template>
                                </span>
                              </span>
                            </template>
                          </span>
                        </div>
                        <div v-else class="detBox">
                          <label>{{ ccitem.title + (ccitem.unit ? "(" + ccitem.unit + ")" : "") }} :</label>
                          <span>
                            <template v-if="ccitem.link">
                              <el-link style="font-size: 12px;  margin-right: 7px;  padding-right: 20px;" type="primary"
                                @click="linkOpen(ccitem, valueindex)"
                                v-for="(  valueItem, valueindex) in ccitem.showValue" :key="valueindex">{{ valueItem }}
                                <div style="overflow: hidden;  position: absolute;  right: 0px;  top: 1px;">
                                  <img
                                    style="vertical-align: middle;  position: relative;  top: -1px;  transform: translateX(30px);  filter: drop-shadow(    #409eff -30px 0px 0px  );"
                                    src="@/assets/icon/linkTab.svg" alt="" width="15px" height="15px"
                                    v-if="valueItem" />
                                </div>
                              </el-link>
                            </template>
                            <template v-if="ccitem.customStyle && ccitem.customStyle[0].urls">
                              <span>
                                <img v-if="ccitem.showValue.join(',') == '百度云'"
                                style="vertical-align: middle;  position: relative;  top: -1px;width: 16px;height:14px;"
                                src="@/assets/icon/cloud/BDC.png" alt="" width="20px" height="20px" />
                              <img v-if="ccitem.showValue.join(',') == '腾讯云'"
                                style="vertical-align: middle;  position: relative;  top: -1px;width: 16px;height:14px;"
                                src="@/assets/icon/cloud/TCC.svg" alt="" width="15px" height="15px" />
                              <img v-if="ccitem.showValue.join(',') == '华为云'"
                                style="vertical-align: middle;  position: relative;  top: -1px;width: 16px;height:14px;"
                                src="@/assets/icon/cloud/HWC.svg" alt="" width="15px" height="15px" />
                              <img v-if="ccitem.showValue.join(',') == '阿里云'"
                                style="vertical-align: middle;  position: relative;  top: -1px;width: 16px;height:14px;"
                                src="@/assets/icon/cloud/ABC.svg" alt="" width="15px" height="15px" />
                              <img v-if="ccitem.showValue.join(',') == '亚马逊' || ccitem.showValue.join(',') == 'AWS'"
                                style="vertical-align: middle;  position: relative;  top: -1px;width: 16px;height:14px;"
                                src="@/assets/icon/cloud/AWS.png" alt="" width="15px" />
                              <span style="margin-left: 2px">{{ ccitem.showValue.join(",") }}</span>
                              </span>
                            </template>
                            <template
                              v-if="ccitem.style != 'textarea' && !ccitem.link && !ccitem.customStyle && ccitem.tagsArr">
                              <span v-if="ccitem.showValue && ccitem.showValue.length != 0">
                                <span slot="content" v-for="(tagItem, tagInx) in ccitem.showValue" :key="tagInx">
                                  <el-tooltip content="Bottom center" placement="bottom" effect="light">
                                    <div slot="content">
                                      <span v-for="(titem, key, tindex) in tagItem" :key="tindex">
                                        {{ tindex != 0 ? ":" + titem : "" + titem }}
                                      </span>
                                    </div>
                                    <el-tag
                                      style="margin-right: 3px;  margin-bottom: 3px;  max-width: 120px;  white-space: nowrap;  overflow: hidden;  text-overflow: ellipsis;"
                                      v-if="Object.keys(tagItem).length != 0">
                                      <span v-for="(titem, key, tindex) in tagItem" :key="tindex">
                                        {{ tindex != 0 ? ":" + titem : "" + titem }}
                                      </span>
                                    </el-tag>
                                  </el-tooltip>
                                </span>
                              </span>
                            </template>
                            <template v-if="ccitem.style == 'textarea'">
                              <el-input readonly :autosize="{ minRows: 2, maxRows: 200 }" type="textarea"
                                v-model="ccitem.showValue" style="width: 440px"></el-input>
                            </template>
                            <template
                              v-if="ccitem.style != 'textarea' && !ccitem.link && !ccitem.customStyle && !ccitem.tagsArr">
                              <span v-if="ccitem.isTime && ccitem.showValue.join(',') != ''">
                                {{ ccitem.showValue.join(",") | filterTimeShow }}
                              </span>
                              <span v-else>
                                {{ ccitem.showValue.join(",") }}
                              </span>
                            </template>
                            <template v-else>
                              <span v-if="!(ccitem.customStyle && ccitem.customStyle[0].urls)">{{ ccitem.showValue.join(",") }}</span>
                            </template>
                          </span>
                        </div>
                      </el-col>
                    </template>
                  </el-row>
                </div>
              </div>
              <div v-if="item.title != '基本信息'">
                <div>
                  <div style="text-align: right; margin-bottom: 5px" v-show="tableHeaderArr.length">
                    <el-dialog :modal="false" :show-close="false" :close-on-press-escape="true"
                      custom-class="advanced_searchDetial" :visible.sync="poHidess"
                      :style="{ '--dynamic-margin-top': eventSearchNumY + 'px', }">
                      <div style="background: #f8f9fa;  padding: 10px;  border-radius: 5px;">
                        <span v-for="(Onetem, leII) in advSearch" :key="leII">
                          <advancedSearchOne :searchMoreArr="searchMoreArr" :tableHeaderArr="tableHeaderArr"
                            :searchLevel="1" :searchLevelData="Onetem" :key="'aa' + leII" :ownIndex="leII"
                            @addNewsearch="addNewsearch" @delNewsearch="delNewsearch" :oneIndex="leII"
                            :advSearch="advSearch" :parentIndex="null" style="margin-top: 10px"
                            v-if="!Onetem.child || Onetem.child.length == 0"></advancedSearchOne>
                          <template v-if="Onetem.child && Onetem.child.length != 0">
                            <div v-if="advSearch[0].category"
                              style="border: 1px solid #e1e1e1;  position: relative;  margin: 10px 0;  font-size: 10px;">
                              <span
                                style="position: absolute;  top: -5px;  left: -1px;  background: #f8f9fa;  font-weight: 600;">
                                {{ advSearch[0].category }}
                              </span>
                            </div>
                            <el-card style="margin-top: 15px" class="gjSearchBody">
                              <div class="gjSearchBodyItem">
                                <div v-for="(levItem, levI) in Onetem.child" :key="levI">
                                  <advancedSearchOne :searchMoreArr="searchMoreArr" :tableHeaderArr="tableHeaderArr"
                                    :searchLevel="2" :searchLevelData="levItem" :ownIndex="levI" :oneIndex="leII"
                                    @addNewsearch="addNewsearch" @delNewsearch="delNewsearch" :parentIndex="leII"
                                    :advSearch="advSearch" v-if="!levItem.child || levItem.child.length == 0"
                                    style="margin-top: 10px"></advancedSearchOne>
                                  <div v-if="levItem.child && levItem.child.length != 0">
                                    <div v-if="advSearch[leII].child[0].category"
                                      style="border: 1px solid #e1e1e1;  position: relative;  margin: 10px 0;  font-size: 10px;">
                                      <span
                                        style="position: absolute;  top: -5px;  left: -1px;  background: #f8f9fa;  font-weight: 600;">
                                        {{ advSearch[leII].child[0].category }}
                                      </span>
                                    </div>
                                    <el-card style="margin-top: 15px" class="gjSearchBody">
                                      <div class="gjSearchBodyItem">
                                        <div v-for="(  levItem2, levI2) in levItem.child" :key="levI2">
                                          <advancedSearchOne :searchMoreArr="searchMoreArr"
                                            :tableHeaderArr="tableHeaderArr" levI levI2 :searchLevel="3"
                                            :searchLevelData="levItem2" :ownIndex="levI2" :parentIndex="levI"
                                            @addNewsearch="addNewsearch" :oneIndex="leII" @delNewsearch="delNewsearch"
                                            :advSearch="advSearch" style="margin-top: 10px"></advancedSearchOne>
                                        </div>
                                      </div>
                                    </el-card>
                                  </div>
                                </div>
                              </div>
                            </el-card>
                          </template>
                        </span>
                        <div style="text-align: right; margin-top: 15px">
                          <el-button @click="resetSearch">重 置</el-button>
                          <el-button type="primary" plain @click="searchTable(false)">搜 索</el-button>
                        </div>
                      </div>
                    </el-dialog>
                    <div style="border: 1px solid #e1e1e1;  float: left;  height: 28px;  border-radius: 4px;">
                      <el-button
                        style="float: left;  padding: 6px 10px !important;  background: #f5f7fa;  color: #409eff;  border: none;  position: relative;  z-index: 1;  border-radius: 4px 0 0 4px;  height: 26px;  border-right: 1px solid #e1e1e1;"
                        icon="el-icon-plus" @click="searchShowStatus" class="myElementMan">高级搜索</el-button>
                      <el-input style="width: 560px; float: left; height: 26px" class="batchSearchClass"
                        placeholder="支持模糊搜索" clearable v-model="batchSearchName" @clear="batchSearchNameTable"
                        @keydown.enter.native="batchSearchNameTable">
                      </el-input>
                    </div>
                    <i class="el-icon-refresh" style="margin-left: 6px;  cursor: pointer;  color: rgb(64, 158, 255);"
                      @click="searchTable('res')"></i>
                    <el-popover placement="right-start" width="300" trigger="click">
                      <div>
                        <el-checkbox-group v-model="checkboxTH" @change="showTH" size="small">
                          <el-row :gutter="20">
                            <template v-for="(item, index) in tableHeaderArr">
                              <el-col :span="12"
                                v-if="item.name != 'id' && item.name != 'project.cloudProjectAccountConfig.cloudProjectConfig.id' && item.name != 'project.cloudProjectAccountConfig.cloudProjectConfig.name' && !item.hide"
                                :key="index">
                                <el-checkbox style="width: 100%; margin-bottom: 7px" :label="item.title"
                                  border></el-checkbox>
                              </el-col>
                            </template>
                          </el-row>
                        </el-checkbox-group>
                        <el-button icon="el-icon-refresh-right" size="mini"
                          style="float: right; cursor: pointer; margin-top: 5px" @click="showTHRefresh">重置</el-button>
                      </div>
                      <i style="margin-left: 7px;  color: #409eff;  cursor: pointer;  font-size: 16px;" slot="reference"
                        class="el-icon-setting"></i>
                    </el-popover>
                  </div>
                </div>
                <div v-if="searchObjTag" style="margin: 20px 0 10px; cursor: pointer">
                  <div v-if="searchObjTag.length" style="position: relative; float: left">
                    <div v-for="(tagsValitem, tagsValindex) in searchObjTag" :key="tagsValindex"
                      style="margin-right: 20px;  background: #ecf5ff;  padding: 5px;  float: left;  margin-bottom: 10px;">
                      <span :key="index + 1" v-for="(item, index) in tagsValitem" @click="enterTagLast(tagsValindex)">
                        <span v-if="!item.child || !item.child.length">
                          <span v-if="tagsValitem[0] && tagsValitem[0].category && index != 0" class="tagTypeClass">
                            {{ tagsValitem[0].category }}
                          </span>
                          <span class="speical_tag" slot="reference" :disable-transitions="false"
                            style="margin: 0 3px; cursor: pointer">
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
                              {{ item.name }}>{{ new Date(item.value).getTime() | filterTimeShow("YYYY-DD-MM") }}
                            </span>
                            <span v-if="item.type == '_HLT'" class="widthwarp">
                              {{ item.name }}&lt;{{ new Date(item.value).getTime() | filterTimeShow }}
                            </span>
                            <span v-if="item.type == '_DLT'" class="widthwarp">
                              {{ item.name }}&lt;{{ new Date(item.value).getTime() | filterTimeShow("YYYY-DD-MM") }}
                            </span>
                            <span v-if="item.type == 'arr'" style="top: -4px; position: relative">
                              <span v-if="item.value && item.value.split('\n').length > 8">
                                {{ item.name + ":[" + item.value.split("\n")[0] + "," + item.value.split("\n")[1] +
                                  "等" + item.value.split("\n").length + "个]" }}
                              </span>
                              <span v-else>
                                {{ item.name + ":[" + item.value + "]" }}
                              </span>
                            </span>
                            <span v-if="item.type == 'selectNull'" class="widthwarp">
                              {{ item.name }}
                              <span v-if="item.value != 'no'">!</span>
                            </span>
                          </span>
                        </span>

                        <span v-else>
                          <span class="tagTypeClass"><span style="margin-right: 3px">
                              {{ tagsValitem[0].category }}
                            </span>(</span>
                          <span :key="index1 + 1" v-for="(itemTwo, index1) in item.child">
                            <span v-if="!itemTwo.child || !itemTwo.child.length">
                              <span v-if="item.child[0].category && index1 != 0" class="tagTypeClass">
                                {{ item.child[0].category }}
                              </span>
                              <span class="speical_tag" slot="reference" :disable-transitions="false"
                                style="margin: 0 3px; cursor: pointer">
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
                                  {{ itemTwo.name }}>{{ new Date(itemTwo.value).getTime() | filterTimeShow("YYYY-DD-MM")
                                  }}
                                </span>
                                <span v-if="itemTwo.type == '_HLT'" class="widthwarp">
                                  {{ itemTwo.name }}&lt;{{ new Date(itemTwo.value).getTime() | filterTimeShow }}
                                </span>
                                <span v-if="itemTwo.type == '_DLT'" class="widthwarp">
                                  {{ itemTwo.name }}&lt;{{ new Date(itemTwo.value).getTime() |
                                    filterTimeShow("YYYY-DD-MM") }}
                                </span>
                                <span v-if="itemTwo.type == 'arr'" style="top: -4px; position: relative">
                                  <span v-if="itemTwo.value && itemTwo.value.split('\n').length > 8">
                                    {{ itemTwo.name + ":[" + itemTwo.value.split("\n")[0] + "," +
                                      itemTwo.value.split("\n")[1] +
                                      "等" + itemTwo.value.split("\n").length + "个]" }}
                                  </span>
                                  <span v-else>
                                    {{ itemTwo.name + ":[" + itemTwo.value + "]" }}
                                  </span>
                                </span>
                                <span v-if="itemTwo.type == 'selectNull'" class="widthwarp">
                                  {{ itemTwo.name }}
                                  <span v-if="itemTwo.value != 'no'">!</span>
                                </span>
                              </span>
                            </span>
                            <span v-else>
                              <span class="tagTypeClass"><span style="margin-right: 3px">
                                  {{ item.child[0].category }}
                                </span>(</span>
                              <span :key="index3 + 1" v-for="(itemThere, index3) in itemTwo.child">
                                <span v-if="itemTwo.child[0].category && index3 != 0" class="tagTypeClass">{{
                                  itemTwo.child[0].category }}</span>
                                <span class="speical_tag" slot="reference" :disable-transitions="false"
                                  style="margin: 0 3px; cursor: pointer"
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
                                      filterTimeShow("YYYY-DD-MM") }}
                                  </span>
                                  <span v-if="itemThere.type == '_HLT'" class="widthwarp">
                                    {{ itemThere.name }}&lt;{{ new Date(itemThere.value).getTime() | filterTimeShow }}
                                  </span>
                                  <span v-if="itemThere.type == '_DLT'" class="widthwarp">
                                    {{ itemThere.name }}&lt;{{ new Date(itemThere.value).getTime() |
                                      filterTimeShow("YYYY-DD-MM") }}
                                  </span>
                                  <span v-if="itemThere.type == 'arr'" style="top: -4px; position: relative">
                                    <span v-if="itemThere.value && itemThere.value.split('\n').length > 8">
                                      {{ itemThere.name + ":[" + itemThere.value.split("\n")[0] + ","
                                        + itemThere.value.split("\n")[1] + "等" + itemThere.value.split("\n").length + "个]"
                                      }}
                                    </span>
                                    <span v-else>
                                      {{ itemThere.name + ":[" + itemThere.value + "]" }}
                                    </span>
                                  </span>
                                  <span v-if="itemThere.type == 'selectNull'" class="widthwarp">
                                    {{ itemThere.name }}
                                    <span v-if="itemThere.value != 'no'">!</span>
                                  </span>
                                </span>
                              </span>
                              <span class="tagTypeClass" style="margin-right: 3px">)</span>
                            </span>
                          </span>
                          <span v-if="tagsValitem[0] && tagsValitem[0].category" class="tagTypeClass"
                            style="margin-right: 3px">)</span>
                        </span>
                      </span>

                      <i class="el-icon-circle-close" style="cursor: pointer;  margin-left: 10px;  font-size: 12px;"
                        @click="refAdvancedSearch(tagsValindex)"></i>
                    </div>
                    <el-link v-if="searchObjTag.length" size="mini"
                      style="cursor: pointer;  margin-left: 10px;  font-size: 12px;  position: absolute;  bottom: 15px;  right: -15px;"
                      type="primary" @click="refAdvancedSearchLast(tagsValindex)">重置</el-link>
                  </div>
                </div>
                <el-table :data="tableData"
                  :style="isCollapseStatus ? 'width: calc(100vw - 80px);' : 'width: calc(100vw - 170px);'"
                  @selection-change="handleSelectionChangeTol" v-loading="tableLoading">
                  <el-table-column v-if="tableHeaderArr.length" type="selection" width="40"></el-table-column>
                  <template v-for="(item, index) in tableHeaderArr">
                    <template v-if="checkboxTH.indexOf(item.title) != -1">
                      <template v-if="item.enum">
                        <el-table-column :key="index" :prop="item.name"
                          :label="item.title + (item.unit ? '(' + item.unit + ')' : '')" :show-overflow-tooltip="true">
                          <template slot-scope="scope">
                            <div>
                              <span v-if="showlistObj.schema == 'JobObject'">
                                <template v-for="(objitem, objindex) in item.enum">
                                  <span :key="objindex">
                                    <template v-if="objitem.style">
                                      <template v-if="objitem.style == 'info'">
                                        <el-link :key="index"
                                          v-if="scope.row.jobStep.job && objitem.en == scope.row.jobStep.job.status"
                                          :underline="false" class="el-icon-remove" style="font-size: 12px"
                                          :type="objitem.style"><span style="margin-left: 5px">
                                            {{ objitem.zh }}</span>
                                        </el-link>
                                      </template>
                                      <template v-else>
                                        <el-link :key="index"
                                          v-if="scope.row.jobStep.job && objitem.en == scope.row.jobStep.job.status"
                                          :underline="false"
                                          :class="'el-icon-' + (objitem.style == 'danger' ? 'error' : objitem.style)"
                                          style="font-size: 12px" :type="objitem.style">
                                          <span style="margin-left: 5px">{{ objitem.zh }}
                                          </span></el-link>
                                      </template>
                                    </template>
                                    <template v-else>
                                      <el-link :key="index" v-if="objitem.en == scope.row[item.name]" :underline="false"
                                        style="font-size: 12px" :type="objitem.style">
                                        <span style="margin-left: 5px">
                                          {{ objitem.zh }}
                                        </span>
                                      </el-link>
                                    </template>
                                  </span>
                                </template>
                              </span>
                              <span v-else>
                                <template v-for="(objitem, objindex) in item.enum">
                                  <span :key="objindex">
                                    <template v-if="objitem.style">
                                      <template v-if="item.name.indexOf('.') != -1">
                                        <el-link :key="index"
                                          v-if="objitem.en == scope.row[item.name.split('.')[0]][item.name.split('.')[1]]"
                                          :underline="false" style="font-size: 12px" :type="objitem.style">
                                          <span style="margin-left: 5px">
                                            {{ objitem.zh }}
                                          </span>
                                        </el-link>
                                      </template>
                                      <template v-else>
                                        <template v-if="objitem.style == 'info'">
                                          <el-link :key="index" v-if="objitem.en == scope.row[item.name]"
                                            :underline="false" class="el-icon-remove" style="font-size: 12px"
                                            :type="objitem.style">
                                            <span style="margin-left: 5px">
                                              {{ objitem.zh }}
                                            </span>
                                          </el-link>
                                        </template>
                                        <template v-else>
                                          <el-link :key="index" v-if="objitem.en == scope.row[item.name]"
                                            :underline="false"
                                            :class="'el-icon-' + (objitem.style == 'danger' ? 'error' : objitem.style)"
                                            style="font-size: 12px" :type="objitem.style"><span
                                              style="margin-left: 5px">
                                              {{ objitem.zh }}
                                            </span>
                                          </el-link>
                                        </template>
                                      </template>
                                    </template>
                                    <template v-else>
                                      <template v-if="item.name.indexOf('.') != -1">
                                        <el-link :key="index"
                                          v-if="objitem.en == scope.row[item.name.split('.')[0]][item.name.split('.')[1]]"
                                          :underline="false" style="font-size: 12px" :type="objitem.style">
                                          <span style="margin-left: 5px">
                                            {{ objitem.zh }}
                                          </span>
                                        </el-link>
                                      </template>
                                      <template v-else>
                                        <el-link :key="index" v-if="objitem.en == scope.row[item.name]"
                                          :underline="false" style="font-size: 12px" :type="objitem.style"><span
                                            style="margin-left: 5px">
                                            {{ objitem.zh }}
                                          </span>
                                        </el-link>
                                      </template>
                                    </template>
                                  </span>
                                </template>
                              </span>
                            </div>
                          </template>
                        </el-table-column>
                      </template>
                      <template v-else>
                        <el-table-column v-if="item.type == 'object' && item.list" :key="index"
                          :label="item.title + (item.unit ? '(' + item.unit + ')' : '')" :show-overflow-tooltip="true">
                          <template slot-scope="scope">
                            <div v-if="scope.row[item.name.split('.')[0]]">
                              {{ scope.row[item.name.split(".")[0]].length }}
                            </div>
                          </template>
                        </el-table-column>
                        <template v-else>
                          <template v-if="item.customStyle">
                            <el-table-column v-if="item.style == 'lineFeedCopy'" :prop="item.name"
                              :label="item.title + (item.unit ? '(' + item.unit + ')' : '')" :key="index" width="250"
                              fixed :show-overflow-tooltip="true">
                              <template slot-scope="scope">
                                <span style="font-size: 12px;  color: #409eff;  cursor: pointer;" type="primary"
                                  @click="reviewDetail(scope.row, $route.name)" :underline="false">
                                  <span class="hoverSpan">
                                    {{ scope.row.name == "" || !scope.row.name ? "无名称" : scope.row.name }}
                                  </span>
                                </span>
                                <br />
                                <span style="color: #ccc" v-for="(eeitem, eeindex) in item.customStyle" :key="eeindex">
                                  {{ scope.row[eeitem.name.replace(/\./g, "")] }}
                                  <br />
                                </span>
                                <i style="cursor: pointer;  vertical-align: middle;  position: absolute;  top: 33%;  z-index: 1111;  right: 0px;  margin-top: -4px;"
                                  class="el-icon-copy-document" @click.stop="_copy(scope.row[item.name])"></i>
                              </template>
                            </el-table-column>
                            <el-table-column v-if="item.style == 'serverType'" :prop="item.name"
                              :label="item.title + (item.unit ? '(' + item.unit + ')' : '')" :key="index"
                              :show-overflow-tooltip="true">
                              <template slot-scope="scope">
                                <span>
                                  <span v-if="scope.row[item.customStyle[0].cpuName.replace(/\./g, '')]">
                                    {{ scope.row[item.customStyle[0].cpuName.replace(/\./g, "")] }}
                                    <span v-if="item.name.indexOf('memoryGB')">GB</span>
                                    <span v-else>MB</span>
                                  </span>
                                  <span v-if="scope.row[item.customStyle[0].memoryName.replace(/\./g, '')]">
                                    {{ scope.row[item.customStyle[0].memoryName.replace(/\./g, "")] }}MB
                                  </span>
                                  <span v-if="scope.row[item.customStyle[0].gpuName.replace(/\./g, '')]">
                                    GPU:{{ scope.row[item.customStyle[0].gpuName.replace(/\./g, "")] }}个
                                  </span>
                                  <span v-if="scope.row[item.customStyle[0].gpuModel.replace(/\./g, '')]">
                                    {{ scope.row[item.customStyle[0].gpuModel.replace(/\./g, "")] }}
                                  </span>
                                  <span v-if="scope.row[item.customStyle[0].diskSizeName.replace(/\./g, '')]">
                                    磁盘{{ scope.row[item.customStyle[0].diskSizeName.replace(/\./g, "")] }}G
                                  </span> </span><br />
                                <span>
                                  {{ scope.row[item.customStyle[1].name.replace(/\./g, "")] }}
                                </span>
                              </template>
                            </el-table-column>
                            <el-table-column v-if="!item.style" :prop="item.name"
                              :label="item.title + (item.unit ? '(' + item.unit + ')' : '')" :key="index"
                              :show-overflow-tooltip="true">
                              <template slot-scope="scope">
                                <span v-for="(eeitem, eeindex) in item.customStyle" :key="eeindex">
                                  <span style="font-size: 12px"
                                    v-if="eeitem.name && scope.row[eeitem.name.replace(/\./g, '')]">

                                    <!-- 有link链接 -->
                                    <span v-if="eeitem.link">
                                      <span style="font-size: 12px;  color: #409eff;  cursor: pointer;" type="primary"
                                        @click="linkOpen(scope.row[eeitem.linkId.replace(/\./g, '')], eeitem.link)"
                                        :underline="false">
                                        <span v-if="eeitem.prefix">
                                          {{ eeitem.prefix }}
                                        </span>
                                        {{ scope.row[eeitem.name.replace(/\./g, "")] }}</span>
                                    </span>
                                    <span v-else>
                                      <span v-if="eeitem.urls">
                                        <img v-if="scope.row[eeitem.name.replace(/\./g, '')] == '百度云'"
                                          style="vertical-align: middle;  position: relative;  top: -1px;"
                                          src="@/assets/icon/cloud/BDC.png" alt="" width="20px" height="20px" />
                                        <img v-if="scope.row[eeitem.name.replace(/\./g, '')] == '腾讯云'"
                                          style="vertical-align: middle;  position: relative;  top: -1px;"
                                          src="@/assets/icon/cloud/TCC.svg" alt="" width="20px" height="20px" />
                                        <img v-if="scope.row[eeitem.name.replace(/\./g, '')] == '华为云'"
                                          style="vertical-align: middle;  position: relative;  top: -1px;"
                                          src="@/assets/icon/cloud/HWC.svg" alt="" width="20px" height="20px" />
                                        <img v-if="scope.row[eeitem.name.replace(/\./g, '')] == '阿里云'"
                                          style="vertical-align: middle;  position: relative;  top: -1px;"
                                          src="@/assets/icon/cloud/ABC.svg" alt="" width="20px" height="20px" />
                                        <img v-if="scope.row[eeitem.name.replace(/\./g, '')] == '亚马逊' || scope.row[eeitem.name.replace(/\./g, '')] == 'AWS'"
                                          style="vertical-align: middle;  position: relative;  top: -1px;"
                                          src="@/assets/icon/cloud/AWS.png" alt="" width="18px" />
                                      </span>
                                      <span v-if="eeitem.prefix">
                                        {{ eeitem.prefix }}
                                      </span>
                                      <span v-if="eeitem.description" style="color: #ccc">
                                        {{ eeitem.description }}
                                      </span>
                                      <span v-if="eeitem.color" :style="{ color: eeitem.color }">
                                        {{ scope.row[eeitem.name.replace(/\./g, "")] }}
                                      </span>
                                      <span v-else>
                                        {{ scope.row[eeitem.name.replace(/\./g, "")] }}
                                      </span>
                                      <span v-if="eeitem.unit">
                                        {{ eeitem.unit }}
                                      </span>
                                    </span>
                                  </span>
                                  <!-- 带括号得 -->
                                  <span v-if="eeitem.subName">
                                    ({{ scope.row[eeitem.subName.replace(/\./g, "")] }})
                                  </span>
                                  <br v-if="eeitem.name && scope.row[eeitem.name.replace(/\./g, '')]" />
                                </span>
                              </template>
                            </el-table-column>
                          </template>
                          <template v-if="!item.customStyle">
                            <el-table-column v-if="item.isTime" :key="index" :prop="item.name"
                              :label="item.title + (item.unit ? '(' + item.unit + ')' : '')"
                              :show-overflow-tooltip="true">
                              <template slot-scope="scoped">
                                <div v-if="item.name == 'jobStep.job.createAt'">
                                  <template v-if="scoped.row.jobStep.job.createAt">
                                    {{ scoped.row.jobStep.job.createAt | filterTimeShow }}
                                  </template>
                                </div>
                                <span v-else>
                                  <span v-if="scoped.row[item.name]">
                                    {{ scoped.row[item.name] | filterTimeShow }}
                                  </span>
                                </span>
                              </template>
                            </el-table-column>
                            <template v-if="item.userMEP">
                              <el-table-column v-if="users.indexOf(username) != -1" :key="index" :prop="item.name"
                                :label="item.title + (item.unit ? '(' + item.unit + ')' : '')"
                                :show-overflow-tooltip="true">
                                <template slot-scope="scoped">
                                  <span v-if="item.name.indexOf('.') != -1">
                                    {{ scoped.row[item.name.replace(/\./g, "")] }}
                                  </span>
                                  <span v-else>
                                    {{ scoped.row[item.name] }}
                                  </span>
                                </template>
                              </el-table-column>
                            </template>
                            <el-table-column v-if="!item.userMEP && !item.isTime" :key="index" :prop="item.name"
                              :label="item.title + (item.unit ? '(' + item.unit + ')' : '')"
                              :show-overflow-tooltip="true">
                              <template slot-scope="scoped">
                                <span v-if="item.discountPrice">
                                  {{ scoped.row.discount * scoped.row.originalPrice }}
                                </span>
                                <span v-else>
                                  <span v-if="item.name.indexOf('.') != -1">
                                    {{ scoped.row[item.name.replace(/\./g, "")] }}
                                  </span>
                                  <span v-else>
                                    <span>{{ scoped.row[item.name] }}</span>
                                  </span>
                                </span>
                              </template>
                            </el-table-column>
                          </template>
                        </template>
                      </template>
                    </template>
                  </template>
                  <el-table-column v-if="tableHeaderArr.length" fixed="right" label="操作" align="center" width="140">
                    <template slot="header" slot-scope="scope">
                      <label v-show="scope">操作</label>
                    </template>
                    <template slot-scope="scope">
                      <div>
                        <span v-if="showlistObj.link">
                          <el-link
                            style="cursor: pointer;  color: #1890ff;  font-size: 12px;  border: none;  margin: 0 5px;"
                            @click="reviewDetail(scope.row)">详情</el-link>
                        </span>
                      </div>
                    </template>
                  </el-table-column>
                </el-table>
                <div style="margin-top: 10px; text-align: center" v-show="tableHeaderArr.length">
                  <el-pagination @current-change="changehandleCurrent" @size-change="handleSizeChange"
                    :current-page="pageNum" :page-size="pageSize" :page-sizes="[10, 20, 30, 100, 150, 200]"
                    layout="total, sizes, prev, pager, next, jumper" :total="total">
                  </el-pagination>
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </div>
    </el-drawer>
    <!-- 新建任务弹框 infoAddAttr 表格行信息 infoAddForm表单提交信息 reviewStatus弹框状态-->
    <div v-if="reviewStatus222">
      <formtem-Dialog :nodeCheckName="nodeCheckName" :detailOrLook="detailOrLook" @addClose="addClose"
        :infoAddAttr="infoReviewAttr" :detailOrsigle="detailOrsigle" :infoAddForm="infoReviewForm"
        :reviewStatus="reviewStatus222" :newTaskName="newTaskName" :historyTableList="historyTableList"
        :historyparamreferenceShow="historyparamreferenceShow" :baseConfigData="baseConfigData"
        :childNodename="childNodename" :checkList="checkList" :FormTemplatjobMode="FormTemplatjobMode"></formtem-Dialog>
    </div>
  </div>
</template>

<script>
import formtemDialog from "@/views/components/formTemp";
import { userMixin } from "@/components/mixin/user";
import { JSONConfig } from "@/components/mixin/JSONConfig";
import { detailMixin } from "@/components/mixin/detailJSON";
import advancedSearchOne from "@/views/components/advancedSearchOne";
import { advSearchHandle } from "@/components/mixin/advSearchHandle";
import { mapGetters } from "vuex";
import Http from "@/components/api/services";
//例如：import 《组件名称》 from '《组件路径》';
export default {
  //import引入的组件需要注入到对象中才能使用
  props: {
    jobIds: {
      type: Number,
    },
    tableDetailStatus: {
      type: String,
    },
  },
  components: {
    formtemDialog,
    advancedSearchOne,
  },
  mixins: [userMixin,JSONConfig, detailMixin, advSearchHandle],
  data() {
    //这里存放数据
    return {
      advanced_distanceTop: "",
      advanced_distanceleft: "163px",
      advSearchIndex: null,
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
      enumAllArr: [],
      draw_direction: "btt",
      objectMessageUp: false,
      tableDetailStatusC: "",
      FormTemplatjobMode: "",
      tableLoading: false,
      dateTimeValue: [],
      pickerOptions: {
        shortcuts: [
          {
            text: "最近半小时",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - (3600 * 1000) / 2);
              picker.$emit("pick", [start, end]);
            },
          },
          {
            text: "最近6小时",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 6);
              picker.$emit("pick", [start, end]);
            },
          },
          {
            text: "最近12小时",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 12);
              picker.$emit("pick", [start, end]);
            },
          },
          {
            text: "最近1天",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24);
              picker.$emit("pick", [start, end]);
            },
          },
          {
            text: "最近7天",
            onClick(picker) {
              const end = new Date();
              const start = new Date();
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7);
              picker.$emit("pick", [start, end]);
            },
          },
        ],
      },
      detail_IDQ: "",
      drawerDetail: true,
      isCollapseStatus: false,
      eventSearchNumY: 0,
      reviewStatus: false,
      activeName: "1",
      activeNameNow: "1",
      showConfigData: [],
      baseConfigData: {},
      propertySchema: "",
      tableData: [],
      tableHeaderArr: [],
      checkboxTH: [],
      checkboxTHBF: [],
      searchObj: {},
      batchSearchName: "",
      pageNum: 1,
      pageSize: 10,
      total: 0,
      showlistObj: {},
      username: "",
      reviewStatus222: false,
      detailOrLook: false,
      detailOrsigle: "",
      newTaskName: "",
      historyTableList: {},
      historyparamreferenceShow: [],
      nodeCheckName: "",
      nodeCheckProductName: "",
      dialogtitle: "",
      infoReviewForm: {},
      infoReviewAttr: [],
      dialogReview: false,
      statusObj: [],
      breadcrumbArr: [],
      childNodename: "",
      checkList: [],
      jobbg: false,
      childDetailCheckList: [],
      searchFormData: {},
      searchObjTag: [],
      classWidth: "",
      searchMoreArr: [],
      cascaderMoreArr: [],
      poHidess: false,
      tagStatus: false,
      tagTitle: "",
      tagAllArr: [],
      batchHostarr1: [],
      showDetNameKey:""
    };
  },
  watch: {
    jobIds: {
      handler: function () {
        this.drawObjectDetail = false;
        this.getRowObj();
      },
      deep: true,
    },
  },
  computed: {
    ...mapGetters(["loginUserName", "formPolicyDomain"]),
  },
  //方法集合
  methods: {
    addNewsearch(type, level, ownIndex, parentIndex, oneIndex, obj) {
      if (level == 1) {
        if (this.advSearch[0].category == "") {
          this.advSearch[0].category = type;
          this.advSearch.push({
            name: "",
            type: "",
            val: "",
            category: type,
            index: [],
            child: [],
            objContent: {},
          });
        } else {
          if (this.advSearch.length == 1) {
            this.advSearch[0].category = type;
            this.advSearch.push({
              name: "",
              type: "",
              val: "",
              category: type,
              index: [],
              child: [],
              objContent: {},
            });
          } else {
            if (this.advSearch[ownIndex].category == type) {
              this.advSearch.push({
                name: "",
                type: "",
                val: "",
                category: type,
                index: [],
                child: [],
                objContent: {},
              });
            } else {
              let arr = [];
              arr.push({
                name: obj.name,
                type: obj.type,
                val: obj.val,
                category: type,
                index: obj.index,
                child: [],
                objContent: obj.objContent,
              });
              arr.push({
                name: "",
                type: "",
                val: "",
                category: type,
                index: [],
                child: [],
                objContent: {},
              });
              this.advSearch[ownIndex].child = arr;
            }
          }
        }
      }
      if (level == 2) {
        if (this.advSearch[parentIndex].child[ownIndex].category == type) {
          this.advSearch[parentIndex].child.push({
            name: "",
            type: "",
            val: "",
            category: type,
            index: [],
            child: [],
            objContent: {},
          });
        } else {
          let arr = [];
          arr.push({
            name: obj.name,
            type: obj.type,
            val: obj.val,
            category: type,
            index: obj.index,
            objContent: obj.objContent,
          });
          arr.push({
            name: "",
            type: "",
            val: "",
            category: type,
            index: [],
            objContent: {},
          });
          this.advSearch[parentIndex].child[ownIndex].child = arr;
        }
      }
      if (level == 3) {
        this.advSearch[oneIndex].child[parentIndex].child.push({
          name: "",
          type: "",
          val: "",
          category:
            this.advSearch[oneIndex].child[parentIndex].child[0].category,
          index: [],
          objContent: {},
        });
      }
      this.advSearch = JSON.parse(JSON.stringify(this.advSearch));
    },
    delNewsearch(level, ownIndex, parentIndex, oneIndex, obj) {
      let advSearchData = JSON.parse(JSON.stringify(this.advSearch));
      if (level == 1) {
        if (advSearchData.length == 1 && ownIndex == 0 && (!obj.child || !obj.child.length)) {
          return;
        }
        advSearchData.splice(ownIndex, 1);
        if (advSearchData.length == 1) {
          if (advSearchData[0].child && advSearchData[0].child.length) {
            advSearchData = advSearchData[0].child;
          }
        }
      }
      if (level == 2) {
        advSearchData[parentIndex].child.splice(ownIndex, 1);
        if (advSearchData[parentIndex].child.length == 1) {
          advSearchData[parentIndex] = {
            name: advSearchData[parentIndex].child[0].name,
            type: advSearchData[parentIndex].child[0].type,
            val: advSearchData[parentIndex].child[0].val,
            category: advSearchData[parentIndex].category,
            index: advSearchData[parentIndex].child[0].index,
            objContent: advSearchData[parentIndex].child[0].objContent,
            child: advSearchData[parentIndex].child[0].child ? advSearchData[parentIndex].child[0].child : [],
          };
        }
      }
      if (level == 3) {
        advSearchData[oneIndex].child[parentIndex].child.splice(ownIndex, 1);
        if (advSearchData[oneIndex].child[parentIndex].child.length == 1) {
          advSearchData[oneIndex].child[parentIndex] = {
            name: advSearchData[oneIndex].child[parentIndex].child[0].name,
            type: advSearchData[oneIndex].child[parentIndex].child[0].type,
            val: advSearchData[oneIndex].child[parentIndex].child[0].val,
            category: advSearchData[oneIndex].child[parentIndex].category,
            index: advSearchData[oneIndex].child[parentIndex].child[0].index,
            objContent: advSearchData[oneIndex].child[parentIndex].child[0].objContent,
            child: advSearchData[oneIndex].child[parentIndex].child[0].child ? advSearchData[oneIndex].child[parentIndex].child[0].child : [],
          };
        }
      }
      this.advSearch = JSON.parse(JSON.stringify(advSearchData));
      this.$forceUpdate();
    },
    amplifyAll() {
      this.tableDetailStatusC = "halfAll";
      this.draw_direction = "rtl";
      this.$emit("tableHeight", "", "center", false);
      sessionStorage.setItem("tableDetailStatusC", this.tableDetailStatusC);
    },
    amplifyScreen() {
      this.draw_direction = "btt";
      this.tableDetailStatusC = "halfScreen";
      this.$emit("tableHeight", "up", "center", false);
      sessionStorage.setItem("tableDetailStatusC", this.tableDetailStatusC);
    },
    amplifyAllBan() {
      this.tableDetailStatusC = "halfScreenBan";
      this.draw_direction = "rtl";
      this.$emit("tableHeight", "", "left", true);
      sessionStorage.setItem("tableDetailStatusC", this.tableDetailStatusC);
    },
    // 标签管理弹框关闭
    addCloseTag(val) {
      this.tagStatus = val;
      this.showTHRefresh();
    },
    //高级搜索,点击菜单发生变化时
    searchForm(item) {
      this.searchFormData = {};
      this.tableHeaderArr.map((key) => {
        if (key.name == this.cascaderMorestr) {
          key.indexSearch = item;
          this.searchFormData = JSON.parse(JSON.stringify(key));
        }
      });
    },
    //高级搜索,默认选择第一个子菜单时
    searchFormchange() {
      this.cascaderMorestrSecond = null;
      this.searchFormData = {};
      this.searchMoreArr.map((item, index) => {
        if (item.name == this.cascaderMorestr) {
          this.searchMoreArrSecond = item.children;
        }
      });
      let firstStr = this.cascaderMorestr;
      this.searchMoreArr.map((key) => {
        if (key.name == firstStr) {
          twoStr = key.children[0].value;
        }
      });
    },
    searchShowStatus(ev) {
      this.eventSearchNumY = ev.clientY + 10;
      const distance = this.$root.$el.querySelector(".myElementMan");
      const distanceTop = distance.getBoundingClientRect().top;
      const distanceLeft = distance.getBoundingClientRect().left;
      this.advanced_distanceTop = Math.ceil(distanceTop) + 35 + "px";
      var advancedClass = document.getElementsByClassName("advanced_searchDetial")[0];
      advancedClass.style.marginLeft = Math.ceil(distanceLeft) - 15 + "px";
      this.advSearchIndex = null;
      this.poHidess = true;
      this.advSearch = [
        {
          name: "",
          type: null,
          val: "",
          category: "",
          index: [],
          child: [],
          objContent: {},
        },
      ];
    },
    entrySreach(item, obj) {
      let oldname = item.enName.split("+")[0];
      let oldsecond = item.enName.split("+")[1] ? item.enName.split("+")[1] : "";
      this.cascaderMorestr = oldname;
      this.searchFormData = {};
      this.searchMoreArr.map((item, index) => {
        if (item.name == this.cascaderMorestr) {
          this.searchMoreArrSecond = item.children;
        }
      });
      let firstStr = this.cascaderMorestr;
      this.searchMoreArr.map((key) => {
        if (key.name == firstStr) {
          twoStr = key.children[0].value;
        }
      });
      this.cascaderMorestrSecond = oldsecond;
      this.tableHeaderArr.map((key) => {
        if (key.name == this.cascaderMorestr) {
          key.indexSearch = this.cascaderMorestrSecond;
          this.searchFormData = JSON.parse(JSON.stringify(key));
        }
      });
      this.poHidess = true;
    },
    resetSearch() {
      this.advSearch = [
        {
          name: "",
          type: null,
          val: "",
          category: "",
          index: [],
          child: [],
          objContent: {},
        },
      ];
    },
    enterTagLast(index) {
      this.advSearchIndex = index;
      this.poHidess = true;
      this.advSearch = this.lastAdvSearch[index];
    },
    refAdvancedSearch(index) {
      this.searchObjTag.splice(index, 1);
      this.lastAdvSearch.splice(index, 1);
      this.getTable();
    },
    refAdvancedSearchLast() {
      this.searchObjTag = [];
      this.lastAdvSearch = [];
      this.getTable();
    },
    addClose(val) {
      this.reviewStatus222 = val;
      this.getTable();
    },
    // 任务类型点击taskname英文名称，请求接口 chinaName中文名称 id项目创建的id
    changeAddnew(taskname, chinaName, id) {
      this.dialogVisible = false;
      //中文名称
      this.newTaskName = chinaName;
      this.nodeCheckName = taskname;
      this.reviewStatus222 = true;
      this.detailOrsigle = this.baseConfigData.id;
    },
    showTH() {
      let checkboxName = `${this.$route.name}${this.activeName}checkboxTHBFArrdetail`;
      localStorage.setItem(checkboxName, this.checkboxTH);
      if (this.checkboxTH.length == 0) {
        this.checkboxTH = JSON.parse(JSON.stringify(this.checkboxTHBF));
      } else {
        this.checkboxTHBF = JSON.parse(JSON.stringify(this.checkboxTH));
      }
    },
    showTHRefresh() {
      let checkboxName = `${this.$route.name}${this.activeName}checkboxTHBFArrdetail`;
      localStorage.removeItem(checkboxName);
      this.getListSearchStr(this.showConfigData.tabs[this.activeName - 1]);
    },
    handleClose() {
      this.$router.replace({
        path: location.pathname,
      });
      this.drawerDetail = false;
      this.$emit("addCloseDraw");
    },
    handleClick() {
      this.searchObjTag = [];
      this.lastAdvSearch = [];
      if (this.activeNameNow != this.activeName) {
        this.activeNameNow = this.activeName;
        if (this.activeName != "1") {
          this.total = 0;
          this.jobbg = this.showConfigData.tabs[this.activeName - 1].jobbg ? this.showConfigData.tabs[this.activeName - 1].jobbg : false;
          this.getListSearchStr(this.showConfigData.tabs[this.activeName - 1], this.activeName);
        }
      }
      this.$router.replace({
        path: location.pathname,
        query: {
          detail_ID: this.$route.query.detail_ID,
          activeName: this.activeName,
        }
      });
      this.childDetailCheckList = [];
    },
    getRowObj() {
      this.showDetNameKey = this.showConfigData.name ? this.showConfigData.name.replace(/\./g, "") : 'name'
      this.detail_IDQ = this.$route.query.detail_ID;
      this.activeNameNow = "1";
      let postData = {
        schema: this.showConfigData.schema,
        where: {
          id: this.jobIds
        }
      };
      Http.getQueryList(postData).then((response) => {
        this.tableData = [];
        if (response.data.data.data) {
          this.baseConfigData = response.data.data.data[0];
          this.showConfigData.tabs[0].groups.map((k, i) => {
            k.properties.map((k1, i1) => {
              let arr = k1.name.split(".");
              let tempObj = [JSON.parse(JSON.stringify(this.baseConfigData))];
              let arr1 = [];
              for (let i2 of arr) {
                let aa1 = [];
                arr1 = [];
                for (let i3 of tempObj) {
                  let aa = i3[i2] ? i3[i2] : "";
                  if (aa instanceof Array) {
                    aa1.push(...aa);
                  } else {
                    if (k1.type == "integer") {
                      aa = aa == "" ? "0" : aa;
                    }
                    aa1.push(aa);
                  }
                  if (k1.link) {
                    arr1.push(i3["id"]);
                  }
                }
                tempObj = aa1;
              }
              if (k1.style == "textarea") {
                k1.showValue = tempObj.join(",");
              } else {
                k1.showValue = tempObj;
              }
              if ((k1.title.indexOf('时间') != -1 || k1.isTime) && !k1.enum) {
                if(tempObj.join(",")){
                  k1.showValue = [this.$moment(tempObj.join(',')).format('YYYY-MM-DD hh:mm:ss')];
                }
              }
              k1.linkID = arr1;
              this.baseConfigData[k1.name.replace(/\./g, "")] = tempObj.join(",");
            });
          });
          this.showConfigData = JSON.parse(
            JSON.stringify(this.showConfigData)
          );
          this.activeName = this.$route.query.activeName ? this.$route.query.activeName : "1";
          if (this.activeName) {
            this.handleClick();
          }
        }
      }
      );
    },
    getListSearchStr(obj, activeNum) {
      if (!obj.properties) {
        return false;
      }
      this.searchMoreArr = [];
      this.checkboxTH = [];
      this.tableHeaderArr = obj.properties;
      this.tableHeaderArr.map((k, i) => {
        k.title = k.title.replace(/\./g, "");
        this.searchObj[k.name] = "";
        if (k.name != "id" && k.default) {
          this.checkboxTH.push(k.title);
        }
        if (k.schema) {
          this.getSelList(k);
        }
        if (k.name && k.name != "id" && !k.hide) {
          this.searchMoreArr.push(k);
        }
        if (k.name) {
          if (k.tagsArr) {
            let str = "";
            k.tagsArr.map((item) => {
              str = k.name + "." + item;
              this.searchObj[str] = "";
            });
          } else {
            this.searchObj[k.name] = "";
          }
        }
        if (k.enum) {
          k.enum.map((euitem, eni) => {
            let valTempStr = `${k.name}:::${euitem.zh}`;
            this.enumAllArr[valTempStr] = euitem.en;
          });
        }
      });
      this.searchMoreArr.map((item) => {
        item.label = item.title;
        item.value = item.name;
        if (item.enum) {
          item.children = [
            {
              value: "select",
              label: "选择",
            },
          ];
        } else {
          if (item.type == "number" || item.type == "integer") {
            item.children = [
              {
                value: "",
                label: "精确",
              },
              {
                value: "_GT",
                label: "大于",
              },
              {
                value: "_LT",
                label: "小于",
              },
            ];
          } else if (item.type == "string") {
            let arr = [];
            if (item.index) {
              if (item.index.indexOf("hash") != -1) {
                arr.push(
                  {
                    value: "",
                    label: "精确",
                  },
                  {
                    value: "arr",
                    label: "批量",
                  }
                );
              }
              if (item.index.indexOf("regexp") != -1) {
                arr.push({
                  value: "_REGEX",
                  label: "模糊",
                });
              }
              if (item.index.indexOf("hour") != -1) {
                arr.push(
                  {
                    value: "_HGT",
                    label: "大于",
                  },
                  {
                    value: "_HLT",
                    label: "小于",
                  }
                );
              }
              if (item.index.indexOf("day") != -1) {
                arr.push(
                  {
                    value: "_DGT",
                    label: "大于",
                  },
                  {
                    value: "_DLT",
                    label: "小于",
                  }
                );
              }
              item.children = arr;
            } else {
              item.children = [
                {
                  value: "",
                  label: "精确",
                },
              ];
            }
          } else {
            item.children = [
              {
                value: "",
                label: "精确",
              },
            ];
          }
        }
        item.children.push({
          value: "selectNull",
          label: "空值",
        });
      });
      this.searchObj = JSON.parse(JSON.stringify(this.searchObj));
      let checkboxName = `${this.$route.name}${this.activeName}checkboxTHBFArrdetail`;
      let pArr = localStorage.getItem(checkboxName);
      if (pArr) {
        if (pArr.indexOf(",") == -1) {
          this.checkboxTH = [pArr];
        } else {
          this.checkboxTH = pArr.split(",");
        }
      } else {
        this.checkboxTH = this.checkboxTH;
      }
      this.checkboxTHBF = JSON.parse(JSON.stringify(this.checkboxTH));
      this.showlistObj = obj;
      this.getTable(activeNum);
    },
    searchTable(type) {
      if (this.advSearch[0].val || this.advSearch[0].category) {
        if (this.advSearchIndex || this.advSearchIndex === 0) {
          this.lastAdvSearch[this.advSearchIndex] = JSON.parse(JSON.stringify(this.advSearch));
        } else {
          this.lastAdvSearch.push(JSON.parse(JSON.stringify(this.advSearch)));
        }
        var searchObjTagLast = JSON.parse(JSON.stringify(this.lastAdvSearch));
        searchObjTagLast.map((lastItem, lastI) => {
          lastItem.map((advItem, advI) => {
            if (advItem.child && advItem.child.length) {
              advItem.child.map((advItem2, adv2) => {
                if (advItem2.child && advItem2.child.length) {
                  advItem2.child.map((advItem3, adv3) => {
                    if (!advItem3.child || !advItem3.child.length) {
                      lastItem[advI].child[adv2].child[adv3] = this.searchTagChange(advItem3);
                    }
                  });
                } else {
                  lastItem[advI].child[adv2] = this.searchTagChange(advItem2);
                }
              });
            } else {
              lastItem[advI] = this.searchTagChange(advItem);
            }
          });
        });
        this.searchObjTag = searchObjTagLast;
      }
      if (!type) {
        this.pageNum = 1;
        this.pageSize = 10;
      }
      this.getTable();
    },
    searchTagChange(addvitem) {
      let eItemTxt = "";
      let searchObjTagss = [];
      if (addvitem.val) {
        this.tableHeaderArr.map((item) => {
          if (item.name == addvitem.name) {
            if (item.enum) {
              eItemTxt += "[";
              item.enum.map((eitem) => {
                if (Array.isArray(addvitem.val)) {
                  addvitem.val.map((kkey, kki) => {
                    if (eitem.en == kkey) {
                      eItemTxt += `${eitem.zh},`;
                    }
                  });
                } else {
                  if (eitem.en == addvitem.val) {
                    eItemTxt = eitem.zh;
                  }
                }
              });
              eItemTxt = eItemTxt.slice(0, eItemTxt.length - 1);
              eItemTxt += "]";
            } else {
              eItemTxt = addvitem.val;
            }
            searchObjTagss = {
              name: addvitem.objContent.label,
              value: eItemTxt,
              type: addvitem.type ? addvitem.type : "select",
              category: addvitem.category,
              child: addvitem.child,
            };
          }
        });
      }
      return searchObjTagss;
    },
    // 搜索tag移除、
    handleCloseTag(item) {
      this.searchObj[item.enName] = "";
      this.searchTable();
    },
    batchSearchNameTable() {
      this.pageNum = 1;
      this.pageSize = 10;
      this.getTable();
    },
    handleSizeChange(val) {
      this.pageSize = val;
      this.pageNum = 1;
      this.getTable();
    },
    changehandleCurrent(val) {
      this.pageNum = val;
      this.getTable();
    },
    getSelList(obj, value) {
      let postData = {};
      let keyV = `${obj.name.split(".")[obj.name.split(".").length - 1]}`;
      if (value && value != "") {
        postData = {
          schema: obj.schema,
          where: {
            [`${keyV}_REGEX:`]: value,
          },
        }
      } else {
        postData = {
          schema: obj.schema,
        }
      }
      Http.getQueryList(postData).then((response) => {
        if (response.data.data.data) {
          if (obj.DataList) {
            obj.DataList = response.data.data.data;
          }
        }
      }
      );
    },
    getTable(activeNum) {
      this.tableLoading = true;
      var name = this.showlistObj.schema;
      this.showlistObj.where.replace("$id", this.baseConfigData.id)
      this.tableData = [];
      var postData = {
        schema: name,
        page_size:this.pageSize,
        page_num:this.pageNum,
        order:"id DESC",
        where:{}
      }
      var or = {}
      if (this.batchSearchName != '') {
        var showlistObjSearch = []
        if (this.showlistObj.search) {
          showlistObjSearch = JSON.parse(JSON.stringify(this.showlistObj.search))
        } else {
          showlistObjSearch = ['name']
        }
        var arr = this.batchSearchName.split(" ").filter(item => item)
        showlistObjSearch.map((item) => {
            var searKey1 = item.split('.').reverse()
            var str11 = ''
            searKey1.map((k1, i1) => {
              if (i1 == 0) {
                if(arr.length>1){
                  str11 = `{"${k1}_IN":${JSON.stringify(arr)}}`
                }else{
                  str11 = `{"${k1}_REGEX":"${arr.join()}"}`
                }
              } else {
                str11 = `{"${k1}":${str11}}`
              }
            })
            var getKet = Object.keys(JSON.parse(str11))[0]
            or[getKet] = JSON.parse(str11)[getKet]
        })
      }
      var and = [];
      if (this.lastAdvSearch.length) {
        this.lastAdvSearch.map((lastItem, lastI) => {
          if (lastItem[0].category) {
            let firstType = lastItem[0].category;
            var arr1 = {
              [firstType]: [],
            };
            lastItem.map((advItem, advI) => {
              if (advItem.child && advItem.child.length) {
                let secondType = advItem.child[0].category;
                var arr2 = {
                  [secondType]: [],
                };
                advItem.child.map((advItem2, adv2) => {
                  if (advItem2.child && advItem2.child.length) {
                    let thirdType = advItem2.child[0].category;
                    var arr3 = {
                      [thirdType]: [],
                    };
                    advItem2.child.map((advItem3, adv3) => {
                      if (!advItem3.child || !advItem3.child.length) {
                        if (!advItem3.val) {
                          return;
                        }
                        if (Object.prototype.toString.call(advItem3.val).indexOf("Array") != -1) {
                          if (advItem3.val.length == 0) {
                            return;
                          }
                        }
                        arr3[thirdType].push(
                          this.advSearchHandle(advItem3, name)
                        );
                      }
                    });
                    arr2[secondType].push(arr3);
                  } else {
                    if (!advItem2.val) {
                      return;
                    }
                    if (Object.prototype.toString.call(advItem2.val).indexOf("Array") != -1) {
                      if (advItem2.val.length == 0) {
                        return;
                      }
                    }
                    arr2[secondType].push(this.advSearchHandle(advItem2, name));
                  }
                });
                arr1[firstType].push(arr2);
              } else {
                if (!advItem.val) {
                  return;
                }
                if (Object.prototype.toString.call(advItem.val).indexOf("Array") != -1) {
                  if (advItem.val.length == 0) {
                    return;
                  }
                }
                arr1[firstType].push(this.advSearchHandle(advItem, name));
              }
            });
            and.push(arr1);
          } else if (!lastItem[0].category && lastItem[0].val) {
            if (Object.prototype.toString.call(lastItem[0].val).indexOf("Array") != -1) {
              if (lastItem[0].val.length == 0) {
                return;
              }
            }
            let arr1 = this.advSearchHandle(lastItem[0], name);
            and.push(arr1);
          }
        });
      }
      if(Object.keys(or).length){
        postData.where.or = or
      }
      if(Object.keys(and).length){
        var arrAnd = []
        and.map((k1, i1) => {
          if(k1.AND){
            arrAnd.push(...k1.AND)
          }else{
            arrAnd.push(k1)
          }
        })
        postData.where.and = arrAnd
      }
      Http.getQueryList(postData).then((response) => {
        if (activeNum && activeNum != this.activeName) {
          return;
        }
        this.tableData = [];
        this.tableLoading = false;
        this.total = response.data.data.total;
        this.tableData = response.data.data.data;
        this.tableHeaderArr.map((k1, i1) => {
          if(k1.name!='id'){
            this.tableData.map((k2, i2) => {
              let arr = k1.name.split(".");
              if (arr.length > 1) {
                if (!k2[arr[0]]) {
                  k2[arr[0]] = {};
                }
              }
              let tempObj = [JSON.parse(JSON.stringify(k2))];
              for (let i2 of arr) {
                let aa1 = [];
                for (let i3 of tempObj) {
                  let aa = i3[i2] ? i3[i2] : "";
                  if (aa instanceof Array) {
                    aa1.push(...aa);
                  } else {
                    if (k1.type == "integer") {
                      aa = aa == "" ? "0" : aa;
                    }
                    aa1.push(aa);
                  }
                }
                tempObj = aa1;
              }
              if ((k1.title.indexOf("时间") != -1 || k1.isTime) && !k1.enum) {
                if (k1.isTimestamp) {
                  if (tempObj.join(",")) {
                    k2[k1.name.replace(/\./g, "")] = this.$moment(Number(tempObj.join(",")) * 1000).format("YYYY-MM-DD HH:mm:ss");
                  }
                } else {
                 
                  if (tempObj.join(",") && tempObj.join(",").length == 10) {
                    k2[k1.name.replace(/\./g, "")] = tempObj.join(",") ? this.$moment(tempObj.join(",")).format("YYYY-MM-DD") : tempObj.join(",");
                  } else {
                    k2[k1.name.replace(/\./g, "")] = tempObj.join(",") ? this.$moment(tempObj.join(",")).format("YYYY-MM-DD HH:mm:ss") : tempObj.join(",");
                  }
                }
              } else {
                k2[k1.name.replace(/\./g, "")] = tempObj.join(",");
              }
            });
          }
        }).catch(() => {
            this.tableLoading = false;
          });
      })
        .catch(() => {
          this.tableLoading = false;
        });
      this.poHidess = false;
      this.advSearch = [
        {
          name: "",
          type: null,
          val: "",
          category: "",
          index: [],
          child: [],
          objContent: {},
        },
      ];
    },
    handleSelectionChangeTol(val) {
      this.childDetailCheckList = val;
    },
    rebootBtnNew(nodeCreate, title, obj) {
      this.childNodename = ''
      this.checkList = []
      this.nodeCheckName = nodeCreate
      this.newTaskName = title
      this.nodeCheckProductName = nodeCreate
      this.dialogtitle = title
      this.reviewStatus222 = true
    },
    rebootBtn(nodeCreate, nodename, rowArr, title) {
      if (nodeCreate != "FormDeleteCloudDNSRecord" && nodeCreate != "FormDeleteCloudPrefixListByManage" && nodeCreate != "FormUpdateCloudPrefixListByManage" && nodeCreate != "FormBatchCreateCloudPrefixList" && nodeCreate != "FormUpdateCloudQuotaConfig") {
        if (nodeCreate != "FormDeleteCloudProductKubeNodePool") {
          if (!rowArr[0].project) {
            this.$message({
              showClose: true,
              message: "所属项目为空，请联系管理员！",
              type: "warning",
            });
            return;
          }
        } else {
          if (!rowArr[0].cluster.project) {
            this.$message({
              showClose: true,
              message: "所属项目为空，请联系管理员！！",
              type: "warning",
            });
            return;
          }
        }
      }
      if (nodeCreate == "FormDeleteCloudPrefixListByManage" || nodeCreate == "FormUpdateCloudPrefixListByManage") {
        let proArr = [];
        rowArr.map((item, index) => {
          console.log(item);
          if (proArr.indexOf(item.account.alias) == -1) {
            proArr.push(item.account.alias);
          }
        });
        if (proArr.length != 1) {
          this.$message({
            showClose: true,
            message: "请选择同一云账号！",
            type: "warning",
          });
          return;
        }
      }
      this.nodeCheckName = nodeCreate;
      this.nodeCheckProductName = nodeCreate;
      this.dialogtitle = title;
      this.newTaskName = title;
      this.childNodename = nodename;
      this.checkList = rowArr;
      this.reviewStatus222 = true;
      this.detailOrsigle = this.baseConfigData.id;
    },
    reProBtn(nodeCreate, nodename, rowArr, title) {
      this.nodeCheckName = nodeCreate;
      this.nodeCheckProductName = nodeCreate;
      this.dialogtitle = title;
      this.newTaskName = title;
      this.checkList = rowArr;
      this.reviewStatus222 = true;
      this.detailOrsigle = this.baseConfigData.id;
    },
    resetForm(formName) {
      if (this.$refs[formName]) {
        this.$refs[formName].resetFields();
      }
    },
    linkOpen(obj, i) {
      window.open(`${location.origin}${obj.link}?detail_ID=${obj.linkID[i]}`);
    },
    reviewDetail(obj) {
      if (!this.jobbg) {
        window.open(`${location.origin}${this.showlistObj.link}?detail_ID=${obj.id}`);
      } else {
        if (obj.jobStep && obj.jobStep.job) {
          window.open(`${location.origin}${this.showlistObj.link}?jobId=${obj.jobStep.job.id}`);
        }
        if (!obj.jobStep && obj.id) {
          window.open(`${location.origin}${this.showlistObj.link}?jobId=${obj.id}`);
        }
      }
    },
  },
  //生命周期 - 挂载完成（可以访问DOM元素）
  created() {
    this.tableDetailStatusC = JSON.parse(JSON.stringify(this.tableDetailStatus));
    this.username = sessionStorage.getItem("username");
    this.isCollapseStatus = sessionStorage.getItem("isCollapse") == 1 ? true : false;
    this.breadcrumbArr.unshift(this.$route.matched[this.$route.matched.length - 1].meta.pageTitle);
    this.breadcrumbArr.unshift(this.$route.matched[this.$route.matched.length - 2].meta.pageTitle);
  },
  mounted() {
    if (this.tableDetailStatusC == "halfScreen") {
      this.draw_direction = "btt";
      this.$emit("tableHeight", "up");
    }
    if (this.tableDetailStatusC == "halfScreenBan") {
      this.draw_direction = "rtl";
      this.$emit("tableHeight", "");
    }
    if (this.tableDetailStatusC == "halfAll") {
      this.draw_direction = "rtl";
      this.$emit("tableHeight", "");
    }
  },
  destroyed() {
    this.$EventBus.$off("isCollapse");
  },
};
</script>
<style scoped>
.batchSearchClass /deep/ .el-input__inner {
  border: none;
  height: 26px;
}

.taskDesk_box {
  width: 100%;
}

.taskDesk_box /deep/.isHightClass {
  background: #e8f4ff;
  color: #1682e6;
}

.taskDesk_box /deep/ .el-card__body {
  padding: 10px 10px 12px;
}

.taskDesk_box /deep/ .el-form-item--mini .el-form-item__label {
  width: 110px;
  min-width: 88px;
  white-space: nowrap;
  text-align: right;
  color: #aca9a9;
  font-size: 13px !important;
  font-weight: 600 !important;
}

.taskDesk_box /deep/ .el-card__body {
  padding: 10px;
}

.taskDesk_box .refreshBox_w {
  display: flex;
  position: relative;
  bottom: -5px;
  z-index: 1;
  justify-content: center;
  top: -7px;
}

.batchClass {
  margin: 20px 0 8px;
  overflow: auto;
}

.detBox {
  display: flex;
  line-height: 16px;
}

.detBox label {
  font-size: 14px;
  font-weight: 600;
  color: #777;
  min-width: 120px;
  display: inline-block;
  text-align: right;
  padding-right: 15px;
}

.detBox span {
  font-size: 12px;
  word-break: break-all;
}

.diaBox111 /deep/ .el-table th {
  padding: 0px !important;
}

.taskDesk_box /deep/ .infoTable .el-table__header {
  display: none;
}

.taskDesk_box /deep/ .infoTable td {
  border-bottom: none;
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
<style lang="scss">
.speical_tag .el-tag__close {
  position: absolute !important;
  top: 6px !important;
  right: 0 !important;
}

.el-tooltip__popper {
  white-space: pre-line !important;
}

.taskDesk_box .uls {
  padding: 0;
  margin: 0;
  overflow: auto;
}

.taskDesk_box .uls li {
  list-style: none;
  padding: 8px 3px;
  font-size: 12px;
  border-bottom: 1px solid #eee;
}

.taskDesk_box .el-table .colorStatus {
  background: #dcecf9 !important;
}

.taskDesk_box .uls li:hover {
  background: #faf9f8;
}

.taskDesk_box .batchH4 {
  position: relative;
  margin: 0;
  padding: 0;
  color: #333333;
  font-size: 14px;
  font-weight: 500;
  text-indent: 20px;
  line-height: 24px;
}

.taskDesk_box {
  .tabOneBox {
    padding: 20px 10px 20px 20px;
    font-size: 13px !important;

    span {
      font-size: 16px;
    }

    .el-link {
      span {
        font-size: 12px;
      }
    }

    .el-col:first-child {
      white-space: nowrap;
      text-align: right;
      color: #aca9a9;
      font-size: 13px !important;
      font-weight: 600 !important;
    }

    .el-row .el-col:last-child {
      color: #333;
      font-size: 12px;
      position: relative;
      bottom: -1px;
    }
  }

  .detailContent {
    padding: 20px 10px 0px 20px;
    flex: 1;
    margin-left: 15px;
    cursor: pointer;
  }

  .detailContent .el-row {
    padding: 15px 0;
    font-size: 12px;
    color: #999;
  }

  .detailContent .el-row .contents {
    color: #000;
    font-size: 13px;
    margin-top: 10px;
    word-wrap: break-word;
  }

  .wordred {
    color: red;
  }

  .titleCount {
    width: 100%;
    background: #eee;
    border-radius: 5px;
    padding: 8px 10px;
    text-align: center;
    font-size: 12px;
    cursor: pointer;
    margin: 15px 0;
  }

  .circle,
  .circleR,
  .circleS,
  .circleF,
  .circleI {
    padding: 0 5px;
    height: 50px;
    position: relative;
    border-radius: 10px;
    border: 1px solid #999;
    text-align: center;
    line-height: 50px;
    font-size: 12px;
    margin: 0 5px;
    cursor: pointer;
  }

  .circleR {
    background: #f00;
    color: #fff;
  }

  .circleS {
    background: #265d2c;
    color: #fff;
  }

  .circleF {
    background: red;
    color: #fff;
  }

  .circleI {
    background: #eee;
  }

  .lines {
    width: 33px;
    height: 1px;
    line-height: 40px;
    background: #333;
    vertical-align: middle;
    position: relative;
    top: 50%;
    margin-top: 0px;
    left: 3px;
  }

  .fuhao {
    font-size: 15px;
    position: relative;
    top: 50%;
    font-weight: 600;
    margin-top: -7.5px;
    cursor: pointer;
  }

  .shuline {
    height: 40px;
    width: 1px;
    background: #999;
    position: relative;
    bottom: 0;
    left: 50%;
    margin-left: -1px;
  }

  .shufu {
    font-size: 14px;
    position: relative;
    bottom: 18px;
  }

  .childBlock {
    width: 100%;
    margin-top: 66px;
    display: -webkit-box;
    display: -ms-flexbox;
    display: flex;
    -webkit-box-pack: center;
    -ms-flex-pack: center;
    justify-content: center;
  }

  .childBlockChild:last-child .lines {
    display: none;
  }

  .childBlockChild:last-child .fuhao:last-child {
    display: none;
  }
}

.informationIfrme {
  width: 100%;
  height: calc(100% + 25px);
  position: relative;
  top: -4px;
  left: 0;
}
</style>
<style>
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
  padding: 0 !important;
}

.el-breadcrumb__inner,
.el-breadcrumb__inner a {
  font-weight: 600 !important;
}

.el-breadcrumb__inner {
  color: #000;
  font-size: 14px;
}

.el-breadcrumb__separator {
  color: #000;
  font-weight: 600 !important;
  font-size: 14px;
}

.el-drawer__body {
  overflow: auto;
  border-left: 1px solid #e9e9e9;
  padding-left: 15px;
}

.el-drawer__header {
  display: none !important;
}

.trangle {
  position: relative;
  display: block;
  width: 0;
  height: 0;
  border-top: 30px solid #409eff;
  border-right: 36px solid transparent;
  cursor: pointer;
  right: 15px;
}

.widthwarp {
  display: inline-block;
  max-width: 250px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  height: 14px;
}

.advanced_searchDetial {
  float: left;
  margin-left: 15px;
  margin-top: var(--dynamic-margin-top) !important;
  position: relative;
  margin-left: 163px;
}

.advanced_searchDetial .el-dialog__header {
  padding: 0 !important;
}

.advanced_searchDetial .el-dialog__body {
  padding: 10px !important;
}

.tagTypeClass {
  font-size: 11px;
  color: #666;
}

.speical_tag {
  color: #409eff;
  font-weight: 600;
  font-size: 12px;
  position: relative;
  top: 4px;
}

.gjSearchBody .gjSearchBodyItem {
  padding: 5px !important;
  background: #f8f9fa !important;
}

.advanced_searchDetial /deep/ .el-dialog__body {
  padding: 10px 30px 20px !important;
  background: #ebeff3 !important;
}

.iscollClass2 {
  width: calc(100vw - 140px);
  left: 140px;
  overflow: auto;
  top: 46px;
}

.iscollClass1 {
  width: calc(100vw - 60px);
  left: 60px;
  overflow: auto;
}

.iscollClass3 {
  width: calc(100vw - 590px);
  left: 590px;
  overflow: auto;
  top: 0;
}

.iscollClass4 {
  width: calc(100vw - 510px);
  left: 520px;
  overflow: auto;
  top: 48px;
}
</style>