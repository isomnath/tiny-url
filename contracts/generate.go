package contracts

type GenerateRequest struct {
	OriginalURL string `json:"original_url"`
}

type GenerateResponse struct {
	TinyURL string `json:"tiny_url"`
}
