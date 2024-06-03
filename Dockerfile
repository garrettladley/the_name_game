FROM golang:1.22-alpine as builder

WORKDIR /app
RUN apk add --no-cache make nodejs npm git

COPY . ./
RUN make install
RUN make build-prod
RUN > /app/.env

FROM scratch
COPY --from=builder /app/bin/the_name_game /the_name_game
COPY --from=builder /app/public /public
COPY --from=builder /app/htmx /htmx
COPY --from=builder /app/.env .env

EXPOSE 3000
ENTRYPOINT [ "./the_name_game" ]