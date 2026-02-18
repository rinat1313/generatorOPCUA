package service

import (
	"encoding/xml"
	"fmt"
	"generatorOPCUA/internal/domen"
	"io/ioutil"
	"os"
	"regexp"
)

func ParsingSolution(path string) (map[string][]domen.TechnologyObject, error) {
	regTemplate := `<TechnologyObject\s+[^>]+>.*?</TechnologyObject>`
	fi, err := os.Open("data/" + path)
	if err != nil {
		fmt.Println(err)
	}
	defer fi.Close()
	text, err := ioutil.ReadAll(fi)
	if err != nil {
		fmt.Println(err)
	}

	//var objects []domen.TechProp
	mapObjectsIsType := make(map[string][]domen.TechnologyObject)
	var object domen.TechnologyObject

	march := regexp.MustCompile(regTemplate)
	matches := march.FindAll(text, -1)
	for _, m := range matches {

		err := xml.Unmarshal(m, &object)
		if err != nil {
			fmt.Printf("Ошибка при unmarshal: %s\n", err)
		}

		if _, ok := mapObjectsIsType[object.Template]; !ok {
			mapObjectsIsType[object.Template] = make([]domen.TechnologyObject, 0)
		}

		mapObjectsIsType[object.Template] = append(mapObjectsIsType[object.Template], object)
	}

	return mapObjectsIsType, nil
}

func GetNameSolution(path string) (string, error) {
	regSolution := "<Solution\\s+[^>]+>.*?</Solution>"
	//regSolution := "<Solution[^>]*\">"
	fi, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer fi.Close()
	text, err := ioutil.ReadAll(fi)
	if err != nil {
		fmt.Println(err)
	}
	var solution domen.Solution
	match := regexp.MustCompile(regSolution)
	matches := match.FindAll(text, -1)
	err = xml.Unmarshal(matches[0], &solution)
	if err != nil {
		fmt.Printf("Ошибка парсинга : %s\n", err)
	}
	return solution.Name + " (" + solution.Description + ")", nil
}

func GetFilesIsDir(path string) ([]domen.FileSolution, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var result []domen.FileSolution

	for _, file := range files {
		if !file.IsDir() {
			var sol domen.FileSolution
			sol.NameFile = file.Name()
			result = append(result, sol)
		}
	}

	return result, nil
}
