package c2api

import "math/rand"

func chunkSlice[T any](input []T, numChunks int) [][]T {
	chunks := make([][]T, numChunks)
	for i := range input {
		chunkIndex := i % numChunks
		chunks[chunkIndex] = append(chunks[chunkIndex], input[i])
	}
	return chunks
}

func randElem[T any](slice *[]T) T {
	return (*slice)[rand.Intn(len(*slice))]
}
