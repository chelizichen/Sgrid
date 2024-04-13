#!/bin/bash  

# if permission denied
# run script with ` chmod +x build.sh ` 
readonly ServerName="SgridPackageServer"
readonly SgridFile="sgrid_app"
# rm
rm ./$ServerName.tar.gz ./SgridFile

# compile
GOOS=linux GOARCH=amd64 go build -o $SgridFile

# build
tar -cvf $ServerName.tar.gz ./sgrid.yaml ./$SgridFile