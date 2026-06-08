import { Redirect } from 'expo-router';
import { getStoredStudent } from '../src/lib/api';
import { useEffect, useState } from 'react';
import { LoadingScreen } from '../src/components/LoadingScreen';

export default function Index() {
  const [student, setStudent] = useState<{ ra: string } | null | undefined>(undefined);

  useEffect(() => {
    getStoredStudent().then(setStudent).catch(() => setStudent(null));
  }, []);

  if (student === undefined) return <LoadingScreen />;
  if (student) return <Redirect href="/(student)" />;
  return <Redirect href="/(auth)/login" />;
}
