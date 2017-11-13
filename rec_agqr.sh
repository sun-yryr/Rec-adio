!/bin/sh

DATE=`date +%Y_%m%d_%H%M`
SAVE_sDIR=`date +%Y%m`
SAVE_DIR="/home/pi/AGQR/${SAVE_sDIR}/"
RECORDING_TIME=5

RTMP_r='--rtmp "rtmpe://fms1.uniqueradio.jp/"'
RTMP_a='-a ?rtmp://fms-base2.mitene.ad.jp/agqr/'
RTMP_f='-f "WIN 16,0,0,257"'
RTMP_W='-W http://www.uniqueradio.jp/agplayerf/LIVEPlayer-HD0318.swf'
RTMP_p='-p http://www.uniqueradio.jp/agplayerf/newplayerf2-win.php'
RTMP_C='-C B:0'
RTMP_y='-y aandg22'
RTMP_B="--stop ${RECORDING_TIME}"
RTMP_o="-o ${SAVE_DIR}agqr-${DATE}.flv"
RTMP_PARAM="${RTMP_r} ${RTMP_a} ${RTMP_f} ${RTMP_W} ${RTMP_p} ${RTMP_C} ${RTMP_y} ${RTMP_B} --live ${RTMP_o}"

# echo agqr-${DATE}.flv
echo "rtmpdump ${RTMP_PARAM}"

#recoding start
rtmpdump ${RTMP_PARAM}