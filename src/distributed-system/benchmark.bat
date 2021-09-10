# 先自己安装 Redis， 到对应目录里面执行

mkdir result
./redis-benchmark.exe -d 10 > ./result/10.txt
./redis-benchmark.exe -d 20 > ./result/20.txt
./redis-benchmark.exe -d 50 > ./result/50.txt
./redis-benchmark.exe -d 100 > ./result/100.txt
./redis-benchmark.exe -d 1024 > ./result/1024.txt
./redis-benchmark.exe -d 5120 > ./result/5120.txt