import { useState } from 'react';
import { Link } from 'react-router-dom';
import { register } from '../api/auth';
import { Spinner } from '../components/Spinner';
import type { WebUser } from '../types';

interface Props {
  onLogin: (user: WebUser) => void;
}

export default function Signup({ onLogin: _onLogin }: Props) {
  const [displayName, setDisplayName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [registered, setRegistered] = useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (password.length < 8) {
      setError('Password must be at least 8 characters');
      return;
    }
    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return;
    }
    setError('');
    setLoading(true);
    try {
      await register(email, password, displayName);
      setRegistered(true);
    } catch (err: unknown) {
      const msg =
        (err as { response?: { data?: { error?: string } } })?.response?.data?.error ?? 'Registration failed';
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
        <div
          className="rounded-t-2xl px-8 py-5 text-center"
          style={{ background: 'var(--kb-red)' }}
        >
          <Link to="/" className="serif inline-flex items-center gap-2 text-2xl font-bold text-white">
            <img src="/kramer-icon.jpg" alt="" className="h-9 w-9 rounded-full object-cover" />
            KramerBot
          </Link>
        </div>

        <div className="rounded-b-2xl border border-t-0 bg-white px-8 py-8 shadow-sm"
          style={{ borderColor: '#e5e7eb' }}>

          {registered ? (
            <div className="space-y-4 text-center">
              <div className="text-4xl">✉️</div>
              <h2 className="serif text-2xl font-bold" style={{ color: 'var(--kb-ink)' }}>
                Check your email
              </h2>
              <p className="text-sm text-slate-500">
                We sent a verification link to <span className="font-medium text-slate-700">{email}</span>.
                Click the link to activate your account.
              </p>
              <p className="text-xs text-slate-400">
                Didn't receive it? Check your spam folder or{' '}
                <button
                  onClick={() => setRegistered(false)}
                  className="font-medium hover:underline"
                  style={{ color: 'var(--kb-red)' }}
                >
                  try again
                </button>.
              </p>
              <Link to="/login" className="btn-red inline-flex w-full justify-center mt-2">
                Back to Sign In
              </Link>
            </div>
          ) : (
            <>
              <h2 className="serif mb-6 text-2xl font-bold" style={{ color: 'var(--kb-ink)' }}>
                Create account
              </h2>

              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label className="mb-1 block text-sm font-medium text-slate-700">Display name</label>
                  <input
                    type="text"
                    value={displayName}
                    onChange={(e) => setDisplayName(e.target.value)}
                    className="input-field"
                    placeholder="Jerry Seinfeld"
                  />
                </div>
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
                  <label className="mb-1 block text-sm font-medium text-slate-700">
                    Password <span className="text-slate-400">(min. 8 chars)</span>
                  </label>
                  <input
                    type="password"
                    required
                    minLength={8}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="input-field"
                    placeholder="••••••••"
                  />
                </div>
                <div>
                  <label className="mb-1 block text-sm font-medium text-slate-700">
                    Confirm password
                  </label>
                  <input
                    type="password"
                    required
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    className="input-field"
                    placeholder="••••••••"
                  />
                </div>

                {error && (
                  <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p>
                )}

                <button type="submit" disabled={loading} className="btn-red w-full justify-center">
                  {loading && <Spinner size="sm" />} Create Account
                </button>
              </form>

              <p className="mt-6 text-center text-sm text-slate-500">
                Already have an account?{' '}
                <Link to="/login" className="font-semibold hover:underline" style={{ color: 'var(--kb-red)' }}>
                  Sign in
                </Link>
              </p>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
