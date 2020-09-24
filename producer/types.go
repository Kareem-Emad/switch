package producer

import (
	"github.com/Kareem-Emad/switch/dal"
	faktory "github.com/contribsys/faktory/client"
)

//ProductionManager the stucture holding all data associated with faktory
type ProductionManager struct {
	FaktoryClient     *faktory.Client
	subcriptionGroups map[string][]dal.Subscriber
}
