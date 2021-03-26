package jsonschema

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/steinfletcher/apitest"
	"github.com/xeipuuv/gojsonschema"
)

// Validate validates the http response body against the provided json schema
func Validate(schema string) apitest.Assert {
	return func(res *http.Response, req *http.Request) error {
		schemaLoader := gojsonschema.NewStringLoader(schema)
		bodyStr, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		responseBodyLoader := gojsonschema.NewBytesLoader(bodyStr)
		result, err := gojsonschema.Validate(schemaLoader, responseBodyLoader)
		if err != nil {
			return err
		}
		if !result.Valid() {
			return fmt.Errorf("invalid json schema. %s", result.Errors())
		}
		return nil
	}
}
