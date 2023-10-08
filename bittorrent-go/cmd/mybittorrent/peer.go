package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type trackerRequest struct {
	InfoHash   [20]byte
	PeerId     string
	Port       int
	Uploaded   int
	Downloaded int
	Left       int
	Compact    int
}

type trackerResponse struct {
	Interval int
	Peers    []peer
}

type peer struct {
	Ip   string
	Port int
}

func getPeers(announce string, request trackerRequest) (trackerResponse, error) {
	u, err := url.Parse(announce)
	if err != nil {
		return trackerResponse{}, err
	}

	q := u.Query()
	q.Set("info_hash", string(request.InfoHash[:]))
	q.Set("peer_id", request.PeerId)
	q.Set("port", strconv.Itoa(request.Port))
	q.Set("uploaded", strconv.Itoa(request.Uploaded))
	q.Set("downloaded", strconv.Itoa(request.Downloaded))
	q.Set("left", strconv.Itoa(request.Left))
	q.Set("compact", strconv.Itoa(request.Compact))
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return trackerResponse{}, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return trackerResponse{}, err
	}

	decoded, err := decodeBencode(string(b))
	if err != nil {
		return trackerResponse{}, err
	}
	dict := decoded.(map[string]interface{})

	responsePeers := dict["peers"].(string)
	peers := make([]peer, len(responsePeers)/6)
	for i := 0; i < len(responsePeers); i += 6 {
		peers[i/6].Ip = fmt.Sprintf("%d.%d.%d.%d", responsePeers[i], responsePeers[i+1], responsePeers[i+2], responsePeers[i+3])
		peers[i/6].Port = int(responsePeers[i+4])<<8 + int(responsePeers[i+5])
	}

	ret := trackerResponse{
		Interval: dict["interval"].(int),
		Peers:    peers,
	}
	return ret, nil
}
