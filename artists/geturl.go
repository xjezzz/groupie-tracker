package artists

import (
	"errors"
	"io/ioutil"
	"net/http"
)

func GetUrl(url string) ([]byte, error) {
	var Error error
	netClient := http.Client{}
	resp, err := netClient.Get(url)
	if err != nil {
		Error = errors.New("error")
		return nil, Error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Error = errors.New("error")
		return nil, Error
	}
	return body, nil
}
