package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Result struct {
	Result     float64 `json:"result"`
	InstanceId string  `json:"instanceId"`
	CallTime   string  `json:"callTime"`
}

func HandlerGateway(w http.ResponseWriter, r *http.Request) {

	queryInfo, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Println("parse queryInfo wrong: ", queryInfo)
		fmt.Fprintln(w, "parse queryInfo wrong")
		return
	}
	fmt.Println("height:%s,weight:%s", queryInfo.Get("height"), queryInfo.Get("weight"))
	strReqUrl := "http://mersher-provider:4540/bmi?height=" + queryInfo.Get("height") + "&weight=" + queryInfo.Get("weight")
	resp, err := http.Get(strReqUrl)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println("body err: ", body, err)
		return
	}
	fmt.Println("body: " + string(body))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(body)

}

func HandlerGateway2(w http.ResponseWriter, r *http.Request) {
	queryInfo, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Println("parse queryInfo wrong: ", queryInfo)
		fmt.Fprintln(w, "parse queryInfo wrong")
		return
	}
	fmt.Println("height:%s,weight:%s", queryInfo.Get("height"), queryInfo.Get("weight"))
	strReqUrl := "http://mersher-provider/bmi?height=" + queryInfo.Get("height") + "&weight=" + queryInfo.Get("weight")
	//strReqUrl := "http://mersher-ht-provider/bmi?height=" + queryInfo.Get("height") + "&weight=" + queryInfo.Get("weight")
	//strReqUrl := "http://mersher-ht-provider:4555/bmi?height=" + queryInfo.Get("height") + "&weight=" + queryInfo.Get("weight")
	proxy, _ := url.Parse("http://127.0.0.1:30101") //将mesher设置为http代理
	httpClient := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	req, err := http.NewRequest(http.MethodGet, strReqUrl, nil)
	if err != nil {
		fmt.Println("make req err: ", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Do err: ", resp, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println("body err: ", body, err)
		return
	}
	fmt.Println("body: " + string(body))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(body)
}

func main() {
	http.HandleFunc("/bmi", HandlerGateway2)
	http.ListenAndServe("192.168.88.64:4538", nil)
}