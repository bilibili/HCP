<template>
  <div class="centent" id="map" ref="echarts">
    <div class="navBarTop">
      <div>
        <el-card>
          <div style="display: flex;flex-direction: row;flex-wrap: nowrap;justify-content: space-between;">
            <h1 style="font-size: 16px; margin: 0">导航</h1>
          </div>
          <div style="display: flex">
            <div style="flex: 1">
              <h4 style="margin-top: 20px;color: rgb(51, 51, 51);font-size: 14px;font-weight: 600;">
                <img src="@/assets/icon/product.svg" width="18px;" style="vertical-align: middle" alt="" />
                常用云产品
                <el-popover placement="bottom-start" width="550" trigger="click" popper-class="popperDiaBox">
                  <span style="font-size: 12px;color: rgb(64, 158, 255);cursor: pointer;font-weight: 500;"
                    slot="reference"><i class="el-icon-plus"></i>添加快捷入口</span>
                  <div style="display: flex">
                    <div class="openPBox">
                      <h4 style="margin: 8px 0">计算</h4>
                      <span style="font-size: 12px;cursor: pointer;display: block;margin-top: 5px;" :underline="false"
                        @click="beLickTo('云服务器')" v-on:mouseover="hoverLikePro('yfwqjia')"
                        v-on:mouseout="hoverLikeProOut('yfwqjia')">
                        <i class="el-icon-star-on openColor"
                          v-if="productLikeArrLast && productLikeArrLast.indexOf('云服务器') != -1"></i>
                        云服务器
                        <i class="el-icon-plus" v-if="showjiaList.yfwqjia" style="margin-left: 3px; color: #409eff"></i>
                      </span>
                    </div>
                  </div>
                </el-popover>
              </h4>
              <div v-if="productLikeArrLast && productLikeArrLast.length != 0">
                <el-col :span="6" v-for="(pItem, index) in productLikeArrLast" :key="index" style="margin-bottom: 10px">
                  <el-tag closable @click="toProduct(pItem)" @close="beLickTo(pItem)" style="width: 93%" class="tagss">
                    {{ pItem }}
                  </el-tag>
                </el-col>
              </div>
            </div>
            <div style="flex: 1; margin-left: 25px">
              <h4 style="margin-top: 20px;color: rgb(51, 51, 51);font-size: 14px;font-weight: 600;">
                <img src="@/assets/icon/handle.svg" alt="" width="20px;" style="vertical-align: middle" />
                常用操作
                <el-popover placement="bottom-start" width="450" trigger="click" popper-class="popperDiaBox2">
                  <template v-if="baseMessage.templates">
                    <el-col :span="12" style="font-size: 12px; cursor: pointer"
                      v-for="(tItem, tindex) in baseMessage.templates" :key="tindex">
                      <div :underline="false" @click="beLickToPro(tItem.title)" v-on:mouseover="hoverLike(tindex)"
                        v-on:mouseout="hoverOut(tindex)">
                        <i class="el-icon-star-on openColor" v-if="projectlikeArrLast &&
                          projectlikeArrLast.indexOf(tItem.title) != -1
                          "></i>
                        {{ tItem.title }}
                        <i class="el-icon-plus" v-if="tItem.showjia"
                          style="margin-left: 3px;cursor: pointer;color: #409eff;"></i>
                      </div>
                    </el-col>
                  </template>
                  <span style="font-size: 12px;color: rgb(64, 158, 255);cursor: pointer;font-weight: 500;"
                    slot="reference">
                    <i class="el-icon-plus"></i>添加快捷入口
                  </span>
                </el-popover>
              </h4>
              <div style="margin-top: -2px">
                <div v-if="projectlikeArrLast && projectlikeArrLast.length != 0">
                  <el-col :span="8" style="margin-bottom: 10px" v-for="(pItem, index) in projectlikeArrLast"
                    :key="index">
                    <el-tag closable @click="taskClick(pItem)" @close="beLickToPro(pItem)" style="width: 93%"
                      class="tagss">
                      {{ pItem }}
                    </el-tag>
                  </el-col>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </div>
    </div>
    <div class="conBottom">
      <!-- 右下视图 -->
      <el-card style="height: 197px">
        <h4 style="padding: 0px 20px; margin: 0px; display: flex">
          <div style="flex: 1; text-align: center">
            <img src="@/assets/icon/person.png" style="width: 42px; margin-right: 10px" alt="" />
          </div>
          <div style="flex: 1.1; margin-left: -18px; position: relative;display: flex;flex-direction: column;align-items: center;">
            {{ username }}
            <div
              style="width: 67px;height: 18px;border: 1px solid rgb(204, 204, 204);border-radius: 20px;font-weight: 500;font-size: 12px;line-height: 17px;text-align: center;margin-top: 5px;">
              {{ loginUserInfo.RoleName }}
            </div>
          </div>
        </h4>
        <div style="display: flex;position: relative;margin-top: 7px;cursor: pointer;">
          <div style="flex: 1; text-align: center" @click="jumpToprojectMes">
            <h5>关联项目</h5>
            <span>{{ baseMessage.countProject }}</span>
          </div>
        </div>
      </el-card>
      <el-card style="height: 197px; margin-top: 10px; cursor: pointer">
        <h1 style="font-size: 16px">工作台</h1>
        <div class="workPlace">
          <div style="flex: 1; text-align: center">我的待办</div>
          <span style="position: absolute; left: 50%; color: #ccc">|</span>
          <div style="flex: 1; text-align: center; color: #409eff" @click="toProduct('我的待办')">
            {{  baseMessage.todo  }}
          </div>
        </div>
        <div class="workPlace">
          <div style="flex: 1; text-align: center">活跃申请</div>
          <span style="position: absolute; left: 50%; color: #ccc">|</span>
          <div style="flex: 1; text-align: center; color: #409eff" @click="toProduct('我的申请')">
            {{ baseMessage.myJob }}
          </div>
        </div>
      </el-card>
      <el-card style="margin-top: 10px; flex: 1">
        <h1 style="font-size: 16px">最新公告</h1>
        <div style="font-size: 12px">暂无~</div>
      </el-card>
    </div>
    <!-- 新建任务弹框 infoAddAttr 表格行信息 infoAddForm表单提交信息 reviewStatus弹框状态-->
    <div v-if="reviewStatus">
      <formtem-Dialog :baseConfigData="baseConfigData" :nodeCheckName="nodeCheckName" :detailOrLook="detailOrLook"
        @addClose="addClose" :infoAddAttr="infoAddAttr" :detailOrsigle="detailOrsigle" :infoAddForm="infoAddForm"
        :reviewStatus="reviewStatus" :newTaskName="newTaskName" :historyTableList="historyTableList"
        :historyparamreferenceShow="historyparamreferenceShow" :childNodename="childNodename" :checkList="checkList"
        :FormTemplatjobMode="FormTemplatjobMode"></formtem-Dialog>
    </div>
  </div>
