#测试 http
./cache-benchmark -type http -n 100000 -r 100000 -t set

#测试 tcp
./cache-benchmark -type tcp -n 100000 -r 100000 -t set

#查看状态
curl http://localhost:6800/status