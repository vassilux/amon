package main

import (
	"fmt"
	"github.com/cloudfoundry/gosigar"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"labix.org/v2/mgo"
	"os"
	"time"
)

func memoryFormat(val uint64) uint64 {
	return val / 1024
}

const output_format = "%-15s %4s %4s %5s %4s %-15s\n"

func formatSize(size uint64) string {
	return sigar.FormatSize(size * 1024)
}

func GetSysInfo() SysInfo {
	sysInfo := SysInfo{}
	sysInfo.Drives = make(map[string]string)
	uptime := sigar.Uptime{}
	uptime.Get()
	mem := sigar.Mem{}
	mem.Get()
	sysInfo.Uptime = uptime.Format()
	sysInfo.MemTotal = memoryFormat(mem.Total)
	sysInfo.MemUsed = memoryFormat(mem.Used)
	sysInfo.MemFree = memoryFormat(mem.Free)
	hostname, err := os.Hostname()
	//
	sysInfo.MemUsed = sysInfo.MemUsed * 100 / sysInfo.MemTotal
	//
	if err == nil {
		sysInfo.Name = hostname
	} else {
		sysInfo.Name = "mistery"
	}
	//
	now := time.Now()
	sysInfo.SysTime = fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())
	fslist := sigar.FileSystemList{}
	fslist.Get()

	for _, fs := range fslist.List {
		dir_name := fs.DirName

		usage := sigar.FileSystemUsage{}

		usage.Get(dir_name)
		sysInfo.Drives[dir_name] = sigar.FormatPercent(usage.UsePercent())
	}

	return sysInfo
}

func GetMySqlStatus(host, user, password string) error {
	db := mysql.New("tcp", "", host, user, password, "asteriskcdrdb")
	//
	err := db.Connect()

	if err != nil {
		return err
	}
	err = db.Close()

	return err
}

func GetMongoStatus(address string) error {
	session, err := mgo.Dial(address)

	if err != nil {

		return err
	}
	session.SetMode(mgo.Monotonic, true)

	defer session.Close()

	return nil
}
