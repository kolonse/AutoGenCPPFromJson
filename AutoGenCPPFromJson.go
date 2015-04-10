// AutoGenCPPFromJson project AutoGenCPPFromJson.go
package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

var globalClassNameMap map[string]int

func formatName(name string) string {
	buff := []byte(name)
	buff[0] &= ^(uint8(32))
	return string(buff)
}

func formatToFunc(key, varType, varName, rapidType string, calssSpace string) []string {
	var parseInt []string
	parseInt = append(parseInt, fmt.Sprintf("bool %s::Parse%s(rapidjson::Value& value)", calssSpace, formatName(key)))
	parseInt = append(parseInt, fmt.Sprintf("{"))
	parseInt = append(parseInt, fmt.Sprintf("\tif(!value.HasMember(\"%s\"))", key))
	parseInt = append(parseInt, fmt.Sprintf("\t{"))
	parseInt = append(parseInt, fmt.Sprintf("\t\tthis ->SetLastError(\"%s not exist\");", key))
	parseInt = append(parseInt, fmt.Sprintf("\t\treturn false;"))
	parseInt = append(parseInt, fmt.Sprintf("\t}"))
	parseInt = append(parseInt, fmt.Sprintf("\tif(!value[\"%s\"].Is%s())", key, rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\t{"))
	parseInt = append(parseInt, fmt.Sprintf("\t\tthis ->SetLastError(\"%s not %s\");", key, rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\t\treturn false;"))
	parseInt = append(parseInt, fmt.Sprintf("\t}"))
	parseInt = append(parseInt, fmt.Sprintf("\tthis ->Set%s(value[\"%s\"].Get%s());", formatName(key), key, rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\treturn true;"))
	parseInt = append(parseInt, fmt.Sprintf("}"))
	parseInt = append(parseInt, fmt.Sprintf(""))
	return parseInt
}

func formatToFuncForArrayMember(key, varType, varName, rapidType string, calssSpace string) []string {
	var parseInt []string
	parseInt = append(parseInt, fmt.Sprintf("bool %s::Parse%s(rapidjson::Value& value)", calssSpace, formatName(key)))
	parseInt = append(parseInt, fmt.Sprintf("{"))
	parseInt = append(parseInt, fmt.Sprintf("\tif(!value.Is%s())", rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\t{"))
	parseInt = append(parseInt, fmt.Sprintf("\t\tthis ->SetLastError(\"%s not %s\");", key, rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\t\treturn false;"))
	parseInt = append(parseInt, fmt.Sprintf("\t}"))
	parseInt = append(parseInt, fmt.Sprintf("\tthis ->Set%s(value.Get%s());", formatName(key), rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\treturn true;"))
	parseInt = append(parseInt, fmt.Sprintf("}"))
	parseInt = append(parseInt, fmt.Sprintf(""))
	return parseInt
}

func formatDouble(key, varType, varName, calssSpace string) []string {
	rapidType := "Double"
	var parseInt []string
	parseInt = append(parseInt, fmt.Sprintf("bool %s::Parse%s(rapidjson::Value& value)", calssSpace, formatName(key)))
	parseInt = append(parseInt, fmt.Sprintf("{"))
	parseInt = append(parseInt, fmt.Sprintf("\tif(!value.HasMember(\"%s\"))", key))
	parseInt = append(parseInt, fmt.Sprintf("\t{"))
	parseInt = append(parseInt, fmt.Sprintf("\t\tthis ->SetLastError(\"%s not exist\");", key))
	parseInt = append(parseInt, fmt.Sprintf("\t\treturn false;"))
	parseInt = append(parseInt, fmt.Sprintf("\t}"))
	parseInt = append(parseInt, fmt.Sprintf("\tif(value[\"%s\"].Is%s())", key, rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\t{"))
	//parseInt = append(parseInt, fmt.Sprintf("\t\tthis ->SetLastError(\"%s not %s\");", key, rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\t\tthis ->Set%s(value[\"%s\"].Get%s());", formatName(key), key, rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\t\treturn true;"))
	parseInt = append(parseInt, fmt.Sprintf("\t}"))

	rapidType = "Int64"
	parseInt = append(parseInt, fmt.Sprintf("\tif(value[\"%s\"].Is%s())", key, rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\t{"))
	//
	parseInt = append(parseInt, fmt.Sprintf("\t\tthis ->Set%s((double)value[\"%s\"].Get%s());", formatName(key), key, rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\t\treturn true;"))
	parseInt = append(parseInt, fmt.Sprintf("\t}"))
	//parseInt = append(parseInt, fmt.Sprintf("\tthis ->Set%s(value[\"%s\"].Get%s());", formatName(key), key, rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\tthis ->SetLastError(\"%s not %s or int\");", key, rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\treturn false;"))
	parseInt = append(parseInt, fmt.Sprintf("}"))
	parseInt = append(parseInt, fmt.Sprintf(""))
	return parseInt
}

func formatDoubleForArray(key, varType, varName string, calssSpace string) []string {
	var parseInt []string
	rapidType := "Double"
	parseInt = append(parseInt, fmt.Sprintf("bool %s::Parse%s(rapidjson::Value& value)", calssSpace, formatName(key)))
	parseInt = append(parseInt, fmt.Sprintf("{"))
	parseInt = append(parseInt, fmt.Sprintf("\tif(value.Is%s())", rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\t{"))
	parseInt = append(parseInt, fmt.Sprintf("\t\tthis ->Set%s(value.Get%s());", formatName(key), rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\t\treturn true;"))
	parseInt = append(parseInt, fmt.Sprintf("\t}"))

	rapidType = "Int64"
	parseInt = append(parseInt, fmt.Sprintf("\tif(value.Is%s())", rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\t{"))
	parseInt = append(parseInt, fmt.Sprintf("\t\tthis ->Set%s((double)value.Get%s());", formatName(key), rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\t\treturn true;"))
	parseInt = append(parseInt, fmt.Sprintf("\t}"))
	parseInt = append(parseInt, fmt.Sprintf("\tthis ->SetLastError(\"%s not %s or int\");", key, rapidType))
	parseInt = append(parseInt, fmt.Sprintf("\treturn false;"))
	parseInt = append(parseInt, fmt.Sprintf("}"))
	parseInt = append(parseInt, fmt.Sprintf(""))
	return parseInt
}

func formatFloat64(key string, value float64, calssSpace string, forWho string) (string, string, []string) {
	var varType, varName string
	var rapidType string
	/** 对不同的类型进行区分 */
	switch value {
	case 0:
		varType = "int"
		varName = "m_n" + formatName(key)
		rapidType = "Int"
	case 1:
		varType = "unsigned int"
		varName = "m_u" + formatName(key)
		rapidType = "Uint"
	case 2:
		varType = "int64_t"
		varName = "m_n64" + formatName(key)
		rapidType = "Int64"
	case 3:
		varType = "uint64_t"
		varName = "m_u64" + formatName(key)
		rapidType = "Uint64"
	case 4:
		varType = "double"
		varName = "m_d" + formatName(key)
		rapidType = "Double"
		if forWho == "object" {
			return varType, varName, formatDouble(key, varType, varName, calssSpace)
		} else {
			return varType, varName, formatDoubleForArray(key, varType, varName, calssSpace)
		}
	default:
		panic(fmt.Sprintf("%s's value must be 0,1,2,3,4 etc.means int32,uint32,int64,uing64,double", key))
	}
	if forWho == "object" {
		return varType, varName, formatToFunc(key, varType, varName, rapidType, calssSpace)
	} else {
		return varType, varName, formatToFuncForArrayMember(key, varType, varName, rapidType, calssSpace)
	}
}

/**
*	return : 变量类型  变量名 变量解析函数
 */
func formatBool(key string, calssSpace string, forWho string) (string, string, []string) {
	varType := "bool"
	varName := "m_b" + formatName(key)
	if forWho == "object" {
		return varType, varName, formatToFunc(key, varType, varName, "Bool", calssSpace)
	} else {
		return varType, varName, formatToFuncForArrayMember(key, varType, varName, "Bool", calssSpace)
	}
}

func formatString(key string, calssSpace string, forWho string) (string, string, []string) {
	varType := "string"
	varName := "m_s" + formatName(key)
	if forWho == "object" {
		return varType, varName, formatToFunc(key, varType, varName, "String", calssSpace)
	} else {
		return varType, varName, formatToFuncForArrayMember(key, varType, varName, "String", calssSpace)
	}
}

/**
*	json对象的成员
*	如果对象是整数int32 那么必须填写0
*				uint32			1
*				int64			2
*				uint64			3
*				double			4
 */
func ProcessObject(className string, tapString string, classSpace string, parseResult map[string]interface{}, hArray *[]string, cppArray *[]string) (string, string) {
	fmt.Println(tapString + "解析对象 " + className)
	classType := string("C") + formatName(className)
	classCountValue, ok := globalClassNameMap[classType]
	if ok { // 如果找到相同的类声明 那么需要将类重新命名
		classType += fmt.Sprintf("_%d", classCountValue)
		globalClassNameMap[classType] = classCountValue + 1
	} else {
		globalClassNameMap[classType] = 0
	}
	var curClassSpace string
	if len(classSpace) == 0 {
		curClassSpace = classType
	} else {
		curClassSpace = classSpace + "::" + classType
	}

	var classHArray []string
	classHArray = append(classHArray, tapString+"class "+classType)
	classHArray = append(classHArray, tapString+"{")

	var classHPublicArray []string
	var classHProtectedArray []string
	var classHPrivateArray []string
	var classCPPArray []string
	/** 构造函数定义 */
	var classConstructFunc []string
	classConstructFunc = append(classConstructFunc, fmt.Sprintf("%s\t%s()", tapString, classType))
	classConstructFunc = append(classConstructFunc, tapString+string("\t{"))

	/** .h 文件初始化规则 */
	classHPublicArray = append(classHPublicArray, tapString+string("public:"))
	classHPublicArray = append(classHPublicArray, tapString+string("\tvoid SetLastError( string para){m_sLastError = para;};"))
	classHPublicArray = append(classHPublicArray, tapString+string("\tstring GetLastError(){return m_sLastError;};"))
	classHPublicArray = append(classHPublicArray, tapString+string("\tbool Unmarshal( const string& sJson );"))
	classHPublicArray = append(classHPublicArray, tapString+string("\tbool Unmarshal( rapidjson::Value& value );"))
	classHPublicArray = append(classHPublicArray, tapString+string("\tstring Marshal();"))
	classHPublicArray = append(classHPublicArray, tapString+string("\tvoid Marshal(rapidjson::Value& value,rapidjson::Document& document);"))
	classHProtectedArray = append(classHProtectedArray, tapString+string("protected:"))
	classHPrivateArray = append(classHPrivateArray, tapString+string("private:"))
	classHPrivateArray = append(classHPrivateArray, tapString+string("\tstring m_sLastError;"))
	/** .cpp文件初始化规则 解析函数生成 */
	classCPPArray = append(classCPPArray, fmt.Sprintf("bool %s::Unmarshal( const string& sJson )", curClassSpace))
	classCPPArray = append(classCPPArray, fmt.Sprintf("{"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\trapidjson::Document parser ;"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\tparser.Parse<0>(sJson.c_str());"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\tif( parser.HasParseError())"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t{"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tthis ->SetLastError(parser.GetParseError());"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t\treturn false;"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t}"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\treturn this ->Unmarshal(*(rapidjson::Value*)&parser);"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("}"))
	classCPPArray = append(classCPPArray, fmt.Sprintf(""))
	classCPPArray = append(classCPPArray, fmt.Sprintf("bool %s::Unmarshal( rapidjson::Value& value )", curClassSpace))
	classCPPArray = append(classCPPArray, fmt.Sprintf("{"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\tif( !value.IsObject())"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t{"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tthis ->SetLastError(\"%s is object,but here not\");", classType))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t\treturn false;"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t}"))
	/** .cpp文件初始化规则 序列化函数生成 */
	var classCPPMarshalFuncArray []string
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("string %s::Marshal()", curClassSpace))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("{"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\trapidjson::Document document ;"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\trapidjson::Value jRoot( rapidjson::kObjectType );"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\tthis ->Marshal(jRoot,document);"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\trapidjson::GenericStringBuffer<rapidjson::UTF8<>> buffer;"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\trapidjson::Writer<rapidjson::GenericStringBuffer<rapidjson::UTF8<>>> writer( buffer );"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\tjRoot.Accept( writer );"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\treturn buffer.GetString();"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("}"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf(""))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("void %s::Marshal( rapidjson::Value& value,rapidjson::Document& document )", curClassSpace))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("{"))
	var classFunIntArray []string
	for key, value := range parseResult {
		/** 增加函数关键字的解析函数声明 */
		classHProtectedArray = append(classHProtectedArray, fmt.Sprintf("%s\tbool Parse%s(rapidjson::Value& value);", tapString, formatName(key)))
		/** 增加函数关键字解析函数实现 */
		var parseInt []string
		var memberName string
		var memberFunc string
		var varName string
		var varType string
		var varTypeValue int
		switch value.(type) {
		case bool:
			varType, varName, parseInt = formatBool(key, curClassSpace, "object")
			varTypeValue = 0
			classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\tvalue.AddMember( \"%s\", this ->%s, document.GetAllocator());", key, varName))
			classConstructFunc = append(classConstructFunc, fmt.Sprintf("%s\t\t%s=true;", tapString, varName))
		case float64:
			varType, varName, parseInt = formatFloat64(key, value.(float64), curClassSpace, "object")
			varTypeValue = 0
			classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\tvalue.AddMember( \"%s\", this ->%s, document.GetAllocator());", key, varName))
			classConstructFunc = append(classConstructFunc, fmt.Sprintf("%s\t\t%s=0;", tapString, varName))
		case string:
			varType, varName, parseInt = formatString(key, curClassSpace, "object")
			varTypeValue = 0
			classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\tvalue.AddMember( \"%s\", this ->%s.c_str(), document.GetAllocator());", key, varName))
		case []interface{}:
			varTypeValue = 1
			var innerClassH []string
			var innerClassCPP []string
			varType, varName = ProcessArray(key, tapString+"\t", curClassSpace, parseResult[key].([]interface{}), &innerClassH, &innerClassCPP)
			var temp []string
			temp = append(temp, classHPublicArray[0:1]...)
			temp = append(temp, innerClassH...)
			temp = append(temp, classHPublicArray[1:]...)
			classHPublicArray = temp
			parseInt = innerClassCPP
			classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t{"))
			classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\trapidjson::Value jArray( rapidjson::kArrayType );"))
			classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\tthis ->%s.Marshal(jArray,document);", varName))
			classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\tvalue.AddMember( \"%s\", jArray, document.GetAllocator());", key))
			classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t}"))
		case map[string]interface{}:
			varTypeValue = 2
			var innerClassH []string
			var innerClassCPP []string
			varType, varName = ProcessObject(key, tapString+"\t", curClassSpace, parseResult[key].(map[string]interface{}), &innerClassH, &innerClassCPP)
			var temp []string
			temp = append(temp, classHPublicArray[0:1]...)
			temp = append(temp, innerClassH...)
			temp = append(temp, classHPublicArray[1:]...)
			classHPublicArray = temp
			parseInt = innerClassCPP
			classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t{"))
			classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\trapidjson::Value jObject( rapidjson::kObjectType );"))
			classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\tthis ->%s.Marshal(jObject,document);", varName))
			classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\tvalue.AddMember( \"%s\", jObject, document.GetAllocator());", key))
			classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t}"))
		}
		memberFunc = fmt.Sprintf("%s\t%s Get%s(){return %s;}", tapString, varType, formatName(key), varName)
		classHPublicArray = append(classHPublicArray, memberFunc)
		memberFunc = fmt.Sprintf("%s\tvoid Set%s(%s para){ %s = para;}", tapString, formatName(key), varType, varName)
		classHPublicArray = append(classHPublicArray, memberFunc)
		memberName = fmt.Sprintf("%s\t%s %s;", tapString, varType, varName)
		classHPrivateArray = append(classHPrivateArray, memberName)
		classFunIntArray = append(classFunIntArray, parseInt...)
		if varTypeValue == 0 {
			classCPPArray = append(classCPPArray, fmt.Sprintf("\tif(!this ->Parse%s(value))", formatName(key)))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t{"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t\treturn false;"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t}"))
		} else if varTypeValue == 1 {
			classCPPArray = append(classCPPArray, fmt.Sprintf("\tif(!value.HasMember(\"%s\"))", key))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t{"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tthis ->SetLastError(\"%s not exist\");", key))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t\treturn false;"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t}"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\tif(!value[\"%s\"].IsArray())", key))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t{"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tthis ->SetLastError(\"%s not Array\");", key))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t\treturn false;"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t}"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\tif(!this ->%s.Unmarshal(value[\"%s\"]))", varName, key))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t{"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tthis ->SetLastError(this ->%s.GetLastError());", varName))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t\treturn false;"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t}"))
		} else {
			classCPPArray = append(classCPPArray, fmt.Sprintf("\tif(!value.HasMember(\"%s\"))", key))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t{"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tthis ->SetLastError(\"%s not exist\");", key))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t\treturn false;"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t}"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\tif(!value[\"%s\"].IsObject())", key))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t{"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tthis ->SetLastError(\"%s not Object\");", key))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t\treturn false;"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t}"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\tif(!this ->%s.Unmarshal(value[\"%s\"]))", varName, key))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t{"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tthis ->SetLastError(this ->%s.GetLastError());", varName))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t\treturn false;"))
			classCPPArray = append(classCPPArray, fmt.Sprintf("\t}"))
		}
	}
	classConstructFunc = append(classConstructFunc, tapString+string("\t}"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("}"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\treturn true;"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("}"))
	classCPPArray = append(classCPPArray, fmt.Sprintf(""))
	classCPPArray = append(classCPPArray, classFunIntArray...)
	*cppArray = append(*cppArray, classCPPArray...)
	*cppArray = append(*cppArray, classCPPMarshalFuncArray...)

	var temp []string
	temp = append(temp, classHPublicArray[0:1]...)
	temp = append(temp, classConstructFunc...)
	temp = append(temp, classHPublicArray[1:]...)
	classHPublicArray = temp
	classHArray = append(classHArray, classHPublicArray...)
	classHArray = append(classHArray, classHProtectedArray...)
	classHArray = append(classHArray, classHPrivateArray...)
	classHArray = append(classHArray, tapString+"};")
	*hArray = append(*hArray, classHArray...)
	return classType, fmt.Sprintf("m_object%s", formatName(className))
}

func ProcessArray(className string, tapString string, classSpace string, parseResult []interface{}, hArray *[]string, cppArray *[]string) (string, string) {
	fmt.Println(tapString + "解析数组 " + className)
	classType := string("C") + formatName(className)
	classCountValue, ok := globalClassNameMap[classType]
	if ok { // 如果找到相同的类声明 那么需要将类重新命名
		classType += fmt.Sprintf("_%d", classCountValue)
		globalClassNameMap[classType] = classCountValue + 1
	} else {
		globalClassNameMap[classType] = 0
	}
	var curClassSpace string
	if len(classSpace) == 0 {
		curClassSpace = classType
	} else {
		curClassSpace = classSpace + "::" + classType
	}

	var classHArray []string
	classHArray = append(classHArray, tapString+"class "+classType)
	classHArray = append(classHArray, tapString+"{")

	var classHPublicArray []string
	var classHProtectedArray []string
	var classHPrivateArray []string
	var classCPPArray []string
	/** .h 文件初始化规则 */
	classHPublicArray = append(classHPublicArray, tapString+string("public:"))
	classHPublicArray = append(classHPublicArray, tapString+string("\tvoid SetLastError( string para){m_sLastError = para;};"))
	classHPublicArray = append(classHPublicArray, tapString+string("\tstring GetLastError(){return m_sLastError;};"))
	classHPublicArray = append(classHPublicArray, tapString+string("\tbool Unmarshal( const string& sJson );"))
	classHPublicArray = append(classHPublicArray, tapString+string("\tbool Unmarshal( rapidjson::Value& value );"))
	classHPublicArray = append(classHPublicArray, tapString+string("\tstring Marshal();"))
	classHPublicArray = append(classHPublicArray, tapString+string("\tvoid Marshal(rapidjson::Value& value,rapidjson::Document& document);"))
	classHProtectedArray = append(classHProtectedArray, tapString+string("protected:"))
	classHPrivateArray = append(classHPrivateArray, tapString+string("private:"))
	classHPrivateArray = append(classHPrivateArray, tapString+string("\tstring m_sLastError;"))
	/** .cpp文件初始化规则 */
	classCPPArray = append(classCPPArray, fmt.Sprintf("bool %s::Unmarshal( const string& sJson )", curClassSpace))
	classCPPArray = append(classCPPArray, fmt.Sprintf("{"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\trapidjson::Document parser ;"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\tparser.Parse<0>(sJson.c_str());"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\tif( parser.HasParseError())"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t{"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tthis ->SetLastError(parser.GetParseError());"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t\treturn false;"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t}"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\treturn this ->Unmarshal(*(rapidjson::Value*)&parser);"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("}"))
	classCPPArray = append(classCPPArray, fmt.Sprintf(""))
	classCPPArray = append(classCPPArray, fmt.Sprintf("bool %s::Unmarshal( rapidjson::Value& value )", curClassSpace))
	classCPPArray = append(classCPPArray, fmt.Sprintf("{"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\tif( !value.IsArray())"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t{"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tthis ->SetLastError(\"%s is array,but here not\");", classType))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t\treturn false;"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t}"))
	/** .cpp文件初始化规则 序列化函数生成 */
	var classCPPMarshalFuncArray []string
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("string %s::Marshal()", curClassSpace))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("{"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\trapidjson::Document document ;"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\trapidjson::Value jRoot( rapidjson::kArrayType );"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\tthis ->Marshal(jRoot,document);"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\trapidjson::GenericStringBuffer<rapidjson::UTF8<>> buffer;"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\trapidjson::Writer<rapidjson::GenericStringBuffer<rapidjson::UTF8<>>> writer( buffer );"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\tjRoot.Accept( writer );"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\treturn buffer.GetString();"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("}"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf(""))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("void %s::Marshal( rapidjson::Value& value,rapidjson::Document& document )", curClassSpace))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("{"))
	var classFunIntArray []string
	/** 先暂时只支持一种类型 支持多种类型的后续在考虑*/
	if len(parseResult) == 0 || len(parseResult) > 1 {
		panic("json 数组暂时只能支持一种类型")
	}

	var varType, varName string
	key := "ArrayMember"
	varTypeValue := 0
	switch parseResult[0].(type) {
	case bool:
		varType, varName, classFunIntArray = formatBool(key, curClassSpace, "array")
	case float64:
		varType, varName, classFunIntArray = formatFloat64(key, parseResult[0].(float64), curClassSpace, "array")
	case string:
		varType, varName, classFunIntArray = formatString(key, curClassSpace, "array")
	case []interface{}:
		varTypeValue = 1
		var innerClassH []string
		var innerClassCPP []string
		varType, varName = ProcessArray(key, tapString+"\t", curClassSpace, parseResult[0].([]interface{}), &innerClassH, &innerClassCPP)
		var temp []string
		temp = append(temp, classHPublicArray[0:1]...)
		temp = append(temp, innerClassH...)
		temp = append(temp, classHPublicArray[1:]...)
		classHPublicArray = temp
		classFunIntArray = innerClassCPP
	case map[string]interface{}:
		varTypeValue = 2
		var innerClassH []string
		var innerClassCPP []string
		varType, varName = ProcessObject(key, tapString+"\t", curClassSpace, parseResult[0].(map[string]interface{}), &innerClassH, &innerClassCPP)
		var temp []string
		temp = append(temp, classHPublicArray[0:1]...)
		temp = append(temp, innerClassH...)
		temp = append(temp, classHPublicArray[1:]...)
		classHPublicArray = temp
		classFunIntArray = innerClassCPP
	}

	/** 对数组元素序列化 */
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\tfor(unsigned int i = 0; i < this ->m_vElemGroup.size(); ++ i)"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t{"))
	switch parseResult[0].(type) {
	case bool:
		classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\tvalue.PushBack( this ->m_vElemGroup[i], document.GetAllocator());"))
	case float64:
		classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\tvalue.PushBack( this ->m_vElemGroup[i], document.GetAllocator());"))
	case string:
		varType, varName, classFunIntArray = formatString(key, curClassSpace, "array")
		classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\tvalue.PushBack( this ->m_vElemGroup[i].c_str(), document.GetAllocator());"))
	case []interface{}:
		classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\t{"))
		classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\t\trapidjson::Value jArray( rapidjson::kArrayType );"))
		classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\t\tthis ->m_vElemGroup[i].Marshal(jArray,document);"))
		classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\t\tvalue.PushBack( jArray, document.GetAllocator());"))
		classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\t}"))
	case map[string]interface{}:
		classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\t{"))
		classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\t\trapidjson::Value jObject( rapidjson::kObjectType );"))
		classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\t\tthis ->m_vElemGroup[i].Marshal(jObject,document);"))
		classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\t\tvalue.PushBack( jObject, document.GetAllocator());"))
		classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t\t}"))
	}
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("\t}"))
	classCPPMarshalFuncArray = append(classCPPMarshalFuncArray, fmt.Sprintf("}"))
	classHProtectedArray = append(classHProtectedArray, fmt.Sprintf("%s\tbool Parse%s(rapidjson::Value& value);", tapString, formatName(key)))
	memberFunc := fmt.Sprintf("%s\t%s Get%s(){return %s;}", tapString, varType, formatName(key), varName)
	classHProtectedArray = append(classHProtectedArray, memberFunc)
	memberFunc = fmt.Sprintf("%s\tvoid Set%s(%s para){ %s = para;}", tapString, formatName(key), varType, varName)
	classHProtectedArray = append(classHProtectedArray, memberFunc)
	memberName := fmt.Sprintf("%s\t%s %s;", tapString, varType, varName)
	classHPrivateArray = append(classHPrivateArray, memberName)
	classHPublicArray = append(classHPublicArray, fmt.Sprintf("%s\tvector<%s> GetElemGroup(){return this ->m_vElemGroup;}", tapString, varType))
	classHPublicArray = append(classHPublicArray, fmt.Sprintf("%s\tvoid SetElemGroup(vector<%s> vElemGroup){this ->m_vElemGroup = vElemGroup;}", tapString, varType))
	classHPrivateArray = append(classHPrivateArray, fmt.Sprintf("%s\tvector<%s> m_vElemGroup;", tapString, varType))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\trapidjson::Document::ValueIterator it = value.Begin();"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\tfor(; it != value.End();it ++)"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t{"))
	if varTypeValue == 0 {
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tif(!this ->Parse%s(*it))", formatName(key)))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t{"))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t\treturn false;"))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t}"))
	} else if varTypeValue == 1 {
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tif(!(*it).IsArray())"))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t{"))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t\tthis ->SetLastError(\"%s not Array\");", key))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t\treturn false;"))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t}"))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tif(!this ->%s.Unmarshal(*it))", varName))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t{"))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t\tthis ->SetLastError(this ->%s.GetLastError());", varName))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t\treturn false;"))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t}"))
	} else {
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tif(!(*it).IsObject())"))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t{"))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t\tthis ->SetLastError(\"%s not Object\");", key))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t\treturn false;"))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t}"))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tif(!this ->%s.Unmarshal(*it))", varName))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t{"))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t\tthis ->SetLastError(this ->%s.GetLastError());", varName))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t\treturn false;"))
		classCPPArray = append(classCPPArray, fmt.Sprintf("\t\t}"))
	}
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t\tthis ->m_vElemGroup.push_back(this ->Get%s());", formatName(key)))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\t}"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("\treturn true;"))
	classCPPArray = append(classCPPArray, fmt.Sprintf("}"))
	classCPPArray = append(classCPPArray, fmt.Sprintf(""))
	classCPPArray = append(classCPPArray, classFunIntArray...)
	*cppArray = append(*cppArray, classCPPArray...)
	*cppArray = append(*cppArray, classCPPMarshalFuncArray...)
	classHArray = append(classHArray, classHPublicArray...)
	classHArray = append(classHArray, classHProtectedArray...)
	classHArray = append(classHArray, classHPrivateArray...)
	classHArray = append(classHArray, tapString+"};")
	*hArray = append(*hArray, classHArray...)
	return classType, fmt.Sprintf("m_array%s", formatName(className))
}

