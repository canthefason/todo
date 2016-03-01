FROM centurylink/ca-certs
MAINTAINER Can Yucel "can.yucel@gmail.com"
EXPOSE 8000

ADD todo-service /

CMD ["/todo-service"]
