package main

import (
	"fmt"
)

func main() {
	var number int64 = 5
	fmt.Println(number)
	fmt.Printf("%064b\n", number)

	result, err := Int64BitToggle(number, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	fmt.Printf("%064b\n", result)

	fmt.Println()

	var number2 int64 = 4
	fmt.Println(number2)
	fmt.Printf("%064b\n", number2)

	result2, err := Int64BitSet(number2, 0, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result2)
	fmt.Printf("%064b\n", result2)

}

func Int64BitToggle(num int64, index int) (result int64, err error) {
	if index >= 64 || index < 0 {
		return num, fmt.Errorf("Given index %d is outside of [0, 63] bounds of Int64", index)
	}
	return num ^ (1 << index), nil
}

func Int64BitSet(num int64, index int, bitValue bool) (result int64, err error) {
	if index >= 64 || index < 0 {
		return num, fmt.Errorf("Given index %d is outside of [0, 63] bounds of Int64", index)
	}
	if bitValue {
		return num | (1 << index), nil // ИЛИ
	} else {
		return num &^ (1 << index), nil // И НЕ
	}
}
