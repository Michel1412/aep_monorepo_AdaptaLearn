import { useState, useEffect } from 'react';
import { View, Text, TextInput, TouchableOpacity, StyleSheet, ScrollView, Linking } from 'react-native';
import { router, useLocalSearchParams } from 'expo-router';
import { useStudyTimer } from '../../../src/hooks/useStudyTimer';
import { QuestionInput } from '../../../src/components/QuestionInput';
import { api } from '../../../src/lib/api';
import type { Activity, AnswerType } from '@adaptalearn/shared-types';

export default function ActivityScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const [activity, setActivity] = useState<Activity | null>(null);
  const [answer, setAnswer] = useState('');
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const { seconds, formatted } = useStudyTimer(true);

  useEffect(() => {
    async function load() {
      try {
        const session = await api.activities.startSession(id!);
        setSessionId(session.session_id);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Erro ao iniciar sessão');
      }
    }
    load();
  }, [id]);

  useEffect(() => {
    async function fetchActivity() {
      const { getStoredStudent } = await import('../../../src/lib/api');
      const student = await getStoredStudent();
      if (!student) return;
      const classes = await api.students.getClasses(student.ra);
      for (const c of classes) {
        const acts = await api.students.getActivities(student.ra, c.id);
        const found = acts.find((a) => a.id === id);
        if (found) { setActivity(found); return; }
      }
    }
    fetchActivity();
  }, [id]);

  async function handleSubmit() {
    if (!answer.trim()) { setError('Resposta obrigatória'); return; }
    if (!sessionId) return;
    setSubmitting(true);
    try {
      const timeSpent = Math.max(seconds, 10);
      await api.activities.submit(id!, sessionId, answer, timeSpent);
      router.replace(`/(student)/activities/${id}/result?time=${timeSpent}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erro ao enviar');
    } finally {
      setSubmitting(false);
    }
  }

  if (!activity) return <Text style={styles.loading}>Carregando...</Text>;

  return (
    <ScrollView style={styles.container}>
      <Text style={styles.title}>{activity.title}</Text>
      <View style={styles.timerBox}>
        <Text style={styles.timerLabel}>Tempo de estudo</Text>
        <Text style={styles.timer}>{formatted}</Text>
      </View>
      <Text style={styles.question}>{activity.question}</Text>
      {activity.complement_type !== 'none' && activity.complement_url && (
        <TouchableOpacity onPress={() => Linking.openURL(activity.complement_url!)}>
          <Text style={styles.link}>Abrir complemento ({activity.complement_type})</Text>
        </TouchableOpacity>
      )}
      <QuestionInput
        answerType={activity.answer_type as AnswerType}
        value={answer}
        onChange={setAnswer}
        options={activity.options ? JSON.parse(activity.options) : undefined}
      />
      {error ? <Text style={styles.error}>{error}</Text> : null}
      <TouchableOpacity style={styles.button} onPress={handleSubmit} disabled={submitting}>
        <Text style={styles.buttonText}>{submitting ? 'Enviando...' : 'Enviar resposta'}</Text>
      </TouchableOpacity>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 16, backgroundColor: '#f9fafb' },
  loading: { padding: 16 },
  title: { fontSize: 20, fontWeight: 'bold', marginBottom: 12 },
  timerBox: { backgroundColor: '#2563eb', padding: 12, borderRadius: 8, alignItems: 'center', marginBottom: 16 },
  timerLabel: { color: '#bfdbfe', fontSize: 12 },
  timer: { color: '#fff', fontSize: 28, fontWeight: 'bold' },
  question: { fontSize: 15, lineHeight: 22, marginBottom: 16 },
  link: { color: '#2563eb', marginBottom: 16 },
  error: { color: '#dc2626', marginBottom: 8 },
  button: { backgroundColor: '#2563eb', padding: 14, borderRadius: 8, alignItems: 'center', marginTop: 16 },
  buttonText: { color: '#fff', fontWeight: '600' },
});
