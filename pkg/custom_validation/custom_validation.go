package custom_validation

import (
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"

	Config "mini-accounting/config"
	Constants "mini-accounting/constants"
	Library "mini-accounting/library"
)

type CustomValidation interface {
	GetValidator() *validator.Validate
	ValidateStruct(st interface{}, tagName string) []map[string]interface{}
	GetCustomErrorMessage(validationErrors validator.ValidationErrors, callback func(anyParam interface{}) interface{}, fields ...string) []map[string]interface{}
	GetCustomFieldName(validationErrors validator.ValidationErrors, tag string, getField func(field string) *reflect.StructField) *[]string
	ConvertStructToInterfaceFields(orig interface{}) interface{}
}

type CustomValidationImpl struct {
	config    Config.Config
	library   Library.Library
	validator *validator.Validate
	st        interface{}
}

func NewCustomValidation(
	config Config.Config,
	library Library.Library,
) CustomValidation {
	instance := CustomValidationImpl{
		config:    config,
		library:   library,
		validator: validator.New(),
	}

	return &instance
}

func (u *CustomValidationImpl) GetValidator() *validator.Validate {
	u.validator.RegisterValidation(Constants.ValidationDateLTETodayIfNotNull, u.DateLTETodayIfNotNull)
	u.validator.RegisterValidation(Constants.ValidationDateLTEParamIfNotNull, u.DateLTEParamIfNotNull)
	u.validator.RegisterValidation(Constants.ValidationDateGTEParamIfNotNull, u.DateGTEParamIfNotNull)
	u.validator.RegisterValidation(Constants.ValidationMinNumeric, u.MinNumeric)
	u.validator.RegisterValidation(Constants.ValidationMinNumericUnless, u.MinNumericUnless)
	u.validator.RegisterValidation(Constants.ValidationNIK, u.ValidateNIK)
	u.validator.RegisterValidation(Constants.ValidationCIF, u.ValidateCIF)
	u.validator.RegisterValidation(Constants.ValidationDataTypeString, u.ValidateDataTypeString)
	u.validator.RegisterValidation(Constants.ValidationDataTypeNumeric, u.ValidateDataTypeNumeric)
	u.validator.RegisterValidation(Constants.ValidationLimitMinNumeric, u.MinLimitNumeric)
	u.validator.RegisterValidation(Constants.ValidationInterfaceTypeRequired, u.ValidateInterfaceTypeRequired)
	u.validator.RegisterValidation(Constants.ValidationValueNumeric, u.ValidateValueNumeric)
	u.validator.RegisterValidation(Constants.ValidationAlphaNumeric, u.ValidationAlphaNumeric)
	u.validator.RegisterValidation(Constants.ValidationUnixTime, u.Unixtime)
	u.validator.RegisterValidation(Constants.ValidationEmailRequiredWithout, u.ValidationEmailRequiredWithout)
	u.validator.RegisterValidation(Constants.ValidationAlphanumericRequiredWithout, u.ValidationAlphaNumericRequiredWithout)
	return u.validator
}

func (u *CustomValidationImpl) ValidateStruct(st interface{}, tagName string) []map[string]interface{} {
	u.st = st
	err := u.GetValidator().Struct(st)

	if err == nil {
		return []map[string]interface{}{}
	}

	validationErrors := err.(validator.ValidationErrors)
	fields := u.GetCustomFieldName(validationErrors, tagName, func(field string) *reflect.StructField {
		fieldStruct, ok := reflect.TypeOf(st).Elem().FieldByName(field)
		if !ok {
			return nil
		}
		return &fieldStruct
	})

	return u.GetCustomErrorMessage(validationErrors, nil, *fields...)
}

func (u *CustomValidationImpl) getFieldValueByTag(v interface{}, tagName, tagValue string) (interface{}, error) {
	r := reflect.ValueOf(v)

	// Dereference pointer if necessary
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}

	t := r.Type()
	for i := 0; i < r.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get(tagName) == tagValue {
			return r.Field(i).Interface(), nil
		}
	}

	return nil, u.library.Errorf("no field found with %s='%s'", tagName, tagValue)
}

func (u *CustomValidationImpl) GetCustomFieldName(validationErrors validator.ValidationErrors, tag string, getField func(field string) *reflect.StructField) *[]string {
	var fields []string
	for _, v := range validationErrors {
		field := getField(v.Field())
		if field == nil {
			fields = append(fields, v.Field())
			continue
		}

		fieldName, ok := (*field).Tag.Lookup(tag)
		if !ok {
			fields = append(fields, v.Field())
			continue
		}

		fields = append(fields, fieldName)
	}

	return &fields
}

