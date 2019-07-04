package config

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

var (
	directory = "../testdata/config"
	wg        sync.WaitGroup
)

func TestNewConfig(t *testing.T) {

	config, err := NewConfig(directory)
	assert.Nil(t, err)

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(directory string) {
			_, err := NewConfig(directory)
			assert.Nil(t, err)
			wg.Done()
		}(directory)
	}

	wg.Wait()

	// 单例测试
	config2, err := NewConfig(directory)
	assert.Nil(t, err)

	assert.Equal(t, config, config2)
}

func TestConfig_Set(t *testing.T) {

	config, err := NewConfig(directory)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		t.Fail()
	}

	tests := []struct {
		File  string
		Key   string
		Value string
	}{
		{"app", "x", "x"},
		{"app", "s1.x", "s1x"},
		{"app", "s1.z.y", "s1xy"},
		{"new", "x", "x"},
		{strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int()), "x", "x"},
	}
	fmt.Printf("out,%p",config)

	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(config *Config) {
			testn := make([]struct {
				File  string
				Key   string
				Value string
			},5)
			copy(testn,tests)

			config, err := NewConfig(directory)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
				t.Fail()
			}
			fmt.Printf("in,%p\n",config)
			for _, test := range testn {
				//fmt.Println(test.file,test.key,test.value,i)
				err = config.Set(test.File+"."+test.Key, test.Value)
				assert.Nil(t, err)
			}

			wg.Done()
		}(config)
	}

	wg.Wait()
}

// 正常数据测试
func TestConfig_Get(t *testing.T) {
	config, err := NewConfig(directory)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		t.Fail()
	}

	tests := []struct {
		file  string
		key   string
		value string
	}{
		{"app", "t_key", "t_value"},
		{"app", "t1.t2", "t2_value"},
		{"new", "nt.nt2.nt3", "nt3_value"},
	}

	for _, test := range tests {
		err := config.Set(test.file+"."+test.key, test.value)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			t.Fail()
		}
	}

	// 当一个值的时候，返回ini.File完整对象
	res, err := config.Get(`app`)
	if _, ok := res.(*ini.File); !ok {
		fmt.Printf("%s\n", err.Error())
		t.Fail()
	}

	// 当是2个的时候，返回默认section的key值
	res1, err := config.Get(`app.t_key`)
	fmt.Printf("\n%s\n", res1.(*ini.Key).Value())
	assert.Equal(t, `t_value`, res1.(*ini.Key).Value())

	// 当是2个的，Key不存在时
	_, err = config.Get(`app.ssssss`)
	if err == nil {
		t.Fail()
	} else if v, ok := err.(*FormatError); !ok {
		fmt.Printf("fail: error is %T", v)
		t.Fail()
	}

	// 当是3个值的时候，返回指定section的key值
	res2, err := config.Get(`app.t1.t2`)
	fmt.Printf("\n%s\n", res2.(*ini.Key).Value())
	assert.Equal(t, `t2_value`, res2.(*ini.Key).Value())

	// 当是4个以上值的时候，返回指定section子section的key值
	res3, err := config.Get(`new.nt.nt2.nt3`)
	fmt.Printf("\n%s\n", res3.(*ini.Key).Value())
	assert.Equal(t, `nt3_value`, res3.(*ini.Key).Value())

	// 优先查找section的拼接
	res4, err := config.Get(`new.nt.nt2`)
	fmt.Printf("\n%T\n", res4)
	if _, ok := res4.(*ini.Section); !ok {
		t.Fail()
	}
}

// default 非正常数据测试
//func TestConfig_GetDefault(t *testing.T) {
//	config, err := NewConfig(directory)
//	if err != nil {
//		fmt.Printf("%s\n", err.Error())
//		t.Fail()
//	}
//
//	// 当一个值的时候，返回ini.File完整对象
//	res := config.GetDefault(`ssss`, `def`)
//
//	assert.Equal(t, `def`, res.(string))
//
//	// 当一个值的时候，返回ini.File完整对象
//	res = config.GetDefault(`test.test`, 0)
//	resInt, _ := res.(*ini.Key).Int()
//	assert.Equal(t, 2, resInt)
//}

func TestConfig_All(t *testing.T) {
	config, err := NewConfig(directory)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		t.Fail()
	}

	for _, v := range config.All() {
		fmt.Printf("%v", v)
	}
}

func TestConfig_Set_Key_Error(t *testing.T) {
	config, err := NewConfig(directory)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		t.Fail()
	}

	err = config.Set("app", "123")
	if err == nil {
		t.Fail()
	}
}

func TestFormatError_Error(t *testing.T) {
	err := &FormatError{message: "abcdef"}
	assert.Equal(t, "abcdef", err.Error())
}

// ======================== Example ======================

func ExampleConfig_Set() {
	config, err := NewConfig(directory)
	if err != nil {
		panic(err.Error())
	}

	err = config.Set("test.a.b.c", "1")
	if err != nil {
		panic(err.Error())
	}
}

func ExampleConfig_Get() {
	config, err := NewConfig(directory)
	if err != nil {
		panic(err.Error())
	}

	value, err := config.Get("app.x")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(value.(*ini.Key).Value())

	// Output:
	// x
}

//
//func ExampleConfig_GetDefault() {
//	config, err := NewConfig(directory)
//	if err != nil {
//		panic(err.Error())
//	}
//
//	value1 := config.GetDefault("app.zzzz", `def`)
//	value2 := config.GetDefault("app.t_key", `def2`)
//
//	fmt.Println(value1, value2)
//
//	// Output:
//	// def t_value
//}

func ExampleConfig_All() {
	config, err := NewConfig(directory)
	if err != nil {
		panic(err.Error())
	}

	configs := config.All()

	fmt.Printf("%T", configs)

	// Output:
	// map[string]*ini.File
}
