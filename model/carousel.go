package model

import "time"

type Carousel struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	ImgPath   string
	WebView   string
	CreatTime time.Time
}
