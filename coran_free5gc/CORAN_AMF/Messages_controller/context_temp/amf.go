package context

import (
	"fmt"
	"net"
	"reflect"
	"sync"

	"github.com/coranlabs/CORAN_LIB_OPENAPI/models"
	//"github.com/coranlabs/CORAN_LIB_UTIL/idgenerator"
)

var (
	amfContext AMFContext
)

const (
	MaxNumOfTAI                       int   = 16
	MaxNumOfBroadcastPLMNs            int   = 12
	MaxNumOfPLMNs                     int   = 12
	MaxNumOfSlice                     int   = 1024
	MaxNumOfAllowedSnssais            int   = 8
	MaxValueOfAmfUeNgapId             int64 = 1099511627775
	MaxNumOfServedGuamiList           int   = 256
	MaxNumOfPDUSessions               int   = 256
	MaxNumOfDRBs                      int   = 32
	MaxNumOfAOI                       int   = 64
	MaxT3513RetryTimes                int   = 4
	MaxT3522RetryTimes                int   = 4
	MaxT3550RetryTimes                int   = 4
	MaxT3560RetryTimes                int   = 4
	MaxT3565RetryTimes                int   = 4
	MAxNumOfAlgorithm                 int   = 8
	DefaultT3502                      int   = 720  // 12 min
	DefaultT3512                      int   = 3240 // 54 min
	DefaultNon3gppDeregistrationTimer int   = 3240 // 54 min
)

type PlmnSupportItem struct {
	PlmnId     *models.PlmnId  `yaml:"plmnId" valid:"required"`
	SNssaiList []models.Snssai `yaml:"snssaiList,omitempty" valid:"required"`
}
type Ladn struct {
	Dnn     string       `yaml:"dnn" valid:"type(string),minstringlength(1),required"`
	TaiList []models.Tai `yaml:"taiList" valid:"required"`
}
type SecurityAlgorithm struct {
	IntegrityOrder []uint8 // slice of security.AlgIntegrityXXX
	CipheringOrder []uint8 // slice of security.AlgCipheringXXX
}
type NetworkName struct {
	Full  string `yaml:"full" valid:"type(string)"`
	Short string `yaml:"short,omitempty" valid:"type(string)"`
}
type IDGenerator struct {
	lock       sync.Mutex
	minValue   int64
	maxValue   int64
	valueRange int64
	offset     int64
	usedMap    map[int64]bool
}
type AMFContext struct {
	EventSubscriptionIDGenerator *IDGenerator
	EventSubscriptions           sync.Map
	UePool                       sync.Map        // map[supi]*AmfUe
	RanUePool                    sync.Map        // map[AmfUeNgapID]*RanUe
	AmfRanPool                   sync.Map        // map[net.Conn]*AmfRan
	LadnPool                     map[string]Ladn // dnn as key
	SupportTaiLists              []models.Tai
	ServedGuamiList              []models.Guami
	PlmnSupportList              []PlmnSupportItem
	RelativeCapacity             int64
	NfId                         string
	Name                         string
	NfService                    map[models.ServiceName]models.NfService // nfservice that amf support
	UriScheme                    models.UriScheme
	BindingIPv4                  string
	SBIPort                      int
	RegisterIPv4                 string
	HttpIPv6Address              string
	TNLWeightFactor              int64
	SupportDnnLists              []string
	AMFStatusSubscriptions       sync.Map // map[subscriptionID]models.SubscriptionData
	NrfUri                       string
	NrfCertPem                   string
	SecurityAlgorithm            SecurityAlgorithm
	NetworkName                  NetworkName
	NgapIpList                   []string // NGAP Server IP
	NgapPort                     int
	T3502Value                   int    // unit is second
	T3512Value                   int    // unit is second
	Non3gppDeregTimerValue       int    // unit is second
	TimeZone                     string // "[+-]HH:MM[+][1-2]", Refer to TS 29.571 - 5.2.2 Simple Data Types

	Locality string

	OAuth2Required bool
}

func GetSelf() *AMFContext {
	return &amfContext
}
func (context *AMFContext) AmfRanFindByConn(conn net.Conn) (*Ran, bool) {
	if value, ok := context.AmfRanPool.Load(conn); ok {
		return value.(*Ran), ok
	}
	return nil, false
}
func (context *AMFContext) NewAmfRan(conn net.Conn) *Ran {
	ran := Ran{}
	ran.SupportedTAList = make([]SupportedTAI, 0, MaxNumOfTAI*MaxNumOfBroadcastPLMNs)
	ran.Conn = conn
	addr := conn.RemoteAddr()
	if addr != nil {
		fmt.Printf("if")

	} else {
		fmt.Printf("else")
	}

	context.AmfRanPool.Store(conn, &ran)
	return &ran
}

type SupportedTAI struct {
	Tai        models.Tai
	SNssaiList []models.Snssai
}

func NewSupportedTAI() (tai SupportedTAI) {
	tai.SNssaiList = make([]models.Snssai, 0, MaxNumOfSlice)
	return
}

func InTaiList(servedTai models.Tai, taiList []models.Tai) bool {
	for _, tai := range taiList {
		if reflect.DeepEqual(tai, servedTai) {
			return true
		}
	}
	return false
}
