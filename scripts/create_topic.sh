#!/bin/bash
kafka-topics --create --topic event --bootstrap-server kafka:29092 --replication-factor 1 --partitions 1 || echo "Topic already exists, ignoring error."
