package tasks

import "fmt"

func MultiplicationTable(num int) {
	for i := 1; i <= num; i++ {
		for j := 1; j <= num; j++ {
			fmt.Print(i*j, "\t")
		}
		fmt.Println()
	}
}