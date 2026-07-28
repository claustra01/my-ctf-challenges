FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir notes
RUN chmod 755 notes

RUN go build -o main .


FROM gcr.io/distroless/base:nonroot AS runner
ENV GOENV=production

COPY --from=build --chown=nonroot:nonroot /app/main /app/main
COPY --from=build --chown=nonroot:nonroot /app/index.html /app/index.html
COPY --from=build --chown=nonroot:nonroot /app/notes /app/notes

EXPOSE 8080

USER nonroot
ENTRYPOINT ["/app/main"]