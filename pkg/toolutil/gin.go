// Author: huaxr
// Time:   2021/12/31 上午10:55
// Git:    huaxr

package toolutil

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetPagination(c *gin.Context) (page int, pageSize int) {
	var err error
	page, err = strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err = strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 0 {
		pageSize = 50
	}
	return
}

func GetOffsetLimit(c *gin.Context) (int, int) {
	page, pageSize := GetPagination(c)
	offset := pageSize * (page - 1)
	return offset, pageSize
}
