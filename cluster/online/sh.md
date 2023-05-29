wget 10.20.109.70:8888/nsq-1.2.0.linux-amd64.go1.12.9.tar.gz

#  flow-service.xxx.com address should be changed
10.20.109.51  nsq etcd nsqlookup nsqadmin grafana
10.20.109.70  nsq etcd nsqlookup infludb
10.20.109.75  nsq etcd nsqlookup 
10.20.109.82  nsq 
10.20.109.147 nsq 
10.20.109.160 nsq 
10.20.109.183 nsq 
10.20.109.215 nsq
10.20.109.233 nsq
10.20.109.234 nsq

groupadd -r dba
useradd -r -g dba -G root tnuser
tar -zxvf nsq-1.2.0.linux-amd64.go1.12.9.tar.gz
ln -s nsq-1.2.0.linux-amd64.go1.12.9 nsq
mkdir -p /usr/local/nsq/data/
chown -R tnuser.dba /usr/local/nsq-1.2.0.linux-amd64.go1.12.9 /usr/local/nsq/data/

# 10.20.109.51
nohup /usr/local/nsq/bin/nsqlookupd -http-address 10.20.109.51:4161 -tcp-address 10.20.109.51:4160 -broadcast-address 10.20.109.51 >./nsqlookup.log 2>&1 &
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.20.109.51:4150 -http-address 10.20.109.51:4151 -lookupd-tcp-address 10.20.109.51:4160 -lookupd-tcp-address 10.20.109.70:4160 -lookupd-tcp-address 10.20.109.75:4160 --auth-http-address=flow-service.xxx.com -broadcast-address 10.20.109.51 -msg-timeout=1h -max-msg-timout=2h >./nsqd.log 2>&1 &
nohup /usr/local/nsq/bin/nsqadmin -lookupd-http-address 10.20.109.51:4161 -lookupd-http-address 10.20.109.70:4161 -lookupd-http-address 10.20.109.75:4161 >./nsqadmin.log 2>&1 &

# 10.20.109.70
nohup /usr/local/nsq/bin/nsqlookupd -http-address 10.20.109.70:4161 -tcp-address 10.20.109.70:4160 -broadcast-address 10.20.109.70 >./nsqlookup.log 2>&1 &
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.20.109.70:4150 -http-address 10.20.109.70:4151 -lookupd-tcp-address 10.20.109.51:4160 -lookupd-tcp-address 10.20.109.70:4161 -lookupd-tcp-address 10.20.109.75:4161 --auth-http-address=flow-service.xxx.com -broadcast-address 10.20.109.70 >./nsqd.log 2>&1 &

# 10.20.109.75
nohup /usr/local/nsq/bin/nsqlookupd -http-address 10.20.109.75:4161 -tcp-address 10.20.109.75:4160 -broadcast-address 10.20.109.75 >./nsqlookup.log 2>&1 &
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.20.109.75:4150 -http-address 10.20.109.75:4151 -lookupd-tcp-address 10.20.109.51:4160 -lookupd-tcp-address 10.20.109.70:4161 -lookupd-tcp-address 10.20.109.75:4161 --auth-http-address=flow-service.xxx.com -broadcast-address 10.20.109.75 >./nsqd.log 2>&1 &

# 10.20.109.82
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.20.109.82:4150 -http-address 10.20.109.82:4151 -lookupd-tcp-address 10.20.109.51:4160 -lookupd-tcp-address 10.20.109.70:4160 -lookupd-tcp-address 10.20.109.75:4160 --auth-http-address=flow-service.xxx.com -broadcast-address 10.20.109.82 >./nsqd.log 2>&1 &

# 10.20.109.147
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.20.109.147:4150 -http-address 10.20.109.147:4151 -lookupd-tcp-address 10.20.109.51:4160 -lookupd-tcp-address 10.20.109.70:4160 -lookupd-tcp-address 10.20.109.75:4160 --auth-http-address=flow-service.xxx.com -broadcast-address 10.20.109.147 >./nsqd.log 2>&1 &

# 10.20.109.160
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.20.109.160:4150 -http-address 10.20.109.160:4151 -lookupd-tcp-address 10.20.109.51:4160 -lookupd-tcp-address 10.20.109.70:4160 -lookupd-tcp-address 10.20.109.75:4160 --auth-http-address=flow-service.xxx.com -broadcast-address 10.20.109.160 >./nsqd.log 2>&1 &

# 10.20.109.183
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.20.109.183:4150 -http-address 10.20.109.183:4151 -lookupd-tcp-address 10.20.109.51:4160 -lookupd-tcp-address 10.20.109.70:4160 -lookupd-tcp-address 10.20.109.75:4160 --auth-http-address=flow-service.xxx.com -broadcast-address 10.20.109.183 >./nsqd.log 2>&1 &

# 10.20.109.215
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.20.109.215:4150 -http-address 10.20.109.215:4151 -lookupd-tcp-address 10.20.109.51:4160 -lookupd-tcp-address 10.20.109.70:4160 -lookupd-tcp-address 10.20.109.75:4160 --auth-http-address=flow-service.xxx.com -broadcast-address 10.20.109.215 >./nsqd.log 2>&1 &

# 10.20.109.233
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.20.109.233:4150 -http-address 10.20.109.233:4151 -lookupd-tcp-address 10.20.109.51:4160 -lookupd-tcp-address 10.20.109.70:4160 -lookupd-tcp-address 10.20.109.75:4160 --auth-http-address=flow-service.xxx.com -broadcast-address 10.20.109.233 >./nsqd.log 2>&1 &

# 10.20.109.234
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.20.109.234:4150 -http-address 10.20.109.234:4151 -lookupd-tcp-address 10.20.109.51:4160 -lookupd-tcp-address 10.20.109.70:4160 -lookupd-tcp-address 10.20.109.75:4160 --auth-http-address=flow-service.xxx.com -broadcast-address 10.20.109.234 >./nsqd.log 2>&1 &
