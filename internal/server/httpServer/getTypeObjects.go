package httpServer

import (
	"encoding/json"
	"fmt"
	"generatorOPCUA/internal/service"
	"io/ioutil"
	"net/http"
)

func GetTypeObjects(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== GetTypeObjects ===")
	fmt.Printf("Метод запроса: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL.String())
	fmt.Printf("Заголовки: %v\n", r.Header)

	// Добавляем CORS заголовки
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

	// Обрабатываем preflight OPTIONS запрос
	if r.Method == "OPTIONS" {
		fmt.Println("Обработка OPTIONS запроса")
		w.WriteHeader(http.StatusOK)
	}

	// Проверяем метод - теперь разрешаем POST
	if r.Method != "POST" {
		fmt.Printf("Ошибка: метод %s не поддерживается\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("Получен POST запрос на загрузку объектов из решения")

	// Читаем тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Ошибка чтения тела: %v\n", err)
		http.Error(w, "Error reading request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Printf("Тело запроса: %s\n", string(body))

	// Парсим JSON запрос
	var requestData struct {
		Filename string `json:"filename"`
	}

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		fmt.Printf("Ошибка парсинга JSON: %v\n", err)
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if requestData.Filename == "" {
		fmt.Println("Ошибка: filename пустой")
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	fmt.Printf("Загрузка объектов из файла: %s\n", requestData.Filename)

	// Получаем данные из файла
	objectsMap, err := service.ParsingSolution(requestData.Filename)
	if err != nil {
		fmt.Printf("Ошибка при парсинге файла: %v\n", err)
		http.Error(w, "Error parsing solution file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Получено объектов: %d\n", len(objectsMap))

	// Преобразуем мапу в нужный формат для фронтенда
	response := make(map[string]interface{})

	// Для каждого типа объекта формируем структуру данных
	for objectType, objects := range objectsMap {
		fmt.Printf("Обработка типа %s, объектов: %d\n", objectType, len(objects))

		typeData := make(map[string]interface{})

		var paramsList []map[string]interface{}
		var protectionsList []map[string]interface{}
		var commandsList []map[string]interface{}
		var objectsList []map[string]interface{}

		for _, obj := range objects {
			// Преобразуем в общий формат для списка объектов
			objMap := map[string]interface{}{
				"id":   obj.Id,
				"name": obj.Name,
				"tag":  obj.Tag,
				"use":  true,
			}
			objectsList = append(objectsList, objMap)

			// Разделяем по типу
			if obj.Template == "param" {
				paramsList = append(paramsList, map[string]interface{}{
					"description":  obj.Description,
					"templateName": obj.Template,
					"opcUaPath":    obj.Tag,
					"use":          true,
				})
			} else if obj.Template == "protection" {
				protectionsList = append(protectionsList, map[string]interface{}{
					"description":  obj.Description,
					"templateName": obj.Template,
					"opcUaPath":    obj.Tag,
					"use":          true,
				})
			} else if obj.Template == "command" {
				commandsList = append(commandsList, map[string]interface{}{
					"description":  obj.Description,
					"templateName": obj.Template,
					"opcUaPath":    obj.Tag,
					"use":          true,
				})
			}
		}

		typeData["params"] = paramsList
		typeData["protections"] = protectionsList
		typeData["commands"] = commandsList
		typeData["objects"] = objectsList

		response[objectType] = typeData
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Printf("Ошибка кодирования JSON: %v\n", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Успешно отправлены данные для %d типов объектов\n", len(response))
}

func init() {
	// Один для OPTIONS (preflight)
	AddNewFunction(CreateHandlerCommand("/GetTypeObjects", "OPTIONS", GetTypeObjects))

	// Второй для POST (основной запрос)
	AddNewFunction(CreateHandlerCommand("/GetTypeObjects", "POST", GetTypeObjects))
}
