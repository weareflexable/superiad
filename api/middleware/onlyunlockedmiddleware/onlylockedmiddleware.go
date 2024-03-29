package onlyunlockedmiddleware

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/TheLazarusNetwork/go-helpers/httpo"
	"github.com/TheLazarusNetwork/go-helpers/logo"
	"github.com/Weareflexable/Superiad/models/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func OnlyUnlocked() gin.HandlerFunc {
	return func(c *gin.Context) {
		ByteBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			httpo.NewErrorResponse(500, "failed to read body").SendD(c)
			logo.Errorf("failed to read body from request: %s", err)
			c.Abort()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(ByteBody))
		var req AccessRequest
		err = c.ShouldBindJSON(&req)
		if err != nil {
			err := fmt.Errorf("body is invalid: %w", err)
			httpo.NewErrorResponse(http.StatusBadRequest, err.Error()).SendD(c)
			c.Abort()
			return
		}
		_user, err := user.GetUser(req.UserId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				httpo.NewErrorResponse(httpo.UserNotFound, "user not found").Send(c, 404)
				c.Abort()
				return
			}

			httpo.NewErrorResponse(500, "failed to fetch user").SendD(c)
			logo.Errorf("failed to fetch user with id: %s, err: %s", req.UserId, err)
			c.Abort()
			return
		}

		if _user.IsUserLocked {
			httpo.NewErrorResponse(httpo.UserLocked, "cannot perform this operation when user is locked").Send(c, http.StatusForbidden)
			c.Abort()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(ByteBody))

		c.Next()
	}
}
