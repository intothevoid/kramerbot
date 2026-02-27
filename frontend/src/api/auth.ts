import api from './client';
import type { APIResponse, AuthResponse } from '../types';

export async function register(email: string, password: string, displayName: string): Promise<AuthResponse> {
  const res = await api.post<APIResponse<AuthResponse>>('/auth/register', {
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
