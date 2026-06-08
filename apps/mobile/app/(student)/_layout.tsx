import { Tabs } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';

export default function StudentLayout() {
  return (
    <Tabs screenOptions={{ tabBarActiveTintColor: '#2563eb' }}>
      <Tabs.Screen
        name="index"
        options={{ title: 'Dashboard', tabBarIcon: ({ color, size }) => <Ionicons name="stats-chart" size={size} color={color} /> }}
      />
      <Tabs.Screen
        name="classes"
        options={{ title: 'Turmas', headerShown: false, tabBarIcon: ({ color, size }) => <Ionicons name="school" size={size} color={color} /> }}
      />
      <Tabs.Screen name="activities" options={{ href: null }} />
    </Tabs>
  );
}
