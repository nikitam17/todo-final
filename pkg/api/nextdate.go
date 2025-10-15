package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"todo-final/pkg/db"
)

const FormatDate = "20060102"

// Функция для сравнения дат (date > now)
func afterNow(date, now time.Time) bool {
	return date.After(now)
}

// Функция, вычисляющая следующую дату для задачи в соответствии с указанным правилом
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(FormatDate, dstart)
	if err != nil {
		return "", err
	}
	if repeat == "" {
		//if afterNow(date, now) {
		if (dstart == now.Format(FormatDate)) && (dstart < time.Now().Format(FormatDate)) {
			return "", fmt.Errorf("дата: %s равна текущей и не указано повторение", dstart)
		}
		return dstart, nil
	}
	rep := strings.Split(repeat, " ")
	switch rep[0] {
	// перенос на указанное число дней
	case "d":
		if len(rep) < 2 {
			return "", fmt.Errorf("не указан интервал в днях")
		}
		interval, err := strconv.Atoi(rep[1])
		if err != nil {
			return "", err
		}
		if interval <= 0 || interval > 400 {
			return "", fmt.Errorf("превышен максимально допустимый интервал (1-400): %d", interval)

		}
		if dstart == now.Format(FormatDate) {
			return dstart, nil
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
	// перенос на год
	case "y":
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
	// назначение на указанные дни недели
	case "w":
		if len(rep) < 2 {
			return "", fmt.Errorf("не указаны дни недели")
		}
		days := strings.Split(rep[1], ",")
		weekdays := make([]int, 0)
		for i := 0; i < len(days); i++ {
			d, err := strconv.Atoi(days[i])
			if err != nil {
				return "", err
			}
			if d > 7 || d < 1 {
				return "", fmt.Errorf("дни недели должны быть от 1 для понедельника до 7 для воскресенья: %d", i)

			}
			weekdays = append(weekdays, d)
		}
		for {
			findwd := false
			date = date.AddDate(0, 0, 1)
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			for i := 0; i < len(weekdays); i++ {
				if weekday == weekdays[i] {
					findwd = true
					break
				}
			}
			if !findwd {
				continue
			}
			if afterNow(date, now) {
				break
			}
		}
	// назначение на указанные месяцы (если есть) и дни месяца
	case "m":
		if len(rep) < 2 {
			return "", fmt.Errorf("не указаны дни месяца")
		}
		days := strings.Split(rep[1], ",")
		monthDays := make([]int, 0)
		monthNumbers := make([]int, 0)
		for i := 0; i < len(days); i++ {
			d, err := strconv.Atoi(days[i])
			if err != nil {
				return "", err
			}
			if d > 31 || d < -2 || d == 0 {
				return "", fmt.Errorf("дни месяца должны быть в интервале от 1 до 31 и -1, -2: %d", i)

			}
			monthDays = append(monthDays, d)
		}
		// если учитываем месяцы
		if len(rep) == 3 {
			months := strings.Split(rep[2], ",")
			for i := 0; i < len(months); i++ {
				m, err := strconv.Atoi(months[i])
				if err != nil {
					return "", err
				}
				if m > 12 || m < 1 {
					return "", fmt.Errorf("номера месяцев должны быть в интервале от 1 до 12: %d", i)

				}
				monthNumbers = append(monthNumbers, m)
			}
		}
		for {
			findMonth := false
			findDay := false
			date = date.AddDate(0, 0, 1)
			day := int(date.Day())
			month := int(date.Month())
			for i := 0; i < len(monthNumbers); i++ {
				if month == monthNumbers[i] {
					findMonth = true
					break
				}
			}
			// если не учитываем месяцы, то идём дальше на дни
			if len(monthNumbers) > 0 && !findMonth {
				continue
			}
			for i := 0; i < len(monthDays); i++ {
				md := monthDays[i]
				if md < 0 {
					firstDayOfNextMonth := date.AddDate(0, 1, -date.Day()+1)
					md = firstDayOfNextMonth.AddDate(0, 0, md).Day()
				}
				if day == md {
					findDay = true
					break
				}
			}
			if !findDay {
				continue
			}
			if afterNow(date, now) {
				break
			}
		}
	default:
		return "", fmt.Errorf("недопустимый символ: %s", rep[0])
	}
	return date.Format(FormatDate), nil
}

// обработичик для функции обработки правил переноса задач
func nextDayHandler(res http.ResponseWriter, req *http.Request) {
	dstart := req.FormValue("date")
	repeat := req.FormValue("repeat")
	strNow := req.FormValue("now")
	var strDate string
	var err error
	var now time.Time
	if strNow == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(FormatDate, strNow)
	}
	if err == nil {
		strDate, err = NextDate(now, dstart, repeat)
	}
	if err == nil {
		res.Write([]byte(strDate))
	} else {
		errString := err.Error()
		res.Write([]byte(errString))
	}
}

// Обработчик работы с задачами
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// получение задачи по id
	case http.MethodGet:
		taskGetIdHandler(w, r)
	// обновление задачи
	case http.MethodPut:
		taskEditHandler(w, r)
	// удаление задачи
	case http.MethodDelete:
		taskDeleteHandler(w, r)
	// добавление задачи
	case http.MethodPost:
		addTaskHandler(w, r)
	}
}

// Обработчик получения ID задачи
func taskGetIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJson(w, task, http.StatusOK)
}

// Обработчик обновления задачи
func taskEditHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	// десериализуем JSON
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	if task.ID == "" {
		writeJson(w, map[string]string{"error": "не указан идентификатор"}, http.StatusInternalServerError)
		return
	}
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "не указан заголовок задачи"}, http.StatusInternalServerError)
		return
	}
	// проверяем поля task
	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	// обновляем задачу в БД
	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJson(w, map[string]string{}, http.StatusOK)
}

// Обработчик удаления задачи
func taskDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "не указан идентификатор"}, http.StatusInternalServerError)
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	writeJson(w, map[string]string{}, http.StatusOK)
}

// Обработчик выполнения задачи
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		task, err := db.GetTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		if task.Repeat == "" {
			err := db.DeleteTask(id)
			if err != nil {
				writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
				return
			}
		} else {
			now := time.Now()
			if task.Date == now.Format(FormatDate) {
				now = now.AddDate(0, 0, 1)
			}
			strDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
				return
			}
			task.Date = strDate
			err = db.UpdateTask(task)
			if err != nil {
				writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
				return
			}
		}
		writeJson(w, map[string]string{}, http.StatusOK)
	} else {
		writeJson(w, map[string]string{"error": "неправильный метод"}, http.StatusMisdirectedRequest)
	}
}
