package register

import (
	"net/http"

	"github.com/TheLazarusNetwork/go-helpers/httpo"
	"github.com/TheLazarusNetwork/go-helpers/logo"
	"github.com/Weareflexable/Superiad/models/user"

	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/register")
	{

		g.GET("", register)
	}
}

func register(c *gin.Context) {
	uid, err := user.AddUser()
	if err != nil {
		httpo.NewErrorResponse(http.StatusInternalServerError, "failed to add user").SendD(c)
		logo.Errorf("failed to add user, error: %s", err)

	} else {
		payload := RegisterPayload{
			Uid: uid,
		}
		httpo.NewSuccessResponseP(200, "user registration successfull", payload).SendD(c)
	}
}
