import { useState, useEffect, useRef } from 'react';
import { Send, CheckCircle, XCircle, RefreshCw } from 'lucide-react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { generateTelegramLink, getTelegramStatus, unlinkTelegram } from '../api/user';
import { Spinner } from './Spinner';

export function TelegramLinker() {
  const qc = useQueryClient();
  const [deepLink, setDeepLink] = useState<string | null>(null);
  const [secondsLeft, setSecondsLeft] = useState(0);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const { data: status, isLoading } = useQuery({
    queryKey: ['telegramStatus'],
    queryFn: getTelegramStatus,
    refetchInterval: deepLink ? 4000 : false,
  });

  useEffect(() => {
    if (status?.linked && deepLink) {
      setDeepLink(null);
      clearInterval(timerRef.current!);
    }
  }, [status?.linked, deepLink]);

  const linkMutation = useMutation({
    mutationFn: generateTelegramLink,
    onSuccess: (data) => {
      setDeepLink(data.deep_link);
      setSecondsLeft(15 * 60);
      timerRef.current = setInterval(() => {
        setSecondsLeft((s) => {
          if (s <= 1) {
            clearInterval(timerRef.current!);
            setDeepLink(null);
            return 0;
          }
          return s - 1;
        });
      }, 1000);
    },
  });

  const unlinkMutation = useMutation({
    mutationFn: unlinkTelegram,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['telegramStatus'] }),
  });

  if (isLoading) return <Spinner size="sm" />;

  if (status?.linked) {
    return (
      <div
        className="flex items-center justify-between rounded-xl border p-4"
        style={{ background: '#f0fdf4', borderColor: '#bbf7d0' }}
      >
        <div className="flex items-center gap-2" style={{ color: '#166534' }}>
          <CheckCircle className="h-5 w-5" />
          <span className="font-medium">
            Linked{status.telegram_username ? ` as @${status.telegram_username}` : ''}
          </span>
        </div>
        <button
          onClick={() => unlinkMutation.mutate()}
          disabled={unlinkMutation.isPending}
          className="flex items-center gap-1 rounded-lg px-3 py-1.5 text-sm disabled:opacity-50"
          style={{ color: 'var(--kb-red)' }}
        >
          <XCircle className="h-4 w-4" /> Unlink
        </button>
      </div>
    );
  }

  const mins = String(Math.floor(secondsLeft / 60)).padStart(2, '0');
  const secs = String(secondsLeft % 60).padStart(2, '0');

  return (
    <div>
      <p className="mb-3 text-sm text-slate-600">
        Link your Telegram account to receive deal notifications.
      </p>

      {deepLink ? (
        <div className="space-y-3">
          <a
            href={deepLink}
            target="_blank"
            rel="noopener noreferrer"
            className="btn-red flex w-full items-center justify-center gap-2"
          >
            <Send className="h-4 w-4" /> Open in Telegram
          </a>
          <div className="flex items-center justify-between text-xs text-slate-500">
            <span className="flex items-center gap-1">Waiting for link… <Spinner size="sm" /></span>
            <span className="font-mono">{mins}:{secs}</span>
          </div>
          <button
            onClick={() => linkMutation.mutate()}
            className="flex items-center gap-1 text-xs hover:underline"
            style={{ color: 'var(--kb-red)' }}
          >
            <RefreshCw className="h-3 w-3" /> Generate new link
          </button>
        </div>
      ) : (
        <button
          onClick={() => linkMutation.mutate()}
          disabled={linkMutation.isPending}
          className="btn-red flex w-full items-center justify-center gap-2 disabled:opacity-60"
        >
          {linkMutation.isPending ? <Spinner size="sm" /> : <Send className="h-4 w-4" />}
          Link Telegram Account
        </button>
      )}
    </div>
  );
}
