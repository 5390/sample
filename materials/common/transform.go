package common

import (
	"encoding/json"
	"genMaterials/log"
	"io"
	"io/ioutil"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

var appLog *log.AppLogger

func String2Int(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		appLog.Error("Convert String to Int Error. : " + err.Error())
		panic(err)
	}
	return i
}

func Int2String(i int) string {
	s := strconv.Itoa(i)
	return s
}

func ByteArray2String(requestBody io.Reader) (string, error) {
	body, err := ioutil.ReadAll(requestBody)
	if err != nil {
		appLog.Error("Convert Array to String Error. : " + err.Error())
		panic(err)
	}

	return string(body), err
}

func Map2Json(m map[string]interface{}) string {
	tmp, err := json.Marshal(m)
	if err != nil {
		appLog.Error("Convert Map to Json Error. : " + err.Error())
		panic(err)
	}
	return string(tmp)
}

func MarshalJson(v interface{}) string {
	var json = jsoniter.ConfigFastest
	buf, err := json.Marshal(v)
	if err != nil {
		appLog.Error("Json Marshall Error. : " + err.Error())
		panic(err)
	}
	// fmt.Printf("buf is", string(buf))
	return string(buf)
}

func Unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func UniqueIntList(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func StringToByte(str string) []byte {
	return []byte(str)
}

func ByteToString(bytes []byte) string {
	return string(bytes)
}
