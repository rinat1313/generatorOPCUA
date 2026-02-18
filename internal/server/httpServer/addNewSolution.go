package httpServer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func AddSolution(w http.ResponseWriter, r *http.Request) {
	// Добавляем CORS заголовки
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

	// Обрабатываем preflight OPTIONS запрос
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Проверяем метод
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Printf("Пришел запрос на добавление файла: %s\n", r.Host)

	// Парсим multipart form (32 MB максимум)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем файл из формы
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error getting file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Printf("Получен файл: %s, размер: %d байт\n", handler.Filename, handler.Size)

	// Создаем директорию data если её нет
	uploadDir := "./data"
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		http.Error(w, "Error creating directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем имя файла без расширения и расширение
	filename := handler.Filename
	ext := filepath.Ext(filename)
	nameWithoutExt := filename[:len(filename)-len(ext)]

	// Генерируем уникальное имя файла
	uniqueFilename, err := generateUniqueFilename(uploadDir, nameWithoutExt, ext)
	if err != nil {
		http.Error(w, "Error generating unique filename: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем файл на диске
	filePath := filepath.Join(uploadDir, uniqueFilename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error creating file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Копируем содержимое
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "success",
		"message":  fmt.Sprintf("Файл %s успешно загружен как %s", handler.Filename, uniqueFilename),
		"filename": uniqueFilename,
		"filepath": filePath,
	})

	fmt.Printf("Файл сохранен как: %s\n", uniqueFilename)
}

// Функция для генерации уникального имени файла
func generateUniqueFilename(dir, nameWithoutExt, ext string) (string, error) {
	// Проверяем существует ли файл с оригинальным именем
	originalPath := filepath.Join(dir, nameWithoutExt+ext)
	_, err := os.Stat(originalPath)

	if os.IsNotExist(err) {
		// Файл не существует, можно использовать оригинальное имя
		return nameWithoutExt + ext, nil
	}

	if err != nil {
		// Другая ошибка при проверке файла
		return "", err
	}

	// Файл существует, ищем свободное имя с номером
	counter := 1
	for {
		// Формируем имя с номером: имя_1.xml, имя_2.xml, и т.д.
		numberedName := fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, ext)
		numberedPath := filepath.Join(dir, numberedName)

		_, err := os.Stat(numberedPath)
		if os.IsNotExist(err) {
			// Нашли свободное имя
			return numberedName, nil
		}

		if err != nil {
			// Другая ошибка при проверке файла
			return "", err
		}

		counter++

		// Защита от бесконечного цикла (максимум 9999 попыток)
		if counter > 9999 {
			return "", fmt.Errorf("слишком много файлов с похожим именем")
		}
	}
}

func init() {
	AddNewFunction(CreateHandlerCommand("/AddSolution", "POST", AddSolution))
}
