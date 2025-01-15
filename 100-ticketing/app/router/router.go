package router

import (
	"100-ticketing/app/controllers"
	"100-ticketing/app/models"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"regexp"
)

func New(mod *models.Models) *http.ServeMux {
	//routing mux
	mux := http.NewServeMux()

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
			Serve(), mod, []string{"admin"}, false),
	)

	// add MiddlewareAuth to auth checking
	mux.Handle("/api/v1/user/me",
		controllers.MiddlewareAuth(controllers.NewController(nil).
			Method("GET", userHandler.Me()). // find user by id
			Serve(), mod, []string{"admin", "user"}, false),
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
			Serve(), mod, []string{"admin"}, false),
	)

	// add MiddlewareAuth to auth checking
	mux.Handle("/api/v1/user",
		controllers.MiddlewareAuth(controllers.NewController(nil).
			Method("GET", userHandler.GetAll()). // get all users
			Serve(), mod, []string{"admin"}, false),
	)

	// create new NewRoute controllers
	newRoute := controllers.NewRoute(mod)
	// routeIdPattern - validate match path route
	routeIdPattern := regexp.MustCompile(`^/api/v1/route/([A-Za-z0-9]+)$`)
	mux.Handle("/api/v1/route/",
		controllers.NewController(routeIdPattern).
			Method("GET", controllers.MiddlewareAuth(newRoute.Get(), mod, []string{"admin", "user"}, false)). // find route
			Method("PUT", controllers.MiddlewareAuth(newRoute.Put(), mod, []string{"admin", "user"}, false)). // find route
			Serve(),
	)

	mux.Handle("/api/v1/route/import",
		controllers.NewController(nil).
			Method("POST", controllers.MiddlewareAuth(newRoute.Import(), mod, []string{"admin"}, false)). // create route
			Serve(),
	)

	mux.Handle("/api/v1/route",
		controllers.NewController(nil).
			Method("POST", controllers.MiddlewareAuth(newRoute.Create(), mod, []string{"admin"}, false)).        // create route
			Method("GET", controllers.MiddlewareAuth(newRoute.GetAll(), mod, []string{"admin", "user"}, false)). // find route
			Serve(),
	)

	// create new NewTickets controllers
	newTickets := controllers.NewTickets(mod)
	// routeIdPattern - validate match path route
	ticketIdPattern := regexp.MustCompile(`^/api/v1/ticket/([A-Za-z0-9]+)$`)
	mux.Handle("/api/v1/ticket/",
		controllers.NewController(ticketIdPattern).
			Method("GET", controllers.MiddlewareAuth(newTickets.Get(), mod, []string{"admin", "user"}, false)).       // find ticket
			Method("PUT", controllers.MiddlewareAuth(newTickets.Put(), mod, []string{"admin", "user"}, false)).       // edit ticket
			Method("DELETE", controllers.MiddlewareAuth(newTickets.Delete(), mod, []string{"admin", "user"}, false)). // edit ticket
			Serve(),
	)

	mux.Handle("/api/v1/ticket",
		controllers.NewController(nil).
			Method("POST", controllers.MiddlewareAuth(newTickets.Create(), mod, []string{"admin", "user"}, false)). // create ticket
			Method("GET", controllers.MiddlewareAuth(newTickets.GetAll(), mod, []string{"admin", "user"}, false)).  // get ticket
			Serve(),
	)

	// create new NewBooking controllers
	newBooking := controllers.NewBooking(mod)
	// bookingIdPattern - validate match path route
	bookingIdPattern := regexp.MustCompile(`^/api/v1/booking/([A-Za-z0-9]+)$`)
	mux.Handle("/api/v1/booking/",
		controllers.NewController(bookingIdPattern).
			Method("GET", controllers.MiddlewareAuth(newBooking.Get(), mod, []string{"admin", "user"}, true)). // find booking
			Serve(),
	)

	mux.Handle("/api/v1/booking",
		controllers.NewController(nil).
			Method("POST", controllers.MiddlewareAuth(newBooking.Create(), mod, []string{"admin", "user"}, true)). // create booking
			Method("GET", controllers.MiddlewareAuth(newBooking.GetAll(), mod, []string{"admin"}, true)).          // get booking
			Serve(),
	)

	mux.Handle("/api/v1/booking/complete",
		controllers.NewController(nil).
			Method("POST", controllers.MiddlewareAuth(newBooking.Complete(), mod, []string{"admin"}, true)). // create booking et booking
			Serve(),
	)

	mux.Handle("/api/v1/booking/cancel",
		controllers.NewController(nil).
			Method("POST", controllers.MiddlewareAuth(newBooking.Cancel(), mod, []string{"admin"}, true)). // create booking et booking
			Serve(),
	)

	// swagger docs server handler
	docsLocation := "./docs/"
	if len(os.Getenv("DOCS_LOCATION")) != 0 {
		docsLocation = os.Getenv("DOCS_LOCATION")
	}
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.yaml"), //The url pointing to API definition
		httpSwagger.PersistAuthorization(true),
	))
	mux.Handle("/docs/", http.StripPrefix("/docs", http.FileServer(http.Dir(docsLocation))))

	// root controllers
	mux.Handle("/", controllers.NewController(nil).
		Method("GET", controllers.NewRoot().Get()).
		Serve())

	return mux
}
