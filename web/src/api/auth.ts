import request from "@/utils/request";

// 定义接口类型
export interface LoginParams {
  username: string;
  password: string;
  captchaId: string;
  captcha: string;
}

export interface LoginResult {
  token: string;
  userInfo: {
    id: number;
    username: string;
    name: string;
    roles: any[];
  };
}

export interface CaptchaResult {
  captchaId: string;
  captchaImg: string;
}

// 获取验证码
export const getCaptcha = () => {
  // request<T> 这里的 T 是返回的数据类型
  return request.post<any, CaptchaResult>("/auth/captcha");
};

// 登录
export const login = (data: LoginParams) => {
  return request.post<any, LoginResult>("/auth/login", data);
};

// 退出登录
export const logout = () => {
  return request.post<any, null>("/auth/logout");
};
