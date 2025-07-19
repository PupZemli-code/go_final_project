package repeat

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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
			return err
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
			fmt.Errorf("недопустимый символ формата: %v", formatRepeat[0])
		}
	}
	return nil
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if err := ValidDstarRepeat(dstart, repeat); err != nil {
		return "", fmt.Errorf("формат repeat не прошел проверку: %w", err)
	}
	/* // В разработке
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("время в переменной dstart не может быть преобразовано в корректную дату — ошибка выполнения time.Parse('20060102', dstart): %w", err)
	}
	formatRepeat := strings.Split(repeat, " ")
	var interval int

	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}
	*/
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}
