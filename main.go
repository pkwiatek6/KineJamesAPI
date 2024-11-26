package main

import (
	"flag"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/pkwiatek6/KineJamesAPI/actions"
	"github.com/pkwiatek6/KineJamesAPI/get"
	"github.com/pkwiatek6/KineJamesAPI/post"
)

var (
	MongoDB_URI     = flag.String("URI", "", "URI of the MongoDB instance")
	Database_Name   = flag.String("dbName", "Characters", "Name of the database to use")
	Collection_Name = flag.String("col", "Sheets", "Name of the collection to use")
)

func main() {
	initLogging()
	mongoClient := &actions.MongoClient{}
	err := mongoClient.ConnectDB(*MongoDB_URI, *Database_Name, *Collection_Name)
	if err != nil {
		log.Panic().Msgf("Could not connect to db. Error: %v, Client: %v", err, mongoClient)
	}
	log.Info().Msg("Connected to DB")

	log.Debug().Msgf("mongoClient: %v", *mongoClient)
	//There might be some bad security here because I couldn't figure out how to pass the secret from the website to the api so the api is unprotected
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "GET", "OPTIONS"},
		AllowHeaders: []string{"content-type,access-control-allow-origin, access-control-allow-headers"},
	}))
	initRoutes(router, mongoClient)
	err = router.Run(":8070")
	if err != nil {
		log.Panic().Msgf("Failed to start Gin server due to: " + err.Error())
	}
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

// dbMidware passes the mongo Client down so i don't tie up resources
// basically if I ever need client in middleware I can get it easily and non-bocking, I think
func dbMidware(client *actions.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("dbClient", client)
		c.Next()
	}
}
