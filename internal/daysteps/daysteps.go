// daysteps отвечает за учёт активности в течение дня.
// Он собирает переданную информацию в виде строк,
// парсит их и выводит информацию о количестве шагов, пройденной дистанции и потраченных калориях.
package daysteps

import (
	"fmt"
	"log"
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

// parsePackage принимает строку с данными,
// возращает количество шагов и время прогулки
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
		return 0, 0, fmt.Errorf("stepsCount should be more than 0, got: %d", stepsCount)
	}

	walkTime, err := time.ParseDuration(splitData[1])
	if err != nil {
		return 0, 0, fmt.Errorf("did not parse walkTime: %w", err)
	}

	if walkTime <= 0 {
		return 0, 0, fmt.Errorf("walkTime should be more than 0, got: %d", walkTime)
	}

	return stepsCount, walkTime, nil
}

// DayActionInfo парсит строку с данными,
// вычисляет дистанцию в километрах и количество потраченных калорий
// и возвращает строку с инофрмацией
func DayActionInfo(data string, weight, height float64) string {
	stepsCount, walkTime, err := parsePackage(data)
	if err != nil {
		log.Printf("error: %s\n", err)
		return ""
	}

	distance := float64(stepsCount) * stepLength / float64(mInKm)

	spentCalories, err := spentcalories.WalkingSpentCalories(
		stepsCount,
		weight,
		height,
		walkTime)
	if err != nil {
		log.Printf("error: %s\n", err)
		return ""
	}

	return dayActionInfoText(stepsCount, distance, spentCalories)
}

// dayActionInfoText возвращает текст вывода информации для пользователя
func dayActionInfoText(stepsCount int, distance, spentCalories float64) string {
	output := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		stepsCount,
		distance,
		spentCalories)
	return output
}
