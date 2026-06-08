import { useEffect, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { View, Text, ScrollView, StyleSheet, RefreshControl } from 'react-native';
import { api, getStoredStudent } from '../../src/lib/api';
import { HoursChart } from '../../src/components/HoursChart';

export default function DashboardScreen() {
  const [student, setStudent] = useState<{ ra: string; name: string } | null>(null);
  const [refreshing, setRefreshing] = useState(false);

  useEffect(() => {
    getStoredStudent().then(setStudent);
  }, []);

  const { data, refetch, isLoading } = useQuery({
    queryKey: ['student-dashboard', student?.ra],
    queryFn: () => api.students.getDashboard(student!.ra),
    enabled: !!student?.ra,
  });

  async function onRefresh() {
    setRefreshing(true);
    await refetch();
    setRefreshing(false);
  }

  if (!student) return null;

  return (
    <ScrollView
      style={styles.container}
      refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
    >
      <Text style={styles.greeting}>Olá, {student.name.split(' ')[0]}!</Text>
      {isLoading ? (
        <Text>Carregando...</Text>
      ) : (
        <>
          <View style={styles.row}>
            <StatCard
              title="Atividades"
              value={`${data?.completed_activities ?? 0} / ${data?.total_activities ?? 0}`}
            />
            <StatCard
              title="Horas"
              value={`${(data?.actual_hours ?? 0).toFixed(1)}h / ${(data?.planned_hours ?? 0).toFixed(1)}h`}
            />
          </View>
          {data && data.subjects.length > 0 && (
            <HoursChart subjects={data.subjects} plannedHours={data.planned_hours} actualHours={data.actual_hours} />
          )}
          <Text style={styles.sectionTitle}>Por exercício</Text>
          {data?.exercise_breakdown.map((ex) => (
            <View key={ex.activity_id} style={styles.exerciseCard}>
              <Text style={styles.exerciseTitle}>{ex.title}</Text>
              <Text style={styles.exerciseMeta}>
                Previsto: {ex.estimated_minutes} min · Gasto: {Math.round(ex.time_spent_seconds / 60)} min
              </Text>
            </View>
          ))}
        </>
      )}
    </ScrollView>
  );
}

function StatCard({ title, value }: { title: string; value: string }) {
  return (
    <View style={styles.statCard}>
      <Text style={styles.statTitle}>{title}</Text>
      <Text style={styles.statValue}>{value}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 16, backgroundColor: '#f9fafb' },
  greeting: { fontSize: 22, fontWeight: 'bold', marginBottom: 16 },
  row: { flexDirection: 'row', gap: 12, marginBottom: 16 },
  statCard: { flex: 1, backgroundColor: '#fff', padding: 16, borderRadius: 12, borderWidth: 1, borderColor: '#e5e7eb' },
  statTitle: { color: '#6b7280', fontSize: 13 },
  statValue: { fontSize: 18, fontWeight: '700', marginTop: 4 },
  sectionTitle: { fontSize: 16, fontWeight: '600', marginVertical: 12 },
  exerciseCard: { backgroundColor: '#fff', padding: 12, borderRadius: 8, marginBottom: 8, borderWidth: 1, borderColor: '#e5e7eb' },
  exerciseTitle: { fontWeight: '600' },
  exerciseMeta: { color: '#6b7280', fontSize: 12, marginTop: 4 },
});
