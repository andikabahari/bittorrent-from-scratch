package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

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
		bencode, err := os.ReadFile(torrentPath)
		if err != nil {
			log.Fatalf("Error reading file: %s", err)
		}

		meta, err := newMetainfo(string(bencode))
		if err != nil {
			log.Fatalf("Error creating metainfo: %v", err)
		}
		fmt.Println("Tracker URL:", meta.Announce)
		fmt.Println("Length:", meta.Info.Length)
		fmt.Println("Info Hash:", meta.infoHash)
		fmt.Println("Piece Length:", meta.Info.PieceLength)
		fmt.Printf("Piece Hashes:\n%s\n", strings.Join(meta.pieceHashes, "\n"))

	case "peers":
		torrentPath := os.Args[2]
		bencode, err := os.ReadFile(torrentPath)
		if err != nil {
			log.Fatalf("Error reading file: %s", err)
		}

		meta, err := newMetainfo(string(bencode))
		if err != nil {
			log.Fatalf("Error creating metainfo: %v", err)
		}

		request := trackerRequest{
			InfoHash:   meta.infoHash,
			PeerId:     "00112233445566778899",
			Port:       6881,
			Uploaded:   0,
			Downloaded: 0,
			Left:       meta.Info.Length,
			Compact:    1,
		}
		response, err := getPeers(meta.Announce, request)
		if err != nil {
			log.Fatalf("Error getting peers: %v", err)
		}
		for _, peer := range response.Peers {
			fmt.Printf("%s:%d\n", peer.Ip, peer.Port)
		}

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}
