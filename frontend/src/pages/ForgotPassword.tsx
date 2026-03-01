import { useState } from 'react';
import { Link } from 'react-router-dom';
import { forgotPassword } from '../api/auth';
import { Spinner } from '../components/Spinner';

export default function ForgotPassword() {
  const [email, setEmail] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<{ message: string; reset_link?: string } | null>(null);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const data = await forgotPassword(email);
      setResult(data);
    } catch {
      setError('Something went wrong. Please try again.');
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
          <h2 className="serif mb-2 text-2xl font-bold" style={{ color: 'var(--kb-ink)' }}>
            Forgot password
          </h2>
          <p className="mb-6 text-sm text-slate-500">
            Enter your email and we'll generate a reset link.
          </p>

          {result ? (
            <div className="space-y-4">
              <p className="rounded-lg bg-green-50 px-4 py-3 text-sm text-green-800">
                {result.message}
              </p>
              {result.reset_link && (
                <div>
                  <p className="mb-1 text-xs font-semibold text-slate-500 uppercase tracking-wide">
                    Reset link (copy this):
                  </p>
                  <a
                    href={result.reset_link}
                    className="block break-all rounded-lg border px-3 py-2 text-xs hover:underline"
                    style={{ borderColor: 'var(--kb-yellow)', color: 'var(--kb-red)' }}
                  >
                    {result.reset_link}
                  </a>
                </div>
              )}
              <Link to="/login" className="btn-red w-full justify-center text-sm">
                Back to Sign In
              </Link>
            </div>
          ) : (
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

              {error && (
                <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p>
              )}

              <button type="submit" disabled={loading} className="btn-red w-full justify-center">
                {loading && <Spinner size="sm" />} Generate Reset Link
              </button>
            </form>
          )}

          {!result && (
            <p className="mt-6 text-center text-sm text-slate-500">
              Remembered it?{' '}
              <Link to="/login" className="font-semibold hover:underline" style={{ color: 'var(--kb-red)' }}>
                Sign in
              </Link>
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
