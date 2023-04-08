#!/bin/bash


filename1='singaporeAirlines.pdf'
filename2='scoot.pdf'
filename3='jetstar.pdf'
url='http://127.0.0.1:8000/read'

curl "$url" \
  --form "data=@$filename1" \
  --form "data=@$filename2" \
  -H "Content-Type: multipart/form-data"