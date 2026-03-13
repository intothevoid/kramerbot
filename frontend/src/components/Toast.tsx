import { useState, useEffect, useRef } from 'react';
import { CheckCircle } from 'lucide-react';

export function useToast() {
  const [toastMsg, setToastMsg] = useState('');
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const showToast = (msg: string) => {
    if (timerRef.current) clearTimeout(timerRef.current);
    setToastMsg(msg);
    timerRef.current = setTimeout(() => setToastMsg(''), 3000);
  };

  return { toastMsg, showToast };
}

export function Toast({ message }: { message: string }) {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    if (message) {
      setVisible(true);
    } else {
      setVisible(false);
    }
  }, [message]);

  if (!message) return null;

  return (
    <div
      className="fixed bottom-5 right-5 z-50 flex items-center gap-2 rounded-xl border px-4 py-3 shadow-lg transition-opacity"
      style={{
        background: '#FFFEF7',
        borderColor: 'var(--kb-red)',
        opacity: visible ? 1 : 0,
      }}
    >
      <CheckCircle className="h-4 w-4 shrink-0" style={{ color: 'var(--kb-red)' }} />
      <span className="text-sm font-medium" style={{ color: 'var(--kb-ink)' }}>
        {message}
      </span>
    </div>
  );
}
