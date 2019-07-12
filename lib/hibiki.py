# -*- coding: utf-8 -*-
import requests
import json
import re
import os
import lib.functions as f
import datetime as DT
import subprocess

class hibiki:
    def __init__(self, keywords, SAVEROOT, dbx):
        self.change_keywords(keywords)
        self.SAVEROOT = SAVEROOT
        self.dbx = dbx

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
                title = title.replace(" ", "_")
                # フォルダの作成
                dir_path = self.SAVEROOT + "/" + title
                f.createSaveDir(dir_path)
                # ファイル重複チェック
                update_date = DT.datetime.strptime(episode["updated_at"].split(" ")[0], "%Y/%m/%d")
                file_name = title + "_" + update_date.strftime("%Y%m%d") + ".m4a"
                file_path = dir_path +"/"+ file_name
                file_path = file_path.replace(" ", "_")
                if file_name in os.listdir(dir_path):
                    continue
                url2 = "https://vcms-api.hibiki-radio.jp/api/v1/programs/" + program.get("access_id")
                res2 = requests.get(url2, headers=headers)
                tmpjson = json.loads(res2.text)
                video_url = api_base + "videos/play_check?video_id=" + str(tmpjson["episode"]["video"]["id"])
                res2 = requests.get(video_url, headers=headers)
                tmpjson = json.loads(res2.text)
                print(tmpjson)
                print(title)
                if (tmpjson.get("playlist_url") is None):
                    continue
                returnData.append(title)
                cwd = 'ffmpeg -i "%s" -acodec copy "%s"' % (tmpjson["playlist_url"], file_path)
                p1 = subprocess.run(cwd, stdout=subprocess.DEVNULL, shell=True)
                dbx_path = "/radio/" + title
                res = self.dbx.files_list_folder('/radio')
                db_list = [d.name for d in res.entries]
                if not title in db_list:
                    self.dbx.files_create_folder(dbx_path)
                dbx_path += "/" + file_name
                fs = open(file_path, "rb")
                self.dbx.files_upload(fs.read(), dbx_path)
                fs.close()
        print("finish")
        return returnData