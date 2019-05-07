# -*- coding: utf-8 -*-
import radiko
import subprocess
from multiprocessing import Process
from signal import SIGINT
import time
import datetime as DT


def main():
    radio = radiko.yradio()
    AuthToken = radio.authorization()
    radio.change_keywords(["声優", "垣花正あなたとハッピー"])
    res = radio.search()
    com = DT.timedelta(seconds=20)
    zero = DT.timedelta()
    while(True):
        now = DT.datetime.now()
        print(res)
        for data in res:
            tmp_time = data["DT_ft"] - now
            if (tmp_time < zero):
                res.remove(data)
            elif (tmp_time < com):
                print(data["ftl"])
                p = Process(target=rec, args=([data, tmp_time.total_seconds(), AuthToken],))
                p.start()
                res.remove(data)
        time.sleep(10)



def rec(data):
    program_data = data[0]
    AuthToken = data[2]
    print(program_data["title"])
    cwd = ('ffmpeg -headers "X-Radiko-AuthToken: %s" -i "https://rpaa.smartstream.ne.jp/so/playlist.m3u8?station_id=%s&l=15&lsid=85605612400545823313125539756301319021&type=b" -acodec copy ./savefile/%s.m4a' % (AuthToken, program_data["station"], program_data["title"]+"_"+program_data["ft"][:12]))
    p1 = subprocess.Popen(cwd, stdout=subprocess.DEVNULL, shell=True)
    time.sleep(program_data["dur"]+data[1]-5)
    #time.sleep(30)
    p1.send_signal(SIGINT)

if __name__ == "__main__":
    main()