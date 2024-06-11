package utils

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

func WrapleRestyError(resp *resty.Response, err error) error {
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
