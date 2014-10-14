package force

import (
	"errors"
	"fmt"
	"net/url"
)

type QueryResult struct {
	Done         bool          `json:"done"`
	TotalSize    int           `json:"totalSize"`
	QueryLocator string        `json:"nextRecordsUrl"`
	Records      []interface{} `json:"records"`
}

func validateConn(conn *Conn) error {
	if conn == nil || len(conn.SessionId) == 0 || len(conn.ServiceUrl) == 0 {
		return errors.New("invalid connection")
	}
	return nil
}

func Query(conn *Conn, query string) (*QueryResult, error) {
	err := validateConn(conn)
	if err != nil {
		return nil, err
	}
	url := conn.ServiceUrl + "/query?q=" + url.QueryEscape(query)
	fmt.Println(url)

	queryResult := QueryResult{}
	err = restGetJSON(conn, url, &queryResult)
	if err != nil {
		return nil, err
	}
	return &queryResult, nil
}

func QueryMore(conn *Conn, queryLocator string) (*QueryResult, error) {
	err := validateConn(conn)
	if err != nil {
		return nil, err
	}
	if len(queryLocator) == 0 {
		return nil, fmt.Errorf("empty query locator")
	}

	url := conn.InstanceUrl + queryLocator
	fmt.Println(url)
	queryResult := QueryResult{}
	err = restGetJSON(conn, url, &queryResult)
	if err != nil {
		return nil, err
	}
	return &queryResult, nil
}
