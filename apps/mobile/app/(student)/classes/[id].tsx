import { useQuery } from '@tanstack/react-query';
import { View, Text, TouchableOpacity, StyleSheet, FlatList } from 'react-native';
import { router, useLocalSearchParams } from 'expo-router';
import { api, getStoredStudent } from '../../../src/lib/api';
import { useEffect, useState } from 'react';

export default function ClassActivitiesScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const [ra, setRa] = useState<string | null>(null);
  useEffect(() => { getStoredStudent().then((s) => setRa(s?.ra ?? null)); }, []);

  const { data: activities, isLoading } = useQuery({
    queryKey: ['class-activities', ra, id],
    queryFn: () => api.students.getActivities(ra!, id!),
    enabled: !!ra && !!id,
  });

  if (isLoading) return <Text style={styles.loading}>Carregando...</Text>;

  return (
    <FlatList
      style={styles.container}
      data={activities}
      keyExtractor={(item) => item.id}
      ListHeaderComponent={<Text style={styles.header}>Atividades da semana</Text>}
      renderItem={({ item }) => (
        <TouchableOpacity
          style={styles.card}
          onPress={() => router.push(`/(student)/activities/${item.id}`)}
        >
          <View style={styles.cardHeader}>
            <Text style={styles.title}>{item.title}</Text>
            <View style={[styles.badge, item.status === 'completed' ? styles.done : styles.pending]}>
              <Text style={styles.badgeText}>{item.status === 'completed' ? 'Concluída' : 'Pendente'}</Text>
            </View>
          </View>
          <Text style={styles.meta}>{item.estimated_minutes} min estimados</Text>
        </TouchableOpacity>
      )}
    />
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 16, backgroundColor: '#f9fafb' },
  loading: { padding: 16 },
  header: { fontSize: 18, fontWeight: '600', marginBottom: 12 },
  card: { backgroundColor: '#fff', padding: 16, borderRadius: 12, marginBottom: 12, borderWidth: 1, borderColor: '#e5e7eb' },
  cardHeader: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center' },
  title: { fontSize: 15, fontWeight: '600', flex: 1 },
  badge: { paddingHorizontal: 8, paddingVertical: 4, borderRadius: 12 },
  pending: { backgroundColor: '#fef3c7' },
  done: { backgroundColor: '#d1fae5' },
  badgeText: { fontSize: 11, fontWeight: '600' },
  meta: { color: '#6b7280', fontSize: 12, marginTop: 8 },
});
