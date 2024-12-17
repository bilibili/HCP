<template>
  <div class="navbar navbarBoxMain">
    <div style="line-height:40px;background:#184375;color:#fff;position: relative;z-index: 9999;padding-top: 5px;">
      <span style="margin-left:8px;height:40px;line-height:40px;margin-left:35px">
        <span @click="toDashboard"
          style="font-weight:800;color:#fff;cursor: pointer;vertical-align: middle;font-size:24px;position: relative;top: -3px;margin-right:80px;"
          class="closeDrawer">
          <span style="font-size: 17px;position: relative;top: -3px;">多云管理平台</span>
        </span>
        <!-- <el-select v-if="$route.meta.projectNoget" v-model="projectCheck" multiple collapse-tags filterable
          @change="changeProject" :class="projectCheck.length > 1 ? 'hideoneacSel acSelColor' : 'acSelColor oneacSel'"
          @focus="focusHide" clearable style="margin-left: 4%;width:130px;top: -2px;" placeholder="请选择项目" size="small">
          <el-option v-for="(item, index) in projectIdArr" :key="index" :label="item.name" :value="item.id">
          </el-option>
        </el-select> -->
      </span>
      <div class="right-menu">
        <span style="color:#fff;float:right;">
          <el-menu class="el-menu-demo navbarMenu" mode="horizontal" @select="handleSelect">
            <template v-for="(item, index) in permissionRoutes[0].children">
              <el-submenu v-if="item.children[0]" :key="index" :index="'cp' + index">
                <template slot="title">{{ item.meta.pageTitle }}</template>
                <template v-for="(itemchild, i) in item.children">
                  <template v-for="(itemchild1, ii) in itemchild.children">
                    <el-menu-item v-if="!itemchild1.meta.hidden" :key="'cp' + index + i + ii"
                      :index="permissionRoutes[0].path + '/' + item.path + '/' + itemchild.path + '/' + itemchild1.path"><span
                        :class="urlHref == permissionRoutes[0].path + '/' + item.path + '/' + itemchild.path + '/' + itemchild1.path ? 'activeS' : ''">{{
                          itemchild1.meta.pageTitle }}</span></el-menu-item>
                  </template>
                </template>
              </el-submenu>
            </template>
          </el-menu>
        </span>
        <div class="avatar-container right-menu-item hover-effect">
          <div class="avatar-wrapper">
           
            <el-popover
              placement="bottom"
              width="90"
              trigger="hover"
              >
              <div style="text-align: center;">
                <el-link type="primary" style="font-size: 12px;" @click="logout">退出登录</el-link>
              </div>
              <span class="closeDrawer" slot="reference" style="color:#fff;font-size:12px;">
                <i class="el-icon-user" style="vertical-align: middle;"></i>
              {{ username }}
            </span>
            </el-popover>

          </div>
        </div>
      </div>
    </div>
    <el-dialog :visible.sync="modelStatus" :close-on-press-escape="false" :close-on-click-modal="false"
      :show-close="false" width="34%" class="sliderClass">
      <div class="wscn-http404">
        <h3 style="margin:0"><i class="el-icon-warning-outline" style="color:#409EFF;margin: 0 5px 0 0;"></i>提示</h3>
        <div class="bullshit" style="text-align:center">
          <div class="bullshit__headline" style="margin:10px 0 10px 0;">该用户账号下暂无任何项目,请
            <el-button @click="taskClick('FormJoinCloudProjectConfig', '加入项目')" type="primary"
              style="font-size:12px;">
              申请加入项目
            </el-button>
          </div>
        </div>
      </div>
    </el-dialog>
    <div v-if="reviewStatus222">
      <formtem-Dialog :nodeCheckName='nodeCheckName' :detailOrLook="detailOrLook" @addClose="addClose"
        :infoAddAttr="infoReviewAttr" :detailOrsigle='detailOrsigle' :infoAddForm="infoReviewForm"
        :reviewStatus='reviewStatus222' :newTaskName="newTaskName" :historyTableList='historyTableList'
        :historyparamreferenceShow='historyparamreferenceShow' :baseConfigData='baseConfigData'
        :childNodename='childNodename' :checkList='checkList' :FormTemplatjobMode='FormTemplatjobMode'></formtem-Dialog>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { userMixin } from '@/components/mixin/user'
