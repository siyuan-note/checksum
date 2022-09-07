package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	p := "."
	entries, err := os.ReadDir(p)
	if nil != err {
		log.Fatalf("read dir failed: %s", err)
	}

	buf := bytes.Buffer{}
	for _, entry := range entries {
		filename := filepath.Join(p, entry.Name())
		if isSkip(filename) {
			continue
		}

		hash, hashErr := sha256Hash(filename)
		if nil != hashErr {
			log.Fatalf("get hash failed: %s", hashErr)
		}
		buf.WriteString(hash + " " + entry.Name() + "\n")
	}

	log.Printf("\n%s", buf.String())
	if err = os.WriteFile(filepath.Join(p, "SHA256SUMS.txt"), buf.Bytes(), 0644); nil != err {
		log.Fatalf("write file failed: %s", err)
	}
}

func sha256Hash(filename string) (ret string, err error) {
	file, err := os.Open(filename)
	if nil != err {
		return
	}
	defer file.Close()

	hash := sha256.New()
	reader := bufio.NewReader(file)
	buf := make([]byte, 1024*1024*4)
	for {
		switch n, readErr := reader.Read(buf); readErr {
		case nil:
			hash.Write(buf[:n])
		case io.EOF:
			return fmt.Sprintf("%x", hash.Sum(nil)), nil
		default:
			return "", err
		}
	}
}

func isSkip(filename string) bool {
	info, err := os.Stat(filename)
	if nil != err {
		log.Fatalf("stat file [%s] failed: %s", filename, err)
	}
	if info.IsDir() {
		return true
	}

	base := filepath.Base(filename)
	if strings.HasPrefix(base, ".") || "SHA256SUMS.txt" == base || "checksum" == base || "checksum.exe" == base {
		return true
	}
	return false
}
