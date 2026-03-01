import api from './client';
import type {
  APIResponse,
  DealsPage,
  OzbDeal,
  AmazonDeal,
  TelegramLinkResponse,
  TelegramStatus,
  WebUser,
} from '../types';

export async function getProfile(): Promise<WebUser> {
  const res = await api.get<APIResponse<WebUser>>('/user/profile');
  return res.data.data!;
}

export async function getKeywords(): Promise<string[]> {
  const res = await api.get<APIResponse<{ keywords: string[] }>>('/user/keywords');
  return res.data.data?.keywords ?? [];
}

export async function updatePreferences(prefs: {
  ozb_good: boolean;
  ozb_super: boolean;
  amz_daily: boolean;
  amz_weekly: boolean;
}): Promise<WebUser> {
  const res = await api.put<APIResponse<WebUser>>('/user/preferences', prefs);
  return res.data.data!;
}

export async function addKeyword(keyword: string): Promise<string[]> {
  const res = await api.post<APIResponse<{ keywords: string[] }>>('/user/keywords', { keyword });
  return res.data.data?.keywords ?? [];
}

export async function removeKeyword(keyword: string): Promise<string[]> {
  const res = await api.delete<APIResponse<{ keywords: string[] }>>(`/user/keywords/${encodeURIComponent(keyword)}`);
  return res.data.data?.keywords ?? [];
}

export async function generateTelegramLink(): Promise<TelegramLinkResponse> {
  const res = await api.post<APIResponse<TelegramLinkResponse>>('/user/telegram/link');
  return res.data.data!;
}

export async function getTelegramStatus(): Promise<TelegramStatus> {
  const res = await api.get<APIResponse<TelegramStatus>>('/user/telegram/status');
  return res.data.data!;
}

export async function unlinkTelegram(): Promise<void> {
  await api.delete('/user/telegram/link');
}

export async function getOzbDeals(type?: string, limit = 50, offset = 0): Promise<DealsPage<OzbDeal>> {
  const params: Record<string, string | number> = { limit, offset };
  if (type) params.type = type;
  const res = await api.get<APIResponse<DealsPage<OzbDeal>>>('/deals/ozbargain', { params });
  return res.data.data!;
}

export async function getAmazonDeals(type?: string, limit = 50, offset = 0): Promise<DealsPage<AmazonDeal>> {
  const params: Record<string, string | number> = { limit, offset };
  if (type) params.type = type;
  const res = await api.get<APIResponse<DealsPage<AmazonDeal>>>('/deals/amazon', { params });
  return res.data.data!;
}
