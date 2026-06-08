package domain

import "context"

type Repository interface {
	Migrate(ctx context.Context) error

	GetTeacherByEmail(ctx context.Context, email string) (*Teacher, error)
	GetTeacherByID(ctx context.Context, id string) (*Teacher, error)
	GetStudentByRA(ctx context.Context, ra string) (*Student, error)
	CreateStudent(ctx context.Context, student *Student) error

	GetClassesByTeacher(ctx context.Context, teacherID string) ([]Class, error)
	GetClassByID(ctx context.Context, id string) (*Class, error)
	GetClassesByStudent(ctx context.Context, ra string) ([]Class, error)

	GetStudentsByClass(ctx context.Context, classID string) ([]Student, error)
	EnrollStudent(ctx context.Context, classID, ra string) (*Student, error)
	IsStudentEnrolled(ctx context.Context, classID, ra string) (bool, error)

	GetActivitiesByClass(ctx context.Context, classID, weekStart string) ([]Activity, error)
	GetActivityByID(ctx context.Context, id string) (*Activity, error)
	CreateActivity(ctx context.Context, activity *Activity) error
	UpdateActivity(ctx context.Context, activity *Activity) error

	CreateSession(ctx context.Context, session *StudySession) error
	GetSession(ctx context.Context, id string) (*StudySession, error)
	EndSession(ctx context.Context, id string) error
	HasSubmission(ctx context.Context, activityID, ra string) (bool, error)
	CreateSubmission(ctx context.Context, submission *Submission) error

	GetStudentDashboard(ctx context.Context, ra, weekStart string) (*StudentDashboard, error)
	GetClassMetrics(ctx context.Context, classID, weekStart string) ([]ClassMetrics, error)
	CountTeacherStats(ctx context.Context, teacherID, weekStart string) (classes, activities, students int, err error)
}
