package address

import "time"

type Address struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	AddressLine1 string    `json:"address_line_1"`
	AddressLine2 string    `json:"address_line_2"`
	PostalCode   string    `json:"postal_code"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	Country      string    `json:"country"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AddressReq struct {
	ID           int64  `json:"id,omitempty"`
	UserID       int64  `json:"user_id,omitempty"`
	AddressLine1 string `json:"address_line_1" validate:"required,min=5,max=250"`
	AddressLine2 string `json:"address_line_2" validate:"max=250"`
	PostalCode   string `json:"postal_code" validate:"required,min=3,max=100"`
	City         string `json:"city" validate:"required,min=3,max=100"`
	State        string `json:"state" validate:"required,min=3,max=100"`
	Country      string `json:"country" validate:"required,min=3,max=100"`
}

type AddressRes struct {
	Count     int        `json:"count"`
	Addresses *[]Address `json:"addresses"`
}
