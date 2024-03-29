package isowner

import (
	"errors"
	"math/big"
	"net/http"

	"github.com/TheLazarusNetwork/go-helpers/httpo"
	"github.com/TheLazarusNetwork/go-helpers/logo"
	"github.com/Weareflexable/Superiad/models/user"
	"github.com/Weareflexable/Superiad/pkg/network/polygon"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/isowner")
	{

		g.POST("", isowner)
	}
}

func isowner(c *gin.Context) {
	var req IsOwnerRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		httpo.NewErrorResponse(http.StatusBadRequest, "body is invalid").SendD(c)

		return
	}
	mnemonic, err := user.GetMnemonic(req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpo.NewErrorResponse(httpo.UserNotFound, "user not found").Send(c, 404)

			return
		}
		httpo.NewErrorResponse(http.StatusInternalServerError, "failed to fetch user").SendD(c)
		logo.Errorf("failed to fetch user with id %v, err %s", req.UserId, err)
		return
	}

	isOwner, err := polygon.ERC721IsOwner(mnemonic, common.HexToAddress(req.ContractAddress), big.NewInt(req.TokenId))

	if err != nil {
		httpo.NewErrorResponse(http.StatusInternalServerError, "failed to call ERC721IsOwner").SendD(c)
		logo.Errorf("failed to call ERC721IsOwner for user with id %v,erc721Address %v,tokenId %v, err %s", req.UserId, req.ContractAddress, req.TokenId, err)
		return
	}
	payload := IsOwnerPayload{
		IsOwner: isOwner,
	}
	httpo.NewSuccessResponseP(200, "Result success", payload).SendD(c)

}
