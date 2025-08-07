package models

import (
	"github.com/uptrace/bun"
	"time"
)

type File struct {
	bun.BaseModel `bun:"table:files"`

	Name     string `json:"name" bun:"name,notnull"`
	FileSize int64  `json:"file_size" bun:"file_size,notnull"`
	FileType string `json:"file_type" bun:"file_type,notnull"`
	FilePath string `json:"file_path" bun:"file_path,notnull"`
	Info     string `json:"info" bun:"info"`

	UploadBy int64 `json:"upload_by" bun:"upload_by"`
	User     *User `json:"user" bun:"rel:belongs-to,join:upload_by=id"`

	CreatedAt time.Time  `json:"created_at" bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time  `json:"updated_at" bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at,nullzero"`
}
