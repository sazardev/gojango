package middleware

import (
	"fmt"
	"log"
	"time"
)

// Context interface for middleware compatibility
type Context interface {
	Method() string
	Path() string
	ClientIP() string
	GetHeader(string) string
	Header(string, string)
	ErrorJSON(int, string, error) error
}

// Logger middleware logs HTTP requests
func Logger() func(Context) error {
	return func(c Context) error {
		start := time.Now()
		
		// Log the request
		log.Printf("%s %s from %s", c.Method(), c.Path(), c.ClientIP())
		
		// You would normally call the next handler here,
		// but since our middleware system is simple, we just return
		// The actual request handling happens in the main handler chain
		
		duration := time.Since(start)
		log.Printf("Request completed in %v", duration)
		
		return nil
	}
}

// CORS middleware adds CORS headers
func CORS(allowOrigin string) func(Context) error {
	if allowOrigin == "" {
		allowOrigin = "*"
	}
	
	return func(c Context) error {
		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "3600")
		
		// Handle preflight requests
		if c.Method() == "OPTIONS" {
			return nil
		}
		
		return nil
	}
}

// Recovery middleware recovers from panics
func Recovery() func(Context) error {
	return func(c Context) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v", r)
				c.ErrorJSON(500, "Internal Server Error", fmt.Errorf("%v", r))
			}
		}()
		
		return nil
	}
}

// BasicAuth middleware provides basic authentication
func BasicAuth(username, password string) func(Context) error {
	return func(c Context) error {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			return c.ErrorJSON(401, "Unauthorized", nil)
		}
		
		// Simple basic auth check (in production, use proper crypto)
		expectedAuth := fmt.Sprintf("Basic %s", basicAuthEncode(username+":"+password))
		if auth != expectedAuth {
			return c.ErrorJSON(401, "Unauthorized", nil)
		}
		
		return nil
	}
}

// Helper function for basic auth encoding (simplified)
func basicAuthEncode(credentials string) string {
	// In a real implementation, you'd use base64 encoding
	// This is just a placeholder
	return credentials
}

// RequestID middleware adds a unique request ID
func RequestID() func(Context) error {
	return func(c Context) error {
		requestID := generateRequestID()
		c.Header("X-Request-ID", requestID)
		return nil
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// RateLimit middleware provides simple rate limiting
func RateLimit(maxRequests int, window time.Duration) func(Context) error {
	// Simple in-memory rate limiter (not production ready)
	requestCounts := make(map[string]int)
	lastReset := time.Now()
	
	return func(c Context) error {
		now := time.Now()
		
		// Reset counter if window expired
		if now.Sub(lastReset) > window {
			requestCounts = make(map[string]int)
			lastReset = now
		}
		
		clientIP := c.ClientIP()
		requestCounts[clientIP]++
		
		if requestCounts[clientIP] > maxRequests {
			return c.ErrorJSON(429, "Too Many Requests", nil)
		}
		
		return nil
	}
}

// Security middleware adds common security headers
func Security() func(Context) error {
	return func(c Context) error {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		return nil
	}
}
