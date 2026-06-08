package domain

import "time"

type AnswerType string

const (
	AnswerMultipleChoice AnswerType = "multiple_choice"
	AnswerTrueFalse      AnswerType = "true_false"
	AnswerShortText      AnswerType = "short_text"
	AnswerEssay          AnswerType = "essay"
)

type ComplementType string

const (
	ComplementNone  ComplementType = "none"
	ComplementVideo ComplementType = "video"
	ComplementImage ComplementType = "image"
	ComplementLink  ComplementType = "link"
)

type Teacher struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type Student struct {
	RA   string `json:"ra"`
	Name string `json:"name"`
}

type Subject struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type Class struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	SubjectID   string `json:"subject_id"`
	TeacherID   string `json:"teacher_id"`
	Series      int    `json:"series"`
	Turma       string `json:"turma"`
	SubjectCode string `json:"subject_code,omitempty"`
	SubjectName string `json:"subject_name,omitempty"`
	StudentCount int   `json:"student_count,omitempty"`
	ActivityCount int  `json:"activity_count,omitempty"`
}

type Enrollment struct {
	ID        string `json:"id"`
	StudentRA string `json:"student_ra"`
	ClassID   string `json:"class_id"`
}

type Activity struct {
	ID               string         `json:"id"`
	ClassID          string         `json:"class_id"`
	Title            string         `json:"title"`
	Question         string         `json:"question"`
	EstimatedMinutes int            `json:"estimated_minutes"`
	AnswerType       AnswerType     `json:"answer_type"`
	ComplementType   ComplementType `json:"complement_type"`
	ComplementURL    *string        `json:"complement_url"`
	WeekStart        string         `json:"week_start"`
	CreatedBy        string         `json:"created_by"`
	Options          *string        `json:"options"`
	Status           string         `json:"status,omitempty"`
}

type StudySession struct {
	ID         string     `json:"id"`
	ActivityID string     `json:"activity_id"`
	StudentRA  string     `json:"student_ra"`
	StartedAt  time.Time  `json:"started_at"`
	EndedAt    *time.Time `json:"ended_at,omitempty"`
}

type Submission struct {
	ID               string    `json:"id"`
	ActivityID       string    `json:"activity_id"`
	StudentRA        string    `json:"student_ra"`
	SessionID        string    `json:"session_id"`
	Answer           string    `json:"answer"`
	TimeSpentSeconds int       `json:"time_spent_seconds"`
	SubmittedAt      time.Time `json:"submitted_at"`
}

type SubjectProgress struct {
	ClassID     string `json:"class_id"`
	ClassName   string `json:"class_name"`
	SubjectName string `json:"subject_name"`
	Completed   int    `json:"completed"`
	Total       int    `json:"total"`
}

type ExerciseBreakdown struct {
	ActivityID       string `json:"activity_id"`
	Title            string `json:"title"`
	EstimatedMinutes int    `json:"estimated_minutes"`
	TimeSpentSeconds int    `json:"time_spent_seconds"`
}

type StudentDashboard struct {
	CompletedActivities int                 `json:"completed_activities"`
	TotalActivities     int                 `json:"total_activities"`
	PlannedHours        float64             `json:"planned_hours"`
	ActualHours         float64             `json:"actual_hours"`
	Subjects            []SubjectProgress   `json:"subjects"`
	ExerciseBreakdown   []ExerciseBreakdown `json:"exercise_breakdown"`
}

type ClassMetrics struct {
	RA                  string `json:"ra"`
	Name                string `json:"name"`
	CompletedActivities int    `json:"completed_activities"`
	TotalTimeSeconds    int    `json:"total_time_seconds"`
}

type CreateActivityInput struct {
	Title            string         `json:"title"`
	Question         string         `json:"question"`
	EstimatedMinutes int            `json:"estimated_minutes"`
	AnswerType       AnswerType     `json:"answer_type"`
	ComplementType   ComplementType `json:"complement_type"`
	ComplementURL    *string        `json:"complement_url"`
	WeekStart        string         `json:"week_start"`
	Options          *string        `json:"options"`
}

type SubmitInput struct {
	Answer           string `json:"answer"`
	TimeSpentSeconds int    `json:"time_spent_seconds"`
}

const MinSessionSeconds = 10
const MaxSessionSeconds = 4 * 60 * 60
