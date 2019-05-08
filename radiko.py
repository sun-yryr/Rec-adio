# -*- coding: utf-8 -*-
import requests
import json
import xml.etree.ElementTree as ET
import re
import base64
import datetime as DT


cast = re.compile("(三森すずこ|夏川椎菜)")

class yradio:
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
        print(res.headers)
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
            print("----------")
            print(res.headers)
            return AuthToken
        else:
            print("Authorization2 Failed")
            return None
