package domen

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type Object struct {
	Id  string
	Tag string
}

type DpeTag struct {
	Dpe string
	Tag string
}

// ControlCommand - команда управления
type ControlCommand struct {
	XMLName      xml.Name `xml:"ControlCommand"`
	Target       string   `xml:"Target,attr"`
	Id           string   `xml:"Id,attr"`
	IdFromSource string   `xml:"IdFromSource,attr"`
	Name         string   `xml:"Name,attr"`
	Description  string   `xml:"Description,attr"`
	Tag          string   `xml:"Tag,attr"`
	ExternalCode string   `xml:"ExternalCode,attr"`
	Template     string   `xml:"Template,attr"`
}

type ObjCommandIdTag struct {
	Id      string
	OpenId  string
	CloseId string
	StopId  string
}

type CommandStruct struct {
	Id             string
	NameTemplate   string
	AfterHeaderTag string
}

type AutomationStruct struct {
	Id             string
	NameTemplate   string
	AfterHeaderTag string
}

type Automation struct {
	XMLName              xml.Name             `xml:"Automation"`
	Target               string               `xml:"Target,attr"`
	Id                   string               `xml:"Id,attr"`
	IdFromSource         string               `xml:"IdFromSource,attr"`
	Name                 string               `xml:"Name,attr"`
	Description          string               `xml:"Description,attr"`
	Tag                  string               `xml:"Tag,attr"`
	ExternalCode         string               `xml:"ExternalCode,attr"`
	Template             string               `xml:"Template,attr"`
	TechnologyProperties []TechnologyProperty `xml:"TechnologyProperty"`
}

type TechnologyProperty struct {
	Id           string `xml:"Id,attr"`
	IdFromSource string `xml:"IdFromSource,attr"`
	Name         string `xml:"Name,attr"`
	Description  string `xml:"Description,attr"`
	Tag          string `xml:"Tag,attr"`
	ExternalCode string `xml:"ExternalCode,attr"`
	Template     string `xml:"Template,attr"`

	// Используем кастомный тип для обработки запятой в числах
	RealValue    CommaFloat64 `xml:"RealValue,attr,omitempty"`
	StringValue  string       `xml:"StringValue,attr,omitempty"`
	IntegerValue int          `xml:"IntegerValue,attr,omitempty"`

	// Вложенные ссылки (например, на задвижки в ValvesNotToClose)
	TechnologyObjectLinks []TechnologyObjectLink `xml:"TechnologyObjectLink"`
}

type TechnologyObjectLink struct {
	Id               string `xml:"Id,attr"`
	IdFromSource     string `xml:"IdFromSource,attr"`
	TechnologyObject string `xml:"TechnologyObject,attr"`
}

type CommaFloat64 float64

func (cf *CommaFloat64) UnmarshalXMLAttr(attr xml.Attr) error {
	s := strings.Replace(attr.Value, ",", ".", 1)
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*cf = CommaFloat64(f)
	return nil
}

type TechProp struct {
	XMLName      xml.Name `xml:"TechnologyProperty"`
	Id           string   `xml:"Id,attr"`
	IdFromSource string   `xml:"IdFromSource,attr"`
	Name         string   `xml:"Name,attr"`
	Description  string   `xml:"Description,attr"`
	Tag          string   `xml:"Tag,attr"`
	ExternalCode string   `xml:"ExternalCode,attr"`
	Template     string   `xml:"Template,attr"`

	// Разные типы значений
	RealValue     string `xml:"RealValue,attr,omitempty"`
	BooleanValue  string `xml:"BooleanValue,attr,omitempty"`
	EnumItemValue string `xml:"EnumItemValue,attr,omitempty"`
	IntegerValue  string `xml:"IntegerValue,attr,omitempty"`
	LinkValue     string `xml:"LinkValue,attr,omitempty"`
	StringValue   string `xml:"StringValue,attr,omitempty"`

	TechnologyPassport *TechPass `xml:"TechnologyPassport,omitempty"`
}

type TechPass struct {
	XMLName      xml.Name `xml:"TechnologyPassport"`
	Id           string   `xml:"Id,attr"`
	IdFromSource string   `xml:"IdFromSource,attr"`
	Name         string   `xml:"Name,attr"`
	Description  string   `xml:"Description,attr"`
	Tag          string   `xml:"Tag,attr"`
	ExternalCode string   `xml:"ExternalCode,attr"`
	Template     string   `xml:"Template,attr"`

	TechnologyProperties []TechProp `xml:"TechnologyProperty"`
}

type TechnologyObject struct {
	XMLName      xml.Name `xml:"TechnologyObject"`
	Id           string   `xml:"Id,attr"`
	IdFromSource string   `xml:"IdFromSource,attr"`
	Name         string   `xml:"Name,attr"`
	Description  string   `xml:"Description,attr"`
	Tag          string   `xml:"Tag,attr"`
	ExternalCode string   `xml:"ExternalCode,attr"`
	Template     string   `xml:"Template,attr"`

	TechnologyProperties []TechProp `xml:"TechnologyProperty"`
}

type Solution struct {
	XMLName             xml.Name `xml:"Solution"`
	Id                  string   `xml:"Id,attr"`
	Name                string   `xml:"Name,attr"`
	Description         string   `xml:"Description,attr"`
	Tag                 string   `xml:"Tag,attr"`
	ExternalCode        string   `xml:"ExternalCode,attr"`
	StepIterationsCount string   `xml:"StepIterationsCount,attr"`
	Company             string   `xml:"Company,attr"`
}

type FileSolution struct {
	NameFile string
	Discript string
}
