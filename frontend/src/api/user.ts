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

export async function addKeyword(keyword: string): Promise<void> {
  await api.post('/user/keywords', { keyword });
}

export async function removeKeyword(keyword: string): Promise<void> {
  await api.delete(`/user/keywords/${encodeURIComponent(keyword)}`);
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
