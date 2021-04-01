package utils

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type Rules map[string][]string

type RulesMap map[string]Rules

var CustomizeMap = make(map[string]Rules)

func NotEmpty() string {
	return "notEmpty"
}


//@function: Verify
//@description: 校验方法
//@param: st interface{}, roleMap Rules(入参实例，规则map)
//@return: err error

func Verify(st interface{}, roleMap Rules) (err error) {
	compareMap := map[string]bool{
		"lt": true,
		"le": true,
		"eq": true,
		"ne": true,
		"ge": true,
		"gt": true,
	}

	typ := reflect.TypeOf(st)
	val := reflect.ValueOf(st)

	kd := val.Kind() // 获取到st对应的类别

	if kd != reflect.Struct {
		return errors.New("expect struct")
	}
	num := val.NumField()

	// 遍历结构体的所有字段
	for i := 0; i < num; i++ {
		tagVal := typ.Field(i)
		val := val.Field(i)
		if len(roleMap[tagVal.Name]) > 0 {
			for _, v := range roleMap[tagVal.Name] {
				switch {
				case v == "notEmpty":
					if isBlank(val) {
						return errors.New(tagVal.Name + "值不能为空")
					}
				case compareMap[strings.Split(v, "=")[0]]:
					if !compareVerify(val, v) {
						return errors.New(tagVal.Name + "长度或值不在合法范围," + v)
					}
				}
			}
		}

	}
	return nil
}

//@function: compareVerify
//@description: 长度和数字的校验方法 根据类型自动校验
//@param: value reflect.Value, VerifyStr string
//@return: bool
func compareVerify(value reflect.Value, VerifyStr string) bool {

	switch value.Kind() {
	case reflect.String, reflect.Slice, reflect.Array:
		return compare(value.Len(), VerifyStr)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compare(value.Uint(), VerifyStr)
	case reflect.Float32, reflect.Float64:
		return compare(value.Float(), VerifyStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compare(value.Int(), VerifyStr)
	default:
		return false
	}

}

func compare(value interface{}, VerifyStr string) bool {
	VerifyStrArr := strings.Split(VerifyStr, "=")
	val := reflect.ValueOf(value)

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		Vint, VErr := strconv.ParseInt(VerifyStrArr[1], 10, 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Int() < Vint
		case VerifyStrArr[0] == "le":
			return val.Int() <= Vint
		case VerifyStrArr[0] == "eq":
			return val.Int() == Vint
		case VerifyStrArr[0] == "ne":
			return val.Int() != Vint
		case VerifyStrArr[0] == "ge":
			return val.Int() >= Vint
		case VerifyStrArr[0] == "gt":
			return val.Int() > Vint
		default:
			return false

		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		Vint, VErr := strconv.Atoi(VerifyStrArr[1])
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Uint() < uint64(Vint)
		case VerifyStrArr[0] == "le":
			return val.Uint() <= uint64(Vint)
		case VerifyStrArr[0] == "eq":
			return val.Uint() == uint64(Vint)
		case VerifyStrArr[0] == "ne":
			return val.Uint() != uint64(Vint)
		case VerifyStrArr[0] == "ge":
			return val.Uint() >= uint64(Vint)
		case VerifyStrArr[0] == "gt":
			return val.Uint() > uint64(Vint)
		default:
			return false
		}
	case reflect.Float32, reflect.Float64:
		VFloat, VErr := strconv.ParseFloat(VerifyStrArr[1], 64)
		if VErr != nil {
			return false
		}
		switch {
		case VerifyStrArr[0] == "lt":
			return val.Float() < VFloat
		case VerifyStrArr[0] == "le":
			return val.Float() <= VFloat
		case VerifyStrArr[0] == "eq":
			return val.Float() == VFloat
		case VerifyStrArr[0] == "ne":
			return val.Float() != VFloat
		case VerifyStrArr[0] == "ge":
			return val.Float() >= VFloat
		case VerifyStrArr[0] == "gt":
			return val.Float() > VFloat
		default:
			return false
		}
	default:
		return false
	}

}

//@function: isBlank
//@description: 非空校验
//@param: value reflect.Value
//@return: bool

func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	// TODO 这个地方还没有弄明白
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
