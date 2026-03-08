import { useState, useEffect } from 'react';
import type { WebUser } from '../types';
import { getProfile } from '../api/user';

export function useAuth() {
  const [user, setUser] = useState<WebUser | null>(null);
  // Start loading only when a token exists; no-token case is already resolved.
  const [loading, setLoading] = useState(() => !!localStorage.getItem('kb_token'));

  useEffect(() => {
    const token = localStorage.getItem('kb_token');
    if (!token) return;
    getProfile()
      .then(setUser)
      .catch(() => {
        localStorage.removeItem('kb_token');
        setUser(null);
      })
      .finally(() => setLoading(false));
  }, []);

  const signOut = () => {
    localStorage.removeItem('kb_token');
    setUser(null);
  };

  return { user, setUser, loading, signOut };
}
