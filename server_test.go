package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Kareem-Emad/switch/auth"
	"github.com/Kareem-Emad/switch/dal"
	"github.com/dgrijalva/jwt-go"
)

// ============================================================
//                        Setup
// ============================================================
func cleanDB() {
	dal.DeleteAllSubscribers()
}

// =============================================================
//                         Tests
// =============================================================
func TestAuthenticationFlow(t *testing.T) {
	/*
		should pass authentication as long as headers are
		set properly with a valid token
	*/
	jwtKey := "my_secret_key"

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &auth.Claims{
		Author: "some_service",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))

	if err != nil {
		t.Fatal(err)
	}

	testToken := fmt.Sprintf("bearer %s", tokenString)

	fakeHeaders := http.Header{}
	fakeHeaders.Add("Authorization", testToken)

	if auth.Authenticate(fakeHeaders) != true {
		t.Fatal("token should pass authentication signature check")
	}
}

func TestGetSubs(t *testing.T) {
	/*
		should be able to fetch all subscribers from db using FetchAllSubscribers
	*/
	subs := dal.FetchAllSubscribers()

	if len(subs) != 0 {
		t.Fatal("subs list should be empty as no subscription created yet")
	}

}

func TestCreateandGetSubs(t *testing.T) {
	/*
		should be able to fetch  all created subscribers from db
	*/
	defer cleanDB()

	dal.CreateSubscriber("True", "https://google.com", "mytopic")
	subs := dal.FetchAllSubscribers()

	if len(subs) != 1 || subs[0].FilterExpression != "True" || subs[0].TargetURL != "https://google.com" || subs[0].Topic != "mytopic" {
		t.Fatal("subs list should be contain 1 created record")
	}

	dal.CreateSubscriber("False", "https://facebook.com", "mytopic")
	subs = dal.FetchAllSubscribers()

	if len(subs) != 2 || subs[1].FilterExpression != "False" || subs[1].TargetURL != "https://facebook.com" || subs[1].Topic != "mytopic" {
		t.Fatal("subs list should contain 2 created records")
	}

}

func TestGroupSubscribers(t *testing.T) {
	/*
		should be able to fetch all subscribers grouped by topic in a map
	*/
	defer cleanDB()
	dal.CreateSubscriber("True", "https://google11.com", "t1")
	dal.CreateSubscriber("True", "https://google12.com", "t1")
	dal.CreateSubscriber("True", "https://google13.com", "t1")
	dal.CreateSubscriber("True", "https://google21.com", "t2")
	dal.CreateSubscriber("True", "https://google22.com", "t2")

	subsMap := dal.GroupSubscribersByTopic()

	if len(subsMap["t1"]) != 3 || len(subsMap["t2"]) != 2 {
		t.Fatal("should be two groups of topics in map with two subs arrays 3,2")
	}
}
