package repository

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	"github.com/adaptalearn/api/internal/domain"
	"github.com/google/uuid"
)

//go:embed migrations/001_init.sql
var migrationSQL string

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) Migrate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, migrationSQL)
	return err
}

func (r *SQLiteRepository) GetTeacherByEmail(ctx context.Context, email string) (*domain.Teacher, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, email, password_hash FROM teachers WHERE email = ?`, email)
	var t domain.Teacher
	if err := row.Scan(&t.ID, &t.Name, &t.Email, &t.PasswordHash); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *SQLiteRepository) GetTeacherByID(ctx context.Context, id string) (*domain.Teacher, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, email, password_hash FROM teachers WHERE id = ?`, id)
	var t domain.Teacher
	if err := row.Scan(&t.ID, &t.Name, &t.Email, &t.PasswordHash); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *SQLiteRepository) GetStudentByRA(ctx context.Context, ra string) (*domain.Student, error) {
	row := r.db.QueryRowContext(ctx, `SELECT ra, name FROM students WHERE ra = ?`, ra)
	var s domain.Student
	if err := row.Scan(&s.RA, &s.Name); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SQLiteRepository) CreateStudent(ctx context.Context, student *domain.Student) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO students (ra, name) VALUES (?, ?)`, student.RA, student.Name)
	return err
}

func (r *SQLiteRepository) GetClassesByTeacher(ctx context.Context, teacherID string) ([]domain.Class, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id, c.name, c.subject_id, c.teacher_id, c.series, c.turma,
		       s.code, s.name,
		       (SELECT COUNT(*) FROM enrollments e WHERE e.class_id = c.id),
		       (SELECT COUNT(*) FROM activities a WHERE a.class_id = c.id)
		FROM classes c
		JOIN subjects s ON s.id = c.subject_id
		WHERE c.teacher_id = ?
		ORDER BY c.series, c.turma`, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []domain.Class
	for rows.Next() {
		var c domain.Class
		if err := rows.Scan(&c.ID, &c.Name, &c.SubjectID, &c.TeacherID, &c.Series, &c.Turma,
			&c.SubjectCode, &c.SubjectName, &c.StudentCount, &c.ActivityCount); err != nil {
			return nil, err
		}
		classes = append(classes, c)
	}
	return classes, rows.Err()
}

func (r *SQLiteRepository) GetClassByID(ctx context.Context, id string) (*domain.Class, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT c.id, c.name, c.subject_id, c.teacher_id, c.series, c.turma, s.code, s.name
		FROM classes c JOIN subjects s ON s.id = c.subject_id WHERE c.id = ?`, id)
	var c domain.Class
	if err := row.Scan(&c.ID, &c.Name, &c.SubjectID, &c.TeacherID, &c.Series, &c.Turma, &c.SubjectCode, &c.SubjectName); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *SQLiteRepository) GetClassesByStudent(ctx context.Context, ra string) ([]domain.Class, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id, c.name, c.subject_id, c.teacher_id, c.series, c.turma, s.code, s.name
		FROM classes c
		JOIN subjects s ON s.id = c.subject_id
		JOIN enrollments e ON e.class_id = c.id
		WHERE e.student_ra = ?
		ORDER BY c.name`, ra)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []domain.Class
	for rows.Next() {
		var c domain.Class
		if err := rows.Scan(&c.ID, &c.Name, &c.SubjectID, &c.TeacherID, &c.Series, &c.Turma, &c.SubjectCode, &c.SubjectName); err != nil {
			return nil, err
		}
		classes = append(classes, c)
	}
	return classes, rows.Err()
}

func (r *SQLiteRepository) GetStudentsByClass(ctx context.Context, classID string) ([]domain.Student, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.ra, s.name FROM students s
		JOIN enrollments e ON e.student_ra = s.ra
		WHERE e.class_id = ? ORDER BY s.name`, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []domain.Student
	for rows.Next() {
		var s domain.Student
		if err := rows.Scan(&s.RA, &s.Name); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, rows.Err()
}

func (r *SQLiteRepository) IsStudentEnrolled(ctx context.Context, classID, ra string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM enrollments WHERE class_id = ? AND student_ra = ?`, classID, ra).Scan(&count)
	return count > 0, err
}

func (r *SQLiteRepository) EnrollStudent(ctx context.Context, classID, ra string) (*domain.Student, error) {
	student, err := r.GetStudentByRA(ctx, ra)
	if err != nil {
		return nil, fmt.Errorf("aluno não encontrado")
	}
	enrolled, err := r.IsStudentEnrolled(ctx, classID, ra)
	if err != nil {
		return nil, err
	}
	if enrolled {
		return nil, fmt.Errorf("aluno já matriculado nesta turma")
	}
	_, err = r.db.ExecContext(ctx, `INSERT INTO enrollments (id, student_ra, class_id) VALUES (?, ?, ?)`,
		uuid.New().String(), ra, classID)
	if err != nil {
		return nil, err
	}
	return student, nil
}

