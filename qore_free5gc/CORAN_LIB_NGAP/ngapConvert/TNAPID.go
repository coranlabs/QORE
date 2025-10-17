package ngapConvert

import (
	"encoding/binary"

	"github.com/coranlabs/CORAN_LIB_NGAP/ngapType"
)

func TNAPIDToInt(id ngapType.TNAPID) (idInt64 int64) {
	idInt64 = int64(binary.BigEndian.Uint16(id.Value))
	return
}

func TNAPIDToNgap(idInt64 uint64) (id ngapType.TNAPID) {
	id.Value = make([]byte, 6)
	binary.BigEndian.PutUint16(id.Value[0:2], uint16(idInt64>>32))
	binary.BigEndian.PutUint32(id.Value[2:], uint32(idInt64))
	return
}
