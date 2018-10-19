package helpers

import (
	"encoding/json"
	// "fmt"
	"time"
	// "github.com/gin-gonic/gin"
)

type EveHelper struct{}

var start time.Time

func (ctrl EveHelper) Panics(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func (ctrl EveHelper) GetTime() string {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	return now.String()
}

//Create the desc & merge with original json. Code adapted from --> https://www.sohamkamani.com/blog/2017/10/18/parsing-json-in-golang/
func (ctrl EveHelper) GenerateJsonDesc(oldJson string, status bool) ([]byte, error) {
	type gen struct {
		ListData  string `json:"listdata"`
		GenTime   string `json:"generatetime"`
		GenStatus bool   `json:"generatestatus"`
	}
	desc_value := &gen{GenTime: ctrl.GetTime(), GenStatus: status}
	temp, newJson := map[string]interface{}{}, map[string]interface{}{}
	json.Unmarshal([]byte(oldJson), &temp)
	newJson["listdata"] = temp
	newJson["generatetime"] = desc_value.GenTime
	newJson["generatestatus"] = desc_value.GenStatus
	return json.Marshal(newJson)
}