func (r *SQLiteRepository) GetActivitiesByClass(ctx context.Context, classID, weekStart string) ([]domain.Activity, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, class_id, title, question, estimated_minutes, answer_type,
		       complement_type, complement_url, week_start, created_by, options
		FROM activities WHERE class_id = ? AND week_start = ?
		ORDER BY title`, classID, weekStart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanActivities(rows)
}

func scanActivities(rows *sql.Rows) ([]domain.Activity, error) {
	var activities []domain.Activity
	for rows.Next() {
		var a domain.Activity
		var complementURL, options sql.NullString
		if err := rows.Scan(&a.ID, &a.ClassID, &a.Title, &a.Question, &a.EstimatedMinutes,
			&a.AnswerType, &a.ComplementType, &complementURL, &a.WeekStart, &a.CreatedBy, &options); err != nil {
			return nil, err
		}
		if complementURL.Valid {
			a.ComplementURL = &complementURL.String
		}
		if options.Valid {
			a.Options = &options.String
		}
		activities = append(activities, a)
	}
	return activities, rows.Err()
}

func (r *SQLiteRepository) GetActivityByID(ctx context.Context, id string) (*domain.Activity, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, class_id, title, question, estimated_minutes, answer_type,
		       complement_type, complement_url, week_start, created_by, options
		FROM activities WHERE id = ?`, id)
	var a domain.Activity
	var complementURL, options sql.NullString
	if err := row.Scan(&a.ID, &a.ClassID, &a.Title, &a.Question, &a.EstimatedMinutes,
		&a.AnswerType, &a.ComplementType, &complementURL, &a.WeekStart, &a.CreatedBy, &options); err != nil {
		return nil, err
	}
	if complementURL.Valid {
		a.ComplementURL = &complementURL.String
	}
	if options.Valid {
		a.Options = &options.String
	}
	return &a, nil
}

func (r *SQLiteRepository) CreateActivity(ctx context.Context, activity *domain.Activity) error {
	if activity.ID == "" {
		activity.ID = uuid.New().String()
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO activities (id, class_id, title, question, estimated_minutes, answer_type,
			complement_type, complement_url, week_start, created_by, options)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		activity.ID, activity.ClassID, activity.Title, activity.Question, activity.EstimatedMinutes,
		activity.AnswerType, activity.ComplementType, activity.ComplementURL, activity.WeekStart,
		activity.CreatedBy, activity.Options)
	return err
}

func (r *SQLiteRepository) UpdateActivity(ctx context.Context, activity *domain.Activity) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE activities SET title=?, question=?, estimated_minutes=?, answer_type=?,
			complement_type=?, complement_url=?, week_start=?, options=?
		WHERE id=?`,
		activity.Title, activity.Question, activity.EstimatedMinutes, activity.AnswerType,
		activity.ComplementType, activity.ComplementURL, activity.WeekStart, activity.Options, activity.ID)
	return err
}

func (r *SQLiteRepository) CreateSession(ctx context.Context, session *domain.StudySession) error {
	if session.ID == "" {
		session.ID = uuid.New().String()
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO study_sessions (id, activity_id, student_ra, started_at) VALUES (?, ?, ?, ?)`,
		session.ID, session.ActivityID, session.StudentRA, session.StartedAt.Format(time.RFC3339))
	return err
}

func (r *SQLiteRepository) GetSession(ctx context.Context, id string) (*domain.StudySession, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, activity_id, student_ra, started_at, ended_at FROM study_sessions WHERE id = ?`, id)
	var s domain.StudySession
	var startedAt string
	var endedAt sql.NullString
	if err := row.Scan(&s.ID, &s.ActivityID, &s.StudentRA, &startedAt, &endedAt); err != nil {
		return nil, err
	}
	s.StartedAt, _ = time.Parse(time.RFC3339, startedAt)
	if endedAt.Valid {
		t, _ := time.Parse(time.RFC3339, endedAt.String)
		s.EndedAt = &t
	}
	return &s, nil
}

func (r *SQLiteRepository) EndSession(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE study_sessions SET ended_at = ? WHERE id = ?`,
		time.Now().UTC().Format(time.RFC3339), id)
	return err
}

