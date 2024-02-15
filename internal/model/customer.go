package model

import "time"

type Customer struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	DateCreated  time.Time `json:"date_created"`
	DateModified time.Time `json:"date_modified"`
}
