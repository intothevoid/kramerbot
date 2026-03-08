import api from './client';
import type { APIResponse, AuthResponse } from '../types';

export async function register(email: string, password: string, displayName: string): Promise<{ message: string }> {
  const res = await api.post<APIResponse<{ message: string }>>('/auth/register', {
    email,
    password,
    display_name: displayName,
  });
  return res.data.data!;
}

export async function login(email: string, password: string): Promise<AuthResponse> {
  const res = await api.post<APIResponse<AuthResponse>>('/auth/login', { email, password });
  return res.data.data!;
}

export async function logout(): Promise<void> {
  await api.post('/auth/logout');
}

export async function verifyEmail(token: string): Promise<AuthResponse> {
  const res = await api.get<APIResponse<AuthResponse>>(`/auth/verify-email?token=${encodeURIComponent(token)}`);
  return res.data.data!;
}

export async function forgotPassword(email: string): Promise<{ message: string }> {
  const res = await api.post<APIResponse<{ message: string }>>('/auth/forgot-password', { email });
  return res.data.data!;
}

export async function resetPassword(token: string, password: string): Promise<{ message: string }> {
  const res = await api.post<APIResponse<{ message: string }>>('/auth/reset-password', { token, password });
  return res.data.data!;
}
