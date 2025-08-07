package models

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID       int64  `json:"id" bun:"id,pk,autoincrement"`
	UserName string `json:"username" bun:"username,unique"`
	Password string `json:"password" bun:"password,notnull"`

	Files []File `json:"files" bun:"rel:has-many,join:id=upload_by"`

	CreatedAt time.Time  `json:"created_at" bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time  `json:"updated_at" bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at,nullzero"`
}
