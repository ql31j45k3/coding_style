package example

import (
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

type APIExampleCond struct {
	dig.In

	R *gin.Engine
}

type taskMap struct {
	sync.RWMutex

	taskID2Detail map[string]taskMapDetail
}

type taskMapDetail struct {
	goCount     int
	finishCount int
}

func (r *taskMap) getAndExist(taskID string) (taskMapDetail, bool) {
	r.RLock()
	defer r.RUnlock()

	v, ok := r.taskID2Detail[taskID]
	return v, ok
}

func (r *taskMap) isExist(taskID string) bool {
	r.RLock()
	defer r.RUnlock()

	_, ok := r.taskID2Detail[taskID]
	return ok
}

func (r *taskMap) set(taskID string, v taskMapDetail) {
	r.Lock()
	defer r.Unlock()

	r.taskID2Detail[taskID] = v
}

type taskData struct {
	taskID string

	task *taskMap
}

type responseGoroutine struct {
	TaskID string `json:"task_id"`
}

type responseTaskStatus struct {
	TaskID string `json:"task_id"`

	Detail responseTaskStatusDetail `json:"detail"`
}

type responseTaskStatusDetail struct {
	GoCount     int `json:"go_count"`
	FinishCount int `json:"finish_count"`
}
