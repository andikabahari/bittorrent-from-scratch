package main

import (
	"fmt"
	"strconv"
	"strings"
)

func encodeBencode(v interface{}) string {
	switch val := v.(type) {
	case int:
		return fmt.Sprintf("i%de", val)
	case string:
		return fmt.Sprintf("%d:%s", len(val), val)
	case []interface{}:
		elems := make([]string, len(val))
		for i, elem := range val {
			elems[i] = encodeBencode(elem)
		}
		return fmt.Sprintf("l%se", strings.Join(elems, ""))
	case map[string]interface{}:
		keys := make([]string, len(val))
		i := 0
		for k := range val {
			keys[i] = k
			i++
		}
		// sort.Strings(keys)
		pairs := make([]string, len(val))
		for _, key := range keys {
			encodedKey := encodeBencode(key)
			encodedValue := encodeBencode(val[key])
			if encodedKey != "" && encodedValue != "" {
				pairs = append(pairs, encodedKey+encodedValue)
			}
		}
		return fmt.Sprintf("d%se", strings.Join(pairs, ""))
	default:
		return ""
	}
}

func decodeBencode(bencodedString string) (interface{}, error) {
	firstChar := bencodedString[0]
	if '0' <= firstChar && firstChar <= '9' {
		return decodeString(bencodedString)
	}
	switch firstChar {
	case 'i':
		return decodeInteger(bencodedString)
	case 'l':
		return decodeList(bencodedString)
	case 'd':
		return decodeDictionary(bencodedString)
	}
	return nil, fmt.Errorf("only strings are supported at the moment")
}

func decodeString(s string) (string, error) {
	i := strings.Index(s, ":")
	length, err := strconv.Atoi(s[:i])
	if err != nil {
		return "", err
	}
	return s[i+1 : i+length+1], nil
}

func decodeInteger(s string) (int, error) {
	firstE := strings.Index(s, "e")
	return strconv.Atoi(s[1:firstE])
}

func decodeList(s string) ([]interface{}, error) {
	list := make([]interface{}, 0)
	s = s[1:]
	for len(s) > 0 && s[0] != 'e' {
		elem, err := decodeBencode(s)
		if err != nil {
			return nil, err
		}
		list = append(list, elem)
		s = s[len(encodeBencode(elem)):]
		if len(s) > 0 && s[0] == 'e' {
			break
		}
	}
	return list, nil
}

func decodeDictionary(s string) (map[string]interface{}, error) {
	dict := make(map[string]interface{})
	s = s[1:]
	for len(s) > 0 && s[0] != 'e' {
		key, err := decodeBencode(s)
		if err != nil {
			return nil, err
		}
		s = s[len(encodeBencode(key)):]
		val, err := decodeBencode(s)
		if err != nil {
			return nil, err
		}
		s = s[len(encodeBencode(val)):]
		dict[key.(string)] = val
		if len(s) > 0 && s[0] == 'e' {
			break
		}
	}
	return dict, nil
}
