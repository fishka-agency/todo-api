package models

import "errors"

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Validate проверка задачи на валидность
// @return error - ошибка
func (t *Task) Validate() error {
	if t.Title == "" { // проверка на наличие названия задачи
		return errors.New("title is required") // возврат ошибки
	}

	// возврат nil, если задача валидна
	return nil
}
