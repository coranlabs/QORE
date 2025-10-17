package ebpf_datapath

import (
	"context"
	"io"
	"net"
	"os"

	"github.com/RoaringBitmap/roaring"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	UPF_config "github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/config"
	"github.com/coranlabs/CORAN_UPF_eBPF/Application_entity/logger"
	"golang.org/x/sys/unix"
)

type EBPF_entity struct {
	Coran_ebpf_datapathObjects

	FarIdTracker *IdTracker
	QerIdTracker *IdTracker
}

type IdTracker struct {
	bitmap  *roaring.Bitmap
	maxSize uint32
}

func NewIdTracker(size uint32) *IdTracker {
	newBitmap := roaring.NewBitmap()
	newBitmap.Flip(0, uint64(size))

	return &IdTracker{
		bitmap:  newBitmap,
		maxSize: size,
	}
}

func (EBPF_controller *EBPF_entity) Unload_and_detach(closers io.Closer) {
	if err := closers.Close(); err != nil {
		logger.EBPF_Datapath.Errorf("Failed to unload and detach ebpf objects: %s", err)
	} else {
		logger.EBPF_Datapath.Infof("Successfully unloaded and detached ebpf objects")
	}
}

func CloseAllObjects(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// type LoaderFunc func(obj interface{}, opts *ebpf.CollectionOptions) error
// type Loader struct {
// 	LoaderFunc
// 	object interface{}
// }

// func LoadAllObjects(opts *ebpf.CollectionOptions, loaders ...Loader) error {
// 	for _, loader := range loaders {
// 		if err := loader.LoaderFunc(loader.object, opts); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func (EBPF_controller *EBPF_entity) Load_and_attach(config *UPF_config.UpfConfig, ctx context.Context) error {
	pinPath := "/sys/fs/bpf/coran_upf_datapath"
	if err := os.MkdirAll(pinPath, os.ModePerm); err != nil {
		logger.EBPF_Datapath.Infof("failed to create bpf fs subpath: %+v", err)
		return err
	}

	collectionOptions := ebpf.CollectionOptions{
		Maps: ebpf.MapOptions{
			// Pin the map to the BPF filesystem and configure the
			// library to automatically re-write it in the BPF
			// program, so it can be re-used if it already exists or
			// create it if not
			PinPath: pinPath,
		},
	}

	err := LoadCoran_ebpf_datapathObjects(&EBPF_controller.Coran_ebpf_datapathObjects, &collectionOptions)
	if err != nil {
		return err
	}

	// err := LoadCoran_ebpf_datapathObjects(&EBPF_controller.Coran_ebpf_datapathObjects, &collectionOptions)
	// if err != nil {
	// 	return err
	// }

	// if err := EBPF_controller.ResizeAllMaps(config.QerMapSize, config.FarMapSize, config.PdrMapSize); err != nil {
	// 	logger.InitLog.Errorf("Failed to resize all maps: %s", err)
	// 	return err
	// }

	go func() {

		for _, ifaceName := range config.InterfaceName {
			iface, err := net.InterfaceByName(ifaceName)
			if err != nil {
				logger.EBPF_Datapath.Errorf("Lookup network iface %q: %s", ifaceName, err.Error())
			}

			// Attach the program.
			l, err := link.AttachXDP(link.XDPOptions{
				Program:   EBPF_controller.HexaDatapathEntrypoint,
				Interface: iface.Index,
				Flags:     flag(config.XDPAttachMode),
			})
			if err != nil {
				logger.EBPF_Datapath.Errorf("Could not attach XDP program: %s", err.Error())
			}
			defer l.Close()

			logger.EBPF_Datapath.Infof("Attached XDP program to iface %q (index %d)", iface.Name, iface.Index)
		}
		logger.EBPF_Datapath.Infof("Sleeper 2 is here")
		//time.Sleep(10000 * time.Minute)

		<-ctx.Done()
	}()

	logger.EBPF_Datapath.Infof("Sleeper 3 is here")
	//time.Sleep(10000 * time.Minute)

	return nil

}

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cflags "$BPF_CFLAGS" -target bpf Coran_ebpf_datapath 	Hexa_datapath_entity/Hexa_datapath.c -- -I. -O2 -Wall -g

