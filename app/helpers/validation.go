package helpers

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

func init() {
	// Initialize English and Indonesian translators
	en := en.New()
	id := id.New()           // Indonesian locale
	uni = ut.New(en, en, id) // Register both English and Indonesian translators

	// Get the translator for the Indonesian language
	Trans, found := uni.GetTranslator("id")
	if !found {
		panic("Translator for 'id' not found")
	}

	// Initialize the validator and register Indonesian translations
	validate = validator.New()
	if err := id_translations.RegisterDefaultTranslations(validate, Trans); err != nil {
		panic("Error registering translations: " + err.Error())
	}
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}

// GetTranslator returns the translator instance
func GetTranslator() ut.Translator {
	Trans, found := uni.GetTranslator("id")
	if !found {
		panic("Translator for 'id' not found")
	}
	return Trans
}
func getJsonTagName(field reflect.StructField) string {
	// Get the "json" tag
	tag := field.Tag.Get("json")
	if tag == "" {
		return field.Name
	}
	return tag
}

func getAliasTagName(field reflect.StructField) string {
	// Get the "alias" tag
	alias := field.Tag.Get("alias")
	if alias == "" {
		return ""
	}
	return alias
}

func ValidateData(data interface{}) (map[string]string, error) {
	err := validate.Struct(data)
	if err != nil {
		validationErrors := make(map[string]string)

		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, err := range ve {
				// Get field by name
				field, found := reflect.TypeOf(data).Elem().FieldByName(err.Field())
				if found {
					// Get the alias or fall back to JSON tag
					alias := getAliasTagName(field)
					fieldName := alias
					if alias == "" {
						fieldName = getJsonTagName(field)
					}

					// Replace the default field name with alias in the error message
					message := err.Translate(GetTranslator())
					translatedMessage := strings.Replace(message, err.Field(), fieldName, 1)
					tagJosnName := getJsonTagName(field)
					validationErrors[tagJosnName] = translatedMessage
				}
			}
			return validationErrors, nil
		} else {
			return nil, fmt.Errorf("validation failed: %v", err)
		}
	}
	return nil, nil
}

func CheckUserName(userName string) string {

	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	phonePattern := `^\+?[0-9]{10,15}$` // Allows optional '+' and numbers between 10-15 digits

	identifier := strings.TrimSpace(userName)
	if identifier == "" {
		return "unknown"
	}

	if match, _ := regexp.MatchString(emailPattern, userName); match {
		return "email"
	}

	if match, _ := regexp.MatchString(phonePattern, userName); match {
		return "phone"
	}

	return "username"
}

func FormatPhoneNumber(phone string, wa bool) string {
	phone = strings.TrimPrefix(phone, "+")
	phone = strings.ReplaceAll(phone, "-", "")

	var phones = phone
	if wa {
		if strings.HasPrefix(phone, "08") {
			phones = "62" + strings.TrimPrefix(phone, "0")
		}
	} else {
		if strings.HasPrefix(phone, "62") {
			phones = "0" + strings.TrimPrefix(phone, "62")
		}
	}
	return phones
}

func CheckFormatPhoneNumber(phone string) bool {
	phone = strings.TrimPrefix(phone, "+")

	if strings.HasPrefix(phone, "62") {
		return true
	}

	if strings.HasPrefix(phone, "08") {
		return true
	}

	return false
}

func DecodeBase64(base64String string) (string, error) {

	decodedBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", fmt.Errorf("error decoding Base64 string: %v", err)
	}

	decodedString := string(decodedBytes)
	return decodedString, nil
}
