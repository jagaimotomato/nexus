import axios, {
  type AxiosInstance,
  type AxiosResponse,
  type InternalAxiosRequestConfig,
} from "axios";
import { message } from "antd";
import useUserStore from "@/store/useUserStore";

interface Result<T = any> {
  code: number;
  msg: string;
  data: T;
}

console.log(import.meta.env.VITE_API_URL);

const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 10000,
  headers: { "Content-Type": "application/json" },
});

service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = useUserStore.getState().token;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

service.interceptors.response.use(
  (response: AxiosResponse<Result>) => {
    const res = response.data;
    if (res.code === 0) {
      return res.data;
    }
    message.error(res.msg || "系统错误");

    if (res.code === 401) {
      useUserStore.getState().logout();
      window.location.href = "/login";
    }
    return Promise.reject(new Error(res.msg || "系统错误"));
  },
  (error) => {
    // 处理 HTTP 网络错误 (如 404, 500)
    console.error("Request Error:", error);
    let msg = "网络连接异常";
    if (error.response) {
      switch (error.response.status) {
        case 401:
          msg = "未授权，请登录";
          break;
        case 403:
          msg = "拒绝访问";
          break;
        case 404:
          msg = "请求资源不存在";
          break;
        case 500:
          msg = "服务器内部错误";
          break;
        default:
          msg = error.message;
      }
    }
    message.error(msg);
    return Promise.reject(error);
  }
);

export default service;
