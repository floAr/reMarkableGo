package rmtemplates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"muzzammil.xyz/jsonc"
)

type Template struct {
	Name       string   `json:"name"`
	Filename   string   `json:"filename"`
	IconCode   string   `json:"iconCode"`
	Categories []string `json:"categories"`
	Landscape  bool     `json:"landscape,omitempty"`
}

type TemplatesMaster struct {
	Templates []Template `json:"templates"`
}

func (tm TemplatesMaster) HasTemplateWithName(name string) bool {
	for i := 0; i < len(tm.Templates); i++ {
		if tm.Templates[i].Name == name {
			return true
		}
	}
	return false
}

func (tm TemplatesMaster) HasTemplateForFile(filename string) bool {
	var extension = filepath.Ext(filename)
	filename = filename[0 : len(filename)-len(extension)]
	for i := 0; i < len(tm.Templates); i++ {
		if tm.Templates[i].Filename == filename {
			return true
		}
	}
	return false
}

func (tm TemplatesMaster) Append(newTemplate Template) TemplatesMaster {
	tm.Templates = append(tm.Templates, newTemplate)
	return tm
}

func LoadTemplateMaster(templatesFile string) TemplatesMaster {
	data, err := ioutil.ReadFile(templatesFile)
	if err != nil {
		fmt.Print(err)
	}
	jc := jsonc.ToJSON(data) // Calling jsonc.ToJSON() to convert JSONC to JSON
	var tMaster TemplatesMaster
	err = json.Unmarshal([]byte(jc), &tMaster)
	if err != nil {
		fmt.Println("error:", err)
	}
	return tMaster
}