func WriteSourceFile(fileName string, fileContentArray []string) {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for _, writeLine := range fileContentArray {
		_, err = f.Write([]byte(writeLine))
		_, err = f.Write([]byte("\n"))
	}
}

/**
*	usage: AutoGenCPPFromJson.exe [json file] [GenName] [GenDir]
*	生成结果:在当前目录下会生成 GenName.h 和 GenName.cpp 文件
 */
func main() {
	arg_num := len(os.Args)
	if arg_num < 4 { // 需要有二个参数
		fmt.Println("args count error")
		panic(-1)
	}

	jsonFile := os.Args[1]
	genName := os.Args[2]
	genDir := os.Args[3]
	if genDir[len(genDir)-1] != '\\' && genDir[len(genDir)-1] != '/' {
		genDir += string("/")
	}

	genDir += genName
	/** 创建目标目录 */
	err := os.Mkdir(genDir, 0777)
	jsonContent, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}
	// 创建 map 实例
	globalClassNameMap = make(map[string]int)

	var parseResult interface{}
	err = json.Unmarshal(jsonContent, &parseResult)
	if err != nil {
		panic(err)
	}
	/** 头文件预写入信息 **/
	genNameCPP_H := string(genDir) + "/" + string(genName) + ".h"
	var genNameCPP_HArray []string
	curTimeMd5 := fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String())))
	defHead := string("JSZHOU2_") + genName + "_" + curTimeMd5

	genNameCPP_HArray = append(genNameCPP_HArray, string("#ifndef ")+defHead)
	genNameCPP_HArray = append(genNameCPP_HArray, string("#define ")+defHead)
	genNameCPP_HArray = append(genNameCPP_HArray, string("#include <string>"))
	genNameCPP_HArray = append(genNameCPP_HArray, string("#include <vector>"))
	genNameCPP_HArray = append(genNameCPP_HArray, string("#include <rapidjson/rapidjson.h>"))
	genNameCPP_HArray = append(genNameCPP_HArray, string("#include <rapidjson/document.h>"))
	genNameCPP_HArray = append(genNameCPP_HArray, string("#include <rapidjson/stringbuffer.h>"))
	genNameCPP_HArray = append(genNameCPP_HArray, string("#include <rapidjson/writer.h>"))

	genNameCPP_HArray = append(genNameCPP_HArray, string("using namespace std;"))
	genNameCPP_HArray = append(genNameCPP_HArray, string("\n"))
	genNameCPP_HArray = append(genNameCPP_HArray, string("\n"))
	/** CPP文件预写入信息 */
	genNameCPP_CPP := string(genDir) + "/" + string(genName) + ".cpp"
	var genNameCPP_CPPArray []string
	genNameCPP_CPPArray = append(genNameCPP_CPPArray, string("#include \"")+string(genName)+".h\"")
	/** 判断结构体为对象类型还是数组类型 */
	switch parseResult.(type) {
	case []interface{}:
		fmt.Println("数组类型解析")
		ProcessArray(string(genName), string(""), string(""), parseResult.([]interface{}), &genNameCPP_HArray, &genNameCPP_CPPArray)
	case map[string]interface{}:
		fmt.Println("对象类型解析")
		ProcessObject(string(genName), string(""), string(""), parseResult.(map[string]interface{}), &genNameCPP_HArray, &genNameCPP_CPPArray)
	default:
		panic("json 格式有问题")
	}
	/** 头文件结尾信息写入 */
	genNameCPP_HArray = append(genNameCPP_HArray, string("\n"))
	genNameCPP_HArray = append(genNameCPP_HArray, string("#endif"))

	/** 生成.h和.cpp文件 */
	WriteSourceFile(genNameCPP_H, genNameCPP_HArray)
	WriteSourceFile(genNameCPP_CPP, genNameCPP_CPPArray)
}