func (u *CustomValidationImpl) GetCustomErrorMessage(validationErrors validator.ValidationErrors, callback func(anyParam interface{}) interface{}, fields ...string) []map[string]interface{} {
	var errors []map[string]interface{}
	for i, verr := range validationErrors {
		fieldName := verr.Field()
		if len(fields) == len(validationErrors) {
			fieldName = fields[i]
		}
		switch verr.Tag() {
		case Constants.ValidationRequired, Constants.ValidationInterfaceTypeRequired:
			errors = append(errors, map[string]interface{}{
				"field":   verr.Field(),
				"message": u.library.Sprintf(Constants.ErrValidationRequired.Error(), fieldName),
			})
		case Constants.ValidationMin:
			minValue := verr.Param()
			errors = append(errors, map[string]interface{}{
				"field":   verr.Field(),
				"message": u.library.Sprintf(Constants.ErrValidationMin.Error(), fieldName, minValue),
			})
		case Constants.ValidationMax:
			errors = append(errors, map[string]interface{}{
				"field":   verr.Field(),
				"message": u.library.Sprintf(Constants.ErrValidationMax.Error(), fieldName),
			})
		case Constants.ValidationOneOf:
			errors = append(errors, map[string]interface{}{
				"field":   verr.Field(),
				"message": u.library.Sprintf(Constants.ErrValidationOneOF.Error(), fieldName),
			})
		case Constants.ValidationDataTypeString, Constants.ValidationDataTypeNumeric:
			errors = append(errors, map[string]interface{}{
				"field":   verr.Field(),
				"message": u.library.Sprintf(Constants.ErrDataTypeInvalid.Error(), fieldName),
			})
		case Constants.ValidationAlphaNumeric:
			errors = append(errors, map[string]interface{}{
				"field":   verr.Field(),
				"message": Constants.ErrAlphaNumeric.Error(),
			})
		case Constants.Email:
			errors = append(errors, map[string]interface{}{
				"field":   verr.Field(),
				"message": Constants.ErrEmail.Error(),
			})
		case Constants.ValidationEmailRequiredWithout:
			errors = append(errors, map[string]interface{}{
				"field":   verr.Field(),
				"message": Constants.ErrEmail.Error(),
			})
		case Constants.ValidationAlphanumericRequiredWithout:
			errors = append(errors, map[string]interface{}{
				"field":   verr.Field(),
				"message": Constants.ErrAlphaNumeric.Error(),
			})
		default:
			errors = append(errors, map[string]interface{}{
				"field":   verr.Field(),
				"message": Constants.ErrSomethingWentWrong.Error(),
			})
		}
	}

	return errors
}

func (u *CustomValidationImpl) ConvertStructToInterfaceFields(orig interface{}) interface{} {
	origType := u.library.ReflectTypeOf(orig)
	if origType.Kind() == reflect.Ptr {
		origType = origType.Elem()
	}

	// Define fields for the new struct
	fields := []reflect.StructField{}
	for i := 0; i < origType.NumField(); i++ {
		field := origType.Field(i)

		// Create a new StructField with type interface{}
		newField := reflect.StructField{
			Name: field.Name,
			Type: u.library.ReflectTypeOf((*interface{})(nil)).Elem(), // Set type to interface{}
			Tag:  field.Tag,                                           // Preserve the tags
		}
		fields = append(fields, newField)
	}

	// Create the new struct type
	newStructType := u.library.ReflectStructOf(fields)

	// Create a new instance of the struct
	return u.library.ReflectNew(newStructType).Interface()
}

// VALIDATE FIELD THAT CONTAINS VALUE LESS THAN OR EQUAL TO TODAY IF NOT NULL
func (u *CustomValidationImpl) DateLTETodayIfNotNull(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if field == "" {
		return true
	}

	dateField, err := u.library.ParseTime(Constants.YYYYMMDD, field)
	if err != nil {
		return false
	}

	today := u.library.GetNow().Truncate(24 * time.Hour) // Set time to the start of the day
	return dateField.Before(today) || dateField.Equal(today)
}

func (u *CustomValidationImpl) DateLTEParamIfNotNull(fl validator.FieldLevel) bool {
	param := fl.Param()
	if param == "" {
		return true
	}

	paramField, err := u.getFieldValueByTag(u.st, "name", param)
	if err != nil {
		return true
	}

	if paramField == nil {
		return true
	}

	value, ok := paramField.(string)
	if !ok {
		value = *(paramField.(*string))
	}

	dateParam, err := u.library.ParseTime(Constants.YYYYMMDD, value)
	if err != nil {
		return true
	}

	field := fl.Field().String()
	if field == "" {
		return true
	}

	dateField, err := u.library.ParseTime(Constants.YYYYMMDD, field)
	if err != nil {
		return true
	}

	return dateField.Before(dateParam) || dateField.Equal(dateParam)
}

