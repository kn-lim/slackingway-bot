package utils_test

import (
	"net/url"
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

func TestPrintPayloadFields(t *testing.T) {
	t.Run("Valid payload", func(t *testing.T) {
		// Create a valid URL-encoded JSON payload
		payload := url.QueryEscape(`{"field1":"value1","field2":2}`)

		expected := `{"field1":"value1","field2":2}`
		got, err := utils.PrintPayloadFields(payload)
		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("Invalid URL encoding", func(t *testing.T) {
		// Create an invalid URL-encoded payload
		payload := "%"

		_, err := utils.PrintPayloadFields(payload)
		assert.NotNil(t, err)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		// Create a valid URL-encoded but invalid JSON payload
		payload := url.QueryEscape(`{"field1":"value1","field2":}`)

		_, err := utils.PrintPayloadFields(payload)
		assert.NotNil(t, err)
	})
}
