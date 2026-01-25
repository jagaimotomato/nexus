import { useState, useEffect } from "react";
import { Button, Form, Input, message, Typography } from "antd"; // 移除 Row, Col，用 Flex 布局替代
import { useNavigate } from "react-router-dom";
import {
  UserOutlined,
  LockOutlined,
  SafetyCertificateOutlined,
  CodeOutlined,
} from "@ant-design/icons";
import useUserStore from "@/store/useUserStore";
import { getCaptcha, login, type LoginParams } from "@/api/auth";

const { Title, Text } = Typography;

const LoginPage = () => {
  const navigate = useNavigate();
  const userStoreLogin = useUserStore((state) => state.login);
  const setUserInfo = useUserStore((state) => state.setUserInfo);

  const [loading, setLoading] = useState(false);
  const [captchaData, setCaptchaData] = useState({ img: "", id: "" });

  const fetchCaptcha = async () => {
    try {
      const data = await getCaptcha();
      setCaptchaData({ id: data.captchaId, img: data.captchaImg });
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    fetchCaptcha();
  }, []);

  const onFinish = async (values: Omit<LoginParams, "captchaId">) => {
    setLoading(true);
    try {
      const res = await login({ ...values, captchaId: captchaData.id });
      message.success("登录成功");
      userStoreLogin(res.token);
      setUserInfo(res.userInfo);
      navigate("/dashboard", { replace: true });
    } catch (error) {
      console.error(error);
      fetchCaptcha();
    } finally {
      setLoading(false);
    }
  };

  return (
    // 容器：全屏 + 渐变背景 + 居中
    <div className="flex h-screen w-full items-center justify-center bg-gradient-to-br from-slate-100 to-slate-300 relative overflow-hidden">
      {/* 背景装饰球 (绝对定位 + 模糊) */}
      <div className="absolute top-1/4 left-1/4 h-72 w-72 rounded-full bg-blue-400 opacity-20 blur-[80px] animate-pulse"></div>
      <div className="absolute bottom-1/4 right-1/4 h-64 w-64 rounded-full bg-purple-400 opacity-20 blur-[80px]"></div>

      {/* 登录卡片：白色背景 + 毛玻璃 + 阴影 */}
      <div className="z-10 w-full max-w-[420px] rounded-xl bg-white/80 p-10 shadow-2xl backdrop-blur-md">
        {/* Header */}
        <div className="mb-8 text-center">
          <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-blue-500/10">
            <CodeOutlined className="text-3xl text-blue-600" />
          </div>
          <Title level={3} className="!mb-1 !mt-0 text-gray-800">
            Nexus Admin
          </Title>
          <Text type="secondary" className="text-gray-500">
            企业级后台管理系统
          </Text>
        </div>

        {/* Form */}
        <Form
          name="login"
          size="large"
          onFinish={onFinish}
          autoComplete="off"
          className="space-y-4" // Tailwind: 子元素垂直间距
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: "请输入用户名" }]}
            className="!mb-4"
          >
            <Input
              prefix={<UserOutlined className="text-gray-400" />}
              placeholder="账户: admin"
              className="hover:border-blue-400 focus:border-blue-500"
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: "请输入密码" }]}
            className="!mb-4"
          >
            <Input.Password
              prefix={<LockOutlined className="text-gray-400" />}
              placeholder="密码: 123456"
            />
          </Form.Item>

          <Form.Item className="!mb-6">
            <div className="flex gap-3">
              <Form.Item
                name="captcha"
                noStyle
                rules={[{ required: true, message: "验证码必填" }]}
              >
                <Input
                  prefix={
                    <SafetyCertificateOutlined className="text-gray-400" />
                  }
                  placeholder="验证码"
                  className="flex-1"
                />
              </Form.Item>

              {/* 验证码图片框 */}
              <div
                onClick={fetchCaptcha}
                className="flex h-[40px] w-32 cursor-pointer items-center justify-center overflow-hidden rounded border border-gray-300 bg-white hover:border-blue-400 transition-colors"
                title="点击刷新"
              >
                {captchaData.img ? (
                  <img
                    src={captchaData.img}
                    alt="code"
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <span className="text-xs text-gray-400">Loading</span>
                )}
              </div>
            </div>
          </Form.Item>

          <Form.Item className="!mb-0">
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              className="h-11 w-full rounded-lg bg-blue-600 text-base font-medium shadow-md hover:bg-blue-500"
            >
              登 录
            </Button>
          </Form.Item>
        </Form>
      </div>
    </div>
  );
};

export default LoginPage;
