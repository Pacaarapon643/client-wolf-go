package util

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// สร้างตัวแปร global
var Validate = validator.New()

// สร้าง struct สำหรับเก็บ Error ที่จะส่งกลับไปหา User
type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func init() {
	// ฟังก์ชันนี้จะรันเองอัตโนมัติ 1 ครั้งเมื่อ package นี้ถูกเรียกใช้
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func ValidateStruct(s interface{}) []map[string]string {
	var errors []map[string]string
	err := Validate.Struct(s)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			// ดึงชื่อ field (ถ้าใช้ RegisterTagNameFunc จะได้เป็นตัวเล็กตาม JSON)
			field := err.Field()

			// สร้างข้อความภาษาไทยตาม Tag
			var message string
			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("กรุณากรอกข้อมูลช่อง %s", field)
			case "email":
				message = "รูปแบบอีเมลไม่ถูกต้อง"
			case "min":
				message = fmt.Sprintf("%s ต้องมีความยาวอย่างน้อย %s ตัวอักษร", field, err.Param())
			case "max":
				message = fmt.Sprintf("%s ต้องมีความยาวไม่เกิน %s ตัวอักษร", field, err.Param())
			case "gte":
				message = fmt.Sprintf("%s ต้องมีค่าตั้งแต่ %s ขึ้นไป", field, err.Param())
			default:
				message = fmt.Sprintf("ข้อมูลช่อง %s ไม่ถูกต้อง", field)
			}

			errorDetail := map[string]string{
				"field":   field,
				"message": message, // ส่ง message ภาษาไทยกลับไปแทน tag
			}
			errors = append(errors, errorDetail)
		}
	}
	return errors
}
