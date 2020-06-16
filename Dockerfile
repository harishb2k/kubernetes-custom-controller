FROM golang:1.13.12-buster
RUN mkdir /code
RUN mkdir /code/controller
WORKDIR /code
COPY controller/controller.go /code/controller
COPY * /code/
CMD ["/bin/bash", "run.sh"]