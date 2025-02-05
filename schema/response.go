package schema

type SuccessResponse struct {
	Success    bool   `json:"success"`
	Payload    any    `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
}

type FailureResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
}
