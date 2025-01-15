package controllers

import (
	"60-upload-file/app/models"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type User struct {
	validator *validator.Validate
	model     *models.Models
}

// NewUser return new object of User
func NewUser(m *models.Models) *User {
	return &User{model: m, validator: validator.New()}
}

// GetAll implement net http handler
// @Tags User
// @Summary Get all user and search all user by name, this endpoint only available for role admin.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name query string false "Name to search"
// @Success 200 {object} models.AllUsers
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBGet
// @Router /api/v1/user [get]
func (c *User) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := c.model.GetAllUsers("users", r.URL.Query().Get("name"))
		if err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}

		WriteResponse(w, http.StatusOK, user, true)
		return
	}
}

// Get implement net http handler
// @Tags User
// @Summary Find user by id, this feature need role admin.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User id to search"
// @Success 200 {object} models.User
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBGet
// @Router /api/v1/user/{id} [get]
func (c *User) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := c.model.FindUserById("users", GetLastPathID(r))
		if err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}

		WriteResponse(w, http.StatusOK, &models.User{
			Id:        user.Id,
			Email:     user.Email,
			Name:      user.Name,
			Roles:     user.Roles,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}, true)

		return
	}
}

// Put implement net http handler
// @Tags User
// @Summary Update user by id, this feature need Role Admin.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User id to edit"
// @Param data body models.UserEditRequest true "request data"
// @Success 200 {object} models.UserEditResponseSuccess
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBUpset
// @Router /api/v1/user/{id} [put]
func (c *User) Put() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData models.UserEditRequest
		err := GetRequestBodyData(r, &requestData)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "failed to get request data", true)
			return
		}

		user, err := c.model.FindUserById("users", GetLastPathID(r))
		if err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}

		// validate name
		err = c.validator.Var(requestData.Name, "required")
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "User validation failed: name: required", true)
			return
		}

		user.Name = requestData.Name
		err = c.model.UpsetUserByEmail("users", &user)
		if err != nil {
			WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to insert or update data to mongodb, %s", err.Error()), true)
			return
		}

		WriteResponse(w, http.StatusOK, "User was updated successfully!", true)

		return
	}
}

// Delete implement net http handler
// @Tags User
// @Summary Delete user by id, this feature need Role Admin.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User id to edit"
// @Success 200 {object} models.UserDeleteResponseSuccess
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBUpset
// @Router /api/v1/user/{id} [delete]
func (c *User) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user, err := c.model.FindUserById("users", GetLastPathID(r))
		if err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}

		if err := c.model.DeleteUserByID("users", user.Id); err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to delete data from mongodb, %s", err.Error()), true)
			return
		}

		WriteResponse(w, http.StatusOK, fmt.Sprintf("User with id %s was deleted successfully!", user.Id), true)

		return
	}
}

// Me implement net http handler
// @Tags User
// @Summary Get current login user profile.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBGet
// @Router /api/v1/user/me [get]
func (c *User) Me() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get userId from header
		userId := r.Header.Get("Userid")
		user, err := c.model.FindUserById("users", userId)
		if err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}

		WriteResponse(w, http.StatusOK, &models.User{
			Id:        user.Id,
			Email:     user.Email,
			Name:      user.Name,
			Roles:     user.Roles,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}, true)

		return
	}
}

