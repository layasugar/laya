package test

import (
	"fmt"
	"github.com/layasugar/laya/gots"
	"github.com/layasugar/laya/gots/queries"
	"testing"
)

type myObject struct {
	Id       int64  `json:"pk1"`
	ComMet   string `json:"commet"`
	Garden   int8   `json:"garden"`
	Number   int32  `json:"number"`
	Phone    string `json:"phone"`
	Username string `json:"username"`
}

var Client = gots.NewClient(
	gots.SetEndPoint("xxxx"),
	gots.SetInstanceName("xxxx"),
	gots.SetAKI("xxx"),
	gots.SetAKS("xxxx"),
)

func TestDeleteRequest_Do(t *testing.T) {
	fmt.Println("开始测试宽表删除")
	resp, err := Client.Table("user").Delete(3).Do()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("删除行数: %d行\r\n", *resp)
}

func TestGetRequest_Do(t *testing.T) {
	fmt.Println("开始测试主键获取一条")
	resp, err := Client.Table("user").Get(1).Do()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if resp != nil {
		for _, v := range resp {
			fmt.Printf("%s: %v    ", v.ColumnName, v.Value)
		}
		fmt.Printf("\r\n")
	}
}

func TestGetRangeRowsRequest_Do(t *testing.T) {
	fmt.Println("开始测试主键范围获取")
	resp, err := Client.Table("user").GetRange(1, 130).Do()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("next_pk: %v\r\n", resp.NextPk)
	for _, v := range resp.Columns {
		for _, va := range v {
			fmt.Printf("%s: %v    ", va.ColumnName, va.Value)
		}
		fmt.Printf("\r\n")
	}
}

func TestGetRowsRequest_Do(t *testing.T) {
	fmt.Println("开始测试主键多条获取")
	resp, err := Client.Table("user").GetRows([]int64{1, 2, 3, 120}).Do()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range resp {
		for _, va := range v {
			fmt.Printf("%s: %v    ", va.ColumnName, va.Value)
		}
		fmt.Printf("\r\n")
	}
}

func TestInsertRequest_Do(t *testing.T) {
	fmt.Println("开始测试结构体添加")
	h := myObject{3, "测试", 20, 12, "123456789101", "李四2"}
	resp, err := Client.Table("user").Insert(&h).Do()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("增加行数: %d行\r\n", *resp)

	fmt.Println("开始测试map添加")
	resp, err = Client.Table("user").Insert(map[string]interface{}{
		"pk1":      3,
		"commet":   "测试map添加",
		"garden":   21,
		"number":   21,
		"phone":    "3213854513",
		"username": "李四3",
	}).Do()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("增加行数: %d行\r\n", *resp)
}

func TestUpdateRequest_Do(t *testing.T) {
	fmt.Println("开始测试结构体修改一条")
	h := myObject{Id: 120, ComMet: "结构体测试修改1条"}
	count, err := Client.Table("dev_order").Update(&h).Do()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("修改行数: %d行\r\n", *count)

	fmt.Println("开始测试map修改一行")
	count, err = Client.Table("dev_order").Update(map[string]interface{}{"commet": "测试修改", "pk1": 1}).Do()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("修改行数: %d行\r\n", *count)
}

func TestWriteRowsRequest_Do(t *testing.T) {
	fmt.Println("开始测试批量修改，插入多条map")
	var puts1 = []map[string]interface{}{
		{"commet": "测试批量添加", "garden": 20, "number": 12, "phone": "123456789101", "username": "李四1101"},
		{"commet": "测试批量添加", "garden": 20, "number": 12, "phone": "123456789101", "username": "李四1102"},
		{"commet": "测试批量添加", "garden": 20, "number": 12, "phone": "123456789101", "username": "李四1103"},
		{"commet": "测试批量添加", "garden": 20, "number": 12, "phone": "123456789101", "username": "李四1104"},
	}
	var putsPks = []int64{1101, 1102, 1103, 1104}

	var updates1 = []map[string]interface{}{
		{"commet": "测试批量修改", "garden": 18},
		{"commet": "测试批量修改", "garden": 18},
	}
	var updatePks = []int64{124, 125}
	resp, err := Client.Table("user").WriteRows().SetPutData(puts1, putsPks).SetUpdateData(updates1, updatePks).Do()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(*resp)

	fmt.Println("开始测试批量修改，删除，插入多条struct")
	var updates = []myObject{
		{Id: 126, ComMet: "测试批量修改"},
		{Id: 127, ComMet: "测试批量修改"},
	}
	var puts = []myObject{
		{1010, "测试批量添加", 20, 12, "123456789101", "李四2"},
		{1011, "测试批量添加", 20, 12, "123456789101", "李四2"},
		{1012, "测试批量添加", 20, 12, "123456789101", "李四2"},
		{1013, "测试批量添加", 20, 12, "123456789101", "李四2"},
	}
	resp, err = Client.Table("user").WriteRows().
		SetPutData(&puts).
		SetDelData([]int64{1, 3}).
		SetUpdateData(&updates).Do()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(*resp)
}

// 单表各种条件查询
func TestSearchRequest_Do(t *testing.T) {
	resp1, err := Client.Table("user").Search().OrderBy("created_at asc,ots_pk_sort").Limit(5).Do()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	resp2, err := Client.Table("user").Search("user_index").OrderBy("created_at asc,ots_pk_sort").Limit(5).Do()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	q1 := queries.Not(queries.TermQuery("username", "tom"))
	q2 := queries.And(queries.TermsQuery("age", 10, 12, 13), queries.RangeQuery("age", ">", 15))
	q3 := queries.Or(q1, q2)
	resp, err := Client.Table("user").Search().Query(q3).Do()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("IsAllSuccess: ", resp.IsAllSuccess) //查看返回结果是否完整。
	fmt.Println("TotalCount: ", resp.TotalCount)     //匹配的总行数。
	fmt.Println("RowCount: ", len(resp.Rows))        //返回的行数。
	fmt.Println(resp)
	fmt.Println(resp1)
	fmt.Println(resp2)
}

func TestSearchRequest_Do1(t *testing.T) {
	q1 := queries.TermQuery("student_id", 862354)
	resp, err := Client.Table("order").Search().Query(q1).Do()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("IsAllSuccess: ", resp.IsAllSuccess) //查看返回结果是否完整。
	fmt.Println("TotalCount: ", resp.TotalCount)     //匹配的总行数。
	fmt.Println("RowCount: ", len(resp.Rows))        //返回的行数。
	fmt.Println(resp)
}
