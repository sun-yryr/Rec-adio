# -*- coding: utf-8 -*-
import os
import json
import requests

line_token = ""

def load_configurations(Path):
    if (os.path.isfile(Path) is False):
        print("config file is not found")
        return None
    f = open(Path, "r")
    tmp = json.load(f)
    return tmp

def createSaveDir(Path):
    if (os.path.isdir(Path) is False):
        os.makedirs(Path)

def is_recording_succeeded(Path):
    size = os.path.getsize(Path)
    if (size >= 1024) and (62914560 > size):
        return True
    else:
        return False

def recording_successful_toline(title):
    headers = {"Authorization": "Bearer %s" % line_token}
    payload = {"message": "\n"+title+" を録音しました!"}
    requests.post("https://notify-api.line.me/api/notify", headers=headers, data=payload)

def recording_failure_toline(title):
    headers = {"Authorization": "Bearer %s" % line_token}
    payload = {"message": "\n"+title+" の録音に失敗しました"}
    requests.post("https://notify-api.line.me/api/notify", headers=headers, data=payload)