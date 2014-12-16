package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

//Cluster informations returned
//It is very simple implementation and the principal job made by crm_mon
//Please see crm_mon documenation and the project INSTALL(md/html/pdf) file for more details
func GetClusterInfo(filepath string) (crmMon CrmMon, err error) {
	xmlFile, err := os.Open(filepath)

	if err != nil {
		return crmMon, err
	}

	defer xmlFile.Close()

	XMLdata, _ := ioutil.ReadAll(xmlFile)

	xml.Unmarshal(XMLdata, &crmMon)

	return crmMon, err
}
