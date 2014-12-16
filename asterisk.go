package main

import (
	gami "code.google.com/p/gami"
	"fmt"
	"net"
	"strings"
	//"time"
)

func ConnectToAsterisk(addr string, port int, username string, password string) (a *gami.Asterisk, con net.Conn, err error) {
	con, err = net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return a, con, err
	}
	//
	a = gami.NewAsterisk(&con, nil)

	err = a.Login(username, password)

	return a, con, err
}

//helper to send an asterisk AMI command with a result contents a message
func sendAsteriskAMICommande(ast *gami.Asterisk, command string) (gami.Message, error) {
	ch := make(chan gami.Message)

	check := func(m gami.Message) {
		ch <- m
	}

	m := gami.Message{"Action": "COMMAND", "Command": command}

	err := ast.SendAction(m, &check)

	if err != nil {
		return nil, err
	}

	return <-ch, nil
}

//Get the asterisk serveur uptime
func getAsteriskSystemUptime(ast *gami.Asterisk) (uptime string, err error) {
	m, err := sendAsteriskAMICommande(ast, "core show uptime")

	if err != nil {
		return "", err
	}
	fmt.Println(m)
	return m["System uptime"], nil
}

func getAsteriskPriSpans(ast *gami.Asterisk) (r []PriSpan, err error) {
	ch := make(chan bool)
	ml := []gami.Message{}
	cscf := func() func(gami.Message) {

		return func(m gami.Message) {
			ml = append(ml, m)
			//fmt.Println(m)
			if m["Event"] == "PRIShowSpansComplete" {
				//
				ch <- true
			}

			if m["Response"] == "Error" {
				ch <- false
			}
		}
	}()
	m := gami.Message{"Action": "PRIShowSpans"}
	ast.HoldCallbackAction(m, &cscf)
	<-ch

	for i := 0; i < len(ml); i++ {
		if ml[i]["Event"] == "PRIShowSpans" {
			r = append(r, PriSpan{
				ml[i]["Span"],
				ml[i]["DChannel"],
				ml[i]["Order"],
				ml[i]["Active"],
				ml[i]["Up"],
			})
		}

		if m["Response"] == "Error" {
			ch <- false
		}

	}

	return r, nil
}

func getAsteriskSipPeers(ast *gami.Asterisk) (r []SipPeer, err error) {
	ch := make(chan bool)
	ml := []gami.Message{}
	cscf := func() func(gami.Message) {

		return func(m gami.Message) {
			//fmt.Println(m)
			ml = append(ml, m)
			if m["Event"] == "PeerlistComplete" {
				//
				ch <- true
			}
		}
	}()

	m := gami.Message{"Action": "Sippeers"}

	ast.HoldCallbackAction(m, &cscf)
	<-ch

	for i := 0; i < len(ml); i++ {
		if ml[i]["Event"] == "PeerEntry" {
			r = append(r, SipPeer{
				ml[i]["IPaddress"],
				ml[i]["IPport"],
				ml[i]["ObjectName"],
				ml[i]["Status"],
			})
		}

		if m["Response"] == "Error" {
			ch <- false
		}

	}

	return r, nil
}

func getSipTrunks(ast *gami.Asterisk) (r []SipTrunk, err error) {

	ch := make(chan bool)

	ml := []gami.Message{}

	cscf := func() func(gami.Message) {

		return func(m gami.Message) {
			//fmt.Println(m)
			ml = append(ml, m)

			if m["Event"] == "RegistrationsComplete" {
				ch <- true
			}
		}
	}()

	m := gami.Message{"Action": "SIPshowregistry"}

	ast.HoldCallbackAction(m, &cscf)

	<-ch
	if len(ml) <= 2 {
		return r, nil
	}

	for i := 0; i < len(ml); i++ {
		if ml[i]["Event"] == "RegistryEntry" {
			r = append(r, SipTrunk{
				ml[i]["Username"],
				ml[i]["State"],
				ml[i]["Host"],
				ml[i]["Port"],
				ml[i]["Refresh"],
				ml[i]["RegistrationTime"],
			})
		}

	}

	return r, nil
}

func getIaxTrunks(ast *gami.Asterisk) (r []IaxTrunk, err error) {
	//
	ch := make(chan bool)

	ml := []gami.Message{}

	cscf := func() func(gami.Message) {

		return func(m gami.Message) {

			ml = append(ml, m)

			if m["Event"] == "RegistrationsComplete" {
				ch <- true
			}
			if m["Response"] == "Error" {
				ch <- false
			}
		}
	}()

	m := gami.Message{"Action": "IAXregistry"}

	ast.HoldCallbackAction(m, &cscf)

	<-ch

	for i := 0; i < len(ml); i++ {
		if ml[i]["Event"] == "RegistryEntry" {
			r = append(r, IaxTrunk{
				ml[i]["Username"],
				ml[i]["Refresh"],
				ml[i]["State"],
				ml[i]["Host"],
			})
		}

		if m["Response"] == "Error" {
			ch <- false
		}

	}

	return r, nil
}

//Get the avtive calls numbers and the processed calls since the asterisk running
func getAsteriskCalls(ast *gami.Asterisk) (activeCalls, processedCalls string, err error) {
	m, err := sendAsteriskAMICommande(ast, "core show calls")

	if err != nil {
		return "", "", err
	}

	cmdData := m["CmdData"]

	words := strings.Fields(cmdData)
	//get tokens for calls indications
	return words[0], words[3], nil
}

//PLEASE see the AsteriskInfo struct about returned informations
func GetAsteriskInfo(ast *gami.Asterisk) (*AsteriskInfo, error) {
	var err error

	astInfo := NewAsteriskInfo()

	astInfo.Uptime, err = getAsteriskSystemUptime(ast)

	if err != nil {
		return astInfo, err
	}

	astInfo.ActiveCalls, astInfo.ProcessedCalls, err = getAsteriskCalls(ast)

	if err != nil {
		return astInfo, err
	}

	astInfo.PriSpans, err = getAsteriskPriSpans(ast)

	if err != nil {
		return astInfo, err
	}

	astInfo.SipPeers, err = getAsteriskSipPeers(ast)

	if err != nil {
		return astInfo, err
	}

	astInfo.IaxTrunks, err = getIaxTrunks(ast)
	if err != nil {
		return astInfo, err
	}

	astInfo.SipTrunks, err = getSipTrunks(ast)
	if err != nil {
		return astInfo, err
	}

	return astInfo, nil
}
