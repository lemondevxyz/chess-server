package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserFromCtx(c *gin.Context) User {
	ctx := &context{
		Context: c,
	}
	u := User{}

	err := ctx.ParseJSON(&u)
	if err != nil {
		return User{}
	}

	return u
}

func UserWall(h handler) gin.HandlerFunc {
	return func(c *gin.Context) {

		u := GetUserFromCtx(c)
		cl := GetClient(u)
		if u == (User{}) || cl == nil {
			c.AbortWithError(http.StatusForbidden, ErrClient)

			return
		}

		ctx := &context{Context: c, c: cl}
		h(ctx)
	}
}

func GameWall(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, ErrClient)
	}

	u := User{}
	err = json.Unmarshal(body, &u)
	if err != nil {
		c.AbortWithError(http.StatusForbidden, ErrClient)
	}

	if GetClient(u) == nil {
		c.AbortWithError(http.StatusForbidden, ErrClient)
	}
}
