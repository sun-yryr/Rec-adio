# -*- coding: utf-8 -*-
import requests
import json
import re
import os
c = re.compile("(ガルパ)")
headers = {"X-Requested-With": "XMLHttpRequest"}
res = requests.get("https://vcms-api.hibiki-radio.jp/api/v1/programs", headers=headers)
js = json.loads(res.text)
for program in js:
    episode = program.get("episode")
    if episode is None:
        continue
    title = program.get("name")
    personality = program.get("cast")
    if (c.search(title)):
        url = "https://vcms-api.hibiki-radio.jp/api/v1/videos/play_check?video_id=" + str(program["id"])
        url2 = "https://vcms-api.hibiki-radio.jp/api/v1/programs/" + program.get("access_id")
        url3 = "https://vcms-api.hibiki-radio.jp/api/v1/product_informations?program_id=" + str(program["id"])
        res2 = requests.get(url2, headers=headers)
        print(res2.text)
