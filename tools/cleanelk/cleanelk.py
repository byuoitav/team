#!/usr/bin/python3

import requests
import os
import re
from datetime import datetime, timedelta

elkaddr = os.environ["ELK_ADDR"]
indexPrefixes = ["av-all-events", "legacy-av-all-events"]
daysToRetain = 14

cutoff = datetime.now() - timedelta(days=daysToRetain)

#Get all of the indexes with the prefixes
for cur in indexPrefixes:
    resp = requests.get(elkaddr + "/_cat/indices/" + cur + "*")
    if resp.status_code / 100 != 2:
        print("error")
        print(resp.text)
        continue
    r = re.compile(cur + "-[\d]{8}")
    matches = r.findall(resp.text)
    for j in matches:
        c = datetime.strptime(j[j.rfind("-")+1:], "%Y%m%d")
        if cutoff > c:
            #We delete
            print("Deleting: ", j)
            m = requests.delete(elkaddr + "/" + j)
