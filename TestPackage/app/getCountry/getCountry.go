package getCountry

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type IPInfo struct {
	IP        string  `json:"ip"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
	Region    string  `json:"region"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (c *Cache) GetLocationInfo(ip string) (string, error) {
	// Check the cache first
	if locationInfo, exists := c.Get(ip); exists {
		return locationInfo, nil
	}

	url := fmt.Sprintf("https://ipwho.is/%s", ip)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get IP info: %s", resp.Status)
	}

	var ipInfo IPInfo
	if err = json.NewDecoder(resp.Body).Decode(&ipInfo); err != nil {
		return "", err
	}

	locationInfo := fmt.Sprintf("%s/%s/%s", ipInfo.Country, ipInfo.Region, ipInfo.City)

	// Store the result in the cache
	c.Set(ip, locationInfo)

	return locationInfo, nil
}

func GetAddress() {
	c := NewCache()
	ip := "8.8.8.8" // Example IP address
	locationInfo, err := c.GetLocationInfo(ip)
	if err != nil {
		log.Fatalf("Error getting location info: %v", err)
	}

	fmt.Printf("Location Info: %s\n", locationInfo)
}
