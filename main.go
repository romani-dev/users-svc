package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"log"

	"github.com/gin-gonic/gin"
)

type Input struct {
	Data []Item `json:"data"`
}

type Item struct {
	ID     int    `json:"ID"`
	Name   string `json:"Name"`
	Image  string `json:"Image"`
	UserID int    `json:"UserID"`
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Items []Item `json:"items"`
}

var db map[int]User
var itemsURL string

func init() {

	db = map[int]User{
		1: {ID: 1, Name: "Felipe", Email: "felipe@gmail.com"},
		2: {ID: 2, Name: "Elisa", Email: "elisa@gmail.com"},
	}

	itemsURL = os.Getenv("ITEMS_URL")

	if itemsURL == "" {
		log.Fatal("environment ITEMS_URL can't be empty")
	}
}

func main() {

	r := gin.Default()

	r.GET("/users", func(c *gin.Context) {

		resp, err := http.Get(itemsURL)
		if err != nil {
			panic(err)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != 200 {
			log.Fatal(string(data))
		}

		var input Input
		err = json.Unmarshal(data, &input)
		if err != nil {
			panic(err)
		}

		finalUsers := []User{}
		for _, user := range db {

			for _, item := range input.Data {
				if item.UserID == user.ID {
					user.Items = append(user.Items, item)
				}
			}

			finalUsers = append(finalUsers, user)
		}
		c.JSON(200, finalUsers)
	})

	r.Run(":8089")
}
