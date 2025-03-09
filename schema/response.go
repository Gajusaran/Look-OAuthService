package schema

type SuccessResponse struct {
	Success    bool `json:"success"`
	Payload    any  `json:"data"`
	StatusCode int  `json:"status"`
}

type FailureResponse struct {
	Success    bool `json:"success"`
	StatusCode int  `json:"status"`
}
