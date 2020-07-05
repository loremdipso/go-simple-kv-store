#!/bin/bash

# empirically found to be the sweet spot for the number of threads
export RAYON_NUM_THREADS=2
printf "\ec"

go run . setblobs '{"id": 123456, "title": "hello", "description": "yoyoyo"}' '{"id": 1456, "title": "hello", "description": "yoyoyo"}'
go run . get 123456
#cargo run -q --release ./model/seeta_fd_frontal_v1.0.bin ./temp/test_1.jpg
