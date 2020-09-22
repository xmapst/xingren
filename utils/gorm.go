package utils

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/noaway/dateparse"
	"reflect"
	"strconv"
	"strings"
)

// StructToMapFilterFields struct转换为map，只支持struct已经写了gorm的column的tag
func StructToMapFilterFields(objPtr interface{}, fileds []string) (map[string]interface{}, error) {
	if len(fileds) < 1 {
		return nil, errors.New("字段不能为空")
	}
	filedMap := make(map[string]struct{}, len(fileds))
	for _, v := range fileds {
		filedMap[v] = struct{}{}
	}
	res := make(map[string]interface{}, len(fileds))
	reflectObjType := reflect.TypeOf(objPtr)
	if reflectObjType.Kind() != reflect.Ptr {
		return nil, errors.New("objPtr只能为指针")
	}
	if reflectObjType.Elem().Kind() != reflect.Struct {
		return nil, errors.New("objPtr只能为struct")
	}
	reflectObjValue := reflect.ValueOf(objPtr).Elem()
	for i := 0; i < reflectObjValue.NumField(); i++ {
		for _, v := range strings.Split(reflectObjValue.Type().Field(i).Tag.Get("gorm"), ";") {
			if strings.Contains(v, "column:") {
				columnNmae := strings.TrimSpace(strings.Replace(v, "column:", "", -1))
				if _, ok := filedMap[columnNmae]; ok {
					res[columnNmae] = reflectObjValue.Field(i).Interface()
					delete(filedMap, columnNmae)
				}
				break
			}
		}
		if len(filedMap) < 1 {
			break
		}
	}
	return res, nil
}

// StructToMap struct转map，先找struct的gorm的column的tag，没有就找json的，在没有就获取字段名称作为key
func StructTMap(objPtr interface{}) (map[string]interface{}, error) {
	reflectObjType := reflect.TypeOf(objPtr)
	if reflectObjType.Kind() != reflect.Ptr {
		return nil, errors.New("objPtr只能为指针")
	}
	if reflectObjType.Elem().Kind() != reflect.Struct {
		return nil, errors.New("objPtr只能为struct")
	}
	reflectObjValue := reflect.ValueOf(objPtr).Elem()
	res := make(map[string]interface{}, reflectObjValue.NumField())

	for i := 0; i < reflectObjValue.NumField(); i++ {
		isGet := false
		for _, v := range strings.Split(reflectObjValue.Type().Field(i).Tag.Get("gorm"), ";") { // 反射获取gorm的column的值
			if strings.Contains(v, "column:") {
				continue
			}
			columnName := strings.TrimSpace(strings.Replace(v, "column:", "", -1))
			if columnName == "" {
				break
			}
			res[columnName] = reflectObjValue.Field(i).Interface()
			isGet = true
			break
		}
		if isGet { // 已经找到gorm column的值
			continue
		}
		for _, v := range strings.Split(reflectObjValue.Type().Field(i).Tag.Get("json"), ",") { // 反射获取json的值
			if v == "-" || v == "omitempty" || v == "" {
				continue
			}
			res[v] = reflectObjValue.Field(i).Interface()
			isGet = true
			break
		}
		if isGet { // 已经找到json的值
			continue
		}
		res[reflectObjValue.Type().Field(i).Name] = reflectObjValue.Field(i).Interface()
	}
	return res, nil
}

// GormNumberColumnToAllName 获取所有的number 类型的 gorm的名称
func GormNumberColumntoAllName(objPtr interface{}) ([]string, error) {
	reflectObjType := reflect.TypeOf(objPtr)
	if reflectObjType.Kind() != reflect.Ptr {
		return nil, errors.New("objPtr只能为指针")
	}
	if reflectObjType.Elem().Kind() != reflect.Struct {
		return nil, errors.New("objPtr只能为struct")
	}
	reflectObjValue := reflect.ValueOf(objPtr).Elem()
	res := make([]string, 0, reflectObjValue.NumField())

	for i := 0; i < reflectObjValue.NumField(); i++ {
		for _, v := range strings.Split(reflectObjValue.Type().Field(i).Tag.Get("gorm"), ";") { // 反射获取gorm的column的值
			if strings.Contains(v, "column:") {
				continue
			}
			columnName := strings.TrimSpace(strings.Replace(v, "column:", "", -1))
			if columnName == "" {
				break
			}
			if !ReflectKindCheckNumber(reflectObjValue.Type().Field(i).Type.Kind()) {
				continue
			}
			res = append(res, columnName)
			break
		}
	}
	return res, nil
}

