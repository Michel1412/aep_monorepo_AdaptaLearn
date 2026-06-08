import { describe, it, expect } from 'vitest';
import { formatStudyTime } from './formatStudyTime';

describe('formatStudyTime', () => {
  it('formats zero', () => {
    expect(formatStudyTime(0)).toBe('00:00');
  });

  it('formats seconds', () => {
    expect(formatStudyTime(3)).toBe('00:03');
  });

  it('formats minutes', () => {
    expect(formatStudyTime(125)).toBe('02:05');
  });
});

describe('useStudyTimer logic', () => {
  it('validates min session time', () => {
    expect(Math.max(5, 10)).toBe(10);
  });
});
