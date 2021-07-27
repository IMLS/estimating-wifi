package tlp

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
)

func cleanupTempFiles() {
	cfg := state.GetConfig()

	f1, err := filepath.Glob(filepath.Join(cfg.Paths.WWW.Root, "*.sqlite*"))
	if err != nil {
		panic(err)
	}
	f2, err := filepath.Glob(filepath.Join(cfg.Paths.WWW.Images, "*.png"))
	if err != nil {
		panic(err)
	}

	for _, f := range append(f1, f2...) {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}
func setup(theTimeIsNow string) {
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(filename)
	path := filepath.Dir(filename)
	configPath := filepath.Join(path, "..", "test", "config.yaml")
	state.UnsafeNewConfigFromPath(configPath)
	cfg := state.GetConfig()
	cfg.RunMode = "test"
	cfg.StorageMode = "sqlite"
	cfg.Databases.ManufacturersPath = filepath.Join(path, "..", "test", "manufacturers.sqlite")
	cfg.Databases.QueuesPath = filepath.Join(path, "..", "test", "www", "queues.sqlite")
	cfg.Databases.DurationsPath = filepath.Join(path, "..", "test", "www", "durations.sqlite")
	cfg.Paths.WWW.Root = filepath.Join(path, "..", "test", "www")
	cfg.Paths.WWW.Images = filepath.Join(path, "..", "test", "www", "images")
	cfg.Device.FCFSId = "ME0000-001"
	cfg.Device.DeviceTag = "testing"

	state.FlushCache()
	log.Println("Calling init config in setup")
	state.InitConfig()
	log.Println("Trying to get session id")
	log.Println("session id is ", cfg.GetCurrentSessionID())
	cfg.Logging.LogLevel = "DEBUG"
	cfg.Log().SetLogLevel(cfg.Logging.LogLevel)

	os.Mkdir(cfg.Paths.WWW.Root, 0755)
	os.Mkdir(cfg.Paths.WWW.Images, 0755)
	mock := clock.NewMock()
	mt, _ := time.Parse("2006-01-02T15:04", theTimeIsNow)
	mock.Set(mt)
	cfg.Clock = mock

	if cfg.Clock == nil {
		cfg.Log().Fatal("clock should not be nil")
	}
	cfg.Log().Debug("mock is now ", cfg.Clock.Now())
}

// type SharkFn func(string) []string
// type MonitorFn func(*models.Device)
// type SearchFn func() *models.Device

func fakeSetMonitorFn(d *models.Device) {

}

func fakeSearchFn() (d *models.Device) {
	d = &models.Device{Exists: true, Logicalname: "fakewan0"}
	return d
}

func fakeSharkFn(dev string) []string {
	return []string{"DE:AD:BE:EF:00:00", "BE:EF:00:00:00:00"}
}
func TestSimpleShark(t *testing.T) {

	setup("1975-10-11T18:00")
	cleanupTempFiles()

	var wg sync.WaitGroup
	wg.Add(1)

	cfg := state.GetConfig()
	pch := make(chan Ping)
	kb := NewKillBroker()
	go kb.Start()
	go func(ch chan Ping) {
		cfg.Log().Debug("pinging")
		pch <- Ping{}
		mock := clock.NewMock()
		mt, _ := time.Parse("2006-01-02T15:04", "1975-10-11T19:00")
		mock.Set(mt)
		cfg.Clock = mock
		pch <- Ping{}
		time.Sleep(2 * time.Second)
		kb.Publish("done")
		wg.Done()
	}(pch)

	// Create a DB for simpleshark to write to.
	db := state.NewSqliteDB(cfg.GetDurationsDatabase().GetPath())
	db.CreateTableFromStruct(structs.Duration{})

	go SimpleShark(kb, pch, db, fakeSetMonitorFn, fakeSearchFn, fakeSharkFn)
	wg.Wait()
}

// 0	1975-10-11T18:00:00Z	1627402397	ME0000-001	testing	0	1975-10-11T18:00:00Z	0	CESTNEPASUNESERIE
