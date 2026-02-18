package httpServer

import (
	"encoding/json"
	"fmt"
	"generatorOPCUA/internal/service"
	"net/http"
)

func GetAllSolution(w http.ResponseWriter, r *http.Request) {
	// Добавляем заголовки CORS для работы с браузером
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

	fmt.Printf("Получен запрос на получение всех Solution\n")

	// Получаем список файлов
	files, err := service.GetFilesIsDir("data")
	if err != nil {
		fmt.Printf("Ошибка получения списка файлов: %v\n", err)
		http.Error(w, "Ошибка чтения директории", http.StatusInternalServerError)
		return
	}

	// Создаем канал для обработки ошибок (если нужно параллельное выполнение)
	type result struct {
		index int
		desc  string
		err   error
	}
	resultChan := make(chan result, len(files))

	// Параллельно получаем описания файлов
	for i, file := range files {
		go func(index int, filename string) {
			desc, err := service.GetNameSolution("data/" + filename)
			resultChan <- result{index: index, desc: desc, err: err}
		}(i, file.NameFile)
	}

	// Собираем результаты
	for i := 0; i < len(files); i++ {
		res := <-resultChan
		if res.err != nil {
			fmt.Printf("Ошибка получения описания для файла %s: %v\n",
				files[res.index].NameFile, res.err)
			// Продолжаем с пустым описанием или пропускаем
			files[res.index].Discript = "Описание недоступно"
		} else {
			files[res.index].Discript = res.desc
		}
	}
	close(resultChan)

	// Устанавливаем заголовки и отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(files); err != nil {
		fmt.Printf("Ошибка кодирования JSON: %v\n", err)
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Успешно отправлено %d решений\n", len(files))

}

func init() {
	AddNewFunction(CreateHandlerCommand("/getAllSolution", "GET", GetAllSolution))
}
