package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/adaptalearn/api/internal/domain"
	"github.com/adaptalearn/api/internal/seed"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo domain.Repository
}

func NewAuthService(repo domain.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) LoginStudent(ctx context.Context, ra string) (*domain.Student, error) {
	student, err := s.repo.GetStudentByRA(ctx, ra)
	if err != nil {
		return nil, errors.New("RA não encontrado")
	}
	return student, nil
}

func (s *AuthService) LoginTeacher(ctx context.Context, email, password string) (*domain.Teacher, error) {
	teacher, err := s.repo.GetTeacherByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(teacher.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("credenciais inválidas")
	}
	return teacher, nil
}

type StudentService struct {
	repo domain.Repository
}

func NewStudentService(repo domain.Repository) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) GetDashboard(ctx context.Context, ra string) (*domain.StudentDashboard, error) {
	return s.repo.GetStudentDashboard(ctx, ra, seed.CurrentWeekStart())
}

func (s *StudentService) GetClasses(ctx context.Context, ra string) ([]domain.Class, error) {
	return s.repo.GetClassesByStudent(ctx, ra)
}

func (s *StudentService) GetActivities(ctx context.Context, ra, classID string) ([]domain.Activity, error) {
	enrolled, err := s.repo.IsStudentEnrolled(ctx, classID, ra)
	if err != nil || !enrolled {
		return nil, errors.New("turma não encontrada ou aluno não matriculado")
	}
	activities, err := s.repo.GetActivitiesByClass(ctx, classID, seed.CurrentWeekStart())
	if err != nil {
		return nil, err
	}
	for i := range activities {
		has, _ := s.repo.HasSubmission(ctx, activities[i].ID, ra)
		if has {
			activities[i].Status = "completed"
		} else {
			activities[i].Status = "pending"
		}
	}
	return activities, nil
}

type ActivityService struct {
	repo domain.Repository
}

func NewActivityService(repo domain.Repository) *ActivityService {
	return &ActivityService{repo: repo}
}

func (s *ActivityService) StartSession(ctx context.Context, activityID, ra string) (*domain.StudySession, error) {
	activity, err := s.repo.GetActivityByID(ctx, activityID)
	if err != nil {
		return nil, errors.New("atividade não encontrada")
	}
	enrolled, err := s.repo.IsStudentEnrolled(ctx, activity.ClassID, ra)
	if err != nil || !enrolled {
		return nil, errors.New("aluno não matriculado nesta turma")
	}
	has, err := s.repo.HasSubmission(ctx, activityID, ra)
	if err != nil {
		return nil, err
	}
	if has {
		return nil, errors.New("atividade já concluída")
	}
	session := &domain.StudySession{ActivityID: activityID, StudentRA: ra}
	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *ActivityService) Submit(ctx context.Context, activityID, sessionID, ra string, input domain.SubmitInput) (*domain.Submission, error) {
	if input.Answer == "" {
		return nil, errors.New("resposta obrigatória")
	}
	if input.TimeSpentSeconds < domain.MinSessionSeconds || input.TimeSpentSeconds > domain.MaxSessionSeconds {
		return nil, fmt.Errorf("tempo inválido: deve estar entre %d e %d segundos", domain.MinSessionSeconds, domain.MaxSessionSeconds)
	}
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, errors.New("sessão não encontrada")
	}
	if session.ActivityID != activityID || session.StudentRA != ra {
		return nil, errors.New("sessão inválida")
	}
	if session.EndedAt != nil {
		return nil, errors.New("sessão já encerrada")
	}
	has, err := s.repo.HasSubmission(ctx, activityID, ra)
	if err != nil {
		return nil, err
	}
	if has {
		return nil, errors.New("atividade já concluída")
	}
	if err := s.repo.EndSession(ctx, sessionID); err != nil {
		return nil, err
	}
	submission := &domain.Submission{
		ActivityID: activityID, StudentRA: ra, SessionID: sessionID,
		Answer: input.Answer, TimeSpentSeconds: input.TimeSpentSeconds,
	}
	if err := s.repo.CreateSubmission(ctx, submission); err != nil {
		return nil, err
	}
	return submission, nil
}

