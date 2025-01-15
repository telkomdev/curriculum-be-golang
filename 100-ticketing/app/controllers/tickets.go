package controllers

import (
	"100-ticketing/app/models"
	"fmt"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type Tickets struct {
	validator *validator.Validate
	model     *models.Models
}

// NewTickets return new object of Tickets
func NewTickets(m *models.Models) *Tickets {
	return &Tickets{model: m, validator: validator.New()}
}

// Create implement net http handler
// @Tags Ticket
// @Summary Create new ticket.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body models.CreateTicketRequest true "request data"
// @Success 200 {object} models.Ticket
// @Failure 400 {object} models.ErrorCreateTicketValidation
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBUpset
// @Router /api/v1/ticket [post]
func (c *Tickets) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData models.CreateTicketRequest
		err := GetRequestBodyData(r, &requestData)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "failed to get request data", true)
			return
		}

		// validation all
		err = c.validator.Struct(requestData)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "Field from, to, userId, price, bookingId, departureTime cannot be empty!", true)
			return
		}

		_, err = time.Parse("2006-01-02T15:04:05.000Z", requestData.DepartureTime)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "departureTime must be format 2006-01-02T15:04:05.000Z", true)
			return
		}

		userId, err := primitive.ObjectIDFromHex(requestData.UserId)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("userId %s not valid : %s", userId, err.Error()), true)
			return
		}

		bookingId, err := primitive.ObjectIDFromHex(requestData.BookingId)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("bookingId %s not valid : %s", bookingId, err.Error()), true)
			return
		}

		id, err := c.model.InsertTicket(&models.Ticket{
			From:          requestData.From,
			To:            requestData.To,
			Price:         requestData.Price,
			DepartureTime: requestData.DepartureTime,
			UserId:        requestData.UserId,
			BookingId:     requestData.BookingId,
		})
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to insert or update data to mongodb, %s", err.Error()), true)
			return
		}

		data, err := c.model.FindTicketById(id)
		if err != nil {
			// write response error
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}

		WriteResponse(w, http.StatusOK, data, true)
		return
	}
}

// GetAll implement net http handler
// @Tags Ticket
// @Summary Find all ticket, or find route with filter by from, to, userId, bookingId.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string false "query to search filter by 'from' location"
// @Param to query string false "query to search filter by 'to' location"
// @Param userId query string false "query to search filter by 'userId' location"
// @Param bookingId query string false "query to search filter by 'bookingId' location"
// @Param page query integer false "Query data by page number"
// @Param size query integer false "Limit size per page"
// @Success 200 {object} models.AllRoutes
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBGet
// @Router /api/v1/ticket [get]
func (c *Tickets) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routes, err := c.model.GetAllTickets(
			r.URL.Query().Get("from"),
			r.URL.Query().Get("to"),
			r.URL.Query().Get("userId"),
			r.URL.Query().Get("bookingId"),
			GetPaginationOption(r))
		if err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}

		WriteResponse(w, http.StatusOK, routes, true)
		return
	}
}

// Get implement net http handler
// @Tags Ticket
// @Summary Get ticket by id.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Ticket id to search"
// @Success 200 {object} models.Ticket
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBGet
// @Router /api/v1/ticket/{id} [get]
func (c *Tickets) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route, err := c.model.FindTicketById(GetLastPathID(r))
		if err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}
		WriteResponse(w, http.StatusOK, &route, true)
		return
	}
}

// Put implement net http handler
// @Tags Ticket
// @Summary Update ticket by id.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Ticket id to update"
// @Param data body models.CreateTicketRequest true "request data"
// @Success 200 {object} models.TicketEditResponseSuccess
// @Failure 400 {object} models.ErrorCreateTicketValidation
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBUpset
// @Router /api/v1/ticket/{id} [put]
func (c *Tickets) Put() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData models.CreateTicketRequest
		err := GetRequestBodyData(r, &requestData)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "failed to get request data", true)
			return
		}

		ticket, err := c.model.FindTicketById(GetLastPathID(r))
		if err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}

		// validation all
		err = c.validator.Struct(requestData)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "Field from, to, userId, price, bookingId, departureTime cannot be empty!", true)
			return
		}

		_, err = time.Parse("2006-01-02T15:04:05.000Z", requestData.DepartureTime)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "departureTime must be format 2006-01-02T15:04:05.000Z", true)
			return
		}

		userId, err := primitive.ObjectIDFromHex(requestData.UserId)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("userId %s not valid : %s", userId, err.Error()), true)
			return
		}

		bookingId, err := primitive.ObjectIDFromHex(requestData.BookingId)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("bookingId %s not valid : %s", bookingId, err.Error()), true)
			return
		}

		err = c.model.UpsetTicketById(&models.Ticket{
			Id:            ticket.Id,
			From:          requestData.From,
			To:            requestData.To,
			Price:         requestData.Price,
			DepartureTime: requestData.DepartureTime,
			UserId:        requestData.UserId,
			BookingId:     requestData.BookingId,
		})

		if err != nil {
			WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to insert or update data to mongodb, %s", err.Error()), true)
			return
		}

		WriteResponse(w, http.StatusOK, "Ticket was updated successfully!", true)
		return
	}
}

// Delete implement net http handler
// @Tags Ticket
// @Summary Delete ticket by id
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Ticket id to delete"
// @Success 200 {object} models.TicketDeleteResponseSuccess
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBUpset
// @Router /api/v1/ticket/{id} [delete]
func (c *Tickets) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ticket, err := c.model.FindTicketById(GetLastPathID(r))
		if err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}

		if err := c.model.DeleteTicketByID(ticket.Id); err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to delete data from mongodb, %s", err.Error()), true)
			return
		}

		WriteResponse(w, http.StatusOK, fmt.Sprintf("Ticket with id %s was deleted successfully!", ticket.Id), true)

		return
	}
}
