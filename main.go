package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"golang.org/x/net/context"
)

const timeout = 5 * time.Second

func main() {
	domain := os.Getenv("CLOUDFLARE_ZONE")
	if domain == "" {
		panic(errors.New("domain is not specified in CLOUDFLARE_ZONE"))
	}
	ip, err := getIP()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Setting IP %s for domain %s ...\n", ip, domain)

	api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	// Fetch the zone ID
	zoneID, err := api.ZoneIDByName(domain)
	if err != nil {
		panic(err)
	}

	records, err := api.DNSRecords(ctx, zoneID, cloudflare.DNSRecord{
		Type: "A",
		Name: domain,
	})
	if err != nil {
		panic(err)
	}
	if len(records) != 1 {
		panic(fmt.Errorf("need to create root domain %s", domain))
	}
	rec := records[0]
	if rec.Content != ip {
		rec.Content = ip
		rec.ModifiedOn = time.Now()
		if err := api.UpdateDNSRecord(ctx, zoneID, rec.ID, rec); err != nil {
			panic(err)
		}
		fmt.Println("Updated!!!")
	} else {
		fmt.Println("No updates")
	}
}

func getIP() (string, error) {
	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Get("http://ip-api.com/json")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	var ip struct {
		Query  string
		Status string
	}
	if err = json.Unmarshal(body, &ip); err != nil {
		return "", err
	}
	if ip.Status != "success" || ip.Query == "" {
		return "", fmt.Errorf("invalid response format")
	}
	return ip.Query, nil
}
