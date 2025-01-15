package controllers

import (
	"100-ticketing/app/models"
	"encoding/csv"
	"fmt"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Route struct {
	validator *validator.Validate
	model     *models.Models
}

// NewRoute return new object of Route
func NewRoute(m *models.Models) *Route {
	return &Route{model: m, validator: validator.New()}
}

// Create implement net http handler
// @Tags Route
// @Summary Create new route, this feature need Role Admin.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body models.CreateRouteRequest true "request data"
// @Success 200 {object} models.Route
// @Failure 400 {object} models.ErrorCreateRouteValidation
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBUpset
// @Router /api/v1/route [post]
func (c *Route) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData models.CreateRouteRequest
		err := GetRequestBodyData(r, &requestData)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "failed to get request data", true)
			return
		}

		err = c.validator.Struct(requestData)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "Field from, to, price, departureTime cannot be empty!", true)
			return
		}

		// validating route departureTime
		_, err = time.Parse("15:04:05", requestData.DepartureTime)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "departureTime must be format 15:04:05", true)
			return
		}

		err = c.model.UpsetRouteByFromAndTo(&models.Route{
			From:          requestData.From,
			To:            requestData.To,
			Price:         requestData.Price,
			DepartureTime: requestData.DepartureTime,
		})
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to insert or update data to mongodb, %s", err.Error()), true)
			return
		}

		data, err := c.model.FindRouteByFromAndTo(requestData.From, requestData.To)
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
// @Tags Route
// @Summary Find all route, or find route with filter by from location or to location.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string false "query to search filter by 'from' location"
// @Param to query string false "query to search filter by 'to' location"
// @Param page query integer false "Query data by page number"
// @Param size query integer false "Limit size per page"
// @Success 200 {object} models.AllRoutes
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBGet
// @Router /api/v1/route [get]
func (c *Route) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routes, err := c.model.GetAllRoutes(r.URL.Query().Get("from"),
			r.URL.Query().Get("to"), GetPaginationOption(r))
		if err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}

		WriteResponse(w, http.StatusOK, routes, true)
		return
	}
}

// Get implement net http handler
// @Tags Route
// @Summary Get route by id.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Route id to search"
// @Success 200 {object} models.Route
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBGet
// @Router /api/v1/route/{id} [get]
func (c *Route) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := GetLastPathID(r)
		route, err := c.model.FindRouteById(id)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				WriteResponse(w, 400, fmt.Sprintf("route with id %s not found", id), true)
				return
			}
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}
		WriteResponse(w, http.StatusOK, &route, true)
		return
	}
}

// Put implement net http handler
// @Tags Route
// @Summary Update route by id.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Route id to update"
// @Param data body models.UpdateRouteRequest true "request data"
// @Success 200 {object} models.RouteEditResponseSuccess
// @Failure 400 {object} models.ErrorUpdateRouteValidation
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBGet
// @Router /api/v1/route/{id} [put]
func (c *Route) Put() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData models.UpdateRouteRequest
		err := GetRequestBodyData(r, &requestData)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "failed to get request data", true)
			return
		}

		route, err := c.model.FindRouteById(GetLastPathID(r))
		if err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}

		// validate name
		err = c.validator.Struct(requestData)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "Field price, departureTime cannot be empty!", true)
			return
		}

		// validating route departureTime
		_, err = time.Parse("15:04:05", requestData.DepartureTime)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "departureTime must be format 15:04:05", true)
			return
		}

		err = c.model.UpsetRouteByFromAndTo(&models.Route{
			From:          route.From,
			To:            route.To,
			Price:         requestData.Price,
			DepartureTime: requestData.DepartureTime,
			CreatedAt:     route.CreatedAt,
			UpdatedAt:     route.UpdatedAt,
		})

		if err != nil {
			WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to insert or update data to mongodb, %s", err.Error()), true)
			return
		}

		WriteResponse(w, http.StatusOK, "Route was updated successfully!", true)
		return
	}
}

// Import implement net http handler
// @Tags Route
// @Summary Import Ticketing Route CSV file to App server, this feature only accessible by role admin.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param routeFile formData file true "CSV route file to import"
// @Success 200 {object} models.RouteImportResponseSuccess
// @Failure 400 {object} models.ErrorCreateRouteValidation
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBUpset
// @Router /api/v1/route/import [post]
func (c *Route) Import() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		file, _, err := r.FormFile("routeFile")
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to get routeFile , %s", err.Error()), true)
			return
		}

		f, err := os.CreateTemp("", "routeFile.*.csv")
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to get routeFile , %s", err.Error()), true)
			return
		}

		defer func(name string) {
			_ = os.Remove(name)
		}(f.Name()) // clean up

		_, _ = io.Copy(f, file)
		if err != nil {
			return
		}

		// open file
		ff, err := os.Open(f.Name())
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to get routeFile , %s", err.Error()), true)
			return
		}

		// remember to close the file at the end of the program
		defer func(ff *os.File) {
			_ = ff.Close()
		}(ff)

		csvReader := csv.NewReader(ff)
		data, err := csvReader.ReadAll()
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to get routeFile , %s", err.Error()), true)
			return
		}

		for i, line := range data {
			if i > 0 { // omit header line
				rData := models.Route{}
				for i, field := range line {
					switch i {
					case 0:
						rData.From = field
					case 1:
						rData.To = field
					case 2:
						if vl, err := strconv.ParseFloat(field, 64); err == nil {
							rData.Price = vl
						}
					case 3:
						rData.DepartureTime = field
					}
				}

				err := c.validator.Struct(rData)
				if err != nil {
					// write response error
					WriteResponse(w, http.StatusBadRequest, "Field from, to, price, departureTime cannot be empty!", true)
					return
				}

				// validating route departureTime
				_, err = time.Parse("15:04:05", rData.DepartureTime)
				if err != nil {
					// write response error
					WriteResponse(w, http.StatusBadRequest, "departureTime must be format 15:04:05", true)
					return
				}

				err = c.model.UpsetRouteByFromAndTo(&rData)
				if err != nil {
					WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to insert or update data to mongodb, %s", err.Error()), true)
					return
				}
			}
		}

		WriteResponse(w, http.StatusOK, "Route data has successfully imported!", true)
		return
	}
}
