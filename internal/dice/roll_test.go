package dice_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kn-lim/slackingway-bot/internal/dice"
)

func TestRoll(t *testing.T) {
	testInput := "2d6+3"

	resultString, resultInt, err := dice.Roll(testInput)

	assert.Nil(t, err)
	assert.NotEmpty(t, resultString)
	assert.NotZero(t, resultInt)
}
