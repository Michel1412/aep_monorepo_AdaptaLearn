package seed

import (
	"context"
	"database/sql"
	"time"

	"github.com/adaptalearn/api/internal/domain"
	"github.com/adaptalearn/api/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Run(ctx context.Context, db *sql.DB) error {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM teachers`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	repo := repository.NewSQLiteRepository(db)
	hash, err := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	teacherID := uuid.New().String()
	_, err = db.ExecContext(ctx, `INSERT INTO teachers (id, name, email, password_hash) VALUES (?, ?, ?, ?)`,
		teacherID, "Maurilio", "maurilio@adaptalearn.local", string(hash))
	if err != nil {
		return err
	}

	subjectID := uuid.New().String()
	_, err = db.ExecContext(ctx, `INSERT INTO subjects (id, code, name) VALUES (?, ?, ?)`,
		subjectID, "TeoCom", "Teoria da Computação")
	if err != nil {
		return err
	}

	students := []domain.Student{
		{RA: "23159293-2", Name: "Vinicius Reginaldo Ferrarini"},
		{RA: "23034350-2", Name: "João Pedro Souza Peixoto Saraiva"},
		{RA: "23220783-2", Name: "Michel Bocchi Junior"},
	}
	for _, s := range students {
		if err := repo.CreateStudent(ctx, &s); err != nil {
			return err
		}
	}

	classIDs := make(map[string]string)
	series := []int{4, 5, 6, 7}
	turmas := []string{"A", "B", "C"}
	for _, s := range series {
		for _, t := range turmas {
			name := formatClassName("TeoCom", s, t)
			id := uuid.New().String()
			classIDs[name] = id
			_, err = db.ExecContext(ctx, `INSERT INTO classes (id, name, subject_id, teacher_id, series, turma) VALUES (?, ?, ?, ?, ?, ?)`,
				id, name, subjectID, teacherID, s, t)
			if err != nil {
				return err
			}
		}
	}

	enrollments := map[string]string{
		"TeoCom-S7-B": "23159293-2",
		"TeoCom-S7-A": "23034350-2",
		"TeoCom-S6-C": "23220783-2",
	}
	for className, ra := range enrollments {
		_, err = db.ExecContext(ctx, `INSERT INTO enrollments (id, student_ra, class_id) VALUES (?, ?, ?)`,
			uuid.New().String(), ra, classIDs[className])
		if err != nil {
			return err
		}
	}

	weekStart := currentWeekStart()
	s7bID := classIDs["TeoCom-S7-B"]
	videoURL := "https://www.youtube.com/watch?v=example"
	linkURL := "https://en.wikipedia.org/wiki/Regular_expression"
	activities := []domain.Activity{
		{
			ClassID: s7bID, Title: "Autômatos finitos determinísticos",
			Question: "Defina autômato finito determinístico (AFD) e dê um exemplo prático.",
			EstimatedMinutes: 30, AnswerType: domain.AnswerEssay,
			ComplementType: domain.ComplementVideo, ComplementURL: &videoURL,
			WeekStart: weekStart, CreatedBy: teacherID,
		},
		{
			ClassID: s7bID, Title: "Linguagem regular e expressões",
			Question: "O que é uma linguagem regular? Cite um exemplo de expressão regular.",
			EstimatedMinutes: 20, AnswerType: domain.AnswerShortText,
			ComplementType: domain.ComplementLink, ComplementURL: &linkURL,
			WeekStart: weekStart, CreatedBy: teacherID,
		},
		{
			ClassID: s7bID, Title: "Máquina de Turing — conceitos",
			Question: "Uma máquina de Turing pode simular qualquer algoritmo computável.",
			EstimatedMinutes: 15, AnswerType: domain.AnswerTrueFalse,
			ComplementType: domain.ComplementNone,
			WeekStart: weekStart, CreatedBy: teacherID,
		},
	}
	for _, a := range activities {
		if err := repo.CreateActivity(ctx, &a); err != nil {
			return err
		}
	}
	return nil
}

func formatClassName(code string, series int, turma string) string {
	return code + "-S" + itoa(series) + "-" + turma
}

func itoa(n int) string {
	if n == 4 {
		return "4"
	}
	if n == 5 {
		return "5"
	}
	if n == 6 {
		return "6"
	}
	return "7"
}

func currentWeekStart() string {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := now.AddDate(0, 0, -(weekday - 1))
	return monday.Format("2006-01-02")
}

func CurrentWeekStart() string {
	return currentWeekStart()
}
