#!/bin/sh

set -uex

go build server.go || true
clang -O2 -Wall -target bpf -c bpf.c -o bpf.o

iptables -t nat -D OUTPUT -j SERVICE-UNMESH || true
iptables -t nat -F SERVICE-UNMESH || true
iptables -t nat -X SERVICE-UNMESH || true
iptables -t nat -N SERVICE-UNMESH
iptables -t nat -A SERVICE-UNMESH -p tcp --dport 3333 -j REDIRECT --to-port 6666
iptables -t nat -I OUTPUT 1 -j SERVICE-UNMESH

rm /sys/fs/bpf/getsockopt || true
bpftool prog load bpf.o /sys/fs/bpf/getsockopt
bpftool cgroup attach /sys/fs/cgroup getsockopt pinned /sys/fs/bpf/getsockopt
