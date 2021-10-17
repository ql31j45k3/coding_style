package tools

import (
	"net/http"
	"sync"
)

var (
	onceTransport    sync.Once
	defaultTransport *http.Transport
)

func GetDefaultTransport() *http.Transport {
	// TODO: 需複製 http.DefaultTransport 不要共用同一個
	onceTransport.Do(func() {
		defaultTransport = http.DefaultTransport.(*http.Transport)

		// 注意: MaxIdleConnsPerHost 為最大空閑連接數
		defaultTransport.MaxIdleConnsPerHost = 20
	})

	return defaultTransport
}
