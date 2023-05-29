
#  flow-service.xxx.com address should be changed
02 10.187.114.196  gitlab 
03 10.187.114.243  nsq promethues
04 10.187.114.198  nsq 
05 10.187.114.236  nsq 
06 10.187.114.212  nsq 
07 10.187.114.211  nsq 
08 10.187.114.201  nsq 
10 10.187.115.5    etcd nsqlookup 
11 10.187.114.187  etcd nsqlookup
12 10.187.114.192  etcd nsqlookup nsqadmin 


# 10.187.114.192
nohup /usr/local/nsq/bin/nsqlookupd -http-address 10.187.114.192:4161 -tcp-address 10.187.114.192:4160 -broadcast-address 10.187.114.192 >./nsqlookup.log 2>&1 &
nohup /usr/local/nsq/bin/nsqadmin -lookupd-http-address 10.187.114.192:4161 -lookupd-http-address 10.187.114.187:4161 -lookupd-http-address 10.187.115.5:4161 >./nsqadmin.log 2>&1 &

# 10.187.114.187
nohup /usr/local/nsq/bin/nsqlookupd -http-address 10.187.114.187:4161 -tcp-address 10.187.114.187:4160 -broadcast-address 10.187.114.187 >./nsqlookup.log 2>&1 &

# 10.187.115.5
nohup /usr/local/nsq/bin/nsqlookupd -http-address 10.187.115.5:4161 -tcp-address 10.187.115.5:4160 -broadcast-address 10.187.115.5 >./nsqlookup.log 2>&1 &

# 10.187.114.196

# 10.187.114.243
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.187.114.243:4150 -http-address 10.187.114.243:4151 -lookupd-tcp-address 10.187.114.192:4160 -lookupd-tcp-address 10.187.114.187:4160 -lookupd-tcp-address 10.187.115.5:4160 --auth-http-address=flow-service.xxx.com -broadcast-address 10.187.114.243 -msg-timeout=1h >./nsqd.log 2>&1 &

# 10.187.114.198
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.187.114.198:4150 -http-address 10.187.114.198:4151 -lookupd-tcp-address 10.187.114.192:4160 -lookupd-tcp-address 10.187.114.187:4160 -lookupd-tcp-address 10.187.115.5:4160 --auth-http-address=flow-service.xxx.com -broadcast-address 10.187.114.198 -msg-timeout=1h >./nsqd.log 2>&1 &

# 10.187.114.236
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.187.114.236:4150 -http-address 10.187.114.236:4151 -lookupd-tcp-address 10.187.114.192:4160 -lookupd-tcp-address 10.187.114.187:4160 -lookupd-tcp-address 10.187.115.5:4160 --auth-http-address=flow-service.xxx.com -broadcast-address 10.187.114.236 -msg-timeout=1h >./nsqd.log 2>&1 &

# 10.187.114.212
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.187.114.212:4150 -http-address 10.187.114.212:4151 -lookupd-tcp-address 10.187.114.192:4160 -lookupd-tcp-address 10.187.114.187:4160 -lookupd-tcp-address 10.187.115.5:4160 --auth-http-address=flow-service.xxx.com -broadcast-address 10.187.114.212 -msg-timeout=1h >./nsqd.log 2>&1 &

# 10.187.114.211
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.187.114.211:4150 -http-address 10.187.114.211:4151 -lookupd-tcp-address 10.187.114.192:4160 -lookupd-tcp-address 10.187.114.187:4160 -lookupd-tcp-address 10.187.115.5:4160 --auth-http-address=flow-service.xxx.com -broadcast-address 10.187.114.211 -msg-timeout=1h >./nsqd.log 2>&1 &

# 10.187.114.201
nohup /usr/local/nsq/bin/nsqd -tcp-address 10.187.114.201:4150 -http-address 10.187.114.201:4151 -lookupd-tcp-address 10.187.114.192:4160 -lookupd-tcp-address 10.187.114.187:4160 -lookupd-tcp-address 10.187.115.5:4160 --auth-http-address=flow-service.xxx.com -broadcast-address 10.187.114.201 -msg-timeout=1h >./nsqd.log 2>&1 &
