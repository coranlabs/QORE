package sbi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getIndexroutes() []Route {
	return []Route{
		{
			"Index",
			http.MethodGet,
			"/",
			func(c *gin.Context) {
				c.JSON(http.StatusOK, "coranlabs")
			},
		},
	}
}
