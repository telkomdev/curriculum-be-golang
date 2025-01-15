package models

import "github.com/golang-jwt/jwt/v4"

var JWtAlg = jwt.SigningMethodHS256

type Index struct {
	Collection string
	Field      string
	Unique     bool
}

const (
	RoleCollectionName    = "roles"
	ItemCollectionName    = "items"
	RouteCollectionName   = "routes"
	UserCollectionName    = "users"
	TicketCollectionName  = "tickets"
	BookingCollectionName = "bookings"
)

// IndexModels - use for mongodb index field, very usefully for sort and filter data
var IndexModels []Index = []Index{
	{Collection: RoleCollectionName, Field: "name", Unique: true},
	{Collection: ItemCollectionName, Field: "name"},
	{Collection: ItemCollectionName, Field: "createdAt"},
	{Collection: RouteCollectionName, Field: "from"},
	{Collection: RouteCollectionName, Field: "to"},
	{Collection: RouteCollectionName, Field: "createdAt"},
	{Collection: UserCollectionName, Field: "name"},
	{Collection: UserCollectionName, Field: "email", Unique: true},
	{Collection: UserCollectionName, Field: "createdAt"},
	{Collection: TicketCollectionName, Field: "createdAt"},
	{Collection: TicketCollectionName, Field: "userId"},
	{Collection: TicketCollectionName, Field: "bookingId"},
	{Collection: TicketCollectionName, Field: "from"},
	{Collection: TicketCollectionName, Field: "to"},
}

// PaginationOption - Limit the Query Size
type PaginationOption struct {
	Page int
	Size int
}

// swagger:model ResponseRoot
type ResponseRoot struct {
	Error   bool   `json:"error" example:"false"`
	Message string `json:"message" example:"This is ticketing app."`
}

// swagger:model ResponseItem
type ResponseItem struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"something when wrong"`
}

// swagger:model CreateItemRequest
type CreateItemRequest struct {
	Name string `json:"name,omitempty" example:"Test Data"`
	Qty  int    `json:"qty,omitempty" example:"1"`
}

// swagger:model Item
type Item struct {
	Id        string `json:"_id,omitempty" bson:"_id,omitempty" example:"6336185fc31ad7ad4022ab87"`
	Name      string `json:"name" example:"Test Data"`
	Qty       int    `json:"qty" example:"1"`
	CreatedAt string `json:"createdAt" example:"2022-09-30T05:12:47.469Z"`
	UpdatedAt string `json:"updatedAt" example:"2022-09-30T05:12:47.469Z"`
}

// swagger:model AllItem
type AllItem struct {
	Count       int    `json:"count" example:"1"`
	TotalPages  int64  `json:"totalPages" example:"10"`
	CurrentPage int64  `json:"currentPage" example:"1"`
	TotalCount  int64  `json:"totalCount" example:"1000"`
	Data        []Item `json:"data"`
}

// swagger:model Unauthorized
type Unauthorized struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"Unauthorized"`
}

// swagger:model ErrorCreateItem
type ErrorCreateItem struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"name and qty cannot be empty"`
}

// swagger:model ErrorMongoDBUpset
type ErrorMongoDBUpset struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"failed to insert or update data to mongodb, {error}"`
}

// swagger:model RequestGetItemInternalServerError
type ErrorMongoDBGet struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"failed to get data from mongodb, {error}"`
}

// Role - type for creating rules object
type Role struct {
	Id   string `json:"_id,omitempty" bson:"_id,omitempty" example:"6336185fc31ad7ad4022ab87"`
	Name string `json:"name" example:"admin"`
}

// swagger:model UserCreateRequest
type UserCreateRequest struct {
	Email    string `json:"email" example:"root@email.com"`
	Name     string `json:"name" example:"Administrator"`
	Password string `json:"password" example:"123456"`
}

// UserList - Use to insert or list from mongodb
type UserList struct {
	Id        string `json:"_id,omitempty" bson:"_id,omitempty" example:"6336185fc31ad7ad4022ab87"`
	Email     string `json:"email" example:"john.doe@email.com"`
	Name      string `json:"name" example:"John Doe"`
	Password  string `json:"password" example:"123456"`
	Roles     []Role `json:"roles,omitempty"`
	CreatedAt string `json:"createdAt" example:"2022-09-30T05:12:47.469Z"`
	UpdatedAt string `json:"updatedAt" example:"2022-09-30T05:12:47.469Z"`
}

// swagger:model UserCreateResponseSuccess
type UserCreateResponseSuccess struct {
	Message string `json:"message" example:"User was registered successfully!"`
}

// swagger:model UserEditResponseSuccess
type UserEditResponseSuccess struct {
	Message string `json:"message" example:"User was updated successfully!"`
}

// swagger:model UserDeleteResponseSuccess
type UserDeleteResponseSuccess struct {
	Message string `json:"message" example:"User with id 631ea95d0770f442fd692fa8 was deleted successfully!"`
}

// swagger:model UserCreateResponseErrorEmailNotValid
type ErrorEmailNotValid struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"User validation failed: email: Please enter a valid email"`
}

