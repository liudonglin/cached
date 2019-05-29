#测试 http
./cache-benchmark -type http -n 100000 -r 100000 -t set

#测试 tcp
./cache-benchmark -type tcp -n 100000 -r 100000 -t set

./cache-benchmark -type tcp -n 100000 -r 100000 -t set -pip 10
./cache-benchmark -type tcp -n 100000 -r 100000 -t get -pip 10

#查看状态
curl http://localhost:6800/status