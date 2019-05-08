# -*- coding: utf-8 -*-
import radiko
import subprocess
from multiprocessing import Process
from signal import SIGINT
import time
import datetime as DT
import json

with open("./settings.json") as f:
    aa = json.load(f)
    SAVEDIR = aa["radio"]["savedir"]


def main():
    radio = radiko.yradio()
    AuthToken = radio.authorization()
    with open("./settings.json") as f:
        jd = json.load(f)
        radio.change_keywords(jd["radio"]["keywords"])
    radiko_data = radio.search()
    com = DT.timedelta(seconds=20)
    zero = DT.timedelta()
    while(True):
        now = DT.datetime.now()
        print(radiko_data)
        if (bool(radiko_data)):
            for data in radiko_data:
                tmp_time = data["DT_ft"] - now
                #if (tmp_time < zero):
                    #radiko_data.remove(data)
                if (tmp_time < com):
                    #print(data["ftl"])
                    p = Process(target=rec_radiko, args=([data, tmp_time.total_seconds(), AuthToken],))
                    p.start()
                    radiko_data.remove(data)
        if (now.hour == 6 and now.minute <= 2 and radio.reload_date != DT.date.today()):
            radio.reload_program()
        time.sleep(10)


def agqr_main():
    agqr_data = create_agqr_data()
    com = DT.timedelta(seconds=20)
    zero = DT.timedelta()
    while(True):
        now = DT.datetime.now()
        print(agqr_data)
        for data in agqr_data:
            if (data["isRepeat"]):
                if (now.weekday() == data["weekday"]):
                    tmp = DT.datetime.strptime(now.strftime("%Y%m%d")+data["start_ts"], "%Y%m%d%H%M%S")
                    tmp_time = tmp - now
                    if (tmp_time > zero and tmp_time < com):
                        p = Process(target=rec_agqr, args=([data, tmp_time.total_seconds(), tmp.strftime("%Y%m%d%H%M%S")],))
                        p.start()
            else:
                tmp = DT.datetime.strptime(data["start_ts"], "%Y%m%d%H%M%S")
                tmp_time = tmp - now
                if (tmp_time > zero and tmp_time < com):
                    p = Process(target=rec_agqr, args=([data, tmp_time.total_seconds(), tmp.strftime("%Y%m%d%H%M%S")],))
                    p.start()
        time.sleep(10)

def create_agqr_data():
    with open("./settings.json") as f:
        jd = json.load(f)
        return jd["AnG"]



def rec_radiko(data):
    program_data = data[0]
    AuthToken = data[2]
    filepath = program_data["title"]+"_"+program_data["ft"][:12]
    print(program_data["title"])
    cwd = ('ffmpeg -headers "X-Radiko-AuthToken: %s" -i "https://rpaa.smartstream.ne.jp/so/playlist.m3u8?station_id=%s&l=15&lsid=85605612400545823313125539756301319021&type=b" -acodec copy "./%s.m4a"' % (AuthToken, program_data["station"], filepath))
    p1 = subprocess.Popen(cwd, stdout=subprocess.DEVNULL, shell=True)
    time.sleep(program_data["dur"]+data[1]-5)
    #time.sleep(30)
    p1.send_signal(SIGINT)

def rec_agqr(data):
    program_data = data[0]
    filepath = program_data["title"] + "_" + data[2][:12]
    cwd = ('rtmpdump --rtmp "rtmpe://fms1.uniqueradio.jp/" ')
    cwd += ('-a ?rtmp://fms-base1.mitene.ad.jp/agqr/ ')
    cwd += ('-f "WIN 16,0,0,257" ')
    cwd += ('-W http://www.uniqueradio.jp/agplayerf/LIVEPlayer-HD0318.swf ')
    cwd += ('-p http://www.uniqueradio.jp/agplayerf/newplayerf2-win.php ')
    cwd += ('-C B:0 ')
    cwd += ('-y aandg22 ')
    cwd += ('--stop %s ' % str(program_data["dur"]+data[1]))
    cwd += ('--live -o %s.flv' % (filepath))
    subprocess.run(cwd, shell=True)
    
    cwd2 = ('avconv -loglevel quiet -i %s.flv -acodec libmp3lame -ab 64k %s.mp3' % (filepath, filepath))
    subprocess.run(cwd2, shell=True)


if __name__ == "__main__":
    ps = [
        Process(target=main),
        Process(target=agqr_main)
    ]
    for i in ps:
        i.start()
    while(True):
        time.sleep(1000)