// swagger:model UserCreateAdminResponseErrorWrongSecretKey
type ErrorWrongSecretKey struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"No secret-key provided or wrong secret-key!"`
}

// swagger:model User
type User struct {
	Id        string `json:"_id,omitempty" bson:"_id,omitempty" example:"6336185fc31ad7ad4022ab87"`
	Email     string `json:"email" example:"john.doe@email.com"`
	Name      string `json:"name" example:"John Doe"`
	Roles     []Role `json:"roles,omitempty"`
	CreatedAt string `json:"createdAt" example:"2022-09-30T05:12:47.469Z"`
	UpdatedAt string `json:"updatedAt" example:"2022-09-30T05:12:47.469Z"`
}

// swagger:model UserAuthRequest
type UserAuthRequest struct {
	Email    string `json:"email" example:"root@email.com"`
	Password string `json:"password" example:"123456"`
}

// swagger:model UserEditRequest
type UserEditRequest struct {
	Name string `json:"name" example:"John Doe"`
}

// swagger:model ErrorUserNotFound
type ErrorUserNotFound struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"User Not Found!"`
}

// swagger:model UserAuthResponse
type UserAuthResponse struct {
	Id          string `json:"_id,omitempty" bson:"_id,omitempty" example:"6336185fc31ad7ad4022ab87"`
	Email       string `json:"email" example:"john.doe@email.com"`
	Name        string `json:"name" example:"John Doe"`
	Roles       []Role `json:"roles,omitempty"`
	AccessToken string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTA2NTQyNjgsIm5iZiI6MTY2NDczNDI2OCwiaWF0IjoxNjY0NzM0MjY4LCJ1c2VySWQiOiI2MzM5Yzg4ZmMzMDMwZjNmM2RmMjUwNGUiLCJ1c2VyUm9sZXMiOlsiYWRtaW4iXX0.DKxQzeaLna3H8MS55nQ2p96KPA_LS3bHhoIrqcNaODQ"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID    string   `json:"userId"`
	UserRoles []string `json:"userRoles"`
}

// swagger:model AllUsers
type AllUsers struct {
	Count       int    `json:"count" example:"1"`
	TotalPages  int64  `json:"totalPages" example:"10"`
	CurrentPage int64  `json:"currentPage" example:"1"`
	TotalCount  int64  `json:"totalCount" example:"1000"`
	Data        []User `json:"data"`
}

// swagger:model Route
type Route struct {
	Id            string  `json:"_id,omitempty" bson:"_id,omitempty" example:"633c692ac31ad7ad4062d0fd"`
	From          string  `json:"from" example:"Malang"`
	To            string  `json:"to" example:"Jakarta"`
	Price         float64 `json:"price" example:"200000"`
	DepartureTime string  `json:"departureTime" example:"09:00:00"`
	CreatedAt     string  `json:"createdAt" example:"2022-09-30T05:12:47.469Z"`
	UpdatedAt     string  `json:"updatedAt" example:"2022-09-30T05:12:47.469Z"`
}

// swagger:model CreateRouteRequest
type CreateRouteRequest struct {
	From          string  `json:"from" example:"Malang" validate:"required"`
	To            string  `json:"to" example:"Jakarta" validate:"required"`
	Price         float64 `json:"price" example:"200000" validate:"required"`
	DepartureTime string  `json:"departureTime" example:"09:00:00" validate:"required"`
}

// swagger:model ErrorCreateRouteValidation
type ErrorCreateRouteValidation struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"Field from, to, price, departureTime cannot be empty!"`
}

// swagger:model AllRoutes
type AllRoutes struct {
	Count       int     `json:"count" example:"1"`
	TotalPages  int64   `json:"totalPages" example:"10"`
	CurrentPage int64   `json:"currentPage" example:"1"`
	TotalCount  int64   `json:"totalCount" example:"1000"`
	Data        []Route `json:"data"`
}

// swagger:model UpdateRouteRequest
type UpdateRouteRequest struct {
	Price         float64 `json:"price" example:"200000" validate:"required"`
	DepartureTime string  `json:"departureTime" example:"09:00:00" validate:"required"`
}

// swagger:model RouteEditResponseSuccess
type RouteEditResponseSuccess struct {
	Message string `json:"message" example:"Route was updated successfully!"`
}

// swagger:model ErrorUpdateRouteValidation
type ErrorUpdateRouteValidation struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"Field price, departureTime cannot be empty!"`
}

// swagger:model RouteImportResponseSuccess
type RouteImportResponseSuccess struct {
	Message string `json:"message" example:"Route data has successfully imported!"`
}