func ReflectKindCheckNumber(kind reflect.Kind) bool {
	temp := []reflect.Kind{
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Float32,
		reflect.Float64,
	}
	for _, v := range temp {
		if v == kind {
			return true
		}
	}
	return false
}

// GormColumnToAllNmae 获取所有gorm的名称
func GormColumnToAllName(objPtr interface{}) ([]string, error) {
	reflectObjType := reflect.TypeOf(objPtr)
	if reflectObjType.Kind() != reflect.Ptr {
		return nil, errors.New("objPtr只能为指针")
	}
	if reflectObjType.Elem().Kind() != reflect.Struct {
		return nil, errors.New("objPtr只能为struct")
	}
	reflectObjValue := reflect.ValueOf(objPtr).Elem()
	res := make([]string, 0, reflectObjValue.NumField())

	for i := 0; i < reflectObjValue.NumField(); i++ {
		for _, v := range strings.Split(reflectObjValue.Type().Field(i).Tag.Get("gorm"), ";") { // 反射获取gorm的column的值
			if !strings.Contains(v, "column:") {
				continue
			}
			columnName := strings.TrimSpace(strings.Replace(v, "column:", "", -1))
			if columnName == "" || columnName == "id" {
				break
			}
			res = append(res, columnName)
			break
		}
	}

	return res, nil
}

//用map的转存到 struct
//支持 bool、int、int8、int16、int32、int64、uint、uint8、uint16、uint32、uint64、float32、float64、time.Time、string、struct和这些类型与slice或者ptr得组合类型
//map得key支持 orm和gorm tag中得column得名称、json tag中得名称、struct中得名字
func MapFullStruct(data map[string]interface{}, objPtr interface{}) error {
	structValue := reflect.ValueOf(objPtr)
	structType := reflect.TypeOf(objPtr)
	if structType.Kind() == reflect.Ptr { //检查obj是不是指针类型得
		structValue = structValue.Elem() //获取指针下得值
		structType = structType.Elem()   //获取指针下得值
	} else {
		return errors.New("Check type error not *Struct")
	}
	if structType.Kind() != reflect.Struct {
		return errors.New("Check type error not Struct")
	}

	// 存储struct name 转换小写的 name 对应的 index
	structNameMap := make(map[string]int, structValue.NumField())
	// 存储struct json 对应的 index
	structJsonTagMap := make(map[string]int, structValue.NumField())
	//存储struct orm 对应的 index
	structOrmTagMap := make(map[string]int, structValue.NumField())
	//存储struct gorm 对应的 index
	structGormTagMap := make(map[string]int, structValue.NumField())

	for i := 0; i < structValue.NumField(); i++ {
		structNameMap[strings.ToLower(structValue.Type().Field(i).Name)] = i
		jsonTagName := strings.TrimSpace(strings.Replace(structValue.Type().Field(i).Tag.Get("json"), ",omitempty", "", -1))
		if jsonTagName != "" {
			structJsonTagMap[jsonTagName] = i
		}
		for _, v := range strings.Split(structValue.Type().Field(i).Tag.Get("orm"), ";") {
			if strings.Contains(v, "column") {
				// 去空, 去末尾的 ), 替换 column(
				ormTagName := strings.TrimSpace(strings.Trim(strings.Replace(v, "column(", "", -1), ")"))
				if ormTagName != "" {
					structOrmTagMap[ormTagName] = i
				}
				break
			}
		}
		for _, v := range strings.Split(structValue.Type().Field(i).Tag.Get("gorm"), ";") {
			if strings.Contains(v, "column:") {

				gormTagName := strings.TrimSpace(strings.Replace(v, "column:", "", -1))
				if gormTagName != "" {
					structGormTagMap[gormTagName] = i
				}
				break
			}
		}
	}

	for fieldName, value := range data {

		fieldIndex := -1
		//判断名字是否符合
		if index, ok := structNameMap[strings.ToLower(fieldName)]; ok {
			fieldIndex = index
		} else if index, ok := structGormTagMap[fieldName]; ok {
			fieldIndex = index
		} else if index, ok = structOrmTagMap[fieldName]; ok {
			fieldIndex = index
		} else if index, ok = structJsonTagMap[fieldName]; ok {
			fieldIndex = index
		}
		if fieldIndex != -1 && structValue.Field(fieldIndex).CanSet() { //找到位置了,并是字段是可以赋值得
			//转换值
			converValue, err := TypeConversion(value, structValue.Field(fieldIndex))
			if err != nil {
				return err
			}
			if converValue.IsValid() { //判断转化后得值是否有效
				structValue.Field(fieldIndex).Set(converValue)
			}
		}
		//不存在该key时则不管
	}
	return nil
}

