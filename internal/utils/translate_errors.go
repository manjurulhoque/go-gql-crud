package utils

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/manjurulhoque/go-gql-crud/internal/db"
	"github.com/manjurulhoque/go-gql-crud/pkg/dbc"
	"reflect"
	"strconv"
	"strings"
)

type FieldError struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (v *FieldError) Error() string {
	return fmt.Sprintf("Field: %s, Message: %s", v.Key, v.Value)
}

var (
	uni *ut.UniversalTranslator
	vl  *validator.Validate
)

func SliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func TranslateError(model interface{}) (errs []FieldError) {
	english := en.New()
	uni = ut.New(english, english)
	trans, _ := uni.GetTranslator("en")

	vl = validator.New()
	_ = enTranslations.RegisterDefaultTranslations(vl, trans)

	_ = vl.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = vl.RegisterTranslation("email_exists", trans, func(ut ut.Translator) error {
		return ut.Add("email_exists", "{0} is already taken", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email_exists", fe.Field())
		return t
	})

	_ = vl.RegisterTranslation("integer", trans, func(ut ut.Translator) error {
		return ut.Add("integer", "{0} must be an integer", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("integer", fe.Field())
		return t
	})

	if registerValidationError := vl.RegisterValidation("passwords_match", func(fl validator.FieldLevel) bool {
		pass1Field := fl.Parent().FieldByName("Password1")
		pass2Field := fl.Parent().FieldByName("Password2")
		return pass1Field.String() == pass2Field.String()
	}); registerValidationError != nil {
		fmt.Println("Error registering emailExists validation")
	}

	_ = vl.RegisterTranslation("passwords_match", trans, func(ut ut.Translator) error {
		return ut.Add("passwords_match", "{0} and {1} must be the same", true) // Registering the error message
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("passwords_match", "Password1", "Password2")
		return t
	})

	if registerValidationError := vl.RegisterValidation("email_exists", func(fl validator.FieldLevel) bool {
		userRepository := db.NewUserRepository(dbc.GetDB())
		user, _ := userRepository.FindUserByEmail(fl.Field().String())
		if reflect.ValueOf(user).IsZero() {
			return true
		}
		return false
	}); registerValidationError != nil {
		fmt.Println("Error registering email_exists validation")
	}

	if registerValidationError := vl.RegisterValidation("integer", func(fl validator.FieldLevel) bool {
		value, err := strconv.Atoi(fl.Field().String())
		if err != nil {
			// handle error
			fmt.Println(err)
		}
		return reflect.TypeOf(value).Kind() == reflect.Int
	}); registerValidationError != nil {
		fmt.Println("Error registering integer validation")
	}

	err := vl.Struct(model)
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)

	for _, e := range validatorErrs {
		//translatedErr := fmt.Errorf(e.Translate(trans))
		translatedErr := FieldError{
			Key:   strings.ToLower(e.Field()),
			Value: e.Translate(trans),
		}
		errs = append(errs, translatedErr)
	}
	return errs
}
