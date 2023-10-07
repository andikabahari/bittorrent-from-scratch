package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	command := os.Args[1]
	switch command {
	case "decode":
		bencodedValue := os.Args[2]
		decoded, _, err := decodeBencode(bencodedValue)
		if err != nil {
			log.Fatalf("Error decoding bencode: %v", err)
		}

		jsonOutput, err := json.Marshal(decoded)
		if err != nil {
			log.Fatalf("Error marshaling json: %v", err)
		}
		fmt.Println(string(jsonOutput))

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func decodeBencode(bencodedString string) (interface{}, int, error) {
	if unicode.IsDigit(rune(bencodedString[0])) {
		return decodeString(bencodedString)
	} else if rune(bencodedString[0]) == 'i' {
		return decodeInteger(bencodedString)
	} else if rune(bencodedString[0]) == 'l' {
		return decodeList(bencodedString)
	} else if rune(bencodedString[0]) == 'd' {
		return decodeDictionary(bencodedString)
	} else {
		fmt.Printf("DI SINI %s\n", string(bencodedString[0]))
		return nil, -1, fmt.Errorf("only strings are supported at the moment")
	}
}

func decodeString(s string) (string, int, error) {
	i := strings.Index(s, ":")
	length, err := strconv.Atoi(s[:i])
	if err != nil {
		return "", -1, err
	}
	offset := i + length
	return s[i+1 : offset+1], offset, nil
}

func decodeInteger(s string) (int, int, error) {
	i := strings.Index(s, "e")
	integer, err := strconv.Atoi(s[1:i])
	if err != nil {
		return 0, -1, err
	}
	offset := i
	return integer, offset, err
}

func decodeList(s string) ([]interface{}, int, error) {
	list := make([]interface{}, 0)
	i := 1
	for i < len(s) && s[i] != 'e' {
		data, offset, err := decodeBencode(s[i:])
		if err != nil {
			return nil, -1, err
		}
		list = append(list, data)
		i += offset + 1
	}
	return list, i, nil
}

func decodeDictionary(s string) (map[string]interface{}, int, error) {
	dict := make(map[string]interface{})
	i := 1
	for i < len(s) && s[i] != 'e' {
		key, offset, err := decodeString(s[i:])
		if err != nil {
			return nil, -1, err
		}
		i += offset + 1
		value, offset, err := decodeBencode(s[i:])
		if err != nil {
			return nil, -1, err
		}
		i += offset + 1
		dict[key] = value
	}
	return dict, i, nil
}