import formtemDialog from "@/views/components/formTemp"
import Http from "@/components/api/services";
export default {
  components: {
    formtemDialog
  },
  mixins: [userMixin],
  data() {
    return {
      modelStatus: false,
      reviewStatus222: false,
      detailOrLook: false,
      nodeCheckName: "",
      infoReviewAttr: [],
      infoReviewForm: {},
      infoAddBFForm: {},
      historyparamreferenceShow: [],
      detailOrsigle: '',
      childNodename: '',
      newTaskName: '',
      historyTableList: {},
      baseConfigData: {},
      FormTemplatjobMode: '',
      username: "",
      projectIdArr: [],
      projectCheck: [],
      urlHref: "",
    }
  },
  watch: {
    '$route': 'routeChange'
  },
  computed: {
    ...mapGetters([
      "permissionRoutes", "loginUserName", "formPolicyDomain"
    ])
  },
  methods: {
    logout(){
      Http.logout().then(res => {
          this.$store.dispatch('user/resetUserName')
          this.$message({
            showClose: true,
            message: "退出成功",
            type: "success",
            duration: 3000
          }); 
          location.href = `${location.origin}/login`;
      })
    },
    addClose() {
      this.reviewStatus222 = false
    },
    taskClick(taskname, chinaName) {
      this.checkList = []
      this.childNodename = ""
      this.modelStatus = false
      //中文名称
      this.newTaskName = chinaName
      this.nodeCheckName = taskname
      this.reviewStatus222 = true
    },
    handleClose() {
      let nowTime = new Date().getTime()
      localStorage.setItem('rateTimeHCP', nowTime)
    },
    toDashboard() {
      this.$router.push({
        path: '/dashboard'
      })
    },
    routeChange() {
      this.urlHref = location.pathname
    },
    handleSelect(key, keyPath) {
      this.$router.push({
        path: key,
      });
    },
    doClosePop(refName, path1, path2, item) {
      var path = item.path
      if (!item.showOne) {
        path += `/${item.children[0].path}`

      }
      this.$router.push({
        path: path,
      });
    },
    changeProject() {
      if (this.projectCheck.length) {
        this.$EventBus.$emit("checkProjectArr", this.projectCheck);
        sessionStorage.setItem("checkProjectArr", this.projectCheck.join(","))
      } else {
        var arr = []
        this.projectIdArr.map((k, i) => {
          arr.push(k.id)
        })
        this.$EventBus.$emit("checkProjectArr", arr);
        sessionStorage.setItem("checkProjectArr", arr.join(","))
      }
    },
    focusHide() {
      this.$EventBus.$emit("closeDrawer", false);
    },
    openInfo(url, status) {
      window.open(url)
    }
  },
  mounted() {
    this.urlHref = location.pathname
    this.routeChange()
  },
  created() {
    sessionStorage.setItem("checkProjectArr", "")
    this.username = sessionStorage.getItem("username")
  },
}
</script>


<style lang="scss" scoped>
.navbar {
  height: 46px;
  overflow: hidden;
  position: relative;
  background: #fff;

  .logo {
    height: 20px;
    vertical-align: middle;
  }

  .logoPng {
    width: 30px;
    vertical-align: middle;
    margin-right: 4px;
  }

  background: #eee;

  .breadcrumb-container {
    float: left;
  }

  .errLog-container {
    display: inline-block;
    vertical-align: top;
  }

  .right-menu {
    float: right;
    height: 100%;
    line-height: 40px;
    display: flex;
    align-items: center;
    &:focus {
      outline: none;
    }
    .right-menu-item {
      display: inline-block;
      padding: 0 8px;
      height: 100%;
      font-size: 16px;
      color: #5a5e66;
      vertical-align: text-bottom;
      &.hover-effect {
        cursor: pointer;
        transition: background 0.3s;
        &:hover {
          background: rgba(0, 0, 0, 0.025);
        }
      }
    }

    .avatar-container {
      margin-right: 10px;
      .avatar-wrapper {
        position: relative;
        top: -3px;
        span {
          margin: 0 10px;
          position: relative;
          top: -1px;
        }
        .el-icon-caret-bottom {
          cursor: pointer;
          position: absolute;
          right: -10px;
          top: 20px;
          font-size: 12px;
        }
      }
    }
  }
}

.toDoBtn:hover {
  border-bottom: 1px solid;
}

.elPopBox {
  .navBoxMain {
    display: flex;
    text-align: center;
    cursor: pointer;
    .navTitle {
      color: #449bd0;
    }
  }
}

.navChildBoxMain {
  flex: 1;

  .navChilddiv {
    div {
      margin-bottom: 12px;
    }
  }

  .navChilddiv div:last-child {
    margin-bottom: 0;
  }
}

.navChildBoxMain:nth-child(2) {
  border-left: 1px solid #e7e7e7;
}

.navBoxMain h4 {
  margin: 0 0 10px 0;
  color: #1890ff;
}

.navBoxMain div {
  line-height: 22px;
}

.activeUi {
  color: #f60 !important;
}

.activeUiBorder {
  color: #f60 !important;
}

