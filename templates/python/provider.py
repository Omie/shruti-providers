#!/usr/bin/python

import json
import os
import requests
import time

from posixpath import join as urljoin

PROVIDER_NAME = "<PROVIDER NAME>"
ENDPOINT_REGISTER = "providers"
ENDPOINT_PUSH = "notifications"

headers = {'content-type': 'application/json'}
visitedStories = {}

def register(shruti_server):
    _data = { "name" : PROVIDER_NAME,
              "display_name" : "<DISPLAY NAME>",
              "description": "<DESCRIPTION>",
              "web_url": "<URL>",
              "icon_url": "<URL>",
              "voice": "<VoiceName>"
        }
    _url = urljoin(shruti_server, ENDPOINT_REGISTER, PROVIDER_NAME)
    print _url
    r = requests.post(_url, data=json.dumps(_data), headers=headers)
    if r.text == 'pq: duplicate key value violates unique constraint "providers_name_key"':
        return 201, "already registered"
    return r.status_code, r.text

def doWork(shruti_server):

    _url = urljoin(shruti_server, ENDPOINT_PUSH)
    print _url

    while True:
        _resp = requests.get("<IMPLEMENT THIS>")

        #submit story to /notifications
        _data = { "title" : "<TITLE TEXT>",
                  "url": "<URL>",
                  "key": PROVIDER_NAME + "<KEY>",
                  "priority": 20,
                  "action": 10,
                  "provider_name": PROVIDER_NAME
                  }
        r = requests.post(_url, data=json.dumps(_data), headers=headers)
        visitedStories["<KEY>"] = True

        msg = "Error submitting story:" + r.text if r.status_code == 500 else "submitted"
        print msg

        time.sleep(15 * 60) # ~15minutes
        print "time to check updates again"

def main():
    shruti_server = None
    shruti_server = os.environ['SHRUTI_SERVER']
    if shruti_server is None:
        print "server url not set"
        return
    status, resp = register(shruti_server)
    print status, resp
    if status == 500:
        print "Error registering:", resp
        return
    doWork(shruti_server)

if __name__ == "__main__":
    main()