func ValidateActivityInput(input domain.CreateActivityInput) error {
	if input.Title == "" || input.Question == "" {
		return errors.New("título e pergunta são obrigatórios")
	}
	if input.EstimatedMinutes <= 0 {
		return errors.New("tempo estimado deve ser positivo")
	}
	validTypes := map[domain.AnswerType]bool{
		domain.AnswerMultipleChoice: true, domain.AnswerTrueFalse: true,
		domain.AnswerShortText: true, domain.AnswerEssay: true,
	}
	if !validTypes[input.AnswerType] {
		return errors.New("tipo de resposta inválido")
	}
	return nil
}

type TeacherService struct {
	repo domain.Repository
}

func NewTeacherService(repo domain.Repository) *TeacherService {
	return &TeacherService{repo: repo}
}

func (s *TeacherService) GetClasses(ctx context.Context, teacherID string) ([]domain.Class, error) {
	return s.repo.GetClassesByTeacher(ctx, teacherID)
}

func (s *TeacherService) GetStudents(ctx context.Context, teacherID, classID string) ([]domain.Student, error) {
	if err := s.ensureClassOwner(ctx, teacherID, classID); err != nil {
		return nil, err
	}
	return s.repo.GetStudentsByClass(ctx, classID)
}

func (s *TeacherService) EnrollStudent(ctx context.Context, teacherID, classID, ra string) (*domain.Student, error) {
	if err := s.ensureClassOwner(ctx, teacherID, classID); err != nil {
		return nil, err
	}
	return s.repo.EnrollStudent(ctx, classID, ra)
}

func (s *TeacherService) CreateActivity(ctx context.Context, teacherID, classID string, input domain.CreateActivityInput) (*domain.Activity, error) {
	if err := ValidateActivityInput(input); err != nil {
		return nil, err
	}
	if err := s.ensureClassOwner(ctx, teacherID, classID); err != nil {
		return nil, err
	}
	activity := &domain.Activity{
		ClassID: classID, Title: input.Title, Question: input.Question,
		EstimatedMinutes: input.EstimatedMinutes, AnswerType: input.AnswerType,
		ComplementType: input.ComplementType, ComplementURL: input.ComplementURL,
		WeekStart: input.WeekStart, CreatedBy: teacherID, Options: input.Options,
	}
	if activity.ComplementType == "" {
		activity.ComplementType = domain.ComplementNone
	}
	if err := s.repo.CreateActivity(ctx, activity); err != nil {
		return nil, err
	}
	return activity, nil
}

func (s *TeacherService) UpdateActivity(ctx context.Context, teacherID, activityID string, input domain.CreateActivityInput) (*domain.Activity, error) {
	if err := ValidateActivityInput(input); err != nil {
		return nil, err
	}
	activity, err := s.repo.GetActivityByID(ctx, activityID)
	if err != nil {
		return nil, errors.New("atividade não encontrada")
	}
	if err := s.ensureClassOwner(ctx, teacherID, activity.ClassID); err != nil {
		return nil, err
	}
	activity.Title = input.Title
	activity.Question = input.Question
	activity.EstimatedMinutes = input.EstimatedMinutes
	activity.AnswerType = input.AnswerType
	activity.ComplementType = input.ComplementType
	activity.ComplementURL = input.ComplementURL
	activity.WeekStart = input.WeekStart
	activity.Options = input.Options
	if err := s.repo.UpdateActivity(ctx, activity); err != nil {
		return nil, err
	}
	return activity, nil
}

func (s *TeacherService) GetMetrics(ctx context.Context, teacherID, classID string) ([]domain.ClassMetrics, error) {
	if err := s.ensureClassOwner(ctx, teacherID, classID); err != nil {
		return nil, err
	}
	return s.repo.GetClassMetrics(ctx, classID, seed.CurrentWeekStart())
}

func (s *TeacherService) GetDashboardStats(ctx context.Context, teacherID string) (classes, activities, students int, err error) {
	return s.repo.CountTeacherStats(ctx, teacherID, seed.CurrentWeekStart())
}

func (s *TeacherService) ensureClassOwner(ctx context.Context, teacherID, classID string) error {
	class, err := s.repo.GetClassByID(ctx, classID)
	if err != nil {
		return errors.New("turma não encontrada")
	}
	if class.TeacherID != teacherID {
		return errors.New("sem permissão para esta turma")
	}
	return nil
}
