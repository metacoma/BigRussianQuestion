FROM ubuntu:latest
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y python3-pip xvfb libxi6 libgconf-2-4 firefox curl x11vnc

ADD https://github.com/mozilla/geckodriver/releases/download/v0.29.1/geckodriver-v0.29.1-linux64.tar.gz /tmp
WORKDIR /tmp
RUN tar xzf geckodriver-v0.29.1-linux64.tar.gz && mv geckodriver /usr/bin/geckodriver 
RUN pip3 install selenium

ADD https://github.com/vi/websocat/releases/download/v1.8.0/websocat_1.8.0_newer_amd64.deb /tmp
WORKDIR /tmp
RUN dpkg -i ./websocat_1.8.0_newer_amd64.deb

ENV DISPLAY=:99
ENV DBUS_SESSION_BUS_ADDRESS=/dev/null

ADD ./run.sh /usr/local/bin/run.sh
ADD ./answers.py /usr/local/bin/answers.py
ADD ./parser.py /usr/local/bin/parser.py
ADD ./send_message.sh /usr/local/bin/send_message.sh

ENTRYPOINT [ "/usr/local/bin/run.sh" ]
