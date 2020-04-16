package validator

import (
	"fmt"
	//ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/firmeve/firmeve/kernel/contract"
	testing2 "github.com/firmeve/firmeve/testing"
	ut "github.com/go-playground/universal-translator"
	//"github.com/go-playground/validator/v10"
	"testing"
)

func app() contract.Application {
	firmeve := testing2.TestingModeFirmeve()
	firmeve.Register(new(Provider), true)
	firmeve.Boot()
	return firmeve
}

func TestValidator_RegisterValidation(t *testing.T) {
	app := app()

	v := app.Get(`validator`).(*Validator)
	err := v.RegisterTranslationValidation(`mobile`, func(fl validator.FieldLevel) bool {
		return true
	}, func(ut ut.Translator) error {
		return ut.Add(`mobile`, `手机号格式不正确`, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(`mobile`, fe.Field())
		return t
	})
	assert.Nil(t, err)

	fmt.Printf("%v", err)
}

// User contains user information
type User struct {
	FirstName      string     `validate:"required"`
	LastName       string     `validate:"required"`
	Age            uint8      `validate:"gte=0,lte=130" alias:"中文" json:"age"`
	Email          string     `validate:"required,email"`
	FavouriteColor string     `validate:"iscolor"`                // alias for 'hexcolor|rgb|rgba|hsl|hsla'
	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
}

// Address houses a users address information
type Address struct {
	Street string `validate:"required"`
	City   string `validate:"required"`
	Planet string `validate:"required"`
	Phone  string `validate:"required"`
}

func TestValidator_Validate(t *testing.T) {
	app := app()

	validator2 := app.Resolve(`validator`).(contract.Validator)

	address := &Address{
		Street: "Eavesdown Docks",
		Planet: "Persphone",
		Phone:  "none",
	}

	user := &User{
		FirstName:      "Badger",
		LastName:       "Smith",
		Age:            135,
		Email:          "Badger.Smith@gmail.com",
		FavouriteColor: "#000-",
		Addresses:      []*Address{address},
	}

	err := validator2.Validate(user)
	assert.NotNil(t, err)
	fmt.Printf("%v", err.(*validateError).errors)
}
