package network

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Status represents the status of network-related checks
type Status struct {
	InternetConnected bool    `json:"internetConnected"`
	PingLatency       float64 `json:"pingLatency"`
	PingStatus        string  `json:"pingStatus"`
	DNSStatus         string  `json:"dnsStatus"`
}

// GetStatus checks internet connectivity, ping latency, and DNS resolution
func GetStatus() Status {
	status := Status{
		InternetConnected: false,
		PingLatency:       0,
		PingStatus:        "Not Available",
		DNSStatus:         "Not Available",
	}

	// Check internet connectivity
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("https://www.google.com")
	if err == nil && resp.StatusCode == 200 {
		status.InternetConnected = true
	}

	// Check ping latency to Google
	pingCmd := exec.Command("ping", "-c", "1", "google.com")
	if output, err := pingCmd.Output(); err == nil {
		// Parse ping output to get latency
		outputStr := string(output)
		if strings.Contains(outputStr, "time=") {
			// Extract time value from ping output
			timeStr := strings.Split(strings.Split(outputStr, "time=")[1], " ")[0]
			// Convert to float64
			if latency, err := parsePingTime(timeStr); err == nil {
				status.PingLatency = latency
				status.PingStatus = fmt.Sprintf("%.1f ms", latency)
			}
		}
	}

	// Check DNS resolution
	if _, err := net.LookupHost("google.com"); err == nil {
		status.DNSStatus = "Working"
	} else {
		status.DNSStatus = "Failed"
	}

	return status
}

// parsePingTime converts ping time string to float64
func parsePingTime(timeStr string) (float64, error) {
	// Remove any non-numeric characters except decimal point
	cleanStr := strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' {
			return r
		}
		return -1
	}, timeStr)

	return strconv.ParseFloat(cleanStr, 64)
}
