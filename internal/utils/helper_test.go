package utils_test

import (
	"testing"

	"github.com/kn-lim/slackingway-bot/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetStructFields(t *testing.T) {
	t.Run("Not a struct", func(t *testing.T) {
		_, err := utils.GetStructFields("not a struct")
		assert.NotNil(t, err)
	})

	t.Run("Successful GetStructFields", func(t *testing.T) {
		testStruct := struct {
			Field1 string
			Field2 int
		}{
			Field1: "test",
			Field2: 1,
		}

		expected := "Field1: test, Field2: 1"
		got, err := utils.GetStructFields(testStruct)
		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	})
}
