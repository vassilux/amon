package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

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
