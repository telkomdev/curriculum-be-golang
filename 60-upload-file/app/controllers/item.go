package controllers

import (
	"60-upload-file/app/models"
	"fmt"
	"net/http"
)

type Items struct {
	model *models.Models
}

func NewItemsList(model *models.Models) *Items {
	return &Items{model: model}
}

// Get Items
// @Tags Items
// @Summary This endpoint will show all item data in JSON format, need too login and authorized as user or admin.
// @Security BearerAuth
// @Success 200 {object} models.AllItem
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBGet
// @Router /api/v1/item [get]
func (c *Items) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := c.model.GetAllItem("items")
		if err != nil {
			// write response error
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}

		// write response ok
		WriteResponse(w, http.StatusOK, res, true)
	}
}

// Post Items
// @Tags Items
// @Summary Request to create new item, need too login and authorized as user or admin.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body models.CreateItemRequest true "request data"
// @Success 200 {object} models.Item
// @Failure 400 {object} models.ErrorCreateItem
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBUpset
// @Router /api/v1/item [post]
func (c *Items) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var requestData models.CreateItemRequest
		err := GetRequestBodyData(r, &requestData)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "failed to get request data", true)
			return
		}

		if len(requestData.Name) == 0 || requestData.Qty < 1 {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "name and qty cannot be empty", true)
			return
		}

		err = c.model.InsertItem("items", &models.Item{
			Name: requestData.Name,
			Qty:  requestData.Qty,
		})

		if err != nil {
			// write response error
			WriteResponse(w, 500, fmt.Sprintf("failed to insert or update data to mongodb, %s", err.Error()), true)
			return
		}

		res, err := c.model.FindItemByName("items", requestData.Name)
		if err != nil {
			// write response error
			WriteResponse(w, 500, fmt.Sprintf("failed to insert or update data to mongodb, %s", err.Error()), true)
			return
		}

		// write response ok
		WriteResponse(w, http.StatusOK, res, true)
	}
}
