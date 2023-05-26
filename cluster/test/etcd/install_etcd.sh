#!/usr/bin/env bash
# ref: https://www.cnblogs.com/ilifeilong/p/11622107.html

groupadd -r dba
useradd -r -g dba -G root tnuser

cd /usr/local/
wget https://github.com/etcd-io/etcd/releases/download/v3.4.1/etcd-v3.4.1-linux-amd64.tar.gz --no-check-certificate
tar -zxf etcd-v3.4.1-linux-amd64.tar.gz
ln -s etcd-v3.4.1-linux-amd64 etcd
mkdir /usr/local/etcd/data
chown -R tnuser.dba /usr/local/etcd-v3.4.1-linux-amd64 /usr/local/etcd/data

# vim /usr/lib/systemd/system/etcd.service
# systemctl daemon-reload
# systemctl start etcd
# systemctl status etcd.service
systemctl daemon-reload
systemctl start etcd