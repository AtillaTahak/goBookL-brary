package url

type URLRequest struct {
	URL       string `json:"url" validate:"required,url"`
	Operation string `json:"operation" validate:"required,oneof=canonical redirection all"`
}

type URLResponse struct {
	ProcessedURL string `json:"processed_url"`
}
