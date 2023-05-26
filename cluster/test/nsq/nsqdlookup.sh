#!/usr/bin/env bash
cd /usr/local/
wget https://s3.amazonaws.com/bitly-downloads/nsq/nsq-1.2.0.linux-amd64.go1.12.9.tar.gz

groupadd -r dba
useradd -r -g dba -G root tnuser
tar -zxvf nsq-1.2.0.linux-amd64.go1.12.9.tar.gz
ln -s nsq-1.2.0.linux-amd64.go1.12.9 nsq
mkdir -p /usr/local/nsq/data/
chown -R tnuser.dba /usr/local/nsq-1.2.0.linux-amd64.go1.12.9 /usr/local/nsq/data/


# 10.90.72.58
nohup /usr/local/nsq/bin/nsqlookupd -http-address 10.90.72.58:4161 -tcp-address 10.90.72.58:4160 -broadcast-address 10.90.72.58 >./nsqlookup.log 2>&1 &

# 10.90.72.94
nohup /usr/local/nsq/bin/nsqlookupd -http-address 10.90.72.94:4161 -tcp-address 10.90.72.94:4160 -broadcast-address 10.90.72.94 >./nsqlookup.log 2>&1 &

# 10.90.72.135
nohup /usr/local/nsq/bin/nsqadmin -lookupd-http-address 10.90.72.58:4161 -lookupd-http-address 10.90.72.94:4161 >./nsqadmin.log 2>&1 &


# 10.90.72.135, auth-add 不能加 schema
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.90.72.135:4150 -http-address 10.90.72.135:4151 -lookupd-tcp-address 10.90.72.58:4160 -lookupd-tcp-address 10.90.72.94:4160  --auth-http-address=api.xueersi.com/orchestration/ -broadcast-address 10.90.72.135 >./nsqd.log 2>&1 &

# 10.90.72.136
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.90.72.136:4150 -http-address 10.90.72.136:4151 -lookupd-tcp-address 10.90.72.58:4160 -lookupd-tcp-address 10.90.72.94:4160  --auth-http-address=api.xueersi.com/orchestration/ -broadcast-address 10.90.72.136 >./nsqd.log 2>&1 &

# 10.90.72.171
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.90.72.171:4150 -http-address 10.90.72.171:4151 -lookupd-tcp-address 10.90.72.58:4160 -lookupd-tcp-address 10.90.72.94:4160  --auth-http-address=api.xueersi.com/orchestration/ -broadcast-address 10.90.72.171 >./nsqd.log 2>&1 &

# 10.90.72.172
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.90.72.172:4150 -http-address 10.90.72.172:4151 -lookupd-tcp-address 10.90.72.58:4160 -lookupd-tcp-address 10.90.72.94:4160  --auth-http-address=api.xueersi.com/orchestration/ -broadcast-address 10.90.72.172 >./nsqd.log 2>&1 &