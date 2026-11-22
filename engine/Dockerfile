FROM alpine:3.22 AS build
RUN apk add --no-cache build-base cmake grpc-dev grpc-plugins protobuf-dev
WORKDIR /src
COPY proto proto
COPY engine/CMakeLists.txt engine/
COPY engine/src engine/src
RUN cmake -S engine -B engine/build -DCMAKE_BUILD_TYPE=Release \
 && cmake --build engine/build -j

FROM alpine:3.22
RUN apk add --no-cache grpc-cpp protobuf libstdc++ \
 && adduser -D -H engine
COPY --from=build /src/engine/build/engine /usr/local/bin/engine
USER engine
EXPOSE 50051
ENTRYPOINT ["/usr/local/bin/engine"]
