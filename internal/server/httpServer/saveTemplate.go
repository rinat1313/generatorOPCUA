package httpServer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// SaveTemplatesRequest структура для входящего запроса
type SaveTemplatesRequest struct {
	ObjectType  string `json:"ObjectType"`
	MappingType string `json:"MappingType"`
	Data        []struct {
		Description string `json:"Description"`
		OpcUaPath   string `json:"OpcUaPath"`
		Use         bool   `json:"Use"`
	} `json:"Data"`
}

// SaveTemplatesResponse структура для ответа
type SaveTemplatesResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Filename string `json:"filename"`
	Count    int    `json:"count"`
}

func SaveTemplates(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== SaveTemplates ===")
	fmt.Printf("Метод запроса: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL.Path)
	fmt.Printf("Заголовки: %v\n", r.Header)

	// Добавляем CORS заголовки для ВСЕХ ответов
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
	w.Header().Set("Access-Control-Max-Age", "3600") // Кэшировать preflight на 1 час

	// Обрабатываем preflight OPTIONS запрос - ЭТО ВАЖНО!
	if r.Method == "OPTIONS" {
		fmt.Println("Обработка OPTIONS запроса для /SaveTemplates")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Проверяем метод
	if r.Method != "POST" {
		fmt.Printf("Ошибка: метод %s не поддерживается\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем JSON запрос
	var req SaveTemplatesRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Printf("Ошибка декодирования JSON: %v\n", err)
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Валидация входных данных
	if req.ObjectType == "" {
		http.Error(w, "ObjectType is required", http.StatusBadRequest)
		return
	}

	if req.MappingType == "" {
		http.Error(w, "MappingType is required", http.StatusBadRequest)
		return
	}

	if len(req.Data) == 0 {
		http.Error(w, "Data is required", http.StatusBadRequest)
		return
	}

	fmt.Printf("Получены данные: ObjectType=%s, MappingType=%s, записей=%d\n",
		req.ObjectType, req.MappingType, len(req.Data))

	// Определяем поддиректорию в зависимости от типа сопоставления
	var subDir string
	switch req.MappingType {
	case "params":
		subDir = "Свойства"
	case "protections":
		subDir = "Защит"
	case "commands":
		subDir = "Команд"
	default:
		http.Error(w, "Invalid MappingType. Must be 'params', 'protections' or 'commands'", http.StatusBadRequest)
		return
	}

	// Создаем полный путь к директории
	templatesDir := "Шаблоны"
	fullPath := filepath.Join(templatesDir, subDir)

	// Создаем директорию, если она не существует
	err = os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		fmt.Printf("Ошибка создания директории: %v\n", err)
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return
	}

	// Формируем имя файла (Тип объекта + .csv)
	filename := req.ObjectType + ".csv"
	filePath := filepath.Join(fullPath, filename)

	fmt.Printf("Сохранение в файл: %s\n", filePath)

	// Открываем файл для записи
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Ошибка создания файла: %v\n", err)
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Создаем CSV writer с табуляцией как разделителем
	writer := csv.NewWriter(file)
	writer.Comma = '\t'
	defer writer.Flush()

	// Записываем данные
	validCount := 0
	for _, item := range req.Data {
		// Пропускаем записи с пустыми полями
		if item.Description == "" || item.OpcUaPath == "" {
			fmt.Printf("Пропущена запись с пустыми полями: %+v\n", item)
			continue
		}

		// Создаем запись для CSV
		record := []string{
			strings.TrimSpace(item.Description),
			strings.TrimSpace(item.OpcUaPath),
		}

		err := writer.Write(record)
		if err != nil {
			fmt.Printf("Ошибка записи в CSV: %v\n", err)
			continue
		}
		validCount++
	}

	// Проверяем, были ли записаны данные
	if validCount == 0 {
		http.Error(w, "No valid data to save", http.StatusBadRequest)
		return
	}

	// Формируем успешный ответ
	response := SaveTemplatesResponse{
		Status:   "success",
		Message:  fmt.Sprintf("Файл успешно сохранен: %s", filename),
		Filename: filename,
		Count:    validCount,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Ошибка кодирования JSON ответа: %v\n", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Успешно сохранено %d записей в файл: %s\n", validCount, filePath)
}

func init() {
	AddNewFunction(CreateHandlerCommand("/SaveTemplates", "POST", SaveTemplates))
	AddNewFunction(CreateHandlerCommand("/SaveTemplates", "OPTIONS", SaveTemplates))
}
