package helpers

import "github.com/gofiber/fiber/v2"

func GetIP(c *fiber.Ctx) string {
	var ipAddr string
	if ips := c.IPs(); len(ips) < 1 {
		ipAddr = c.IP()
	} else {
		ipAddr = ips[0]
	}
	return ipAddr
}
