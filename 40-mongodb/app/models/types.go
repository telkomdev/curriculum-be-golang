package models

// swagger:model ResponseRoot
type ResponseRoot struct {
	Error   bool   `json:"error" example:"false"`
	Message string `json:"message" example:"Hello. Try GET/POST to /api/v1/item"`
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
	Count int    `json:"count" example:"1"`
	Data  []Item `json:"data"`
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
