package intel

import (
	"bufio"
	"dnslookup/common"
	"dnslookup/entity"
	"dnslookup/utils"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func processLine(line string) string {
	index := strings.Index(line, ";")
	if index == 0 {
		return ""
	}

	index = strings.Index(line, "!")
	if index == 0 {
		return ""
	}

	index = strings.Index(line, "#")
	if index == 0 {
		return ""
	} else if index >= 0 {
		line = line[0:index]
	}

	index = strings.Index(line, "0.0.0.0")
	if index >= 0 {
		line = line[index+8:]
	}

	index = strings.Index(line, "127.0.0.1")
	if index >= 0 {
		line = line[index+10:]
	}

	index = strings.Index(line, "^")
	if index >= 0 {
		line = line[:index]
	}

	index = strings.Index(line, "||")
	if index >= 0 {
		line = line[index+2:]
	}

	index = strings.Index(line, "://")
	if index >= 0 {
		line = line[index+3:]
	}

	index = strings.Index(line, "/")
	if index >= 0 {
		line = line[:index]
	}

	index = strings.Index(line, "CNAME")
	if index >= 0 {
		line = line[:index]
	}

	line = strings.TrimSpace(line)
	success := utils.ValidateDomain(line)
	if success == false {
		fmt.Println("skipping: ", line)
		return ""
	}

	return line
}

func updateIntel(store *entity.Storage, filename string, link string, category int) {
	slog.Info("Updating intel", "link", link, "category", category)
	success := utils.DownloadLink(filename, link)
	if success == false {
		slog.Error("Unable to download from link", "link", link, "category", category)
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		slog.Error("Unable to read from intel", "link", link, "category", category)
		return
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := processLine(scanner.Text())
		if len(line) > 0 {
			utils.AddDomain(store, line, category)
		}
	}
}

func getCategory(category string) int {
	switch category {
	case "MAL_DOMAIN":
		{
			return common.MAL_DOMAIN
		}
	case "DDNS_DOMAIN":
		{
			return common.DDNS_DOMAIN
		}
	case "TUNNEL_ENDPOINTS":
		{
			return common.TUNNEL_ENDPOINTS
		}
	case "CLICK_TRACKER":
		{
			return common.CLICK_TRACKER
		}
	case "SHORTERN_DOMAIN":
		{
			return common.SHORTERN_DOMAIN
		}
	case "DECENTRALIZED_WEB_GATEWAY":
		{
			return common.DECENTRALIZED_WEB_GATEWAY
		}
	case "NO_SAFESEARCH":
		{
			return common.NO_SAFESEARCH
		}
	case "DNS_BYPASS":
		{
			return common.DNS_BYPASS
		}
	case "TRACKING_DOMAINS":
		{
			return common.TRACKING_DOMAINS
		}
	case "BIN_DOMAINS":
		{
			return common.BIN_DOMAINS
		}
	case "PIRACY_DOMAINS":
		{
			return common.PIRACY_DOMAINS
		}
	case "AD_DOMAINS":
		{
			return common.AD_DOMAINS
		}
	case "PHISHING_DOMAIN":
		{
			return common.PHISHING_DOMAIN
		}
	case "BOTNET_C2":
		{
			return common.BOTNET_C2
		}
	case "DGA_DOMAIN":
		{
			return common.DGA_DOMAIN
		}
	default:
		return -1
	}
}

func UpdateIntel(store *entity.Storage, config entity.Configuration) {
	for cat, source := range config.Sources {
		category := getCategory(cat)
		for index, link := range source.Feeds {
			filename := fmt.Sprintf("category_%s_feed_%d", cat, index)
			updateIntel(store, filename, link, category)
		}
	}
}
