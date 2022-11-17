package get

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkwiatek6/KineJamesAPI/actions"
	"github.com/rs/zerolog/log"
)

func GetCharacterByName(ginCtx *gin.Context) {
	userid := ginCtx.Param("userid")
	name := ginCtx.Param("name")
	client := ginCtx.MustGet("dbClient").(*actions.MongoClient)
	character, err := client.GetCharacterByName(name, userid)
	if err != nil {
		log.Err(err).Msgf("Could not find Character %v, from user %v", name, userid)
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "Could find Character"})
		return
	}
	ginCtx.JSON(http.StatusOK, character)
}

func GetAllCharacters(ginCtx *gin.Context) {
	userid := ginCtx.Param("userid")
	client := ginCtx.MustGet("dbClient").(*actions.MongoClient)
	characters, err := client.GetAllCharactersFromPlayer(userid)
	if err != nil {
		log.Err(err).Msgf("Could not find Characters from user %v", userid)
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "Could find Characters from user"})
		return
	}
	ginCtx.JSON(http.StatusOK, characters)
}

func GetAllNames(ginCtx *gin.Context) {
	userid := ginCtx.Param("userid")
	client := ginCtx.MustGet("dbClient").(*actions.MongoClient)
	names, err := client.GetAllCharacterNamesFromPlayer(userid)
	if err != nil {
		log.Err(err).Msgf("Could not find Character names from user %v", userid)
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": "Could find Character names from user"})
		return
	}
	ginCtx.JSON(http.StatusOK, names)
}
