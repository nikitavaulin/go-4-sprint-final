package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	splitData := strings.Split(data, ",")
	if parsedCount := len(splitData); parsedCount != 2 {
		return 0, 0, fmt.Errorf("wrong parsedCount, got: %d wanted: %d", parsedCount, 2)
	}

	stepsCount, err := strconv.Atoi(splitData[0])
	if err != nil {
		return 0, 0, fmt.Errorf("did not parse stepsCount: %w", err)
	}

	if stepsCount <= 0 {
		return 0, 0, fmt.Errorf("stepsCount should be more than 0, got: %d wanted", stepsCount)
	}

	walkTime, err := time.ParseDuration(splitData[1])
	if err != nil {
		return 0, 0, fmt.Errorf("did not parse walkTime: %w", err)
	}

	return stepsCount, walkTime, nil
}

func DayActionInfo(data string, weight, height float64) string {
	stepsCount, walkTime, err := parsePackage(data)
	if err != nil {
		fmt.Printf("error: %w\n", err) // mb fix
		return ""
	}

	distance := float64(stepsCount) * stepLength / float64(mInKm)

	spentCalories, err := spentcalories.WalkingSpentCalories(
		stepsCount,
		weight,
		height,
		walkTime)
	if err != nil {
		fmt.Printf("error: %w\n", err) // mb fix
		return ""
	}

	output := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила: %.2f\nВы сожгли %.2f ккал.\n",
		stepsCount,
		distance,
		spentCalories)

	return output
}