// TypeConversion 类型转换
func TypeConversion(value interface{}, toValue reflect.Value) (reflect.Value, error) {
	valueReflect := reflect.ValueOf(value)
	if !valueReflect.IsValid() { //过滤无效值
		return valueReflect, nil
	}
	if toValue.Kind() != reflect.Array && toValue.Kind() != reflect.Slice && valueReflect.Type().Kind() == toValue.Kind() { //如果类型和需要转的类型一致时则直接返回
		return valueReflect, nil
	}
	resValue := reflect.New(toValue.Type()).Elem() //创建一个类型和转换后的类型一样的值
	valueStr := ""
	//一定要类型断言float的相关类型，否则%+v会出现科学计数的方式，导致转换失败
	if _, ok := value.(float64); ok {
		valueStr = fmt.Sprintf("%f", value)
	} else if _, ok := value.(float64); ok {
		valueStr = fmt.Sprintf("%f", value)
	} else if _, ok := value.([]byte); ok {
		valueStr = string(value.([]byte))
	} else {
		valueStr = fmt.Sprintf("%+v", value)
	}
	switch toValue.Kind() {
	case reflect.Bool:
		if strings.ToLower(valueStr) == "on" || strings.ToLower(valueStr) == "yes" || strings.ToLower(valueStr) == "y" {
			resValue.SetBool(true)
			break
		}
		if strings.ToLower(valueStr) == "off" || strings.ToLower(valueStr) == "no" || strings.ToLower(valueStr) == "n" {
			resValue.SetBool(false)
			break
		}
		b, err := strconv.ParseBool(valueStr)
		if err != nil {
			//默认转失败时，尝试转换为float，然后采用非0为true，0为false
			if valueFloat, err := strconv.ParseFloat(valueStr, 64); err == nil {
				resValue.SetBool(valueFloat != 0)
				break
			}
			return resValue, err
		}
		resValue.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		//判断是否包含小数点
		if strings.Contains(valueStr, ".") {
			//去除小数点后面的值
			valueStr = valueStr[:strings.Index(valueStr, ".")]
		}
		x, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return resValue, err
		}
		resValue.SetInt(x)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		//判断是否包含小数点
		if strings.Contains(valueStr, ".") {
			//去除小数点后面的值
			valueStr = valueStr[:strings.Index(valueStr, ".")]
		}
		x, err := strconv.ParseUint(valueStr, 10, 64)
		if err != nil {
			return resValue, err
		}
		resValue.SetUint(x)
	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return resValue, err
		}
		resValue.SetFloat(x)
	case reflect.Interface:
		resValue.Set(reflect.ValueOf(value))
	case reflect.String:
		resValue.SetString(valueStr)
	case reflect.Struct:
		switch toValue.Type().String() {
		case "time.Time":
			reconvTime, err := dateparse.ParseAny(valueStr)
			if err != nil {
				return resValue, err
			}
			resValue.Set(reflect.ValueOf(reconvTime))
		default:
			//校验值是否是map[string]interface{}
			if _, ok := value.(map[string]interface{}); !ok {
				return resValue, fmt.Errorf("%+v 必须为map[string]interface{}", value)
			}
			//反射值得到下一级得类型，并创建根据下一级得类型创建reflect.value
			tt := reflect.New(reflect.TypeOf(toValue.Interface()))
			if err := MapFullStruct(value.(map[string]interface{}), tt.Interface()); err != nil {
				return resValue, err
			}
			resValue.Set(tt.Elem())
		}
	case reflect.Ptr:
		//获取值得类型到下一级得类型，并创建根据类型创建reflect.value
		tt := reflect.New(toValue.Type().Elem())
		//继续递归
		converValue, err := TypeConversion(value, tt.Elem())
		if err != nil {
			return resValue, err
		}
		if converValue.IsValid() { //判断转化后得值是否有效
			tt.Elem().Set(converValue)
		}
		resValue.Set(tt)
	case reflect.Slice:
		//校验值是否是slice
		if valueReflect.Type().Kind() != reflect.Slice {
			return resValue, fmt.Errorf("%+v 必须为slice", value)
		}

		reflectValue := reflect.ValueOf(value)
		//创建切片
		sliceValue := reflect.MakeSlice(toValue.Type(), 0, reflectValue.Len())
		for i := 0; i < reflectValue.Len(); i++ { //反射切片中，遍历切片得值
			//toValue.Type().Elem() 拿到slice得类型
			//根据slice得类型创建相对应得值
			//递归转换切片得中值
			converValue, err := TypeConversion(reflectValue.Index(i).Interface(), reflect.New(toValue.Type().Elem()).Elem())
			if err != nil {
				return resValue, err
			}
			if converValue.IsValid() { //判断转化后得值是否有效
				//将转换好得值添加到切片中
				sliceValue = reflect.Append(sliceValue, converValue)
			}
		}
		resValue.Set(sliceValue)
	}

	return resValue, nil
}

