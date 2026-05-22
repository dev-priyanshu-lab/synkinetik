package intel

import (
	"dnslookup/common"
	"dnslookup/entity"
	"dnslookup/utils"
	"strings"
)

func checkCategory(values *[]string, value int, category int) {
	if value&category > 0 {
		*values = append(*values, common.CATEGORIES[category])
	}
}

func getCategories(category int) string {
	var values []string
	checkCategory(&values, category, common.MAL_DOMAIN)
	checkCategory(&values, category, common.DDNS_DOMAIN)
	checkCategory(&values, category, common.TUNNEL_ENDPOINTS)
	checkCategory(&values, category, common.CLICK_TRACKER)
	checkCategory(&values, category, common.SHORTERN_DOMAIN)
	checkCategory(&values, category, common.DECENTRALIZED_WEB_GATEWAY)
	checkCategory(&values, category, common.NO_SAFESEARCH)
	checkCategory(&values, category, common.DNS_BYPASS)
	checkCategory(&values, category, common.TRACKING_DOMAINS)
	checkCategory(&values, category, common.BIN_DOMAINS)
	checkCategory(&values, category, common.PIRACY_DOMAINS)
	checkCategory(&values, category, common.AD_DOMAINS)
	checkCategory(&values, category, common.PHISHING_DOMAIN)
	checkCategory(&values, category, common.BOTNET_C2)
	checkCategory(&values, category, common.DGA_DOMAIN)
	return strings.Join(values, " ")
}

func Lookup(store *entity.Storage, domain string) (bool, string) {
	category := ""
	success := true
	value := utils.GetDomain(store, domain)
	if value == 0 {
		success = false
	} else {
		category = getCategories(value)
	}

	return success, category
}
