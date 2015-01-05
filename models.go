package main

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type SysInfo struct {
	Uptime   string
	MemTotal uint64
	MemUsed  uint64
	MemFree  uint64
	Name     string
	SysTime  string
	Drives   map[string]string
}

func (s *SysInfo) String() string {
	var val string
	val = fmt.Sprintf("sysUptime:%s\nmemTotal:%10d\nmemUsed:%10d\nmemFree:%10d\n", s.Uptime, s.MemTotal, s.MemUsed, s.MemFree)
	return val
}

func (s *SysInfo) GblString() string {
	var val string
	val = fmt.Sprintf("SysemTime = [%s]\nHostname = [%s]\nUptime = [%s]\nMemUsed = [%02d]", s.SysTime, s.Name, s.Uptime, s.MemUsed)
	return val
}

type CrmNode struct {
	XMLName          xml.Name `xml:"node"`
	Name             string   `xml:"name,attr"`
	Id               string   `xml:"id,attr"`
	Online           string   `xml:"online,attr"`
	Standby          string   `xml:"standby,attr"`
	StandbyOnFail    string   `xml:"standby_onfail,attr"`
	Pending          string   `xml:"pending,attr"`
	Unclean          string   `xml:"unclean,attr"`
	Shutdown         string   `xml:"shutdown,attr"`
	ExpectedUp       string   `xml:"expected_up,attr"`
	ResourcesRunning string   `xml:"resources_running,attr"`
}

type CrmNodes struct {
	XMLName xml.Name  `xml:"nodes"`
	Nodes   []CrmNode `xml:"node"`
}

type CrmResources struct {
	XMLName   xml.Name      `xml:"resources"`
	Resources []CrmResource `xml:"resource"`
}

type CrmResource struct {
	XMLName       xml.Name `xml:"resource"`
	Id            string   `xml:"id,attr"`
	ResourceAgent string   `xml:"resource_agent,attr"`
	Role          string   `xml:"role,attr"`
	Active        string   `xml:"active,attr"`
	Node          CrmNode  `xml:"node"`
}

type CrmCurrentDc struct {
	XMLName    xml.Name `xml:"current_dc"`
	Present    string   `xml:"present,attr"`
	Name       string   `xml:"name,attr"`
	Id         string   `xml:"id,attr"`
	WithQuorum string   `xml:"with_quorum,attr"`
}

type CrmSummary struct {
	XMLName   xml.Name     `xml:"summary"`
	CurrentDC CrmCurrentDc `xml:"current_dc"`
}

type CrmMon struct {
	XMLName   xml.Name     `xml:"crm_mon"`
	Summary   CrmSummary   `xml:"summary"`
	Nodes     CrmNodes     `xml:"nodes"`
	Resources CrmResources `xml:"resources"`
}

func (c *CrmMon) String() string {
	var nodes []string
	var resources []string

	for i := 0; i < len(c.Nodes.Nodes); i++ {
		nodes = append(nodes, c.Nodes.Nodes[i].String())
	}

	for i := 0; i < len(c.Resources.Resources); i++ {
		resources = append(nodes, c.Resources.Resources[i].String())
	}

	return fmt.Sprintf("Nodes\n%s\nResources\n%s", strings.Join(nodes, ""), strings.Join(resources, ""))
}

func (c *CrmMon) GblString() string {
	return fmt.Sprintf("ActiveNode = [%s]", c.Summary.CurrentDC.Name)
}

func (c *CrmNodes) String() string {

	var nodes []string

	for i := 0; i < len(c.Nodes); i++ {
		nodes = append(nodes, c.Nodes[i].String())
	}

	return strings.Join(nodes, "")
}

func (crmNode CrmNode) String() string {
	return fmt.Sprintf("\t Name : %s - Id : %s - Online : %s - Standby : %s \n", crmNode.Name, crmNode.Id, crmNode.Online, crmNode.Standby)
}

func (crmResource CrmResource) String() string {
	return fmt.Sprintf("\t Id : %s - Role : %s - Active : %s\n", crmResource.Id, crmResource.Role, crmResource.Active)
}

type PriSpan struct {
	Span     string
	DChannel string
	Order    string
	Active   string
	Up       string
}

func (priSpan *PriSpan) String() string {
	return fmt.Sprintf("Span : %s\tDChannel:%s\tOrder:%s\tActive:%s\tUp:%s\n",
		priSpan.Span, priSpan.DChannel, priSpan.Order, priSpan.Active, priSpan.Up)
}

type SipPeer struct {
	IPaddress string
	IPport    string
	Name      string
	Status    string
}

func (sipPeer *SipPeer) String() string {
	return fmt.Sprintf("Name : %s\tIPaddress:%s\tIPport:%s\tStatus:%s\n",
		sipPeer.Name, sipPeer.IPaddress, sipPeer.IPport, sipPeer.Status)
}

