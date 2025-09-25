// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// that can be found in the LICENSE file.

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Define a simple GET endpoint
	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
