package passgen

import "errors"

var errMemoNotImplemented = errors.New("memo mode not implemented yet")

func generateMemo(_ Options) (string, error) {
	return "", errMemoNotImplemented
}
