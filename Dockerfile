#run backend build
FROM golang:1.22-alpine AS backend-build

WORKDIR /go-lang-api/toDo
COPY ./toDo/go.mod ./toDo/go.sum ./

RUN go mod download
COPY ./toDo ./

RUN go build -o toDo 

# run frontend build
FROM node:18-alpine AS frontend-build
WORKDIR /go-lang-api/todo-frontend
COPY ./todo-frontend/package.json ./todo-frontend/package-lock.json ./
RUN npm install
COPY ./todo-frontend ./
RUN npm run build

# lightweight image
FROM alpine:latest
WORKDIR /app
COPY --from=backend-build /go-lang-api/toDo .
COPY --from=frontend-build /go-lang-api/todo-frontend/build ./todo-frontend/build

EXPOSE 8090
CMD ["./toDo"]