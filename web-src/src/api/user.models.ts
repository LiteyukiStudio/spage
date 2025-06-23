/**
 * This file defines the request models used in the API.
 */
import type { CaptchaProvider } from "@/types/captcha";

import type { BaseResponse } from "./base.models";
import type { Org } from "./org.models";
import type { Project } from "./project.models";

/**
 * 用户 User
 */
export interface User {
  avatarUrl?: string;
  description: string;
  displayName?: string;
  email?: null | string;
  id: number;
  language: string;
  name: string;
  organizations: number[];
  projectLimit: number;
  role: "admin" | "user";
  [property: string]: unknown;
}

export interface LoginRequest {
  username: string;
  password: string;
  captchaToken: string;
  remember: boolean;
}

export interface LoginResponse extends BaseResponse {
  token: string;
  refreshToken: string;
  userId: number;
}

export interface RegisterRequest {
  username: string;
  password: string;
  email: string;
  captchaToken: string;
  verifyCode: string;
}

export interface RegisterResponse extends BaseResponse {
  token: string;
  userId: number;
}

export interface UserResponse extends BaseResponse {
  user: User | null;
}

export interface UserProjectsResponse extends BaseResponse {
  projects: Project[];
  total: number;
}

export interface UserOrganizationsResponse extends BaseResponse {
  organizations: Org[];
}

export interface CaptchaConfig extends BaseResponse {
  provider: CaptchaProvider;
  siteKey: string;
  url?: string;
}

export interface OidcConfig {
  name: string;
  displayName: string;
  icon: string;
  loginUrl: string;
}
