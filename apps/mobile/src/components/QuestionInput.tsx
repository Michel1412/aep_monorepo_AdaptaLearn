import { View, Text, TextInput, TouchableOpacity, StyleSheet } from 'react-native';
import type { AnswerType } from '@adaptalearn/shared-types';

interface Props {
  answerType: AnswerType;
  value: string;
  onChange: (v: string) => void;
  options?: string[];
}

export function QuestionInput({ answerType, value, onChange, options }: Props) {
  if (answerType === 'true_false') {
    return (
      <View style={styles.row}>
        {['true', 'false'].map((opt) => (
          <TouchableOpacity
            key={opt}
            style={[styles.option, value === opt && styles.selected]}
            onPress={() => onChange(opt)}
          >
            <Text style={value === opt ? styles.selectedText : undefined}>
              {opt === 'true' ? 'Verdadeiro' : 'Falso'}
            </Text>
          </TouchableOpacity>
        ))}
      </View>
    );
  }

  if (answerType === 'multiple_choice' && options) {
    return (
      <View>
        {options.map((opt) => (
          <TouchableOpacity
            key={opt}
            style={[styles.option, value === opt && styles.selected]}
            onPress={() => onChange(opt)}
          >
            <Text style={value === opt ? styles.selectedText : undefined}>{opt}</Text>
          </TouchableOpacity>
        ))}
      </View>
    );
  }

  return (
    <TextInput
      style={[styles.input, answerType === 'essay' && styles.textarea]}
      value={value}
      onChangeText={onChange}
      placeholder="Sua resposta..."
      multiline={answerType === 'essay'}
      numberOfLines={answerType === 'essay' ? 6 : 1}
    />
  );
}

const styles = StyleSheet.create({
  input: { backgroundColor: '#fff', borderWidth: 1, borderColor: '#d1d5db', borderRadius: 8, padding: 12 },
  textarea: { minHeight: 120, textAlignVertical: 'top' },
  row: { flexDirection: 'row', gap: 12 },
  option: { flex: 1, padding: 14, borderRadius: 8, borderWidth: 1, borderColor: '#d1d5db', backgroundColor: '#fff', alignItems: 'center', marginBottom: 8 },
  selected: { borderColor: '#2563eb', backgroundColor: '#eff6ff' },
  selectedText: { color: '#2563eb', fontWeight: '600' },
});
