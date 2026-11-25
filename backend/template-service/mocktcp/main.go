// Package main is a generic mock TCP server for non-HTTP downstreams —
// Mockoon's equivalent for raw-socket protocols (see
// document/MOCKING_GUIDE.md §3). Trimmed and genericized from
// as400-adapter's mockas400/main.go (379 lines, EBCDIC/MBASE-specific) —
// go read that file directly if your protocol is actually AS400 MBASE;
// this version only keeps the pattern that generalizes: read a
// fixed-length header, extract an operation code at a known offset,
// return a canned response for that code, round-robin through multiple
// fixtures per code if more than one exists.
//
// Usage:
//
//	go run ./mocktcp/ -port 9100 -fixtures ./mocktcp/fixtures -code-offset 0 -code-length 4
package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

type opResponses struct {
	responses [][]byte
	counter   atomic.Int64
}

func (r *opResponses) next() []byte {
	idx := r.counter.Add(1) - 1
	return r.responses[idx%int64(len(r.responses))]
}

type fixtureStore struct {
	mu   sync.RWMutex
	data map[string]*opResponses
}

// loadFixtures reads fixtures/<op-code>/response.txt (single) or
// fixtures/<op-code>/response_N.txt (round-robin, N starting at 1).
func (s *fixtureStore) load(dir string) error {
	s.data = make(map[string]*opResponses)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		opCode := e.Name()
		files, err := filepath.Glob(filepath.Join(dir, opCode, "response*.txt"))
		if err != nil || len(files) == 0 {
			continue
		}
		sort.Strings(files)
		var responses [][]byte
		for _, f := range files {
			data, err := os.ReadFile(f)
			if err != nil {
				return err
			}
			responses = append(responses, data)
		}
		s.data[opCode] = &opResponses{responses: responses}
	}
	return nil
}

func (s *fixtureStore) responseFor(opCode string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.data[opCode]
	if !ok {
		return nil, false
	}
	return r.next(), true
}

func main() {
	port := flag.String("port", "9100", "TCP port to listen on")
	fixturesDir := flag.String("fixtures", "./mocktcp/fixtures", "Path to fixture directory, one subfolder per op code")
	codeOffset := flag.Int("code-offset", 0, "Byte offset of the operation code within each frame")
	codeLength := flag.Int("code-length", 4, "Byte length of the operation code")
	flag.Parse()

	store := &fixtureStore{}
	if err := store.load(*fixturesDir); err != nil {
		log.Fatalf("failed to load fixtures from %s: %v", *fixturesDir, err)
	}
	log.Printf("mock TCP server listening on :%s (%d op codes loaded)", *port, len(store.data))

	ln, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}
	defer ln.Close() //nolint:errcheck

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go handleConn(conn, store, *codeOffset, *codeLength)
	}
}

func handleConn(conn net.Conn, store *fixtureStore, codeOffset, codeLength int) {
	defer conn.Close() //nolint:errcheck

	buf := make([]byte, 64*1024)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		log.Printf("read error: %v", err)
		return
	}
	if n < codeOffset+codeLength {
		log.Printf("frame too short to contain an op code (%d bytes)", n)
		return
	}

	opCode := strings.TrimSpace(string(buf[codeOffset : codeOffset+codeLength]))
	resp, ok := store.responseFor(opCode)
	if !ok {
		log.Printf("no fixture for op code %q", opCode)
		return
	}
	if _, err := conn.Write(resp); err != nil {
		log.Printf("write error: %v", err)
	}
}
