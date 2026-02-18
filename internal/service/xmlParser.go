package service

import (
	"encoding/xml"
	"fmt"
	"generatorOPCUA/internal/domen"
	"os"
	"regexp"
)

func ParseCommandXML() map[string]domen.ObjCommandIdTag {
	var result []domen.ControlCommand
	data, err := os.ReadFile("data/Solution.xml")
	if err != nil {
		panic(err)
	}

	templateCommand := "<ControlCommand\\s+[^>]*?/>"
	re := regexp.MustCompile(templateCommand)
	matches := re.FindAllString(string(data), -1)

	fmt.Printf("Найдено объектов: %d\n", len(matches))

	for _, command := range matches {
		var tmp domen.ControlCommand
		err := xml.Unmarshal([]byte(command), &tmp)
		if err != nil {
			panic(err)
		}
		result = append(result, tmp)
	}
	mappingCommand := map[string]domen.ObjCommandIdTag{}
	for _, command := range result {
		if command.Template != "ValveOpen" && command.Template != "ValveClose" && command.Template != "ValveStop" {
			continue
		}

		tmp := domen.ObjCommandIdTag{}
		tmp.Id = command.Target

		if value, ok := mappingCommand[command.Target]; ok {
			tmp = value
		}

		if command.Template == "ValveOpen" {
			if len(tmp.OpenId) == 0 {
				tmp.OpenId = command.Id
			}
		} else if command.Template == "ValveClose" {
			if len(tmp.CloseId) == 0 {
				tmp.CloseId = command.Id
			}
		} else if command.Template == "ValveStop" {
			if len(tmp.StopId) == 0 {
				tmp.StopId = command.Id
			}
		}
		mappingCommand[command.Target] = tmp
	}

	fmt.Printf("Итоговое количество задвижек: %d\n", len(mappingCommand))
	return mappingCommand
}

func ParseCommandXMLToNameObj(name string) map[string]map[string]domen.CommandStruct {
	var result []domen.ControlCommand
	data, err := os.ReadFile("data/Solution.xml")
	if err != nil {
		panic(err)
	}
	templateCommand := "<ControlCommand\\s+[^>]*?/>"
	if name == "NPS" {
		templateCommand = "<ControlCommand\\b[^>]*>(?:.*?</ControlCommand>|/>)"
	}
	re := regexp.MustCompile(templateCommand)
	matches := re.FindAllString(string(data), -1)

	for _, command := range matches {
		var tmp domen.ControlCommand
		err := xml.Unmarshal([]byte(command), &tmp)
		if err != nil {
			fmt.Println(err)
		}
		result = append(result, tmp)
	}

	mapComm, err := MakeObjectsCommand("ШаблоныКомманд/"+name+".csv", '\t')

	if err != nil {
		fmt.Println(err)
	}

	var commands = make(map[string]map[string]domen.CommandStruct, len(result))

	fmt.Printf("Всего распарсеных комманд: %d\n", len(mapComm))
	fmt.Printf("Найдено объектов: %d\n", len(result))

	for _, command := range result {
		if _, ok := mapComm[command.Template]; !ok {
			continue
		}
		tmp, _ := mapComm[command.Template]
		tmp.Id = command.Id
		m, ok := commands[command.Target]
		if !ok {
			m = make(map[string]domen.CommandStruct)
		}
		m[command.Id] = tmp
		commands[command.Target] = m
	}

	fmt.Printf("Итоговое количество комманд: %d\n", len(commands))
	return commands
}

func ParseAutomationXMLToNameObj(name string) map[string]map[string]domen.AutomationStruct {
	var result []domen.Automation
	// <Automation\s+([^>]+)>(.*?)/>

	data, err := os.ReadFile("data/Solution.xml")
	if err != nil {
		panic(err)
	}

	templateAutomation := "<Automation\\b[^>]*>.*?</Automation>"
	re := regexp.MustCompile(templateAutomation)
	matches := re.FindAllString(string(data), -1)
	for _, automation := range matches {
		var tmp domen.Automation
		err := xml.Unmarshal([]byte(automation), &tmp)
		if err != nil {
			fmt.Printf("Ошибка с текстом: %s\n", automation)
		}
		result = append(result, tmp)
	}

	// start

	inf, err := MakeObjects("objects/"+name+".csv", '\t')

	unicMap := make(map[string]domen.Automation)

	for _, automation := range result {
		for _, value := range inf {
			if automation.Target == value.Id {
				if _, ok := unicMap[automation.Template]; !ok {
					unicMap[automation.Template] = automation
				}
			}
		}
	}

	// end

	mapAutomation, err := MakeObjectsAutomation("ШаблоныЗащит/"+name+".csv", '\t')
	if err != nil {
		fmt.Println(err)
	}
	var automations = make(map[string]map[string]domen.AutomationStruct)
	for _, automation := range result {
		if _, ok := mapAutomation[automation.Template]; !ok {
			continue
		}
		tmp, _ := mapAutomation[automation.Template]
		tmp.Id = automation.Id
		m, ok := automations[automation.Target]
		if !ok {
			m = make(map[string]domen.AutomationStruct)
		}
		m[automation.Id] = tmp
		automations[automation.Target] = m
	}

	return automations
}

