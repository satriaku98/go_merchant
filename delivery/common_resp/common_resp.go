package common_resp

type SuccessResponse struct {
	Success bool        `json:"sucess"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var (
	statusSuccess = true
	statusFailed  = false
	successMsg    = "Success"
	failedMsg     = "Failed"
)

func SuccessMessage(message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Success: statusSuccess,
		Message: message,
		Data:    data,
	}
}

func FailedMessage(message string) *ErrorResponse {
	return &ErrorResponse{
		Success: statusFailed,
		Message: message,
	}
}
