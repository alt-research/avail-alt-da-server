FROM golang:1.23.4
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make da-server
ARG ADDR=0.0.0.0
ARG PORT=8080
ARG AVAIL_RPC=http://localhost:9933
ARG AVAIL_SEED=""
ARG AVAIL_APPID=0
ENV ADDR=${ADDR} \
    PORT=${PORT} \
    AVAIL_RPC=${AVAIL_RPC} \
    AVAIL_SEED=${AVAIL_SEED} \
    AVAIL_APPID=${AVAIL_APPID}
EXPOSE ${PORT}
EXPOSE 8080
EXPOSE 433

# Set default values for environment variables
ENV ADDR=0.0.0.0
ENV PORT=8080
ENV AVAIL_RPC=http://localhost:9933
ENV AVAIL_SEED=""
ENV AVAIL_APPID=0
ENV AVAIL_TIMEOUT=100

EXPOSE ${PORT}

# Print environment variables and run the application
CMD ./bin/avail-da-server --addr=$ADDR --port=$PORT --avail.rpc="$AVAIL_RPC" --avail.seed="$AVAIL_SEED" --avail.appid=$AVAIL_APPID --avail.timeout=$AVAIL_TIMEOUT
