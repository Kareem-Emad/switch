package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	faktory "github.com/contribsys/faktory/client"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Message structure expected from publishers
type Message struct {
	Payload string
	Topic   string
}

// Subscriber the struct holding info about each topic subscriber
type Subscriber struct {
	gorm.Model
	FilterExpression string
	TargetURL        string
	Topic            string
}

// Claims the struct holding the expected fields encoded in incoming jwt tokens
type Claims struct {
	Author string `json:"author"`
	jwt.StandardClaims
}

func fetchAllSubscribers() []Subscriber {
	var subs []Subscriber

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("[Switch_DB_Driver]  failed to connect database | error: %s", err))
	}

	db.AutoMigrate(&Subscriber{})
	// db.Create(&Subscriber{FilterExpression: "True", TargetURL: "https://09d60f3cb7c49e6e36d7885105edab9c.m.:domain.net", Topic: "test"})
	db.Find(&subs)
	fmt.Println(fmt.Sprintf("[Switch_DB_Driver] fetched all subscribers | total count %d", len(subs)))
	return subs
}
func preLoadSubscribers() map[string][]Subscriber {
	subscriberTopicGroups := make(map[string][]Subscriber)

	subs := fetchAllSubscribers()

	for _, sub := range subs {
		subscriberTopicGroups[sub.Topic] = append(subscriberTopicGroups[sub.Topic], sub)
	}
	return subscriberTopicGroups
}

func isValidJwtToken(token string) bool {
	return !(len(token) <= 0 || strings.Count(token, ".") != 2)
}

func verifyTokenSignature(token string) bool {

	var jwtKey = []byte("my_secret_key")
	claims := &Claims{}
	if !isValidJwtToken(token) {
		return false
	}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !tkn.Valid {
		return false
	}

	fmt.Println(fmt.Sprintf("[Switch_Auth] token verified for author %s", claims.Author))
	return true
}

func authenticate(headers http.Header) bool {
	return verifyTokenSignature(getjwtToken(headers))
}

func getjwtToken(headers http.Header) string {
	authStringSlices := strings.Split(headers.Get("Authorization"), "bearer ")

	if len(authStringSlices) < 2 {
		return ""
	}
	return authStringSlices[1]
}

func main() {

	r := mux.NewRouter()
	fmt.Println("[Switch] Starting server ....")

	subscriberTopicGroups := preLoadSubscribers()
	faktoryClient, err := faktory.Open()
	if err != nil {
		panic(fmt.Sprintf("[Switch_Faktory_Producer] Failed to establish connection to faktory server | error: %s", err))
	}

	r.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		var message Message

		if authenticate(r.Header) == false {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(fmt.Sprintf("[Switch] getting subcribers for topic %s | found %d subscribers", message.Topic, len(subscriberTopicGroups[message.Topic])))

		for _, sub := range subscriberTopicGroups[message.Topic] {
			job := faktory.NewJob("my_queue", sub.TargetURL, sub.FilterExpression, message.Payload)
			err = faktoryClient.Push(job)

			if err != nil {
				fmt.Println(fmt.Sprintf("[Switch_Faktory_Producer] Failed to push new job in queue | error: %s", err))
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Println(fmt.Sprintf("[Switch_Faktory_Producer] Sucessfully insereted job in queue with specs {%s, %s, %s}", sub.TargetURL, sub.FilterExpression, message.Payload))
		}
	}).Methods("POST")

	fmt.Println("[Switch]  Listening to port 8080")
	http.ListenAndServe(":8080", r)
}
