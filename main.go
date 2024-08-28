package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://andy:wUkXAQBDu6pIRajJ@cluster0.nxesq.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx) // disconnect the db after the code is terminated

	// err = client.Ping(ctx, readpref.Primary())		//connect to the online via the provided link
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// databases, err := client.ListDatabaseNames(ctx, bson.M{})		//checkavailable databases
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(databases)

	quickstartDatabase := client.Database("quickstart")
	podcastsCollection := quickstartDatabase.Collection("podcast")
	episodesCollection := quickstartDatabase.Collection("episodes")

	// podcastResult, err := podcastsCollection.InsertOne(ctx, bson.D{
	// 	{"title", "the polyglot Developer Podcast"},
	// 	{"author", "Nic Raboy"},
	// 	{"tags", bson.A{"development", "programming", "coding"}},
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(podcastResult.InsertedID)

	// episodeResult, err := episodesCollection.InsertMany(ctx, []interface{}{
	// 	bson.D{
	// 		{"podcast", podcastResult.InsertedID},
	// 		{"title", "episode #1"},
	// 		{"description", "This is the first episode"},
	// 		{"duration", 24},
	// 	},
	// 	bson.D{
	// 		{"podcast", podcastResult.InsertedID},
	// 		{"title", "episode #2"},
	// 		{"description", "This is the second episode"},
	// 		{"duration", 34},
	// 	},
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(episodeResult.InsertedIDs...)

	cursor, err := episodesCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	// this is not to be used for larger documents with many items
	// var episodes []bson.M
	// if err = cursor.All(ctx, &episodes); err != nil {
	// 	log.Fatal(err)
	// }
	// // fmt.Println(episodes)
	// for _, episode := range episodes {
	// 	fmt.Println(episode["title"])
	// }
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var episode bson.M
		if err = cursor.Decode(&episode); err != nil {
			log.Fatal(err)
		}
		// fmt.Println(episode)
	}

	var podcast bson.M
	if err = podcastsCollection.FindOne(ctx, bson.M{}).Decode(&podcast); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(podcast)

	filterCursor, err := episodesCollection.Find(ctx, bson.M{"duration": 34})
	if err != nil {
		log.Fatal(err)
	}
	var episodesFiltered []bson.M
	if err = filterCursor.All(ctx, &episodesFiltered); err != nil {
		log.Fatal(err)
	}
	// fmt.Println(episodesFiltered)

	opts := options.Find()
	opts.SetSort(bson.D{{"duration", 1}})
	sortCursor, err := episodesCollection.Find(ctx, bson.D{
		{"duration", bson.D{
			{"$gt", 23},
		}},
	}, opts)
	if err != nil {
		log.Fatal(err)
	}
	var episodesSorted []bson.D
	if err = sortCursor.All(ctx, &episodesSorted); err != nil {
		log.Fatal(err)
	}
	fmt.Println(episodesSorted)
}
