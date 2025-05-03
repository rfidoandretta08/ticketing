package dto

type TicketRequest struct {
	EventID uint `json:"event_id" binding:"required"`
	Qty     int  `json:"qty" binding:"required"`
}

type PaymentUpdateResponse struct {
	ID            uint   `json:"id"`
	Status        string `json:"status"`
	PaymentStatus string `json:"payment_status"`
}

type TicketResponse struct {
	ID            uint    `json:"id"`
	EventName     string  `json:"event_name"`
	EventDate     string  `json:"event_date"`
	Location      string  `json:"location"`
	Price         float64 `json:"price"`
	Status        string  `json:"status"`
	PaymentStatus string  `json:"payment_status"`
	BookingDate   string  `json:"booking_date"`
	Qty           int     `json:"quantity"`  // Menyertakan Quantity
	SubTotal      float64 `json:"sub_total"` // Menyertakan SubTotal
}
