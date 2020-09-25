package dal // data access layer

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var tag = "[Switch_DB_Driver]"

// opens a db connection with the specified backend(sqlite, mysql, postgresql)
func openConnection() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("%s  failed to connect database | error: %s", tag, err))
	}
	return db
}

// FetchAllSubscribers lists all subscribers in the subcribers table
func FetchAllSubscribers() []Subscriber {
	var subs []Subscriber

	db := openConnection()

	db.AutoMigrate(&Subscriber{})
	db.Find(&subs)
	fmt.Println(fmt.Sprintf("%s fetched all subscribers | total count %d", tag, len(subs)))
	return subs
}

// CreateSubscriber creates a new subscriber and returns the status of the creation as bool
func CreateSubscriber(filterExp string, url string, topic string) bool {
	db := openConnection()

	result := db.Create(&Subscriber{FilterExpression: filterExp, TargetURL: url, Topic: topic})
	return result.Error == nil
}

// GroupSubscribersByTopic returns a map of topics each key holding an array of associated subscribers
func GroupSubscribersByTopic() map[string][]Subscriber {
	subscriberTopicGroups := make(map[string][]Subscriber)

	subs := FetchAllSubscribers()

	for _, sub := range subs {
		subscriberTopicGroups[sub.Topic] = append(subscriberTopicGroups[sub.Topic], sub)
	}
	return subscriberTopicGroups
}

// DeleteAllSubscribers deletes all subscribers' records from database
func DeleteAllSubscribers() bool {
	db := openConnection()

	result := db.Where("1 = 1").Delete(&Subscriber{})
	return result.Error == nil
}
