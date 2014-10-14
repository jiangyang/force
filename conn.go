package force

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type OrgType int

const (
	PRODUCTION = iota
	SANDBOX
)

type Conn struct {
	SessionId         string `xml:"Body>loginResponse>result>sessionId"`
	Sandbox           bool   `xml:"Body>loginResponse>result>sandbox"`
	ServerUrl         string `xml:"Body>loginResponse>result>serverUrl"`
	ServiceUrl        string
	InstanceUrl       string
	MetadataServerUrl string `xml:"Body>loginResponse>result>metadataServerUrl"`
	UserId            string `xml:"Body>loginResponse>result>userId"`
}

type soapError struct {
	FaultCode   string `xml:"Body>Fault>faultcode"`
	FaultString string `xml:"Body>Fault>faultstring"`
}

// username password and version
// version is without "v" e.g. 28.0
// org type: PRODUCTION or SANDBOX
func Login(username, password, ver string, orgType OrgType) (*Conn, error) {
	if len(username) == 0 || len(password) == 0 {
		return nil, errors.New("empty credentials")
	}

	if len(ver) == 0 {
		ver = "28.0"
	}

	template := `<?xml version="1.0" encoding="utf-8" ?>
        <env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema"
            xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
            xmlns:env="http://schemas.xmlsoap.org/soap/envelope/">
            <env:Body>
                <n1:login xmlns:n1="urn:partner.soap.sforce.com">
                    <n1:username>%s</n1:username>
                    <n1:password>%s</n1:password>
                </n1:login>
            </env:Body>
        </env:Envelope>`

	reqBody := fmt.Sprintf(template, username, password)

	endpointf := ""
	switch orgType {
	case PRODUCTION:
		endpointf = "https://login.salesforce.com/services/Soap/u/%s"
	case SANDBOX:
		endpointf = "https://test.salesforce.com/services/Soap/u/%s"
	default:
		return nil, errors.New("unsupported org type")
	}

	endpoint := fmt.Sprintf(endpointf, ver)
	// fmt.Println("end point:", endpoint)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml")
	req.Header.Set("SOAPAction", "login")
	client := &http.Client{}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	resTxt, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	e := soapError{}
	xml.Unmarshal(resTxt, &e)
	if len(e.FaultCode) > 0 {
		return nil, fmt.Errorf("%s - %s", e.FaultCode, e.FaultString)
	}

	c := &Conn{}
	xml.Unmarshal(resTxt, c)
	// serviceUrl is the base url for rest
	c.ServiceUrl = c.ServerUrl[0:strings.Index(c.ServerUrl, "Soap")] + "data/v" + ver
	// instanceUrl is the base url
	c.InstanceUrl = c.ServerUrl[0:strings.Index(c.ServerUrl, "/services")]
	return c, nil
}
