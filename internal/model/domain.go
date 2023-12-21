package model

import "time"

type User struct {
	ID               string
	Version          int32
	Login            string
	Enabled          bool
	HashedPassword   string
	Permissions      []string
	CreatedDate      time.Time
	LastModifiedDate time.Time
}
