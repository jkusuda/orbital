.PHONY: build run clean

build:
	cmake -S . -B build
	cmake --build build

run: build
	./build/engine

clean:
	rm -rf build
