package api

import (
	"fmt"
	"sort"
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
		if formatRepeat[1] == "" || len(formatRepeat[1]) == 0 {
			return fmt.Errorf("некоректные параметры для правила")
		}

		days := strings.Split(formatRepeat[1], ",")

		for _, day := range days {
			dayNum, err := strconv.Atoi(day)
			if err != nil {
				return fmt.Errorf("ошибка форматирования strconv.Atoi(formatRepeat[1]): %w", err)
			}
			if dayNum < 1 || dayNum > 7 {
				return fmt.Errorf("некорректный день недели: %v", dayNum)
			}
		}

		for _, dayStr := range days {
			dayNum, err := strconv.Atoi(dayStr)
			if err != nil {
				return fmt.Errorf("ошибка конвертации: %w", err)
			}
			if dayNum > 7 {
				return fmt.Errorf("первышено максимально допустимый интрервал для %v > 7", dayNum)
			}
		}
	case "m":
		if len(formatRepeat) == 3 {
			monthsSlice := strings.Split(formatRepeat[2], ",")
			for _, v := range monthsSlice {
				numMonth, err := strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("ошибка форматирования strconv.Atoi(formatRepeat[1]): %w", err)
				}
				if numMonth > 12 {
					return fmt.Errorf("недоступный месяц под номером: %v", numMonth)
				}
			}
		}
		days := strings.Split(formatRepeat[1], ",")
		var daysInt []int
		for _, v := range days {

			vInt, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("ошибка форматирования strconv.Atoi(v): %w", err)
			}

			daysInt = append(daysInt, int(vInt))
		}
		sort.Ints(daysInt)
		for _, v := range daysInt {
			if v > 31 {
				return fmt.Errorf("номер дня в инструкции > 31: %v, %v", v, daysInt)
			}
		}

	default:
		if formatRepeat[0] != "y" {
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
	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("время в переменной dstart не может быть преобразовано в корректную дату — ошибка выполнения time.Parse('20060102', dstart): %w", err)
	}
	formatRepeat := strings.Split(repeat, " ")

	// Обработка всех случаев
	switch formatRepeat[0] {
	case "d":
		days, err := strconv.Atoi(formatRepeat[1])
		if err != nil {
			return "", fmt.Errorf("ошибка форматирования strconv.Atoi(formatRepeat[1]): %w", err)
		}

		maxIterations := 100
		for i := 0; i < maxIterations; i++ {
			// if !now.Before(date) {
			// 	return date.Format(dateFormat), nil
			// }
			date = date.AddDate(0, 0, days)
			if afterNow(date, now) {
				return date.Format(dateFormat), nil
			}
		}
		Logger.Printf("не удалось найти дату в пределах %d итераций", maxIterations)
		return "", fmt.Errorf("не удалось найти дату в пределах %d итераций", maxIterations)

	case "w":
		var daysWeek [7]bool
		days := strings.Split(formatRepeat[1], ",")
		for _, v := range days {
			daysInt, err := strconv.Atoi(v)
			if err != nil {
				return "", err
			}
			daysWeek[daysInt-1] = true
		}

		for {
			date = date.AddDate(0, 0, 1)
			weekday := date.Weekday()
			v := GetDayNumberByString(weekday.String())

			if daysWeek[v] {
				if afterNow(date, now) {
					return date.Format(dateFormat), nil
				}
			}
		}

	case "m":
		var day [32]bool
		var months [13]bool
		// заполняем доступные месяцы
		if len(formatRepeat) == 3 {
			monthsSlice := strings.Split(formatRepeat[2], ",")
			if monthsSlice[0] == "" {
				for i, _ := range months {
					months[i] = true
				}
			}
			for _, v := range monthsSlice {
				numMonth, err := strconv.Atoi(v)
				if err != nil {
					return "", fmt.Errorf("ошибка форматирования strconv.Atoi(formatRepeat[1]): %w", err)
				}
				months[numMonth-1] = true
			}
		} else {
			for i, _ := range months {
				months[i] = true
			}
		}
		var days []int
		for _, v := range strings.Split(formatRepeat[1], ",") {
			dayNum, err := strconv.Atoi(v)
			if err != nil {
				return "", fmt.Errorf("ошибка форматирования strconv.Atoi(formatRepeat[1]): %w", err)
			}

			days = append(days, dayNum)
		}
		// Запись дней в слайс day
		for _, dayNum := range days {
			// Запись значений > 0
			if dayNum > 0 {
				day[dayNum-1] = true
			}
			// Запись значений < 0
			if dayNum < 0 {
				if dayNum < -2 {
					return "", fmt.Errorf("запрос на -3 день недоступен")
				}
				y, m, _ := date.Date()
				// logger.Printf("m=%d", m)

				firstOfNextMonth := time.Date(y, m+1, 1, 0, 0, 0, 0, time.UTC)
				// logger.Printf("первый день следующего месяца %v", firstOfNextMonth.Format(dateFormat))

				// dayMinus хранит номер дня высчитаный из инструкции -1,-2
				dayMinus := firstOfNextMonth.AddDate(0, 0, dayNum)
				// logger.Printf("дата после вычитания dayMinus[%v]", dayMinus.Format(dateFormat))
				// dayCount номер последнего или предпоследнего дня
				_, _, dayCount := dayMinus.Date()

				day[dayCount-1] = true
			}
		}
		for i := 0; i < 700; i++ {
			date = date.AddDate(0, 0, 1)
			_, m, d := date.Date()
			if months[m-1] {
				if day[d-1] {
					if afterNow(date, now) {
						return date.Format(dateFormat), nil
					}
				}
			}
			if i >= 700 {
				return "", fmt.Errorf("превышено число итераций цикла i > 700, %v\n%v", months, day)
			}
		}

	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(dateFormat), nil
	}
	return "", nil
}

// afterNow проверяет что date > now
func afterNow(date, now time.Time) bool {
	return date.After(now)
}

// GetDayNumberByString компенсирует разницу начала отсчета недели
func GetDayNumberByString(day string) int {
	dayMap := map[string]int{
		"Sunday":    6,
		"Monday":    0,
		"Tuesday":   1,
		"Wednesday": 2,
		"Thursday":  3,
		"Friday":    4,
		"Saturday":  5,
	}
	if number, ok := dayMap[day]; ok {
		return number
	}
	return 0
}
