import { useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { resetPassword } from '../api/auth';
import { Spinner } from '../components/Spinner';

export default function ResetPassword() {
  const [params] = useSearchParams();
  const token = params.get('token') ?? '';
  const navigate = useNavigate();

  const [password, setPassword] = useState('');
  const [confirm, setConfirm] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (password !== confirm) {
      setError('Passwords do not match');
      return;
    }
    if (password.length < 8) {
      setError('Password must be at least 8 characters');
      return;
    }
    if (!token) {
      setError('Missing reset token — please use the link from your reset email.');
      return;
    }
    setError('');
    setLoading(true);
    try {
      await resetPassword(token, password);
      setSuccess(true);
      setTimeout(() => navigate('/login'), 3000);
    } catch (err: unknown) {
      const msg =
        (err as { response?: { data?: { error?: string } } })?.response?.data?.error ??
        'Invalid or expired token. Please request a new reset link.';
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
          <h2 className="serif mb-6 text-2xl font-bold" style={{ color: 'var(--kb-ink)' }}>
            Reset password
          </h2>

          {success ? (
            <div className="space-y-4 text-center">
              <p className="rounded-lg bg-green-50 px-4 py-3 text-sm text-green-800">
                ✅ Password updated! Redirecting to sign in…
              </p>
              <Link to="/login" className="btn-red w-full justify-center text-sm">
                Go to Sign In
              </Link>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="mb-1 block text-sm font-medium text-slate-700">New password</label>
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
                <label className="mb-1 block text-sm font-medium text-slate-700">Confirm password</label>
                <input
                  type="password"
                  required
                  value={confirm}
                  onChange={(e) => setConfirm(e.target.value)}
                  className="input-field"
                  placeholder="••••••••"
                />
              </div>

              {error && (
                <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p>
              )}

              <button type="submit" disabled={loading} className="btn-red w-full justify-center">
                {loading && <Spinner size="sm" />} Update Password
              </button>
            </form>
          )}
        </div>
      </div>
    </div>
  );
}
