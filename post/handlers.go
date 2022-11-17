package post

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkwiatek6/KineJamesAPI/actions"
	"github.com/pkwiatek6/KineJamesAPI/structs"
	"github.com/rs/zerolog/log"
)

func SaveCharacter(ginCtx *gin.Context) {
	var character structs.Character
	client := ginCtx.MustGet("dbClient").(*actions.MongoClient)
	err := ginCtx.BindJSON(&character)
	if err != nil {
		log.Err(err).Msg("Could not bind JSON to Character struct")
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "Could not bind JSON to Character struct"})
		return
	}
	err = client.SaveCharacterToUser(character)
	if err != nil {
		log.Err(err).Msg("Could not save character to user")
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "Could not save Character to database"})
		return
	}
	ginCtx.JSON(http.StatusCreated, gin.H{"status": "Saved Character to database"})

}
