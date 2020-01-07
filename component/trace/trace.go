package trace

import (
	"github.com/ebar-go/ego/utils/strings"
	"github.com/petermattis/goid"
	"sync"
)

var (
	traceIds = map[int64]string{}
	rwm      sync.RWMutex
)

func NewId() string {
	return "TraceId" + strings.UUID()
}

// SetTraceId
func SetTraceId(id string) {
	goID := getGoroutineId()
	rwm.Lock()
	defer rwm.Unlock()

	traceIds[goID] = id
}

// GetTraceId
func GetTraceId() string {
	goID := getGoroutineId()
	rwm.RLock()
	defer rwm.RUnlock()

	return traceIds[goID]
}

// DeleteTraceId
func DeleteTraceId() {
	goID := getGoroutineId()
	rwm.Lock()
	defer rwm.Unlock()

	delete(traceIds, goID)
}

func getGoroutineId() int64 {
	return goid.Get()
}
