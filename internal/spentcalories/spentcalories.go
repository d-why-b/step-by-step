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
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// Переменные ошибок
var (
	ErrStringFormat = errors.New("неверный формат строки")
	ErrSteps        = errors.New("неверное количество шагов")
	ErrTime         = errors.New("неверное время")
	ErrWeight       = errors.New("неверный вес")
	ErrHeigt        = errors.New("неверный рост")
	ErrSpeed        = errors.New("неверная скорость")
	ErrTraningType  = errors.New("неизвестный тип тренировки")
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parseData := strings.Split(data, ",")
	if len(parseData) != 3 {
		return 0, "", 0, ErrStringFormat // Возвращается пустая строка, вместо типа тренировки, для сигнализации о том, что произошла ошибка.
	}

	steps, err := strconv.Atoi(parseData[0])
	if err != nil {
		return 0, "", 0, err // Возвращается пустая строка, вместо типа тренировки, для сигнализации о том, что произошла ошибка.
	}

	if steps <= 0 {
		return 0, "", 0, ErrSteps
	}

	time, err := time.ParseDuration(parseData[2])
	if err != nil {
		return 0, "", 0, err // Возвращается пустая строка, вместо типа тренировки, для сигнализации о том, что произошла ошибка.
	}

	if time <= 0 {
		return 0, "", 0, ErrTime
	}

	return steps, parseData[1], time, nil
}

func distance(steps int, height float64) float64 {
	strideLength := height * stepLengthCoefficient
	distance := float64(steps) * strideLength / mInKm
	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	distance := distance(steps, height)

	if duration <= 0 {
		return 0
	}

	return distance / duration.Hours()
}

// TrainingInfo принимает информацию о тренировки и парсит в строку.
// Функция принимает:
// количество шагов, тип тренировки, время тренировки в виде строки формата "1234,Ходьба,2h30m"; вес; рост
// Возвращает строку в формате:
// "Тип тренировки: Бег
// Длительность: 0.75 ч.
// Дистанция: 10.00 км.
// Скорость: 13.34 км/ч
// Сожгли калорий: 18621.75" и ошибку.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, traning, time, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if steps <= 0 {
		log.Println(ErrSteps.Error())
		return "", ErrSteps
	}

	if time <= 0 {
		log.Println(ErrTime.Error())
		return "", ErrTime
	}

	if weight <= 0 {
		log.Println(ErrWeight.Error())
		return "", ErrWeight
	}

	if height <= 0 {
		log.Println(ErrHeigt.Error())
		return "", ErrHeigt
	}

	switch traning {
	case "Бег":
		distance := distance(steps, height)
		averageSpeed := meanSpeed(steps, height, time)
		caloriesBurned, err := RunningSpentCalories(steps, weight, height, time)

		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", traning, time.Hours(), distance, averageSpeed, caloriesBurned), nil

	case "Ходьба":
		distance := distance(steps, height)
		averageSpeed := meanSpeed(steps, height, time)
		caloriesBurned, err := WalkingSpentCalories(steps, weight, height, time)

		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", traning, time.Hours(), distance, averageSpeed, caloriesBurned), nil
	default:
		return "", ErrTraningType
	}
}

func calculateSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, ErrSteps
	}

	if duration <= 0 {
		return 0, ErrTime
	}

	if weight <= 0 {
		return 0, ErrWeight
	}

	if height <= 0 {
		return 0, ErrHeigt
	}

	averageSpeed := meanSpeed(steps, height, duration)
	if averageSpeed <= 0 {
		fmt.Println(ErrSpeed)
	}

	caloriesBurned := weight * averageSpeed * duration.Minutes() / minInH
	return caloriesBurned, nil
}

// RunningSpentCalories считает и возвращает количество пораченых каллорий после пробежки.
// На вход ожидается:
// количество шагов; вес; рост; время пробежки.
// Фукция возвращает количество потраченных калорий и ошибку.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := calculateSpentCalories(steps, weight, height, duration)
	return calories, err
}

// WalkingSpentCalories считает и возвращает количество потрачеых каллорий после прогулки.
// На вход ожидается:
// количество шагов; вес; рост; время прогулки.
// Функция возвращает количество пораченых калорий и ошибку.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := calculateSpentCalories(steps, weight, height, duration)
	return calories * walkingCaloriesCoefficient, err
}
