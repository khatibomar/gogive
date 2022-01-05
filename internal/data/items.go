package data

import "time"

type Item struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"-"`
	Title      string    `json:"title"`
	Categories []string  `json:"categories,omitempty"`
	Version    int32     `json:"version"`
}