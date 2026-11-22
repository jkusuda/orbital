cmake_minimum_required(VERSION 3.20)
project(orbital_engine LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(Protobuf CONFIG REQUIRED)
find_package(gRPC CONFIG REQUIRED)

set(PROTO_DIR ${CMAKE_CURRENT_SOURCE_DIR}/../proto)

# A small library holding the generated code, built from the .proto
add_library(engine_proto ${PROTO_DIR}/orbital/v1/engine.proto)
target_link_libraries(engine_proto PUBLIC protobuf::libprotobuf gRPC::grpc++)
target_include_directories(engine_proto PUBLIC ${CMAKE_CURRENT_BINARY_DIR})

# Messages -> engine.pb.h / engine.pb.cc
protobuf_generate(TARGET engine_proto
  IMPORT_DIRS ${PROTO_DIR}
  PROTOC_OUT_DIR ${CMAKE_CURRENT_BINARY_DIR})

# Service -> engine.grpc.pb.h / engine.grpc.pb.cc
protobuf_generate(TARGET engine_proto
  LANGUAGE grpc
  GENERATE_EXTENSIONS .grpc.pb.h .grpc.pb.cc
  PLUGIN "protoc-gen-grpc=\$<TARGET_FILE:gRPC::grpc_cpp_plugin>"
  IMPORT_DIRS ${PROTO_DIR}
  PROTOC_OUT_DIR ${CMAKE_CURRENT_BINARY_DIR})

# Your program
add_executable(engine src/main.cpp src/propagate.cpp)
target_link_libraries(engine PRIVATE engine_proto)
