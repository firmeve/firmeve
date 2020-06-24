package binding

import (
	json2 "encoding/json"
	"errors"
	"github.com/firmeve/firmeve/kernel/contract"
	mock2 "github.com/firmeve/firmeve/testing/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

var (
	app contract.Application
)

/*
<input type="text" name="name" value="joeybloggs"/>
  <input type="text" name="age" value="3"/>
  <input type="text" name="gender" value="Male"/>
  <input type="text" name="address[0].name" value="26 Here Blvd."/>
  <input type="text" name="address[0].phone" value="9(999)999-9999"/>
  <input type="text" name="address[1].name" value="26 There Blvd."/>
  <input type="text" name="address[1].phone" value="1(111)111-1111"/>
  <input type="text" name="active" value="true"/>
  <input type="text" name="map_example[key]" value="value"/>
  <input type="text" name="nested_map[key][key]" value="value"/>
  <input type="text" name="nested_array[0][0]" value="value"/>
  <input type="submit"/>
*/

// Address contains address information
type Address struct {
	Name  string
	Phone string
}

type AddressForm struct {
	Name  string `form:"name"`
	Phone string `form:"phone"`
}

// User contains user information
type User struct {
	Name        string
	Age         uint8
	Gender      string
	Active      bool
	Address     []Address
	MapExample  map[string]string
	NestedMap   map[string]map[string]string
	NestedArray [][]string
}

func TestBindJSON(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	http := mock2.NewMockHttpProtocol(mockCtrl)
	http.EXPECT().IsMethod(`GET`).Return(false)
	http.EXPECT().ContentType().Return(`application/json`)
	user := &User{
		Name:   "hello",
		Age:    123,
		Gender: "gender",
		Active: false,
		Address: []Address{{
			Name:  "address",
			Phone: "18712345234",
		}},
		MapExample: map[string]string{"key": "value"},
		NestedMap:  map[string]map[string]string{"nested": {"n_key": "n_value"}},
		NestedArray: [][]string{
			{"a", "b"},
			{"d", "e"},
		},
	}
	v, err := json2.Marshal(user)
	assert.Nil(t, err)
	http.EXPECT().Message().Return(v, nil)

	user2 := new(User)
	err = Bind(http, user2)
	assert.Nil(t, err)
	//fmt.Printf("%#v", user2)
	assert.Equal(t, *user, *user2)
}

func TestBindQuery(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	http := mock2.NewMockHttpProtocol(mockCtrl)
	http.EXPECT().IsMethod(`GET`).Return(true)
	http.EXPECT().Values().Return(map[string][]string{
		"name":  {"simon"},
		"phone": {"123"},
	})
	address := new(AddressForm)
	err := Bind(http, address)
	assert.Nil(t, err)
	assert.Equal(t, "simon", address.Name)
	assert.Equal(t, "123", address.Phone)
}

func TestBindForm(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	http := mock2.NewMockHttpProtocol(mockCtrl)
	http.EXPECT().IsMethod(`GET`).Return(false)
	http.EXPECT().ContentType().Return(contract.HttpMimeForm)
	form := url.Values{
		"Name":                     []string{"hello"},
		"Age":                      []string{"123"},
		"Gender":                   []string{"gender"},
		"Address[0].Name":          []string{"address"},
		"Address[0].Phone":         []string{"18712345234"},
		"active":                   []string{"false"},
		"MapExample[key]":          []string{"value"},
		"NestedMap[nested][n_key]": []string{"n_value"},
		"NestedArray[0][0]":        []string{"a"},
		"NestedArray[0][1]":        []string{"b"},
		"NestedArray[1][0]":        []string{"d"},
		"NestedArray[1][1]":        []string{"e"},
	}
	gomock.InOrder(
		//http.EXPECT().Values().Return(form),
		http.EXPECT().Values().Return(form),
	)
	//http.EXPECT().Values().Return(form)

	user := &User{
		Name:   "hello",
		Age:    123,
		Gender: "gender",
		Active: false,
		Address: []Address{{
			Name:  "address",
			Phone: "18712345234",
		}},
		MapExample: map[string]string{"key": "value"},
		NestedMap:  map[string]map[string]string{"nested": {"n_key": "n_value"}},
		NestedArray: [][]string{
			{"a", "b"},
			{"d", "e"},
		},
	}

	user2 := new(User)
	err := Bind(http, user2)
	assert.Nil(t, err)
	assert.Equal(t, *user, *user2)
}

func TestFormError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	protocol := mock2.NewMockProtocol(mockCtrl)
	err := Form.Protocol(protocol, nil)
	assert.Equal(t, true, errors.Is(ProtocolTypeError, err))
}
