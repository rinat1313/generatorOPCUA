package service

import (
	"encoding/csv"
	"errors"
	"generatorOPCUA/internal/domen"
	"io"
	"os"
	"strings"
)

func makeObjArray(text []string) (domen.Object, error) {
	if len(text) != 2 {
		return domen.Object{}, errors.New("несоответствует размер массива")
	}
	result := domen.Object{}
	result.Id = text[0]
	result.Tag = text[1]
	return result, nil
}

func MakeObjects(path string, comma rune) ([]domen.Object, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	r.Comma = comma

	var objects []domen.Object
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		info, err := makeObjArray(record)
		if err != nil {
			return nil, err
		}
		objects = append(objects, info)
	}
	return objects, nil
}

func MakeObjectsToString(body string, comma rune) ([]domen.Object, error) {
	var objects []domen.Object
	reader := csv.NewReader(strings.NewReader(body))
	reader.Comma = comma
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		info, err := makeObjArray(record)
		if err != nil {
			return nil, err
		}
		objects = append(objects, info)
	}
	return objects, nil
}

func makeDpeTag(text []string) (domen.DpeTag, error) {
	if len(text) != 2 {
		return domen.DpeTag{}, nil
	}
	result := domen.DpeTag{Dpe: text[0], Tag: text[1]}
	return result, nil
}

func MakeDpeTags(path string, comma rune) (map[string]string, error) {
	f, err := os.Open("Выборка/" + path + ".csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = comma

	var result = make(map[string]string)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		info, err := makeDpeTag(record)
		if err != nil {
			return nil, err
		}

		if !strings.Contains(info.Dpe, ".") {
			if _, ok := result[info.Tag]; !ok {
				result[info.Tag] = info.Dpe
			} else {
				if len(result[info.Tag]) > len(info.Dpe) {
					result[info.Tag] = info.Dpe
				}
			}
		}

	}

	return result, nil
}

func MakeCommandInfo(text []string) (domen.CommandStruct, error) {
	if len(text) != 2 {
		return domen.CommandStruct{}, nil
	}
	result := domen.CommandStruct{}
	result.NameTemplate = text[0]
	result.AfterHeaderTag = text[1]
	return result, nil
}

func MakeObjectsCommand(path string, comma rune) (map[string]domen.CommandStruct, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(f)
	r.Comma = comma
	var commands = make(map[string]domen.CommandStruct)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		info, err := MakeCommandInfo(record)
		if err != nil {
			return nil, err
		}
		commands[info.NameTemplate] = info
	}
	return commands, nil
}

func MakeObjectCommandToString(body string, comma rune) (map[string]domen.CommandStruct, error) {
	var commands = make(map[string]domen.CommandStruct)
	for _, line := range strings.Split(body, "\n") {
		splitLine := strings.Split(line, string(comma))
		info, err := MakeCommandInfo(splitLine)
		if err != nil {
			return nil, err
		}
		commands[info.NameTemplate] = info
	}
	return commands, nil
}

func MakeAutomationInfo(text []string) (domen.AutomationStruct, error) {
	if len(text) != 2 {
		return domen.AutomationStruct{}, nil
	}
	result := domen.AutomationStruct{}
	result.NameTemplate = text[0]
	result.AfterHeaderTag = text[1]
	return result, nil
}

func MakeObjectsAutomation(path string, comma rune) (map[string]domen.AutomationStruct, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(f)
	r.Comma = comma
	var automations = make(map[string]domen.AutomationStruct)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
		}
		info, err := MakeAutomationInfo(record)
		if err != nil {
			return nil, err
		}
		automations[info.NameTemplate] = info
	}
	return automations, nil
}

func MakeObjectsAutomationToString(body string, comma rune) (map[string]domen.AutomationStruct, error) {
	var automations = make(map[string]domen.AutomationStruct)
	for _, line := range strings.Split(body, "\n") {
		splitLine := strings.Split(line, string(comma))
		info, err := MakeAutomationInfo(splitLine)
		if err != nil {
			return nil, err
		}
		automations[info.NameTemplate] = info
	}
	return automations, nil
}
