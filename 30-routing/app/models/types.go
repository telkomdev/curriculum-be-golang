package models

// swagger:model ResponseRoot
type ResponseRoot struct {
	Error   bool   `json:"error" example:"false"`
	Message string `json:"message" example:"Hello. Try GET to /api/v1/item"`
}

// swagger:model ResponseItem
type ResponseItem struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"something when wrong"`
}
