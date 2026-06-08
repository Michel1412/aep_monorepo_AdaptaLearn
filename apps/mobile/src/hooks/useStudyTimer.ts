import { useEffect, useRef, useState } from 'react';
import { formatStudyTime } from './formatStudyTime';

export function useStudyTimer(autoStart = true) {
  const [seconds, setSeconds] = useState(0);
  const [running, setRunning] = useState(autoStart);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    if (running) {
      intervalRef.current = setInterval(() => setSeconds((s) => s + 1), 1000);
    }
    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, [running]);

  function pause() { setRunning(false); }
  function resume() { setRunning(true); }
  function reset() { setSeconds(0); setRunning(false); }

  return { seconds, running, pause, resume, reset, formatted: formatStudyTime(seconds) };
}
