package ipp

import "fmt"

// IsNotExistsError checks a given error whether a printer or class does not exist
func IsNotExistsError(err error) bool {
	if err == nil {
		return false
	}

	return err.Error() == "The printer or class does not exist."
}

// IPPError used for non ok ipp status codes
type IPPError struct {
	Status  int16
	Message string
}

func (e IPPError) Error() string {
	return fmt.Sprintf("ipp status: %d, message: %s", e.Status, e.Message)
}

// HTTPError used for non 200 http codes
type HTTPError struct {
	Code int
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("got http code %d", e.Code)
}
