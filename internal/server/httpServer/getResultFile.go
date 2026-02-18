package httpServer

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"generatorOPCUA/internal/domen"
	"generatorOPCUA/internal/service"
	"io/ioutil"
	"net/http"
	"strings"
)

// Структуры для входящего запроса
type GenerateRequest struct {
	Filename   string                   `json:"filename"`
	ObjectType string                   `json:"objectType"`
	Options    GenerateOptions          `json:"options"`
	Objects    []GenerateObject         `json:"objects"`
	Mappings   map[string][]MappingItem `json:"mappings,omitempty"`
}

type GenerateOptions struct {
	HeadTags    bool `json:"headTags"`
	Properties  bool `json:"properties"`
	Commands    bool `json:"commands"`
	Protections bool `json:"protections"`
}

type GenerateObject struct {
	ID  string `json:"id"`
	Tag string `json:"tag"`
}

type MappingItem struct {
	Description string `json:"description"`
	Tag         string `json:"tag"`
}

// Структура для ответа (XML)
type ResultXML struct {
	XMLName     xml.Name  `xml:"Result"`
	ObjectType  string    `xml:"ObjectType"`
	HeadTags    []HeadTag `xml:"HeadTags>HeadTag,omitempty"`
	Properties  []Mapping `xml:"Properties>Property,omitempty"`
	Commands    []Mapping `xml:"Commands>Command,omitempty"`
	Protections []Mapping `xml:"Protections>Protection,omitempty"`
}

type HeadTag struct {
	ObjectID string `xml:"ObjectID,attr"`
	Tag      string `xml:",chardata"`
}

type Mapping struct {
	Description string `xml:"Description,attr"`
	Tag         string `xml:",chardata"`
}

func GetResultFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("=== GetResultFile ===")
	fmt.Printf("Метод запроса: %s\n", r.Method)
	fmt.Printf("URL: %s\n", r.URL.Path)

	// Добавляем CORS заголовки
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

	// Обрабатываем preflight OPTIONS запрос
	if r.Method == "OPTIONS" {
		fmt.Println("Обработка OPTIONS запроса")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Проверяем метод
	if r.Method != "POST" {
		fmt.Printf("Ошибка: метод %s не поддерживается\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Читаем тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Ошибка чтения тела: %v\n", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//fmt.Printf("Тело запроса: %s\n", string(body))

	// Декодируем JSON запрос
	var req GenerateRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		fmt.Printf("Ошибка декодирования JSON: %v\n", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Валидация
	if req.Filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	if req.ObjectType == "" {
		http.Error(w, "ObjectType is required", http.StatusBadRequest)
		return
	}

	if len(req.Objects) == 0 {
		http.Error(w, "Objects list is required", http.StatusBadRequest)
		return
	}

	fmt.Printf("Получен запрос на генерацию: ObjectType=%s, файл=%s, объектов=%d\n",
		req.ObjectType, req.Filename, len(req.Objects))
	fmt.Printf("Опции: HeadTags=%v, Properties=%v, Commands=%v, Protections=%v\n",
		req.Options.HeadTags, req.Options.Properties, req.Options.Commands, req.Options.Protections)

	var commands map[string]map[string]domen.CommandStruct
	var automations map[string]map[string]domen.AutomationStruct

	commandsResult := ""

	if req.Options.Commands {
		bodyResult := ""
		comma := '\t'
		for _, command := range req.Mappings["commands"] {
			bodyResult += command.Description + string(comma) + command.Tag + "\n"
		}
		commands = service.ParseCommandToNameObj(req.ObjectType, req.Filename, bodyResult, comma)
	}

	if req.Options.Properties {
		bodyResult := ""
		comma := '\t'
		for _, command := range req.Mappings["properties"] {
			bodyResult += command.Description + string(comma) + command.Tag + "\n"
		}
		automations = service.ParseAutomationToBody(req.ObjectType, req.Filename, bodyResult, comma)
	}

	headerResult := ""

	if req.Options.HeadTags {
		for _, object := range req.Objects {
			headerResult += "\t<object id=\"" + object.ID + "\" tag=\"" + object.Tag + "\" />\n"
		}
	}

	if req.Options.Commands {
		for _, object := range req.Objects {
			if command, ok := commands[object.ID]; ok {
				for _, value := range command {
					commandsResult += "\t<object id=\"" + value.Id + "\" tag=\"" + object.Tag + value.AfterHeaderTag + "\" />\n"
				}
			}
		}
	}

	automatoResult := ""

	if req.Options.Protections {
		for _, object := range req.Objects {
			if automation, ok := automations[object.ID]; ok {
				for _, value := range automation {
					automatoResult += "\t<object id=\"" + value.Id + "\" tag=\"" + object.Tag + value.AfterHeaderTag + "\" />\n"
				}
			}
		}
	}

	propertiesResult := ""
	if req.Options.Properties {
		for _, object := range req.Objects {
			for _, afterOPCUA := range req.Mappings["params"] {
				propertiesResult += "\t<property id=\"" + object.ID + "\" key=\"" + afterOPCUA.Description + "\" tag=\"" + object.Tag + "." + afterOPCUA.Tag + "\" />\n"
			}
		}
	}

	// =====================================================
	// ЗДЕСЬ ВЫЗОВ ВАШЕЙ ФУНКЦИИ ДЛЯ ГЕНЕРАЦИИ XML ФАЙЛА
	// =====================================================
	// Предположим, у вас есть функция GenerateXML, которая принимает данные
	// и возвращает сгенерированный XML в виде []byte
	//
	// Пример:
	// xmlData, err := service.GenerateResultXML(req)
	// if err != nil {
	//     fmt.Printf("Ошибка генерации XML: %v\n", err)
	//     http.Error(w, "Error generating XML", http.StatusInternalServerError)
	//     return
	// }
	// =====================================================

	// Для примера создадим XML из структур

	resultXml := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<root>\n  <objects>\n"
	resultXml += headerResult
	resultXml += commandsResult
	resultXml += automatoResult
	resultXml += "  </objects>\n  <properties>\n"
	resultXml += propertiesResult
	resultXml += "  </properties>\n"
	resultXml += "</root>"

	fmt.Printf("Сгенерирован XML размером %d байт\n", len(resultXml))

	// Формируем имя файла
	filename := fmt.Sprintf("result_%s_%s.xml",
		strings.ReplaceAll(req.ObjectType, " ", "_"),
		strings.ReplaceAll(strings.ReplaceAll(req.Filename, ".xml", ""), " ", "_"))

	// Устанавливаем заголовки для скачивания файла
	w.Header().Set("Content-Type", "application/xml")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(resultXml)))
	w.WriteHeader(http.StatusOK)

	// Отправляем файл - просто отправляем строку как массив байт
	_, err = w.Write([]byte(resultXml))
	if err != nil {
		fmt.Printf("Ошибка отправки файла: %v\n", err)
	} else {
		fmt.Printf("Файл %s успешно отправлен клиенту\n", filename)
	}
}

func init() {
	AddNewFunction(CreateHandlerCommand("/GetResultFile", "POST", GetResultFile))
	AddNewFunction(CreateHandlerCommand("/GetResultFile", "OPTIONS", GetResultFile))
}
