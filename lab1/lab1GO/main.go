package main

import (
	"fmt"
	"go-app/tasks"
)

func main() {
	fmt.Println("Task 1")
	tasks.MultiplicationTable(7)

	fmt.Println("Task2")
	test := []int{1, 2, 3, 4, 7, 2}
	tasks.Increment(test)

	fmt.Println("Task 3")
	arr1 := []interface{}{115, 101, 122, 105, 122}
	result1 := tasks.IsVow(arr1)
	fmt.Println("Результат: ", result1)
}
