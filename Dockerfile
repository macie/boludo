FROM python:3-alpine

LABEL Name=boludo \
      Version=0.0.1
EXPOSE 3001

WORKDIR /src
ADD . /src

# python dependencies:
RUN python3 -m pip install -r requirements.txt
