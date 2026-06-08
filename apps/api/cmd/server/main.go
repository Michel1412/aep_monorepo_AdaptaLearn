package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/adaptalearn/api/internal/handler"
	"github.com/adaptalearn/api/internal/repository"
	"github.com/adaptalearn/api/internal/seed"
	"github.com/adaptalearn/api/internal/service"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

func main() {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = filepath.Join("data", "adaptalearn.db")
	}
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	repo := repository.NewSQLiteRepository(db)
	if err := repo.Migrate(ctx); err != nil {
		log.Fatal("migrate:", err)
	}
	if err := seed.Run(ctx, db); err != nil {
		log.Fatal("seed:", err)
	}

	authSvc := service.NewAuthService(repo)
	studentSvc := service.NewStudentService(repo)
	teacherSvc := service.NewTeacherService(repo)
	activitySvc := service.NewActivityService(repo)

	r := gin.Default()
	h := handler.New(authSvc, studentSvc, teacherSvc, activitySvc)
	h.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("AdaptaLearn API running on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
