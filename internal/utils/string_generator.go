package utils

import (
	"fmt"
	"strings"
	"time"
)

func GenerateInvoiceCode() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("INV-%d", timestamp)
}

func GenerateShopName(shopName string) string {
	if strings.Contains(shopName, "@") {
		parts := strings.Split(shopName, "@")
		if len(parts) > 0 {
			return "toko-" + parts[0]
		}
	}

	return strings.ReplaceAll(shopName, " ", "")
}
