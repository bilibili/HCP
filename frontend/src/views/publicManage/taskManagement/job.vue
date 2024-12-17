<template>
  <div class="autoBox installati_box">
    <div>
      <div
        style="border-bottom: 1px solid #e6e6e6;  padding-bottom: 10px;  margin-bottom: 10px;  margin-top: 8px;  font-size: 14px;">
        <span style="font-weight: 600">
          {{ $route.meta.topTitle }}
          <i class="el-icon-arrow-right" style="font-size: 14px; font-weight: 600; margin: 0 6px"
            v-if="$route.meta.topTitle"></i>
          {{ $route.meta.pageTitle }}
        </span>
      </div>
      <div
        style="display: flex;  margin-bottom: 10px;  align-items: center;  margin-top: -3px;  justify-content: space-between;">
        <el-button type="primary" @click="changeAddnew">发起任务</el-button>
        <div style="display: flex">
          <el-tag class="speical_tag" :key="index + 1" v-for="(item, index) in searchObjTag" closable
            :disable-transitions="false" style="margin-right: 3px" @close="handleCloseTag(item)">
            <span v-if="searchObj.status > 2">
              {{ item.name + ":[" + batcharrLength[0] + "," + batcharrLength[1] + "等" + searchObj.status + "个]" }}
            </span>
            <span v-else>
              {{ item.name + ":[" + item.value + "]" }}
            </span>
          </el-tag>
          <el-popover placement="bottom" trigger="click" :popper-class="classWidth" @hide="closePopover"
            v-model="poHide">
            <div style="display: flex">
              <div style="flex: 1">
                <el-cascader-panel ref="cascade" :options="searchMoreArr" @change="searchForm" v-model="cascaderMoreArr"
                  :checkStrictly="true" @expand-change="searchFormchange">
                  <template #default="{ node, data }">
                    <span>{{ data.label }}</span>
                    <i v-if="node.isLeaf" class="el-icon-arrow-right el-cascader-node__postfix"
                      style="top: 10px; right: 10px"></i>
                  </template>
                </el-cascader-panel>
              </div>
              <div style="flex: 2.5; display: flex; margin-top: 10px" v-if="classWidth == 'bigwidth'">
                <el-form :model="searchObj" ref="searchObjForm" label-width="0px" class="demo-ruleForm">
                  <el-form-item>
                    <el-select multiple placeholder="请选择" filterable clearable style="width: 100%"
                      v-model="searchObj.status">
                      <el-option v-for="(oitem, index) in sreachFormArr" :value="oitem.en" :label="oitem.zh"
                        :key="index"></el-option>
                    </el-select>
                  </el-form-item>
                  <div style="text-align: right">
                    <el-button type="primary" plain @click="textSearchJob(1)">搜 索</el-button>
                  </div>
                </el-form>
              </div>
            </div>
            <el-button style="font-size: 12px;  margin-right: 10px;  position: relative;  top: 1px;"
              @click="searchShowStatus" slot="reference">高级搜索<i style="margin-left: 5px"
                class="el-icon-caret-bottom"></i></el-button>
          </el-popover>
          <div style="margin-left: 10px">
            <el-input clearable @clear="textSearchJob(1)" style="width: 250px" v-model="sreachForm.Keyword"
              placeholder="请输入任务名称、对象名称进行搜索" @keydown.enter.native="textSearchJob(1)"></el-input>
          </div>
        </div>
      </div>
      <div>
        <el-table :data="looktableData" stripe style="width: 100%">
          <el-table-column type="selection" width="40"></el-table-column>
          <el-table-column prop="title" label="任务标题" sortable>
            <template slot-scope="scoped">
              <el-link type="primary" style="margin: 0px 5px; font-size: 12px"
                :underline="false" @click="goDetialpage(scoped.row)">{{ scoped.row.title }}</el-link>
            </template>
          </el-table-column>
          <el-table-column prop="user.name" label="申请人">
            <template slot-scope="scoped">
              <span v-if="scoped.row.user">
                {{ scoped.row.user.name }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="create_time" label="创建时间" sortable width="160">
            <template slot-scope="scoped">
              <span v-if="scoped.row.create_time">
                {{ scoped.row.create_time | filterTimeShow }}
              </span>
              <span v-else> - </span>
            </template>
          </el-table-column>
          <el-table-column prop="endTime" label="完成时间" sortable width="160">
            <template slot-scope="scoped">
              <span v-if="scoped.row.endTime">
                {{ scoped.row.endTime | filterTimeShow }}
              </span>
              <span v-else> - </span>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" sortable width="90">
            <template slot-scope="scoped">
              <span v-if="scoped.row.status == 'running'">
                <el-link :underline="false" type="warning" class="el-icon-time" style="font-size: 12px"><span
                    style="margin-left: 5px">{{ organizationArr[scoped.row.status] }}</span></el-link>
              </span>
              <span v-if="scoped.row.status == 'success' || scoped.row.status == 'success_ack'">
                <el-link :underline="false" type="success" class="el-icon-success" style="font-size: 12px"><span
                    style="margin-left: 5px">{{ organizationArr[scoped.row.status] }}</span></el-link>
              </span>
              <span
                v-if="scoped.row.status == 'idle' || scoped.row.status == 'auth' || scoped.row.status == 'auth_wait' || scoped.row.status == 'auth_req' || scoped.row.status == 'auth_wait_sync' || scoped.row.status == 'create_wait' || scoped.row.status == 'creating'">
                <el-link :underline="false" type="info" class="el-icon-remove" style="font-size: 12px"><span
                    style="margin-left: 5px">{{ organizationArr[scoped.row.status] }}</span></el-link>
              </span>
              <span v-if="scoped.row.status == 'failure' || scoped.row.status == 'auth_deny'">
                <el-link :underline="false" type="danger" class="el-icon-error" style="font-size: 12px"><span
                    style="margin-left: 5px">{{ organizationArr[scoped.row.status] }}</span></el-link>
              </span>
              <span v-if="scoped.row.status == 'confirming'">
                <el-link :underline="false" type="warning" class="el-icon-question" style="font-size: 12px"><span
                    style="margin-left: 5px">{{ organizationArr[scoped.row.status] }}</span></el-link>
              </span>
              <span v-if="!organizationArr[scoped.row.status]">
                <el-link :underline="false" type="info" class="el-icon-question" style="font-size: 12px"><span
                    style="margin-left: 5px">{{ organizationArr[scoped.row.status] }}</span></el-link>
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="progress" label="进度" align="center" width="100">
          </el-table-column>
          <el-table-column fixed="right" label="操作" align="center" width="160">
            <template slot-scope="scoped">
              <el-link style="font-size: 12px; color: #409eff; margin: 0 2px"
                @click="goDetialpage(scoped.row)">详情</el-link>
            </template>
          </el-table-column>
        </el-table>
        <div style="text-align: center; margin: 10px 0 25px 0">
          <el-pagination @current-change="changehandleCurrent" @size-change="handleSizeChange" :page-size="pageSize"
            :current-page="pageNum" :page-sizes="[10, 20, 30, 50, 100, 150, 200]"
            layout="total, sizes, prev, pager, next, jumper" :total="totle">
          </el-pagination>
        </div>
      </div>
      <!-- 选择任务类型弹框 -->
      <el-dialog title="选择任务类型" :visible.sync="dialogVisible" class="diaBox diaboxTi" width="700px">
        <div style="display: flex; padding: 10px 20px 20px 20px">
          <div style="width: 100px; min-width: 100px; border-right: 1px solid #ccc">
            <el-badge value="hot" class="item">
              <div>
                <div @click="activeTitle = '常用任务'" class="itemHover itemHover1"
                  :style="activeTitle == '常用任务' ? 'line-height: 22px;margin-bottom:18px;font-size: 15px;cursor: pointer;color: #409EFF;color: rgb(64, 158, 255);border-bottom: 2px solid rgba(64, 158, 255,.7);font-weight:600;' : 'line-height: 22px;margin-bottom:18px;font-size: 15px;cursor: pointer;'">
                  常用任务
                </div>
              </div>
            </el-badge>
            <div>
              <div class="itemHover itemHover1" @click="activeTitle = '系统管理'"
                :style="activeTitle == '系统管理' ? 'line-height: 22px;margin-bottom:18px;font-size: 15px;cursor: pointer;color: #409EFF;color: rgb(64, 158, 255);border-bottom: 2px solid rgba(64, 158, 255,.7);font-weight:600;display: inline-block;' : 'line-height: 22px;margin-bottom:18px;font-size: 15px;cursor: pointer;'">
                系统管理
              </div>
            </div>
          </div>
          <div style="flex-grow: 1; margin-left: 40px" v-if="activeTitle == '常用任务'">
            <div>
              <el-badge value="hot">
                <div
                  style="line-height: 30px;  font-size: 14px;  color: #000;  margin-bottom: 5px;  color: rgba(64, 158, 255, 1);  font-weight: 600;">
                  资源编排
                </div>
              </el-badge>
              <div>
                <span class="itemHover"
                  style="line-height: 30px;  cursor: pointer;  font-size: 13px;  margin-right: 40px;"
                  @click="taskClick('FormCreateCloudProduct', '云产品新建')">新建</span>
                <span class="itemHover"
                  style="line-height: 30px;  cursor: pointer;  font-size: 13px;  margin-right: 40px;"
                  @click="taskClick('FormConfigCloudProduct', '云产品改配','FormConfigCloudServer')">改配</span>
                <span class="itemHover"
                  style="line-height: 30px;  cursor: pointer;  font-size: 13px;  margin-right: 40px;"
                  @click="taskClick('FormDeleteCloudProduct', '云产品清退','FormDeleteCloudServer')">清退</span>
              </div>
              <div
                style="line-height: 30px;  font-size: 14px;  color: #000;  margin-bottom: 5px;  color: rgba(64, 158, 255, 1);  font-weight: 600;  margin-top: 10px;">
                云服务器
              </div>
              <div>
                <span class="itemHover"
                  style="line-height: 30px;  cursor: pointer;  font-size: 13px;  margin-right: 40px;"
                  @click="taskClick('FormRebootCloudServer', '云产品重启')">重启</span>
                <span class="itemHover"
                  style="line-height: 30px;  cursor: pointer;  font-size: 13px;  margin-right: 40px;"
                  @click="taskClick('FormReinstallCloudServer', '云产品重装')">重装</span>
                <span class="itemHover"
                  style="line-height: 30px;  cursor: pointer;  font-size: 13px;  margin-right: 40px;"
                  @click="taskClick('FormRenameBatchCloudServer', '云产品改名','FormRenameCloudServer')">改名</span>
              </div>
            </div>
          </div>
          <div style="flex-grow: 1; margin-left: 40px" v-if="activeTitle == '系统管理'">
            <div>
              <div
                style="line-height: 30px;  font-size: 14px;  color: #000;  margin-bottom: 5px;  color: rgba(64, 158, 255, 1);  font-weight: 600;">
                项目管理
              </div>
              <div>
                <span class="itemHover"
                  style="line-height: 30px;  cursor: pointer;  font-size: 13px;  margin-right: 40px;"
                  @click="taskClick('FormAdminCreateCloudProjectConfig', '项目新建')">项目新建</span>
                <span class="itemHover"
                  style="line-height: 30px;  cursor: pointer;  font-size: 13px;  margin-right: 40px;"
                  @click="taskClick('FormAdminAppendCloudProjectAccountConfig', '账号新建')">账号新建</span>
              </div>
              <div
                style="line-height: 30px;  font-size: 14px;  color: #000;  margin-bottom: 5px;  color: rgba(64, 158, 255, 1);  font-weight: 600;  margin-top: 10px;">
                账号管理
              </div>
              <div>
                <span class="itemHover"
                  style="line-height: 30px;  cursor: pointer;  font-size: 13px;  margin-right: 40px;"
                  @click="taskClick('FormAdminCreateUser', '用户新建')">用户新建</span>
                <span class="itemHover" style="line-height: 30px; cursor: pointer; font-size: 13px"
                  @click="taskClick('FormJoinCloudProjectConfig', '用户加入项目')">用户加入项目</span>
              </div>
            </div>
          </div>
        </div>
      </el-dialog>
      <!-- 新建任务弹框 infoAddAttr 表格行信息 infoAddForm表单提交信息 reviewStatus弹框状态-->
      <div v-if="reviewStatus">
        <formtem-Dialog :baseConfigData="baseConfigData" :nodeCheckName="nodeCheckName" :detailOrLook="detailOrLook"
          @addClose="addClose" :infoAddAttr="infoAddAttr" :detailOrsigle="detailOrsigle" :infoAddForm="infoAddForm"
          :reviewStatus="reviewStatus" :newTaskName="newTaskName" :historyTableList="historyTableList"
          :historyparamreferenceShow="historyparamreferenceShow" :childNodename="childNodename" :checkList="checkList"
          :FormTemplatjobMode="FormTemplatjobMode" :jobIdStatus="jobIdStatus"></formtem-Dialog>
      </div>
      <div v-if="drawDetail">
        <draw-log @addCloseDraw="addCloseDraw" :jobIds="jobIds"></draw-log>
      </div>
    </div>
  </div>
</template>

<script>
//例如：import 《组件名称》 from '《组件路径》';
import formtemDialog from "@/views/components/formTemp";
import drawLog from "./taskDesk";
import { userMixin } from "@/components/mixin/user";
import { jobReview } from "@/components/mixin/jobReview";
import { mapGetters } from "vuex";
import Http from "@/components/api/services";
export default {
  //import引入的组件需要注入到对象中才能使用
  props: [],
  components: {
    formtemDialog,
    drawLog,
  },
  mixins: [userMixin, jobReview],
  computed: {
    ...mapGetters(["loginUserName", "formPolicyDomain"]),
  },
  data() {
    //这里存放数据
    return {
      poHide: false,
      searchMoreArr: [
        {
          value: "statusP",
          label: "状态",
          children: [
            {
              value: "status",
              label: "选择",
            },
          ],
        },
      ],
      cascaderMoreArr: [],
      searchObj: {},
      searchObjTag: [],
      sreachFormArr: [
        {
          en: "auth_wait",
          zh: "审批中",
        },
        {
          en: "auth_deny",
          zh: "拒绝",
        },
        {
          en: "running",
          zh: "进行中",
        },
        {
          en: "failure",
          zh: "失败",
        },
        {
          en: "success",
          zh: "成功",
        },
      ],
      FormTemplatjobMode: "",
      childNodename: "",
      checkList: [],
      baseConfigData: {},
      historyTableList: {},
      historyparamreferenceShow: [],
      looktableData: [],
      sreachForm: {
        Keyword: "",
      },
      pageSize: 20,
      pageNum: 1,
      totle: 0,
      organizationArr: {
        success: "成功",
        success_ack: "结单",
        failure: "失败",
        running: "进行中",
        auth_wait: "审核中",
        auth_req: "审核中",
        auth_wait_sync: "审核中",
        auth: "审核中",
        idle: "等待",
        confirming: "待确认",
        auth_deny: "已拒绝",
        ordinary: "普通工单",
        stability: "稳定性工单",
        maintenance: "维护工单",
        create_wait: "创建中",
        creating: "创建中",
      },
      newTaskName: "",
      nodeCheckName: "",
      infoAddForm: {},
      infoAddBFForm: {},
      infoAddAttr: [],
      dialogVisible: false,
      reviewStatus: false,
      detailOrLook: false,
      detailOrsigle: "",
      drawDetail: false,
      jobIds: "",
      workType: "",
      username: "",
      activeTitle: "常用任务",
      classWidth: "",
    };
  },
  watch: {
    $route: {
      handler() {
        if (this.$route.query.jobId) {
          let JobCode = this.$route.query.jobId;
          if (this.$route.query.objectId) {
            let oId = this.$route.query.objectId;
            this.$router.replace({
              path: location.pathname,
              query: {
                jobId: JobCode,
                objectId: oId,
              },
            });
          } else {
            this.$router.replace({
              path: location.pathname,
              query: {
                jobId: JobCode,
              },
            });
          }
          this.jobIds = JobCode;
          this.drawDetail = true;
        }
      },
      deep: true,
    },
  },
  //方法集合
  methods: {
    searchForm() { },
    //高级搜索,默认选择第一个子菜单时
    searchFormchange(item) {
      this.cascaderMoreArr = [item, this.searchMoreArr[0].children[0].value];
      this.classWidth = "bigwidth";
    },
    closePopover() {
      this.classWidth = "smallwidth";
    },
    searchShowStatus() {
      this.classWidth = "smallwidth";
      this.cascaderMoreArr = null;
      this.$refs.cascade.activePath = [];
      this.$refs.cascade.clearCheckedNodes = null;
      this.$refs.cascade.checkedValue = [];
      this.searchFormData = {};
    },
    // 搜索tag移除、
    handleCloseTag(item) {
      this.searchObj.status = [];
      this.textSearchJob();
    },
    // 发起按钮请求都有哪些新建任务
    changeAddnew() {
      this.dialogVisible = true;
    },
    // 任务类型点击taskname英文名称，请求接口 chinaName中文名称
    taskClick(taskname, chinaName,childname) {
      this.checkList = [];
      this.childNodename = childname || '';
      this.dialogVisible = false;
      //中文名称
      this.newTaskName = chinaName;
      this.nodeCheckName = taskname;
      if (taskname == "FormRenameBatchCloudServer") {
        this.childNodename = "FormRenameCloudServer";
      }
      this.detailOrLook = false;
      this.historyTableList = {};
      this.reviewStatus = true;
    },
    // 发起任务弹框关闭
    addClose(val) {
      this.reviewStatus = val;
      this.textSearchJob(1);
    },
    // 详情遮罩层弹框关闭
    addCloseDraw() {
      this.drawDetail = false;
    },
    // 查看任务列表按钮+接口
    textSearchJob(Size) {
      if (Size) {
        this.pageNum = 1;
      }
      this.looktableData = [];

      var postData = {
        "schema": "job",
        "order":"id DESC",
        page_size:this.pageSize,
        page_num:this.pageNum,
        where:{}
      }

      this.searchObjTag = [];
      if (this.searchObj.status && this.searchObj.status.length > 0) {
        postData.where.status_IN = this.searchObj.status;
        let eItemTxt = "";
        this.searchObj.status.map((item, ii) => {
          this.sreachFormArr.map((sitem, si) => {
            if (sitem.en == item) {
              eItemTxt += `${sitem.zh},`;
            }
          });
        });
        eItemTxt = eItemTxt.slice(0, eItemTxt.length - 1);
        this.searchObjTag.push({
          name: "状态",
          value: eItemTxt,
        });
      }
      Http.getJobList(postData).then((response) => {
        this.totle = response.data.data.total
        this.looktableData =  response.data.data.data || [];
      })
    },
    // 详情按钮
    goDetialpage(row) {
      this.$router.replace({
        path: location.pathname,
        query: {
          jobId: row.id,
        },
      });
      this.jobIds = row.id;
      this.drawDetail = true;
    },
    changehandleCurrent(val) {
      this.pageNum = val;
      this.textSearchJob();
    },
    handleSizeChange(val) {
      this.pageNum = 1;
      this.pageSize = val;
      this.textSearchJob();
    },
  },
  created() {
    this.pageSize = Number(this.$route.query._pageSize) ? Number(this.$route.query._pageSize) : this.pageSize;
    this.pageNum = Number(this.$route.query._pageNumber) ? Number(this.$route.query._pageNumber) : this.pageNum;
    if (this.$route.query.jobId && !this.$route.query.objectId) {
      let JobCode = Number(this.$route.query.jobId);
      this.$router.replace({
        path: location.pathname,
        query: {
          jobId: JobCode,
        },
      });
      this.jobIds = JobCode;
      this.drawDetail = true;
    }
    if (this.$route.query.jobId && this.$route.query.objectId) {
      let JobCode = this.$route.query.jobId;
      let ObjCode = this.$route.query.objectId;
      this.$router.replace({
        path: location.pathname,
        query: {
          jobId: JobCode,
          objectId: ObjCode,
        },
      });
      this.jobIds = JobCode;
      this.drawDetail = true;
    }
  },
  //生命周期 - 挂载完成（可以访问DOM元素）
  mounted() {
    // 查看任务列表
    if (this.$route.query.nodeName && this.$route.query.nodeTitle) {
      this.taskClick(this.$route.query.nodeName, this.$route.query.nodeTitle);
    }
    this.workType = this.$route.name;
    this.username = sessionStorage.getItem("username");
    this.textSearchJob();
  },
};
</script>
<style scoped lang="scss">
.installati_box {
  width: 100%;

  .diaBox {
    .title {
      width: 100%;
      height: 30px;
      border-bottom: 1px solid #666;
      font-size: 16px;
      color: #1890ff;
    }

    .liTitle {
      margin-top: 15px;
      cursor: pointer;
      padding-left: 15px;
    }
  }

  .batchH4 {
    margin: 0;
    padding: 0;
    color: #333333;
    font-size: 15px;
    font-weight: 600;
    text-indent: 20px;
    border-left: 4px solid #23ade5;
    line-height: 24px;
  }
}

.popper__arrow {
  display: none !important;
}

.el-select-dropdown {
  max-width: 500px !important;
}

.el-cascader-panel.is-bordered {
  border: none !important;
}

.el-cascader-menu {
  border: none !important;
}
</style>
<style>
.smallwidth {
  width: 206px !important;
}

.bigwidth {
  width: 595px !important;
  right: 3% !important;
  left: unset !important;
}
</style>
<style scoped>
.installati_box /deep/ .el-dialog__body {
  padding: 0px 20px 10px;
}

.itemHover1 {
  color: #000;
}

.itemHover:hover {
  color: #409eff;
}

.diaboxTi /deep/ .el-dialog__title {
  font-size: 16px;
}
</style>
