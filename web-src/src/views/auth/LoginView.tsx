"use client";

import { useEffect, useState } from "react";
import { Link, useNavigate, useLocation } from "react-router-dom";
import { t } from "i18next";
import { CircleUserRound, ShieldCheck } from "lucide-react";
import { getCaptchaConfig, getOidcConfig, getUser, login } from "@/api/user.api";
import { OidcConfig } from "@/api/user.models";
import AIOCaptchaWidget from "@/components/captcha/AIOCaptcha";
import { useDevice } from "@/contexts/DeviceContext";
import { CaptchaProps, CaptchaProvider } from "@/types/captcha";

export default function LoginView() {
  const { isMobile } = useDevice();
  // 替换 Next.js 路由为 React Router
  const navigate = useNavigate();
  const location = useLocation();

  // 使用 URLSearchParams 解析查询参数
  const searchParams = new URLSearchParams(location.search);
  const redirectUrl = searchParams.get("redirect") || "/";

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [captchaToken, setCaptchaToken] = useState("");
  const [captchaProps, setCaptchaProps] = useState<CaptchaProps | null>(null);
  const [captchaKey, setCaptchaKey] = useState(Date.now());
  const [oidcConfigs, setOidcConfigs] = useState<OidcConfig[]>([]);
  const [remember, setRemember] = useState(false);

  useEffect(() => {
    getUser()
      .then(() => {
        navigate(redirectUrl);
      })
      .catch((error) => {
        console.error("Error fetching user data:", error);
      });
    getCaptchaConfig()
      .then((response) => {
        const config = response.data;
        setCaptchaProps({
          provider: config.provider as CaptchaProvider,
          siteKey: config.siteKey || "",
          url: config.url || "",
          onSuccess: (token: string) => setCaptchaToken(token),
          onError: () => setError("login.captcha.failed"),
        });
      })
      .catch(() => setError("login.captcha.fetchFailed"));
  }, [redirectUrl, navigate]);
  console.log("captcha Token:", captchaToken);
  useEffect(() => {
    getOidcConfig()
      .then((response) => {
        const configs = response.data.oidcConfigs || [];
        setOidcConfigs(configs);
      })
      .catch((error) => {
        console.error("获取 OIDC 配置失败:", error);
        setError("login.oidc.fetchFailed");
      });
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    login({ username, password, captchaToken, remember })
      .then(() => {
        // 使用 navigate 替代 router.push
        navigate(redirectUrl);
      })
      .catch((err) => {
        setError(t("login.failed") + ": " + err.response?.data?.message || err.message);
        console.error("登录失败:", err);
        setCaptchaToken("");
        setCaptchaKey(Date.now());
      });
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-100 dark:bg-gray-900">
      <form
        onSubmit={handleSubmit}
        className={
          isMobile
            ? "w-full min-h-screen flex flex-col justify-center px-10 bg-white dark:bg-gray-800"
            : "w-[400px] bg-white dark:bg-gray-800 rounded-3xl shadow-2xl px-8 py-8 flex flex-col justify-center"
        }
      >
        <h2 className="text-3xl font-bold text-center text-gray-900 dark:text-white mb-6 tracking-tight">
          {t("login.login")}
        </h2>
        {/* 账号密码输入框 */}
        <div className="space-y-4 mb-4">
          {/* 用户名输入框 */}
          <div className="relative w-full">
            <span className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">
              <CircleUserRound className="w-5 h-5" />
            </span>
            <input
              type="text"
              placeholder={t("login.usernameOrEmail")}
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full pl-10 pr-4 py-2 rounded-3xl border border-gray-300 dark:border-gray-700 bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          {/* 密码输入框 */}
          <div className="relative w-full">
            <span className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">
              <ShieldCheck className="w-5 h-5" />
            </span>
            <input
              type="password"
              placeholder={t("login.password")}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full pl-10 pr-4 py-2 rounded-3xl border border-gray-300 dark:border-gray-700 bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
        {/* 记住设备和忘记密码 */}
        <div className="flex items-center justify-between mb-2 px-1 text-sm">
          <label className="flex items-center gap-1 text-gray-500 dark:text-gray-400">
            <input
              type="checkbox"
              className="accent-blue-600"
              checked={remember}
              onChange={(e) => setRemember(e.target.checked)}
            />
            {t("login.remember")}
          </label>
          <Link to="/forgot-password" className="text-blue-600 hover:underline">
            {t("login.forgotPassword")}
          </Link>
        </div>
        {/* Captcha组件 */}
        <div
          className={`flex justify-center (${![CaptchaProvider.DISABLE, CaptchaProvider.RECAPTCHA].includes(captchaProps?.provider ?? CaptchaProvider.DISABLE) ? "my-4" : ""})`}
        >
          {captchaProps && <AIOCaptchaWidget
            key={captchaKey}
            {...(captchaProps || {
              provider: CaptchaProvider.DISABLE,
              siteKey: "",
              url: "",
              onSuccess: (token: string) => setCaptchaToken(token),
              onError: () => {
                setError("login.captcha.failed");
                setCaptchaKey(Date.now());
              },
            })}
          />}
        </div>
        {error && <div className="text-red-500 text-sm text-center mb-2">{t(error)}</div>}
        {/* 登录按钮 */}
        <button
          type="submit"
          disabled={!captchaToken}
          className={
            "w-full py-2 rounded-3xl font-semibold text-lg transition-colors mt-2 " +
            (captchaToken
              ? "bg-blue-600 hover:bg-blue-700 text-white cursor-pointer"
              : "bg-blue-200 text-white cursor-not-allowed")
          }
        >
          {t("login.login")} {captchaToken ? "" : t("login.captcha.processing")}
        </button>
        {/* oidc登录按钮 */}
        {oidcConfigs.map((config) => (
          <Link
            key={config.name}
            to={config.loginUrl}
            className="w-full py-2 mt-4 rounded-3xl bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white text-center font-semibold hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
          >
            <div className="flex items-center justify-center gap-2">
              <div className="flex-shrink-0 w-6 h-6 relative">
                {/* 替换 Next.js Image 为常规 img 标签 */}
                <img
                  src={config.icon}
                  alt={config.displayName}
                  className="object-contain w-full h-full"
                />
              </div>
              <span className="text-base">
                {t("login.oidc.use", { provider: t(config.displayName) })}
              </span>
            </div>
          </Link>
        ))}
      </form>
    </div>
  );
}