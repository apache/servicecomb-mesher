package health

import (
	"errors"
	"github.com/go-chassis/foundation/httpclient"
	"github.com/go-mesh/mesher/proxy/config"
	"github.com/go-mesh/openlogging"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

//HTTPCheck checks http service
func HTTPCheck(check *config.HealthCheck, address string) error {
	c, err := httpclient.GetURLClient(&httpclient.DefaultURLClientOption)
	if err != nil {
		openlogging.Error("can not get http client: " + err.Error())
		//must not return error, because it is mesher error
		return nil
	}
	var url = "http://" + address
	if check.URI != "" {
		url = url + check.URI
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		openlogging.Error("can not get http req: " + err.Error())
		//must not return error, because it is mesher error
		return nil
	}
	resp, err := c.Do(req)
	if err != nil {
		openlogging.Error("server can not be connected: " + err.Error())
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if check.Match != nil {
		if check.Match.Status != "" {
			n, _ := strconv.Atoi(check.Match.Status)
			if resp.StatusCode != n {
				return errors.New("status is not " + check.Match.Status)
			}
		}
		if check.Match.Body != "" {
			re := regexp.MustCompile(check.Match.Body)
			if !re.Match(body) {
				return errors.New("body does not match " + check.Match.Body)
			}
		}
	} else {
		if resp.StatusCode == 200 {
			return nil
		}
	}
	return nil
}
