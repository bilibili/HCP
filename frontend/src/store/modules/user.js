import Layout from "@/layout";
import Http from "@/components/api/services";
import Cookies from "js-cookie";

const getDefaultState = () => {
  return {
    token: "",
    name: "",
    avatar: "",
    routes: [],
    addRoutes: [],
    globalarr: [],
    loginUserName: "",
    loginUserID: "",
    loginUserInfo: {},
    formPolicyDomain: {},
    manageUser: [
      "admin",
      "root",
    ],
    manageUserGroup: [],
  };
};

const state = getDefaultState();

export const loadView = (view, obj) => {
  // 路由懒加载
  if (!view) {
    if (!obj) {
      return (resolve) => require([`@/views/components/productlist`], resolve);
    }
    return (resolve) => require([`@/layout`], resolve);
  }
  return (resolve) => require([`@/views/${view}`], resolve);
};

const routerList = [
  {
    path: "/general",
    component: Layout,
    icon: "el-icon-cpu",
    redirect: "/general/job/list/apply",
    meta: { title: "publicManage", icon: "index", pageTitle: "通用工具" },
    children: [
      {
        path: "job",
        name: "taskManagement",
        redirect: "/general/job/list/apply",
        component: () => import("@/views/publicManage"),
        meta: { title: "taskManagement", icon: "index", pageTitle: "操作管理" },
        children: [
          {
            path: "list",
            name: "taskmanageSetting",
            icon: "el-icon-setting",
            redirect: "/general/job/list/apply",
            component: () =>
              import("@/views/publicManage/taskManagement"),
            meta: { title: "manage", icon: "index", pageTitle: "操作管理" },
            children: [
              {
                path: "apply",
                name: "myapply",
                component: () =>
                  import("@/views/publicManage/taskManagement/job"),
                meta: {
                  title: "apply",
                  icon: "index",
                  pageTitle: "我的申请",
                  topTitle: "操作管理",
                  projectNoget: false,
                },
              },
              {
                path: "all",
                name: "myjob",
                component: () =>
                  import("@/views/publicManage/taskManagement/job"),
                meta: {
                  title: "job",
                  icon: "index",
                  pageTitle: "全部任务",
                  topTitle: "操作管理",
                  projectNoget: false,
                },
              },
            ],
          },
        ],
      },
      {
        path: "manage",
        name: "sysManage",
        redirect: "/general/manage/project/list",
        component: () => import("@/views/publicManage"),
        meta: { title: "sysManage", icon: "index", pageTitle: "系统管理" },
        children: [
          {
            path: "project",
            name: "sysmanage",
            icon: "el-icon-film",
            redirect: "/general/manage/project/list",
            component: () => import("@/views/publicManage"),
            meta: {
              title: "manage",
              icon: "index",
              pageTitle: "项目管理",
              topTitle: "系统管理",
              projectNoget: false,
            },
            children: [
              {
                path: "list",
                name: "projectManage",
                component: () => import("@/views/components/productlist"),
                meta: {
                  title: "project",
                  icon: "index",
                  pageTitle: "项目",
                  topTitle: "系统管理",
                  projectNoget: false,
                },
              },
              {
                path: "account",
                name: "accountManage",
                component: () => import("@/views/components/productlist"),
                meta: {
                  title: "account",
                  icon: "index",
                  pageTitle: "账号配置",
                  topTitle: "系统管理",
                  projectNoget: false,
                },
              },
            ],
          },
          {
            path: "account",
            name: "sysUsermanage",
            icon: "el-icon-bank-card",
            redirect: "/general/manage/account/user",
            component: () => import("@/views/components"),
            meta: { title: "manage", icon: "index", pageTitle: "公有云管理" },
            children: [
              {
                path: "CloudProvider",
                name: "CloudProvider",
                component: () => import("@/views/components/productlist"),
                meta: {
                  title: "project",
                  icon: "index",
                  pageTitle: "云厂商",
                  topTitle: "公有云管理",
                  projectNoget: false,
                },
              },
              {
                path: "list",
                name: "Useraccount",
                component: () => import("@/views/components/productlist"),
                meta: {
                  title: "project",
                  icon: "index",
                  pageTitle: "云账号",
                  topTitle: "公有云管理",
                  projectNoget: false,
                },
              },
              {
                path: "user",
                name: "User",
                component: () => import("@/views/components/productlist"),
                meta: {
                  title: "project",
                  icon: "index",
                  pageTitle: "用户",
                  topTitle: "公有云管理",
                  projectNoget: false,
                },
              },
            ],
          },
        ],
      },
    ],
  },
  {
    path: "/compute",
    component: Layout,
    icon: "el-icon-cpu",
    redirect: "/compute/server/instance/list",
    meta: { title: "index", icon: "index", pageTitle: "计算" },
    children: [
      {
        path: "server",
        name: "host",
        redirect: "/compute/server/instance/list",
        component: () => import("@/views/components"),
        meta: {
          title: "index",
          icon: "index",
          pageTitle: "云服务器",
          remark: "弹性可伸缩的计算服务",
        },
        children: [
          {
            path: "instance",
            name: "hostManage",
            icon: "el-icon-cpu",
            redirect: "/compute/server/instance/list",
            component: () => import("@/views/components"),
            meta: { title: "index", icon: "index", pageTitle: "计算" },
            children: [
              {
                path: "list",
                name: "cloudServer",
                component: () => import("@/views/components/productlist"),
                meta: {
                  title: "general",
                  icon: "index",
                  pageTitle: "云服务器",
                  topTitle: "计算",
                  projectNoget: true,
                },
              },
            ],
          },
          {
            path: "image",
            name: "serverJX",
            icon: "el-icon-money",
            redirect: "/compute/host/image/list",
            component: () => import("@/views/components"),
            meta: { title: "index", icon: "index", pageTitle: "镜像" },
            children: [
              {
                path: "list",
                name: "serverImage",
                component: () => import("@/views/components/productlist"),
                meta: {
                  title: "general",
                  icon: "index",
                  pageTitle: "系统镜像",
                  topTitle: "镜像",
                  projectNoget: false,
                },
              },
            ],
          },
          {
            path: "network",
            name: "servernetwork",
            icon: "el-icon-connection",
            redirect: "/compute/host/network/VPC",
            component: () => import("@/views/components"),
            meta: { title: "index", icon: "index", pageTitle: "网络" },
            children: [
              {
                path: "vpc",
                name: "VPC",
                component: () => import("@/views/components/productlist"),
                meta: {
                  title: "general",
                  icon: "index",
                  pageTitle: "VPC",
                  topTitle: "网络",
                  projectNoget: false,
                },
              },
              {
                path: "subnet",
                name: "subNet",
                component: () => import("@/views/components/productlist"),
                meta: {
                  title: "general",
                  icon: "index",
                  pageTitle: "子网",
                  topTitle: "网络",
                  projectNoget: false,
                },
              },
              {
                path: "securityGroup",
                name: "securityGrouplist",
                component: () => import("@/views/components/productlist"),
                meta: {
                  title: "general",
                  icon: "index",
                  pageTitle: "安全组",
                  topTitle: "网络",
                  projectNoget: true,
                },
              },
            ],
          }
        ],
      },
    ],
  },
];

