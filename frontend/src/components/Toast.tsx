import { CheckCircle } from 'lucide-react';

export function Toast({ message }: { message: string }) {
  if (!message) return null;

  return (
    <div
      className="fixed bottom-5 right-5 z-50 flex items-center gap-2 rounded-xl border px-4 py-3 shadow-lg"
      style={{ background: '#FFFEF7', borderColor: 'var(--kb-red)' }}
    >
      <CheckCircle className="h-4 w-4 shrink-0" style={{ color: 'var(--kb-red)' }} />
      <span className="text-sm font-medium" style={{ color: 'var(--kb-ink)' }}>
        {message}
      </span>
    </div>
  );
}
