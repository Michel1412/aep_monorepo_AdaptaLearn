package handler

import (
	"net/http"

	"github.com/adaptalearn/api/internal/domain"
	"github.com/adaptalearn/api/internal/middleware"
	"github.com/adaptalearn/api/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	auth    *service.AuthService
	student *service.StudentService
	teacher *service.TeacherService
	activity *service.ActivityService
}

func New(auth *service.AuthService, student *service.StudentService, teacher *service.TeacherService, activity *service.ActivityService) *Handler {
	return &Handler{auth: auth, student: student, teacher: teacher, activity: activity}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.Use(middleware.CORS())
	api := r.Group("/api/v1")

	api.POST("/auth/student", h.loginStudent)
	api.POST("/auth/teacher", h.loginTeacher)

	students := api.Group("/students/:ra", middleware.AuthRequired(), middleware.RequireRole("student"))
	{
		students.GET("/dashboard", h.getDashboard)
		students.GET("/classes", h.getStudentClasses)
		students.GET("/classes/:classId/activities", h.getStudentActivities)
	}

	activities := api.Group("/activities", middleware.AuthRequired())
	{
		activities.POST("/:id/sessions/start", middleware.RequireRole("student"), h.startSession)
		activities.POST("/:id/sessions/:sessionId/submit", middleware.RequireRole("student"), h.submitActivity)
		activities.PUT("/:id", middleware.RequireRole("teacher"), h.updateActivity)
	}

	teachers := api.Group("/teachers/me", middleware.AuthRequired(), middleware.RequireRole("teacher"))
	{
		teachers.GET("/dashboard", h.getTeacherDashboard)
		teachers.GET("/classes", h.getTeacherClasses)
		teachers.GET("/classes/:classId/students", h.getClassStudents)
		teachers.POST("/classes/:classId/students", h.enrollStudent)
		teachers.GET("/classes/:classId/metrics", h.getClassMetrics)
		teachers.POST("/classes/:classId/activities", h.createActivity)
	}
}

func (h *Handler) loginStudent(c *gin.Context) {
	var req struct {
		RA string `json:"ra" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RA obrigatório"})
		return
	}
	student, err := h.auth.LoginStudent(c.Request.Context(), req.RA)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	token, _ := middleware.GenerateToken("student", student.RA)
	c.JSON(http.StatusOK, gin.H{"token": token, "student": student})
}

func (h *Handler) loginTeacher(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email e senha obrigatórios"})
		return
	}
	teacher, err := h.auth.LoginTeacher(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	token, _ := middleware.GenerateToken("teacher", teacher.ID)
	c.JSON(http.StatusOK, gin.H{"token": token, "teacher": gin.H{"id": teacher.ID, "name": teacher.Name, "email": teacher.Email}})
}

func (h *Handler) getDashboard(c *gin.Context) {
	ra := c.Param("ra")
	if c.GetString("subjectID") != ra {
		c.JSON(http.StatusForbidden, gin.H{"error": "sem permissão"})
		return
	}
	dashboard, err := h.student.GetDashboard(c.Request.Context(), ra)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dashboard)
}

func (h *Handler) getStudentClasses(c *gin.Context) {
	ra := c.Param("ra")
	if c.GetString("subjectID") != ra {
		c.JSON(http.StatusForbidden, gin.H{"error": "sem permissão"})
		return
	}
	classes, err := h.student.GetClasses(c.Request.Context(), ra)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classes)
}

func (h *Handler) getStudentActivities(c *gin.Context) {
	ra := c.Param("ra")
	classID := c.Param("classId")
	if c.GetString("subjectID") != ra {
		c.JSON(http.StatusForbidden, gin.H{"error": "sem permissão"})
		return
	}
	activities, err := h.student.GetActivities(c.Request.Context(), ra, classID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activities)
}

func (h *Handler) startSession(c *gin.Context) {
	activityID := c.Param("id")
	ra := c.GetString("subjectID")
	session, err := h.activity.StartSession(c.Request.Context(), activityID, ra)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"session_id": session.ID, "started_at": session.StartedAt})
}

func (h *Handler) submitActivity(c *gin.Context) {
	activityID := c.Param("id")
	sessionID := c.Param("sessionId")
	ra := c.GetString("subjectID")
	var input domain.SubmitInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido"})
		return
	}
	submission, err := h.activity.Submit(c.Request.Context(), activityID, sessionID, ra, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"submission_id": submission.ID, "time_spent_seconds": submission.TimeSpentSeconds,
		"submitted_at": submission.SubmittedAt,
	})
}

func (h *Handler) getTeacherDashboard(c *gin.Context) {
	teacherID := c.GetString("subjectID")
	classes, activities, students, err := h.teacher.GetDashboardStats(c.Request.Context(), teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"classes": classes, "activities": activities, "students": students})
}

func (h *Handler) getTeacherClasses(c *gin.Context) {
	teacherID := c.GetString("subjectID")
	classes, err := h.teacher.GetClasses(c.Request.Context(), teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classes)
}

func (h *Handler) getClassStudents(c *gin.Context) {
	teacherID := c.GetString("subjectID")
	classID := c.Param("classId")
	students, err := h.teacher.GetStudents(c.Request.Context(), teacherID, classID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, students)
}

func (h *Handler) enrollStudent(c *gin.Context) {
	teacherID := c.GetString("subjectID")
	classID := c.Param("classId")
	var req struct {
		RA string `json:"ra" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "RA obrigatório"})
		return
	}
	student, err := h.teacher.EnrollStudent(c.Request.Context(), teacherID, classID, req.RA)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, student)
}

func (h *Handler) getClassMetrics(c *gin.Context) {
	teacherID := c.GetString("subjectID")
	classID := c.Param("classId")
	metrics, err := h.teacher.GetMetrics(c.Request.Context(), teacherID, classID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, metrics)
}

func (h *Handler) createActivity(c *gin.Context) {
	teacherID := c.GetString("subjectID")
	classID := c.Param("classId")
	var input domain.CreateActivityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido"})
		return
	}
	activity, err := h.teacher.CreateActivity(c.Request.Context(), teacherID, classID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, activity)
}

func (h *Handler) updateActivity(c *gin.Context) {
	teacherID := c.GetString("subjectID")
	activityID := c.Param("id")
	var input domain.CreateActivityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payload inválido"})
		return
	}
	activity, err := h.teacher.UpdateActivity(c.Request.Context(), teacherID, activityID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activity)
}
