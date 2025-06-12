package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Image struct {
	ID  int    `json:"id"`
	Src string `json:"src"`
}

type Product struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Price        string   `json:"price"`
	SalePrice    string   `json:"sale_price"`
	RegularPrice string   `json:"regular_price"`
	Images       []string `json:"images"`
	Categories   []string `json:"categories"`
	Description  string   `json:"description"`
	ShortDesc    string   `json:"short_description"`
	StockStatus  string   `json:"stock_status"`
	SKU          string   `json:"sku"`
	OnSale       bool     `json:"on_sale"`
	Permalink    string   `json:"permalink"`
}

func GetProducts(c *gin.Context) {
	// Load environment variables
	key := os.Getenv("WOOCOMMERCE_KEY")
	secret := os.Getenv("WOOCOMMERCE_SECRET")
	baseURL := os.Getenv("WOOCOMMERCE_URL")

	// fmt.Println("### Woo Key:", key)
	// fmt.Println("### Woo Secret:", secret)
	// fmt.Println("### Woo URL:", baseURL)

	url := fmt.Sprintf("%s/wp-json/wc/v3/products", baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.SetBasicAuth(key, secret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch products"})
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	var rawProducts []map[string]interface{} // raw response before cleaning
	err = json.Unmarshal(body, &rawProducts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	var products []Product
	for _, item := range rawProducts {
		var imageUrls []string
		if images, ok := item["images"].([]interface{}); ok {
			for _, img := range images {
				if imgMap, ok := img.(map[string]interface{}); ok {
					if src, ok := imgMap["src"].(string); ok {
						imageUrls = append(imageUrls, src)
					}
				}
			}
		}

		// Clean up the product data
		var categories []string
		if cats, ok := item["categories"].([]interface{}); ok {
			for _, cat := range cats {
				if catMap, ok := cat.(map[string]interface{}); ok {
					categories = append(categories, catMap["name"].(string))
				}
			}
		}

		product := Product{
			ID:           int(item["id"].(float64)),
			Name:         item["name"].(string),
			Price:        item["price"].(string),
			SalePrice:    item["sale_price"].(string),
			RegularPrice: item["regular_price"].(string),
			Images:       imageUrls,
			Categories:   categories,
			Description:  item["description"].(string),
			ShortDesc:    item["short_description"].(string),
			StockStatus:  item["stock_status"].(string),
			SKU:          item["sku"].(string),
			OnSale:       item["on_sale"].(bool),
			Permalink:    item["permalink"].(string),
		}

		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}
