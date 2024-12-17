<template>
  <section class="app-main">
    <div @click="toDashboard" class="menuActive2">
      <span style="position: relative;top: 1px;font-size:12px">
        <span style="margin-right:5px">总览</span>
        <span>|</span>
      </span>
    </div>
    <div style="display:flex;position: relative;">
      <div>
        <el-menu class="el-menu-vertical-demo" text-color="#000" ref="cmenunav"
          active-text-color="rgba(24, 144, 255, 0.7)">
          <el-menu-item @click="todraw" index="openNav" :class="drawer ? 'menuActive menuActive1' : 'menuActive1'">
            <div>
              <span slot="title" style="position: relative;top: 0px;font-size:12px;" class="openDrawer">产品
                <i :class="drawer ? 'el-icon-arrow-up openDrawer' : 'el-icon-arrow-down openDrawer'"
                  style="font-size: 12px;position: absolute;top: 1px;"></i>
              </span>
            </div>
          </el-menu-item>
        </el-menu>
      </div>
      <div v-if="hideStatus" :class="isCollapse ? 'navBox1' : 'navBox2'">
        <template>
          <el-menu
            :style="isCollapse ? 'height: calc(100vh - 46px);overflow-y: auto;overflow-x: hidden;width:60px;background:#cbdcf8;' : 'height: calc(100vh - 46px);width:140px;background:#cbdcf8;overflow-y: auto;overflow-x: hidden;'"
            :default-active="activeNav" :default-openeds="activeopeneds" class="el-menu-vertical-demo"
            :collapse="isCollapse" text-color="#000" ref="menunav" active-text-color="rgba(24, 144, 255, 0.7)">
            <template v-for="(item, index) in childRouterArr">
              <template v-if="item.children">
                <el-submenu :index="parentUrl + '/' + item.path" :key="index">
                  <template slot="title">
                    <i :class="item.icon"></i>
                    <span>{{ item.meta.pageTitle }}</span>
                  </template>

                  <template v-for="(childItem, childIndex) in item.children">
                    <el-menu-item v-if="!childItem.meta.hidden" @click="linkRouter(item, childItem)" :key="childIndex"
                      :index="parentUrl + '/' + item.path + '/' + childItem.path">
                      <span slot="title" class="closeDrawer">{{ childItem.meta.pageTitle }}</span>
                    </el-menu-item>
                  </template>
                </el-submenu>
              </template>
              <template v-else>
                <el-menu-item v-if="!item.meta.hidden" @click="linkRouter(item)" :key="index"
                  :index="parentUrl + '/' + item.path">
                  <i :class="item.icon"></i>
                  <span slot="title" class="closeDrawer">{{ item.meta.pageTitle }}</span>
                </el-menu-item>
              </template>
            </template>
            <template v-if="locationUrl.indexOf('/general/job') == -1">
              <template v-for="(item, index) in pRouterArr">
                <el-menu-item v-if="!item.meta.hidden" @click="czBtn(item)" :key="index + 'p'"
                  :index="pUrl + '/' + item.path">
                  <i :class="item.icon"></i>
                  <span slot="title" class="closeDrawer czItem">{{ item.meta.pageTitle }}</span>
                </el-menu-item>
              </template>
            </template>
          </el-menu>
        </template>
        <i @click="openIsCollapse"
          :class="isCollapse ? 'el-icon-s-unfold isCollapse1' : 'el-icon-s-fold isCollapse2'"></i>
      </div>
      <transition name="fade-transform" mode="in-out">
        <div
          :style="$route.name != 'DashboardActive' ? (isCollapse ? 'flex-grow: 1;width: calc(100% - 60px);position: absolute;background: #fff;height: 100vh;right: 0px;z-index:9999;overflow: auto' : 'flex-grow: 1;width: calc(100% - 140px);position: absolute;background: #fff;height: 100vh;right: 0px;z-index:9999;overflow: auto') : 'flex-grow: 1;width: calc(100% - 0px);position: absolute;background: #fff;height: 100vh;right: 0px;z-index:9999;overflow: auto'"
          v-if="drawer" ref="drawerDialog">
          <div style="padding:20px 0 20px 20px;">
            <div style="display:flex;height:calc(100vh - 90px)">
              <div style="width:140px;min-width:140px;height:100%;border-right:1px solid #ccc;padding: 0px 15px 0;">
                <template v-for="(item, index) in permissionRoutes">
                  <div v-if="index == 0 && !item.hidden" :key="index">
                    <div style="text-indent: 7px;">
                      <div
                        style="line-height:30px;border-bottom: 1px solid #ccc;font-size:16px;margin-bottom:10px;font-weight:600;color:#666;">
                        {{ item.meta.pageTitle }}</div>
                      <template v-for="(childItem, childIndex) in item.children">
                        <div
                          :class="parentUrl == item.path + '/' + childItem.path ? 'navItem navItemActive' : 'navItem'"
                          style="line-height:36px;font-size:14px;cursor: pointer;"
                          @click="linkRouterNav(item, childItem)" :key="childIndex">{{ childItem.meta.pageTitle }}</div>
                      </template>
                    </div>
                  </div>
                </template>
              </div>
              <div style="flex-grow:1;height:100%;overflow: auto;">
                <div style="padding:0 0 0 20px;">
                  <template v-for="(item, index) in permissionRoutes">
                    <div v-if="index != 0" :key="index">
                      <div v-if="!item.meta.hidden">
                        <div style="text-indent: 7px;overflow: hidden;">
                          <div
                            style="line-height:26px;font-size:16px;margin-bottom:12px;font-weight:600;color:#666;border-left:5px solid #409EFF;text-indent: 20px;">
                            {{ item.meta.pageTitle }}</div>
                          <template v-for="(childItem, childIndex) in item.children">
                            <el-card shadow="hover"
                              :class="!childItem.children || !childItem.children ? 'cardItem' : 'cardItem cardItem1'"
                              style="line-height:26px;font-size:14px;cursor: pointer;min-width:250px;float:left;margin:0 12px 15px 0"
                              :key="childIndex">
                              <div style="padding:10px;position: relative;font-weight:600;"
                                @click="linkRouterNav(item, childItem)"
                                :class="parentUrl == item.path + '/' + childItem.path ? 'navItem navItemActive' : 'navItem'">
                                <div>{{ childItem.meta.pageTitle }}</div>
                                <div v-if="childItem.meta.remark" style="font-size:13px;color:#999;">
                                  {{ childItem.meta.remark }}</div>
                                <div v-else style="font-size:12px;color:#999;">{{ childItem.meta.pageTitle }}介绍内容描述！
                                </div>
                                <div v-if="!childItem.children || !childItem.children"
                                  style="position: absolute;width: 100%;height: 100%;right: 0;top: 0;background: #fff;opacity: 0.6;display: flex;justify-content: center;align-items: center;">
                                </div>
                              </div>
                            </el-card>
                          </template>
                        </div>
                      </div>
                      <div v-else :key="index">
                        <div style="text-indent: 7px;overflow: hidden;">
                          <div
                            style="line-height:26px;font-size:16px;margin-bottom:12px;font-weight:600;color:#666;border-left:5px solid #409EFF;text-indent: 20px;">
                            {{ routesArr[index].meta.pageTitle }}</div>
                          <template v-for="(childItem, childIndex) in routesArr[index].children">
                            <el-card shadow="hover" :class="!childItem.redirect ? 'cardItem' : 'cardItem cardItem1'"
                              style="line-height:26px;font-size:14px;cursor: pointer;min-width:250px;float:left;margin:0 12px 15px 0"
                              :key="childIndex">
                              <div style="padding:10px;position: relative;font-weight:600;"
                                @click="linkRouterNav1(childItem)"
                                :class="locationUrl == childItem.redirect ? 'navItem navItemActive' : 'navItem'">
                                <div>{{ childItem.meta.pageTitle }}</div>
                                <div v-if="childItem.meta.remark" style="font-size:13px;color:#999;">
                                  {{ childItem.meta.remark }}</div>
                                <div v-else style="font-size:12px;color:#999;">{{ childItem.meta.pageTitle }}介绍内容描述！
                                </div>
                                <div v-if="!childItem.redirect"
                                  style="position: absolute;width: 100%;height: 100%;right: 0;top: 0;background: #fff;opacity: 0.6;display: flex;justify-content: center;align-items: center;">
                                </div>
                              </div>
                            </el-card>
                          </template>
                        </div>
                      </div>
                    </div>
                  </template>
                </div>
              </div>
            </div>
          </div>
        </div>
      </transition>
      <transition name="fade-transform" mode="out-in">
        <div class="boxMainShow"
          :style="isCollapse ? 'background:#fff;display:flex;flex-grow:1;width: calc(100vw - 80px);' : 'background:#fff;display:flex;flex-grow:1;width: calc(100vw - 170px);'">
          <router-view :key="key" />
        </div>
      </transition>
    </div>
  </section>
