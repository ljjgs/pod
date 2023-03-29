package common

import (
	"reflect"
)

// 将 src 结构体对象中的字段值拷贝到 dest 结构体对象中对应的同名字段中。
func SwapTo(src interface{}, dest interface{}) error {
	srcValue := reflect.ValueOf(src)
	destValue := reflect.ValueOf(dest)

	// 判断 src 和 dest 是否为空或不可设置
	if !srcValue.IsValid() || !destValue.IsValid() || !destValue.Elem().CanSet() {
		return nil
	}

	// 获取 src 和 dest 的类型信息
	srcType := reflect.TypeOf(src).Elem()
	destType := reflect.TypeOf(dest).Elem()

	// 遍历 src 的所有字段
	for i := 0; i < srcType.NumField(); i++ {
		srcField := srcType.Field(i)
		destField, ok := destType.FieldByName(srcField.Name)
		if !ok {
			continue // 如果 dest 中不存在同名字段，则跳过
		}
		if !destField.Type.AssignableTo(srcField.Type) {
			continue // 如果 dest 中同名字段的类型不兼容，则跳过
		}
		srcFieldValue := srcValue.Elem().FieldByName(srcField.Name)
		if !srcFieldValue.IsValid() {
			continue // 如果 src 中该字段值无效，则跳过
		}
		destFieldValue := destValue.Elem().FieldByName(srcField.Name)
		if !destFieldValue.IsValid() || !destFieldValue.CanSet() {
			continue // 如果 dest 中该字段值无效或不可设置，则跳过
		}
		destFieldValue.Set(srcFieldValue) // 将 src 中该字段值拷贝到 dest 中
	}

	return nil
}
