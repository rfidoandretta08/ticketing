package dto

type EventRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Location    string  `json:"location" binding:"required"`
	DateTime    string  `json:"date_time" binding:"required"`
	Capacity    int     `json:"capacity" binding:"required,min=1"`
	Price       float64 `json:"price" binding:"required,min=0"`
}

type EventResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	DateTime    string  `json:"date_time"`
	Capacity    int     `json:"capacity"`
	Available   int     `json:"available"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
}
