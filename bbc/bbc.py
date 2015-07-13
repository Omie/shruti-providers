#!/usr/bin/python

import feedparser
import json
import os
import requests
import time

from posixpath import join as urljoin

PROVIDER_NAME = "bbc"
ENDPOINT_REGISTER = "providers"
ENDPOINT_PUSH = "notifications"

BBC_FEED_URL = 'http://feeds.bbci.co.uk/news/world/rss.xml?edition=uk'

headers = {'content-type': 'application/json'}
visitedStories = {}

def register(shruti_server):
    _data = { "name" : PROVIDER_NAME,
              "display_name" : "BBC",
              "description": "world news headlines from BBC",
              "web_url": "http://bbc.com",
              "icon_url": "http://static.bbci.co.uk/frameworks/barlesque/2.83.10/desktop/3.5/img/blq-blocks_grey_alpha.png"
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
        feed = feedparser.parse(BBC_FEED_URL)
        for entry in feed['entries'][:3]:
            # extract news id from its guid
            _guid = entry["guid"]
            story_id = _guid.split("/")[-1] #last part
            print story_id
            visited = visitedStories.get(story_id, False)
            if not visited:
                _data = { "title" : entry["title"],
                          "url": entry["link"],
                          "key": PROVIDER_NAME + story_id,
                          "priority": 20,
                          "action": 10,
                          "provider": PROVIDER_NAME
                          }
                r = requests.post(_url, data=json.dumps(_data), headers=headers)
                visitedStories[story_id] = True
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

