package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

var RedirectAttemptedError = errors.New("redirect")

func main() {
	engine := gin.Default()
	engine.GET("/", func(c *gin.Context) {
		url := c.Query("url")
		if url == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"redirectLink": nil,
				"status":       "error",
				"message":      "The url parameter is not specified.",
			})
			return
		}
		if link, err := checkRedirect(url); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"status":       "success",
				"redirectLink": link,
				"message":      "Everything was completed successfully.",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"redirectLink": nil,
				"status":       "error",
				"message":      "An error has occurred in the server.",
			})
			log.Print(err)
		}
	})
	engine.Run(":80")
}

func checkRedirect(link string) (string, error) {
	client := &http.Client{
		Timeout: time.Duration(3) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return RedirectAttemptedError
		},
	}
	resp, err := client.Head(link)
	defer resp.Body.Close()
	var loc string
	if urlError, ok := err.(*url.Error); ok && urlError.Err == RedirectAttemptedError {
		loc = fmt.Sprint(resp.Header["Location"][0])
	}
	return loc, nil
}
