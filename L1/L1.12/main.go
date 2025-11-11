package main

import (
	"fmt"
)

func main() {
	sequence := []string{"cat", "cat", "dog", "cat", "tree"}

	set := make(map[string]struct{})

	for _, item := range sequence {
		set[item] = struct{}{}
	}

	fmt.Println("Словарь:", set)

	var result []string
	for key := range set {
		result = append(result, key)
	}

	fmt.Println("Уникальные элементы:", result)
}
