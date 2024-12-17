import axios from 'axios';

// 创建axios实例
const service = axios.create({
  timeout: 30000,
  withCredentials: true,
});

// 用于存储 Vue 实例
let vueInstance = null; 

// 设置 Vue 实例
export function setVueInstance(instance) {
  vueInstance = instance;
}

// 拦截请求
service.interceptors.request.use(
  config => {
    return config;
  },
  err => {
    console.log("请求错误:", err);
    if (vueInstance) {
      vueInstance.$message.error("请求错误");
    }
    return Promise.reject(err);
  },
);

// 拦截响应
service.interceptors.response.use(
  response => {
    // 怎么在这里拿到请求的url呢？
    if (response?.data?.code !== 0 && response?.data?.code != 200 && response.config.url!=="/api/v1/cloud/user/getUserInfo") {
      vueInstance.$message({
        showClose: true,
        message: response?.data?.message,
        type: 'error'
      });
    }
    if (response?.data?.code === 0 || response?.data?.code === 200) {
      return response;
    } else {
      return Promise.reject(response?.data);
    }
  },
  err => {
    const response = err.response;
    vueInstance.$message({
      showClose: true,
      type: 'error',
      message: "服务器请求错误",
    });
    return Promise.reject(err);
  },
);

export default service;
