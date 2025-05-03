package dto

type SummaryReportResponse struct {
	TotalEvents  int     `json:"total_events"`
	TotalTickets int     `json:"total_tickets"`
	TotalRevenue float64 `json:"total_revenue"`
	Upcoming     int     `json:"upcoming_events"`
	Ongoing      int     `json:"ongoing_events"`
	Completed    int     `json:"completed_events"`
}

type EventReportResponse struct {
	EventName     string  `json:"event_name"`
	TotalCapacity int     `json:"total_capacity"`
	TicketsSold   int     `json:"tickets_sold"`
	Revenue       float64 `json:"revenue"`
	OccupancyRate float64 `json:"occupancy_rate"`
}
