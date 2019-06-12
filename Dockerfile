FROM golang:1.12-alpine AS build
LABEL maintainer="Vitalii Kozlovskyi <ubombi@gmail.com>"

# to update dependencies run `docker build --no-cache .`
RUN apk add --no-cache git ca-certificates

WORKDIR /src/
ENV GOBIN=/bin


# Only re-built when dependencies changed and gomod files were updated
COPY go.mod go.sum ./
RUN go mod download

# Code-generation should be placed here. In separate layer

# This layer is rebuilt when a file changes in the project directory
COPY . ./
RUN go build -ldflags "-extldflags '-static'" -tags netgo -o /bin/server


### Main stage
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /bin/server /bin/server
EXPOSE 8008
ENTRYPOINT ["server"]
