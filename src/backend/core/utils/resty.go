package utils

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

// WrapRestyError wraps a Resty response error into a standard Go error.
// The function takes a pointer to a Resty response and an error as input.
// If the error is not nil, it returns the error.
// If the response is not an error, it returns nil.
// If the response contains an error message, it creates a new error using fmt.Errorf and returns it.
// Otherwise, it returns nil.
func WrapRestyError(resp *resty.Response, err error) error {
	if err != nil {
		return err
	}
	if !resp.IsError() {
		return nil
	}
	errMsg := string(resp.Body())
	if errMsg != "" {
		return fmt.Errorf(errMsg)
	}
	return nil
}
