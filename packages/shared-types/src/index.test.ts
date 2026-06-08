import { describe, it, expect } from 'vitest';
import { CreateActivitySchema, formatClassName } from './index';

describe('CreateActivitySchema', () => {
  it('validates a correct activity', () => {
    const result = CreateActivitySchema.safeParse({
      title: 'Test',
      question: 'What is AFD?',
      estimated_minutes: 30,
      answer_type: 'essay',
      complement_type: 'none',
      week_start: '2026-06-02',
    });
    expect(result.success).toBe(true);
  });

  it('rejects empty title', () => {
    const result = CreateActivitySchema.safeParse({
      title: '',
      question: 'Q',
      estimated_minutes: 10,
      answer_type: 'essay',
      complement_type: 'none',
      week_start: '2026-06-02',
    });
    expect(result.success).toBe(false);
  });
});

describe('formatClassName', () => {
  it('formats class name correctly', () => {
    expect(formatClassName('TeoCom', 7, 'B')).toBe('TeoCom-S7-B');
  });
});
