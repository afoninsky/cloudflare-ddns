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
	domainWild := fmt.Sprintf("*.%s", domain)
	if domain == "" {
		panic(errors.New("domain is not specified in CLOUDFLARE_ZONE"))
	}
	ip, err := getIP()
	if err != nil {
		panic(err)
	}

	api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	zoneID, err := api.ZoneIDByName(domain)
	if err != nil {
		panic(err)
	}

	records, err := api.DNSRecords(ctx, zoneID, cloudflare.DNSRecord{
		Type: "A",
	})
	if err != nil {
		panic(err)
	}
	found := false
	for _, rec := range records {
		if rec.Name != domain && rec.Name != domainWild {
			continue
		}
		found = true
		if rec.Content == ip {
			continue
		}
		rec.Content = ip
		rec.ModifiedOn = time.Now()
		if err := api.UpdateDNSRecord(ctx, zoneID, rec.ID, rec); err != nil {
			panic(err)
		}
	}
	if !found {
		panic(fmt.Sprintf("neither %s not %s not found", domain, domainWild))
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
