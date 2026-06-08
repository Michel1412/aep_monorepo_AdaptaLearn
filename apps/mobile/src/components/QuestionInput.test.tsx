import { describe, it, expect } from 'vitest';
import type { AnswerType } from '@adaptalearn/shared-types';

function getAnswerOptions(answerType: AnswerType, options?: string[]) {
  if (answerType === 'true_false') return ['true', 'false'];
  if (answerType === 'multiple_choice' && options) return options;
  return null;
}

describe('QuestionInput logic', () => {
  it('returns true/false options', () => {
    expect(getAnswerOptions('true_false')).toEqual(['true', 'false']);
  });

  it('returns multiple choice options', () => {
    expect(getAnswerOptions('multiple_choice', ['A', 'B'])).toEqual(['A', 'B']);
  });

  it('returns null for essay', () => {
    expect(getAnswerOptions('essay')).toBeNull();
  });
});
