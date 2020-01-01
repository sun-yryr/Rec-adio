# -*- coding: utf-8 -*-
import requests
import json
import re
import datetime as DT
import lib.functions as f
import time
import subprocess
import os


class agqr:
    AGQR_URL = "https://agqr.sun-yryr.com/api/today"
    def __init__(self):
        res = requests.get(self.AGQR_URL)
        res.encoding = "utf-8"
        self.isKeyword = False
        self.reload_date = DT.date.today()
        self.program_agqr = json.loads(res.text)

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
            print(word)
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
        #print(program_data)
        wait_start_time = data[1]
        SAVEROOT = data[2]

        dir_path = SAVEROOT + "/" + program_data["title"].replace(" ", "_")
        f.createSaveDir(dir_path)

        file_path = dir_path + "/" + program_data["title"].replace(" ", "_") + "_" + program_data["ft"][:12]
        cwd  = ('rtmpdump -r rtmp://fms-base1.mitene.ad.jp/agqr/aandg1b ')
        cwd += ('--stop %s ' % str(program_data["dur"]*60))
        cwd += ('--live -o "%s.flv"' % (file_path))
        time.sleep(wait_start_time)
        #rtmpdumpは時間指定の終了ができるので以下を同期処理にする
        subprocess.run(cwd, shell=True)
        #変換をする
        cwd2 = ('ffmpeg -loglevel error -i "%s.flv" -vn -acodec copy "%s.m4a"' % (file_path, file_path))
        subprocess.run(cwd2, shell=True)
        print("agqr finish!")
        if (f.is_recording_succeeded(file_path)):
            f.recording_successful_toline(program_data["title"])
            fs = open(file_path+".m4a", "rb")
            f.DropBox.upload(program_data["title"], program_data["ft"], fs.read())
            fs.close()
        else:
            f.recording_failure_toline(program_data["title"])
        os.remove(file_path + ".flv")