package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Result struct{
	Result float64  `json:"result"`
	InstanceId string `json:"instanceId"`
	CallTime string `json:"callTime"`
}

func HandlerCalculator(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandlerCalculator1 begin")
	queryInfo, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		fmt.Println("parse queryInfo wrong: ", queryInfo)
		fmt.Fprintln(w, "parse queryInfo wrong")
		return
	}
	fmt.Println("height:%s,weight:%s", queryInfo.Get("height"), queryInfo.Get("weight"))
	ddwHeight, err := strconv.ParseFloat(queryInfo.Get("height"), 64)
	if err != nil {
		fmt.Println("para height wrong: %s", queryInfo.Get("height"))
		fmt.Fprintln(w, "para height wrong")
	}
	if ddwHeight <= 0 {
		time.Sleep(10 * time.Second)
		panic("err input height: ")
	}
	ddwHeight /= 100
	ddwWeight, err := strconv.ParseFloat(queryInfo.Get("weight"), 64)
	if err != nil {
		fmt.Println("para height wrong: %s", queryInfo.Get("weight"))
		fmt.Fprintln(w, "para height wrong")
	}
	if ddwWeight <= 0 {
		time.Sleep(10 * time.Second)
		panic("err input weight: ")
	}
	ddwBmi := ddwWeight / (ddwHeight * ddwHeight)
	fmt.Println("bmi:", ddwBmi)
	var result Result
	result.Result, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", ddwBmi), 64)
	strTime := time.Now().Format("2006-01-02 15:04:05")
	arrTime := strings.Split(strTime, " ")
	result = Result{result.Result, "goHttpServer", arrTime[1]}
	bResult, err := json.Marshal(result)
	if err != nil {
		fmt.Println("result err ", result, string(bResult), err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("result ", result, string(bResult))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(bResult)

}

func main() {
	http.HandleFunc("/bmi", HandlerCalculator)
	http.ListenAndServe("127.0.0.1:4537", nil)
}