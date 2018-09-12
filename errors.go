package golang_commons

type ErrorResponse struct {
	Message string   `json:"message"`
	ExtCode int      `json:"extCode"`
	Errors  []string `json:"errors"`
	Servlet string   `json:"servlet"`
}
