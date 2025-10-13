FROM ubuntu:latest
WORKDIR /app
COPY todo-final .env ./
COPY web ./web/

ENV TODO_PORT=7540
ENV TODO_DBFILE="scheduler.db"
ENV TODO_PASSWORD="123456"

CMD ["./todo-final", "-logtostderr=true"]