</template>

<script>
import { userMixin } from '@/components/mixin/user';
import formtemDialog from '@/views/components/formTemp';
import { mapGetters } from 'vuex';
import Http from "@/components/api/services";
export default {
  components: {
    formtemDialog,
  },
  mixins: [userMixin],
  computed: {
    ...mapGetters(['loginUserName', 'formPolicyDomain',"loginUserInfo"]),
  },
  data() {
    return {
      bulletinBoard: {},
      datasuccess: false,
      searchProduct: '',
      activeName: 'first',
      FormTemplatjobMode: '',
      timeValues: '',
      productLoading: false,
      trendLoading: false,
      drawDetail: false,
      handleObject: {
        '通用云产品创建': 'FormCreateCloudProduct',
        '通用云产品改配': 'FormConfigCloudProduct',
        '通用云产品清退': 'FormDeleteCloudProduct',
      },
      handleProduct: {
        '云服务器': '/compute/server/instance/list',
      },
      searchTop: '',
      username: '',
      baseMessage: {
        countProject: 0,
        myJob: 0,
        todo: 0,
        templates: [
          {
            name: 'FormCreateCloudProduct',
            title: '通用云产品创建',
          },
          {
            name: 'FormConfigCloudProduct',
            title: '通用云产品改配',
          },
          {
            name: 'FormDeleteCloudProduct',
            title: '通用云产品清退',
          },
        ],
      },
      projectArr: [],
      productLikeArr: [],
      productLikeArrLast: [],
      projectlikeArr: [],
      projectlikeArrLast: [],
      nodeCheckName: '',
      infoAddForm: {},
      infoAddBFForm: {},
      infoAddAttr: [],
      dialogVisible: false,
      reviewStatus: false,
      detailOrLook: false,
      detailOrsigle: '',
      drawDetail: false,
      childNodename: '',
      checkList: [],
      baseConfigData: {},
      historyTableList: {},
      historyparamreferenceShow: [],
      projectArrsList: [],
      providerList: [],
      productList: {
        products: [
          {
            name: 'CloudServer',
            title: '云服务器',
          }
        ],
      },
      productListLast: [],
      showjiaList: {
        yfwqjia: false,
        ljsjia: false,
        vpcjia: false,
        zwjia: false,
        zsjia: false,
        xxdl1jia: false,
        aqzjia: false,
        fzjhjia: false,
        eipjia: false,
        cdnjia: false,
        zs: false,
        rdsjia: false,
        resjia: false,
        dxccjia: false,
        yypjia: false,
        rqfwjia: false,
        byjs: false,
        MongoDBjia: false,
        natjia: false,
        DNSjia: false,
        dkjia: false,
      },
      proValueArr1: [],
      projectSelArr: [],
    };
  },
  created() {
    this.username = sessionStorage.getItem('username');
    this.projectArrsList = sessionStorage.getItem('checkProjectArr') ? sessionStorage.getItem('checkProjectArr').split(',') : [];
    let pArr = localStorage.getItem('productLikeArr');
    if (pArr) {
      if (pArr.indexOf(',') == -1) {
        this.productLikeArr = [pArr];
      } else {
        this.productLikeArr = pArr.split(',');
      }
    } else {
      this.productLikeArr = [];
    }
    this.productLikeArrLast = this.productLikeArr;
    let proArr = localStorage.getItem('projectlikeArr');
    if (proArr) {
      if (proArr.indexOf(',') == -1) {
        this.projectlikeArr = [proArr];
      } else {
        this.projectlikeArr = proArr.split(',');
      }
    } else {
      this.projectlikeArr = [];
    }
    this.projectlikeArrLast = this.projectlikeArr;
  },
  mounted() {
    console.log(this.loginUserInfo);
    this.timeValues = [
      this.$moment().subtract(7, 'day').format('YYYY-MM-DD'),
      this.$moment().format('YYYY-MM-DD'),
    ];
    this.$EventBus.$on('checkProjectArr', (value) => {
      this.projectArrsList = value;
    });
    Http.getDashCount().then((res) => {
      this.baseMessage.countProject = res.data.data.project_num
      this.baseMessage.todo = res.data.data.todo_num
      this.baseMessage.myJob = res.data.data.application_num
    })
  },
  methods: {
    changeClor(index) {
      this.productListLast[index].style = true;
    },
    changeclorold(index) {
      this.productListLast[index].style = false;
    },
    hoverLike(index) {
      this.baseMessage.templates[index].showjia = JSON.parse(
        JSON.stringify(true)
      );
      this.$forceUpdate();
    },
    hoverOut(index) {
      this.baseMessage.templates[index].showjia = false;
    },
    hoverLikePro(name) {
      this.showjiaList[name] = true;
      this.$forceUpdate();
    },
    hoverLikeProOut(name) {
      this.showjiaList[name] = false;
    },
    jumpToprojectMes() {
      this.$router.push({
        path: '/general/manage/project/list',
        query: {
          status: 'dashboardproject',
        },
      });
    },
    jumpTocuser() {
      this.$router.push({
        path: '/general/manage/account/list',
        query: {
          status: 'dashboardaccount',
        },
      });
    },
    // 收藏按钮 云产品
    beLickTo(name) {
      this.productLikeArr = localStorage.getItem('productLikeArr');
      if (this.productLikeArr) {
        let prodNew = [];
        if (this.productLikeArr.indexOf(',') == -1) {
          prodNew = [this.productLikeArr];
        } else {
          prodNew = this.productLikeArr.split(',');
        }
        if (prodNew.indexOf(name) == -1) {
          prodNew.push(name);
        } else {
          for (let proItem in prodNew) {
            if (prodNew[proItem] == name) {
              prodNew.splice(proItem, 1);
            }
          }
        }
        this.productLikeArr = prodNew;
      } else {
        this.productLikeArr = [];
        this.productLikeArr.push(name);
      }
      localStorage.setItem('productLikeArr', this.productLikeArr);
      this.productLikeArrLast = JSON.parse(JSON.stringify(this.productLikeArr));
    },
    // 收藏按钮 项目
    beLickToPro(name) {
      this.projectlikeArr = localStorage.getItem('projectlikeArr');
      if (this.projectlikeArr) {
        let prodNew = [];
        if (this.projectlikeArr.indexOf(',') == -1) {
          prodNew = [this.projectlikeArr];
        } else {
          prodNew = this.projectlikeArr.split(',');
        }
        if (prodNew.indexOf(name) == -1) {
          prodNew.push(name);
        } else {
          for (let proItem in prodNew) {
            if (prodNew[proItem] == name) {
              prodNew.splice(proItem, 1);
            }
          }
        }
        this.projectlikeArr = prodNew;
      } else {
        this.projectlikeArr = [];
        this.projectlikeArr.push(name);
      }
      localStorage.setItem('projectlikeArr', this.projectlikeArr);
      this.projectlikeArrLast = JSON.parse(JSON.stringify(this.projectlikeArr));
    },
    // 常用云产品跳转
    toProduct(name) {
      if (this.handleProduct[name]) {
        let jumpurls = this.handleProduct[name];
        this.$router.push({
          path: jumpurls,
          query: {
            projectId: this.proValueArr1,
          },
        });
      }
    },
    // 常用操作跳表单
    // 任务类型点击taskname英文名称，请求接口 chinaName中文名称
    taskClick(chinaName) {
      this.checkList = [];
      this.childNodename = '';
      this.dialogVisible = false;
      //中文名称
      this.newTaskName = chinaName;
      this.nodeCheckName = this.handleObject[chinaName];
      this.reviewStatus = true;
    },
    // 发起任务弹框关闭
    addClose(val) {
      this.reviewStatus = val;
    },
    // 详情遮罩层弹框关闭
    addCloseDraw() {
      this.drawDetail = false;
    },
  },
  destroyed() {
    this.$EventBus.$off('checkProjectArr');
  },
};
</script>
<style lang="scss" scoped>
.hintCardBoxMain {
  display: flex;
  justify-content: center;
}

