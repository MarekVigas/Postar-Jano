package resources

type StatusResponse struct {
	Status string `json:"status"`
}

func NewOkStatusResponse() StatusResponse {
	return StatusResponse{Status: "ok"}
}

func NewErrorStatusResponse() StatusResponse {
	return StatusResponse{Status: "err"}
}
