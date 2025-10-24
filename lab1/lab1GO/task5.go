package tasks

func IsVow(arr []interface{}) []interface{} {

vowels := map[int]string {
	97: "a",
	101: "e",
	105: "i",
	111: "o",
	117: "u",
}

for i ,val := range arr {
	if num, ok := val.(int); ok {
		if vowel, exists := vowels[num]; exists {
			arr[i] = vowel
		}
	}
}

	return arr
}