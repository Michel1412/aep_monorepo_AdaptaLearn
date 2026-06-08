import { Stack } from 'expo-router';

export default function ClassesLayout() {
  return (
    <Stack>
      <Stack.Screen name="index" options={{ title: 'Minhas Turmas' }} />
      <Stack.Screen name="[id]" options={{ title: 'Atividades da Semana' }} />
    </Stack>
  );
}
