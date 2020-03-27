package ipp

import "fmt"

func IsNotExistsError(err error) bool {
	if err == nil {
		return false
	}

	return err.Error() == "The printer or class does not exist."
}

type IPPError struct {
	Status  int16
	Message string
}

func (e IPPError) Error() string {
	return fmt.Sprintf("ipp status: %d, message: %s", e.Status, e.Message)
}

type HTTPError struct {
	Code int
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("got http code %d", e.Code)
}
