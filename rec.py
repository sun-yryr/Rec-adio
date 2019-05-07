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
                p = Process(target=rec_radiko, args=([data, tmp_time.total_seconds(), AuthToken],))
                p.start()
                res.remove(data)
        time.sleep(10)



def rec_radiko(data):
    program_data = data[0]
    AuthToken = data[2]
    print(program_data["title"])
    cwd = ('ffmpeg -headers "X-Radiko-AuthToken: %s" -i "https://rpaa.smartstream.ne.jp/so/playlist.m3u8?station_id=%s&l=15&lsid=85605612400545823313125539756301319021&type=b" -acodec copy ./savefile/%s.m4a' % (AuthToken, program_data["station"], program_data["title"]+"_"+program_data["ft"][:12]))
    p1 = subprocess.Popen(cwd, stdout=subprocess.DEVNULL, shell=True)
    time.sleep(program_data["dur"]+data[1]-5)
    #time.sleep(30)
    p1.send_signal(SIGINT)

def rec_agqr(data):
    program_data = data[0]
    filepath = program_data["title"] + "_" + program_data["ft"][:12]
    cwd = ('rtmpdump --rtmp "rtmpe://fms1.uniqueradio.jp/" ')
    cwd += ('-a ?rtmp://fms-base1.mitene.ad.jp/agqr/ ')
    cwd += ('-f "WIN 16,0,0,257" ')
    cwd += ('-W http://www.uniqueradio.jp/agplayerf/LIVEPlayer-HD0318.swf ')
    cwd += ('-p http://www.uniqueradio.jp/agplayerf/newplayerf2-win.php ')
    cwd += ('-C B:0 ')
    cwd += ('-y aandg22 ')
    cwd += ('--stop %s ' % program_data["dur"])
    cwd += ('--live -o %s.flv' % (filepath))
    subprocess.run(cwd, shell=True)
    
    cwd2 = ('avconv -loglevel quiet -i %s.flv -acodec libmp3lame -ab 64k %s.mp3' % (filepath, filepath))
    subprocess.run(cwd2, shell=True)


if __name__ == "__main__":
    main()