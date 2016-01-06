package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

func perror(err error) {
	if err != nil {
		panic(err)
	}
}

func fakeDial(proto, addr string) (conn net.Conn, err error) {
	return net.Dial("unix", "/var/run/docker.sock")
}

func getJSON(c http.Client, url string, target *interface{}) error {
	r, err := c.Get(url)
	perror(err)
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(&target)
}

func main() {
	tr := &http.Transport{
		Dial: fakeDial,
	}
	client := &http.Client{Transport: tr}

	var data interface{}
	err := getJSON(*client, "http://unix.sock/version", &data) //containers/json
	perror(err)
	fmt.Println(data)
}
