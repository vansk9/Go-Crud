package web

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/pkg/errors"
)

func registerTranslator() (*validator.Validate, error) {
	v := validator.New()

	english := en.New()
	translator = ut.New(english, english)

	if err := registerValidation(v); err != nil {
		return nil, err
	}

	translatorEn, ok := translator.GetTranslator("en")
	if !ok {
		return nil, errors.New("cannot found message translator")
	}

	if err := enTranslations.RegisterDefaultTranslations(v, translatorEn); nil != err {
		return nil, errors.Wrap(err, "registering default translation")
	}

	if err := registerTranslation(v, translatorEn); err != nil {
		return nil, err
	}

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v, nil
}

var (
	translator *ut.UniversalTranslator
	validate   *validator.Validate
)

func init() {
	var err error
	validate, err = registerTranslator()
	if nil != err {
		panic(err)
	}
}

func Validator() *validator.Validate {
	return validate
}

func registerValidation(v *validator.Validate) error {

	return nil
}

func registerTranslation(v *validator.Validate, translatorEn ut.Translator) error {

	return nil
}
