package service

import (
	"testing"

	"github.com/adaptalearn/api/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestValidateActivityInput(t *testing.T) {
	valid := domain.CreateActivityInput{
		Title: "Test", Question: "Q?", EstimatedMinutes: 10,
		AnswerType: domain.AnswerEssay,
	}
	assert.NoError(t, ValidateActivityInput(valid))

	invalid := domain.CreateActivityInput{Title: "", Question: "Q", EstimatedMinutes: 10, AnswerType: domain.AnswerEssay}
	assert.Error(t, ValidateActivityInput(invalid))

	invalidTime := domain.CreateActivityInput{Title: "T", Question: "Q", EstimatedMinutes: 0, AnswerType: domain.AnswerEssay}
	assert.Error(t, ValidateActivityInput(invalidTime))
}

func TestSubmitTimeValidation(t *testing.T) {
	input := domain.SubmitInput{Answer: "test", TimeSpentSeconds: 5}
	assert.True(t, input.TimeSpentSeconds < domain.MinSessionSeconds)

	input.TimeSpentSeconds = 100
	assert.True(t, input.TimeSpentSeconds >= domain.MinSessionSeconds)
	assert.True(t, input.TimeSpentSeconds <= domain.MaxSessionSeconds)
}
