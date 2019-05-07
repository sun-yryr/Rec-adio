# -*- coding: utf-8 -*-
import requests
import json
from bs4 import BeautifulSoup

def main():
    url = "https://www.agqr.jp/timetable/streaming-sp.php"
    res = requests.get(url)
    res.encoding = 'utf-8'
    soup = BeautifulSoup(res.text, 'lxml')
    root = soup.find_all("table", attrs={"class": "timetb-ag"})
    for tr in root.find_all("tr"):
        for td in tr.find_all("td"):
            print()

if __name__ == "__main__":
    main()