// Auth implement net http handler
// @Tags User
// @Summary Login user.
// @Accept json
// @Produce json
// @Param data body models.UserAuthRequest true "request data"
// @Success 200 {object} models.UserAuthResponse
// @Failure 400 {object} models.ErrorUserNotFound
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBGet
// @Router /api/v1/user/auth [post]
func (c *User) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData models.UserAuthRequest
		err := GetRequestBodyData(r, &requestData)
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "failed to get request data", true)
			return
		}

		// validate email
		err = c.validator.Var(requestData.Email, "email")
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "User validation failed: email: Please enter a valid email", true)
			return
		}

		// validate password
		err = c.validator.Var(requestData.Password, "required")
		if err != nil {
			// write response error
			WriteResponse(w, http.StatusBadRequest, "User validation failed: password: required", true)
			return
		}

		// find user by email
		user, err := c.model.FindUserByEmail("users", requestData.Email)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// write response error
				WriteResponse(w, http.StatusBadRequest, "User Not Found!", true)
				return
			}

			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}

		// validating password
		if user.Password != requestData.Password {
			// write response error
			WriteResponse(w, http.StatusUnauthorized, "Wrong password!", true)
			return
		}

		var roles []string
		for _, s := range user.Roles {
			roles = append(roles, s.Name)
		}

		// generated token
		expired, _ := time.ParseDuration("7200h")
		claims := &models.UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
				NotBefore: jwt.NewNumericDate(time.Now().UTC()),
				ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expired)),
			},
			UserID:    user.Id,
			UserRoles: roles,
		}

		token := jwt.NewWithClaims(models.JWtAlg, claims)
		signedToken, err := token.SignedString([]byte(models.Secret))
		if err != nil {
			WriteResponse(w, 500, fmt.Sprintf("failed to generated token, %s", err.Error()), true)
			return
		}

		WriteResponse(w, http.StatusOK, &models.UserAuthResponse{
			Id:          user.Id,
			Email:       user.Email,
			Name:        user.Name,
			Roles:       user.Roles,
			AccessToken: signedToken,
		}, true)

		return
	}
}

// Create implement net http handler
// @Tags User
// @Summary Create new user, need role admin bearer token.
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body models.UserCreateRequest true "request data"
// @Success 200 {object} models.UserCreateResponseSuccess
// @Failure 400 {object} models.ErrorEmailNotValid
// @Failure 401 {object} models.Unauthorized
// @Failure 500 {object} models.ErrorMongoDBUpset
// @Router /api/v1/user/create [post]
func (c *User) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.registerUser(c.model, "user", w, r)
	}
}

// CreateAdmin implement net http handler
// @Tags User
// @Summary Create new user admin, need super admin secret key.
// @Accept json
// @Produce json
// @Param secret-key header string true "Your super admin secret-key"
// @Param data body models.UserCreateRequest true "request data"
// @Success 200 {object} models.UserCreateResponseSuccess
// @Failure 400 {object} models.ErrorEmailNotValid
// @Failure 403 {object} models.ErrorWrongSecretKey
// @Failure 500 {object} models.ErrorMongoDBUpset
// @Router /api/v1/user/create/admin [post]
func (c *User) CreateAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.registerUser(c.model, "admin", w, r)
	}
}

// registerUser - function to register user, used by UserCreateAdmin or UserCreate
func (c *User) registerUser(model *models.Models, useRole string, w http.ResponseWriter, r *http.Request) {
	var requestData models.UserCreateRequest
	err := GetRequestBodyData(r, &requestData)
	if err != nil {
		// write response error
		WriteResponse(w, http.StatusBadRequest, "failed to get request data", true)
		return
	}

	// validate email
	err = c.validator.Var(requestData.Email, "email")
	if err != nil {
		// write response error
		WriteResponse(w, http.StatusBadRequest, "User validation failed: email: Please enter a valid email", true)
		return
	}

	// validate name
	err = c.validator.Var(requestData.Name, "required")
	if err != nil {
		// write response error
		WriteResponse(w, http.StatusBadRequest, "User validation failed: name: required", true)
		return
	}

	// validate password
	err = c.validator.Var(requestData.Password, "required")
	if err != nil {
		// write response error
		WriteResponse(w, http.StatusBadRequest, "User validation failed: password: required", true)
		return
	}

	// find user by email
	_, err = model.FindUserByEmail("users", requestData.Email)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
			return
		}
	}

	// return error, user already registered
	if err == nil {
		WriteResponse(w, http.StatusBadRequest, "User validation failed: email: already registered", true)
		return
	}

	// find role
	role, err := model.FindRoleByName("roles", useRole)
	if err != nil {
		WriteResponse(w, 500, fmt.Sprintf("failed to get data from mongodb, %s", err.Error()), true)
		return
	}

	err = model.UpsetUserByEmail("users", &models.UserList{
		Email:    requestData.Email,
		Name:     requestData.Name,
		Password: requestData.Password,
		Roles:    []models.Role{role}})
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to insert or update data to mongodb, %s", err.Error()), true)
		return
	}

	WriteResponse(w, http.StatusOK, "User was registered successfully!", true)
}
