import Http from "./request";
// 获取签名
function getSign() {
  return Http.post(`/api/v1/cloud/user/getSign`);
}
// 用户登录
function login(data) {
  return Http.post(`/api/v1/cloud/login`,data);
}
// 退出登录
function logout() {
  return Http.post(`/api/v1/cloud/logout`);
}
// 获取用户权限
function getUserAuth() {
  return Http.post(`/api/v1/cloud/user/getUserInfo`);
}
// 获取仪表盘数据
function getDashCount() {
  return Http.post(`/api/v1/cloud/user/overview`);
}
// 获取模型列表
function getQueryList(data) {
  return Http.post(`/api/v1/cloud/db/query`, data);
}
// 获取表单模板
function getFormTemplate(data) {
  return Http.post(`/api/v1/cloud/db/query/formTemplate`, data);
}
// 创建任务
function createJob(data) {
  return Http.post(`/api/v1/cloud/job/create`, data);
}
// 获取任务列表
function getJobList(data) {
  return Http.post(`/api/v1/cloud/db/query`, data);
}
// 获取任务详情
function getJobDetail(data) {
  return Http.post(`/api/v1/cloud/job/get`, data);
}
export default {
  getSign,
  login,
  logout,
  getUserAuth,
  getDashCount,
  getQueryList,
  getFormTemplate,
  createJob,
  getJobList,
  getJobDetail
};
