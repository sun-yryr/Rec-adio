#!/usr/bin/python3
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

SAVEROOT = ""
T_BASELINE = DT.timedelta(seconds=60)
T_ZERO = DT.timedelta()


def main_radiko():
    Radiko = radiko.radiko()
    Radiko.change_keywords(keywords)
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
                    p = Process(target=radiko.rec, args=([data, tmp_time.total_seconds(), AuthToken, SAVEROOT],))
                    p.start()
                    radiko_data.remove(data)
        if (now.hour == 6 and now.minute <= 5 and Radiko.reload_date != DT.date.today()):
            Radiko.reload_program()
            radiko_data = Radiko.search()
        time.sleep(60)

def main_agqr():
    Agqr = agqr.agqr()
    Agqr.change_keywords(keywords)
    agqr_data = Agqr.search()
    #print(agqr_data)
    while(True):
        now = DT.datetime.now()
        if (bool(agqr_data)):
            for data in agqr_data:
                tmp_time = data["DT_ft"] - now
                if (tmp_time < T_ZERO):
                    agqr_data.remove(data)
                elif (tmp_time < T_BASELINE):
                    p = Process(target=Agqr.rec, args=([data, tmp_time.total_seconds(), SAVEROOT],))
                    p.start()
                    agqr_data.remove(data)
        if (now.hour == 0 and now.minute <= 5 and Agqr.reload_date != DT.date.today()):
            Agqr.reload_program()
            agqr_data = Agqr.search()
        time.sleep(60)


def main_onsen_hibiki():
    Onsen = onsen.onsen(keywords, SAVEROOT)
    Hibiki = hibiki.hibiki(keywords, SAVEROOT)
    while(True):
        now = DT.datetime.now()
        if (now.hour == 7 and now.minute <= 5 and Onsen.reload_date != DT.date.today()):
            titles = Onsen.rec()
            titles.extend(Hibiki.rec())
            if (bool(titles)):
                f.LINE.recording_successful_toline("、".join(titles))
            else:
                print("in onsen, hibiki. there aren't new title.")
        time.sleep(300)

if __name__ == "__main__":
    config = f.load_configurations()
    SAVEROOT = f.createSaveDirPath(config["all"]["savedir"])
    print("SAVEROOT : " + SAVEROOT)
    keywords = config["all"]["keywords"]
    ps = [
        Process(target=main_radiko),
        Process(target=main_agqr),
        Process(target=main_onsen_hibiki)
    ]
    for i in ps:
        i.start()
    while(True):
        time.sleep(1000)
