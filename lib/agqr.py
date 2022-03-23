# -*- coding: utf-8 -*-
import requests
import json
import re
import datetime as DT
import time
import subprocess
import os
from . import functions as f


class agqr:
    AGQR_URL = "https://agqr.sun-yryr.com/api/today"
    def __init__(self):
        self.reload_program()
        self.isKeyword = False

    def reload_program(self):
        res = requests.get(self.AGQR_URL)
        res.encoding = "utf-8"
        self.program_agqr = json.loads(res.text)
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

    def delete_keywords(self):
        self.change_keywords([])

    def search(self):
        if (self.isKeyword is False): return []
        res = []
        for prog in self.program_agqr:
            ck = False
            title = prog.get("title")
            pfm = prog.get("pfm")
            if (self.keyword.search(title)):    ck = True
            if (ck is False) and (pfm is not None):
                if (self.keyword.search(pfm)):  ck = True
            if (ck):
                res.append({
                    "title": title.replace(" ", "_"),
                    "ft": prog.get("ft"),
                    "DT_ft": DT.datetime.strptime(prog.get("ft"), "%Y%m%d%H%M"),
                    "to": prog.get("to"),
                    "dur": int(prog.get("dur")),
                    "pfm": pfm
                })
        if bool(res): return res
        else: return []

    def rec(self, data):
        program_data = data[0]
        wait_start_time = data[1]
        SAVEROOT = data[2]

        print(program_data["title"])

        dir_name = program_data["title"].replace(" ", "_")
        dir_path = SAVEROOT + "/" + dir_name
        f.createSaveDir(dir_path)

        file_path = dir_path + "/" + program_data["title"].replace(" ", "_") + "_" + program_data["ft"][:12]

        # recording....
        # TODO: 度々変更されるので環境変数から読み込む
        url = "https://fms2.uniqueradio.jp/agqr10/aandg1.m3u8"
        duration = int(program_data["dur"]) * 60
        cwd = ('ffmpeg -loglevel error -i "%s" -movflags faststart -t %s  "%s.mp3"' % (url, duration, file_path))
        time.sleep(wait_start_time)
        print("Agqr: recording start")
        subprocess.run(cwd, stdin=subprocess.PIPE, stdout=subprocess.DEVNULL, shell=True)
        print("Agqr: finished!")
        time.sleep(10)

        if (f.is_recording_succeeded(file_path)):
            f.Mysql.insert(
                title= program_data["title"].replace(" ", "_"),
                pfm= program_data["pfm"],
                timestamp= program_data["ft"] + "00",
                station= "AGQR",
                uri= "http://example.com" # TODO: 削除する
            )

# サンプル: pipenv run python -m lib.agqr "rec-sample" "2020-01-01" "1" "./"
if __name__ == "__main__":
    import sys
    args = sys.argv
    title = args[1]
    ft = args[2]
    dur = (int) (args[3])
    path = args[4]
    
    program_data = {
        "title": title,
        "ft": ft,
        "dur": dur,
        "pfm": "sample"
    }
    
    client = agqr()
    client.rec([
        program_data,
        0,
        path
    ])
