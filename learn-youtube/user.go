package main

import "math/rand"

func randomUser() string {
	names := []string{"Alice", "Bob", "Charlie", "David", "Eve"}
	randomIndex := rand.Intn(len(names))
	return names[randomIndex]
}
