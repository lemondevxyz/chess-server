package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toms1441/chess/internal/game"
)

type handler func(c *context)

type context struct {
	*gin.Context
	c *game.Client
	g *game.Game
}

func (c *context) ParseJSON(obj interface{}) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   fmt.Errorf("cannot read body"),
		})

		return err
	}

	err = json.Unmarshal(body, obj)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Errorf("json: %w", err),
		})

		return err
	}

	err = v.Struct(obj)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Errorf("validator: %w", err),
		})

		return err
	}

	return nil
}
