#!/bin/sh

cd bittorrent-tester
go build -o ../bittorrent-go/test.out ./cmd/tester
cp -r ./torrents ../bittorrent-go/
cp -r ./response ../bittorrent-go/

cd ../bittorrent-go
CODECRAFTERS_SUBMISSION_DIR=$(pwd) \
CODECRAFTERS_TEST_CASES_JSON=`cat ../test_cases.json` \
./test.out

rm ./test.out
rm -r ./torrents
rm -r ./response