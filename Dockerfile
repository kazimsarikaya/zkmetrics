FROM alpine as build
MAINTAINER Kazım SARIKAYA <kazimsarikaya@sanaldiyar.com>
COPY . /code
WORKDIR /code
RUN apk add go git && \
    go mod vendor && \
    go build -mod vendor && \
    ls -ltrh

FROM alpine
COPY --from=build /code/zkmetrics .
CMD /zkmetrics
