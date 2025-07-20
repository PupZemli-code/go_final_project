package repeat

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ValidDstarRepeat производит проверки входящих данных dstart, repeat
func ValidDstarRepeat(dstart string, repeat string) error {
	if repeat == "" {
		return fmt.Errorf("в параметре repeat — пустая строка")
	}
	formatRepeat := strings.Split(repeat, " ")

	switch formatRepeat[0] {
	case "d":
		if len(formatRepeat) <= 1 {
			return fmt.Errorf("не указан интервал в днях")
		}
		interval, err := strconv.Atoi(formatRepeat[1])
		if err != nil {
			return fmt.Errorf("ошибка форматирования strconv.Atoi(formatRepeat[1]): %w", err)
		}
		if interval > 400 {
			return fmt.Errorf("первышено максимально допустимый интрервал для %v = %v", formatRepeat, interval)
		}
	case "w":
		if len(formatRepeat) <= 1 {
			return fmt.Errorf("не указаны дни повторений")
		}
		interval, err := strconv.Atoi(formatRepeat[1])
		if err != nil {
			return err
		}
		if interval > 7 {
			return fmt.Errorf("первышено максимально допустимый интрервал для %v = %v", formatRepeat, interval)
		}
	default:
		if formatRepeat[0] != "m" && formatRepeat[0] != "y" {
			return fmt.Errorf("недопустимый символ формата: %v", formatRepeat[0])
		}
	}
	return nil
}

// NextDate возвращает строку с датой в формате 20060102
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if err := ValidDstarRepeat(dstart, repeat); err != nil {
		return "", fmt.Errorf("формат repeat не прошел проверку: %w", err)
	}
	// В разработке
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("время в переменной dstart не может быть преобразовано в корректную дату — ошибка выполнения time.Parse('20060102', dstart): %w", err)
	}
	formatRepeat := strings.Split(repeat, " ")
	var days int

	switch formatRepeat[0] {
	case "d":
		days, err = strconv.Atoi(formatRepeat[1])
		if err != nil {
			return "", fmt.Errorf("ошибка форматирования strconv.Atoi(formatRepeat[1]): %w", err)
		}

		for {
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format("20060102"), nil

	case "w":
		return "", fmt.Errorf("формат еще не подерживается")
	case "m":
		return "", fmt.Errorf("формат еще не подерживается")
	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format("20060102"), nil
	default:
		return "", nil
	}
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}
