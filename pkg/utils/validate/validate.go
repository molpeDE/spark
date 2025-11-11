package validate

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	go_validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"golang.org/x/exp/maps"
)

var validator = go_validator.New()
var eng = en.New()
var uni = ut.New(eng, eng)
var trans, _ = uni.GetTranslator("en")

func init() {
	en_translations.RegisterDefaultTranslations(validator, trans)
}

func ValidateStruct(theStruct any) error {
	if err := validator.Struct(theStruct); err != nil {
		if validateErr, ok := err.(go_validator.ValidationErrors); ok {
			return fmt.Errorf("%s", strings.Join(maps.Values(validateErr.Translate(trans)), "\n"))
		} else if _, ok := err.(*go_validator.InvalidValidationError); ok {
			return nil // arg wasnt a struct, lgtm
		} else {
			log.Println(err)
			return err
		}
	}

	return nil
}
