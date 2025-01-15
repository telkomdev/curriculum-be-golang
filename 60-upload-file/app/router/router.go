package router

import (
	"60-upload-file/app/controllers"
	"60-upload-file/app/models"
	"net/http"
	"regexp"
)

func New(mod *models.Models) *http.ServeMux {
	//routing mux
	mux := http.NewServeMux()

	// create new NewUser controllers
	itemHandler := controllers.NewItemsList(mod)

	// create new NewItemsList controllers
	// create new items controllers, add MiddlewareAuth to auth checking
	// add MiddlewareAuth to auth checking
	mux.Handle("/api/v1/item",
		controllers.MiddlewareAuth(controllers.NewController(nil).
			Method("GET", itemHandler.Get()).   // list all items
			Method("POST", itemHandler.Post()). // add item
			Serve(), mod, []string{"admin", "user"}),
	)

	// create new NewUser controllers
	userHandler := controllers.NewUser(mod)

	// create new UserCreateAdmin controllers
	// add MiddlewareSuperAdminAuth to auth checking
	mux.Handle("/api/v1/user/create/admin",
		controllers.MiddlewareSuperAdminAuth(controllers.NewController(nil).
			Method("POST", userHandler.CreateAdmin()). // find user by id
			Serve()),
	)

	// create new UserCreate controllers
	// add MiddlewareAuth to auth checking
	mux.Handle("/api/v1/user/create",
		controllers.MiddlewareAuth(controllers.NewController(nil).
			Method("POST", userHandler.Create()). // find user by id
			Serve(), mod, []string{"admin"}),
	)

	// add MiddlewareAuth to auth checking
	mux.Handle("/api/v1/user/me",
		controllers.MiddlewareAuth(controllers.NewController(nil).
			Method("GET", userHandler.Me()). // find user by id
			Serve(), mod, []string{"admin", "user"}),
	)

	// add MiddlewareAuth to auth checking
	mux.Handle("/api/v1/user/auth",
		controllers.NewController(nil).
			Method("POST", userHandler.Auth()).
			Serve(),
	)

	// add MiddlewareAuth to auth checking
	// userIdPattern - validate match path user
	userIdPattern := regexp.MustCompile(`^/api/v1/user/([A-Za-z0-9]+)$`)
	mux.Handle("/api/v1/user/",
		controllers.MiddlewareAuth(controllers.NewController(userIdPattern).
			Method("GET", userHandler.Get()).       // find user by id
			Method("PUT", userHandler.Put()).       // edit user by id
			Method("DELETE", userHandler.Delete()). // delete user by id
			Serve(), mod, []string{"admin"}),
	)

	// add MiddlewareAuth to auth checking
	mux.Handle("/api/v1/user",
		controllers.MiddlewareAuth(controllers.NewController(nil).
			Method("GET", userHandler.GetAll()). // get all users
			Serve(), mod, []string{"admin"}),
	)

	// create new NewRoute controllers
	newRoute := controllers.NewRoute(mod)
	// routeIdPattern - validate match path route
	routeIdPattern := regexp.MustCompile(`^/api/v1/route/([A-Za-z0-9]+)$`)
	mux.Handle("/api/v1/route/",
		controllers.NewController(routeIdPattern).
			Method("GET", controllers.MiddlewareAuth(newRoute.Get(), mod, []string{"admin", "user"})). // find route
			Method("PUT", controllers.MiddlewareAuth(newRoute.Put(), mod, []string{"admin", "user"})). // find route
			Serve(),
	)

	mux.Handle("/api/v1/route/import",
		controllers.NewController(nil).
			Method("POST", controllers.MiddlewareAuth(newRoute.Import(), mod, []string{"admin"})). // create route
			Serve(),
	)

	mux.Handle("/api/v1/route",
		controllers.NewController(nil).
			Method("POST", controllers.MiddlewareAuth(newRoute.Create(), mod, []string{"admin"})).        // create route
			Method("GET", controllers.MiddlewareAuth(newRoute.GetAll(), mod, []string{"admin", "user"})). // find route
			Serve(),
	)

	// root controllers
	mux.Handle("/", controllers.NewController(nil).
		Method("GET", controllers.NewRoot().Get()).
		Serve())

	return mux
}
