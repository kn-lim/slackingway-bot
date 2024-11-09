package utils_test

import (
	"testing"

	"github.com/kn-lim/slackingway-bot/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestPrintStructFields(t *testing.T) {
	t.Run("Not a struct", func(t *testing.T) {
		_, err := utils.PrintStructFields("not a struct")
		assert.NotNil(t, err)
	})

	t.Run("Successful GetStructFields", func(t *testing.T) {
		testStruct := struct {
			Field1 string `json:"field1"`
			Field2 int    `json:"field2"`
		}{
			Field1: "test",
			Field2: 1,
		}

		expected := `{"field1":"test","field2":1}`
		got, err := utils.PrintStructFields(testStruct)
		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	})
}
