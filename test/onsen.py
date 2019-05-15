# -*- coding: utf-8 -*-
import requests
import json
import re

path = "./conf/config.json"
fs = open(path, "r")
js = json.load(fs)
word = "("
for keyword in js["Onsen"]["keywords"]:
    word += keyword
    word += "|"
word = word.rstrip("|")
word += ")"
c = re.compile(word)

res = requests.get("http://www.onsen.ag/api/shownMovie/shownMovie.json")
res.encoding = "utf-8"
programs = json.loads(res.text)
for program in programs["result"]:
    url = "http://www.onsen.ag/data/api/getMovieInfo/%s" % program
    res2 = requests.get(url)
    prog = json.loads(res2.text[9:len(res2.text)-3])
    title = prog.get("title")
    personality = prog.get("personality")
    update_DT = prog.get("update")
    count = prog.get("count")
    if (title is not None and personality is not None and update_DT !=""):
        if (c.search(personality)):
            print(update_DT, title, personality)
            """movie_url = prog["moviePath"]["pc"]
            res3 = requests.get(movie_url)
            f = open(path+title+".mp3", "wb")
            f.write(res3.content)
            f.close()"""
