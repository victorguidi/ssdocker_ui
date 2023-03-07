#Build the frontend
FROM docker.io/node:18-alpine AS frontend

WORKDIR /app

RUN npm i -g pnpm

COPY ./frontend/package.json ./

RUN pnpm install

COPY ./frontend ./

RUN pnpm build

FROM docker.io/golang:1.20-alpine AS backend

WORKDIR /app

COPY ./backend/go.mod ./

RUN go mod download

COPY ./backend/ ./

COPY --from=frontend /app/static ./static

RUN go build -o main ./src

FROM alpine:latest

WORKDIR /app

COPY --from=backend /app/static ./static

COPY --from=backend /app/main ./

COPY --from=backend /app/src/config.yml ./

EXPOSE 8080

CMD ["./main", "--c", "config.yml"]

