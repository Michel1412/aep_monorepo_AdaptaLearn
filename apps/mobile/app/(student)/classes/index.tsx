import { useQuery } from '@tanstack/react-query';
import { View, Text, TouchableOpacity, StyleSheet, FlatList } from 'react-native';
import { router } from 'expo-router';
import { api, getStoredStudent } from '../../../src/lib/api';
import { useEffect, useState } from 'react';

export default function ClassesScreen() {
  const [ra, setRa] = useState<string | null>(null);
  useEffect(() => { getStoredStudent().then((s) => setRa(s?.ra ?? null)); }, []);

  const { data: classes, isLoading } = useQuery({
    queryKey: ['student-classes', ra],
    queryFn: () => api.students.getClasses(ra!),
    enabled: !!ra,
  });

  if (isLoading) return <Text style={styles.loading}>Carregando...</Text>;

  return (
    <FlatList
      style={styles.container}
      data={classes}
      keyExtractor={(item) => item.id}
      renderItem={({ item }) => (
        <TouchableOpacity
          style={styles.card}
          onPress={() => router.push(`/(student)/classes/${item.id}`)}
        >
          <Text style={styles.name}>{item.name}</Text>
          <Text style={styles.subject}>{item.subject_name}</Text>
        </TouchableOpacity>
      )}
    />
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 16, backgroundColor: '#f9fafb' },
  loading: { padding: 16 },
  card: { backgroundColor: '#fff', padding: 16, borderRadius: 12, marginBottom: 12, borderWidth: 1, borderColor: '#e5e7eb' },
  name: { fontSize: 16, fontWeight: '600' },
  subject: { color: '#6b7280', marginTop: 4 },
});
