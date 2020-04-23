package main

import (
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
)

var redisClient *redis.Client

func main() {
	addPin("1", "Pin 1")
	addPin("2", "Pin 2")
	addPin("3", "Pin 3")
	addPin("4", "Pin 4")
	addPin("5", "Pin 5")
	addPin("6", "Pin 6")
	addPin("7", "Pin 7")
	addPin("8", "Pin 8")
	addPin("9", "Pin 9")
	addPin("10", "Pin 10")
	addPin("11", "Pin 11")
	addPin("12", "Pin 12")
	addPin("13", "Pin 13")

	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6320",
		DB:   0,
	})

	r := setupRouter()

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

type Pin struct {
	ID    string
	Name  string
	Type  string
	State bool
	Order int
}

var pins = make(map[string]*Pin)

type UpdateRequest struct {
	Items []struct {
		ID    string `json:"id"`
		State bool   `json:"state"`
	} `json:"items"`
}

func addPin(id, name string) {
	pins[id] = &Pin{ID: id, Name: name, Order: len(pins)}
}

func getPins() []*Pin {
	result := make([]*Pin, 0)

	for _, p := range pins {
		if err := redisClient.HGet("raspi", p.ID).Scan(&p.State); err != nil {
			if err != redis.Nil {
				log.Println("read state", err)
			}
		}
		result = append(result, p)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Order < result[j].Order
	})

	return result
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.GET("/raspi/pins", func(c *gin.Context) {
		p := make([]gin.H, 0)

		for _, pin := range getPins() {
			p = append(p, gin.H{
				"id":    pin.ID,
				"name":  pin.Name,
				"state": pin.State,
			})
		}

		c.JSON(http.StatusOK, gin.H{"items": p})
	})

	r.POST("/raspi/pins", func(c *gin.Context) {
		var json UpdateRequest

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, p := range json.Items {
			if pin, ok := pins[p.ID]; ok {
				pin.State = p.State
				redisClient.HSet("raspi", pin.ID, pin.State)
			}
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	return r
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
