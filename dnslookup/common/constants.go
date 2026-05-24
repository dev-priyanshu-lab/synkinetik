package common

const (
	MAL_DOMAIN                = 1
	DDNS_DOMAIN               = 2
	TUNNEL_ENDPOINTS          = 4
	CLICK_TRACKER             = 16
	SHORTERN_DOMAIN           = 32
	DECENTRALIZED_WEB_GATEWAY = 64
	NO_SAFESEARCH             = 128
	DNS_BYPASS                = 256
	TRACKING_DOMAINS          = 512
	BIN_DOMAINS               = 1024
	PIRACY_DOMAINS            = 2048
	AD_DOMAINS                = 4096
	PHISHING_DOMAIN           = 8192
	BOTNET_C2                 = 16384
	DGA_DOMAIN                = 32768
	REDIS_ADDR                = "localhost:6379"
	REDIS_SOCKET              = "/var/run/redis/redis.sock"
	DNS_LOOKUP_CONFIG_FILE    = "/etc/synkinetik/lookup.conf"
	APP_LOG_FILE              = "/var/log/synkinetik/resolver/lookup.log"
)

var CATEGORIES = map[int]string{
	MAL_DOMAIN:                "Malicious Domain",
	DDNS_DOMAIN:               "DDNS Domain",
	TUNNEL_ENDPOINTS:          "Tunnel Endpoints",
	CLICK_TRACKER:             "Click Tracker",
	SHORTERN_DOMAIN:           "Shortener Domain",
	DECENTRALIZED_WEB_GATEWAY: "Decentralized Web Gateway",
	NO_SAFESEARCH:             "No SafeSearch",
	DNS_BYPASS:                "DNS Bypass",
	TRACKING_DOMAINS:          "Tracking Domains",
	BIN_DOMAINS:               "Bin Domains",
	PIRACY_DOMAINS:            "Piracy Domains",
	AD_DOMAINS:                "Ad Domains",
	PHISHING_DOMAIN:           "Phishing Domain",
	BOTNET_C2:                 "Botnet C2",
	DGA_DOMAIN:                "DGA Domain",
}
