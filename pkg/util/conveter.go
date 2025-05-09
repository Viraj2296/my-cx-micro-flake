package util

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"time"
)

func ToBinaryRunes(s string) string {
	var buffer bytes.Buffer
	for _, runeValue := range s {
		fmt.Fprintf(&buffer, "%b", runeValue)
	}
	return fmt.Sprintf("%s", buffer.Bytes())
}

func ToBinaryBytes(s string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(s); i++ {
		fmt.Fprintf(&buffer, "%b", s[i])
	}
	return fmt.Sprintf("%s", buffer.Bytes())
}

func ToEncodeBase64String(_zInput string) string {
	return b64.StdEncoding.EncodeToString([]byte(_zInput))
}

func ToDecodeBase64String(_zInput string) ([]byte, error) {
	return b64.StdEncoding.DecodeString(_zInput)
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func getSize(content io.Seeker) (int64, error) {
	size, err := content.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, err
	}
	_, err = content.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return size, nil
}

func InterfaceToFloat(interfaceValue interface{}) float64 {
	fmt.Println("reflect.TypeOf(interfaceValue)", reflect.TypeOf(interfaceValue))
	switch v := interfaceValue.(type) {
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case float64:
		return v
	case string:
		floatValue, err := strconv.ParseFloat(v, 64)

		if err != nil {
			return 0
		}
		return floatValue
	case []uint8:
		var rawData []byte
		for _, element := range rawData {
			rawData = append(rawData, byte(element))
		}
		stringData := string(rawData)
		floatValue, err := strconv.ParseFloat(stringData, 64)
		if err != nil {
			return 0
		}
		return floatValue
	default:
		objectType := fmt.Sprintf("%T", v)
		fmt.Println("invalid data type [", objectType, "]")
		return 0
	}
}
func InterfaceToTime(value time.Time) string {
	// Format the time as a string using a specific layout
	// Adjust the layout based on your preferred date and time format
	return value.Format("2006-01-02 15:04:05")
}
func InterfaceToInt(interfaceValue interface{}) int {
	switch v := interfaceValue.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		intValue, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return intValue
	case []uint8:
		var rawData []byte
		for _, element := range rawData {
			rawData = append(rawData, byte(element))
		}
		stringData := string(rawData)
		intValue, err := strconv.Atoi(stringData)
		if err != nil {
			return 0
		}
		return intValue
	default:
		objectType := fmt.Sprintf("%T", v)
		fmt.Println("invalid data type [", objectType, "]")
		return 0
	}
}

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}
func RemoveFromIntArray(sourceArray interface{}, removeElementsArray interface{}) []int {
	var sourceIntArray []int

	switch x := sourceArray.(type) {
	case []interface{}:
		for _, i := range x {
			sourceIntArray = append(sourceIntArray, InterfaceToInt(i))
		}
		switch removeElementInterface := removeElementsArray.(type) {
		case []interface{}:
			var removeElementArray []int
			for _, removeElement := range removeElementInterface {
				removeElementArray = append(removeElementArray, InterfaceToInt(removeElement))
			}
			var remainingIntArray []int
			for _, sourceElement := range sourceIntArray {
				index := sort.SearchInts(removeElementArray, sourceElement)

				if index != len(removeElementArray) {
					// not found
					remainingIntArray = append(remainingIntArray, sourceElement)
				}
			}

			return remainingIntArray

		case interface{}:
			fmt.Println("sourceIntArray:", sourceIntArray)
			fmt.Println("InterfaceToInt(removeElementsArray):", InterfaceToInt(removeElementsArray))
			var remainingIntArray []int
			for _, sourceElement := range sourceIntArray {
				if sourceElement != InterfaceToInt(removeElementsArray) {
					remainingIntArray = append(remainingIntArray, sourceElement)
				}
			}
			return remainingIntArray
		}
	default:
		return sourceIntArray
	}
	return sourceIntArray
}

func InterfaceToIntArray(interfaceArray interface{}) []int {
	var convertedIntArray []int
	switch x := interfaceArray.(type) {
	case []interface{}:
		for _, i := range x {
			convertedIntArray = append(convertedIntArray, InterfaceToInt(i))
		}
	default:
		return convertedIntArray
	}
	return convertedIntArray
}
func InterfaceToStringArray(interfaceArray interface{}) []string {
	var convertedIntArray []string
	switch x := interfaceArray.(type) {
	case []interface{}:
		for _, i := range x {
			convertedIntArray = append(convertedIntArray, InterfaceToString(i))
		}
	default:
		fmt.Println("InterfaceToStringArray, invalid type :", reflect.TypeOf(interfaceArray))
		return convertedIntArray
	}
	return convertedIntArray
}

func AppendToObjectArray(src interface{}, appendData interface{}) []interface{} {
	var modifiedArray []interface{}
	switch x := src.(type) {
	case []interface{}:
		for _, i := range x {
			modifiedArray = append(modifiedArray, i)
		}
		modifiedArray = append(modifiedArray, appendData)
	default:
		return modifiedArray
	}
	return modifiedArray
}

// AppendToIntArray this function add just int  or array of int to an existing array
func AppendToIntArray(existingDataArray interface{}, appendData interface{}) []int {
	var modifiedArray []int
	switch x := existingDataArray.(type) {
	case []interface{}:
		for _, i := range x {
			modifiedArray = append(modifiedArray, InterfaceToInt(i))
		}
		switch appendInterface := appendData.(type) {
		case []interface{}:
			for _, appendValue := range appendInterface {
				modifiedArray = append(modifiedArray, InterfaceToInt(appendValue))
			}
		case interface{}:
			modifiedArray = append(modifiedArray, InterfaceToInt(appendData))
		}
	default:
		return RemoveDuplicateInt(modifiedArray)
	}
	return RemoveDuplicateInt(modifiedArray)
}

