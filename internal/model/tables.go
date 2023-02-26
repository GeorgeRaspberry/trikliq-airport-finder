package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Role            string `json:"role,omitempty"`
	Password        string `json:"password,omitempty"`
	Salutation      string `json:"salutation,omitempty"`
	FirstName       string `json:"firstName,omitempty"`
	LastName        string `json:"lastName,omitempty"`
	Street          string `json:"street,omitempty"`
	ZipCode         string `json:"zipCode,omitempty"`
	City            string `json:"city,omitempty"`
	Email           string `json:"email,omitempty"`
	PhoneNumber     string `json:"phoneNumber,omitempty"`
	Country         string `json:"country,omitempty"`
	DeliveryAddress string `json:"deliveryAddress,omitempty"`
	DeliveryCountry string `json:"deliveryCountry,omitempty"`
	DeliveryStreet  string `json:"deliveryStreet,omitempty"`
	DeliveryZipCode string `json:"deliveryZipCode,omitempty"`
	DeliveryCity    string `json:"deliveryCity,omitempty"`
	UserGroup       string `json:"userGroup,omitempty"`
	Newsletter      string `json:"newsletter,omitempty"`
}
