package helpers

import "github.com/gin-gonic/gin"

type error struct {
	Error string
}

func MyAbort(c *gin.Context, str string) {
	c.AbortWithStatusJSON(400, error{Error: str})
}
