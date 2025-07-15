package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Auth string `json:"auth"`
	IP   string `json:"ip"`
}

type Endpoint struct {
	Name    string `json:"name"`
	Zone    string `json:"zone"`
	ID      string `json:"id"`
	Proxied bool   `json:"proxied"`
}

type RequestData struct {
	Content string   `json:"content"`
	Name    string   `json:"name"`
	Proxied bool     `json:"proxied"`
	Type    string   `json:"type"`
	Comment string   `json:"comment"`
	ID      string   `json:"id"`
	Tags    []string `json:"tags"`
	TTL     int      `json:"ttl"`
}

func main() {
	path, err := os.Executable()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	path = path[:len(path)-8]
	fmt.Println("path:", path)

	config := loadConfig(path)
	ip := getIP()

	if *ip == config.IP {
		fmt.Println("IP unchanged!")
		return
	}

	fmt.Println("Updating DNS...")

	config.IP = *ip
	err = saveConfig(config, path)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	endpoints := loadEndpoints(path)

	for _, endpoint := range *endpoints {
		fmt.Println(endpoint.Name)
		err := sendData(config.Auth, *ip, endpoint)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
	}
	fmt.Println("Done!")
}

func loadConfig(wd string) *Config {
	content, err := os.ReadFile(wd + "/data/config.json")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var c Config
	err = json.Unmarshal(content, &c)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return &c
}

func saveConfig(config *Config, wd string) error {
	file, _ := json.MarshalIndent(config, "", " ")

	err := os.WriteFile(wd+"/data/config.json", file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func loadEndpoints(wd string) *[]Endpoint {
	content, err := os.ReadFile(wd + "/data/endpoints.json")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var endpoints []Endpoint
	err = json.Unmarshal(content, &endpoints)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return &endpoints
}

func getIP() *string {
	resp, err := http.Get("https://ipinfo.io/ip")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	ip := string(body)
	return &ip
}

func sendData(auth string, ip string, endpoint Endpoint) error {
	url := "https://api.cloudflare.com/client/v4/zones/" + endpoint.Zone + "/dns_records/" + endpoint.ID
	currentTime := time.Now()
	data := RequestData{
		Content: ip,
		Name:    endpoint.Name,
		Proxied: endpoint.Proxied,
		Type:    "A",
		Comment: "Mercury::" + currentTime.Format("2006.01.02 15:04:05"),
		ID:      endpoint.ID,
		Tags:    []string{},
		TTL:     1,
	}

	bdata, _ := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bdata))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
