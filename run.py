#!/bin/python3.6
# -*- coding: utf-8 -*-
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
    Radiko.change_keywords(config["radiko"]["keywords"])
    radiko_data = Radiko.search()
    while(True):
        now = DT.datetime.now()
        #print(radiko_data)
        if (bool(radiko_data)):
            for data in radiko_data:
                tmp_time = data["DT_ft"] - now
                if (tmp_time < T_ZERO):
                    radiko_data.remove(data)
                elif (tmp_time < T_BASELINE):
                    AuthToken = Radiko.authorization()
                    p = Process(target=radiko.rec, args=([data, tmp_time.total_seconds(), AuthToken, SAVEROOT, dbx],))
                    p.start()
                    radiko_data.remove(data)
        if (now.hour == 6 and now.minute <= 5 and Radiko.reload_date != DT.date.today()):
            Radiko.reload_program()
            radiko_data = Radiko.search()
        time.sleep(50)

def main_agqr():
    agqr_data = config["AnG"]
    while(True):
        now = DT.datetime.now()
        #print(agqr_data)
        for data in agqr_data:
            if (data["isRepeat"]):
                # 毎週録音
                # 録音は今日?
                if (now.weekday() == data["weekday"]):
                    # 録音開始時間の生成
                    tmp = DT.datetime.strptime(now.strftime("%Y%m%d")+data["start_ts"], "%Y%m%d%H%M%S")
                    # 現在時刻との差分
                    tmp_time = tmp - now
                    # 過去でなく、数分前か
                    if (tmp_time > T_ZERO and tmp_time < T_BASELINE):
                        p = Process(target=agqr.rec, args=([data, tmp_time.total_seconds(), tmp.strftime("%Y%m%d%H%M%S"), SAVEROOT, dbx],))
                        p.start()
            else:
                tmp = DT.datetime.strptime(data["start_ts"], "%Y%m%d%H%M%S")
                tmp_time = tmp - now
                if (tmp_time > T_ZERO and tmp_time < T_BASELINE):
                    p = Process(target=agqr.rec, args=([data, tmp_time.total_seconds(), tmp.strftime("%Y%m%d%H%M%S"), SAVEROOT, dbx],))
                    p.start()
        time.sleep(50)

def main_onsen_hibiki():
    Onsen = onsen.onsen(config["Onsen"]["keywords"], SAVEROOT, dbx)
    Hibiki = hibiki.hibiki(config["Hibiki"]["keywords"], SAVEROOT, dbx)
    while(True):
        now = DT.datetime.now()
        if (now.hour == 7 and now.minute <= 5 and Onsen.reload_date != DT.date.today()):
            titles = Onsen.rec()
            titles.extend(Hibiki.rec())
            if (bool(titles)):
                f.recording_successful_toline("、".join(titles))
            else:
                print("in onsen, hibiki. there aren't new title.")
        time.sleep(300)

if __name__ == "__main__":
    ROOT = (__file__.replace("/run.py", ""))
    if (ROOT == __file__):
        ROOT = ROOT.replace("run.py", ".")
    config = f.load_configurations(ROOT + "/conf/config.json")
    if (config is None):
        exit(code=-1)
    f.line_token = config["all"]["line_token"]
    SAVEROOT = config.get("all").get("savedir")
    if (os.path.isfile(SAVEROOT) is None):
        SAVEROOT = ROOT + "/savefile"
    dbx = dropbox.Dropbox(config["all"]["dbx_token"])
    dbx.users_get_current_account()
    res = dbx.files_list_folder('')
    db_list = [d.name for d in res.entries]
    if not "radio" in db_list:
        dbx.files_create_folder("radio")
    ps = [
        Process(target=main_radiko),
        Process(target=main_agqr),
        Process(target=main_onsen_hibiki)
    ]
    for i in ps:
        i.start()
    while(True):
        time.sleep(1000)
