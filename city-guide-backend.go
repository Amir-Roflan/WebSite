package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

const yandexAPIKey = "f4fce057-d47d-487e-bd40-3edc90203535" // <-- сюда вставь свой ключ

func main() {
	r := gin.Default()

	r.GET("/coords", func(c *gin.Context) {
		city := c.Query("city")
		if city == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "city parameter is required"})
			return
		}

		geoURL := "https://geocode-maps.yandex.ru/1.x/"
		params := url.Values{}
		params.Set("apikey", yandexAPIKey)
		params.Set("format", "json")
		params.Set("geocode", city)

		fullURL := geoURL + "?" + params.Encode()

		resp, err := http.Get(fullURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch from Yandex"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read response body"})
			return
		}

		fmt.Println("🟡 Ответ от Яндекса:")
		fmt.Println(string(body)) // Для отладки, можно убрать после теста

		var geoResp struct {
			Response struct {
				GeoObjectCollection struct {
					FeatureMember []struct {
						GeoObject struct {
							Point struct {
								Pos string `json:"pos"`
							} `json:"Point"`
						} `json:"GeoObject"`
					} `json:"featureMember"`
				} `json:"GeoObjectCollection"`
			} `json:"response"`
		}

		err = json.Unmarshal(body, &geoResp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse JSON from Yandex"})
			return
		}

		if len(geoResp.Response.GeoObjectCollection.FeatureMember) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "city not found"})
			return
		}

		coords := geoResp.Response.GeoObjectCollection.FeatureMember[0].GeoObject.Point.Pos
		var lon, lat string
		_, err = fmt.Sscanf(coords, "%s %s", &lon, &lat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse coordinates"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"lat": lat,
			"lon": lon,
		})
	})

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
