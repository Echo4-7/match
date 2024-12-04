package service

import (
	"Fire/dao"
	"Fire/pkg/e"
	"Fire/pkg/util/log"
	"Fire/serializer"
	"context"
)

type CarouselService struct {
}

func (service *CarouselService) ListPosters(ctx context.Context) serializer.Response { // TODO：minio
	carouselDao := dao.NewCarouselDao(ctx)
	carousels, err := carouselDao.ListCarousel()
	if err != nil {
		log.LogrusObj.Infoln("ListCarousel failed:", err)
		return serializer.HandleError(e.ServerBusy)
	}
	return serializer.BuildListResponse(serializer.BuildCarousels(carousels), len(carousels))
}
