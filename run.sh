#!/bin/bash

printf "\ec"
go run . setblobs '{"id": 123456, "title": "hello", "description": "yoyoyo"}' '{"id": 1456, "title": "hello", "description": "yoyoyo"}'
go run . set some_key some_value some_other_key
go run . get 123456 some_key
go run . get some_other_key # should be missing, since there was no value
