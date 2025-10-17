package sbi

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/coranlabs/CORAN_SCP/Application_entity/logger"
	"github.com/gin-gonic/gin"
)

func (s *Server) forwardRequestAusf(c *gin.Context, destination string) {
	// Read the request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// Restore the request body so it can be used again
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	destinationURI := destination + c.Request.URL.Path

	// Create a new request with the original method, URL, and body
	newRequest, err := http.NewRequest(c.Request.Method, destinationURI, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new request"})
		return
	}

	// Copy headers from the original request
	newRequest.Header = c.Request.Header

	client := &http.Client{}

	// Send the request and wait for the response
	resp, err := client.Do(newRequest)
	if err != nil {
		log.Printf("Error forwarding request to new destination: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Set response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// Send the response back to the client
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)

	// Log the successful forwarding
	logger.ScpLog.Tracef("Request forwarded to: %s, Status: %d", destinationURI, resp.StatusCode)
	logger.ScpLog.Infof("Request forwarded to: %s, Status: %d", destination, resp.StatusCode)
}

func (s *Server) forwardRequestPcf(c *gin.Context, destination string) {
	// Read the request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// Restore the request body so it can be used again
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	destinationURI := destination + c.Request.URL.Path

	// Create a new request with the original method, URL, and body
	newRequest, err := http.NewRequest(c.Request.Method, destinationURI, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new request"})
		return
	}

	// Copy headers from the original request
	newRequest.Header = c.Request.Header

	client := &http.Client{}

	// Send the request and wait for the response
	resp, err := client.Do(newRequest)
	if err != nil {
		log.Printf("Error forwarding request to new destination: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Set response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// Send the response back to the client
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)

	// Log the successful forwarding
	logger.ScpLog.Tracef("Request forwarded to: %s, Status: %d", destinationURI, resp.StatusCode)
	logger.ScpLog.Infof("Request forwarded to: %s, Status: %d", destination, resp.StatusCode)
}

func (s *Server) forwardRequestAmf(c *gin.Context, destination string) {
	// Read the request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// Restore the request body so it can be used again
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	destinationURI := destination + c.Request.URL.Path

	// Create a new request with the original method, URL, and body
	newRequest, err := http.NewRequest(c.Request.Method, destinationURI, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new request"})
		return
	}

	// Copy headers from the original request
	newRequest.Header = c.Request.Header

	client := &http.Client{}

	// Send the request and wait for the response
	resp, err := client.Do(newRequest)
	if err != nil {
		log.Printf("Error forwarding request to new destination: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Set response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// Send the response back to the client
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)

	// Log the successful forwarding
	logger.ScpLog.Tracef("Request forwarded to: %s, Status: %d", destinationURI, resp.StatusCode)
	logger.ScpLog.Infof("Request forwarded to: %s, Status: %d", destination, resp.StatusCode)
}

func (s *Server) forwardRequestNssf(c *gin.Context, destination string) {
	// Read the request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// Restore the request body so it can be used again
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	destinationURI := destination + c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		destinationURI += "?" + c.Request.URL.RawQuery
	}

	// Create a new request with the original method, URL, and body
	newRequest, err := http.NewRequest(c.Request.Method, destinationURI, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new request"})
		return
	}

	// Copy headers from the original request
	newRequest.Header = c.Request.Header

	client := &http.Client{}

	// Send the request and wait for the response
	resp, err := client.Do(newRequest)
	if err != nil {
		log.Printf("Error forwarding request to new destination: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Set response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}
	// Send the response back to the client
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)

	// Log the successful forwarding
	logger.ScpLog.Tracef("Request forwarded to: %s, Status: %d", destinationURI, resp.StatusCode)
	logger.ScpLog.Infof("Request forwarded to: %s, Status: %d", destination, resp.StatusCode)
}

func (s *Server) forwardRequestSmf(c *gin.Context, destination string) {
	// Read the request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// Restore the request body so it can be used again
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	destinationURI := destination + c.Request.URL.Path

	// Create a new request with the original method, URL, and body
	newRequest, err := http.NewRequest(c.Request.Method, destinationURI, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new request"})
		return
	}

	// Copy headers from the original request
	newRequest.Header = c.Request.Header

	client := &http.Client{}

	// Send the request and wait for the response
	resp, err := client.Do(newRequest)
	if err != nil {
		log.Printf("Error forwarding request to new destination: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Set response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// Send the response back to the client
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)

	// Log the successful forwarding
	logger.ScpLog.Tracef("Request forwarded to: %s, Status: %d", destinationURI, resp.StatusCode)
	logger.ScpLog.Infof("Request forwarded to: %s, Status: %d", destination, resp.StatusCode)
}

func (s *Server) forwardRequestUdm(c *gin.Context, destination string) {
	// Read the request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// Restore the request body so it can be used again
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	destinationURI := destination + c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		destinationURI += "?" + c.Request.URL.RawQuery
	}

	// Create a new request with the original method, URL, and body
	newRequest, err := http.NewRequest(c.Request.Method, destinationURI, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new request"})
		return
	}

	// Copy headers from the original request
	newRequest.Header = c.Request.Header

	client := &http.Client{}

	// Send the request and wait for the response
	resp, err := client.Do(newRequest)
	if err != nil {
		log.Printf("Error forwarding request to new destination: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Set response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// Send the response back to the client
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)

	// Log the successful forwarding
	logger.ScpLog.Tracef("Request forwarded to: %s, Status: %d", destinationURI, resp.StatusCode)
	logger.ScpLog.Infof("Request forwarded to: %s, Status: %d", destination, resp.StatusCode)
}

func (s *Server) forwardRequestUdr(c *gin.Context, destination string) {
	// Read the request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// Restore the request body so it can be used again
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	destinationURI := destination + c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		destinationURI += "?" + c.Request.URL.RawQuery
	}

	// Create a new request with the original method, URL, and body
	newRequest, err := http.NewRequest(c.Request.Method, destinationURI, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new request"})
		return
	}

	// Copy headers from the original request
	newRequest.Header = c.Request.Header

	client := &http.Client{}

	// Send the request and wait for the response
	resp, err := client.Do(newRequest)
	if err != nil {
		log.Printf("Error forwarding request to new destination: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Set response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// Send the response back to the client
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)

	// Log the successful forwarding
	logger.ScpLog.Tracef("Request forwarded to: %s, Status: %d", destinationURI, resp.StatusCode)
	logger.ScpLog.Infof("Request forwarded to: %s, Status: %d", destination, resp.StatusCode)
}

func (s *Server) forwardRequestChf(c *gin.Context, destination string) {
	// Read the request body
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// Restore the request body so it can be used again
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	destinationURI := destination + c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		destinationURI += "?" + c.Request.URL.RawQuery
	}

	// Create a new request with the original method, URL, and body
	newRequest, err := http.NewRequest(c.Request.Method, destinationURI, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new request"})
		return
	}

	// Copy headers from the original request
	newRequest.Header = c.Request.Header

	client := &http.Client{}

	// Send the request and wait for the response
	resp, err := client.Do(newRequest)
	if err != nil {
		log.Printf("Error forwarding request to new destination: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Set response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}

	// Send the response back to the client
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)

	// Log the successful forwarding
	logger.ScpLog.Tracef("Request forwarded to: %s, Status: %d", destinationURI, resp.StatusCode)
	logger.ScpLog.Infof("Request forwarded to: %s, Status: %d", destination, resp.StatusCode)
}
