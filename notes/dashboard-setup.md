#!/usr/bin/env python

import json

searchresults = json.loads('@option.searchresults@')

d = {}
for bucket in searchresults["aggregations"]["unique_hostname"]["buckets"]:
    d[bucket["key"]] = bucket["last_timestamp"]["value_as_string"]

j = json.dumps(d)
print(j)

#r = requests.post(postURL, j) #print(r.status_code, r.reason)
#print(r.status_code, r.reson)
