package ipp

func IsNotExistsError(err error) bool {
	if err == nil {
		return false
	}

	return err.Error() == "The printer or class does not exist."
}
