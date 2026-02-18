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

// TemplateData структура для хранения данных шаблона
type TemplateData struct {
	Description string `json:"description"`
	Template    string `json:"templateName"`
	OPCUAPath   string `json:"opcUaPath"`
	Use         bool   `json:"use"`
}

// TemplateResponse структура ответа
type TemplateResponse map[string]map[string][]TemplateData

func GetTemplate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== GetTemplate ===")
	fmt.Printf("Метод запроса: %s\n", r.Method)

	// Добавляем CORS заголовки
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

	// Обрабатываем preflight OPTIONS запрос
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Проверяем метод
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Базовая директория с шаблонами
	templatesDir := "Шаблоны"

	// Получаем список всех типов из директории Свойства
	types, err := getTemplateTypes(filepath.Join(templatesDir, "Свойства"))
	if err != nil {
		fmt.Printf("Ошибка получения списка типов: %v\n", err)
		http.Error(w, "Error reading templates directory", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Найдены типы: %v\n", types)

	// Формируем ответ
	response := make(TemplateResponse)

	for _, templateType := range types {
		// Загружаем свойства
		properties, err := loadTemplateData(templatesDir, "Свойства", templateType)
		if err != nil {
			fmt.Printf("Ошибка загрузки свойств для %s: %v\n", templateType, err)
		}

		// Загружаем команды
		commands, err := loadTemplateData(templatesDir, "Команд", templateType)
		if err != nil {
			fmt.Printf("Ошибка загрузки команд для %s: %v\n", templateType, err)
		}

		// Загружаем защиты
		protections, err := loadTemplateData(templatesDir, "Защит", templateType)
		if err != nil {
			fmt.Printf("Ошибка загрузки защит для %s: %v\n", templateType, err)
		}

		// Создаем мапу для типа
		typeData := make(map[string][]TemplateData)

		if len(properties) > 0 {
			typeData["params"] = properties
		}
		if len(commands) > 0 {
			typeData["commands"] = commands
		}
		if len(protections) > 0 {
			typeData["protections"] = protections
		}

		// Добавляем в ответ, если есть хотя бы один тип данных
		if len(typeData) > 0 {
			response[templateType] = typeData
		}
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Ошибка кодирования JSON: %v\n", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Успешно отправлены шаблоны для %d типов\n", len(response))
}

// getTemplateTypes получает список типов из директории
func getTemplateTypes(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var types []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
			// Убираем расширение .csv
			typeName := strings.TrimSuffix(file.Name(), ".csv")
			types = append(types, typeName)
		}
	}

	return types, nil
}

// loadTemplateData загружает данные из CSV файла
func loadTemplateData(baseDir, category, templateType string) ([]TemplateData, error) {
	// Формируем путь к файлу
	filePath := filepath.Join(baseDir, category, templateType+".csv")

	fmt.Printf("Загрузка файла: %s\n", filePath)

	// Открываем файл
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Создаем CSV reader
	reader := csv.NewReader(file)
	reader.Comma = '\t' // Используем табуляцию как разделитель
	reader.TrimLeadingSpace = true

	// Читаем все записи
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []TemplateData
	for _, record := range records {
		if len(record) < 2 {
			continue // Пропускаем некорректные строки
		}

		// Очищаем поля от лишних пробелов
		description := strings.TrimSpace(record[0])
		opcUaPath := strings.TrimSpace(record[1])

		// Пропускаем пустые строки
		if description == "" || opcUaPath == "" {
			continue
		}

		data = append(data, TemplateData{
			Description: description,
			Template:    templateType, // Используем имя файла как template
			OPCUAPath:   opcUaPath,
			Use:         true, // По умолчанию включено
		})
	}

	fmt.Printf("Загружено %d записей из %s\n", len(data), filePath)
	return data, nil
}

// Альтернативная версия с поддержкой разных разделителей
func loadTemplateDataFlexible(baseDir, category, templateType string) ([]TemplateData, error) {
	filePath := filepath.Join(baseDir, category, templateType+".csv")

	fmt.Printf("Загрузка файла: %s\n", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Пробуем разные разделители
	delimiters := []rune{'\t', ',', ';', ' '}

	var data []TemplateData
	var lastErr error

	for _, delim := range delimiters {
		// Возвращаемся в начало файла
		file.Seek(0, 0)

		reader := csv.NewReader(file)
		reader.Comma = delim
		reader.TrimLeadingSpace = true
		reader.FieldsPerRecord = -1 // Разрешаем переменное количество полей

		records, err := reader.ReadAll()
		if err != nil {
			lastErr = err
			continue
		}

		data = make([]TemplateData, 0)
		for _, record := range records {
			if len(record) < 2 {
				continue
			}

			description := strings.TrimSpace(record[0])
			opcUaPath := strings.TrimSpace(record[1])

			if description == "" || opcUaPath == "" {
				continue
			}

			data = append(data, TemplateData{
				Description: description,
				Template:    templateType,
				OPCUAPath:   opcUaPath,
				Use:         true,
			})
		}

		if len(data) > 0 {
			fmt.Printf("Успешно загружено %d записей с разделителем '%c'\n", len(data), delim)
			return data, nil
		}
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return data, nil
}

// Функция для регистрации обработчика
func init() {
	// Если используется ваша система регистрации:
	AddNewFunction(CreateHandlerCommand("/GetTemplate", "GET", GetTemplate))
}