</template>

<script>
import { mapGetters } from 'vuex'
export default {
  name: 'AppMain',
  computed: {
    cachedViews() {
      return this.$store.state.tagsView.cachedViews
    },
    key() {
      return this.$route.path
    },
    title() {
      return this.$route.meta.pageTitle
    },
    ...mapGetters([
      "permissionRoutes"
    ])
  },
  data() {
    return {
      disRemark: '',
      slidervalue: 8,
      marks: {
        0: '0',
        1: '1',
        2: '2',
        3: '3',
        4: '4',
        5: '5',
        6: '6',
        7: '7',
        8: '8',
        9: '9',
        10: '10',
      },
      childRouterArr: [],
      activeopeneds: [],
      urlPath: "",
      pUrl: "",
      pRouterArr: [],
      activeNav: "",
      parentUrl: "",
      isCollapse: false,
      drawer: false,
      workStatus: false,
      locationUrl: "",
      routesArr: [],
      hideStatus: false
    }
  },
  watch: {
    '$route': 'routeChange',
    drawer() {
      if (!this.drawer) {
        this.$nextTick(() => {
          if (this.$refs.menunav) {
            this.$refs.menunav.activeIndex = this.activeNav
          }
        })
      }
    }
  },
  methods: {
    changeSlidervalue() {
      var colorStr = this.slidervalue < 6 ? '#ff2400' : (this.slidervalue < 8 ? '#ff8c00' : '#36bf36')
      for (var i = 0; i <= 10; i++) {
        if (i <= this.slidervalue) {
          this.marks[i] = {
            style: {
              color: colorStr,
            },
            label: `${i}`
          }
        } else {
          this.marks[i] = `${i}`
        }
      }
    },
    handleClose() {
      let nowTime = new Date().getTime()
      localStorage.setItem('rateTimeHCP', nowTime)
    },
    routeChange(obj) {
      this.hideStatus = false
      if (this.$route.name != 'DashboardActive') {
        this.hideStatus = true
      }
      this.locationUrl = location.pathname
      if (this.$route.matched.length > 2) {
        this.parentUrl = ""
        var arr = this.$route.matched[1].path.split("/").filter(item => item)
        this.permissionRoutes.map((k, i) => {
          if (k.path == '/' + arr[0]) {
            this.parentUrl += k.path
            k.children.map((k1, i1) => {
              if (k1.path == arr[1]) {
                this.parentUrl += "/" + k1.path
                this.childRouterArr = JSON.parse(JSON.stringify(k1.children))
              }
            })
          }
        })
      }
      this.activeNav = this.$route.matched[this.$route.matched.length - 1].path
      this.$nextTick(() => {
        if (this.$refs.menunav) {
          this.$refs.menunav.activeIndex = this.activeNav
        }
        this.activeopeneds.push(this.$route.matched[this.$route.matched.length - 2].path)
      })
    },
    linkRouter(obj, obj1, status) {
      this.workStatus = status
      if (status) {
        var parentUrl = this.pUrl
      } else {
        var parentUrl = this.parentUrl
      }
      if (obj1) {
        this.activeNav = parentUrl + "/" + obj.path + "/" + obj1.path
        this.$router.push({
          path: parentUrl + "/" + obj.path + "/" + obj1.path,
        });
      } else {
        this.activeNav = parentUrl + "/" + obj.path
        this.$router.push({
          path: parentUrl + "/" + obj.path,
        });
      }
    },
    pushRouter() {
      this.$router.push({
        path: "/work/index",
      });
    },
    todraw() {
      this.drawer = !this.drawer
    },
    toDashboard() {
      this.$router.push({
        path: '/dashboard'
      })
    },
    openIsCollapse() {
      this.isCollapse = !this.isCollapse
      this.$EventBus.$emit("isCollapse", this.isCollapse ? 1 : 0);
      sessionStorage.setItem("isCollapse", this.isCollapse ? 1 : 0)
    },
    linkRouterNav(obj, obj1) {
      if (!obj1.children || !obj1.children.length) {
        this.$message({
          showClose: true,
          message: "正在努力开发中，敬请期待……",
          type: 'info'
        });
        return;
      }
      this.$router.push({
        path: obj1.redirect,
      });
    },
    linkRouterNav1(obj1) {
      if (!obj1.redirect) {
        this.$message({
          showClose: true,
          message: "正在努力开发中，敬请期待……",
          type: 'info'
        });
        return;
      }
      this.$router.push({
        path: obj1.redirect,
      });
    },
    czBtn(obj) {
      this.$router.push({
        path: obj.redirect,
      });
    }
  },
  mounted() {
    this.routeChange(1)
    if (this.$route.name == 'DashboardActive') {
      this.drawer = false
    }
    this.$EventBus.$on("closeDrawer", (value) => {
      this.drawer = value
    });
    this.pUrl = this.permissionRoutes[0].path
    this.permissionRoutes[0].children.map((k, i) => {
      if (k.path == 'job') {
        this.pUrl += "/" + k.path
        this.pRouterArr = k.children
      }
    })
    sessionStorage.setItem("isCollapse", this.isCollapse ? 1 : 0)
  },
  created() {
    this.locationUrl = location.pathname
    document.addEventListener("click", (e) => {
      if (this.$refs.drawerDialog) {
        let aa = this.$refs.drawerDialog.contains(e.target)
        if (e.target.nodeName == "LI") {
          if (e.target.firstChild.nodeName == "DIV") {
            if (e.target.firstChild.firstChild.className.indexOf("openNav")) {
              this.drawer = true
              return
            }
          }
        }
        if (e.target.className.indexOf("openDrawer") != -1) {
          return;
        }
        if (aa) {
          this.drawer = false
          return
        }
        if (e.target.className.indexOf("closeDrawer") != -1) {
          this.drawer = false
          return
        }
        if (e.target.className.indexOf("el-menu") != -1) {
          this.drawer = false
          return
        }
        if (Object.keys(JSON.stringify(e.target.firstChild)).length == 0) {
          this.drawer = false
          return
        }


        if (e.target.className.indexOf("el-icon-arrow-down") != -1) {
          this.drawer = false
          return
        }

        if (!e.target.firstChild.className) {
          this.drawer = false
          return
        }

        if (e.target.firstChild) {
          if (e.target.firstChild.className.indexOf("openDrawer") == -1) {
            if (!aa) {
              this.drawer = false
            }
          }
        } else if (e.target.className.indexOf("openDrawer") == -1) {
          if (!aa) {
            this.drawer = false
          }
        }
      }
    })
  },
}
</script>

