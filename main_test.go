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
	//go am.start()

	/*conn, err := net.Dial("tcp", "localhost:42420")
	if err != nil {
		am.stop()
		c.Log("Can't connect to mock: ", err)
		c.Fail()
	}
	a = gami.NewAsterisk(&conn, nil)
	a.Login("admin", "admin")*/
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

func (a *amock) start() {
	for {
		//fmt.Println("Starting mock ")

		dh := func(m gami.Message) {
			fmt.Println(m)
		}
		ast.DefaultHandler(&dh)

		ch := make(chan gami.Message)
		check := func(m gami.Message) {
			fmt.Printf("DbGetResponse %s", m)
			ch <- m
		}

		ast.RegisterHandler("DbGetResponse", &check)
		ast.DbPut("test", "newkey", "1000", nil)
		ast.DbGet("test", "newkey", nil)
		fmt.Println("Before r = <-ch 1")
		r := <-ch
		fmt.Println("After r = <-ch 1")
		if r["Val"] != "1000" {
			fmt.Printf("c.Fail()")
		}
		ast.UnregisterHandler("DbGetResponse")
		ast.DbDelTree("test", "", nil)
		ast.DbGet("test", "newkey", &check)
		fmt.Println("Before r = <-ch")
		r = <-ch
		fmt.Println("After r = <-ch ")
		if r["Response"] != "Error" {
			fmt.Printf("c.Fail()")
		}
		fmt.Println("OutTestDb ")

		conn, err := a.ln.Accept()
		if err != nil {
			break
		}
		handleConnection(bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)))
	}
}

func (a *amock) stop() {
	a.ln.Close()
}

func handleConnection(rw *bufio.ReadWriter) {

	m := gami.Message{}
	db := make(map[string]map[string]string)
	confs := make(map[string]int)

	for {
		b, _, err := rw.ReadLine()
		if err != nil {
			break
		}
		if len(b) == 0 {
			r := gami.Message{
				"ActionID": m["ActionID"],
			}
			switch m["Action"] {
			case "Login":
				r["Response"] = "Success"
			case "Originate":
				if strings.HasPrefix(m["Channel"], "fakeconference/") {
					arr := strings.Split(m["Channel"], "/")
					confs[arr[1]]++
					r["Response"] = "Success"
				} else {
					r["Response"] = "Error"
				}
			case "Redirect":
				r["Response"] = "Error"
			case "DBPut":
				r["Response"] = "Success"
				fam := m["Family"]
				if _, ok := db[fam]; !ok {
					db[fam] = make(map[string]string)
				}
				db[fam][m["Key"]] = m["Value"]
			case "DBDelTree":
				r["Response"] = "Success"
				delete(db, m["Family"])
			case "DBGet":
				ok := false
				var val string
				if val, ok = db[m["Family"]][m["Key"]]; ok {
					r["Response"] = "Success"
					writeMessage(rw, r)
					r = gami.Message{
						"Val":    val,
						"Event":  "DbGetResponse",
						"Family": m["Family"],
						"Key":    m["Key"],
					}
				} else {
					r["Response"] = "Error"
				}
			case "ConfbridgeList":
				if _, ok := confs[m["Conference"]]; ok {
					for i := 0; i < confs[m["Conference"]]; i++ {
						mem := gami.Message{
							"Member":   fmt.Sprint(i),
							"Channel":  "sip/" + fmt.Sprint(i),
							"ActionID": m["ActionID"],
						}
						writeMessage(rw, mem)
					}
					r["Event"] = "ConfbridgeListComplete"
					r["EventList"] = "Complete"
				} else {
					r["Response"] = "Error"
				}
			case "ConfbridgeKick":
				if _, ok := confs[m["Conference"]]; ok {
					delete(confs, m["Conference"])
					r["Response"] = "Success"
				} else {
					r["Response"] = "Error"
				}
			case "UserEvent":
				r = m
				r["Event"] = "UserEvent"
			default:
				continue
			}
			writeMessage(rw, r)
			m = gami.Message{}
		} else {
			arr := strings.Split(string(b), ":")
			if len(arr) < 2 {
				m["Unknown"] += arr[0]
			} else {
				m[arr[0]] = arr[1]
			}
		}
	}
}

func writeMessage(w *bufio.ReadWriter, m gami.Message) {
	raw := ""
	for k, v := range m {
		raw += k + ":" + v + "\r\n"
	}
	raw += "\r\n"
	w.WriteString(raw)
	w.Flush()
}

func sleep(s int64) {
	time.Sleep(time.Duration(s) * time.Second)
}
