#!/usr/bin/python
from BaseHTTPServer import BaseHTTPRequestHandler,HTTPServer
import urllib
import os
import json
import time

def Calculator(dwHeight, dwWeight):
    ddwHeight = float(dwHeight)
    ddwWeight = float(dwWeight)
    ddwHeight /= 100
    ddwBmi = ddwWeight / (ddwHeight * ddwHeight)
    return float('%.2f' % ddwBmi)

class CalculatorHandler(BaseHTTPRequestHandler):
    #Handler for the GET requests
    def do_GET(self):
        mpath,margs=urllib.splitquery(self.path)
        print mpath
        print margs
        margs.replace("&&", "&")
        inputPara = margs.split("&", 1)
        print inputPara
        if len(inputPara) != 2:
            print 'err input:' + inputPara
            return
        arrHeight = str(inputPara[0]).split("=")
        arrWeight = str(inputPara[1]).split("=")
        if len(arrHeight) != 2 or  len(arrWeight) != 2 :
            return
        dwHeight = float(arrHeight[1])
        dwWeight = float(arrWeight[1])
        ddwBmi = Calculator(dwHeight, dwWeight)
        print "calculator result:" + str(ddwBmi)
        self.send_response(200)
        self.send_header('Content-type','application/json')
        self.send_header('Access-Control-Allow-Origin', '*')
        self.end_headers()
        # Send the html message
        date = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime())
        arrDate = str.split(date, " ")
        result = {"result": ddwBmi, "instanceId": "pythonServer", "callTime": str(arrDate[1])}
        strResult = json.dumps(result)
        print "json result:" + strResult
        self.wfile.write(strResult)
        return

try:
    server = HTTPServer(('127.0.0.1', 4540), CalculatorHandler)
    print 'http server begin:\n'
    server.serve_forever()

except KeyboardInterrupt:
    print 'stop'
    server.socket.close()