.bottomBorder {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 3px;
  background: rgba(24, 144, 255, 0.7) !important;
}

.uiWidth {
  padding: 0 7px;
  color: #000;
  vertical-align: baseline;
}
</style>
<style scoped>
.navbar /deep/ .el-dialog__body {
  padding: 17px 20px 30px !important;
}

.navbar /deep/ .el-select__tags span {
  max-width: 44px !important;
  overflow: hidden !important;
  white-space: nowrap !important;
  text-overflow: ellipsis !important;
  padding: 0 3px;
}

.el-textarea {
  top: 30px;
}

.title-age {
  display: inline-block;
  padding: 20px 0 10px;
  margin-right: 20px;
  cursor: pointer;
  border-bottom: 3px solid #409eff;
  color: #409eff;
}

.sliderClass /deep/ .sliderClass1 .el-slider__bar {
  background: #ff2400 !important;
}

.sliderClass /deep/ .sliderClass2 .el-slider__bar {
  background: #ff8c00 !important;
}

.sliderClass /deep/ .sliderClass3 {
  width: 97%;
}

.sliderClass /deep/ .sliderClass3 .el-slider__bar {
  background: #36bf36 !important;
}

.sliderClass /deep/ .sliderClass1 .el-slider__button {
  border-color: red !important;
}

.sliderClass /deep/ .sliderClass2 .el-slider__button {
  border-color: #ff8c00 !important;
}

.sliderClass /deep/ .sliderClass3 .el-slider__button {
  border-color: #36bf36 !important;
}

.sliderClass /deep/ .sliderClass3 .el-slider__marks-text {
  white-space: nowrap !important;
}

.sliderClass /deep/ .el-slider__marks-text {
  width: 20px;
  text-align: center;
}

.navbar /deep/ .hideoneacSel .el-select__tags span span.el-tag:first-child {
  max-width: 30px !important;
}

.navbar /deep/ .oneacSel .el-select__tags span span.el-tag:first-child {
  max-width: 70px !important;
}

.navbar /deep/.el-select__input {
  margin-left: 15px !important;
}

.navbar /deep/ .el-select__tags span.el-tag {
  position: relative;
}

.navbar /deep/ .el-select__tags span.el-tag:first-child {
  padding-right: 7px;
  padding-left: 3px;
}

.navbar /deep/ .el-select__tags span.el-tag i {
  position: absolute;
  right: 0;
  top: 5px;
}

.navbar /deep/ .acSelColor.el-select .el-input .el-select__caret {
  position: relative;
  top: -2px;
  font-size: 12px;
}

.navbar /deep/ .acSelColor.el-select .el-select__input {
  color: #fff !important;
}

.navbarMenu.el-menu.el-menu--horizontal {
  border-bottom: none;
}

.navbarMenu {
  position: relative;
  top: -3px;
}

.navbarMenu /deep/ .el-submenu .el-submenu__title {
  height: 40px !important;
  line-height: 38px !important;
  padding: 0 7px !important;
  color: #fff;
  background: none !important;
  position: relative;
  top: 2px;
  font-size: 12px !important;
}

.navbarMenu /deep/ .el-submenu.is-active .el-submenu__title {
  border-bottom: none !important;
}

.navbarMenu /deep/ i {
  color: #fff !important;
}

.navbarBoxMain /deep/ .el-dialog__header {
  padding: 0 8px;
  margin: 0;
}

.navbarBoxMain /deep/ .diaBox .el-dialog__header {
  padding: 20px 20px 10px;
}
</style>
<style>
::-webkit-scrollbar-track-piece {
  background-color: #f8f8f8;
}

::-webkit-scrollbar {
  width: 8px;

  height: 8px;
}

::-webkit-scrollbar-thumb {
  background-color: #ddd;
  background-clip: padding-box;
  min-height: 28px;
}

::-webkit-scrollbar-thumb:hover {
  background-color: #c5c5c5;
}

.el-menu::-webkit-scrollbar-track-piece {
  background: #cbdcf8;
}

.el-menu::-webkit-scrollbar-thumb {
  background-color: #f7f8fa;
  background-clip: padding-box;
  min-height: 28px;
}

.el-menu::-webkit-scrollbar-thumb:hover {
  background-color: #ddd;
}

.el-menu--horizontal .el-menu--popup {
  min-width: auto !important;
}

.el-menu--horizontal .el-menu-item.is-active {
  color: #909399 !important;
}

.el-menu--horizontal .el-menu-item:hover {
  color: rgba(24, 144, 255, 0.7) !important;
}

.hoverHint:hover {
  color: #409eff !important;
  border: 1px solid #409eff !important;
}

.activeS {
  color: #409eff;
}
</style>