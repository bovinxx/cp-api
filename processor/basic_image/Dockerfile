FROM ubuntu:latest

RUN apt-get update && apt-get install -y python3
RUN ln -sf /usr/bin/python3 /usr/bin/python

RUN apt-get install -y bash

RUN apt-get install -y gcc
RUN ln -sf /usr/bin/gcc /usr/bin/cc

CMD ["bash"]