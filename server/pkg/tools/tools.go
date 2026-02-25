package tools

import (
	"biz-auto-api/pkg/consts"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
	"unicode"
)

// 生成密码
func GetEncryptedPassword(str string) string {
	str = strings.ToUpper(consts.Salt) + str + consts.Salt
	return MD5Str(str)
}

func MD5Str(src string) string {
	h := md5.New()
	h.Write([]byte(src)) // 需要加密的字符串为
	return hex.EncodeToString(h.Sum(nil))
}

func FuzzyQuery(search string) string {
	return fmt.Sprintf("%%%v%%", search)
}

func StructToMap(s interface{}, tag string, includeZeroValue bool, excludeFields ...string) map[string]interface{} {
	if len(strings.TrimSpace(tag)) == 0 {
		tag = "map"
	}
	m := make(map[string]interface{})
	val := reflect.ValueOf(s)
	typ := reflect.TypeOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		// 如果字段在需要排除的字段中则跳过
		if InSlice(field.Name, excludeFields) {
			continue
		}
		if !field.IsExported() {
			continue
		}
		fieldVal := val.Field(i).Interface()

		dbTag := field.Tag.Get(tag)
		if dbTag == "" {
			dbTag = CamelToSnake(field.Name)
		}
		if includeZeroValue || !reflect.DeepEqual(fieldVal, reflect.Zero(field.Type).Interface()) {
			m[dbTag] = fieldVal
		}
	}
	return m
}

func CamelToSnake(camelCase string) string {
	var snakeCase bytes.Buffer
	for i, r := range camelCase {
		if unicode.IsUpper(r) {
			if i > 0 {
				snakeCase.WriteRune('_')
			}
			snakeCase.WriteRune(unicode.ToLower(r))
		} else {
			snakeCase.WriteRune(r)
		}
	}
	return snakeCase.String()
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func TimeFormatStr(UTCtimeStr string) string {
	if len(UTCtimeStr) == 0 {
		return UTCtimeStr
	}
	t, err := time.Parse(time.RFC3339, UTCtimeStr)
	if err != nil {
		return UTCtimeStr
	}
	return t.Format("2006-01-02 15:04:05")
}

func BoolToString(b bool) string {
	if b {
		return "是"
	}
	return "否"
}

func HasElem(elem string, elems []string) bool {
	for _, existElem := range elems {
		if existElem == elem {
			return true
		}
	}
	return false
}
func Retry[R any](times int, sleep time.Duration, fn func() (R, error)) (res R, err error) {
	if times <= 0 {
		times = 1
	}
	for i := 1; i <= times; i++ {
		res, err = fn()
		// 没报错直接返回结果
		if err == nil {
			return
		}
		// 最后一次
		if i >= times {
			break
		}
		time.Sleep(sleep * time.Duration(i))
	}
	return
}
