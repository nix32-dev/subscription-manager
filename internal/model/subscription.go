package model

type Subscription struct {
	ID      int    `json:"ID" swaggerignore:"true"`
	User_id string `json:"user_id" swaggerignore:"true"`

	// @example Yandex Plus
	Service_name *string `json:"service_name"`

	// @example 1000
	// @minimum 1
	Price *int `json:"price"`

	// @example 09-2023
	Start_date *string `json:"start_date"`

	// @example 01-2024
	End_date *string `json:"end_date"`
}
