package main

import (
	"crypto/sha1"
	"encoding/hex"
)

type metainfo struct {
	Announce string
	Info     info

	infoHash    string
	pieceHashes []string
}

type info struct {
	Length      int
	Name        string
	PieceLength int
	Pieces      string
}

func newMetainfo(bencode string) (metainfo, error) {
	decoded, err := decodeBencode(bencode)
	if err != nil {
		return metainfo{}, err
	}

	dict1 := decoded.(map[string]interface{})
	dict2 := dict1["info"].(map[string]interface{})
	meta := metainfo{
		Announce: dict1["announce"].(string),
		Info: info{
			Length:      dict2["length"].(int),
			Name:        dict2["name"].(string),
			PieceLength: dict2["piece length"].(int),
			Pieces:      dict2["pieces"].(string),
		},
	}

	infoBencode := encodeBencode(dict2)
	checksum := sha1.Sum([]byte(infoBencode))
	meta.infoHash = hex.EncodeToString(checksum[:])

	pieceHashes := make([]string, 0, len(meta.Info.Pieces)/20)
	for i := 0; i < len(meta.Info.Pieces); i += 20 {
		pieceHash := hex.EncodeToString([]byte(meta.Info.Pieces[i : i+20]))
		pieceHashes = append(pieceHashes, pieceHash)
	}
	meta.pieceHashes = pieceHashes

	return meta, nil
}
