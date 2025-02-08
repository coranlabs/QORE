package sbi

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getPcfRoutesRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodGet,
			Pattern: "/policies/:polAssoId",
			// APIFunc: s.HTTPPoliciesPolAssoIdGet,
		},
		{
			Method:  http.MethodDelete,
			Pattern: "/policies/:polAssoId",
			// APIFunc: s.HTTPPoliciesPolAssoIdDelete,
		},
		{
			Method:  http.MethodPost,
			Pattern: "/policies/:polAssoId/update",
			// APIFunc: s.HTTPPoliciesPolAssoIdUpdatePost,
		},
		{
			Method:  http.MethodPost,
			Pattern: "/policies",
			APIFunc: s.HTTPPoliciesPostForwardPcf,
		},
	}
}

func (s *Server) HTTPPoliciesPostForwardPcf(c *gin.Context) {
	action := c.Param("action") // This gives you whatever comes after /nausf-auth/v1/
	log.Println("Action:", action)
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	AusfUri := "http://pcf.free5gc.org:8000" + c.Request.URL.Path

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
