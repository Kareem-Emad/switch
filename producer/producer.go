package producer

import (
	"fmt"

	"github.com/Kareem-Emad/switch/dal"
	faktory "github.com/contribsys/faktory/client"
)

var tag = "[Switch_Faktory_Producer]"

// InitalizeFaktoryConnection starts a connection with the faktory server
func (pm *ProductionManager) InitalizeFaktoryConnection() {
	faktoryClient, err := faktory.Open()
	if err != nil {
		panic(fmt.Sprintf("%s Failed to establish connection to faktory server | error: %s", tag, err))
	}
	pm.FaktoryClient = faktoryClient

}

// SeedTopicSubcriptionMap grabs a fresh copy of all subscribers from the Database grouped by to topic
func (pm *ProductionManager) SeedTopicSubcriptionMap() {
	pm.subcriptionGroups = dal.GroupSubscribersByTopic()
}

// PushNewJob adds a new job into the queue
func (pm *ProductionManager) PushNewJob(url string, filterExp string, payload string) bool {
	job := faktory.NewJob(producerQueueNamespace, url, filterExp, payload)
	job.Queue = producerQueue // Kind of a lazy update because the queue param is not respected(always 'default')

	err := pm.FaktoryClient.Push(job)
	if err != nil {
		fmt.Println(fmt.Sprintf("%s Failed to push new job in queue | error: %s", tag, err))
		return false
	}

	fmt.Println(fmt.Sprintf("%s Sucessfully insereted job in queue with specs {%s, %s, %s}", tag, url, filterExp, payload))
	return true
}

// GenerateJobsForTopic generates a job for each subscriber attached to this topic with the provided data
func (pm *ProductionManager) GenerateJobsForTopic(topic string, payload string) bool {
	fmt.Println(fmt.Sprintf("%s getting subcribers for topic %s | found %d subscribers", tag, topic, len(pm.subcriptionGroups[topic])))

	for _, sub := range pm.subcriptionGroups[topic] {
		status := pm.PushNewJob(sub.TargetURL, sub.FilterExpression, payload)
		if !status {
			return false
		}
	}
	fmt.Println(fmt.Sprintf("%s Successfully Triggered %d subcribers for topic %s", tag, len(pm.subcriptionGroups[topic]), topic))

	return true
}
