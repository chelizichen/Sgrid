<h2 align="center">Sgrid</h2>

[中文介绍]('./readme_zn.md')

<h3 align="center">Agile cluster management services</h3>

***
<h4 align="center">HomePage</h4>
<img src="./note/grid0424.png" />

<h4 align="center">Release</h4>
<img src="./note/release.png" />

<h4 align="center">Logger</h4>
<img src="./note/logger.png" />

***
Compile

````shell
#  ***************** proto ***************
protoc --go_out=. --go-grpc_out=. SgridPackage.proto

protoc --go_out=. SgridPackage.proto

#  ***************** proto ***************

#  ***************** client ***************
cd client 

npm i

./build.sh

cd ..
#  ***************** client ***************

#  ***************** server ***************
./prod.sh
#  ***************** server ***************
````

***
Sgrid Application Support:

1. Multi node group management
2. Using GRPC to build backend services and Protobuf for communication
3. Multi language support (Node, Java, Go)
4. Support horizontal and vertical expansion, and fast start stop services
5. Dynamic configuration file
6. Continuous heartbeat detection
7. Version control, documentation
8. Multi node log monitoring
9. Multiple services can be started simultaneously within a process, and multiple addresses can be listened to.
