package tools

import (
	"github.com/pkg/errors"
	"reflect"
)

func SetField(obj interface{}, fieldName string, fieldValue interface{}) (bool, error) {
	// 获取结构体的反射值和类型
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()
	// 遍历结构体字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == fieldName {
			// 获取字段的反射值
			fv := v.Field(i)
			if !fv.IsValid() || !fv.CanSet() {
				return false, errors.Errorf("cannot set field %s", fieldName)
			}
			// 将值 fieldValue 设置到字段中
			cVal := reflect.ValueOf(fieldValue)
			if fv.Type() != cVal.Type() {
				return false, errors.Errorf("type mismatch: expected %v, got %v", fv.Type(), cVal.Type())
			}
			fv.Set(cVal)
			return true, nil
		}
	}
	return false, nil
}
