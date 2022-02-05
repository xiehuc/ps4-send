// Package main ...
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var (
	port   = flag.Int("p", 8000, "local serve ip+port")
	remote = flag.String("r", "127.0.0.1", "ps4 remote address")
	src    = flag.String("s", "./", "src dir")
	local  string
)

type installApi struct {
	Type     string   `json:"type"`
	Packages []string `json:"packages"`
}

func main() {
	flag.Parse()
	local = localip()

	var lst []string

	err := filepath.Walk(*src, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".pkg") {
			lst = append(lst, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		time.Sleep(2 * time.Second) // wait untial http bootup
		for _, v := range lst {
			send([]string{v})
		}
	}()

	http.Handle("/", http.FileServer(http.Dir(*src)))
	if err := http.ListenAndServe(fmt.Sprint("0.0.0.0:", *port), nil); err != nil {
		log.Fatal(err)
	}
}

func send(files []string) {
	prefix := fmt.Sprint("http://", *remote, ":12800")

	var post installApi
	for i := range files {
		post.Packages = append(post.Packages, fmt.Sprint("http://", local, ":", *port, "/", files[i]))
	}
	post.Type = "direct"
	body, _ := json.Marshal(&post)

	log.Printf("curl -v '%s/api/install' --data '%s' ", prefix, string(body))

	resp, err := http.Post(prefix+"/api/install", "application/json", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	rbody, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(rbody))
}

func localip() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	addr := conn.LocalAddr().(*net.UDPAddr)
	return addr.IP.String()
}
