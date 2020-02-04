FROM alpine as gobuild
RUN apk add go git

FROM gobuild as build
COPY . /code
WORKDIR /code
RUN go mod vendor && \
    go build -mod vendor

FROM alpine
MAINTAINER Kazım SARIKAYA <kazimsarikaya@sanaldiyar.com>
COPY --from=build /code/zkmetrics .
CMD /zkmetrics
