package utils

import (
	"dnslookup/common"
	"dnslookup/entity"
	"strconv"
	"strings"
)

func GetDomain(store *entity.Storage, domain string) int {
	value, err := store.Client.Get(common.Context, domain).Result()
	if err != nil {
		return 0
	}

	category, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return category
}

func isTldLevelOneWildcard(domain string) bool {
	d := domain
	if strings.HasPrefix(d, "*.") {
		d = d[2:]
	}

	return !strings.Contains(d, ".")
}

func AddDomain(store *entity.Storage, domain string, category int) bool {
	if isTldLevelOneWildcard(domain) {
		return false
	}

	value := GetDomain(store, domain)
	category = category | value
	err := store.Client.Set(common.Context, domain, category, 0).Err()
	if err != nil {
		return false
	}

	return true
}
