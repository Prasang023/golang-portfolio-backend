package main

import (
"log"
"os"
"github.com/joho/godotenv"
"github.com/Prasang023/mongo-go/controllers"
"context"
"fmt"
"go.mongodb.org/mongo-driver/bson"
"go.mongodb.org/mongo-driver/mongo"
"github.com/gin-gonic/gin"
"go.mongodb.org/mongo-driver/mongo/options"
)	

func main(){

	// r := httprouter.New()
	err := godotenv.Load(".env")
	if err != nil{
	log.Fatalf("Error loading .env file: %s", err)
	}
	router := gin.Default()
	router.Use(CORSMiddleware())
	uc := controllers.NewAPIController(getSession())
	router.GET("/api/", func(c *gin.Context){
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/api/projects", uc.GetProjects)
	router.GET("/api/project/:id", uc.GetProjectById)
	// r.GET("/api/projects/:id", uc.GetProjectById )
	// r.POST("/user", uc.CreateUser)
	// r.DELETE("/user/:id", uc.DeleteUser)
	// http.ListenAndServe("localhost:9000", r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func getSession() (*mongo.Client){
	mongo_uri := os.Getenv("MONGOURI")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongo_uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	fmt.Println(client)
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 	panic(err)
	// 	}
	// }()
	// Send a ping to confirm a successful connection
	if err := client.Database("test").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client
}