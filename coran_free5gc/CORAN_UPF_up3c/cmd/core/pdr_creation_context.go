package core

import (
	"errors"
	"fmt"
	"net"

	"github.com/coranlabs/CORAN_GO_PFCP/ie"
	"github.com/coranlabs/HEXA_UPF/cmd/core/service"
	"github.com/coranlabs/HEXA_UPF/ebpf"
	"github.com/coranlabs/HEXA_UPF/internal/logger"
	config "github.com/coranlabs/HEXA_UPF/src/config"
	"github.com/rs/zerolog/log"
)

type PDRCreationContext struct {
	Session         *Session
	ResourceManager *service.ResourceManager
	TEIDCache       map[uint8]uint32
}

func NewPDRCreationContext(session *Session, resourceManager *service.ResourceManager) *PDRCreationContext {
	return &PDRCreationContext{
		Session:         session,
		ResourceManager: resourceManager,
		TEIDCache:       make(map[uint8]uint32),
	}
}

// extractPDR takes a PDR IE and a SPDRInfo object, and extracts the relevant info from the PDR IE
// and stores it in the SPDRInfo object. It also allocates a TEID if necessary and stores it in the
// SPDRInfo object.
func (pdrContext *PDRCreationContext) extractPDR(pdr *ie.IE, spdrInfo *SPDRInfo) error {
	logger.PDRLog.Tracef("in extractPDR")

	/*************  ✨ Codeium Command 🌟  *************/
	if outerHeaderRemoval, err := pdr.OuterHeaderRemovalDescription(); err == nil {
		logger.PDRLog.Tracef("outerHeaderRemoval: %d", outerHeaderRemoval)

		spdrInfo.PdrInfo.OuterHeaderRemoval = outerHeaderRemoval
	}
	if farid, err := pdr.FARID(); err == nil {
		logger.PDRLog.Tracef("farid: %d", farid)
		spdrInfo.PdrInfo.FarId = farid
		//spdrInfo.PdrInfo.FarId = pdrContext.getFARID(farid)

	}
	if qerid, err := pdr.QERID(); err == nil {
		logger.QERLog.Tracef("qerid: %d", qerid)
		spdrInfo.PdrInfo.QerId = pdrContext.getQERID(qerid)
	}
	/******  355cf5ca-37b9-4b21-b7e9-c5e5ba413c9a  *******/

	pdi, err := pdr.PDI()
	if err != nil {
		log.Error().Msgf("PDI err: %v", err)
		return err
	}

	if sdfFilter, err := pdr.SDFFilter(); err == nil {
		if sdfFilter.FlowDescription == "" {
			logger.Pfcplog.Warnf("SDFFilter is empty")
		} else if sdfFilterParsed, err := ParseSdfFilter(sdfFilter.FlowDescription); err == nil {
			spdrInfo.PdrInfo.SdfFilter = &sdfFilterParsed
		} else {
			log.Error().Msgf("SDFFilter err: %v", err)
			return err
		}
	}

	if teidPdiId := findIEindex(pdi, 21); teidPdiId != -1 { // IE Type F-TEID
		if fteid, err := pdi[teidPdiId].FTEID(); err == nil {
			var teid = fteid.TEID
			if fteid.HasCh() {
				var allocate = true
				if fteid.HasChID() {
					if teidFromCache, ok := pdrContext.hasTEIDCache(fteid.ChooseID); ok {
						allocate = false
						teid = teidFromCache
						spdrInfo.Allocated = true
					}
				}
				if allocate {
					allocatedTeid, err := pdrContext.getFTEID(pdrContext.Session.RemoteSEID, spdrInfo.PdrID)
					if err != nil {
						log.Error().Msgf("AllocateTEID err: %v", err)
						return fmt.Errorf("can't allocate TEID: %s", causeToString(ie.CauseNoResourcesAvailable))
					}
					teid = allocatedTeid
					spdrInfo.Allocated = true
					if fteid.HasChID() {
						pdrContext.setTEIDCache(fteid.ChooseID, teid)
					}
				}
			}
			spdrInfo.Teid = teid
			return nil
		}
		return fmt.Errorf("F-TEID IE is missing")
	} else if ueIP, err := pdr.UEIPAddress(); err == nil {
		if config.Conf.FeatureUEIP && hasCHV4(ueIP.Flags) {
			if ip, err := pdrContext.getIP(); err == nil {
				ueIP.IPv4Address = cloneIP(ip)
				spdrInfo.Allocated = true
			} else {
				log.Error().Msg(err.Error())
			}
		}
		if ueIP.IPv4Address != nil {
			spdrInfo.Ipv4 = cloneIP(ueIP.IPv4Address)
		} else if ueIP.IPv6Address != nil {
			spdrInfo.Ipv6 = cloneIP(ueIP.IPv6Address)
		} else {
			return fmt.Errorf("UE IP Address IE is missing")
		}

		return nil
	} else {
		logger.Pfcplog.Infof("Both F-TEID IE and UE IP Address IE are missing")
		return nil
	}
}