<style lang="scss" scoped>
.app-main {
  min-height: calc(100vh - 46px);
  width: 100%;
  position: relative;
  overflow: hidden;
  background: rgb(203, 220, 248);
}

.fixed-header + .app-main {
  padding-top: 53px;
}

.hasTagsView {
  .app-main {
    min-height: calc(100vh - 84px);
  }

  .fixed-header + .app-main {
    padding-top: 94px;
  }
}
</style>

<style lang="scss">
.el-popup-parent--hidden {
  .fixed-header {
    padding-right: 15px;
  }
}

.activeUiNav {
  color: red;
}
</style>
<style scoped>
.navBarChild {
  transition: left 0.25s;
}

.navBarChild /deep/ .el-drawer__body {
  overflow: auto;
}

.el-submenu .el-menu-item {
  min-width: 140px;
}

.el-submenu .el-menu-item span {
  color: #666;
}

.el-submenu .el-menu-item.is-active span {
  color: rgba(24, 144, 255, 0.7) !important;
}

.cardIActive.el-card {
  border: 1px solid rgba(64, 158, 255, 0.6);
}

.cardItem1 {
  background-image: url("../../assets/img/recombg.jpg");
  background-repeat: no-repeat;
  background-size: cover;
}

.navItemActive {
  color: #409eff;
}

