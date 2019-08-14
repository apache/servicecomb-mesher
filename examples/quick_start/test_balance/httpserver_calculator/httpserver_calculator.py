#!/usr/bin/python
'''
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
'''

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
        if dwHeight < 0 or dwWeight < 0 :
            time.sleep(6)
            raise RuntimeError('para Error')
            return 
        ddwBmi = Calculator(dwHeight, dwWeight)
        print "calculator result:" + str(ddwBmi)
        self.send_response(200)
        self.send_header('Content-type','application/json')
        self.send_header('Access-Control-Allow-Origin', '*')
        self.end_headers()
        # Send the html message
        date = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime())
        arrDate = str.split(date, " ")
        result = {"result": ddwBmi, "instanceId": "pythonServer2", "callTime": str(arrDate[1])}
        strResult = json.dumps(result)
        print "json result:" + strResult
        self.wfile.write(strResult)
        return

try:
    server = HTTPServer(('127.0.0.1', 4537), CalculatorHandler)
    print 'http server begin:\n'
    server.serve_forever()

except KeyboardInterrupt:
    print 'stop'
    server.socket.close()