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
	Name string `json:"name"`
	Zone string `json:"zone"`
	ID   string `json:"id"`
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
	fmt.Println("Updating DNS...")
	config := loadConfig()
	ip := getIP()

	if *ip == config.IP {
		fmt.Println("IP unchanged!")
		return
	}

	config.IP = *ip
	_ = saveConfig(config)

	endpoints := loadEndpoints()

	for _, endpoint := range *endpoints {
		fmt.Println(endpoint.Name)
		err := sendData(config.Auth, *ip, endpoint)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Done!")
}

func loadConfig() *Config {
	content, err := os.ReadFile("./data/config.json")
	if err != nil {
		panic(err)
	}

	var c Config
	err = json.Unmarshal(content, &c)
	if err != nil {
		panic(err)
	}

	return &c
}

func saveConfig(config *Config) error {
	file, _ := json.MarshalIndent(config, "", " ")

	err := os.WriteFile("data/config.json", file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func loadEndpoints() *[]Endpoint {
	content, err := os.ReadFile("./data/endpoints.json")
	if err != nil {
		panic(err)
	}

	var endpoints []Endpoint
	err = json.Unmarshal(content, &endpoints)
	if err != nil {
		panic(err)
	}

	return &endpoints
}

func getIP() *string {
	resp, err := http.Get("https://ipinfo.io/ip")
	if err != nil {
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
		Proxied: false,
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
