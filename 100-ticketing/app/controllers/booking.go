package controllers

import (
	"100-ticketing/app/adapter/services"
	"100-ticketing/app/models"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type Booking struct {
	validator *validator.Validate
	model     *models.Models
}

// NewBooking return new object of Booking
func NewBooking(m *models.Models) *Booking {
	return &Booking{model: m, validator: validator.New()}
}

// Get implement net http handler
// @Tags Booking
// @Summary Get booking by id.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Booking id to search"
// @Success 200 {object} models.Booking
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBGet
// @Router /api/v1/booking/{id} [get]
func (c *Booking) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route, err := c.model.FindBookingById(GetLastPathID(r))
		if err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}
		WriteResponse(w, http.StatusOK, &route, true)
		return
	}
}

// GetAll implement net http handler
// @Tags Booking
// @Summary Get all bookings, this feature need Role Admin.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query integer false "Query data by page number"
// @Param size query integer false "Limit size per page"
// @Success 200 {object} models.AllBookings
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBGet
// @Router /api/v1/booking [get]
func (c *Booking) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routes, err := c.model.GetAllBookings(GetPaginationOption(r))
		if err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}

		WriteResponse(w, http.StatusOK, routes, true)
		return
	}
}

// Create implement net http handler
// @Tags Booking
// @Summary Create new booking
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body models.CreateBookingRequest true "request data"
// @Success 200 {object} models.Booking
// @Failure 400 {object} models.ErrorRouteNotFound
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBUpset
// @Router /api/v1/booking [post]
func (c *Booking) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData models.CreateBookingRequest
		err := GetRequestBodyData(r, &requestData)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "failed to get request data", true)
			return
		}

		err = c.validator.Struct(requestData)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "Field routeId, quantity, scheduleDate cannot be empty!", true)
			return
		}

		// validating schedule date
		_, err = time.Parse("2006-01-02", requestData.ScheduleDate)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "scheduleDate must be format 2006-01-02", true)
			return
		}

		// validation route, call external services
		route, err := services.Load().GetRouteByID(r.Header.Get("Authorization"), requestData.RouteId)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, err.Error(), true)
			return
		}

		// validating route departureTime
		_, err = time.Parse("15:04:05", route.DepartureTime)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "route departureTime must be format 15:04:05", true)
			return
		}

		departureTime := fmt.Sprintf("%sT%s.000Z", requestData.ScheduleDate, route.DepartureTime)
		bookingId, err := c.model.InsertBooking(&models.Booking{
			Quantity:      requestData.Quantity,
			TotalPrice:    route.Price * float64(requestData.Quantity),
			DepartureTime: departureTime,
			User: models.BookingUser{
				Id:    r.Header.Get("Userid"),
				Name:  r.Header.Get("UserName"),
				Email: r.Header.Get("UserEmail"),
			},
		})

		if err != nil {
			WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to insert or update data to mongodb, %s", err.Error()), true)
			return
		}

		bookingData, err := c.model.FindBookingById(bookingId)
		if err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}

		for i := 0; i < requestData.Quantity; i++ {
			ticket, err := services.Load().CreateNewTicket(r.Header.Get("Authorization"), models.CreateTicketRequest{
				From:          route.From,
				To:            route.To,
				Price:         route.Price,
				DepartureTime: departureTime,
				UserId:        r.Header.Get("Userid"),
				BookingId:     bookingId,
			})
			if err != nil {
				WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to create ticket, %s", err.Error()), true)
				return
			}
			bookingData.Tickets = append(bookingData.Tickets, ticket)
		}

		if err := c.model.UpdateBookingById(&bookingData); err != nil {
			WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to insert or update data to mongodb, %s", err.Error()), true)
			return
		}

		WriteResponse(w, http.StatusOK, bookingData, true)
		return
	}
}

// Complete implement net http handler
// @Tags Booking
// @Summary Update payment status to 1 (completed), this feature need Role Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body models.BookingEditRequest true "request data"
// @Success 200 {object} models.Booking
// @Failure 404 {object} models.ErrorBookingPaymentStatusValidation
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBUpset
// @Router /api/v1/booking/complete [post]
func (c *Booking) Complete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Finish(1, w, r)
	}
}

// Cancel implement net http handler
// @Tags Booking
// @Summary Update payment status to 2 (cancelled), this feature need Role Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body models.BookingEditRequest true "request data"
// @Success 200 {object} models.Booking
// @Failure 404 {object} models.ErrorBookingPaymentStatusValidation
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBUpset
// @Router /api/v1/booking/cancel [post]
func (c *Booking) Cancel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Finish(2, w, r)
	}
}

func (c *Booking) Finish(status int, w http.ResponseWriter, r *http.Request) {
	var requestData models.BookingEditRequest
	err := GetRequestBodyData(r, &requestData)
	if err != nil {
		// write response error
		WriteResponse(w, http.StatusBadRequest, "failed to get request data", true)
		return
	}

	err = c.validator.Struct(requestData)
	if err != nil {
		// write response error
		WriteResponse(w, http.StatusBadRequest, "Field _id cannot be empty!", true)
		return
	}

	bookingData, err := c.model.FindBookingById(requestData.Id)
	if err != nil {
		WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
		return
	}

	if bookingData.PaymentStatus != 0 {
		WriteResponse(w, 404, "Error, complete payment only can be done when payment status is 0 (Created)", true)
		return
	}

	bookingData.PaymentStatus = status

	if err := c.model.UpdateBookingById(&bookingData); err != nil {
		WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to insert or update data to mongodb, %s", err.Error()), true)
		return
	}

	WriteResponse(w, http.StatusOK, bookingData, true)
	return
}
