#!/bin/bash

DATE=`date "+%Y%m%d"`
SAVE_DIR="/home/pi/hdd1/sun-mba/radio/"
RECTIME=`expr $1 \* 60`

RTMP_r='"rtmpe://fms1.uniqueradio.jp/"'
RTMP_a='?rtmp://fms-base1.mitene.ad.jp/agqr/'
RTMP_f='"WIN 16,0,0,257"'
RTMP_W='http://www.uniqueradio.jp/agplayerf/LIVEPlayer-HD0318.swf'
RTMP_p='http://www.uniqueradio.jp/agplayerf/newplayerf2-win.php'
RTMP_C='B:0'
RTMP_y='aandg22'

RTMP_o="${SAVE_DIR}$2_${DATE}"

RTMP_PARAM="--rtmp ${RTMP_r} -a ${RTMP_a} -f ${RTMP_f} -W ${RTMP_W} -p ${RTMP_p} -C ${RTMP_C} -y ${RTMP_y} --stop ${RECTIME} --live -o ${RTMP_o}.flv"

#echo ${RTMP_PARAM}

rtmpdump ${RTMP_PARAM}

avconv -loglevel quiet -i ${RTMP_o}.flv -acodec libmp3lame -ab 64k ${RTMP_o}.mp3

rm ${RTMP_o}.flv

URL="https://notify-api.line.me/api/notify"
MSG="超A&G URL変更"
SIZE=`wc -c ${RTMP_o}.mp3 |cut -d ' ' -f1 -`
if [ $SIZE -lt 500 ]; then
    curl -X POST -H 'Authorization: Bearer AmCprNqcCYuboLMCqvNRHD3y78lNUnCkKpJkvgJsYO8' -F "message=$MSG"  $URL
fi