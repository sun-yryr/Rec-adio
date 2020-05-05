FROM python:3.6.10-slim-buster

ENV DEBIAN_FRONTEND noninteractive
ENV DEBCONF_NOWARNINGS yes

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        build-essential \
        ca-certificates \
        git \
        ffmpeg \
        rtmpdump && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    pip install pipenv

RUN git clone https://github.com/sun-yryr/Rec-adio.git /app

WORKDIR /app

RUN pipenv install

CMD ["pipenv", "run", "start"]