func ParseCommandToNameObj(nameTemplate, nameSolution, body string, comma rune) map[string]map[string]domen.CommandStruct {
	var result []domen.ControlCommand
	data, err := os.ReadFile("data/" + nameSolution)
	if err != nil {
		fmt.Println(err)
	}
	templateCommand := "<ControlCommand\\s+[^>]*?/>"
	fmt.Printf("Проверка комманд\n")
	fmt.Println("Тип команды: " + nameTemplate)
	if nameTemplate == "PumpStation" || nameTemplate == "ShopBooster" {
		fmt.Printf("Добавляем новый regex для парсинга\n")
		templateCommand = "<ControlCommand\\b[^>]*>(?:.*?</ControlCommand>|/>)"
	}

	re := regexp.MustCompile(templateCommand)
	matches := re.FindAllString(string(data), -1)
	for _, command := range matches {
		var tmp domen.ControlCommand
		err := xml.Unmarshal([]byte(command), &tmp)
		if err != nil {
			fmt.Println(err)
		}
		result = append(result, tmp)
	}
	mapComm, err := MakeObjectCommandToString(body, comma)
	if err != nil {
		fmt.Println(err)
	}
	var commands = make(map[string]map[string]domen.CommandStruct, len(result))
	fmt.Printf("Всего распарсеных комманд: %d\n", len(mapComm))
	fmt.Printf("Найдено объектов: %d\n", len(result))

	for _, command := range result {
		if _, ok := mapComm[command.Template]; !ok {
			continue
		}
		tmp, _ := mapComm[command.Template]
		tmp.Id = command.Id
		m, ok := commands[command.Target]
		if !ok {
			m = make(map[string]domen.CommandStruct)
		}
		m[command.Id] = tmp
		commands[command.Target] = m
	}

	fmt.Printf("Итоговое количество комманд: %d\n", len(commands))
	return commands
}

func GetBodySolution(path, regex string) ([]string, error) {
	data, err := os.ReadFile("data/" + path)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(regex)
	matches := re.FindAllString(string(data), -1)
	return matches, nil
}

func ParseAutomationToBody(nameTemplate, nameSolution, body string, comma rune) map[string]map[string]domen.AutomationStruct {
	fmt.Printf("Старт генерации защит для типа объектов: %s\n", nameTemplate)
	var result []domen.Automation

	templateAutomation := "<Automation\\b[^>]*>.*?</Automation>"
	matches, err := GetBodySolution(nameSolution, templateAutomation)
	if err != nil {
		fmt.Println(err)
	}

	for _, automation := range matches {
		var tmp domen.Automation
		err := xml.Unmarshal([]byte(automation), &tmp)

		if err != nil {
			fmt.Printf("Ошибка с текстом: %s\n", automation)
		}
		result = append(result, tmp)
	}

	// start

	inf, err := MakeObjectsToString(body, comma)

	fmt.Printf("Количество объектов объектов: %d\n", len(inf))

	unicMap := make(map[string]domen.Automation)

	for _, automation := range result {
		for _, value := range inf {
			if automation.Target == value.Id {
				if _, ok := unicMap[automation.Template]; !ok {
					unicMap[automation.Template] = automation
				}
			}
		}
	}

	mapAutomation, err := MakeObjectsAutomationToString(body, comma)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Размер мапы: %d\n", len(mapAutomation))
	var automations = make(map[string]map[string]domen.AutomationStruct)
	for _, automation := range result {
		if _, ok := mapAutomation[automation.Template]; !ok {
			continue
		}
		tmp, _ := mapAutomation[automation.Template]
		tmp.Id = automation.Id
		m, ok := automations[automation.Target]
		if !ok {
			m = make(map[string]domen.AutomationStruct)
		}
		m[automation.Id] = tmp
		automations[automation.Target] = m
	}

	return automations
}
