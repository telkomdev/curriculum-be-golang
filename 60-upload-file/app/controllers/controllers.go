package controllers

import (
	"60-upload-file/app/models"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"regexp"
	"strings"
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

// GetRequestBodyData - get raw data from post json body
func GetRequestBodyData(r *http.Request, requestData interface{}) error {
	all, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(all, &requestData)
}

// GetLastPathID - get id from last path
func GetLastPathID(r *http.Request) (id string) {
	res := strings.Split(r.URL.Path, "/")
	if len(res) != 0 {
		return res[len(res)-1]
	}
	return id
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

// MiddlewareSuperAdminAuth checking required for super admin auth
func MiddlewareSuperAdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// always allow options method
		if r.Method == "OPTIONS" {
			WriteResponse(w, http.StatusOK, "", false)
			return
		}

		// check if token exist in header
		authHeader := r.Header.Get("Secret-Key")
		if len(authHeader) == 0 {
			WriteResponse(w, 403, "No secret-key provided!", true)
			return
		}

		if authHeader != models.SecretKey {
			WriteResponse(w, 403, "Wrong secret-key!", true)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// MiddlewareAuth checking required auth before routing
func MiddlewareAuth(next http.Handler, model *models.Models, roles []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// always allow options method
		if r.Method == "OPTIONS" {
			WriteResponse(w, http.StatusOK, "", false)
			return
		}

		// check if token exist in header
		authHeader := r.Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer ") {
			WriteResponse(w, 401, "Unauthorized, token must be start with 'Bearer '", true)
			return
		}

		tokens := strings.Split(authHeader, "Bearer ")
		if len(tokens) != 2 {
			WriteResponse(w, 401, "Unauthorized, token segment failed", true)
			return
		}

		// validating  JWT Token
		var claim models.UserClaims
		_, err := jwt.ParseWithClaims(tokens[1], &claim, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("token not valid : signing method invalid")
			}

			if err := claim.Valid(); err != nil {
				return nil, fmt.Errorf("token not valid : %s", err)
			}

			return []byte(models.Secret), nil
		})

		if err != nil {
			WriteResponse(w, 401, fmt.Sprintf("Unauthorized, %s", err.Error()), true)
			return
		}

		// creating allow roles
		var rolesMap = make(map[string]bool)
		for _, s := range roles {
			rolesMap[s] = true
		}

		// validate roles
		for _, s := range claim.UserRoles {
			if rolesMap[s] {
				// check if user still valid
				_, err := model.FindUserById("users", claim.UserID)
				if err != nil {
					if err == mongo.ErrNoDocuments {
						WriteResponse(w, 401, "Unauthorized, user not found", true)
						return
					}

					WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
					return
				}
				// inject data to request
				r.Header.Set("Userid", claim.UserID)
				next.ServeHTTP(w, r)
				return
			}
		}

		WriteResponse(w, 401, "Unauthorized, roles not allowed", true)
		return

	})
}
