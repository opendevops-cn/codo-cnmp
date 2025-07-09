package utils

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct{}

func (conf *Config) LoadEnv() {
	var (
		field,
		levelField reflect.Value
		fieldInfo,
		levelFieldInfo reflect.StructField
	)

	v := reflect.ValueOf(conf).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo = v.Type().Field(i)
		if fieldInfo.Type.Kind() == reflect.Struct {
			stf := v.Field(i).NumField()

			for ii := 0; ii < stf; ii++ {
				levelFieldInfo = v.Field(i).Type().Field(ii)
				levelField = v.Field(i).Field(ii)
				setLevelField(&levelFieldInfo, &levelField)
			}
			continue
		}
		field = v.Field(i)
		setLevelField(&fieldInfo, &field)
	}
}

func StrNotIn(str string, list []string) bool {
	for _, v := range list {
		if str == v {
			return false
		}
	}
	return true
}

func StrToInt(str string, defaultVal int) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return defaultVal
	}
	return val
}

func StrToBool(str string) bool {
	str = strings.TrimSpace(str)
	strLower := strings.ToLower(str)
	if strLower == "true" || strLower == "yes" {
		return true
	}
	if strLower == "false" || strLower == "no" {
		return false
	}
	if v, err := strconv.Atoi(str); err == nil {
		return v > 0
	}
	return StrNotIn(str, []string{"", "0"})
}

func setLevelField(rsf *reflect.StructField, rst *reflect.Value) {
	name := rsf.Tag.Get("env")
	if name == "" {
		return
	}
	envValue := os.Getenv(name)
	if envValue != "" {
		switch rsf.Type.Kind() {
		case reflect.String:
			rst.SetString(envValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			rst.SetInt(int64(StrToInt(envValue, 0)))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			rst.SetUint(uint64(StrToInt(envValue, 0)))
		case reflect.Bool:
			rst.SetBool(StrToBool(envValue))
		default:
			fmt.Printf("not support type: %s", rsf.Type.Kind())
		}
	}
}
