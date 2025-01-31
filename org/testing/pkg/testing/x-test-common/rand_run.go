package xtestcommon

import (
	"sync"
	"time"

	"github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
)

const defaultRunRandID = models.RunRandID("t")

var runRandID = defaultRunRandID
var runRandIDLock sync.Mutex

const runRandTimeFormat = "0601021504" // yyMMDDhhmm

func GetOrInitRunRandID() models.RunRandID {
	runRandIDLock.Lock()
	defer runRandIDLock.Unlock()
	if runRandID == defaultRunRandID {
		initializeRunRandID()
	}
	return runRandID
}

func initializeRunRandID() {
	s := models.ScenarioRandID(time.Now().UTC().Format(runRandTimeFormat))
	runRandID = models.RunRandID("t" + s)
}