const mutations = {
  RESET_STATE: (state) => {
    Object.assign(state, getDefaultState());
  },
  SET_TOKEN: (state, token) => {
    state.token = token;
  },
  SET_NAME: (state, name) => {
    state.name = name;
  },
  SET_AVATAR: (state, avatar) => {
    state.avatar = avatar;
  },
  SET_ROUTES: (state, routes) => {
    state.routes = routes;
    state.addRoutes = routes;
  },
  SET_GLOBALARR: (state, arr) => {
    state.globalarr = arr;
  },
  SET_USER: (state, str) => {
    state.loginUserName = str;
  },
  SET_USER_ID: (state, str) => {
    state.loginUserID = str;
  },
  SET_USER_INFO: (state, str) => {
    state.loginUserInfo = str;
  },
  SET_FORMPOLICYDOMAIN: (state, str) => {
    state.formPolicyDomain = str;
  },
  SET_MANAGEUSER: (state, str) => {
    state.manageUser = str;
  },
  SET_MANAGEUSERGROUP: (state, str) => {
    state.manageUserGroup = str;
  },
};

const actions = {
  generateRoutes({ commit }, roles) {
    return Http.getUserAuth().then((response) => {
        console.log(response);
        var username = response.data.data.Name;
        sessionStorage.setItem("username", username);
        sessionStorage.setItem("isCollapse", 0);
        Cookies.set("username", username);
        commit("SET_USER", username);
        commit("SET_USER_ID", response.data.data.ID);
        commit("SET_USER_INFO", response.data.data);
        if(response.data.data.Role=='2'){
          routerList[0].children[1].children[1].children.splice(2,1)
        }
        if(response.data.data.Role=='3'){
          routerList[0].children[1].children.splice(1,1)
        }
        commit("SET_ROUTES", routerList);
      })
      .catch((err) => {
        commit("SET_ROUTES", [
          {
            path: "/",
            component: Layout,
            redirect: "/dashboard",
            name: "DashboardHide",
            hidden: true,
            meta: { hidden: true },
            children: [],
          },
        ]);
      });
  },
  resetUserName({ commit }){
    commit('SET_USER_ID', "")
    commit('SET_USER', "")
    commit('SET_USER_INFO', "")
  },
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
};
