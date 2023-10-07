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

	case "info":
		torrentPath := os.Args[2]
		bencoded, err := os.ReadFile(torrentPath)
		if err != nil {
			log.Fatalf("Error reading file: %s", err)
		}

		decoded, _, err := decodeBencode(string(bencoded))
		if err != nil {
			log.Fatalf("Error decoding bencode: %v", err)
		}

		jsonData, err := json.Marshal(decoded)
		if err != nil {
			log.Fatalf("Error marshaling json: %v", err)
		}

		meta := metainfo{}
		err = json.Unmarshal(jsonData, &meta)
		if err != nil {
			log.Fatalf("Error unmarshaling json: %v", err)
		}
		fmt.Printf("Tracker URL: %s\n", meta.Announce)
		fmt.Printf("Length: %d\n", meta.Info.Length)

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

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
		return nil, -1, fmt.Errorf("only strings are supported at the moment")
	}
}

func decodeString(s string) (string, int, error) {
	i := strings.Index(s, ":")
	length, err := strconv.Atoi(s[:i])
	if err != nil {
		return "", -1, err
	}
	return s[i+1 : i+length+1], i + length, nil
}

func decodeInteger(s string) (int, int, error) {
	i := strings.Index(s, "e")
	integer, err := strconv.Atoi(s[1:i])
	if err != nil {
		return 0, -1, err
	}
	return integer, i, err
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

type metainfo struct {
	Announce  string `json:"announce"`
	CreatedBy string `json:"created by"`
	Info      struct {
		Length      int    `json:"length"`
		Name        string `json:"name"`
		PieceLength int    `json:"piece length"`
		Pieces      string `json:"pieces"`
	} `json:"info"`
}
