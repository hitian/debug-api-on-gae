// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

// This function's name is a must. App Engine uses it to drive the requests properly.
func init() {
	// Starts a new Gin instance with no middle-ware
	r := gin.New()

	// Define your handlers
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/ip", func(c *gin.Context) {
		ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
		c.String(http.StatusOK, ip)
	})

	r.GET("/ua", func(c *gin.Context) {
		c.String(http.StatusOK, c.GetHeader("User-Agent"))
	})

	r.GET("/headers", func(c *gin.Context) {
		var headers = Resp{}
		for headerKey, headerValue := range c.Request.Header {
			if strings.HasPrefix(headerKey, "X-Appengine") || strings.HasPrefix(headerKey, "X-Cloud") || strings.HasPrefix(headerKey, "X-Google") {
				continue
			}
			headers[headerKey] = strings.Join(headerValue, ",")
		}
		c.JSON(http.StatusOK, headers)
	})

	r.GET("/loc", func(c *gin.Context) {
		var resp = Resp{}
		resp["city"] = c.GetHeader("X-Appengine-City")
		resp["citylatlong"] = c.GetHeader("X-Appengine-Citylatlong")
		resp["country"] = c.GetHeader("X-Appengine-Country")
		c.JSON(http.StatusOK, resp)
	})

	r.POST("/post", func(c *gin.Context) {
		var headers = Resp{}
		for headerKey, headerValue := range c.Request.Header {
			headers[headerKey] = strings.Join(headerValue, ",")
		}
		body, err := ioutil.ReadAll(c.Request.Body)
		if err == nil {
			headers["body"] = string(body)
		}
		c.JSON(http.StatusOK, headers)
	})

	r.GET("/date", func(c *gin.Context) {
		c.String(http.StatusOK, time.Now().Format(time.RFC3339))
	})

	r.GET("/timestamp", func(c *gin.Context) {
		c.String(http.StatusOK, fmt.Sprintf("%d", time.Now().Unix()))
	})

	r.GET("/check_status", func(c *gin.Context) {
		type respType map[string]int
		var resp = respType{}
		resp["status"] = 1
		c.JSON(http.StatusOK, resp)
	})

	r.GET("/cookies", func(c *gin.Context) {
		list := make([]string, 0)
		for _, cookie := range c.Request.Cookies() {
			list = append(list, fmt.Sprintf("%s : %s \n", cookie.Name, cookie.Value))
		}
		if len(list) == 0 {
			c.String(http.StatusOK, "No Cookies.")
			return
		}
		c.String(http.StatusOK, "<pre>"+strings.Join(list, "")+"</pre>")
	})

	r.GET("/cookie_set/:name/:cookie", func(c *gin.Context) {
		name := c.Param("name")
		cookie := c.Param("cookie")
		if len(cookie) < 1 {
			c.String(http.StatusBadRequest, "please add cookie to url /cookie_set/:name/:cookie")
			return
		}
		c.SetCookie(name, cookie, 3600, "/", c.Request.Host, true, false)
		c.String(http.StatusOK, "setCookie OK.")
	})

	r.GET("/generate_204", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// Handle all requests using net/http
	http.Handle("/", r)
}

func main() {
	appengine.Main()
}

// Resp common response struct
type Resp map[string]string
