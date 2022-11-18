package main

import (
	"flag"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/pkwiatek6/KineJamesAPI/actions"
	"github.com/pkwiatek6/KineJamesAPI/get"
	"github.com/pkwiatek6/KineJamesAPI/post"
)

func main() {
	initLogging()
	mongoClient := new(actions.MongoClient)
	mongoClient.Database = os.Getenv("Database")
	mongoClient.Collection = os.Getenv("Collection")

	var err error
	mongoClient.Client, err = mongoClient.ConnectDB(os.Getenv("URI"))
	if err != nil {
		log.Panic().Msgf("Could not connect to db. Error: %v, Client: %v", err, mongoClient)
	} else {
		log.Info().Msg("Connected to DB")
	}
	log.Debug().Msgf("mongoClient: %v", mongoClient)
	//There might be some bad security here because I couldn't figure out how to pass the secret from the website to the api so the api is unprotected
	router := gin.Default()
	initRoutes(router, mongoClient)
	router.Run(":8080")
}

func initLogging() {
	isDebugEnabled := flag.Bool("debug", false, "Turns on debug mode")
	flag.Parse()

	if *isDebugEnabled {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Info().Msgf("Log level set to %s", zerolog.GlobalLevel().String())
}

func initRoutes(router *gin.Engine, client *actions.MongoClient) {
	router.Use(dbMidware(client))
	log.Debug().Msgf("Usings dbMidware with %v", client)
	saveHandlers := router.Group("/save")
	{
		saveHandlers.POST("/character", post.SaveCharacter)
	}

	loadHandlers := router.Group("/load")
	{
		loadHandlers.GET("/character/:name/:userid", get.GetCharacterByName)
		loadHandlers.GET("/allcharacters/:userid", get.GetAllCharacters)
		loadHandlers.GET("/allnames/:userid", get.GetAllNames)
	}

}

//dbMidware passes the mongo Client down so i don't tie up resources
//basically if I ever need client in middleware I can get it easily and non-bocking, I think
func dbMidware(client *actions.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("dbClient", client)
		c.Next()
	}
}
