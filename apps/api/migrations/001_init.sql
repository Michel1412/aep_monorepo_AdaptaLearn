-- Teachers
CREATE TABLE IF NOT EXISTS teachers (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

-- Students
CREATE TABLE IF NOT EXISTS students (
    ra TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

-- Subjects
CREATE TABLE IF NOT EXISTS subjects (
    id TEXT PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL
);

-- Classes
CREATE TABLE IF NOT EXISTS classes (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    subject_id TEXT NOT NULL REFERENCES subjects(id),
    teacher_id TEXT NOT NULL REFERENCES teachers(id),
    series INTEGER NOT NULL,
    turma TEXT NOT NULL
);

-- Enrollments
CREATE TABLE IF NOT EXISTS enrollments (
    id TEXT PRIMARY KEY,
    student_ra TEXT NOT NULL REFERENCES students(ra),
    class_id TEXT NOT NULL REFERENCES classes(id),
    UNIQUE(student_ra, class_id)
);

-- Activities
CREATE TABLE IF NOT EXISTS activities (
    id TEXT PRIMARY KEY,
    class_id TEXT NOT NULL REFERENCES classes(id),
    title TEXT NOT NULL,
    question TEXT NOT NULL,
    estimated_minutes INTEGER NOT NULL,
    answer_type TEXT NOT NULL,
    complement_type TEXT NOT NULL DEFAULT 'none',
    complement_url TEXT,
    week_start TEXT NOT NULL,
    created_by TEXT NOT NULL REFERENCES teachers(id),
    options TEXT
);

-- Study Sessions
CREATE TABLE IF NOT EXISTS study_sessions (
    id TEXT PRIMARY KEY,
    activity_id TEXT NOT NULL REFERENCES activities(id),
    student_ra TEXT NOT NULL REFERENCES students(ra),
    started_at TEXT NOT NULL,
    ended_at TEXT
);

-- Submissions
CREATE TABLE IF NOT EXISTS submissions (
    id TEXT PRIMARY KEY,
    activity_id TEXT NOT NULL REFERENCES activities(id),
    student_ra TEXT NOT NULL REFERENCES students(ra),
    session_id TEXT NOT NULL REFERENCES study_sessions(id),
    answer TEXT NOT NULL,
    time_spent_seconds INTEGER NOT NULL,
    submitted_at TEXT NOT NULL,
    UNIQUE(activity_id, student_ra)
);

CREATE INDEX IF NOT EXISTS idx_activities_class ON activities(class_id);
CREATE INDEX IF NOT EXISTS idx_activities_week ON activities(week_start);
CREATE INDEX IF NOT EXISTS idx_enrollments_class ON enrollments(class_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_student ON enrollments(student_ra);
