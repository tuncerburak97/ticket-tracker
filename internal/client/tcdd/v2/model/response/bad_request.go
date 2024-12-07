package response

type BadRequestResponse struct {
	Timestamp string  `json:"timestamp"`
	TraceID   string  `json:"traceId"`
	Status    string  `json:"status"`
	Type      string  `json:"type"`
	Code      int     `json:"code"`
	Message   string  `json:"message"`
	Detail    string  `json:"detail"`
	Title     *string `json:"title"` // Use a pointer to allow null values
}
