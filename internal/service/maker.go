package service

import (
	"fmt"
	"os"
	"strings"
)

func StartMakeObjectInCsvFile(path string, comm rune, templateType string) {
	inf, err := MakeObjects(path, comm)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("первый элемент до: %s\n", inf[0])

	mapTagDpe, err := MakeDpeTags(templateType, comm)
	if err != nil {
		fmt.Println(err)
	}

	for num, value := range inf {

		if newTag, ok := mapTagDpe[value.Tag]; ok {
			value.Tag = newTag
			inf[num] = value
		}
	}

	fmt.Printf("первый элемент после: %s\n", inf[0])

	fmt.Printf("Выбранный шаблон: %s\n", templateType)

	resultText := ""

	objectTmp, err := os.ReadFile("templates/" + templateType + ".txt")
	if err != nil {
		fmt.Println(err)
	}

	resultTextCommand := ""
	objCommandTmp, err := os.ReadFile("ШаблоныКомманд/" + templateType + ".txt")
	if err != nil {
		fmt.Println(err)
	}

	//resultObjectsIdAndTags :=""

	fmt.Printf("Количество тегов: %d\n", len(inf))
	count := 0

	commands := ParseCommandXML()
	fmt.Printf("Полученно команд: %d\n", len(commands))
	for _, obj := range inf {

		tmpObj := string(objectTmp)
		tmpObj = strings.Replace(tmpObj, "<ID>", obj.Id, -1)
		tmpObj = strings.Replace(tmpObj, "<TAG>", obj.Tag, -1)
		resultText += tmpObj + "\n"

		if command, ok := commands[obj.Id]; ok {
			tmpObjCom := string(objCommandTmp)
			tmpObjCom = strings.Replace(tmpObjCom, "<IDOPEN>", command.OpenId, -1)
			tmpObjCom = strings.Replace(tmpObjCom, "<IDCLOSE>", command.CloseId, -1)
			tmpObjCom = strings.Replace(tmpObjCom, "<IDSTOP>", command.StopId, -1)
			tmpObjCom = strings.Replace(tmpObjCom, "<TAG>", obj.Tag, -1)
			resultTextCommand += tmpObjCom
		}

		count++
	}
	fmt.Printf("Итоговое количество тегов: %d\n", count)

	err = CreateFile("result/"+templateType+".txt", resultText)
	if err != nil {
		fmt.Println(err)
	}

	err = CreateFile("result/"+templateType+"Command.txt", resultTextCommand)
	if err != nil {
		fmt.Println(err)
	}
}

func StartMakeObjectInCsvFileAndCommand(path string, comm rune, templateType string) {
	inf, err := MakeObjects(path, comm)
	if err != nil {
		fmt.Println(err)
	}

	mapTagDpe, err := MakeDpeTags(templateType, comm)
	if err != nil {
		fmt.Println(err)
	}

	for num, value := range inf {

		if newTag, ok := mapTagDpe[value.Tag]; ok {
			value.Tag = newTag
			inf[num] = value
		}
	}

	fmt.Printf("Выбранный шаблон: %s\n", templateType)

	resultText := ""

	objectTmp, err := os.ReadFile("templates/" + templateType + ".txt")
	if err != nil {
		fmt.Println(err)
	}

	resultTextCommand := ""
	resultTextAutomation := ""
	resultObjectsIdAndTag := ""

	fmt.Printf("Количество тегов: %d\n", len(inf))
	count := 0

	commands := ParseCommandXMLToNameObj(templateType)
	fmt.Printf("Количество комманд: %d\n", len(commands))

	automations := ParseAutomationXMLToNameObj(templateType)
	fmt.Printf("Количество защит: %d\n", len(automations))

	for _, obj := range inf {
		if obj.Tag == "" {
			continue
		}
		tmpObj := string(objectTmp)
		tmpObj = strings.Replace(tmpObj, "<ID>", obj.Id, -1)
		tmpObj = strings.Replace(tmpObj, "<TAG>", obj.Tag, -1)
		resultText += tmpObj + "\n"

		if command, ok := commands[obj.Id]; ok {
			for _, value := range command {
				resultTextCommand += "<object id=\"" + value.Id + "\" tag=\"" + obj.Tag + value.AfterHeaderTag + "\" />\n"
			}
		}

		if automation, ok := automations[obj.Id]; ok {
			for _, value := range automation {
				resultTextAutomation += "<object id=\"" + value.Id + "\" tag=\"" + obj.Tag + value.AfterHeaderTag + "\" />\n"
			}
		}

		resultObjectsIdAndTag += "<object id=\"" + obj.Id + "\" tag=\"" + obj.Tag + "\" />\n"

		count++
	}
	fmt.Printf("Итоговое количество тегов: %d\n", count)

	err = CreateFile("result/"+templateType+"properties.txt", resultText)
	if err != nil {
		fmt.Println(err)
	}

	err = CreateFile("result/"+templateType+"Command.txt", resultTextCommand)
	if err != nil {
		fmt.Println(err)
	}

	err = CreateFile("result/"+templateType+"Automation.txt", resultTextAutomation)
	if err != nil {
		fmt.Println(err)
	}

	err = CreateFile("result/"+templateType+"Objects.txt", resultObjectsIdAndTag)
	if err != nil {
		fmt.Println(err)
	}
}
