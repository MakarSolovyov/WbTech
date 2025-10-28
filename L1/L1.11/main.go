package main

import (
	"fmt"
	"math/rand"
)

func main() {
	//sequence1 := []int{1, 2, 3}
	//sequence2 := []int{2, 3, 4}

	sequence1 := generateSequence(5, 10, 1, 10)
	sequence2 := generateSequence(5, 10, 1, 10)

	intersection := intersect(sequence1, sequence2)

	fmt.Println(sequence1)
	fmt.Println(sequence2)
	fmt.Println(intersection)

}

func generateSequence(minElements int, maxElements int, minValue int, maxValue int) (result []int) {
	for i := 0; i <= minElements+rand.Intn(maxElements-minElements+1); i++ {
		result = append(result, (minValue + rand.Intn(maxValue-minValue+1)))
	}
	return result
}

func intersect(sequence1, sequence2 []int) (result []int) {
	dict := make(map[int]bool) //[int]bool - для исключения дублирования key - значение, val - маркер включения в результат

	for _, val := range sequence1 {
		dict[val] = true
	}

	for _, val := range sequence2 {
		if dict[val] {
			result = append(result, val)
			dict[val] = false
		}
	}
	return result
}
