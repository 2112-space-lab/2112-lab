package integration

import (
	"testing"
)

func TestAppDatabase(t *testing.T) {
	performAppTestSuite(t, "database")
}
