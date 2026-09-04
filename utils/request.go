package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func DecodeJSON(r *http.Request, destination any) error {
	decoder := json.NewDecoder(r.Body)

	// Запрещаем неизвестные поля в JSON.
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		return err
	}

	// Проверяем, что после первого JSON-объекта больше ничего нет.
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return errors.New("тело запроса должно содержать один JSON-объект")
	}

	return nil
}

func ParseIDFromPath(path string, prefix string) (int, error) {
	idString := strings.TrimPrefix(path, prefix)

	if idString == path || idString == "" {
		return 0, errors.New("ID отсутствует")
	}

	if strings.Contains(idString, "/") {
		return 0, errors.New("некорректный путь")
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, errors.New("ID должен быть числом")
	}

	if id <= 0 {
		return 0, errors.New("ID должен быть положительным")
	}

	return id, nil
}
