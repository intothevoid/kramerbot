import { useState, useRef } from 'react';

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
