#!/usr/bin/python3
# -*- coding: utf-8 -*-
# テストするときは，Rec-adio直下に置いてテストしてください
from lib import radiko, agqr, onsen, hibiki, functions as f
import subprocess
from multiprocessing import Process
from signal import SIGINT
import time
import datetime as DT
import json
import requests
import re
import os
import dropbox

SAVEROOT = ""
T_BASELINE = DT.timedelta(seconds=60)
T_ZERO = DT.timedelta()


def main_radiko():
    Radiko = radiko.radiko()
    AuthToken = Radiko.authorization()
    program = {
        "title": "test_rec",
        "ft": "20190829181110",
        "dur": 60,
        "station": "QRR"
    }
    p = Process(target=radiko.rec, args=([program, 0, AuthToken, SAVEROOT, dbx],))
    p.start()
    t = 0
    while(True):
        print(t)
        t += 10
        time.sleep(10)

if __name__ == "__main__":
    ROOT = (__file__.replace("/radiko_test.py", ""))
    if (ROOT == __file__):
        ROOT = ROOT.replace("radiko_test.py", ".")
    config = f.load_configurations(ROOT + "/conf/config.json")
    if (config is None):
        exit(code=-1)
    f.line_token = config["all"]["line_token"]
    SAVEROOT = config.get("all").get("savedir")
    if (os.path.isfile(SAVEROOT) is None) or (SAVEROOT == ""):
        SAVEROOT = ROOT + "/savefile"
    print("SAVEROOT : " + SAVEROOT)
    dbx = dropbox.Dropbox(config["all"]["dbx_token"])
    dbx.users_get_current_account()
    res = dbx.files_list_folder('')
    db_list = [d.name for d in res.entries]
    if not "radio" in db_list:
        dbx.files_create_folder("radio")
    ps = [
        Process(target=main_radiko)
    ]
    for i in ps:
        i.start()
    while(True):
        time.sleep(1000)
