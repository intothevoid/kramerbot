import { useEffect, useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { verifyEmail } from '../api/auth';
import { Spinner } from '../components/Spinner';
import type { WebUser } from '../types';

interface Props {
  onLogin: (user: WebUser) => void;
}

export default function VerifyEmail({ onLogin }: Props) {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [errorMsg, setErrorMsg] = useState('');

  useEffect(() => {
    const token = searchParams.get('token');
    if (!token) {
      setErrorMsg('No verification token found in the link.');
      setStatus('error');
      return;
    }

    verifyEmail(token)
      .then((data) => {
        localStorage.setItem('kb_token', data.token);
        onLogin(data.user);
        setStatus('success');
        setTimeout(() => navigate('/dashboard', { replace: true }), 1500);
      })
      .catch(() => {
        setErrorMsg('This verification link is invalid or has expired. Please sign up again or request a new link.');
        setStatus('error');
      });
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

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

        <div className="rounded-b-2xl border border-t-0 bg-white px-8 py-10 shadow-sm text-center"
          style={{ borderColor: '#e5e7eb' }}>

          {status === 'loading' && (
            <div className="space-y-4">
              <Spinner size="lg" />
              <p className="text-sm text-slate-500">Verifying your email…</p>
            </div>
          )}

          {status === 'success' && (
            <div className="space-y-4">
              <div className="text-4xl">✅</div>
              <h2 className="serif text-2xl font-bold" style={{ color: 'var(--kb-ink)' }}>
                Email verified!
              </h2>
              <p className="text-sm text-slate-500">
                You're now signed in. Redirecting to your dashboard…
              </p>
            </div>
          )}

          {status === 'error' && (
            <div className="space-y-4">
              <div className="text-4xl">❌</div>
              <h2 className="serif text-2xl font-bold" style={{ color: 'var(--kb-ink)' }}>
                Verification failed
              </h2>
              <p className="rounded-lg bg-red-50 px-4 py-3 text-sm text-red-700">{errorMsg}</p>
              <Link to="/login" className="btn-red inline-flex w-full justify-center">
                Back to Sign In
              </Link>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
