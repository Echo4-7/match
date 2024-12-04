package v1

import (
	"Fire/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListCarousel(c *gin.Context) {
	var listCarousel service.CarouselService
	if err := c.ShouldBind(&listCarousel); err == nil {
		res := listCarousel.ListPosters(c.Request.Context())
		c.JSON(http.StatusOK, res)
		return
	} else {
		c.JSON(http.StatusBadRequest, err)
		return
	}
}
