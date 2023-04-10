#!/bin/bash


filename1='singaporeAirlines.pdf'
filename2='scoot.pdf'
filename3='jetstar.pdf'
url='http://127.0.0.1:8001/read'

curl "$url" \
  --form "data=@$filename1" \
  -H "Content-Type: multipart/form-data"