.hintCardBoxMain .svgImg {
  width: 60px;
  height: 60px;
  margin-right: 20px;
  color: #46bae9 !important;
}

.hintCardBoxMain .svgImg1 {
  width: 60px;
  height: 60px;
  margin-right: 20px;
  color: #fb7299 !important;
}

.hintCardBoxMain .svgImg2 {
  width: 60px;
  height: 60px;
  margin-right: 20px;
  color: #ffbd00 !important;
}

.hintcardBox {
  height: 70px;
  color: #666;
  padding-top: 5px; //   text-align: center;
}

.hintcardBox span {
  font-size: 13px;
  margin-left: 5px;
  font-weight: 800;
}

.hintcardBox p {
  font-size: 13px;
  margin: 0;
  margin-top: 10px;
}

.hintcardBox div {
  position: relative;
  font-size: 32px;
  line-height: 32px;
  font-weight: 700;
  height: 32px;
  min-width: 100px;
  position: relative;
}

.el-button {
  padding: 7px 8px !important;
  background: rgba(203, 220, 248, 0.2);
  margin-left: 0px !important;
  margin-right: 10px;

  span {
    font-size: 12px !important;
  }
}

.centent {
  width: calc(100vw);
  height: calc(100vh - 47px); 
  background: rgba(203, 220, 248, 0.2);
  padding: 10px;
  display: flex;
  overflow: auto;
}

