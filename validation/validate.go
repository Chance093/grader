package validation

import (
	"errors"
	"log"
	"strconv"
)

func ValidateWeights(weights []string) error {
	var sum int
	for _, weight := range weights {
		weightInt, err := strconv.Atoi(weight)
		if err != nil {
			log.Fatalf("String to int conversion failed %v\n", err)
		}

		sum += weightInt
	}

	if sum != 100 {
		return errors.New("Test, Quiz, and Homework weights should add to 100. Please try again.")
	}

	return nil
}

func ValidatePoints(totalPoints, correctPoints string) error {
	totalPointsInt, err := strconv.Atoi(totalPoints)
	if err != nil {
		log.Fatalf("Failed to convert string to int: %s", err.Error())
	}

	correctPointsInt, err := strconv.Atoi(correctPoints)
	if err != nil {
		log.Fatalf("Failed to convert string to int: %s", err.Error())
	}

	if totalPointsInt < correctPointsInt {
		return errors.New("Total points must be greater than or equal to correct points")
	}

	return nil
}
