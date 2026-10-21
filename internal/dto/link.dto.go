package dto

type CreateLinkRequest struct {
	OriginalURL string  `json:"original_url" binding:"required,url" example:"https://www.example.com/very/long/url"`
	Slug        *string `json:"slug,omitempty" binding:"omitempty,alphanum", example:"abc123"`
}
type CreateLinkResponse struct {
	ID          int    `json:"id"`
	OriginalURL string `json:"original_url"`
	Slug        string `json:"slug"`
	ShortURL    string `json:"short_url"`
}
