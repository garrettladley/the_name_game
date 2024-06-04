FROM golang:1.22-alpine as builder

WORKDIR /app
RUN apk add --no-cache make nodejs npm git

COPY . ./
RUN make install
RUN make build-prod

FROM scratch
COPY --from=builder /app/bin/the_name_game /the_name_game

EXPOSE 3000
ENTRYPOINT [ "./the_name_game" ]