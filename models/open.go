package models

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type Open struct{}

// 常量
const (
	//字母
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//数组
	number = "0123456789"
	//6位表示字母索引
	letterIdBits = 6
	//与letterIdBits一样多的所有1位
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

// 随机种子
var src = rand.NewSource(time.Now().UnixNano())

// 取随机字符(n:长度)
func Rand_String(n int) string {
	b := make([]byte, n)
	// 一个rand.Int63（）生成63个随机位，足够letterIdMax字母使用！
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

// 取随机数(n:长度)
func Rand_Number(n int) string {
	b := make([]byte, n)
	// 一个rand.Int63（）生成63个随机位，足够letterIdMax字母使用！
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(number) {
			b[i] = number[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

// GetSign 获取签名()
func GetSign(data map[string]interface{}, secret string) string {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var originStr string
	for _, v := range keys {
		originStr += fmt.Sprintf("%v=%v&", v, data[v])
	}
	originStr += fmt.Sprintf("secret=%v", secret)
	//转小写
	sign := strings.ToLower(originStr)
	return GetMD5(sign)
}

// GetSignRes 获取返回值的签名()
func GetSignRes(i interface{}) string {
	var data map[string]interface{}
	if reflect.TypeOf(i) == reflect.TypeOf(map[string]interface{}{}) {
		data = i.(map[string]interface{})
	} else {
		data = StructToMap(i) //结构体转map
	}
	var keys []string //排序后的键名数组
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys) //排序
	var originStr string
	for _, v := range keys {
		originStr += fmt.Sprintf("%v=%v&", v, data[v])
	}
	originStr += fmt.Sprintf("secret=%v", Genket)
	//转小写
	sign := strings.ToLower(originStr) //签名
	//fmt.Printf("sign = %v\n", sign)
	return GetMD5(sign)
}

// CheckSign 校验签名()
func CheckSign(data map[string]interface{}, secret string) bool {
	//fmt.Printf("CheckSign() data = %v\n", data)
	// 对data按照字母顺序进行排序并组合
	var keys []string
	for k := range data["Data"].(map[string]interface{}) {
		keys = append(keys, k)
	}
	// 排序
	sort.Strings(keys)

	var originStr string
	for _, v := range keys {
		if v == "time_stamp" {
			originStr += fmt.Sprintf("%v=%d&", v, int64(data["Data"].(map[string]interface{})[v].(float64)))
		} else {
			originStr += fmt.Sprintf("%v=%v&", v, data["Data"].(map[string]interface{})[v])
		}
	}
	originStr += fmt.Sprintf("secret=%v", secret)
	signs := strings.ToLower(originStr) //转小写
	MD5 := GetMD5(signs)
	//fmt.Println("拼接后的签名 = ", signs)
	//fmt.Println("拼接后的签名MD5 = ", MD5)
	return MD5 == data["Sign"]
}

// GetMD5 取MD5(字符串) 字符串
func GetMD5(s string) string {
	has := md5.Sum([]byte(s))
	MD5 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return MD5
}

// 获取UUID
func GetUUID() string {
	u4 := uuid.New()
	return u4.String()
}

// 获取卡密时间(s int, i int)
// s: 0-年,1-日,2-月,3-年
func Get_Carmi_Time(s int, i int) int64 {
	var t int
	switch s {
	case 3: //"年"
		t = i * 31536000
	case 2: //"月"
		t = i * 2592000
	case 1: //"日"
		t = i * 86400
	case 0: //"时"
		t = i * 3600
	}
	return int64(t)
}

// ValidateJSON 校验json是否有效
func ValidateJSON(data []byte) bool {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return false
	}
	return true
}

// JsonToMap json转Map
func JsonToMap(data []byte) map[string]interface{} {
	// 将单引号替换为双引号（注意：这可能会引起一些未预见的替换问题，应确保输入格式的一致性）
	modifiedData := strings.Replace(string(data), "'", "\"", -1)
	dec := json.NewDecoder(bytes.NewReader([]byte(modifiedData)))
	dec.UseNumber() // 让解码器使用json.Number类型
	var result map[string]interface{}
	if err := dec.Decode(&result); err != nil {
		fmt.Printf("JSON 解析失败: %v\n", err)
		return nil
	}
	return result
}

// JsonToStruct json转结构体
func JsonToStruct(data []byte, result interface{}) error {
	err := json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("JSON 解析失败: ", err.Error())
		return err
	}
	return nil
}

// MapToJson 将map转换为json
func MapToJson(m interface{}) []byte {
	// 将 map 转换为 JSON 字符串
	jsonData, err := json.Marshal(m)
	if err != nil {
		fmt.Println("错误 MapToJson(): ", err.Error())
	}
	return jsonData
}

// StructToMap 将结构体转换为map
func StructToMap(p interface{}) map[string]interface{} {

	t := reflect.TypeOf(p)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	v := reflect.ValueOf(p)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	result := make(map[string]interface{}) //创建一个map
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		result[field.Name] = value.Interface() //将每个字段的值赋值给map
	}
	return result
}

// BytesToMap 转换为map
func BytesToMap(b []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		fmt.Printf("错误 BytesToMap(): %s \n", err.Error())
		return nil, err
	}
	return m, nil
}