.cardItem /deep/ .el-card__body {
  padding: 0;
}

.navBox1 {
  position: relative;
  animation: navBox1 0.15s;
}

.isCollapse1 {
  position: absolute;
  right: 15px;
  bottom: 7px;
  line-height: 24px;
  color: #909399;
  font-size: 22px;
  cursor: pointer;
}

.isCollapse2 {
  position: absolute;
  left: 20px;
  bottom: 7px;
  line-height: 24px;
  color: #909399;
  font-size: 22px;
  cursor: pointer;
}

@keyframes navBox1 {
  0% {
    min-width: 140px;
  }

  25% {
    min-width: 120px;
  }

  50% {
    min-width: 100px;
  }

  75% {
    min-width: 80px;
  }

  100% {
    min-width: 60px;
  }
}
</style>
<style>
.menuActive1 {
  position: fixed !important;
  left: 210px !important;
  top: 0px !important;
  color: #fff !important;
  z-index: 9999 !important;
  height: 45px !important;
  line-height: 45px !important;
  padding: 0px !important;
  cursor: pointer !important;
  font-size: 12px !important;
}

.menuActive2 {
  position: fixed !important;
  left: 170px !important;
  top: 0px !important;
  color: #fff !important;
  z-index: 9999 !important;
  height: 45px !important;
  line-height: 45px !important;
  padding: 0px !important;
  cursor: pointer !important;
  font-size: 12px !important;
}

