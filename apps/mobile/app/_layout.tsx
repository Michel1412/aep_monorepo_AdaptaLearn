import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Stack } from 'expo-router';
import { useEffect } from 'react';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { loadToken } from '../src/lib/api';

const queryClient = new QueryClient();

export default function RootLayout() {
  useEffect(() => {
    loadToken();
  }, []);

  return (
    <SafeAreaProvider>
      <QueryClientProvider client={queryClient}>
        <Stack screenOptions={{ headerShown: false }}>
          <Stack.Screen name="index" />
          <Stack.Screen name="(auth)" />
          <Stack.Screen name="(student)" />
        </Stack>
      </QueryClientProvider>
    </SafeAreaProvider>
  );
}
