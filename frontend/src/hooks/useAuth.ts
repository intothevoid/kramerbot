import { useState, useEffect } from 'react';
import type { WebUser } from '../types';
import { getProfile } from '../api/user';

export function useAuth() {
  const [user, setUser] = useState<WebUser | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('kb_token');
    if (!token) {
      setLoading(false);
      return;
    }
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
