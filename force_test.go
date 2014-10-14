package force

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

type credsConf struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Version  string `json:"version"`
}

func getCredsFromFile() (credsConf, error) {
	creds := credsConf{}
	file, err := os.Open("creds.json")
	if err != nil {
		fmt.Println("cannot open file containing test creds, needs 'creds.json' in current directory")
		return creds, err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&creds)
	if err != nil {
		fmt.Println("error parsing creds file, check format")
	}
	return creds, err
}

var _conn *Conn

func getTestConn() *Conn {
	if _conn == nil {
		creds, _ := getCredsFromFile()
		_conn, _ = Login(creds.Username, creds.Password, creds.Version, PRODUCTION)
	}
	return _conn
}

func TestLogin_invalid(t *testing.T) {
	fmt.Println("testing login invalid")
	conn, err := Login("", "", "", PRODUCTION)
	fmt.Println(err)
	if err == nil || conn != nil {
		t.Fatal("should have returned an error")
	}
}

func TestLogin_valid_production(t *testing.T) {
	fmt.Println("testing login valid production")
	creds, err := getCredsFromFile()
	if err != nil {
		t.Fatal("could not fetch test creds")
	}

	conn, err := Login(creds.Username, creds.Password, creds.Version, PRODUCTION)
	if err != nil || len(conn.SessionId) == 0 {
		t.Log(err)
		t.Fatal("should not have errored")
	}
}

func TestQuery_invalid(t *testing.T) {
	fmt.Println("testing invalid query")
	conn := getTestConn()
	_, err := Query(conn, "select foo from bar")
	fmt.Println(err)
	if err == nil {
		t.Fatal("should have thrown error")
	}
}

func TestQuery_valid_onepage(t *testing.T) {
	fmt.Println("testing valid query with one page")
	conn := getTestConn()
	r, err := Query(conn, "select Id,Name from account limit 5")
	fmt.Printf("%#v \n", r)
	if err != nil {
		t.Fatal("query should have succeeded")
	}
	if !r.Done || len(r.QueryLocator) > 0 {
		t.Fatal("query should be DONE")
	}
}

func TestQuery_valid_pages(t *testing.T) {
	fmt.Println("testing valid query with multi pages")
	conn := getTestConn()
	r, err := Query(conn, "select Id,Name,description from account")
	fmt.Println(r.Done, r.QueryLocator)
	if err != nil {
		t.Fatal("query should have succeeded")
	}
	if r.Done || len(r.QueryLocator) == 0 {
		t.Fatal("query should NOT be DONE")
	}
	r, err = QueryMore(conn, r.QueryLocator)
	// fmt.Println(r)
	if err != nil {
		t.Fatal(err)
	}
}
