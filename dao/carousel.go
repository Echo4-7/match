package dao

import (
	"Fire/model"
	"context"
	"gorm.io/gorm"
)

type CarouselDao struct {
	*gorm.DB
}

func NewCarouselDao(ctx context.Context) *CarouselDao {
	return &CarouselDao{NewDBClient(ctx)}
}

func NewCarouselByDB(db *gorm.DB) *CarouselDao {
	return &CarouselDao{db}
}

func (dao *CarouselDao) ListCarousel() (carousels []model.Carousel, err error) { // TODO
	err = dao.DB.Model(&model.Carousel{}).Find(&carousels).Error
	return
}
