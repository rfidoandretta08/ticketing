package controller

import (
	"net/http"
	"strconv"

	"ticketing/dto"
	"ticketing/middleware"
	"ticketing/service"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	ticketService service.TicketService
}

func NewTicketController(ticketService service.TicketService) *TicketController {
	return &TicketController{ticketService: ticketService}
}

func (c *TicketController) PurchaseTicket(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.TicketRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := c.ticketService.PurchaseTicket(userID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, ticket)
}

func (c *TicketController) GetUserTickets(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	tickets, pagination, err := c.ticketService.GetUserTickets(userID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       tickets,
		"pagination": pagination,
	})
}

func (c *TicketController) GetTicketByID(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ticketID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket ID"})
		return
	}

	ticket, err := c.ticketService.GetTicketByID(userID, uint(ticketID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ticket)
}

func (c *TicketController) CancelTicket(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ticketID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket ID"})
		return
	}

	if err := c.ticketService.CancelTicket(userID, uint(ticketID)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ticket cancelled successfully"})
}

func (c *TicketController) UpdatePayment(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	ticketID, _ := strconv.Atoi(ctx.Param("id"))

	result, err := c.ticketService.UpdatePayment(userID, uint(ticketID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *TicketController) CancelPayment(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	ticketID, _ := strconv.Atoi(ctx.Param("id"))

	result, err := c.ticketService.CancelPayment(userID, uint(ticketID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *TicketController) GetAllTickets(ctx *gin.Context) {
	role := middleware.GetUserRole(ctx)
	if role != "admin" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	tickets, pagination, err := c.ticketService.GetAllTickets(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":       tickets,
		"pagination": pagination,
	})
}
