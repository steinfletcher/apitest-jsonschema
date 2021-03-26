package jsonschema_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/steinfletcher/apitest"
	jsonschema "github.com/steinfletcher/apitest-jsonschema"
	"github.com/steinfletcher/apitest/mocks"
)

func TestValidate_MatchesSchema(t *testing.T) {
	apitest.New().
		HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`{
			  "firstName": "John",
			  "lastName": "Doe",
			  "age": 21
			}`))
			writer.WriteHeader(http.StatusOK)
		}).
		Get("/").
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonschema.Validate(schema)).
		End()
}

func TestValidate_FailsToMatchSchema(t *testing.T) {
	mockVerifier := mocks.MockVerifier{
		NoErrorFn: func(t apitest.TestingT, err error, msgAndArgs ...interface{}) bool {
			if err == nil {
				t.Fatal("expected an error")
				return false
			}
			if !strings.Contains(err.Error(), "firstName is required") {
				t.Fatal("unexpected error")
			}
			return false
		},
	}

	apitest.New().
		Verifier(mockVerifier).
		HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			_, _ = writer.Write([]byte(`{
			  "firstNames": "John",
			  "lastName": "Doe",
			  "age": 21
			}`))
			writer.WriteHeader(http.StatusOK)
		}).
		Get("/").
		Expect(t).
		Assert(jsonschema.Validate(schema)).
		End()
}

const schema = `{
  "$id": "https://example.com/person.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Person",
  "type": "object",
  "required": [ "firstName", "lastName", "age" ],
  "properties": {
    "firstName": {
      "type": "string",
      "description": "The person's first name."
    },
    "lastName": {
      "type": "string",
      "description": "The person's last name."
    },
    "age": {
      "description": "Age in years which must be equal to or greater than zero.",
      "type": "integer",
      "minimum": 0
    }
  }
}`
