# -*- coding: utf-8 -*-
import requests
import json
import xml.etree.ElementTree as ET
import re
import base64
import datetime as DT
import lib.functions as f
import time
import subprocess
from signal import SIGINT

class radiko:
    RADIKO_URL = "http://radiko.jp/v3/program/today/JP13.xml"
    def __init__(self):
        res = requests.get(self.RADIKO_URL)
        res.encoding = "utf-8"
        self.isKeyword = False
        self.reload_date = DT.date.today()
        self.program_radiko = ET.fromstring(res.text)

    def reload_program(self):
        res = requests.get(self.RADIKO_URL)
        res.encoding = "utf-8"
        self.program_radiko = ET.fromstring(res.text)
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
        for station in self.program_radiko.findall("stations/station"):
            for prog in station.findall("./progs/prog"):
                ck = False
                title = prog.find("title").text
                info = prog.find("info").text
                desc = prog.find("desc").text
                pfm = prog.find("pfm").text
                if (self.keyword.search(title)):    ck = True
                if (ck is False) and (info is not None):
                    if (self.keyword.search(info)): ck = True
                if (ck is False) and (pfm is not None):
                    if (self.keyword.search(pfm)):  ck = True
                if (ck is False) and (desc is not None):
                    if (self.keyword.search(desc)): ck = True
                if (ck):
                    res.append({
                        "station": station.get("id"),
                        "title": prog.find("title").text.replace(" ", "_"),
                        "ft": prog.get("ft"),
                        "DT_ft": DT.datetime.strptime(prog.get("ft"), "%Y%m%d%H%M%S"),
                        "to": prog.get("to"),
                        "ftl": prog.get("ftl"),
                        "tol": prog.get("tol"),
                        "dur": int(prog.get("dur"))
                    })
        if bool(res): return res
        else: return []

    def authorization(self):
        auth1_url = "https://radiko.jp/v2/api/auth1"
        auth2_url = "https://radiko.jp/v2/api/auth2"
        auth_key = "bcd151073c03b352e1ef2fd66c32209da9ca0afa"
        headers = {
            "X-Radiko-App": "pc_html5",
            "X-Radiko-App-Version": "0.0.1",
            "X-Radiko-User": "sunyryr",
            "X-Radiko-Device": "pc"
        }
        res = requests.get(auth1_url, headers=headers)
        if (res.status_code != 200):
            print("Authorization1 Failed")
            return None
        #print(res.headers)
        AuthToken = res.headers["X-RADIKO-AUTHTOKEN"]
        KeyLength = int(res.headers["X-Radiko-KeyLength"])
        KeyOffset = int(res.headers["X-Radiko-KeyOffset"])
        tmp_authkey = auth_key[KeyOffset:KeyOffset+KeyLength]
        AuthKey = base64.b64encode(tmp_authkey.encode('utf-8')).decode('utf-8')
        #print(AuthKey)
        headers = {
            "X-Radiko-AuthToken": AuthToken,
            "X-Radiko-PartialKey": AuthKey,
            "X-Radiko-User": "sunyryr",
            "X-Radiko-Device": "pc"
        }
        res = requests.get(auth2_url, headers=headers)
        if (res.status_code == 200):
            #print("----------")
            #print(res.headers)
            return AuthToken
        else:
            print("Authorization2 Failed")
            return None

def rec(data):
    program_data = data[0]
    wait_start_time = data[1]
    AuthToken = data[2]
    SAVEROOT = data[3]
    dbx = data[4]
    #ディレクトリの作成
    dir_path = SAVEROOT + "/" + program_data["title"]
    f.createSaveDir(dir_path)
    dbx_path = "/radio/" + program_data["title"]
    res = dbx.files_list_folder('/radio')
    db_list = [d.name for d in res.entries]
    if not program_data["title"] in db_list:
        dbx.files_create_folder(dbx_path)
    #保存先パスの作成
    file_path = dir_path + "/" + program_data["title"]+"_"+program_data["ft"][:12]
    dbx_path += "/" +program_data["title"]+"_"+program_data["ft"][:12]+ ".m4a"
    #print(program_data["title"])
    #stream urlの取得
    url = 'http://f-radiko.smartstream.ne.jp/%s/_definst_/simul-stream.stream/playlist.m3u8' % program_data["station"]
    m3u8 = gen_temp_chunk_m3u8_url(url, AuthToken)
    #コマンドの実行
    time.sleep(wait_start_time)
    cwd = ('ffmpeg -headers "X-Radiko-AuthToken: %s" -i "%s" -acodec copy "%s.m4a"' % (AuthToken, m3u8, file_path))
    p1 = subprocess.Popen(cwd, stdout=subprocess.DEVNULL, shell=True)
    time.sleep(program_data["dur"]-10)
    p1.send_signal(SIGINT)
    time.sleep(5)
    if (f.is_recording_succeeded(file_path+ ".m4a")):
        f.recording_successful_toline(program_data["title"])
        fs = open(file_path+".m4a", "rb")
        dbx.files_upload(fs.read(), dbx_path)
        fs.close()
    else:
        f.recording_failure_toline(program_data["title"])

def gen_temp_chunk_m3u8_url( url, AuthToken ):
    headers =  {
        "X-Radiko-AuthToken": AuthToken,
    }
    res  = requests.get(url, headers=headers)
    res.encoding = "utf-8"
    body = res.text
    lines = re.findall( '^https?://.+m3u8$' , body, flags=(re.MULTILINE) )
    return lines[0]