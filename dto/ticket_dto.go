package dto

type TicketRequest struct {
	EventID uint `json:"event_id" binding:"required"`
}

type TicketResponse struct {
	ID          uint    `json:"id"`
	EventName   string  `json:"event_name"`
	EventDate   string  `json:"event_date"`
	Location    string  `json:"location"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
	BookingDate string  `json:"booking_date"`
}
