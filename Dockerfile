FROM golang:1.10

# Env
ENV DATAMESH_HOME=/go/src/github.com/nicolaferraro/datamesh
ENV DATAMESH_DATA=/var/datamesh/data
ENV DATAMESH_LOG_LEVEL=1

#Data
VOLUME $DATAMESH_DATA

# Ports
EXPOSE 6543

# Source
WORKDIR $DATAMESH_HOME
COPY . .

# Build
RUN go build

# Run
CMD ./datamesh -logtostderr -v $DATAMESH_LOG_LEVEL -dir $DATAMESH_DATA server
