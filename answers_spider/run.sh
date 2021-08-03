#!/bin/sh
Xvfb :99 -screen 0 1024x748x24+32 -nolisten tcp &
x11vnc -display $DISPLAY -bg -forever -nopw -quiet -listen 0.0.0.0 -xkb &
python3 /usr/local/bin/answers.py #| python3 /usr/local/bin/parser.py |  websocat -s 0.0.0.0:1234
#python3 /usr/local/bin/answers.py | python3 /usr/local/bin/parser.py | xargs -I{} sh -c '/usr/local/bin/send_message.sh {}'
