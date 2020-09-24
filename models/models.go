package models

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"json2proto/tools"
	"os"
	"reflect"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type JsonToProtoOpts struct {
	MicroBasePath string
	JsonPath      string
	ProtoPath     string
}

func Json2Proto(m *JsonToProtoOpts) error {
	//检查json文件
	_, err := os.Stat(m.JsonPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("JsonPath does not exist.")
			return errors.New("JsonPath does not exist.")
		} else {
			fmt.Println("os.Stat fail " + err.Error())
			return err
		}
	}

	//读json文件
	jFile, err := os.Open(m.JsonPath)
	if err != nil {
		panic(err)
	}
	defer jFile.Close()
	data, err := ioutil.ReadAll(jFile)
	dataStr := strings.Replace(string(data), "\n", "", -1)

	//unmarshall
	var jsonMap = make(map[string]interface{})
	if err = json.Unmarshal([]byte(dataStr), &jsonMap); err != nil {
		fmt.Printf("Unmarshal fail %v", err)
		return err
	}

	//解析json
	var protoContent []string
	protoContent = append(protoContent, "message toProto {\n")
	AnalysisJson(jsonMap, &protoContent, 1)
	protoContent = append(protoContent, "}\n")
	//解析json
	//protoContent := AnalysisJson(dataStr)

	/*write proto*/
	//默认当前文件
	if m.ProtoPath == "" {
		m.ProtoPath = "toProto.proto"
	}

	pFile, err := os.OpenFile(m.ProtoPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666) //文件尾部追加和不存在则创建
	if err != nil {
		fmt.Printf("proto文件打开失败%v, 文件路径:[%s]", err, m.ProtoPath)
		return err
	}
	defer pFile.Close()

	//写入proto文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(pFile)
	for _, c := range protoContent {
		write.WriteString(c)
	}

	//Flush将缓存的文件真正写入到proto文件中
	err = write.Flush()
	if err != nil {
		fmt.Printf("flush fail %v", err)
		return err
	}

	fmt.Printf("write in proto file:[%s]\n", m.ProtoPath)
	return nil
}

//解析jsonmap
func AnalysisJson(jsonMap map[string]interface{}, protoContent *[]string, tabsCount int) {
	var num int
	var tmp string
	var dataType string
	var conType string
	//var tabsCount int

	var tabs string
	for i := 0; i < tabsCount; i++ {
		tabs += "  "
	}

	for name, v := range jsonMap {
		//编辑必要信息
		num++
		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			dataType = "string"
			break

		case reflect.Int:
			dataType = "uint32"
			break

		case reflect.Map:
			conType = "message"
			dataType = tools.StrFirstToUpper(name)
			break

		case reflect.Slice:
			conType = "repeated"
			dataType = tools.StrFirstToUpper(name)
			break
		}

		//拼写语句
		switch conType {
		case "message":
			conType = ""
			*protoContent = append(*protoContent, tabs+"message "+dataType+" {\n")
			AnalysisJson(v.(map[string]interface{}), protoContent, tabsCount+1)
			*protoContent = append(*protoContent, tabs+"}\n")
			tmp = fmt.Sprintf("%s%s %s = %d [json_name = \"%s\"]; \n", tabs, dataType, name, num, name)
			break

		case "repeated":
			conType = ""
			*protoContent = append(*protoContent, tabs+"message "+dataType+" {\n")
			val := reflect.ValueOf(v)
			for i := 0; i < val.Len(); i++ {
				AnalysisJson(val.Index(i).Interface().(map[string]interface{}), protoContent, tabsCount+1)
			}
			*protoContent = append(*protoContent, tabs+"}\n")
			tmp = fmt.Sprintf("%srepeated %s %s = %d [json_name = \"%s\"];\n", tabs, dataType, name, num, name)
			break

		default:
			tmp = fmt.Sprintf("%s%s %s = %d [json_name = \"%s\"]; \n", tabs, dataType, name, num, name)
		}

		*protoContent = append(*protoContent, tmp)
	}

	return
}
