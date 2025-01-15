package controllers

import (
	"30-routing/app/models"
	"encoding/json"
	"net/http"
	"regexp"
)

func WriteResponse(w http.ResponseWriter, code int, message interface{}, isJson bool) {
	var isError bool
	if isJson {
		w.Header().Set("Content-Type", "application/json")
	}

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Secret-Key, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.WriteHeader(code)

	if code != 200 {
		isError = true
	}

	if isJson {
		if vals, ok := message.(string); ok {
			res, _ := json.Marshal(models.ResponseItem{Error: isError, Message: vals})
			_, _ = w.Write(res)
		} else {
			res, _ := json.Marshal(message)
			_, _ = w.Write(res)
		}
	} else {
		if vals, ok := message.(string); ok {
			_, _ = w.Write([]byte(vals))
		}
	}

}

type Controller struct {
	// using for saving object handler
	match   *regexp.Regexp
	methods map[string]http.Handler
}

// NewController - creating new controllers
func NewController(match *regexp.Regexp) *Controller {
	return &Controller{
		match:   match,
		methods: map[string]http.Handler{},
	}
}

func (c *Controller) Method(method string, next http.Handler) *Controller {
	c.methods[method] = next
	return c
}

func (c *Controller) Serve() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// always allow options method
		if r.Method == "OPTIONS" {
			WriteResponse(w, http.StatusOK, "", false)
			return
		}

		if c.match != nil {
			matches := c.match.MatchString(r.URL.Path)
			if !matches {
				WriteResponse(w, http.StatusNotFound, "Not Found", false)
				return
			}
		}

		if c.methods[r.Method] != nil {
			c.methods[r.Method].ServeHTTP(w, r)
			return
		}

		WriteResponse(w, http.StatusMethodNotAllowed, "Method not allowed", false)
		return

	})
}
