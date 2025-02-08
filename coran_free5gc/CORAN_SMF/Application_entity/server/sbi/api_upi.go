package sbi

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/coranlabs/CORAN_SMF/Application_entity/pkg/factory"
	smf_context "github.com/coranlabs/CORAN_SMF/Messages_handling_entity/context"
)

func (s *Server) getUPIRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodGet,
			Pattern: "/",
			APIFunc: func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "Service Available"})
			},
		},
		{
			Method:  http.MethodGet,
			Pattern: "/upNodesLinks",
			APIFunc: s.GetUpNodesLinks,
		},
		{
			Method:  http.MethodPost,
			Pattern: "/upNodesLinks",
			APIFunc: s.PostUpNodesLinks,
		},
		{
			Method:  http.MethodDelete,
			Pattern: "/upNodesLinks/:upNodeRef",
			APIFunc: s.DeleteUpNodeLink,
		},
	}
}

func (s *Server) GetUpNodesLinks(c *gin.Context) {
	upi := smf_context.GetSelf().UserPlaneInformation
	upi.Mu.RLock()
	defer upi.Mu.RUnlock()

	nodes := upi.UpNodesToConfiguration()
	links := upi.LinksToConfiguration()

	json := &factory.UserPlaneInformation{
		UPNodes: nodes,
		Links:   links,
	}

	c.JSON(http.StatusOK, json)
}

func (s *Server) PostUpNodesLinks(c *gin.Context) {
	upi := smf_context.GetSelf().UserPlaneInformation
	upi.Mu.Lock()
	defer upi.Mu.Unlock()

	var json factory.UserPlaneInformation
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	upi.UpNodesFromConfiguration(&json)
	upi.LinksFromConfiguration(&json)

	for _, upf := range upi.UPFs {
		// only associate new ones
		if upf.UPF.UPFStatus == smf_context.NotAssociated {
			upf.UPF.Ctx, upf.UPF.CancelFunc = context.WithCancel(context.Background())
			go s.Processor().ToBeAssociatedWithUPF(smf_context.GetSelf().Ctx, upf.UPF)
		}
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (s *Server) DeleteUpNodeLink(c *gin.Context) {
	// current version does not allow node deletions when ulcl is enabled
	if smf_context.GetSelf().ULCLSupport {
		c.JSON(http.StatusForbidden, gin.H{})
	} else {
		upNodeRef := c.Params.ByName("upNodeRef")
		upi := smf_context.GetSelf().UserPlaneInformation
		upi.Mu.Lock()
		defer upi.Mu.Unlock()
		if upNode, ok := upi.UPNodes[upNodeRef]; ok {
			if upNode.Type == smf_context.UPNODE_UPF {
				go s.Processor().ReleaseAllResourcesOfUPF(upNode.UPF)
			}
			upi.UpNodeDelete(upNodeRef)
			upNode.UPF.CancelFunc()
			c.JSON(http.StatusOK, gin.H{"status": "OK"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{})
		}
	}
}
