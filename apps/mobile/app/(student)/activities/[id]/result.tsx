import { View, Text, TouchableOpacity, StyleSheet } from 'react-native';
import { router, useLocalSearchParams } from 'expo-router';

export default function ResultScreen() {
  const { time } = useLocalSearchParams<{ time: string }>();
  const minutes = Math.round(Number(time) / 60);

  return (
    <View style={styles.container}>
      <Text style={styles.icon}>✓</Text>
      <Text style={styles.title}>Resposta enviada!</Text>
      <Text style={styles.subtitle}>Tempo de estudo registrado: {minutes} min</Text>
      <TouchableOpacity style={styles.button} onPress={() => router.replace('/(student)')}>
        <Text style={styles.buttonText}>Voltar ao dashboard</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, justifyContent: 'center', alignItems: 'center', padding: 24, backgroundColor: '#f9fafb' },
  icon: { fontSize: 48, color: '#16a34a', marginBottom: 16 },
  title: { fontSize: 22, fontWeight: 'bold', marginBottom: 8 },
  subtitle: { color: '#6b7280', marginBottom: 32 },
  button: { backgroundColor: '#2563eb', paddingHorizontal: 24, paddingVertical: 14, borderRadius: 8 },
  buttonText: { color: '#fff', fontWeight: '600' },
});
