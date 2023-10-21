package c2api

// import (
// 	_ "embed"
// )

// //go:embed routes.txt
// var RawRoutesDoc string

// const R_FETCH_MESSAGES = 0
// const R_MESSAGE_RESPOND = 1

// var RouteMap map[int][]string = func() (rMap map[int][]string) {
// 	registeredObfuscatedRoutes := []int{
// 		R_FETCH_MESSAGES, R_MESSAGE_RESPOND,
// 	}

// 	rMap = make(map[int][]string, len(registeredObfuscatedRoutes))

// 	wlChunks := chunkSlice(strings.Split(RawRoutesDoc, "\n"), len(registeredObfuscatedRoutes))

// 	for chunkIndex, chunk := range wlChunks {
// 		rMap[registeredObfuscatedRoutes[chunkIndex]] = chunk
// 	}

// 	return rMap
// }()
