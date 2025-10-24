package tasks

import "fmt"

func Increment(nums []int) []int {
	result := make([]int, len(nums))
	for i := 0; i < len(nums); i++ {
		sum := nums[i] + (i + 1)
		
		if sum > 9 {
			result[i] = sum % 10
		} else {
			result[i] = sum
		}
		fmt.Print(result[i], " ")
	}
	fmt.Println()
	return result
}
