# MQCache Go implementation
[![Go Report](https://goreportcard.com/badge/github.com/satmaelstorm/mqcache)](https://goreportcard.com/report/github.com/satmaelstorm/mqcache) 
[![GoDoc](https://godoc.org/github.com/satmaelstorm/mqcache?status.svg)](http://godoc.org/github.com/satmaelstorm/mqcache)
[![Coverage Status](https://coveralls.io/repos/github/satmaelstorm/mqcache/badge.svg?branch=master)](https://coveralls.io/github/satmaelstorm/mqcache?branch=master) 
![Go](https://github.com/satmaelstorm/mqcache/workflows/Go/badge.svg)

Multi-Queue Cache implementation on Go

[Original article](https://www.usenix.org/legacy/events/usenix01/full_papers/zhou/zhou.pdf)

Implementation of the in-memory cache based on the eviction algorithm MQCache. 
The size of the cache can be specified in both elements and bytes.
