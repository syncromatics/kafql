// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package models

import (
	"time"
)

type Record struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	Offset    int       `json:"offset"`
}

type Topic struct {
	Name string `json:"name"`
}
