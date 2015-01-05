package main

import (
	"bufio"
	gami "code.google.com/p/gami"
	"fmt"
	check "gopkg.in/check.v1"
	"net"
	"strings"
	"testing"
	"time"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	check.TestingT(t)
}

type UnitSuite struct{}

type IntegrationSuite struct{}

var _ = check.Suite(&IntegrationSuite{})
var ast *gami.Asterisk
var con net.Conn
var err error

var am *amock

func (s *IntegrationSuite) SetUpSuite(c *check.C) {

	l, err := net.Listen("tcp", "localhost:42420")
	if err != nil {
		c.Log("Can't start Asterisk mock: ", err)
		c.Fail()
	}

	c.Log("SetUp test suite")
	ast, con, err = ConnectToAsterisk("192.168.3.20", 5038, "astmanager", "lepanos")

	if ast == nil {
		c.Log("Can't connect to astersik: ", err)
		c.Fail()
	}

	am = &amock{l}

	sleep(2)
}

func (i *IntegrationSuite) TestGetAsteriskInfoDb(c *check.C) {
	astInfo, err := GetAsteriskInfo(ast)
	c.Logf("AsterisInfo : \n%s\n", astInfo)

	if err != nil {
		c.Log(err)
		c.Fail()
	}

	if len(astInfo.Uptime) == 0 {
		c.Fail()
	}
}

func (s *IntegrationSuite) TearDownSuite(c *check.C) {
	c.Log("TearDown")
	ast.Logoff()
	con.Close()
	am.stop()
}

type amock struct {
	ln net.Listener
}

func (a *amock) stop() {
	a.ln.Close()
}
