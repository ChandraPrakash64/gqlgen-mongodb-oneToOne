// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Address struct {
	ID     string  `json:"_id"`
	City   *string `json:"city,omitempty"`
	Pin    *string `json:"pin,omitempty"`
	UserID *string `json:"userId,omitempty"`
}

type Mutation struct {
}

type Query struct {
}

type User struct {
	ID      string   `json:"_id"`
	Name    *string  `json:"name,omitempty"`
	Address *Address `json:"address,omitempty"`
}
