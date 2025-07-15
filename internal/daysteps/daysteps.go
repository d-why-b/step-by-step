package daysteps

import (
	"errors"
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

// Переменные ошибок
var (
	ErrStringFormat = errors.New("неверный формат строки")
	ErrSteps        = errors.New("неверное количество шагов")
	ErrTime         = errors.New("неверное время")
)

func parsePackage(data string) (int, time.Duration, error) {
	parseData := strings.Split(data, ",")
	if len(parseData) != 2 {
		return 0, 0, ErrStringFormat
	}

	steps, err := strconv.Atoi(parseData[0])
	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {
		return 0, 0, ErrSteps
	}

	time, err := time.ParseDuration(parseData[1])
	if err != nil {
		return 0, 0, err
	}

	if time <= 0 {
		return 0, 0, ErrTime
	}

	return steps, time, nil
}

// DayActionInfo парсит данные о тренировке, весе и росте, преобразует их в строку и возвращает её.
func DayActionInfo(data string, weight, height float64) string {
	steps, time, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if steps <= 0 {
		log.Println("")
		return ""
	}
	if time <= 0 {
		log.Println("")
		return ""
	}

	distance := float64(steps) * stepLength / mInKm

	caloriesBurned, err := spentcalories.WalkingSpentCalories(steps, weight, height, time)
	if err != nil {
		log.Println(err)
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, caloriesBurned)
}
