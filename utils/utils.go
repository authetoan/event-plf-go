package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseQueryParamAsInt(c *gin.Context, param string, defaultValue int) int {
	valueStr := c.Query(param)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
