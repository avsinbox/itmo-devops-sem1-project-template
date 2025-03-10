package models

import "time"

type Item struct {
	ID         int64
	Name       string
	Category   string
	Price      float64
	CreateDate time.Time
}
