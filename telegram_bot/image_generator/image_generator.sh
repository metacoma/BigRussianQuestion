#!/bin/sh

FONT=./OldCyr_Bold.ttf
FONT=./Kalinov.otf
FONT=./cyrillic_old.ttf
#FONT=courier

# https://legacy.imagemagick.org/Usage/blur/

blur() {
  convert png:- -channel RGBA -blur 2x2 png:-
} 

fill_white() {
  convert png:- -fill white -colorize 80% png:-
} 

add_text() {
  convert png:- \( -background none -fill black -font $FONT -pointSize 60 -size 800x -gravity center caption:"$*" \) -gravity center -composite png:-
}

add_khokhoma() {
  convert png:- \( kho2.jpg -resize 200x200 -extent 200x200 xc:none  -fuzz 10% -transparent white \) -gravity NorthEast -geometry -40x+40 -composite -compose CopyOpacity -shave 1 png:-
} 

add_watermark() {
  convert png:- \( -background none -fill \#888 -font $FONT -pointSize 25 -gravity SouthWest -annotate +10+10 "t.me/BigRussianQuestion" \) png:-
}

any2png() {
  convert - -resize '1098x598^' png:-
}


any2png | fill_white | blur | add_text "$*" | add_watermark | add_khokhoma
