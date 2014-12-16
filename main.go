package main

import (
	gami "code.google.com/p/gami"
	"flag"
	"fmt"
	"github.com/bmizerany/pat"
	log "github.com/cihub/seelog"
	"net"
	"net/http"
	"os"
)

const (
	VERSION = "0.0.1"
)

var (
	config  *Config
	version = flag.Bool("version", false, "show version")
)

func loadLogger() {
	pwd, errPath := os.Getwd()
	if errPath != nil {
		fmt.Println(errPath)
		os.Exit(1)
	}

	logFilePath := fmt.Sprintf("%s/%s", pwd, "conf/logger.xml")
	logger, err := log.LoggerFromConfigAsFile(logFilePath)

	if err != nil {
		log.Error("Can not load the logger configuration file, Please check if the file logger.xml exists on current directory", err)
		os.Exit(1)
	} else {
		log.ReplaceLogger(logger)
		logger.Flush()
	}

}

func cleanup() {
	log.Info("Bye buddy. Live well\n")
	log.Flush()

}

func init() {
	loadLogger()
}

func getStatusGbl(w http.ResponseWriter, r *http.Request) {
	sysInfo := GetSysInfo()
	mysqlDbStatus := 0
	mongoDbStatus := 0

	err := GetMySqlStatus(config.DbMySqlHost, config.DbMySqlUser, config.DbMySqlPassword)

	if err == nil {
		mysqlDbStatus = 1
	} else {
		log.Errorf("Failed get mysql informations : %s", err)
	}

	err = GetMongoStatus(config.MongoHost)

	if err == nil {
		mongoDbStatus = 1
	} else {
		log.Errorf("Failed get mongodb informations : %s", err)
	}

	var crmMon CrmMon
	crmMon, err = GetClusterInfo(config.CrmMonFile)

	if err != nil {
		log.Errorf("Failed get cluster informations : %s", err)
	}

	var ast *gami.Asterisk
	var con net.Conn
	var astInfo *AsteriskInfo

	ast, con, err = ConnectToAsterisk(config.AsteriskAddr, config.AsteriskPort, config.AsteriskUser, config.AsteriskPassword)

	if err != nil {
		astInfo = NewAsteriskInfo()

		response := fmt.Sprintf("%s\nMysqlDb = [%d]\nMongoDb = [%d]\n%s\n%s",
			sysInfo.GblString(), mysqlDbStatus, mongoDbStatus, crmMon.GblString(), astInfo.GblString())

		w.Write([]byte(response))

		return
	}

	astInfo, err = GetAsteriskInfo(ast)

	defer ast.Logoff()
	defer con.Close()

	response := fmt.Sprintf("%s\nMysqlDb = [%d]\nMongoDb = [%d]\n%s\n%s",
		sysInfo.GblString(), mysqlDbStatus, mongoDbStatus, crmMon.GblString(), astInfo.GblString())
	w.Write([]byte(response))
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><body><H1><center>Thank to use AMON version  " + VERSION + "(not production ready)</center></b><center><H2>Live well</H2></center></body></html>"))

}

func main() {
	flag.Parse()
	//
	if *version {
		fmt.Printf("Version : %s\n", VERSION)
		fmt.Println("Get fun! Live well !")
		return
	}
	var err error
	config, err = NewConfig()
	if err != nil {
		log.Criticalf("Error : %s.", err)
		return
	}
	mux := pat.New()
	mux.Get("/getstatus/gbl", http.HandlerFunc(getStatusGbl))

	mux.Get("/index", http.HandlerFunc(getIndex))

	mux.Get("/about", http.HandlerFunc(getIndex))

	http.Handle("/", mux)

	log.Infof("Start listening on %s. Live well", config.ListenAddr)
	if err := http.ListenAndServe(config.ListenAddr, nil); err != nil {
		log.Criticalf("Error ListenAndServe:", err)
		return
	}

}