.navBarTop {
  height: calc(100vh - 70px);
  flex: 4;
  width: 100%;
  display: flex;
  flex-direction: column;
}

.conBottom {
  flex: 1;
  display: flex;
  flex-direction: column;
  margin-left: 15px;
  height: calc(100vh - 70px);
}

.projectPage {
  border: 1px solid rgba(203, 220, 248, 0.1);
  padding: 5px 10px 15px;
  background: rgba(203, 220, 248, 0.2);
  cursor: pointer;
}

.workPlace {
  padding: 15px;
  margin-top: 15px;
  position: relative;
  background: rgba(203, 220, 248, 0.2);
  display: flex;
  font-size: 13px;
}

.product_box {
  display: flex;
  position: relative;
  border: 1px solid rgba(203, 220, 248, 0.2);
  width: 90%;
  border-radius: 3px;
  float: left;
  margin: 7px 10px;
  height: 68px;
  background: rgba(203, 220, 248, 0.2);
}

.openPBox {
  flex: 1;
}

.popperDiaBox .el-link {
  font-size: 12px;
  display: block;
  margin-top: 10px;
}

.el-link.el-link--default:hover {
  color: #606266 !important;
}

.openColor {
  color: rgb(24, 67, 117);
}

.el-tag {
  height: 30px;
  line-height: 30px;
  position: relative;
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
  color: #606266;
}
</style>
<style>
.tagss .el-icon-close {
  position: absolute;
  right: -1px;
  top: 7px;
}

.popperDiaBox /deep/ .popper__arrow,
.popperDiaBox2 /deep/ .popper__arrow {
  display: none !important;
}

.centent /deep/ .el-button span {
  font-size: 12px !important;
}

.timeCheck /deep/ .el-input__inner {
  height: 32px !important;
}

.timeCheck /deep/ .el-input__prefix {
  top: 2px !important;
}
</style>