// for ticket sections

// swagger:model Ticket
type Ticket struct {
	Id            string  `json:"_id,omitempty" bson:"_id,omitempty" example:"6336185fc31ad7ad4022ab87"`
	From          string  `json:"from" example:"Malang" validate:"required"`
	To            string  `json:"to" example:"Jakarta" validate:"required"`
	Price         float64 `json:"price" example:"200000" validate:"required"`
	DepartureTime string  `json:"departureTime" example:"2022-10-01T12:00:00.000Z" validate:"required"`
	UserId        string  `json:"userId,omitempty" example:"632169497f0236bfb3e85412"`
	BookingId     string  `json:"bookingId,omitempty" example:"6336185fc31ad7ad4022ab87"`
	CreatedAt     string  `json:"createdAt" example:"2022-09-30T05:12:47.469Z"`
	UpdatedAt     string  `json:"updatedAt" example:"2022-09-30T05:12:47.469Z"`
}

// swagger:model CreateTicketRequest
type CreateTicketRequest struct {
	From          string  `json:"from" example:"Malang" validate:"required"`
	To            string  `json:"to" example:"Jakarta" validate:"required"`
	Price         float64 `json:"price" example:"200000" validate:"required"`
	DepartureTime string  `json:"departureTime" example:"2022-10-01T12:00:00.000Z" validate:"required"`
	UserId        string  `json:"userId,omitempty" example:"632169497f0236bfb3e854754"`
	BookingId     string  `json:"bookingId,omitempty" example:"6336185fc31ad7ad4022ab87"`
}

// swagger:model ErrorCreateTicketValidation
type ErrorCreateTicketValidation struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"Field from, to, userId, price, bookingId, departureTime cannot be empty!"`
}

// swagger:model TicketEditResponseSuccess
type TicketEditResponseSuccess struct {
	Message string `json:"message" example:"Ticket was updated successfully!"`
}

// swagger:model TicketDeleteResponseSuccess
type TicketDeleteResponseSuccess struct {
	Message string `json:"message" example:"Ticket with id 631ea95d0770f442fd692fa8 was deleted successfully!"`
}

// swagger:model AllTickets
type AllTickets struct {
	Count       int      `json:"count" example:"1"`
	TotalPages  int64    `json:"totalPages" example:"10"`
	CurrentPage int64    `json:"currentPage" example:"1"`
	TotalCount  int64    `json:"totalCount" example:"1000"`
	Data        []Ticket `json:"data"`
}

// swagger:model CreateBookingRequest
type CreateBookingRequest struct {
	RouteId      string `json:"routeId,omitempty" example:"633c692ac31ad7ad4062d0fd" validate:"required"`
	Quantity     int    `json:"quantity" example:"1" validate:"required"`
	ScheduleDate string `json:"scheduleDate" example:"2022-10-01" validate:"required"`
}

// swagger:model BookingUser
type BookingUser struct {
	Id    string `json:"_id,omitempty" bson:"_id,omitempty" example:"6336185fc31ad7ad4022ab87"`
	Email string `json:"email" example:"user@email.com"`
	Name  string `json:"name" example:"Regular User"`
}

// swagger:model Booking
type Booking struct {
	Id            string      `json:"_id,omitempty" bson:"_id,omitempty" example:"6336185fc31ad7ad4022ab87"`
	Quantity      int         `json:"quantity" example:"1"`
	TotalPrice    float64     `json:"totalPrice" example:"500000"`
	PaymentStatus int         `json:"paymentStatus" example:"0"`
	DepartureTime string      `json:"departureTime" example:"2022-09-30T05:12:47.469Z"`
	User          BookingUser `json:"user"`
	Tickets       []Ticket    `json:"tickets"`
	CreatedAt     string      `json:"createdAt" example:"2022-09-30T05:12:47.469Z"`
	UpdatedAt     string      `json:"updatedAt" example:"2022-09-30T05:12:47.469Z"`
}

// swagger:model ErrorRouteNotFound
type ErrorRouteNotFound struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"Route Not Found!"`
}

// swagger:model AllBookings
type AllBookings struct {
	Count       int       `json:"count" example:"1"`
	TotalPages  int64     `json:"totalPages" example:"10"`
	CurrentPage int64     `json:"currentPage" example:"1"`
	TotalCount  int64     `json:"totalCount" example:"1000"`
	Data        []Booking `json:"data"`
}

// swagger:model ErrorBookingPaymentStatusValidation
type ErrorBookingPaymentStatusValidation struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"Error, complete payment only can be done when payment status is 0 (Created)"`
}

// swagger:model BookingEditRequest
type BookingEditRequest struct {
	Id string `json:"_id,omitempty" bson:"_id,omitempty" example:"6336185fc31ad7ad4022ab87" validate:"required"`
}