func BytesToMap1(b []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		fmt.Printf("错误 BytesToMap(): %s \n", err.Error())
		return nil, err
	}

	// 增加日志记录以检查数据结构
	fmt.Printf("解析后的数据: %v\n", m)

	// 检查关键字段的类型，防止类型转换错误
	if timeStampValue, ok := m["time_stamp"]; ok {
		switch v := timeStampValue.(type) {
		case float64:
			fmt.Printf("time_stamp 类型为 float64，值为: %f\n", v)
			m["time_stamp"] = int64(v)
		case string:
			fmt.Printf("time_stamp 类型为 string，值为: %s\n", v)
			// 添加从 string 到 int64 的转换逻辑
			// 示例: int64Value, err := strconv.ParseInt(v, 10, 64)
		default:
			fmt.Printf("time_stamp 类型为 %T，值为: %v\n", v, v)
		}
	} else {
		fmt.Println("time_stamp 字段不存在或类型不匹配")
	}

	return m, nil
}

// InterfaceToMap 将interface转换为map
func InterfaceToMap(b interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal(b.([]byte), &m); err != nil {
		fmt.Printf("错误 InterfaceToMap(): %s \n", err.Error())
		return nil, err
	}
	return m, nil
}

// StructToByte 将结构体转换为Byte
func StructToByte(s interface{}) []byte {
	bytes, err := json.Marshal(s) //将结构体转换为json
	if err != nil {
		fmt.Println("错误 json.Marshal(): ", err)
		return nil
	}
	return bytes
}

// ToInt 将interface转换为int
func ToInt(i interface{}) int {
	var val int
	_, err := fmt.Sscanf(fmt.Sprint(i), "%d", &val)
	if err != nil {
		fmt.Println("interface转换为int失败: ", err.Error())
	}
	return val
}

func StringToInt(i string) (int, error) {
	return strconv.Atoi(i)
}

// InterfaceToByte 将interface转换为Byte
func InterfaceToByte(i interface{}) []byte {
	jsonData, err := json.Marshal(i)
	if err != nil {
		fmt.Println("错误 InterfaceToByte(): ", err.Error())
	}
	return jsonData
}

// InterfaceToString 将interface转换为String
func InterfaceToString(i interface{}) string {
	b, err := json.Marshal(i)
	if err != nil {
		fmt.Println("错误 InterfaceToString(): ", err.Error())
	}
	return string(b)
}

// Generate_Password 给密码就行加密操作
func Generate_Password(userPassword string) (string, error) {
	str, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("错误 Generate_Password(): ", err.Error())
		return "", err
	}
	return string(str), err
}

// Validate_Password 密码比对
func Validate_Password(userPassword string, hashed string) (isOK bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		fmt.Println("错误 Validate_Password(): ", err.Error())
		return false, errors.New("密码比对错误！")
	}
	return true, nil
}

// GetError 自定义错误消息
func GetError(errs validator.ValidationErrors, r interface{}) string {
	s := reflect.TypeOf(r)
	for _, fieldError := range errs {
		filed, _ := s.FieldByName(fieldError.Field())
		errTag := fieldError.Tag() + "_err"
		// 获取对应binding的错误消息
		errTagText := filed.Tag.Get(errTag)
		// 获取统一错误消息
		errText := filed.Tag.Get("err")
		if errTagText != "" {
			return errTagText
		}
		if errText != "" {
			return errText
		}
		return fieldError.Field() + ":" + fieldError.Tag()
	}
	return ""
}

// PrintStruct 打印结构体
func PrintStruct(s interface{}) {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Struct {
		typeOfS := v.Type()
		fmt.Printf("----------------------结构体内容----------------------\n")
		for i := 0; i < v.NumField(); i++ {
			field := typeOfS.Field(i)
			value := v.Field(i).Interface()
			fmt.Printf("%s: %v \n", field.Name, value)
		}
		fmt.Printf("----------------------------------------------------\n")
	} else {
		fmt.Println("Not a struct")
	}
}
