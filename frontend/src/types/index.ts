export interface WebUser {
  id: string;
  email: string;
  display_name: string;
  telegram_chat_id?: number;
  telegram_username?: string;
  created_at: string;
  updated_at: string;
}

export interface OzbDeal {
  id: string;
  title: string;
  url: string;
  time: string;
  upvotes: string;
  dealage: string;
  dealtype: number;
}

export interface AmazonDeal {
  id: string;
  title: string;
  url: string;
  time: string;
  image: string;
  dealtype: number;
}

export interface APIResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

export interface AuthResponse {
  token: string;
  user: WebUser;
}

export interface DealsPage<T> {
  deals: T[];
  total: number;
}

export interface TelegramLinkResponse {
  token: string;
  deep_link: string;
  expires_at: string;
}

export interface TelegramStatus {
  linked: boolean;
  telegram_username?: string;
}