type IaxTrunk struct {
	Username string
	Refresh  string
	State    string
	Host     string
}

func (i *IaxTrunk) String() string {
	return fmt.Sprintf("Username : %s\tRefresh : %s\tState : %s\tHost : %s\n",
		i.Username, i.Refresh, i.State, i.Host)
}

type SipTrunk struct {
	Username         string
	State            string
	Host             string
	Port             string
	Refresh          string
	RegistrationTime string
}

func (s *SipTrunk) String() string {
	return fmt.Sprintf("Username : %s\tState : %s\tHost : %s\tPort : %s\tRefresh : %s\tRegistrationTime:%s\n",
		s.Username, s.State, s.Host, s.Port, s.Refresh, s.RegistrationTime)
}

type AsteriskInfo struct {
	StartUptime    string
	LastReload     string
	ActiveCalls    string
	ProcessedCalls string
	PriSpans       []PriSpan
	SipPeers       []SipPeer
	IaxTrunks      []IaxTrunk
	SipTrunks      []SipTrunk
}

func NewAsteriskInfo() *AsteriskInfo {
	astInfo := &AsteriskInfo{}
	return astInfo
}

func (a *AsteriskInfo) String() string {
	s := fmt.Sprintf("StartUptime:%s\nLastRelaod:%s\nActiveCalls:%s\nProcessedCalls:%s\n",
		a.StartUptime, a.LastReload, a.ActiveCalls, a.ProcessedCalls)

	var priSpans []string
	var sipPeers []string
	var iaxTrunks []string
	var sipTrunks []string

	for i := 0; i < len(a.PriSpans); i++ {
		priSpans = append(priSpans, a.PriSpans[i].String())
	}

	for i := 0; i < len(a.SipPeers); i++ {
		sipPeers = append(sipPeers, a.SipPeers[i].String())
	}

	for i := 0; i < len(a.IaxTrunks); i++ {
		iaxTrunks = append(iaxTrunks, a.IaxTrunks[i].String())
	}

	for i := 0; i < len(a.SipTrunks); i++ {
		sipTrunks = append(sipTrunks, a.SipTrunks[i].String())
	}

	s = fmt.Sprintf("%s\nPriSpans\n%s\nSipPeer\n%s\nSipTrunks\n%s\nIaxTrunks\n%s\n", s,
		strings.Join(priSpans, ""), strings.Join(sipPeers, ""), strings.Join(sipTrunks, ""), strings.Join(iaxTrunks, ""))
	return s

}

func (a *AsteriskInfo) GblString() string {
	var priSpans []string
	var sipPeers []string
	var voipTrunks []string

	priSpans = append(priSpans, fmt.Sprintf("Spans = [%d]\n", len(a.PriSpans)))

	for i := 0; i < len(a.PriSpans); i++ {
		upDown := "Down"

		if a.PriSpans[i].Up == "Yes" {
			upDown = "Up"
		}

		tmp := fmt.Sprintf("SPAN = [%s:%s]\n", a.PriSpans[i].Span, upDown)

		priSpans = append(priSpans, tmp)
	}

	for i := 0; i < len(a.SipPeers); i++ {
		status := "Off"

		if strings.Contains(a.SipPeers[i].Status, "OK (") {
			status = "On"
		}

		peer := fmt.Sprintf("Peer = [%s:%s]\n", a.SipPeers[i].Name, status)

		sipPeers = append(sipPeers, peer)
	}

	voipTrunks = append(voipTrunks, fmt.Sprintf("Trunks = [%d]\n", len(a.IaxTrunks)+len(a.SipTrunks)))

	for i := 0; i < len(a.IaxTrunks); i++ {
		voipTrunk := fmt.Sprintf("Trunk = [%s:IAX:%s]\n", a.IaxTrunks[i].Username, a.IaxTrunks[i].State)

		voipTrunks = append(voipTrunks, voipTrunk)

	}

	for i := 0; i < len(a.SipTrunks); i++ {
		voipTrunk := fmt.Sprintf("Trunk = [%s:SIP:%s]\n", a.SipTrunks[i].Username, a.SipTrunks[i].State)

		voipTrunks = append(voipTrunks, voipTrunk)

	}

	s := fmt.Sprintf("StartUptime = [%s]\nLastReload = [%s]\nActiveCalls = [%s]\nProcessedCalls = [%s]\n%s%s%s\n",
		a.StartUptime, a.LastReload, a.ActiveCalls, a.ProcessedCalls, strings.Join(priSpans, ""), strings.Join(voipTrunks, ""), strings.Join(sipPeers, ""))
	return s
}
