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

func SetTracyHeadersRequest(req *http.Request, c *gin.Context) {
	reqID := c.GetHeader("x-request-id")
	traceID := c.GetHeader("x-b3-traceid")
	spanID := c.GetHeader("x-b3-spanid")
	parentSpanID := c.GetHeader("x-b3-parentspanid")
	sampled := c.GetHeader("x-b3-sampled")
	flags := c.GetHeader("x-b3-flags")
	spanContext := c.GetHeader("x-ot-span-context")

	req.Header.Set("x-request-id", reqID)
	req.Header.Set("x-b3-traceid", traceID)
	req.Header.Set("x-b3-spanid", spanID)
	req.Header.Set("x-b3-parentspanid", parentSpanID)
	req.Header.Set("x-b3-sampled", sampled)
	req.Header.Set("x-b3-flags", flags)
	req.Header.Set("x-ot-span-context", spanContext)
}

func SetTracyHeadersResponse(c *gin.Context) {
	reqID := c.GetHeader("x-request-id")
	traceID := c.GetHeader("x-b3-traceid")
	spanID := c.GetHeader("x-b3-spanid")
	parentSpanID := c.GetHeader("x-b3-parentspanid")
	sampled := c.GetHeader("x-b3-sampled")
	flags := c.GetHeader("x-b3-flags")
	spanContext := c.GetHeader("x-ot-span-context")

	c.Header("x-request-id", reqID)
	c.Header("x-b3-traceid", traceID)
	c.Header("x-b3-spanid", spanID)
	c.Header("x-b3-parentspanid", parentSpanID)
	c.Header("x-b3-sampled", sampled)
	c.Header("x-b3-flags", flags)
	c.Header("x-ot-span-context", spanContext)
}

func main() {

	r := gin.Default()

	r.GET("/users", func(c *gin.Context) {

		req, err := http.NewRequest("GET", itemsURL, nil)
		if err != nil {
			panic(err)
		}
		SetTracyHeadersRequest(req, c)
		client := &http.Client{}
		resp, err := client.Do(req)

		// resp, err := http.Get(itemsURL)
		if err != nil {
			panic(err)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != 200 {
			panic(string(data))
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

		SetTracyHeadersResponse(c)
		c.JSON(200, finalUsers)
	})

	r.GET("/headers", func(c *gin.Context) {
		c.JSON(200, c.Request.Header)
	})

	r.Run(":8089")
}
