package handlers

import (
	"net/http"

	"github.com/reagent/copyright/app"
	"github.com/reagent/copyright/models"
)

func AccountPost(ctx *app.Context) {
	var (
		err     error
		account models.Account
	)

	err = ctx.UnmarshalTo(&account)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.H{"message": err.Error()})
		return
	}

	err = account.Create(ctx.DB())

	if err != nil {
		ctx.JSON(http.StatusBadRequest, app.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

func AccountGet(ctx *app.Context) {
	id, ok := ctx.Param("id")

	if ok {
		if account, _ := models.Find(ctx.DB(), id); account != nil {
			ctx.JSON(http.StatusOK, account)
			return
		}
	}

	ctx.JSON(http.StatusNotFound, app.H{"message": "Not found"})
}
