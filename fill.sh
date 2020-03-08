#!/bin/bash
time for i in {1..3000}; do echo "set azaza$i 123_$i" | nc localhost 8080; done