func (u *CustomValidationImpl) DateGTEParamIfNotNull(fl validator.FieldLevel) bool {
	param := fl.Param()
	if param == "" {
		return true
	}

	paramField, err := u.getFieldValueByTag(u.st, "name", param)
	if err != nil {
		return true
	}

	if paramField == nil {
		return true
	}

	value, ok := paramField.(string)
	if !ok {
		value = *(paramField.(*string))
	}

	dateParam, err := u.library.ParseTime(Constants.YYYYMMDD, value)
	if err != nil {
		return true
	}

	field := fl.Field().String()
	if field == "" {
		return true
	}

	dateField, err := u.library.ParseTime(Constants.YYYYMMDD, field)
	if err != nil {
		return true
	}

	return dateField.After(dateParam) || dateField.Equal(dateParam)
}

func (u *CustomValidationImpl) MinNumeric(fl validator.FieldLevel) bool {
	param := fl.Param()

	if param == "" {
		return false
	}

	numericParam, err := u.library.Atoi(param)
	if err != nil {
		return false
	}

	field := fl.Field()

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(field.Int()) >= float64(numericParam)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(field.Uint()) >= float64(numericParam)
	case reflect.Float32, reflect.Float64:
		return field.Float() >= float64(numericParam)
	}

	return false
}

func (u *CustomValidationImpl) ValidateInterfaceTypeRequired(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		return field.String() != "" // Non-empty string
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return field.Int() != 0 // Non-zero integer
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return field.Uint() != 0 // Non-zero unsigned integer
	case reflect.Float32, reflect.Float64:
		return field.Float() != 0.0 // Non-zero float
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		return field.Len() > 0 // Non-empty collection
	case reflect.Ptr, reflect.Interface:
		return !field.IsNil() // Non-nil pointer or interface
	default:
		return !field.IsZero() // General non-zero value
	}
}

func (u *CustomValidationImpl) MinLimitNumeric(fl validator.FieldLevel) bool {
	param := fl.Param()

	if param == "" {
		return false
	}

	numericParam, err := u.library.Atoi(param)
	if err != nil {
		return false
	}

	field := fl.Field()

	if field.Kind() != reflect.Float64 {
		return float64(field.Int()) != 0 && float64(field.Int()) >= float64(numericParam)
	}

	return field.Float() != 0 && field.Float() >= float64(numericParam)
}

func (u *CustomValidationImpl) MinNumericUnless(fl validator.FieldLevel) bool {
	param := fl.Param()
	if param == "" {
		return false
	}

	var min, unless int
	_, err := fmt.Sscanf(param, "%d:%d", &min, &unless)
	if err != nil {
		return false
	}

	field := int(fl.Field().Int())

	return field == unless || field >= min
}

func (u *CustomValidationImpl) ValidateNIK(fl validator.FieldLevel) bool {
	param := fl.Field().String()

	regex := `^\d{16}$`

	re := regexp.MustCompile(regex)

	return re.MatchString(param)
}

func (u *CustomValidationImpl) ValidateCIF(fl validator.FieldLevel) bool {
	param := fl.Field().String()

	regex := `^\d{10}$`

	re := regexp.MustCompile(regex)

	return re.MatchString(param)
}

func (u *CustomValidationImpl) ValidateDataTypeString(fl validator.FieldLevel) bool {
	_, ok := fl.Field().Interface().(string)
	return ok
}

func (u *CustomValidationImpl) ValidateDataTypeNumeric(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, ok := fl.Field().Interface().(int)
		return ok
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		_, ok := fl.Field().Interface().(uint)
		return ok
	case reflect.Float32, reflect.Float64:
		_, ok := fl.Field().Interface().(float64)
		return ok
	}

	return false
}

func (u *CustomValidationImpl) ValidateValueNumeric(fl validator.FieldLevel) bool {
	param := fl.Field().String()

	regex := `^\d+$`

	re := regexp.MustCompile(regex)

	return re.MatchString(param)
}

func (u *CustomValidationImpl) ValidationAlphaNumeric(fl validator.FieldLevel) bool {
	param := fl.Field().String()

	regex := `^[a-zA-Z0-9]+$`

	re := regexp.MustCompile(regex)

	return re.MatchString(param)
}

func (u *CustomValidationImpl) Unixtime(fl validator.FieldLevel) bool {
	time := fl.Field().Int()
	if time <= 0 {
		return false
	}
	return true
}

func (u *CustomValidationImpl) ValidationEmailRequiredWithout(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	if email != "" {
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(emailRegex, email)
		if !matched {
			return false
		}
	}

	parent := fl.Parent()

	otherFieldName := fl.Param()

	otherField := parent.FieldByName(otherFieldName)
	if !otherField.IsValid() {
		return false
	}

	if email == "" && otherField.String() == "" {
		return false
	}

	return true
}

func (u *CustomValidationImpl) ValidationAlphaNumericRequiredWithout(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if value != "" {
		regex := `^[a-zA-Z0-9]+$`
		matched, _ := regexp.MatchString(regex, value)
		if !matched {
			return false
		}
	}

	parent := fl.Parent()

	otherFieldName := fl.Param()

	otherField := parent.FieldByName(otherFieldName)
	if !otherField.IsValid() {
		return false
	}

	if value == "" && otherField.String() == "" {
		return false
	}

	return true
}
