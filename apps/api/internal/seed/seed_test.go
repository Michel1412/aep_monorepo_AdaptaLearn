package seed

import (
	"context"
	"database/sql"
	"testing"

	"github.com/adaptalearn/api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestSeedData(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	ctx := context.Background()
	repo := repository.NewSQLiteRepository(db)
	require.NoError(t, repo.Migrate(ctx))
	require.NoError(t, Run(ctx, db))

	var teacherCount, classCount int
	require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM teachers`).Scan(&teacherCount))
	require.NoError(t, db.QueryRow(`SELECT COUNT(*) FROM classes`).Scan(&classCount))
	assert.Equal(t, 1, teacherCount)
	assert.Equal(t, 12, classCount)

	var enrolled bool
	err = db.QueryRow(`
		SELECT COUNT(*) > 0 FROM enrollments e
		JOIN classes c ON c.id = e.class_id
		WHERE e.student_ra = ? AND c.name = ?`, "23159293-2", "TeoCom-S7-B").Scan(&enrolled)
	require.NoError(t, err)
	assert.True(t, enrolled)
}
