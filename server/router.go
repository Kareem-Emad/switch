package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Kareem-Emad/switch/auth"
	"github.com/Kareem-Emad/switch/producer"
	"github.com/gorilla/mux"
)

var tag = "[Switch_Server]"

// Start starts an http server with the api routes loaded as specified
func Start(pm producer.ProductionManager) {
	fmt.Println(fmt.Sprintf("%s Starting server ....", tag))
	r := mux.NewRouter()

	r.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		var message Message

		if auth.Authenticate(r.Header) == false {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		pm.GenerateJobsForTopic(message.Topic, message.Payload)
	}).Methods("POST")

	r.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		if auth.Authenticate(r.Header) == false {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		pm.SeedTopicSubcriptionMap()
	}).Methods("GET")

	fmt.Println(fmt.Sprintf("%s  Listening to port %s", tag, serverPort))
	http.ListenAndServe(fmt.Sprintf(":%s", serverPort), r)
}
