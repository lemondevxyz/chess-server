package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toms1441/chess/internal/game"
)

func cmdHandler(c *context) {
	cmd := game.Command{}

	err := c.ParseJSON(cmd)
	if err != nil {
		return
	}

	cl := c.c

	err = cl.Do(cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Errorf("command: %w", err),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
