package golang_commons

const SUCCESS_CODE = 0
const DEFAULT_ERROR_CODE = -1

type ErrorResponse struct {
	Message string   `json:"message"`
	ExtCode int      `json:"extCode"`
	Errors  []string `json:"errors"`
	Servlet string   `json:"servlet"`
}
