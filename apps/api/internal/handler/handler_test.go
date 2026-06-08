package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adaptalearn/api/internal/repository"
	"github.com/adaptalearn/api/internal/seed"
	"github.com/adaptalearn/api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func setupTestRouter(t *testing.T) (*gin.Engine, string) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	ctx := context.Background()
	repo := repository.NewSQLiteRepository(db)
	require.NoError(t, repo.Migrate(ctx))
	require.NoError(t, seed.Run(ctx, db))

	authSvc := service.NewAuthService(repo)
	studentSvc := service.NewStudentService(repo)
	teacherSvc := service.NewTeacherService(repo)
	activitySvc := service.NewActivityService(repo)

	r := gin.New()
	h := New(authSvc, studentSvc, teacherSvc, activitySvc)
	h.RegisterRoutes(r)

	var teacherID string
	err = db.QueryRow(`SELECT id FROM teachers WHERE email = ?`, "maurilio@adaptalearn.local").Scan(&teacherID)
	require.NoError(t, err)
	return r, teacherID
}

func TestLoginStudent(t *testing.T) {
	r, _ := setupTestRouter(t)
	body, _ := json.Marshal(map[string]string{"ra": "23159293-2"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/student", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp["token"])
}

func TestLoginTeacher(t *testing.T) {
	r, _ := setupTestRouter(t)
	body, _ := json.Marshal(map[string]string{"email": "maurilio@adaptalearn.local", "password": "senha123"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/teacher", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestTeacherClassesCount(t *testing.T) {
	r, _ := setupTestRouter(t)
	token := loginTeacherToken(t, r)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/teachers/me/classes", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var classes []map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &classes))
	assert.Len(t, classes, 12)
}

func TestStudentDashboard(t *testing.T) {
	r, _ := setupTestRouter(t)
	token := loginStudentToken(t, r)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/students/23159293-2/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var dashboard map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &dashboard))
	assert.Equal(t, float64(3), dashboard["total_activities"])
}

func loginTeacherToken(t *testing.T, r *gin.Engine) string {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"email": "maurilio@adaptalearn.local", "password": "senha123"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/teacher", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp["token"].(string)
}

func loginStudentToken(t *testing.T, r *gin.Engine) string {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"ra": "23159293-2"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/student", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp["token"].(string)
}
