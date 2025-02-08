package sbi

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const NrfUri = "http://nrf.free5gc.org:8000"

func (s *Server) getAusfRoutesRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, "coranlabs ausf routes")
			},
		},
		{
			Method:  http.MethodPost,
			Pattern: "/ue-authentications/:authCtxId/eap-session",
			APIFunc: s.EapAuthMethodPostForwardAusf,
		},
		{
			Method:  http.MethodPost,
			Pattern: "/ue-authentications",
			APIFunc: s.UeAuthenticationsPostForwardAusf,
		},
		{
			Method:  http.MethodPut,
			Pattern: "/ue-authentications/:authCtxId/5g-aka-confirmation",
			APIFunc: s.UeAuthenticationsAuthCtxID5gAkaConfirmationPutForwardAusf,
		},
	}
}

func (s *Server) UeAuthenticationsPostForwardAusf(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	AusfUri := "http://ausf.free5gc.org:8000" + c.Request.URL.Path

	// AusfUri := s.getNFUri()

	newRequest, err := http.NewRequest(c.Request.Method, AusfUri, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}

	newRequest.Header = c.Request.Header

	client := &http.Client{}

	resp, err := client.Do(newRequest)
	if err != nil {
		log.Printf("Error forwarding request to new destination: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Set(key, value)
		}
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

func (s *Server) UeAuthenticationsAuthCtxID5gAkaConfirmationPutForwardAusf(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	AusfUri := "http://ausf.free5gc.org:8000" + c.Request.URL.Path

	newRequest, err := http.NewRequest(c.Request.Method, AusfUri, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}

	newRequest.Header = c.Request.Header

	client := &http.Client{}

	resp, err := client.Do(newRequest)
	if err != nil {
		log.Printf("Error forwarding request to new destination %s: %v", AusfUri, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body from destination %s: %v", AusfUri, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Set(key, value)
		}
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

func (s *Server) EapAuthMethodPostForwardAusf(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	AusfUri := "http://ausf.free5gc.org:8000" + c.Request.URL.Path

	// AusfUri := s.getNFUri()

	newRequest, err := http.NewRequest(c.Request.Method, AusfUri, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Failed to create new request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}

	newRequest.Header = c.Request.Header

	client := &http.Client{}

	resp, err := client.Do(newRequest)
	if err != nil {
		log.Printf("Error forwarding request to new destination: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Set(key, value)
		}
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}
