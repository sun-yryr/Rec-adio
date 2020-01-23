# -*- coding: utf-8 -*-
import requests
import json
import re
import os
import lib.functions as f
import datetime as DT
import subprocess

class hibiki:
    def __init__(self, keywords, SAVEROOT):
        self.change_keywords(keywords)
        self.SAVEROOT = SAVEROOT

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

    def rec(self):
        api_base = "https://vcms-api.hibiki-radio.jp/api/v1/"
        headers = {"X-Requested-With": "XMLHttpRequest"}
        # 番組の取得
        res = requests.get(api_base+"programs", headers=headers)
        prog_data = json.loads(res.text)
        returnData = []
        for program in prog_data:
            episode = program.get("episode")
            if episode is None:
                continue
            title = program.get("name")
            personality = program.get("cast")
            if (self.keyword.search(title) or self.keyword.search(personality)):
                # フォルダの作成
                dir_path = self.SAVEROOT + "/" + title.replace(" ", "_")
                f.createSaveDir(dir_path)
                # ファイル重複チェック
                update_date = DT.datetime.strptime(episode["updated_at"].split(" ")[0], "%Y/%m/%d")
                file_name = title.replace(" ", "_") + "_" + update_date.strftime("%Y%m%d") + ".m4a"
                file_path = dir_path +"/"+ file_name
                if f.did_record_prog(file_path, title, update_date.strftime("%Y%m%d")):
                    continue
                url2 = "https://vcms-api.hibiki-radio.jp/api/v1/programs/" + program.get("access_id")
                res2 = requests.get(url2, headers=headers)
                tmpjson = json.loads(res2.text)
                if (tmpjson.get("episode") is None) or (tmpjson["episode"].get("video") is None):
                    continue
                video_url = api_base + "videos/play_check?video_id=" + str(tmpjson["episode"]["video"]["id"])
                res2 = requests.get(video_url, headers=headers)
                tmpjson = json.loads(res2.text)
                print(title)
                if (tmpjson.get("playlist_url") is None):
                    continue
                returnData.append(title)
                cwd = 'ffmpeg -loglevel error -i "%s" -acodec copy "%s"' % (tmpjson["playlist_url"], file_path)
                subprocess.run(cwd, shell=True)

                # fs = open(file_path, "rb")
                # f.DropBox.upload(title, update_date.strftime("%Y%m%d"), fs.read())
                # fs.close()
                url = f.Swift.upload_file(filePath=file_path)
                f.Mysql.insert(
                    title= title,
                    pfm= personality,
                    timestamp= update_date.strftime("%Y%m%d"),
                    station= "hibiki",
                    uri= url,
                )
                if (f.Swift.hadInit):
                    cmd = 'rm "%s"' % (file_path)
                    subprocess.run(cmd, shell=True)
        print("finish")
        return returnData
