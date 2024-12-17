<template>
  <div class="taskDesk_box">
    <el-drawer :visible.sync="drawerDetail" :direction="drawDirection" :modal="false" :size="drawSize"
      :wrapperClosable="false" :modal-append-to-body="false"
      :style="isCollapseStatus ? 'width: calc(100vw - 60px);left: 60px;overflow: auto;' : 'width: calc(100vw - 140px);left: 140px;overflow: auto;'"
      custom-class="spiceDrawLog" :before-close="handleClose">
      <div style="height: 15px; display: none">
        <span class="trangle" @click="closeThis">
          <i class="el-icon-close"
            style="color: rgb(255, 255, 255);  position: absolute;  top: -28px;  right: -19px;"></i>
        </span>
      </div>
      <div
        style="height: 100%;  cursor: pointer;  width: 15px;  background: rgb(247, 247, 247);  position: absolute;  z-index: 999;  display: flex;  justify-content: center;  align-items: center;  left: 0;  top: 0px;"
        @click="closeThis">
        <i class="el-icon-arrow-right" style="color: #409eff"> </i>
      </div>
      <div style="display: flex; position: relative">
        <el-card class="tabOneBox">
          <span>任务信息</span>
          <el-row style="margin-top: 10px">
            任务ID：{{ detaileMessage.job_id }}
          </el-row>
          <el-row style="margin-top:10px">
              任务表单：<el-link type="primary" v-if="Object.keys(historyTableList).length" @click="addOpen">查看</el-link>
          </el-row>
        </el-card>
        <el-card class="detailContent">
          <div style="display: flex; justify-content: space-between">
            <span>任务状态</span>
          </div>
          <el-row :gutter="20">
            <el-col :span="5">
              任务类型
              <div class="contents">{{ detaileMessage.job_type }}</div>
            </el-col>
            <el-col :span="5">
              任务状态
              <div class="contents">
                <span v-if="detaileMessage.job_status == '进行中'">
                  <el-link :underline="false" type="warning" class="el-icon-time" style="font-size: 12px"><span
                      style="margin-left: 5px">{{ detaileMessage.job_status }}</span></el-link>
                </span>
                <span v-if="detaileMessage.job_status == '成功' || detaileMessage.job_status == '结单'">
                  <el-link :underline="false" type="success" class="el-icon-success" style="font-size: 12px"><span
                      style="margin-left: 5px">{{ detaileMessage.job_status }}</span></el-link>
                </span>
                <span
                  v-if="detaileMessage.job_status == '等待' || detaileMessage.job_status == '审核中' || detaileMessage.job_status == '创建中'">
                  <el-link :underline="false" type="info" class="el-icon-remove" style="font-size: 12px"><span
                      style="margin-left: 5px">{{ detaileMessage.job_status }}</span></el-link>
                </span>
                <span v-if="detaileMessage.job_status == '失败' || detaileMessage.job_status == '已拒绝'">
                  <el-link :underline="false" type="danger" class="el-icon-error" style="font-size: 12px"><span
                      style="margin-left: 5px">{{ detaileMessage.job_status }}</span></el-link>
                </span>
                <span v-if="detaileMessage.job_status == '待确认'">
                  <el-link :underline="false" type="warning" class="el-icon-question" style="font-size: 12px"><span
                      style="margin-left: 5px">{{ detaileMessage.job_status }}</span></el-link>
                </span>
                <span v-if="detaileMessage.job_status == '关闭'">
                  <el-link :underline="false" type="info" class="el-icon-remove" style="font-size: 12px"><span
                      style="margin-left: 5px">{{ detaileMessage.job_status }}</span></el-link>
                </span>
              </div>
            </el-col>
            <el-col :span="5">
              <div>执行情况</div>
              <div style="margin-top: 8px; line-height: 15px">
                <span @click="clickTotal('')" style="display: block; font-size: 12px">总量：<span
                    style="margin: 0 3px">{{ detaileMessage.total ? detaileMessage.total : 0 }}</span></span>
                <span @click="clickTotal('success')" style="display: block; color: #265d2c; font-size: 12px">成功：<span
                    style="margin: 0 3px">{{ detaileMessage.success ? detaileMessage.success : 0 }}</span></span>
                <span @click="clickTotal('failure')" style="display: block; color: red; font-size: 12px">失败：<span
                    style="margin: 0 3px">{{ detaileMessage.failure ? detaileMessage.failure : 0 }}</span></span>
                <span @click="clickTotal('idle')" style="display: block; color: #4169e1; font-size: 12px">待执行：<span
                    style="margin: 0 3px">{{ detaileMessage.idle ? detaileMessage.idle : 0 }}</span></span>
                <span @click="clickTotal('confirming')"
                  style="display: block; color: #f4a460; font-size: 12px">待确认：<span
                    style="margin: 0 3px">{{ detaileMessage.confirming ? detaileMessage.confirming : 0 }}</span></span>
              </div>
            </el-col>
            <el-col :span="4">
              申请人
              <div class="contents">{{ detaileMessage.user }}</div>
            </el-col>
            <el-col :span="5">
              申请时间
              <div class="contents">{{ detaileMessage.create_time }}</div>
            </el-col>
          </el-row>
          <i class="el-icon-circle-close"
            style="color: red;  position: absolute;  top: 16px;  right: 14px;  cursor: pointer;  font-size: 17px;"
            @click="closeThis"></i>
        </el-card>
      </div>
      <div style="display: flex; justify-content: space-between; margin-top: 15px">
        <div>
        </div>
        <div>
          <!-- <el-select style="margin: 0 8px; width: 130px" v-model="sreachForm.Object_status" filterable clearable
            @clear="getObjectList" @change="getObjectList" placeholder="对象状态">
            <template v-for="(item,key,index) in organizationArr">
                <el-option :key="index" v-if="item!='审核中'&&item!='创建中'"  :label="item" :value="key">
                </el-option>
            </template>
          </el-select>
          <el-input style="width: 200px" clearable @clear="getObjectList" v-model="sreachForm.Object_code"
            placeholder="请输入对象编号" @keydown.enter.native="getObjectList"></el-input> -->
          <i @click="refreshBox" class="el-icon-refresh"
            style="cursor: pointer;  font-size: 20px;  vertical-align: middle;  color: #999;  margin: 0 8px;"></i>
        </div>
      </div>
      <div style="margin-top: 20px" v-if="srechTabShow">
        <el-card style="cursor: pointer">
          <h4 class="batchH4" @click="searchStatus">
            搜索结果:{{ searchNum }}台
            <span style="float: right; cursor: pointer">
              <i class="el-icon-arrow-down" v-if="srechTabStatus"></i>
              <i class="el-icon-arrow-up" v-if="!srechTabStatus"></i>
            </span>
          </h4>
        </el-card>
        <el-table @selection-change="searchSelectiontryAgain" :data="searchDataObject" style="margin-top: 8px"
          v-if="srechTabStatus">
          <el-table-column type="selection" width="40"></el-table-column>
          <el-table-column prop="code" label="对象编号" align="center" width="180">
            <template slot-scope="scoped">
              <div style="display: flex;  justify-content: center;  align-items: center;">
                <span style="font-size: 12px">
                  {{ scoped.row.code }}
                </span>
                <el-tooltip v-if="scoped.row.errors" class="item" placement="top" effect="light">
                  <div slot="content">
                    <div v-for="(item, key, index) in scoped.row.errors" :key="index">
                      {{ item }}
                    </div>
                  </div>
                  <i v-if="scoped.row.status == '成功'" style="color: #b0e0e6; font-size: 16px; margin-left: 5px"
                    class="el-icon-warning"></i>
                  <i v-if="scoped.row.status == '失败'" style="color: #e60000; font-size: 16px; margin-left: 5px"
                    class="el-icon-warning"></i>
                  <i v-if="scoped.row.status == '进行中'" style="color: #ffcc00; font-size: 16px; margin-left: 5px"
                    class="el-icon-warning"></i>
                  <i v-if="scoped.row.status == '等待'" style="color: #b0e0e6; font-size: 16px; margin-left: 5px"
                    class="el-icon-warning"></i>
                </el-tooltip>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="对象名称" align="center" prop="name"></el-table-column>
          <el-table-column label="开始时间" align="center" prop="start_time">
            <template slot-scope="scoped">
              <span v-if="scoped.row.start_time.indexOf('0001') == 0"> - </span>
              <span v-else>
                {{ scoped.row.start_time }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="end_time" label="结束时间" align="center">
            <template slot-scope="scoped">
              <span v-if="scoped.row.end_time.indexOf('0001') == 0">
                -
              </span>
              <span v-else>
                {{ scoped.row.end_time }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="procedure_code" label="执行动作集" align="center">
          </el-table-column>
          <el-table-column label="对象状态" align="left" width="100">
            <template slot-scope="scoped">
              <span v-if="scoped.row.status == '进行中'">
                <span
                  style="background: #f7ab40;  color: #fff;  padding: 4px;  border-radius: 4px;  vertical-align: middle;  font-size: 12px;">{{
                    scoped.row.status }}</span>
                <i class="el-icon-loading" style="font-size: 17px;  margin-left: 8px;  vertical-align: middle;"></i>
              </span>
              <span style="background: #265d2c;color: #fff;padding: 4px;border-radius: 4px;font-size: 12px;"
                v-if="scoped.row.status == '成功'">{{ scoped.row.status }}</span>
              <span style="background: #666;color: #fff;padding: 4px;border-radius: 4px;font-size: 12px;"
                v-if="scoped.row.status == '等待' || scoped.row.status == '审核中' || scoped.row.status == '创建中'">{{
                  scoped.row.status }}</span>
              <span style="background: red;  color: #fff;  padding: 4px;  border-radius: 4px;  font-size: 12px;"
                v-if="scoped.row.status == '失败' || scoped.row.status == '已拒绝'">{{ scoped.row.status }}</span>
              <span style="background: #f4a460;  color: #fff;  padding: 4px;  border-radius: 4px;  font-size: 12px;"
                v-if="scoped.row.status != '确认中' && scoped.row.status != '成功' && scoped.row.status != '等待' && scoped.row.status != '审核中' && scoped.row.status != '失败' && scoped.row.status != '进行中' && scoped.row.status != '创建中'">{{
                  scoped.row.status }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="address" label="操作" width="200" align="center"></el-table-column>
        </el-table>
        <div style="margin-top: 10px; text-align: center" id="tablePagina">
          <el-pagination @current-change="changehandleCurrentS" @size-change="handleSizeChangeS"
            :current-page="SpageNum" :page-size="SpageSize" :page-sizes="[20, 50, 100, 150, 200]"
            layout="total, sizes, prev, pager, next, jumper" :total="searchNum">
          </el-pagination>
        </div>
      </div>
      <template v-else>
        <div class="batchClass" v-for="(item, index) in totaltableData" :key="index">
          <el-card style="cursor: pointer">
            <h4 class="batchH4">
              <span style="float: right; cursor: pointer; margin-top: 24px">
                <el-tooltip class="item" effect="light" content="开始" placement="top-end">
                  <el-button circle type="primary" icon="el-icon-arrow-right" style="vertical-align: middle"
                    v-if="item.status == '已暂停' && detaileMessage.job_status != '审核中' && detaileMessage.job_status != '创建中'"
                    @click="lockClick(item.id)"></el-button>
                </el-tooltip>
                <el-tooltip class="item" effect="light" content="跳过" placement="top-end">
                  <el-button circle type="primary" icon="el-icon-d-arrow-right" style="vertical-align: middle"
                    v-if="item.status == '已暂停' && detaileMessage.job_status != '审核中' && detaileMessage.job_status != '创建中'"
                    @click="JumpNewClick(item.id)"></el-button>
                </el-tooltip>
                <i class="el-icon-arrow-down" style="font-size: 24px; vertical-align: middle"
                  @click="zsClick(item.id, index)" v-if="item.zkStatus"></i>
                <i class="el-icon-arrow-up" style="font-size: 24px; vertical-align: middle"
                  @click="zsClick(item.id, index)" v-if="!item.zkStatus"></i>
              </span>
              <span>{{ item.title }}</span>
              <div style="font-size: 12px; color: #999" class="stepStatus">
                批次状态：
                <span v-if="item.status == '进行中'">
                  <el-link :underline="false" type="warning" class="el-icon-time" style="font-size: 12px"><span
                      style="margin-left: 5px">{{ item.status }}</span></el-link>
                </span>
                <span v-if="item.status == '成功'">
                  <el-link :underline="false" type="success" class="el-icon-success" style="font-size: 12px"><span
                      style="margin-left: 5px">{{ item.status }}</span></el-link>
                </span>
                <span
                  v-if="item.status == '等待' || item.status == '审核中' || item.status == '已暂停' || item.status == '创建中'">
                  <el-link :underline="false" type="info" class="el-icon-remove" style="font-size: 12px"><span
                      style="margin-left: 5px">{{ item.status }}</span></el-link>
                </span>
                <span v-if="item.status == '失败' || item.status == '已拒绝'">
                  <el-link :underline="false" type="danger" class="el-icon-error" style="font-size: 12px"><span
                      style="margin-left: 5px">{{ item.status }}</span></el-link>
                </span>
                <span v-if="item.status == '待确认'">
                  <el-link :underline="false" type="warning" class="el-icon-question" style="font-size: 12px"><span
                      style="margin-left: 5px">{{ item.status }}</span></el-link>
                </span>
                <span v-if="item.status == '关闭'">
                  <el-link :underline="false" type="info" class="el-icon-remove" style="font-size: 12px"><span
                      style="margin-left: 5px">{{
                        item.status
                      }}</span></el-link>
                </span>
              </div>
              <div style="font-size: 12px; color: #999">
                执行动作集：{{ item.procedure_code }}，对象数量：{{ item.total ? item.total : 0 }}
              </div>
            </h4>
          </el-card>
          <div v-if="item.zkStatus && item.objects">
            <el-table :data="item.objects" style="margin-top: 8px"
              @selection-change="handleSelectiontryAgain($event, index, item)">
              <el-table-column type="selection" width="40"></el-table-column>
              <el-table-column prop="code" label="对象编号" align="center" width="180">
                <template slot-scope="scoped">
                  <div style="display: flex;  justify-content: center;  align-items: center;">
                    <span style="font-size: 12px">
                      {{ scoped.row.code }}
                    </span>
                    <el-tooltip v-if="scoped.row.errors" class="item" placement="top" effect="light">
                      <div slot="content">
                        <div v-for="(item, key, index) in scoped.row.errors" :key="index">
                          {{ item }}
                        </div>
                      </div>
                      <i v-if="scoped.row.status == '成功'" style="color: #b0e0e6;  font-size: 16px;  margin-left: 5px;"
                        class="el-icon-warning"></i>
                      <i v-if="scoped.row.status == '失败'" style="color: #e60000;  font-size: 16px;  margin-left: 5px;"
                        class="el-icon-warning"></i>
                      <i v-if="scoped.row.status == '进行中'" style="color: #ffcc00;  font-size: 16px;  margin-left: 5px;"
                        class="el-icon-warning"></i>
                      <i v-if="scoped.row.status == '等待'" style="color: #b0e0e6;  font-size: 16px;  margin-left: 5px;"
                        class="el-icon-warning"></i>
                    </el-tooltip>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="对象名称" align="center" prop="name"></el-table-column>
              <el-table-column label="开始时间" align="center" prop="start_time">
                <template slot-scope="scoped">
                  <span v-if="scoped.row.start_time.indexOf('0001') == 0">
                    -
                  </span>
                  <span v-else>
                    {{ scoped.row.start_time }}
                  </span>
                </template>
              </el-table-column>
              <el-table-column prop="end_time" label="结束时间" align="center">
                <template slot-scope="scoped">
                  <span v-if="scoped.row.end_time.indexOf('0001') == 0">
                    -
                  </span>
                  <span v-else>
                    {{ scoped.row.end_time }}
                  </span>
                </template>
              </el-table-column>
              <el-table-column label="对象状态" align="left" width="100">
                <template slot-scope="scoped">
                  <span v-if="scoped.row.status == '进行中'">
                    <el-link :underline="false" type="warning" class="el-icon-time" style="font-size: 12px"><span
                        style="margin-left: 5px">{{ scoped.row.status }}</span>
                    </el-link>
                  </span>
                  <span v-if="scoped.row.status == '成功'">
                    <el-link :underline="false" type="success" class="el-icon-success" style="font-size: 12px"><span
                        style="margin-left: 5px">{{ scoped.row.status }}</span></el-link>
                  </span>
                  <span v-if="scoped.row.status == '等待' || scoped.row.status == '审核中' || scoped.row.status == '创建中'">
                    <el-link :underline="false" type="info" class="el-icon-remove" style="font-size: 12px"><span
                        style="margin-left: 5px">{{ scoped.row.status }}</span></el-link>
                  </span>
                  <span v-if="scoped.row.status == '失败' || scoped.row.status == '已拒绝'">
                    <el-link :underline="false" type="danger" class="el-icon-error" style="font-size: 12px"><span
                        style="margin-left: 5px">{{ scoped.row.status }}</span></el-link>
                  </span>
                  <span v-if="scoped.row.status == '待确认'">
                    <el-link :underline="false" type="warning" class="el-icon-question" style="font-size: 12px"><span
                        style="margin-left: 5px">{{ scoped.row.status }}</span></el-link>
                  </span>
                  <span v-if="scoped.row.status == '关闭'">
                    <el-link :underline="false" type="info" class="el-icon-remove" style="font-size: 12px"><span
                        style="margin-left: 5px">{{ scoped.row.status }}</span></el-link>
                  </span>
                  <span
                    v-if="scoped.row.status != '待确认' && scoped.row.status != '成功' && scoped.row.status != '等待' && scoped.row.status != '审核中' && scoped.row.status != '失败' && scoped.row.status != '进行中' && scoped.row.status != '关闭' && scoped.row.status != '创建中'">
                    <el-link :underline="false" type="info" class="el-icon-question" style="font-size: 12px"><span
                        style="margin-left: 5px">{{ scoped.row.status }}</span></el-link>
                  </span>
                </template>
              </el-table-column>
              <el-table-column prop="address" label="操作" width="200" align="center"></el-table-column>
            </el-table>
            <div style="margin-top: 10px; text-align: center" id="tablePagina">
              <el-pagination @current-change="changehandleCurrent($event, index)"
                @size-change="handleSizeChange($event, index)" :current-page="tableDatapage[index].stepPageNum"
                :page-size="tableDatapage[index].stepPageSize" :page-sizes="[20, 50, 100, 150, 200]"
                layout="total, sizes, prev, pager, next, jumper" :total="item.total">
              </el-pagination>
            </div>
          </div>
        </div>
      </template>
    </el-drawer>
    <!-- 回显任务弹框 infoAddAttr 表格行信息 infoAddForm表单提交信息 reviewStatus弹框状态-->
    <div v-if="reviewStatus">
      <formtem-Dialog :baseConfigData="baseConfigData" :nodeCheckName="nodeCheckName" :detailOrLook="detailOrLook"
        @addClose="addClose" :infoAddAttr="infoAddAttr" :detailOrsigle="detailOrsigle" :infoAddForm="infoAddForm"
        :reviewStatus="reviewStatus" :newTaskName="newTaskName" :historyTableList="historyTableList"
        :historyparamreferenceShow="historyparamreferenceShow" :childNodename="childNodename" :checkList="checkList"
        :jobIdStatus="detaileMessage.job_status" :FormTemplatjobMode="FormTemplatjobMode"></formtem-Dialog>
    </div>
  </div>
</template>

<script>
import formtemDialog from "@/views/components/formTemp";
import { userMixin } from "@/components/mixin/user";
import { mapGetters } from "vuex";
import Http from '@/components/api/services';
//例如：import 《组件名称》 from '《组件路径》';
export default {
  mixins: [userMixin],
  //import引入的组件需要注入到对象中才能使用
  props: {
    jobIds: {
      type: Number,
    },
  },
  components: {
    formtemDialog,
  },
  filters: {},
  computed: {
    ...mapGetters([
      "sidebar",
      "loginUserName",
      "formPolicyDomain",
    ]),
  },
  data() {
    //这里存放数据
    return {
      tableDatapage: [{
        stepPageNum: 1,
        stepPageSize: 20,
      },],
      StepNum: 0,
      SpageNum: 1,
      SpageSize: 20,
      dealerArr: [],
      FormTemplatjobMode: "",
      FormTemplatworkflow: "",
      usernameAll: "",
      baseConfigData: {},
      historyparamreferenceShow: [],
      historyTableList: {},
      newTaskName: "",
      nodeCheckName: "",
      detailOrLook: true,
      detailOrsigle: "",
      reviewStatus: false,
      formId: "",
      infoAddForm: {},
      infoAddAttr: [],
      drawObjectDetail: false,
      searchJumpArr: [],
      jumpArr: [],
      jumorRetrypArr: [],
      batchJumpTo: true,
      batchJumpToRey: true,
      indexBtchJump: 0,
      srechTabStatus: true,
      srechTabShow: false,
      sreachForm: {
        object_code_in:[],
        Object_status: "",
        Object_code: "",
      },
      sys_bs_in: "",
      timeValueArrObject: [],
      timeValueArr: [],
      totaltableData: [],
      exportsExcelData: [],
      searchDataObject: [],
      searchNum: 0,
      organizationArr: {
        "success": "成功",
        "failure": "失败",
        "running": "进行中",
        "auth_wait": "审核中",
        "auth_req": "审核中",
        "auth_wait_sync": "审核中",
        "auth": "审核中",
        "idle": "等待",
        "confirming": "待确认",
        "auth_deny": "已拒绝",
        "create_wait": "创建中",
        "creating": "创建中",
      },
      detaileMessage: {},
      batchIdArr: [],
      batchIdArrSelect: [],
      batchIdArrSelectObj: {},
      batchSearchSelect: [],
      object_fields: {},
      drawerDetail: true,
      drawSize: "100%",
      drawDirection: "rtl",
      isCollapseStatus: false,
      childNodename: "",
      checkList: [],
    };
  },
  watch: {
    jobIds: {
      handler: function () {
        this.drawObjectDetail = false;
        this.getJobList(this.StepNum, false, true);
      },
      deep: true,
    },
  },
  //方法集合
  methods: {
    handleSizeChange(val, index) {
      this.tableDatapage[index].stepPageSize = val;
      this.tableDatapage[index].stepPageNum = 1;
      this.StepNum = index;
      this.getJobList(index);
    },
    changehandleCurrent(val, index) {
      this.tableDatapage[index].stepPageNum = val;
      this.StepNum = index;
      this.getJobList(index);
    },
    handleSizeChangeS(val) {
      this.SpageSize = val;
      this.SpageNum = 1;
      this.getObjectList();
    },
    changehandleCurrentS(val) {
      this.SpageNum = val;
      this.getObjectList();
    },
    // 点击总数联动对象状态
    clickTotal(data) {
      this.sreachForm.Object_status = data;
      this.getObjectList();
    },
    // 已暂停按钮点击
    lockClick(id) {
      this.$confirm(`是否确定开始？`, {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).then(() => {
        let getData = {
          id: id,
          op_code: "idle",
          user: this.usernameAll,
        };
        this.Http("root", `/api/v1/mars/step/operate`, "post", getData).then((response) => {
          if (response.data.errors) {
            this.$message({
              showClose: true,
              message: "数据请求失败,请联系管理员!\n" +
                response.data.message,
              type: "error",
            });
            return;
          }
          this.$message({
            message: "操作成功",
            type: "success",
            showClose: true,
          });
        })
          .catch((err) => {
            console.log(err);
          });
      });
    },
    // 暂停状态下跳过按钮
    JumpNewClick(id) {
      this.$confirm(`是否确定跳过？`, {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }).then(() => {
        let getData = {
          id: id,
          op_code: "success",
          user: this.usernameAll,
        };
        this.Http("root", `/api/v1/mars/step/operate`, "post", getData).then((response) => {
          if (response.data.errors) {
            this.$message({
              showClose: true,
              message: "数据请求失败,请联系管理员!\n" +
                response.data.message,
              type: "error",
            });
            return;
          }
          this.$message({
            message: "操作成功",
            type: "success",
            showClose: true,
          });
        })
          .catch((err) => {
            console.log(err);
          });
      });
    },
    //点击左边获取任务详情
    getJobList(index, type, temStatus) {
      let searchFormData = {};
      if (index) {
        searchFormData = {
          id: this.jobIds,
          page_num: this.tableDatapage[index].stepPageNum ? this.tableDatapage[index].stepPageNum : 1,
          page_size: this.tableDatapage[index].stepPageSize ? this.tableDatapage[index].stepPageSize : 20,
          step_num: this.StepNum + 1,
        };
      } else {
        searchFormData = {
          Id: this.jobIds,
          page_num: this.tableDatapage[0].stepPageNum,
          page_size: this.tableDatapage[0].stepPageSize,
          step_num: this.StepNum + 1,
        };
      }
      this.exportsExcelData = [];
      this.object_fields = {};
      this.jumpArr = [];
      this.jumorRetrypArr = [];
      this.batchIdArr = [];
      this.searchJumpArr = [];
      this.batchJumpToRey = true;
      this.batchJumpTo = true;
      this.batchIdArrSelect = [];
      this.batchIdArrSelectObj = [];
      this.totaltableData = [];
      Http.getJobDetail(searchFormData).then((response) => {
        this.detaileMessage = response.data.data;
        var tableData = response.data.data.steps ? JSON.parse(JSON.stringify(response.data.data.steps)) : [];
        tableData.map((item, index) => {
          this.totaltableData[index] = item;
          if (type) {
            this.tableDatapage.push({
              stepPageNum: 1,
              stepPageSize: 20,
            });
          } else if (index == this.StepNum) {
            this.tableDatapage[index]["stepPageNum"] = this.tableDatapage[index].stepPageNum ? this.tableDatapage[index].stepPageNum : 1;
            this.tableDatapage[index]["stepPageSize"] = this.tableDatapage[index].stepPageSize ? this.tableDatapage[index].stepPageSize : 20;
          }
          if (this.StepNum == index) {
            this.totaltableData[index].zkStatus = true;
          } else {
            this.totaltableData[index].zkStatus = false;
          }
          this.searchJumpArr = item.jump_options;
          this.object_fields = item.object_fields;
        });
      })
      .catch((err) => {
        this.totaltableData = [];
        console.log(err);
      });
    },
    //点击左边获取对象列表
    getObjectList() {
      if (this.sreachForm.Object_status == "" && this.sreachForm.Object_code == "" && this.sys_bs_in == "") {
        this.batchSearchSelect = [];
        this.batchIdArrSelect = [];
        this.handleChecbox();
        this.$nextTick(() => {
          this.srechTabShow = false;
        });
      } else {
        this.batchIdArrSelect = [];
        this.batchSearchSelect = [];
        this.handleChecbox();
        let bsNewArr = [];
        let bsbatchNewArr = [];
        if (!this.sys_bs_in) {
          bsNewArr = [];
        } else {
          bsNewArr = this.sys_bs_in.split("\n");
          bsNewArr.map((k, i) => {
            if (k.replace(/\s/g, "")) {
              bsbatchNewArr.push(k.replace(/\s/g, ""));
            }
          });
        }
        this.sreachForm.object_code_in = bsbatchNewArr ? bsbatchNewArr : [];
        let searchFormDatas = {
          job_id: this.jobIds,
          ...this.sreachForm,
          page_num: this.SpageNum,
          page_size: this.SpageSize,
        };
        this.searchDataObject = [];
        this.Http("root", `/api/v1/mars/object/list`, "post", searchFormDatas).then((response) => {
          if (response.data.code == 0) {
            this.searchDataObject = response.data.data.objects;
            this.searchNum = response.data.data.count ? response.data.data.count : "0";
          } else {
            this.$message({
              showClose: true,
              message: "搜索失败",
              type: "error",
            });
          }
          this.srechTabShow = true;
        })
          .catch((err) => {
            this.$message({
              showClose: true,
              message: "搜索失败",
              type: "error",
            });
            this.srechTabShow = true;
          });
      }
    },
    // 批量搜索bs重置
    bscodeBatch() {
      this.sys_bs_in = "";
      this.getObjectList();
    },
    // 刷新按钮
    refreshBox() {
      this.$nextTick(() => {
        this.getJobList(this.StepNum, false, true);
        this.srechTabShow = false;
        this.sreachForm = {
          Object_status: "",
          Object_code: "",
          sreachForm: "",
        };
      });
    },
    // 搜索展开或者收起
    searchStatus() {
      this.srechTabStatus = !this.srechTabStatus;
      if (this.batchSearchSelect.length != 0) {
        this.batchSearchSelect = [];
        this.handleChecbox();
      }
    },
    // 展开或者收起
    zsClick(code, index) {
      this.StepNum = index;
      if (
        this.totaltableData[index] &&
        this.totaltableData[index].zkStatus === true
      ) {
        let arr = JSON.parse(JSON.stringify(this.totaltableData));
        this.$nextTick(() => {
          arr.map((item) => {
            if (item.id == code && item.zkStatus) {
              item.zkStatus = false;
            }
          });
        });
        this.totaltableData = arr;
        return;
      }
      this.getJobList(index);
      if (this.batchIdArrSelect && this.batchIdArrSelect.length != 0) {
        this.batchIdArrSelect[index] = [];
        this.handleChecbox();
      }
    },
    // 搜索表格复选
    searchSelectiontryAgain(row) {
      this.batchSearchSelect = row;
      let selectArr = [];
      let procedure_lastArrsel = [];
      this.indexBtchJump = "";
      this.batchIdArr = [];
      let goToHttpSel = true;
      if (this.batchSearchSelect.length != 0) {
        this.batchSearchSelect.map((item) => {
          if (procedure_lastArrsel.indexOf(item.procedure_last) == -1) {
            //是否是同一个节点下
            procedure_lastArrsel.push(item.procedure_last);
          }
          selectArr.push(item);
        });
      }
      // 如果是同一批次下且搜索数据不为空
      if (this.batchIdArrSelect.length == 0) {
        // 同一节点下
        if (procedure_lastArrsel.length == 1) {
          this.indexBtchJump = this.batchSearchSelect[0].id;
          this.batchSearchSelect.map((aitem) => {
            //同一批次节点判断数据是否是失败状态 如果不是失败状态则return
            if (aitem.status == "失败" || aitem.status == "待确认") {
              this.batchJumpToRey = false;
              this.batchJumpTo = false;
              goToHttpSel = true;
              return;
            } else if (aitem.status == "进行中") {
              this.batchJumpToRey = true;
              this.batchJumpTo = false;
              goToHttpSel = true;
              return;
            } else {
              this.batchJumpToRey = true;
              this.batchJumpTo = true;
              //没有失败节点则状态为true
              goToHttpSel = false;
            }
          });
          if (goToHttpSel) {
            let getData = {
              ObjectId: this.indexBtchJump
            };
            this.Http("root", `/api/v1/mars/object/get`, "get", getData).then((response) => {
              if (response.data.code == 0) {
                if (response.data.data) {
                  this.jumorRetrypArr = response.data.data.fail_options;
                  this.jumpArr = response.data.data.jump_options;
                  if (
                    this.jumorRetrypArr &&
                    this.jumorRetrypArr.length != 0
                  ) {
                    this.batchJumpToRey = this.batchJumpToRey;
                  }
                  if (this.jumpArr.length != 0) {
                    this.batchJumpTo = this.batchJumpTo;
                  }
                }
              }
            })
              .catch((err) => {
                console.log(err);
              });
          } else {
            this.batchJumpTo = true;
            this.batchJumpToRey = true;
          }
        } else {
          this.batchJumpTo = true;
          this.batchJumpToRey = true;
        }
      } else {
        this.batchJumpTo = true;
        this.batchJumpToRey = true;
      }
      selectArr.map((arrItem) => {
        if (this.batchIdArr.indexOf(arrItem.id) == -1) {
          this.batchIdArr.push(arrItem.id);
        }
      });
    },
    // 表格复选框
    handleSelectiontryAgain(sel, index, item) {
      this.batchIdArrSelect[index] = sel;
      if (sel.length) {
        this.batchIdArrSelectObj["ASD" + index] = 1;
      } else {
        delete this.batchIdArrSelectObj["ASD" + index];
      }
      this.exportsExcelData = sel;
      this.handleChecbox();
    },
    // 整合复选框数据
    handleChecbox() {
      this.indexBtchJump = "";
      let arr = [];
      this.batchIdArr = [];
      let goToHttp = false;
      let procedure_lastArr = [];
      if (Object.keys(this.batchIdArrSelectObj).length > 1) {
        this.$message({
          showClose: true,
          message: "请选择同一批次服务器",
          type: "warning",
        });
        return false;
      } else {
        if (this.batchIdArrSelect.length != 0) {
          this.batchIdArrSelect.map((item, index) => {
            item.map((Citem) => {
              if (procedure_lastArr.indexOf(Citem.procedure_last) == -1) {
                //是否是同一个节点下
                procedure_lastArr.push(Citem.procedure_last);
              }
            });
            item.map((childItem) => {
              this.indexBtchJump = item[0].id;
              arr.push(childItem);
            });
          });
        }
      }
      // 如果是同一批次下且搜索数据不为空
      if (this.batchSearchSelect.length == 0) {
        // 同一节点下
        this.batchIdArrSelect.map((aitem) => {
          aitem.map((bitem) => {
            //同一批次节点判断数据是否是失败状态 如果不是失败状态则return
            if (bitem.status == "失败" || bitem.status == "待确认") {
              if (procedure_lastArr.length == 1) {
                this.batchJumpToRey = false;
              }
              this.batchJumpTo = false;
              goToHttp = true;
              return;
            } else if (bitem.status == "进行中") {
              this.batchJumpToRey = true;
              this.batchJumpTo = false;
              goToHttp = true;
              return;
            } else {
              this.batchJumpToRey = true;
              this.batchJumpTo = true;
              //没有失败节点则状态为true
              goToHttp = false;
            }
          });
        });
        if (goToHttp) {
          let getData = {
            ObjectId: this.indexBtchJump
          };
          this.Http("root", `/api/v1/mars/object/get`, "get", getData).then((response) => {
            if (response.data.code == 0) {
              if (response.data.data) {
                this.jumorRetrypArr = response.data.data.fail_options;
                this.jumpArr = response.data.data.jump_options;
                if (this.jumorRetrypArr && this.jumorRetrypArr.length != 0) {
                  this.batchJumpToRey = this.batchJumpToRey;
                }
                if (this.jumpArr.length != 0) {
                  this.batchJumpTo = this.batchJumpTo;
                }
              }
            }
          })
            .catch((err) => {
              console.log(err);
            });
        } else {
          this.batchJumpTo = true;
          this.batchJumpToRey = true;
        }
      } else {
        this.batchJumpTo = true;
        this.batchJumpToRey = true;
      }
      arr.map((arrItem) => {
        if (this.batchIdArr.indexOf(arrItem.id) == -1) {
          this.batchIdArr.push(arrItem.id);
        }
      });
    },
    // 发起任务弹框关闭
    addClose() {
      this.reviewStatus = false;
    },
    // 关闭详情
    closeThis() {
      this.$router.replace({
        path: location.pathname,
      });
      this.drawerDetail = false;
      this.$emit("addCloseDraw");
    },
    handleClose() {
      this.$router.replace({
        path: location.pathname,
      });
      this.drawerDetail = false;
      this.$emit("addCloseDraw");
    },
  },
  //生命周期 - 挂载完成（可以访问DOM元素）
  created() {
    this.pageSize = Number(this.$route.query._pageSize) ? Number(this.$route.query._pageSize) : this.pageSize;
    this.pageNum = Number(this.$route.query._pageNumber) ? Number(this.$route.query._pageNumber) : this.pageNum;
  },
  mounted() {
    this.usernameAll = sessionStorage.getItem("username");
    let startTime = new Date().getTime() - 24 * 60 * 60 * 1000 * 30;
    let endTime = new Date().getTime();
    this.timeValueArrObject = [
      this.$moment(startTime).format("YYYY-MM-DD HH:mm:ss"),
      this.$moment(endTime).format("YYYY-MM-DD HH:mm:ss"),
    ];
    this.timeValueArr = [
      this.$moment(startTime).format("YYYY-MM-DD HH:mm:ss"),
      this.$moment(endTime).format("YYYY-MM-DD HH:mm:ss"),
    ];
    //查看对象列表
    this.getJobList(false, true, true);
    this.$EventBus.$on("isCollapse", (value) => {
      this.isCollapseStatus = value ? true : false;
    });
  },
  created() {
    this.isCollapseStatus = sessionStorage.getItem("isCollapse") == 1 ? true : false;
  },
};
</script>
<style scoped>
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

.batchClass .stepStatus span {
  text-indent: 0px;
}

.batchClass /deep/ .el-link {
  display: inline-block;
}

.batchClass /deep/ .el-button {
  padding: 9px !important;
}
</style>
<style lang="scss">
.demoBatch_ruleForm .el-form-item__label {
  font-weight: 600 !important;
  font-size: 12px;
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
    // width: 75px;
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
</style>
<style>
.el-drawer__body {
  overflow: auto;
  border-left: 1px solid #e9e9e9;
  padding-left: 15px;
}

.el-drawer__header {
  display: none !important;
}

.el-drawer__open .el-drawer.rtl {
  top: 46px !important;
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
</style>