func InterfaceArrayToCommaSeperatedString(interfaceValue interface{}) string {
	var commaSeperatedString string
	switch x := interfaceValue.(type) {
	case []interface{}:
		for index, i := range x {
			if index == len(x)-1 {
				commaSeperatedString = commaSeperatedString + InterfaceToString(i)
			} else {
				commaSeperatedString = commaSeperatedString + InterfaceToString(i) + ","
			}
		}
	case []int:
		for index, i := range x {
			if index == len(x)-1 {
				commaSeperatedString = commaSeperatedString + InterfaceToString(i)
			} else {
				commaSeperatedString = commaSeperatedString + InterfaceToString(i) + ","
			}
		}

	default:
		objectType := fmt.Sprintf("%T", x)
		fmt.Println("invalid data type in InterfaceArrayToCommaSeperatedString [", objectType, "]")
		return commaSeperatedString
	}
	return commaSeperatedString
}

func InterfaceToBool(interfaceValue interface{}) bool {
	switch v := interfaceValue.(type) {
	case int:

		if v == 0 {
			return false
		} else {
			return true
		}
	case int32:
		if v == 0 {
			return false
		} else {
			return true
		}
	case int64:
		if v == 0 {
			return false
		} else {
			return true
		}
	case bool:
		if v == true {
			return true
		} else {
			return false
		}
	case float64:
		intValue := int(v)
		if intValue == 0 {
			return false
		} else {
			return true
		}
	case string:
		if v == "false" {
			return false
		} else if v == "true" {
			return true
		}
		return false
	case []uint8:

		rawString := string(v)
		if rawString == "false" {
			return false
		} else if rawString == "true" {
			return true
		}
		return false
	case uint8:
		stringData := string(v)
		if stringData == "false" {
			return false
		} else if stringData == "true" {
			return true
		}
		return false
	default:
		objectType := fmt.Sprintf("%T", v)
		fmt.Println("invalid data type in InterfaceToBool [", objectType, "]")
		return false
	}
}

func MapInterfaceToString(objectFields map[string]interface{}, key string) string {

	if object, ok := objectFields[key]; ok {
		//do something here
		switch v := object.(type) {
		case int:
			return strconv.Itoa(v)
		case int32:
			return strconv.Itoa(int(v))
		case int64:
			return strconv.Itoa(int(v))
		case float64:
			intValue := int(v)
			return strconv.Itoa(intValue)
		case string:
			return v
		case []uint8:
			var rawData []byte
			for _, element := range v {
				rawData = append(rawData, byte(element))
			}
			return string(rawData)
		case uint8:
			return string(v)
		default:
			objectType := fmt.Sprintf("%T", v)
			fmt.Println("invalid data type in InterfaceToString [", objectType, "]")
			return ""
		}
	} else {
		return ""
	}

}

func InterfaceToString(interfaceValue interface{}) string {
	switch v := interfaceValue.(type) {
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.Itoa(int(v))
	case int64:
		return strconv.Itoa(int(v))
	case float64:
		intValue := int(v)
		return strconv.Itoa(intValue)
	case string:
		return v
	case []uint8:
		var rawData []byte
		for _, element := range v {
			rawData = append(rawData, byte(element))
		}
		return string(rawData)
	case uint8:
		return string(v)
	default:
		objectType := fmt.Sprintf("%T", v)
		fmt.Println("invalid data type in InterfaceToString [", objectType, "]")
		return ""
	}
}

func Unique(src interface{}) interface{} {
	srcValueOf := reflect.ValueOf(src)
	dstValueOf := reflect.MakeSlice(srcValueOf.Type(), 0, 0)
	visited := make(map[interface{}]struct{})
	for i := 0; i < srcValueOf.Len(); i++ {
		elementValueOf := srcValueOf.Index(i)
		if _, ok := visited[elementValueOf.Interface()]; ok {
			continue
		}
		visited[elementValueOf.Interface()] = struct{}{}
		dstValueOf = reflect.Append(dstValueOf, elementValueOf)
	}
	return dstValueOf.Interface()
}

func RemoveDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := make([]int, 0)
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func RemoveDuplicateString(src []string) []string {

	// Create a map to track the strings that have been seen
	seen := make(map[string]bool)

	// Create a new slice to hold the unique strings
	var unique []string

	// Iterate through the slice of strings
	for _, str := range src {
		// If the string has not been seen, add it to the map and the slice of unique strings
		if !seen[str] {
			seen[str] = true
			unique = append(unique, str)
		}
	}

	return unique
}

// difference returns the elements in `a` that aren't in `b`.
func Difference(a, b []int) []int {
	mb := make(map[int]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []int
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func DifferenceUsers(a, b []int) []int {
	var diff []int
	mb := make(map[int]struct{}, len(b))
	ma := make(map[int]struct{}, len(a))
	for _, x := range b {
		mb[x] = struct{}{}
	}

	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}

	for _, x := range a {
		ma[x] = struct{}{}
	}

	for _, x := range b {
		if _, found := ma[x]; !found {
			diff = append(diff, x)
		}
	}

	return diff
}

func IsElementExistIntArray(intSlice []int, searchKey int) bool {
	for _, item := range intSlice {
		if item == searchKey {
			return true
		}
	}
	return false
}

func IsStringExist(strSlice []string, searchString string) bool {
	for _, item := range strSlice {
		if item == searchString {
			return true
		}
	}
	return false
}
