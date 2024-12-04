package serializer

import (
	"Fire/model"
	"time"
)

type Carousel struct {
	ID        int       `json:"id"`
	ImgPath   string    `json:"img_path"`
	WebView   string    `json:"web_view"`
	CreatTime time.Time `json:"creat_time"`
}

func BuildCarousel(item *model.Carousel) Carousel {
	return Carousel{
		ID:        item.ID,
		ImgPath:   item.ImgPath,
		WebView:   item.WebView,
		CreatTime: item.CreatTime,
	}
}

func BuildCarousels(items []model.Carousel) (carousels []Carousel) {
	for _, item := range items {
		carousel := BuildCarousel(&item)
		carousels = append(carousels, carousel)
	}
	return carousels
}
