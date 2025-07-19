package dto

import "time"

type MovieListResponse struct {
	Movies []Movie `json:"movies"`
}
type Movie struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseYear time.Time `json:"release_year"`
}
