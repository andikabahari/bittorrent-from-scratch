package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	switch command := os.Args[1]; command {
	case "decode":
		decodeCmd()
	case "info":
		infoCmd()
	case "peers":
		peersCmd()
	case "handshake":
		handshakeCmd()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s", command)
		os.Exit(1)
	}
}

func decodeCmd() {
	bencodedValue := os.Args[2]
	decoded, err := decodeBencode(bencodedValue)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding bencode: %v", err)
		os.Exit(1)
	}

	jsonOutput, err := json.Marshal(decoded)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling json: %v", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonOutput))
}

func infoCmd() {
	torrentPath := os.Args[2]
	meta, err := parseTorrent(torrentPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing torrent file: %v", err)
		os.Exit(1)
	}
	fmt.Println("Tracker URL:", meta.Announce)
	fmt.Println("Length:", meta.Info.Length)
	fmt.Println("Info Hash:", meta.infoHash())
	fmt.Println("Piece Length:", meta.Info.PieceLength)
	fmt.Printf("Piece Hashes:\n%s\n", strings.Join(meta.pieceHashes(), "\n"))
}

func peersCmd() {
	torrentPath := os.Args[2]
	meta, err := parseTorrent(torrentPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing torrent file: %v", err)
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "Error getting peers: %v", err)
		os.Exit(1)
	}
	for _, peer := range response.Peers {
		fmt.Printf("%s:%d\n", peer.Ip, peer.Port)
	}
}

func handshakeCmd() {
	torrentPath := os.Args[2]
	meta, err := parseTorrent(torrentPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing torrent file: %v", err)
		os.Exit(1)
	}
	infoHash, err := hex.DecodeString(meta.infoHash())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding string: %v", err)
		os.Exit(1)
	}

	addr := os.Args[3]
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to the addr: %v", err)
		os.Exit(1)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error closing connection: %v", err)
			os.Exit(1)
		}
	}()

	msg := []byte{protocolStringLen}
	msg = append(msg, []byte(protocolString)...)
	msg = append(msg, reservedBytes...)
	msg = append(msg, []byte(infoHash)...)
	msg = append(msg, []byte(selfPeerId)...)
	writeLen, err := conn.Write(msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing data to the connection: %v", err)
		os.Exit(1)
	}

	b := make([]byte, writeLen)
	_, err = conn.Read(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading data from the connection: %v", err)
		os.Exit(1)
	}

	// 1 byte of protocol string length
	// 19 bytes of protocol string
	// 8 reserved bytes
	// 20 bytes of sha1 info hash
	offset := 1 + 19 + 8 + 20
	peerId := hex.EncodeToString(b[offset:])
	fmt.Println("Peer ID:", peerId)
}
