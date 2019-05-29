项目依赖memberlist作为gossip协议库</br>
go get github.com/hashicorp/memberlist</br>

依赖consistent作为一致性散列库</br>
go get stathat.com/c/consistent</br>

macos上编译</br>
export GOOS=linux</br>
export GOARCH=amd64</br>
go build -o cached-server</br>