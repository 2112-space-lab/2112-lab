package xtestcommon

import (
	"math/rand"

	"github.com/org/2112-space-lab/org/testing/pkg/testing/x-test-common/models"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
const randScenarioLength = 8

func GenerateScenarioRandID() models.ScenarioRandID {
	s := randStringBytesRmndr(randScenarioLength)
	return models.ScenarioRandID(s)
}

func randStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
