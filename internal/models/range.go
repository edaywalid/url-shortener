package models

type Request struct {
	OriginalURL string `json:"original_url"`
}

type Response struct {
	ShortURL string `json:"short_url"`
}

type Range struct {
	Start   uint64 `json:"start"`
	End     uint64 `json:"end"`
	Current uint64 `json:"current"`
}
