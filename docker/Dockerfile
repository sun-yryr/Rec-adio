FROM python:3.9-slim-buster

WORKDIR /app

RUN apt update \
    && apt install -y ffmpeg rtmpdump

COPY . .

RUN pip install pipenv
RUN pipenv install
