import { View, Text, StyleSheet } from 'react-native';
import type { StudentDashboard } from '@adaptalearn/shared-types';

type Subject = StudentDashboard['subjects'][number];

interface Props {
  subjects: Subject[];
  plannedHours: number;
  actualHours: number;
}

export function HoursChart({ subjects, plannedHours, actualHours }: Props) {
  const maxHours = Math.max(plannedHours, actualHours, 1);

  return (
    <View style={styles.container}>
      <Text style={styles.title}>Horas de estudo</Text>
      <View style={styles.barRow}>
        <Text style={styles.label}>Previsto</Text>
        <View style={styles.barBg}>
          <View style={[styles.barFill, styles.planned, { width: `${(plannedHours / maxHours) * 100}%` }]} />
        </View>
        <Text style={styles.value}>{plannedHours.toFixed(1)}h</Text>
      </View>
      <View style={styles.barRow}>
        <Text style={styles.label}>Realizado</Text>
        <View style={styles.barBg}>
          <View style={[styles.barFill, styles.actual, { width: `${(actualHours / maxHours) * 100}%` }]} />
        </View>
        <Text style={styles.value}>{actualHours.toFixed(1)}h</Text>
      </View>
      {subjects.map((s) => (
        <View key={s.class_id} style={styles.subjectRow}>
          <Text style={styles.subjectName}>{s.class_name}</Text>
          <Text style={styles.subjectProgress}>{s.completed}/{s.total}</Text>
        </View>
      ))}
    </View>
  );
}

const styles = StyleSheet.create({
  container: { backgroundColor: '#fff', padding: 16, borderRadius: 12, borderWidth: 1, borderColor: '#e5e7eb', marginBottom: 16 },
  title: { fontWeight: '600', marginBottom: 12 },
  barRow: { flexDirection: 'row', alignItems: 'center', marginBottom: 8, gap: 8 },
  label: { width: 70, fontSize: 12, color: '#6b7280' },
  barBg: { flex: 1, height: 12, backgroundColor: '#e5e7eb', borderRadius: 6, overflow: 'hidden' },
  barFill: { height: '100%', borderRadius: 6 },
  planned: { backgroundColor: '#9ca3af' },
  actual: { backgroundColor: '#2563eb' },
  value: { width: 40, fontSize: 12, textAlign: 'right' },
  subjectRow: { flexDirection: 'row', justifyContent: 'space-between', paddingVertical: 6, borderTopWidth: 1, borderTopColor: '#f3f4f6' },
  subjectName: { fontSize: 13 },
  subjectProgress: { fontSize: 13, color: '#6b7280' },
});
