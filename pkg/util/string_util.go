package util

/*
   Author 		: Sri
   Email 		: k.srijeyanthan@beat.com
   Description : Common class to handle random token generations and encrypt or decrypt text
*/

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"math"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var zAlphaNumeric = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
var numeric = []rune("1234567890")

func FnGenerateNRandomANString(_iNumberOfCharacters int) string {
	rand.Seed(time.Now().UnixNano())
	zGeneratedBytes := make([]rune, _iNumberOfCharacters)
	for i := range zGeneratedBytes {
		zGeneratedBytes[i] = zAlphaNumeric[rand.Intn(len(zAlphaNumeric))]
	}
	return string(zGeneratedBytes)
}

func FnGenerateNString(_iNumberOfCharacters int) string {
	rand.Seed(time.Now().UnixNano())
	zGeneratedBytes := make([]rune, _iNumberOfCharacters)
	for i := range zGeneratedBytes {
		zGeneratedBytes[i] = numeric[rand.Intn(len(numeric))]
	}
	return string(zGeneratedBytes)
}

func FnGenerateNStringWithPrefix(prefix string, _iNumberOfCharacters int) string {
	rand.Seed(time.Now().UnixNano())
	zGeneratedBytes := make([]rune, _iNumberOfCharacters)
	for i := range zGeneratedBytes {
		zGeneratedBytes[i] = numeric[rand.Intn(len(numeric))]
	}
	return prefix + string(zGeneratedBytes)
}
func FnGenerateNRandomString(_iNumberOfCharacters int) string {
	rand.Seed(time.Now().UnixNano())
	zGeneratedBytes := make([]rune, _iNumberOfCharacters)
	for i := range zGeneratedBytes {
		zGeneratedBytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(zGeneratedBytes)
}

func Has(targetString string, listOfString []string) bool {
	for _, element := range listOfString {
		if element == targetString {
			return true
		}
	}
	return false
}

func HasInt(targetInt int, listOfInt []int) bool {
	for _, element := range listOfInt {
		if element == targetInt {
			return true
		}
	}
	return false
}
func GenerateBasicAuth(userName, password string) string {
	auth := userName + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func GetMD5Hash(inputData string) string {
	hash := md5.Sum([]byte(inputData))
	return hex.EncodeToString(hash[:])
}

func GetMD5OfUUID() string {
	inputData := uuid.New().String()
	hash := md5.Sum([]byte(inputData))
	return hex.EncodeToString(hash[:])
}
func GetLineCount(inputData string) int {
	iNumberOfLines := strings.Count(inputData, "\n")
	if len(inputData) > 0 && !strings.HasSuffix(inputData, "\n") {
		iNumberOfLines++
	}
	return iNumberOfLines
}

// Contains tells whether a contains x.
func JSONArrayContains(inputByteArray []byte, x string) bool {
	var arrayOfString []string
	json.Unmarshal(inputByteArray, &arrayOfString)
	for _, n := range arrayOfString {
		if x == n {
			return true
		}
	}
	return false
}

// Contains tells whether a contains x.
func IsSuperAdmin(userRole string) bool {
	if userRole == "super-admin" {
		return true
	}
	return false
}

// Contains tells whether a contains x.
func IsAdmin(userRole string) bool {
	if userRole == "admin" {
		return true
	}
	return false
}

// Contains tells whether a contains x.
func StringArrayContains(arrayList []string, element string) bool {
	for _, n := range arrayList {
		if element == n {
			return true
		}
	}
	return false
}

// Contains tells whether a contains x.
func StringContainsWithPos(arrayList []string, element string) (bool, int) {

	for position, n := range arrayList {
		if element == n {
			return true, position
		}
	}
	return false, -1
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func GetLastString(ss []string) string {
	return ss[len(ss)-1]
}

func ConnectWordWithUnderscore(s string) string {
	rr := make([]rune, 0, len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			rr = append(rr, r)
		} else {
			rr = append(rr, '_')
		}
	}
	return string(rr)
}

// ToSnake convert the given string to snake case following the Golang format:
// acronyms are converted to lower-case and preceded by an underscore.
func ToSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}

// const below are used in packages
const (
	OddError       = "odd number rule provided please provide in even count"
	SelectCapital  = "([a-z])([A-Z])"
	ReplaceCapital = "$1 $2"
)

func caseHelper(input string, isCamel bool, rule ...string) []string {
	if !isCamel {
		re := regexp.MustCompile(SelectCapital)
		input = re.ReplaceAllString(input, ReplaceCapital)
	}
	input = strings.Join(strings.Fields(strings.TrimSpace(input)), " ")
	if len(rule) > 0 && len(rule)%2 != 0 {
		panic(errors.New(OddError))
	}
	rule = append(rule, ".", " ", "_", " ", "-", " ")

	replacer := strings.NewReplacer(rule...)
	input = replacer.Replace(input)
	words := strings.Fields(input)
	return words
}

func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
func CamelCase(input string, rule ...string) string {

	// removing excess space
	wordArray := caseHelper(input, true, rule...)
	for i, word := range wordArray {
		wordArray[i] = UcFirst(word)
	}
	camelString := strings.Join(wordArray, "")
	// this is what programming need
	return LcFirst(camelString)
}
func CalculatePercentage(number1, number2 int) float64 {
	if number2 == 0 {
		return 0 // Avoid division by zero
	}
	// Calculate the percentage
	percentage := (float64(number1) / float64(number2)) * 100

	// Round down to two decimal places
	return math.Floor(percentage*100) / 100
}
func ToLowerCase(value string) string {
	return strings.ToLower(value)
}