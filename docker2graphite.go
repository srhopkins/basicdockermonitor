package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
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

func safeString(s string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Fatalln(err)
	}
	safe := reg.ReplaceAllString(s, "_")
	safe = strings.ToLower(strings.Trim(safe, "_"))
	return safe
}

func main() {
	var prefix = flag.String("prefix", "docker2graphite", "Graphite prefix for metric names e.g. docker2graphite.prod")
	flag.Parse()

	tr := &http.Transport{
		Dial: fakeDial,
	}
	client := &http.Client{Transport: tr}

	var data interface{}
	err := getJSON(*client, "http://unix.sock/containers/json", &data)
	perror(err)

	m := data.([]interface{})
	for _, i := range m {
		d := i.(map[string]interface{})
		for k, v := range d {
			if k == "Names" {
				z := v.([]interface{})
				fmt.Printf("%v.%v.running 1\n", *prefix, safeString(z[0].(string)))
			}
		}
	}
}