/*
gorm 不固定条件查询封装
and 添加
where := []interface{}{
    []interface{}{"id", "=", 1},
    []interface{}{"username", "xmapst"},
}
结构体条件
where := user.User{ID: 1, UserName: "xmapst"}
in,or条件
where := []interface{}{
    []interface{}{"id", "in", []int{1, 2}},
    []interface{}{"username", "=", "xmapst", "or"},
}
not in,or条件
where := []interface{}{
    []interface{}{"id", "not in", []int{1}},
    []interface{}{"username", "=", "xmapst", "or"},
}
map条件
where := map[string]interface{}{"id": 1, "username": "xmapst"}
and,or混合条件
where := []interface{}{
    []interface{}{"id", "in", []int{1, 2}},
    []interface{}{"username = ? or nickname = ?", "xmapst", "jun"},
}
*/
// BuildWhere 构建where条件
func BuildWhere(db *gorm.DB, where interface{}) (*gorm.DB, error) {
	var err error
	if where == nil {
		return db, nil
	}
	t := reflect.TypeOf(where).Kind()
	if t == reflect.Struct || t == reflect.Map {
		db = db.Where(where)
	} else if t == reflect.Slice {
		for _, item := range where.([]interface{}) {
			item := item.([]interface{})
			column := item[0]
			if reflect.TypeOf(column).Kind() == reflect.String {
				count := len(item)
				if count == 1 {
					return nil, errors.New("切片长度不能小于2")
				}
				columnstr := column.(string)
				// 拼接参数形式
				if strings.Index(columnstr, "?") > -1 {
					db = db.Where(column, item[1:]...)
				} else {
					cond := "and" //cond
					opt := "="
					_opt := " = "
					var val interface{}
					if count == 2 {
						opt = "="
						val = item[1]
					} else {
						opt = strings.ToLower(item[1].(string))
						_opt = " " + strings.ReplaceAll(opt, " ", "") + " "
						val = item[2]
					}

					if count == 4 {
						cond = strings.ToLower(strings.ReplaceAll(item[3].(string), " ", ""))
					}

					/*
					   '=', '<', '>', '<=', '>=', '<>', '!=', '<=>',
					   'like', 'like binary', 'not like', 'ilike',
					   '&', '|', '^', '<<', '>>',
					   'rlike', 'regexp', 'not regexp',
					   '~', '~*', '!~', '!~*', 'similar to',
					   'not similar to', 'not ilike', '~~*', '!~~*',
					*/

					if strings.Index(" in notin ", _opt) > -1 {
						// val 是数组类型
						column = columnstr + " " + opt + " (?)"
					} else if strings.Index(" = < > <= >= <> != <=> like likebinary notlike ilike rlike regexp notregexp", _opt) > -1 {
						column = columnstr + " " + opt + " ?"
					}

					if cond == "and" {
						db = db.Where(column, val)
					} else {
						db = db.Or(column, val)
					}
				}
			} else if t == reflect.Map /*Map*/ {
				db = db.Where(item)
			} else {
				db, err = BuildWhere(db, item)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		return nil, errors.New("参数有误")
	}
	return db, nil
}
