package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type containers struct {
	ID         string                   `json:"Id"`
	Names      []string                 `json:"Names"`
	Image      string                   `json:"Image"`
	Command    string                   `json:"Command"`
	Created    int                      `json:"Created"`
	Status     string                   `json:"Status"`
	Ports      []map[string]interface{} `json:"Ports"`
	Labels     map[string]string        `json:"Labels"`
	SizeRw     int                      `json:"SizeRw"`
	SizeRootFs int                      `json:"SizeRootFs"`
}

type containerCheckData struct {
	Epoch int64
	Host  string
	Name  string
}

func perror(err error) {
	if err != nil {
		panic(err)
	}
}

func fakeDial(proto, addr string) (conn net.Conn, err error) {
	return net.DialTimeout("unix", "/var/run/docker.sock", time.Duration(0))
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
	//var prefix = flag.String("prefix", "docker2graphite", "Graphite prefix for metric names e.g. docker2graphite.prod")
	//flag.Parse()

	tr := &http.Transport{
		Dial: fakeDial,
	}
	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}

	//var running interface{}
	var c []containers
	err := getJSON(*client, "http://unix.sock/containers/json", &c)
	perror(err)

	//var m []containers
	if err := json.Unmarshal(running, &m); err != nil {
		panic(err)
	}
	for _, i := range m {
		fmt.Println(i.ID)
	}

	//var data interface{}
	//err := getJSON(*client, "http://unix.sock/containers/877fa0af21e6e26f271f2aeff671fe9a9025c42f431af1441a650e88b4b87bff/stats?stream=false", &data)
	//perror(err)

	//m := data.(map[string]interface{}) //[]interface{}
	//for _, i := range m {
	//	fmt.Println(prefix, i)
	//d := i.(map[string]interface{})
	//for k, v := range d {
	//	if k == "Names" {
	//		z := v.([]interface{})
	//		fmt.Printf("%v.%v.running 1\n", *prefix, safeString(z[0].(string)))
	//	}
	//}
	//}
}
