FROM golang:1.24-alpine AS builder

WORKDIR /usr/local/src/

RUN apk --no-cache add bash make gcc gettext musl-dev

COPY ["go.mod", "go.sum", "./"]

RUN go mod download

COPY . .
RUN make build


FROM alpine AS runner

COPY --from=builder /usr/local/src/build/houseservice .
COPY .env .

CMD [ "/houseservice" ]