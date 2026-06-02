package generator

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func RandomNumber(min, max int) int {

	// Generate a random number in the specified range
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		fmt.Println("Error generating random number:", err)
		return 500
	}

	// Add the minimum value to the random number to bring it into the desired range
	randomNumber = randomNumber.Add(randomNumber, big.NewInt(int64(min)))

	// Convert the result to an integer
	randomInt := int(randomNumber.Int64())

	return randomInt

}
