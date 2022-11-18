package actions

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/pkwiatek6/KineJamesAPI/structs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Client     *mongo.Client
	Database   string
	Collection string
}

func (c MongoClient) ConnectDB(URI string) (*mongo.Client, error) {
	if URI == "" {
		return nil, errors.New("URI env variable not set")
	}
	clientOptions := options.Client().ApplyURI(URI)

	returnClient, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	err = returnClient.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}
	return returnClient, nil
}

//Returns a document from connected db looking in c.Collection from c.Database
//name string is the name of the character you are looking for
// user string is the user that is looking for the character
// returns mongo.ErrNoDocuments if nothing is found, if multiple are found it returns one
func (c MongoClient) GetCharacterByName(name string, user string) (*structs.Character, error) {
	filter := bson.M{"name": name, "user": user}
	log.Debug().Msgf("filter: %v", filter)
	log.Debug().Msgf("c values: %v", c)
	collection := c.Client.Database(c.Database).Collection(c.Collection)
	var character structs.Character
	err := collection.FindOne(context.TODO(), filter).Decode(&character)
	if err != nil {
		return nil, err
	} else {
		return &character, nil
	}
}

//Gets all Characters that are from the give user
func (c MongoClient) GetAllCharactersFromPlayer(user string) (map[string]*structs.Character, error) {
	var results []structs.Character
	var toReturn = make(map[string]*structs.Character)
	filter := bson.M{"user": user}
	log.Debug().Msgf("c values: %v", c)
	cursor, err := c.Client.Database(c.Database).Collection(c.Collection).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, err
	}
	for key, character := range results {
		toReturn[character.User] = &results[key]
	}
	return toReturn, nil
}

//Gets all Characters names that are from the give user.
//Used if you only want names, uses less memory than GetAllCharacters
//Returns an array of strings containing all names
func (c MongoClient) GetAllCharacterNamesFromPlayer(user string) ([]string, error) {
	var results []structs.Character
	filter := bson.M{"user": user}
	log.Debug().Msgf("c values: %v", c)
	cursor, err := c.Client.Database(c.Database).Collection(c.Collection).Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("results values: %v", results)
	var toReturn = make([]string, len(results), cap(results))
	for key, character := range results {
		log.Debug().Msgf("Key %v, Value %v", key, character.Name)
		toReturn[key] = character.Name
	}
	return toReturn, nil
}

//Saves the given Character struct to the db, if there wasn't a character sheet it makes one
func (c MongoClient) SaveCharacterToUser(character structs.Character) error {
	collection := c.Client.Database(c.Database).Collection(c.Collection)
	filter := bson.M{"name": character.Name, "user": character.User}
	update := bson.M{"$set": character}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	} else if updateResult.MatchedCount == 0 {
		log.Debug().Msgf("Filter %v returned no results", filter)
		log.Info().Msg("Failed to find matching document, making a new one")
	} else if updateResult.MatchedCount == 1 {
		log.Debug().Msgf("Character %v was updated with values %v", character.Name, update)
		log.Info().Msg("Document Was updated")
		return nil
	}
	//Creates a new document if there wasn't one already
	insertResult, err := collection.InsertOne(context.TODO(), character)
	if err != nil {
		return err
	}
	log.Debug().Msgf("Inserted post with ID:%v", insertResult.InsertedID)
	return nil
}
