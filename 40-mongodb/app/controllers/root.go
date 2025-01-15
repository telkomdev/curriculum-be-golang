package controllers

import (
	"40-mongodb/app/models"
	"net/http"
)

type Root struct{}

func NewRoot() *Root {
	return &Root{}
}

// Get implement net http handler
// @Tags Root
// @Summary Response this request with Hello Route.
// @Success 200 {object} models.ResponseRoot
// @Router / [get]
func (c *Root) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteResponse(w, http.StatusOK, models.ResponseRoot{Message: "Hello. Try GET/POST to /api/v1/item"}, true)
	}
}
