FROM swift:5.6-focal

WORKDIR /app

RUN apt update \
    && apt install -y ffmpeg

COPY application/Package.* .
RUN swift package resolve

COPY application/ .
RUN swift build -c release
RUN cp "$(swift build --package-path /app -c release --show-bin-path)/Run" ./

ENTRYPOINT ["./Run"]
CMD ["serve", "--env", "production", "--hostname", "0.0.0.0", "--port", "8080"]

