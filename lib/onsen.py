# -*- coding: utf-8 -*-
import requests
import json
import re
import os
import datetime as DT
import subprocess
from . import functions as f

class onsen:
    def __init__(self, keywords, SAVEROOT):
        self.change_keywords(keywords)
        self.SAVEROOT = SAVEROOT
        self.reload_date = DT.date.today()

    def change_keywords(self, keywords):
        if bool(keywords):
            word = "("
            for keyword in keywords:
                word += keyword
                word += "|"
            word = word.rstrip("|")
            word += ")"
            self.isKeyword = True
            self.keyword = re.compile(word)
        else:
            self.isKeyword = False

    def rec(self):
        self.reload_date = DT.date.today()
        res = requests.get("http://www.onsen.ag/api/shownMovie/shownMovie.json")
        res.encoding = "utf-8"
        programs = json.loads(res.text)
        returnData = []
        for program in programs["result"]:
            url = "http://www.onsen.ag/data/api/getMovieInfo/%s" % program
            res2 = requests.get(url)
            prog = json.loads(res2.text[9:len(res2.text)-3])
            title = prog.get("title")
            personality = prog.get("personality")
            update_DT = prog.get("update")
            count = prog.get("count")
            if (title is not None and personality is not None and update_DT !=""):
                if (self.keyword.search(title) or self.keyword.search(personality)):
                    movie_url = prog["moviePath"]["pc"]
                    if (movie_url == ""):
                        continue
                    # title の長さ
                    title = title[:30]
                    # フォルダの作成
                    dir_name = title.replace(" ", "_")
                    dir_path = self.SAVEROOT + "/" + dir_name
                    f.createSaveDir(dir_path)
                    # ファイル重複チェック
                    file_name = title.replace(" ", "_") +"#"+ count + ".mp3"
                    file_path = dir_path +"/"+ file_name
                    if (f.did_record_prog(file_path, title, update_DT)):
                        continue
                    # print(prog["update"], prog["title"], prog["personality"])
                    returnData.append(title)
                    res3 = requests.get(movie_url)
                    fs = open(file_path, "wb")
                    fs.write(res3.content)
                    fs.close()
                    f.Mysql.insert(
                        title= title,
                        pfm= personality,
                        timestamp= update_DT,
                        station= "onsen",
                        uri= "http://example.com"
                    )
        return returnData
