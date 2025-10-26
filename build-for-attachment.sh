#!/bin/bash
go clean
go build -gcflags="all=-N -l"