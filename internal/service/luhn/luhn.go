package luhn

func IsValid(number int) bool {
	if number < 1 {
		return false
	}
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		curr := number % 10

		if i%2 == 0 {
			curr = curr * 2
			if curr > 9 {
				curr = curr%10 + curr/10
			}
		}

		luhn += curr
		number = number / 10
	}
	return luhn % 10
}
