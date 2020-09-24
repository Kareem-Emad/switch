package dal

import "gorm.io/gorm"

// Subscriber the struct holding info about each topic subscriber
type Subscriber struct {
	gorm.Model
	FilterExpression string
	TargetURL        string
	Topic            string
}
