package nasConvert

import (
	"encoding/hex"

	"github.com/coranlabs/CORAN_LIB_NAS/logger"
	"github.com/coranlabs/CORAN_LIB_NAS/nasType"
	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
)

// TS24.501 9.11.2.8 S-NSSAI
func SnssaiToModels(nasSnssai *nasType.SNSSAI) (snssai models.Snssai) {
	if nasSnssai.GetLen() == uint8(4) {
		sD := nasSnssai.GetSD()
		snssai.Sd = hex.EncodeToString(sD[:])
	}
	snssai.Sst = int32(nasSnssai.GetSST())
	return
}

func SnssaiToNas(snssai models.Snssai) []uint8 {
	var buf []uint8

	if snssai.Sd == "" {
		buf = append(buf, 0x01)
		buf = append(buf, uint8(snssai.Sst))
	} else {
		buf = append(buf, 0x04)
		buf = append(buf, uint8(snssai.Sst))
		if byteArray, err := hex.DecodeString(snssai.Sd); err != nil {
			logger.ConvertLog.Warnf("Decode snssai.sd failed: %+v", err)
		} else {
			buf = append(buf, byteArray...)
		}
	}
	return buf
}

func RejectedSnssaiToNas(snssai models.Snssai, rejectCause uint8) []uint8 {
	var rejectedSnssai []uint8

	if snssai.Sd == "" {
		rejectedSnssai = append(rejectedSnssai, (0x01<<4)+rejectCause)
		rejectedSnssai = append(rejectedSnssai, uint8(snssai.Sst))
	} else {
		rejectedSnssai = append(rejectedSnssai, (0x04<<4)+rejectCause)
		rejectedSnssai = append(rejectedSnssai, uint8(snssai.Sst))
		if sDBytes, err := hex.DecodeString(snssai.Sd); err != nil {
			logger.ConvertLog.Warnf("Decode snssai.sd failed: %+v", err)
		} else {
			rejectedSnssai = append(rejectedSnssai, sDBytes...)
		}
	}

	return rejectedSnssai
}
