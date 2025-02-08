package rest

import (
	"net/http"
	"strconv"
	"unsafe"

	"github.com/coranlabs/HEXA_UPF/ebpf"
	"github.com/coranlabs/HEXA_UPF/internal/logger"
	"github.com/gin-gonic/gin"
)

type PdrElement struct {
	Id                 uint32 `json:"id"`
	OuterHeaderRemoval uint8  `json:"outer_header_removal"`
	FarId              uint32 `json:"far_id"`
	QerId              uint32 `json:"qer_id"`
}

func (h *ApiHandler) getUplinkPdrValue(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.MainLog.Tracef("Not an integer id: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var value ebpf.PdrInfo
	if err = h.BpfObjects.IpEntrypointObjects.PdrMapUplinkIp4.Lookup(uint32(id), unsafe.Pointer(&value)); err != nil {
		logger.MainLog.Tracef("Error reading map: %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, PdrElement{
		Id:                 uint32(id),
		OuterHeaderRemoval: value.OuterHeaderRemoval,
		FarId:              value.FarId,
		QerId:              value.QerId,
	})
}

func (h *ApiHandler) setUplinkPdrValue(c *gin.Context) {
	var pdrElement PdrElement
	if err := c.BindJSON(&pdrElement); err != nil {
		logger.MainLog.Tracef("Parsing request body error: %s", err.Error())
		return
	}

	var value = ebpf.PdrInfo{
		OuterHeaderRemoval: pdrElement.OuterHeaderRemoval,
		FarId:              pdrElement.FarId,
		QerId:              pdrElement.QerId,
	}

	if err := h.BpfObjects.IpEntrypointObjects.PdrMapUplinkIp4.Put(uint32(pdrElement.Id), unsafe.Pointer(&value)); err != nil {
		logger.MainLog.Tracef("Error writting map: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, pdrElement)
}
