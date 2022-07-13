FROM swift:5.6-focal

WORKDIR /app

RUN apt update \
    && apt install -y ffmpeg curl

RUN curl -fsSL https://deb.nodesource.com/setup_16.x | bash - \
    && apt install -y nodejs \
    && npm i -g space-dl

RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp \
    && chmod a+rx /usr/local/bin/yt-dlp

COPY application/Package.* .
RUN swift package resolve

COPY application/ .
RUN swift build -c release
RUN cp "$(swift build --package-path /app -c release --show-bin-path)/Run" ./

ENTRYPOINT ["./Run"]
CMD ["serve", "--env", "production", "--hostname", "0.0.0.0", "--port", "8080"]
