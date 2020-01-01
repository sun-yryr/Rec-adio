# -*- coding: utf-8 -*-
import os
import json
import requests
import dropbox

def load_configurations():
    ROOT = (__file__.replace("/lib/functions.py", ""))
    if (ROOT == __file__):
        ROOT = ROOT.replace("functions.py", ".")
    Path = ROOT + "/conf/config.json"
    if (os.path.isfile(Path) is False):
        print("config file is not found")
        return None
    f = open(Path, "r")
    tmp = json.load(f)
    return tmp

def createSaveDirPath(path = ""):
    ROOT = (__file__.replace("/lib/functions.py", ""))
    if (ROOT == __file__):
        ROOT = ROOT.replace("functions.py", ".")
    if path == "":
        Path = ROOT + "/savefile"
    else:
        Path = path
    if not os.path.isdir(Path):
        os.makedirs(Path)
    return Path

def createSaveDir(Path):
    if (os.path.isdir(Path) is False):
        os.makedirs(Path)

def is_recording_succeeded(Path):
    m4a_path = Path + ".m4a"
    if (os.path.isfile(m4a_path)):
        size = os.path.getsize(m4a_path)
        print("rec size = " + str(size))
        if (size >= 1024):
            return True
        else:
            return False
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

class DBXController():
    def __init__(self):
        tmpconf = load_configurations()
        if (tmpconf is None) or (tmpconf["all"]["dbx_token"] == ""):
            self.hadInit = False
            return
        self.dbx = dropbox.Dropbox(tmpconf["all"]["dbx_token"])
        self.dbx.users_get_current_account()
        res = self.dbx.files_list_folder('')
        db_list = [d.name for d in res.entries]
        if not "radio" in db_list:
            self.dbx.files_create_folder("radio")

    def upload(self, title, ft, fileData):
        if not self.hadInit:
            return
        dbx_path = "/radio/" + title
        # dropboxにフォルダを作成する
        res = self.dbx.files_list_folder('/radio')
        db_list = [d.name for d in res.entries]
        if not title in db_list:
            self.dbx.files_create_folder(dbx_path)
        dbx_path += "/" +title + "_" + ft[:12]+ ".m4a"
        self.dbx.files_upload(fileData, dbx_path)
    
    def upload_onsen(self, title, count, fileData):
        if not self.hadInit:
            return
        dbx_path = "/radio/" + title
        # dropboxにフォルダを作成する
        res = self.dbx.files_list_folder('/radio')
        db_list = [d.name for d in res.entries]
        if not title in db_list:
            self.dbx.files_create_folder(dbx_path)
        dbx_path += "/" + title + "#" + count + ".mp3"
        self.dbx.files_upload(fileData, dbx_path)

DropBox = DBXController()