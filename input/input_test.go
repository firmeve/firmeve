package input

import (
	"testing"
)

type MockJson struct {
	Id int
	Type int `json:",int"`
	Type2 int `json:",string"`
	TypeText string `json:"type_text"`
	OrderNo string `json:"order_no"`
}

func TestNew(t *testing.T) {
	//jsonString := `{"id":4787,"type":1,"type2":"1","type_text":"新订单","order_no":"2019112914110451","quantity":3,"total_price":"1995.00","actually_price":"1995.00","order_price":"1995.00","offer_price":"0.00","pay_price":"721.00","exchange_quantity":0,"balance":"1274.00","exchange_price":"0.00","status":1,"status_text":"进行中","payment_status":1,"payment_status_text":"未付款","logistics_status":1,"logistics_status_text":"未接单","created_at":"2019-11-29 14:11:04","completed_at":"1970-01-01 08:00:00","address":{"name":"cw","mobile":"13787172457","address":"addresses"},"coupons":[1,2,3],"transaction":{"id":4728,"transaction_no":"1011175007864509939","pay_type_text":"miniprogram","status_text":"进行中","payed_at":"1970-01-01 08:00:00","trade_no":""},"products":[{"title":"products","image":"https://vgou-static.oss-cn-hangzhou.aliyuncs.com/JZvQrJ3.jpeg","single_price":"665.00","total_price":"1995.00","quantity":3,"product":{"area_price":"665.00","title":"派威动力电池72V20A","image":"https://vgou-static.oss-cn-hangzhou.aliyuncs.com/2019-06-15/tekzOLPk35BzGvhZ8dhmcLjKrknsDxEWxJZvQrJ3.jpeg"},"exchange_single_price":"0.00","exchange_total_price":"0.00","exchange_quantity":0}],"exchanges":["a1","a2","a3"],"refund":null,"user":{"id":743,"mobile":"13787111457","shop_name":"11111","nickname":"11111"}}`
	//json := parser.NewJSON([]byte(jsonString))
	//
	//input := New(json)
	//var mockJson MockJson
	//input.Bind(&mockJson)
	//fmt.Printf("%#v",mockJson)
	//fmt.Println("===========")
	//fmt.Println(input.GetFloat(`type`) == 1)
	//fmt.Println(input.GetSliceString(`exchanges`))
}