func (pdrContext *PDRCreationContext) deletePDR(spdrInfo SPDRInfo, mapOperations ebpf.ForwardingPlaneController) error {
	if spdrInfo.Ipv4 != nil {
		if err := mapOperations.DeletePdrDownlink(spdrInfo.Ipv4); err != nil {
			return fmt.Errorf("Can't delete IPv4 PDR: %s", err.Error())
		}
	} else if spdrInfo.Ipv6 != nil {
		if err := mapOperations.DeleteDownlinkPdrIp6(spdrInfo.Ipv6); err != nil {
			return fmt.Errorf("Can't delete IPv6 PDR: %s", err.Error())
		}
	} else {
		if _, ok := pdrContext.TEIDCache[uint8(spdrInfo.Teid)]; !ok {
			if err := mapOperations.DeletePdrUplink(spdrInfo.Teid); err != nil {
				return fmt.Errorf("Can't delete GTP PDR: %s", err.Error())
			}
			pdrContext.TEIDCache[uint8(spdrInfo.Teid)] = 0
		}
	}
	if spdrInfo.Teid != 0 {
		pdrContext.ResourceManager.FTEIDM.ReleaseTEID(pdrContext.Session.RemoteSEID)
	}
	return nil
}

func (pdrContext *PDRCreationContext) getFARID(farid uint32) uint32 {
	return pdrContext.Session.GetFar(farid).GlobalId
}

func (pdrContext *PDRCreationContext) getQERID(qerid uint32) uint32 {
	return pdrContext.Session.GetQer(qerid).GlobalId
}

func (pdrContext *PDRCreationContext) getFTEID(seID uint64, pdrID uint32) (uint32, error) {
	if pdrContext.ResourceManager == nil || pdrContext.ResourceManager.FTEIDM == nil {
		return 0, errors.New("FTEID manager is nil")
	}

	allocatedTeid, err := pdrContext.ResourceManager.FTEIDM.AllocateTEID(seID, pdrID)
	if err != nil {
		log.Error().Msgf("AllocateTEID err: %v", err)
		return 0, fmt.Errorf("Can't allocate TEID: %s", causeToString(ie.CauseNoResourcesAvailable))
	}
	return allocatedTeid, nil
}

func (pdrContext PDRCreationContext) getIP() (net.IP, error) {
	if pdrContext.ResourceManager == nil || pdrContext.ResourceManager.IPAM == nil {
		return nil, errors.New("IP address manager is nil")
	}
	allocatedIP, err := pdrContext.ResourceManager.IPAM.AllocateIP(pdrContext.Session.RemoteSEID)
	if err != nil {
		return nil, fmt.Errorf("can't allocate IP: %s", causeToString(ie.CauseNoResourcesAvailable))
	}
	return allocatedIP, nil
}

func (pdrContext *PDRCreationContext) hasTEIDCache(chooseID uint8) (uint32, bool) {
	teid, ok := pdrContext.TEIDCache[chooseID]
	return teid, ok
}

func (pdrContext *PDRCreationContext) setTEIDCache(chooseID uint8, teid uint32) {
	pdrContext.TEIDCache[chooseID] = teid
}

func hasCHV4(flags uint8) bool {
	return flags&(1<<4) != 0
}
