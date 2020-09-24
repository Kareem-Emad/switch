package producer

import "os"

var producerQueue = os.Getenv("SWITCH_PRODUCTION_QUEUE")
var producerQueueNamespace = os.Getenv("SWITCH_PRODUCTION_QUEUE_NAMESPACE")
