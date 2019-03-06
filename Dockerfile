FROM archlinux/base

LABEL Name=boludo \
      Version=0.0.1
EXPOSE 3001

# configure OS
ENV PATH /usr/local/bin:$PATH
ENV LANG C.UTF-8
RUN pacman -Sy

# prepare user
RUN useradd -m boludo
WORKDIR /home/boludo

# install OS packages
RUN pacman -S --noconfirm swi-prolog
RUN pacman -S --noconfirm python python-pip

# install python dependencies
COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt

#COPY ./src src

RUN chown -R boludo:boludo ./
USER boludo

CMD [ "/bin/bash" ]
