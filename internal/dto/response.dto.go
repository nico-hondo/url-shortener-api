package dto

type Response struct {
	Success bool
	Data    any
	Msg     string
}

type ErrorResponse struct {
	Success bool
	Data    any
	Msg     string
}
