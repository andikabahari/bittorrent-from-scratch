package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	selfPeerId        = "00112233445566778899"
	selfPort          = 6881
	protocolStringLen = 19
	protocolString    = "BitTorrent protocol"
)

var reservedBytes = []byte{0, 0, 0, 0, 0, 0, 0, 0}

func main() {
	command := os.Args[1]
	switch command {
	case "decode":
		bencodedValue := os.Args[2]
		decoded, err := decodeBencode(bencodedValue)
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
		meta, err := parseTorrent(torrentPath)
		if err != nil {
			log.Fatalf("Error parsing torrent file: %v", err)
		}
		fmt.Println("Tracker URL:", meta.Announce)
		fmt.Println("Length:", meta.Info.Length)
		fmt.Println("Info Hash:", meta.infoHash())
		fmt.Println("Piece Length:", meta.Info.PieceLength)
		fmt.Printf("Piece Hashes:\n%s\n", strings.Join(meta.pieceHashes(), "\n"))

	case "peers":
		torrentPath := os.Args[2]
		meta, err := parseTorrent(torrentPath)
		if err != nil {
			log.Fatalf("Error parsing torrent file: %v", err)
		}

		request := trackerRequest{
			InfoHash:   meta.infoHash(),
			PeerId:     selfPeerId,
			Port:       selfPort,
			Uploaded:   0,
			Downloaded: 0,
			Left:       meta.Info.Length,
			Compact:    1,
		}
		response, err := fetchTracker(meta.Announce, request)
		if err != nil {
			log.Fatalf("Error getting peers: %v", err)
		}
		for _, peer := range response.Peers {
			fmt.Printf("%s:%d\n", peer.Ip, peer.Port)
		}

	case "handshake":
		torrentPath := os.Args[2]
		meta, err := parseTorrent(torrentPath)
		if err != nil {
			log.Fatalf("Error parsing torrent file: %v", err)
		}
		infoHash, err := hex.DecodeString(meta.infoHash())
		if err != nil {
			log.Fatalf("Error decoding string: %v", err)
		}

		addr := os.Args[3]
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Fatalf("Error connecting to the addr: %v", err)
		}
		defer func() {
			err := conn.Close()
			if err != nil {
				log.Fatalf("Error closing connection: %v", err)
			}
		}()

		msg := []byte{protocolStringLen}
		msg = append(msg, []byte(protocolString)...)
		msg = append(msg, reservedBytes...)
		msg = append(msg, []byte(infoHash)...)
		msg = append(msg, []byte(selfPeerId)...)
		writeLen, err := conn.Write(msg)
		if err != nil {
			log.Fatalf("Error writing data to the connection: %v", err)
		}

		b := make([]byte, writeLen)
		_, err = conn.Read(b)
		if err != nil {
			fmt.Println(err)
			return
		}

		// 1 byte of protocol string length
		// 19 bytes of protocol string
		// 8 reserved bytes
		// 20 bytes of sha1 info hash
		offset := 1 + 19 + 8 + 20
		peerId := hex.EncodeToString(b[offset:])
		fmt.Println("Peer ID:", peerId)

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
