package force

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RESTError struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode`
}

func restGetJSON(conn *Conn, url string, output interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", conn.SessionId))
	client := &http.Client{}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	fmt.Println(res.StatusCode)
	decoder := json.NewDecoder(res.Body)
	if res.StatusCode != 200 {
		errs := []RESTError{}
		err := decoder.Decode(&errs)
		if err != nil || len(errs) == 0 {
			resTxt, err2 := ioutil.ReadAll(res.Body)
			if err2 != nil {
				return errors.New("unable to read response body")
			}
			return errors.New(string(resTxt))
		}
		return errors.New(errs[0].Message)
	}

	err = decoder.Decode(&output)
	return err
}
