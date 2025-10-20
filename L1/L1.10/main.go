package main

import (
	"fmt"
	"math/rand"
	"sort"
)

func main() {
	//sequence := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
	sequence := generateSequence(5, 100, -50.00, 50.00)
	groups := groupValues(sequence)

	for key, val := range groups {
		fmt.Println(key, val)
	}
}

func generateSequence(minElements int, maxElements int, minValue float64, maxValue float64) (result []float64) {
	for i := 0; i <= minElements+rand.Intn(maxElements-minElements+1); i++ {
		result = append(result, (minValue + rand.Float64()*(maxValue-minValue)))
	}
	return result
}

func groupValues(sequence []float64) map[int][]float64 {
	groups := make(map[int][]float64)

	sort.Float64s(sequence)

	for _, val := range sequence {
		key := int(val/10) * 10
		groups[int(key)] = append(groups[int(key)], val)
	}

	return groups
}