::-webkit-scrollbar-track-piece {
  background: none !important;
}

.navBox2 /deep/ .menuActive1 i {
  color: #fff;
  font-size: 12px;
  top: 1px;
  position: relative;
}

.navBox2 /deep/ .menuActive2 i {
  color: #fff;
  font-size: 12px;
  top: 1px;
  position: relative;
}

.menuActive2:hover {
  background: none !important;
}

.menuActive1:hover {
  background: none !important;
}

.menuActive1.is-active {
  background: none !important;
}

.menuActive2.is-active {
  background: none !important;
}

.menuActive2:visited {
  background: none !important;
}

.menuActive {
  color: rgba(24, 144, 255, 0.7) !important;
}

.navBox2 /deep/ .menuActive.menuActive1 i {
  color: rgba(24, 144, 255, 0.7) !important;
}

.el-submenu__title {
  padding: 0 0 0 20px !important;
}

.el-submenu__icon-arrow {
  right: 5px !important;
}

.czItem {
  font-size: 14px !important;
}
</style>
<style scoped>
.rateClass {
  position: fixed;
  left: 50px;
  bottom: 11px;
  color: #000;
  font-size: 12px;
  z-index: 11111111;
  cursor: pointer;
}

.sliderClass /deep/ .sliderClass3 {
  width: 97%;
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

.el-textarea {
  top: 30px;
}

.sliderClass /deep/ .el-dialog__header {
  border-bottom: 3px dashed #eee;
  padding: 0 8px;
  margin: 0;
}

.navbar /deep/ .el-dialog__body {
  padding: 17px 20px 30px !important;
}
</style>