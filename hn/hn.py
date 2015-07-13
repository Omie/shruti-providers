#!/usr/bin/python

import json
import os
import requests
import time

from posixpath import join as urljoin

PROVIDER_NAME = "hackernews1"
ENDPOINT_REGISTER = "providers"
ENDPOINT_PUSH = "notifications"

HN_ALL_STORIES = 'https://hacker-news.firebaseio.com/v0/topstories.json'
HN_SINGLE_STORY = 'https://hacker-news.firebaseio.com/v0/item/{0}.json'

headers = {'content-type': 'application/json'}
visitedStories = {}

def register(shruti_server):
    _data = { "name" : PROVIDER_NAME,
              "display_name" : "Hacker News",
              "description": "updates around the globe, mostly tech related",
              "web_url": "https://news.ycombinator.com/",
              "icon_url": "https://news.ycombinator.com/favicon.ico"
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
        _resp = requests.get(HN_ALL_STORIES).json()
        for storyid in _resp[:10]:
            visited = visitedStories.get(storyid, False)
            if not visited:
                # make another request to get details of the story
                _resp = requests.get(HN_SINGLE_STORY.format(storyid)).json()

                _data = { "title" : _resp['title'],
                          "url": _resp['url'],
                          "key": PROVIDER_NAME + str(_resp['id']),
                          "priority": 20,
                          "action": 10,
                          "provider": PROVIDER_NAME
                          }
                r = requests.post(_url, data=json.dumps(_data), headers=headers)

                msg = "Error submitting story:" + r.text if r.status_code == 500 else "submitted"
                print msg

            visitedStories[storyid] = True
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

