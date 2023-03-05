package main

var ErrInvalid = apiError{
	Err:    "invalid",
	Status: 400,
}

type apiError struct {
	Err    string
	Status int
}

func (e apiError) Error() string {
	return e.Err
}
