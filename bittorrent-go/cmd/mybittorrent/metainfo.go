package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"

	"github.com/jackpal/bencode-go"
)

type metainfo struct {
	Announce  string `json:"announce" bencode:"announce"`
	CreatedBy string `json:"created by" bencode:"created by"`
	Info      info   `json:"info" bencode:"info"`

	infoHash    [20]byte
	pieceHashes []string
}

type info struct {
	Length      int    `json:"length" bencode:"length"`
	Name        string `json:"name" bencode:"name"`
	PieceLength int    `json:"piece length" bencode:"piece length"`
	Pieces      string `json:"pieces" bencode:"pieces"`
}

func newMetainfo(r io.Reader) (metainfo, error) {
	meta := metainfo{}
	err := bencode.Unmarshal(r, &meta)
	if err != nil {
		return metainfo{}, err
	}

	var buf bytes.Buffer
	err = bencode.Marshal(&buf, meta.Info)
	if err != nil {
		log.Fatalf("Error marshaling bencode: %v", err)
	}
	checksum := sha1.Sum(buf.Bytes())
	meta.infoHash = checksum

	pieceHashes := make([]string, 0, len(meta.Info.Pieces)/20)
	for i := 0; i < len(meta.Info.Pieces); i += 20 {
		pieceHash := hex.EncodeToString([]byte(meta.Info.Pieces[i : i+20]))
		pieceHashes = append(pieceHashes, pieceHash)
	}
	meta.pieceHashes = pieceHashes

	return meta, nil
}