func Setup_eBPF(config *UPF_config.UpfConfig, ctx context.Context) (*EBPF_entity, error) {

	err := unix.Setrlimit(unix.RLIMIT_MEMLOCK, &unix.Rlimit{
		Cur: unix.RLIM_INFINITY,
		Max: unix.RLIM_INFINITY,
	})

	if err != nil {
		logger.EBPF_Datapath.Fatal("Resource limit incrementation has failed due to ", err)
		return nil, err
	}

	EBPF_controller := &EBPF_entity{
		FarIdTracker: NewIdTracker(UPF_config.Conf.FarMapSize),
		QerIdTracker: NewIdTracker(UPF_config.Conf.QerMapSize),
	}

	if err := EBPF_controller.Load_and_attach(&UPF_config.Conf, ctx); err != nil {
		logger.EBPF_Datapath.Errorf("Loading bpf objects failed: %s", err.Error())
		return nil, err
	}

	return EBPF_controller, nil
}

// func (EBPF_controller *EBPF_entity) ResizeAllMaps(qerMapSize uint32, farMapSize uint32, pdrMapSize uint32) error {
// 	//QER
// 	if err := ResizeEbpfMap(&EBPF_controller.QER_map, EBPF_controller.HexaDatapathEntrypoint, qerMapSize); err != nil {
// 		logger.InitLog.Infof("Failed to resize QER map: %s", err)
// 		return err
// 	}

// 	//FAR
// 	if err := ResizeEbpfMap(&EBPF_controller.FAR_map, EBPF_controller.HexaDatapathEntrypoint, farMapSize); err != nil {
// 		logger.InitLog.Infof("Failed to resize FAR map: %s", err)
// 		return err
// 	}

// 	// PDR
// 	if err := ResizeEbpfMap(&EBPF_controller.PDR_downlinkMap, EBPF_controller.HexaDatapathEntrypoint, pdrMapSize); err != nil {
// 		logger.InitLog.Infof("Failed to resize PDR map: %s", err)
// 		return err
// 	}

// 	if err := ResizeEbpfMap(&EBPF_controller.PDR_uplinkMap, EBPF_controller.HexaDatapathEntrypoint, pdrMapSize); err != nil {
// 		logger.InitLog.Infof("Failed to resize PDR map: %s", err)
// 		return err
// 	}

// 	return nil
// }

// func ResizeEbpfMap(eMap **ebpf.Map, eProg *ebpf.Program, newSize uint32) error {
// 	mapInfo, err := (*eMap).Info()
// 	if err != nil {
// 		logger.InitLog.Infof("Failed get ebpf map info: %s", err)
// 		return err
// 	}
// 	mapInfo.MaxEntries = newSize
// 	// Create a new MapSpec using the information from MapInfo
// 	mapSpec := &ebpf.MapSpec{
// 		Name:       mapInfo.Name,
// 		Type:       mapInfo.Type,
// 		KeySize:    mapInfo.KeySize,
// 		ValueSize:  mapInfo.ValueSize,
// 		MaxEntries: mapInfo.MaxEntries,
// 		Flags:      mapInfo.Flags,
// 	}
// 	if err != nil {
// 		logger.InitLog.Infof("Failed to close old ebpf map: %s, %+v", err, *eMap)
// 		return err
// 	}

// 	// Unpin the old map
// 	err = (*eMap).Unpin()
// 	if err != nil {
// 		logger.InitLog.Infof("Failed to unpin old ebpf map: %s, %+v", err, *eMap)
// 		return err
// 	}

// 	// Close the old map
// 	err = (*eMap).Close()
// 	if err != nil {
// 		logger.InitLog.Infof("Failed to close old ebpf map: %s, %+v", err, *eMap)
// 		return err
// 	}

// 	// Old map will be garbage collected sometime after this point

// 	*eMap, err = ebpf.NewMapWithOptions(mapSpec, ebpf.MapOptions{})
// 	if err != nil {
// 		logger.InitLog.Infof("Failed to create resized ebpf map: %s", err)
// 		return err
// 	}
// 	err = eProg.BindMap(*eMap)
// 	if err != nil {
// 		logger.InitLog.Infof("Failed to bind resized ebpf map: %s", err)
// 		return err
// 	}
// 	return nil
// }

func flag(XDP_attach_mode string) link.XDPAttachFlags {
	switch XDP_attach_mode {
	case "generic":
		return link.XDPGenericMode
	case "native":
		return link.XDPDriverMode
	case "offload":
		return link.XDPOffloadMode
	default:
		return link.XDPGenericMode
	}
}

func (t *IdTracker) Release(id uint32) {
	if id >= t.maxSize {
		return
	}

	t.bitmap.Add(id)
}
