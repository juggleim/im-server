package hbtest

import (
	"context"
	"fmt"

	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
)

var client gohbase.Client
var tableName string = "jim_hismsgs"

func init() {
	client = gohbase.NewClient("127.0.0.1:2181")
}

func TestHbase() {
	// scan, err := hrpc.NewScanStr(context.Background(), tableName, hrpc.NumberOfRows(2))
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}

func GetData() {
	get, err := hrpc.NewGetStr(context.Background(), tableName, "testrow")
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := client.Get(get)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, cell := range result.Cells {
		fmt.Println(string(cell.Qualifier), string(cell.Value))
	}
}

func PutData() {
	//put data
	rowKey := "testrow"
	values := make(map[string]map[string][]byte)
	values["jim"] = make(map[string][]byte)
	values["jim"]["123"] = []byte("value1")
	putRequest, err := hrpc.NewPutStr(context.Background(), tableName, rowKey, values)
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := client.Put(putRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result, result.String())

}

func CreateTable(tname string) {

	adminClient := gohbase.NewAdminClient("127.0.0.1:2181")

	fmt.Println(adminClient)

	crt := hrpc.NewCreateTable(context.Background(), []byte(tname), map[string]map[string]string{
		"jim": {
			"MIN_VERSIONS": "1",
		},
	})
	err := adminClient.CreateTable(crt)
	fmt.Println(err)
}
