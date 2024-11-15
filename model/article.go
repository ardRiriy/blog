package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Article struct {
	bun.BaseModel `bun:"table:articles"`

	Name      string    `bun:",pk"` // urlSuffixを持つ(key)
	FilePath  string    `bun:"file_path,notnull"`
	Tags      []string  `bun:"type:text[]"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}
