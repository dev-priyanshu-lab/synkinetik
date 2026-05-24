package process

import (
	"dnslookup/entity"
	"dnslookup/intel"
	"encoding/json"
	"log/slog"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
)

func lookup(store *entity.Storage, domain string) entity.Response {
	var response entity.Response
	lookup := domain
	index := 0
	for {
		success, category := intel.Lookup(store, lookup)
		if success == true {
			response.Category = category
			response.Tag = []string{}
			return response
		}

		index = strings.Index(lookup, ".")
		if index == -1 {
			break
		}
		lookup = lookup[index+1:]
	}

	return response
}

func StartWebServer(store *entity.Storage, config entity.Configuration) {
	router := gin.Default()
	router.GET("/query", func(c *gin.Context) {
		if domain, found := c.GetQuery("domain"); found {
			response := lookup(store, domain)
			c.JSON(200, gin.H{
				"message": "pong",
				"success": "true",
				"data":    response,
			})
		} else {
			c.JSON(200, gin.H{
				"message": "missing domain query parameter",
				"success": "false",
			})
		}
	})

	router.Run(config.Address)
}

func StartUDPServer(store *entity.Storage, config entity.Configuration) {
	addr, err := net.ResolveUDPAddr("udp", config.Address)
	if err != nil {
		slog.Error("failed to resolve UDP address", "address", config.Address, "error", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		slog.Error("failed to start UDP server", "address", config.Address, "error", err)
		return
	}

	defer conn.Close()

	buf := make([]byte, 8192)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			slog.Error("error reading UDP packet", "error", err)
			continue
		}

		domain := strings.TrimSpace(string(buf[:n]))
		slog.Debug("received query", "client", clientAddr, "domain", domain)

		var resp map[string]interface{}
		if domain == "" {
			slog.Warn("empty domain received", "client", clientAddr)
			resp = map[string]interface{}{
				"message": "missing domain",
				"success": false,
			}
		} else {
			resp = map[string]interface{}{
				"data": lookup(store, domain),
			}
		}

		data, err := json.Marshal(resp)
		if err != nil {
			slog.Error("error marshalling response", "client", clientAddr, "error", err)
			continue
		}

		if _, err := conn.WriteToUDP(data, clientAddr); err != nil {
			slog.Error("error sending response", "client", clientAddr, "error", err)
		}
	}
}
