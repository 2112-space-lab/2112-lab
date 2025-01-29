package handlers

import (
	"fmt"
	"strconv"

	"github.com/Elbujito/2112/src/app-service/internal/config"
)

// TaskName alias definition
type TaskName string

// Task definition
type Task struct {
	Name         TaskName
	Description  string
	RequiredArgs []string
}

// TaskEnv definition
type TaskEnv struct {
	Env *config.SEnv
}

// SetEnv set env variables from config
func (te *TaskEnv) SetEnv(env *config.SEnv) {
	te.Env = env
}

// ParseIntArg helpers to parse int from command line arguments
func ParseIntArg(args map[string]string, key string) (int, error) {
	value, ok := args[key]
	if !ok || value == "" {
		return 0, fmt.Errorf("missing required argument '%s'", key)
	}
	return strconv.Atoi(value)
}
