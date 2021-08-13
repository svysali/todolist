FROM golang:1.16-alpine3.14 AS build

WORKDIR /todolist

EXPOSE 8080

COPY . .

RUN go build -o ./dist/todolist *.go


FROM alpine:3.14

COPY --from=build /todolist/dist/todolist /app/todolist
COPY --from=build /todolist/db/migrations /app/db/migrations
COPY --from=build /todolist/build/entrypoint.sh /app/entrypoint.sh

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]


