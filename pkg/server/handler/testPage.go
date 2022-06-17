package handler

import (
	"github.com/ankogit/go-telegram-bot-template/pkg/server/handler/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getTestPage(c *gin.Context) {
	c.JSON(http.StatusOK, response.DataResponse{Data: "Hello world"})
}
