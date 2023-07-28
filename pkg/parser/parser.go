package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type FlowTrackingData struct {
	APTracking   map[string]int `json:"ap_tracking"`
	ReferenceIDS map[string]int `json:"reference_ids"`
}

type Instructions struct {
	EndCol    int               `json:"end_col"`
	EndLine   int               `json:"end_line"`
	InputFile map[string]string `json:"input_file"`
	StartCol  int               `json:"start_col"`
	StartLine int               `json:"start_line"`
}

type InstructionLocation struct {
	AccessibleScopes []string         `json:"accessible_scopes"`
	FlowTrackingData FlowTrackingData `json:"flow_tracking_data"`
	Hints            []string         `json:"hints"`
	Inst             Instructions     `json:"inst"`
}

type DebugInfo struct {
	FileContents        map[string]string              `json:"file_contents"`
	InstructionLocation map[string]InstructionLocation `json:"instruction_locations"`
}

type Identifier struct {
	FullName   string         `json:"full_name"`
	Members    map[string]any `json:"members"`
	Size       int            `json:"size"`
	Decorators []string       `json:"decorators"`
	PC         int            `json:"pc"`
	Type       string         `json:"type"`
	CairoType  string         `json:"cairo_type"`
	Value      int            `json:"value"`
}

type ApTrackingData struct {
	Group  int `json:"group"`
	Offset int `json:"offset"`
}

type Reference struct {
	ApTrackingData ApTrackingData `json:"ap_tracking_data"`
	Pc             int            `json:"pc"`
	Value          string         `json:"value"`
}

type ReferenceManager struct {
	References []Reference `json:"references"`
}

type CompiledJson struct {
	Attributes       []string              `json:"attributes"`
	Builtins         []string              `json:"builtins"`
	CompilerVersion  string                `json:"compiler_version"`
	Data             []string              `json:"data"`
	DebugInfo        DebugInfo             `json:"debug_info"`
	Hints            map[string]string     `json:"hints"`
	Identifiers      map[string]Identifier `json:"identifiers"`
	MainScope        string                `json:"main_scope"`
	Prime            string                `json:"prime"`
	ReferenceManager ReferenceManager      `json:"reference_manager"`
}

func Parse(jsonPath string) CompiledJson {
	jsonFile, err := os.Open(jsonPath)

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	var cJson CompiledJson

	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &cJson)

	if err != nil {
		fmt.Println(err)
	}

	return cJson

}
