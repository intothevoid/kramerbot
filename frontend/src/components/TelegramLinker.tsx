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
    refetchInterval: deepLink ? 4000 : false, // poll while link is pending
  });

  // Stop polling and clear the deep link once linked.
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
      <div className="flex items-center justify-between rounded-xl border border-green-200 bg-green-50 p-4">
        <div className="flex items-center gap-2 text-green-700">
          <CheckCircle className="h-5 w-5" />
          <span className="font-medium">
            Linked{status.telegram_username ? ` as @${status.telegram_username}` : ''}
          </span>
        </div>
        <button
          onClick={() => unlinkMutation.mutate()}
          disabled={unlinkMutation.isPending}
          className="flex items-center gap-1 rounded-lg px-3 py-1.5 text-sm text-red-600 hover:bg-red-100 disabled:opacity-50"
        >
          <XCircle className="h-4 w-4" /> Unlink
        </button>
      </div>
    );
  }

  const mins = String(Math.floor(secondsLeft / 60)).padStart(2, '0');
  const secs = String(secondsLeft % 60).padStart(2, '0');

  return (
    <div className="rounded-xl border border-slate-200 bg-slate-50 p-4">
      <p className="mb-3 text-sm text-slate-600">
        Link your Telegram account to receive deal notifications directly in Telegram.
      </p>

      {deepLink ? (
        <div className="space-y-3">
          <a
            href={deepLink}
            target="_blank"
            rel="noopener noreferrer"
            className="flex w-full items-center justify-center gap-2 rounded-xl bg-indigo-600 px-4 py-3 font-semibold text-white hover:bg-indigo-700"
          >
            <Send className="h-4 w-4" /> Open in Telegram
          </a>
          <div className="flex items-center justify-between text-xs text-slate-500">
            <span>Waiting for link… <Spinner size="sm" /></span>
            <span className="font-mono">{mins}:{secs}</span>
          </div>
          <button
            onClick={() => linkMutation.mutate()}
            className="flex items-center gap-1 text-xs text-indigo-600 hover:underline"
          >
            <RefreshCw className="h-3 w-3" /> Generate new link
          </button>
        </div>
      ) : (
        <button
          onClick={() => linkMutation.mutate()}
          disabled={linkMutation.isPending}
          className="flex w-full items-center justify-center gap-2 rounded-xl bg-indigo-600 px-4 py-3 font-semibold text-white hover:bg-indigo-700 disabled:opacity-60"
        >
          {linkMutation.isPending ? <Spinner size="sm" /> : <Send className="h-4 w-4" />}
          Link Telegram Account
        </button>
      )}
    </div>
  );
}
