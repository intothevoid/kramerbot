import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { login } from '../api/auth';
import { Spinner } from '../components/Spinner';
import type { WebUser } from '../types';

interface Props {
  onLogin: (user: WebUser) => void;
}

export default function Login({ onLogin }: Props) {
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const data = await login(email, password);
      localStorage.setItem('kb_token', data.token);
      onLogin(data.user);
      navigate('/dashboard');
    } catch (err: unknown) {
      const msg =
        (err as { response?: { data?: { error?: string } } })?.response?.data?.error ?? 'Login failed';
      setError(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      className="flex min-h-screen items-center justify-center px-4"
      style={{ background: 'var(--kb-cream)' }}
    >
      <div className="w-full max-w-md">
        {/* Header stripe */}
        <div
          className="rounded-t-2xl px-8 py-5 text-center"
          style={{ background: 'var(--kb-red)' }}
        >
          <Link to="/" className="serif inline-flex items-center gap-2 text-2xl font-bold text-white">
            <img src="/kramer-icon.jpg" alt="" className="h-9 w-9 rounded-full object-cover" />
            KramerBot - Aussie Deals
          </Link>
        </div>

        <div className="rounded-b-2xl border border-t-0 bg-white px-8 py-8 shadow-sm"
          style={{ borderColor: '#e5e7eb' }}>
          <h2 className="serif mb-6 text-2xl font-bold" style={{ color: 'var(--kb-ink)' }}>
            Sign in
          </h2>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="mb-1 block text-sm font-medium text-slate-700">Email</label>
              <input
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="input-field"
                placeholder="you@example.com"
              />
            </div>
            <div>
              <div className="mb-1 flex items-center justify-between">
                <label className="text-sm font-medium text-slate-700">Password</label>
                <Link
                  to="/forgot-password"
                  className="text-xs font-medium hover:underline"
                  style={{ color: 'var(--kb-red)' }}
                >
                  Forgot password?
                </Link>
              </div>
              <input
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="input-field"
                placeholder="••••••••"
              />
            </div>

            {error && (
              <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p>
            )}

            <button type="submit" disabled={loading} className="btn-red w-full justify-center">
              {loading && <Spinner size="sm" />} Sign In
            </button>
          </form>

          <p className="mt-6 text-center text-sm text-slate-500">
            No account?{' '}
            <Link to="/signup" className="font-semibold hover:underline" style={{ color: 'var(--kb-red)' }}>
              Create one
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
