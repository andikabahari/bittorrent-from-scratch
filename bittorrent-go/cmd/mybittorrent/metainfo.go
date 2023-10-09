package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"os"
)

type metainfo struct {
	Announce string
	Info     info
}

type info struct {
	Length      int
	Name        string
	PieceLength int
	Pieces      string
}

func parseTorrent(src string) (metainfo, error) {
	bencode, err := os.ReadFile(src)
	if err != nil {
		return metainfo{}, err
	}

	decoded, err := decodeBencode(string(bencode))
	if err != nil {
		return metainfo{}, err
	}

	metaDict, ok := decoded.(map[string]interface{})
	if !ok {
		return metainfo{}, errors.New("decoded value should be a dictionary")
	}

	announce, ok := metaDict["announce"].(string)
	if !ok {
		return metainfo{}, errors.New("announce not found")
	}

	infoDict, ok := metaDict["info"].(map[string]interface{})
	if !ok {
		return metainfo{}, errors.New("info not found")
	}

	length, ok := infoDict["length"].(int)
	if !ok {
		return metainfo{}, errors.New("length not found")
	}

	name, ok := infoDict["name"].(string)
	if !ok {
		return metainfo{}, errors.New("name not found")
	}

	pieceLength, ok := infoDict["piece length"].(int)
	if !ok {
		return metainfo{}, errors.New("piece length not found")
	}

	pieces, ok := infoDict["pieces"].(string)
	if !ok {
		return metainfo{}, errors.New("pieces not found")
	}

	meta := metainfo{
		Announce: announce,
		Info: info{
			Length:      length,
			Name:        name,
			PieceLength: pieceLength,
			Pieces:      pieces,
		},
	}
	return meta, nil
}

func (m *metainfo) infoHash() string {
	v := make(map[string]interface{})
	v["length"] = m.Info.Length
	v["name"] = m.Info.Name
	v["piece length"] = m.Info.PieceLength
	v["pieces"] = m.Info.Pieces
	infoBencode := encodeBencode(v)
	checksum := sha1.Sum([]byte(infoBencode))
	return hex.EncodeToString(checksum[:])
}

func (m *metainfo) pieceHashes() []string {
	pieceHashes := make([]string, 0, len(m.Info.Pieces)/20)
	for i := 0; i < len(m.Info.Pieces); i += 20 {
		pieceHash := hex.EncodeToString([]byte(m.Info.Pieces[i : i+20]))
		pieceHashes = append(pieceHashes, pieceHash)
	}
	return pieceHashes
}
