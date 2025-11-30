package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

var (
	ErrValueLessZero       error = errors.New("should be more than 0")
	ErrUnknownTrainingType error = errors.New("неизвестный тип тренировки")
)

func parseTraining(data string) (int, string, time.Duration, error) {
	splitData := strings.Split(data, ",")
	if parsedCount := len(splitData); parsedCount != 3 {
		return 0, "", 0, fmt.Errorf("wrong parsedCount, got: %d wanted: %d", parsedCount, 3)
	}

	stepsCount, err := strconv.Atoi(splitData[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("did not parse stepsCount: %w", err)
	}

	if stepsCount <= 0 {
		return 0, "", 0, fmt.Errorf("stepsCount got: %d, %w", stepsCount, ErrValueLessZero)
	}

	activityTime, err := time.ParseDuration(splitData[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("did not parse activityTime: %w", err)
	}
	if activityTime <= 0 {
		return 0, "", 0, fmt.Errorf("wrong value of activityTime: %w", ErrValueLessZero)
	}

	return stepsCount, splitData[1], activityTime, nil // splitData[1] = activity name

}

func distance(steps int, height float64) float64 {
	var distanceKm float64
	stepLength := height * stepLengthCoefficient
	distanceKm = stepLength * float64(steps) / mInKm
	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0 // mb log
	}
	distanceKm := distance(steps, height)
	return distanceKm / duration.Hours() // average speed

}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, activityTime, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("TrainingInfo: %w", err)
	}

	distanceKm := distance(steps, height)
	avgSpeed := meanSpeed(steps, height, activityTime)

	var spentCalories float64

	switch activity {
	case "Бег":
		spentCalories, err = RunningSpentCalories(steps, weight, height, activityTime)
	case "Ходьба":
		spentCalories, err = WalkingSpentCalories(steps, weight, height, activityTime)

	default:
		return "", fmt.Errorf("TrainingInfo: %w", ErrUnknownTrainingType)
	}
	if err != nil {
		return "", fmt.Errorf("TrainingInfo: %w", err)
	}

	return trainingInfoText(activity, activityTime.Hours(), distanceKm, avgSpeed, spentCalories), nil

}

func trainingInfoText(activity string, activityTimeHours, distanceKm, avgSpeed, spentCalories float64) string {
	output := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, activityTimeHours, distanceKm, avgSpeed, spentCalories)
	return output
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateParamsSpentCalories(steps, weight, height, duration); err != nil {
		return 0, fmt.Errorf("RunningSpentCalories: %w", err)
	}
	avgSpeed := meanSpeed(steps, height, duration)
	spentCalories := weight * avgSpeed * duration.Minutes() / float64(minInH) // (weight * meanSpeed * durationInMinutes) / minInH
	return spentCalories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateParamsSpentCalories(steps, weight, height, duration); err != nil {
		return 0, fmt.Errorf("WalkingSpentCalories: %w", err)
	}
	avgSpeed := meanSpeed(steps, height, duration)
	spentCalories := weight * avgSpeed * duration.Minutes() / float64(minInH) // (weight * meanSpeed * durationInMinutes) / minInH
	return spentCalories * walkingCaloriesCoefficient, nil
}

func validateParamsSpentCalories(steps int, weight, height float64, duration time.Duration) error {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return fmt.Errorf("some values of arguments are wrong, %w", ErrValueLessZero)
	}
	return nil
}