func (r *SQLiteRepository) HasSubmission(ctx context.Context, activityID, ra string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM submissions WHERE activity_id = ? AND student_ra = ?`,
		activityID, ra).Scan(&count)
	return count > 0, err
}

func (r *SQLiteRepository) CreateSubmission(ctx context.Context, submission *domain.Submission) error {
	if submission.ID == "" {
		submission.ID = uuid.New().String()
	}
	if submission.SubmittedAt.IsZero() {
		submission.SubmittedAt = time.Now().UTC()
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO submissions (id, activity_id, student_ra, session_id, answer, time_spent_seconds, submitted_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		submission.ID, submission.ActivityID, submission.StudentRA, submission.SessionID,
		submission.Answer, submission.TimeSpentSeconds, submission.SubmittedAt.Format(time.RFC3339))
	return err
}

func (r *SQLiteRepository) GetStudentDashboard(ctx context.Context, ra, weekStart string) (*domain.StudentDashboard, error) {
	dashboard := &domain.StudentDashboard{}

	rows, err := r.db.QueryContext(ctx, `
		SELECT a.id, a.title, a.estimated_minutes,
		       COALESCE((SELECT time_spent_seconds FROM submissions sub WHERE sub.activity_id = a.id AND sub.student_ra = ?), 0),
		       c.id, c.name, s.name,
		       CASE WHEN sub.id IS NOT NULL THEN 1 ELSE 0 END
		FROM activities a
		JOIN classes c ON c.id = a.class_id
		JOIN subjects s ON s.id = c.subject_id
		JOIN enrollments e ON e.class_id = c.id AND e.student_ra = ?
		LEFT JOIN submissions sub ON sub.activity_id = a.id AND sub.student_ra = ?
		WHERE a.week_start = ?
		ORDER BY c.name, a.title`, ra, ra, ra, weekStart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subjectMap := make(map[string]*domain.SubjectProgress)
	for rows.Next() {
		var activityID, title, classID, className, subjectName string
		var estimatedMinutes, timeSpent, completed int
		if err := rows.Scan(&activityID, &title, &estimatedMinutes, &timeSpent, &classID, &className, &subjectName, &completed); err != nil {
			return nil, err
		}
		dashboard.TotalActivities++
		dashboard.PlannedHours += float64(estimatedMinutes) / 60.0
		dashboard.ActualHours += float64(timeSpent) / 3600.0
		if completed == 1 {
			dashboard.CompletedActivities++
		}
		dashboard.ExerciseBreakdown = append(dashboard.ExerciseBreakdown, domain.ExerciseBreakdown{
			ActivityID: activityID, Title: title, EstimatedMinutes: estimatedMinutes, TimeSpentSeconds: timeSpent,
		})
		sp, ok := subjectMap[classID]
		if !ok {
			sp = &domain.SubjectProgress{ClassID: classID, ClassName: className, SubjectName: subjectName}
			subjectMap[classID] = sp
		}
		sp.Total++
		if completed == 1 {
			sp.Completed++
		}
	}
	for _, sp := range subjectMap {
		dashboard.Subjects = append(dashboard.Subjects, *sp)
	}
	return dashboard, rows.Err()
}

func (r *SQLiteRepository) GetClassMetrics(ctx context.Context, classID, weekStart string) ([]domain.ClassMetrics, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT s.ra, s.name,
		       COALESCE((SELECT COUNT(*) FROM submissions sub
		                 JOIN activities a ON a.id = sub.activity_id
		                 WHERE sub.student_ra = s.ra AND a.class_id = ? AND a.week_start = ?), 0),
		       COALESCE((SELECT SUM(sub.time_spent_seconds) FROM submissions sub
		                 JOIN activities a ON a.id = sub.activity_id
		                 WHERE sub.student_ra = s.ra AND a.class_id = ? AND a.week_start = ?), 0)
		FROM students s
		JOIN enrollments e ON e.student_ra = s.ra
		WHERE e.class_id = ?
		ORDER BY s.name`, classID, weekStart, classID, weekStart, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []domain.ClassMetrics
	for rows.Next() {
		var m domain.ClassMetrics
		if err := rows.Scan(&m.RA, &m.Name, &m.CompletedActivities, &m.TotalTimeSeconds); err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, rows.Err()
}

func (r *SQLiteRepository) CountTeacherStats(ctx context.Context, teacherID, weekStart string) (classes, activities, students int, err error) {
	err = r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM classes WHERE teacher_id = ?`, teacherID).Scan(&classes)
	if err != nil {
		return
	}
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM activities a JOIN classes c ON c.id = a.class_id
		WHERE c.teacher_id = ? AND a.week_start = ?`, teacherID, weekStart).Scan(&activities)
	if err != nil {
		return
	}
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT e.student_ra) FROM enrollments e
		JOIN classes c ON c.id = e.class_id WHERE c.teacher_id = ?`, teacherID).Scan(&students)
	return
}
