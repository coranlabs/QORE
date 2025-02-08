package sbi

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getUdmRoutesRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodPost,
			Pattern: "/:supi/auth-events",
			APIFunc: s.HandleConfirmAuthForwardUdm,
		},
	}
}

func (s *Server) HandleConfirmAuthForwardUdm(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	udmuri := "http://udm.free5gc.org:8000" + c.Request.URL.Path

	// AusfUri := s.getNFUri()

	newRequest, err := http.NewRequest(c.Request.Method, udmuri, bytes.NewBuffer(bodyBytes))
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
