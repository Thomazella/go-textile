package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// addBlockFlags godoc
// @Summary Add a flag
// @Description Adds a flag to a thread block
// @Tags blocks
// @Produce application/json
// @Param id path string true "block id"
// @Success 201 {object} pb.Flag "flag"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /blocks/{id}/flags [post]
func (a *api) addBlockFlags(g *gin.Context) {
	id := g.Param("id")

	thrd := a.getBlockThread(g, id)
	if thrd == nil {
		return
	}

	hash, err := thrd.AddFlag(id)
	if err != nil {
		a.abort500(g, err)
		return
	}

	flag, err := a.node.Flag(hash.B58String())
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusCreated, flag)
}

// lsBlockFlags godoc
// @Summary List flags
// @Description Lists flags on a thread block
// @Tags blocks
// @Produce application/json
// @Param id path string true "block id"
// @Success 200 {object} pb.FlagList "flags"
// @Failure 500 {string} string "Internal Server Error"
// @Router /blocks/{id}/flags [get]
func (a *api) lsBlockFlags(g *gin.Context) {
	id := g.Param("id")

	flags, err := a.node.Flags(id)
	if err != nil {
		a.abort500(g, err)
		return
	}

	pbJSON(g, http.StatusOK, flags)
}

// getBlockFlag godoc
// @Summary Get thread flag
// @Description Gets a thread flag by block ID
// @Tags blocks
// @Produce application/json
// @Param id path string true "block id"
// @Success 200 {object} pb.Flag "flag"
// @Failure 400 {string} string "Bad Request"
// @Router /blocks/{id}/flag [get]
func (a *api) getBlockFlag(g *gin.Context) {
	info, err := a.node.Flag(g.Param("id"))
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	pbJSON(g, http.StatusOK, info)
}
