package main

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Email string             `bson:"email"`
}

func main() {
	// Set up MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb+srv://moiz697:jijjkklkikk129@cluster0.xdxgkmq.mongodb.net/test")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("test")
	userCollection := db.Collection("users")

	// Set up Gin router
	router := gin.Default()

	// Create user endpoint
	router.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := userCollection.InsertOne(context.Background(), user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		user.ID = result.InsertedID.(primitive.ObjectID)
		c.JSON(http.StatusCreated, user)
	})

	// Get all users endpoint
	router.GET("/users", func(c *gin.Context) {
		var users []User
		cursor, err := userCollection.Find(context.Background(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(context.Background())
		for cursor.Next(context.Background()) {
			var user User
			if err := cursor.Decode(&user); err != nil {
				log.Fatal(err)
			}
			users = append(users, user)
		}
		c.JSON(http.StatusOK, users)
	})

	// Get user by ID endpoint
	router.GET("/users/:id", func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user User
		err = userCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	})
	// Update user endpoint
	router.PUT("/users/:id", func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		update := bson.M{
			"$set": bson.M{
				"name":  user.Name,
				"email": user.Email,
			},
		}
		_, err = userCollection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})

	// Delete user endpoint
	router.DELETE("/users/:id", func(c *gin.Context) {
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err = userCollection.DeleteOne(context.Background(), bson.M{"_id": id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})

	log.Fatal(router.Run(":8888"))

}
