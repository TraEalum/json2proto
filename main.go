package main

import (
	"fmt"
	models "json2proto/models"
	"os"

	"github.com/droundy/goopt"
)

var (
	jsonPath  = goopt.String([]string{"-j", "--json path"}, "", "json 文件路径")
	protoPath = goopt.String([]string{"-p", "--proto path"}, "", "生成的proto 文件路径")
)

func init() {
	// Setup goopts
	goopt.Description = func() string {
		return "json to proto."
	}
	goopt.Version = "2.0"
	goopt.Summary = `toproto -s servicename`

	//Parse options
	goopt.Parse(nil)
}

func main() {
	//获取GO_MICRO_BASE
	var err error
	var autogenOpts *models.JsonToProtoOpts

	microBasePath := os.Getenv("GO_MICRO_BASE")
	if microBasePath == "" {
		fmt.Println("Error! Please config your $GO_MICRO_BASE first!")
		goto fail
	}

	if microBasePath[len(microBasePath)-1:] != "/" {
		microBasePath = microBasePath + "/"
	}

	//检查参数
	if *jsonPath == "" {
		fmt.Println("jsonFile can not be null")
		goto fail
	}

	autogenOpts = &models.JsonToProtoOpts{
		MicroBasePath: microBasePath,
		JsonPath:      *jsonPath,
		ProtoPath:     *protoPath,
	}

	//根据json生成对应proto
	err = models.Json2Proto(autogenOpts)
	if err != nil {
		fmt.Println("Json2Proto failed. " + err.Error())
		goto fail
	}

	fmt.Println("JsonToProto success.")
	return

fail:
	fmt.Println("JsonToProto failed.")
}
