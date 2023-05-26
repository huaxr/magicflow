#!/bin/bash

# 1 disc node，2 ram node：
#reset first node
echo "Reset first rabbitmq node."
docker exec rabbitmq1 /bin/bash -c 'rabbitmqctl stop_app'
docker exec rabbitmq1 /bin/bash -c 'rabbitmqctl reset'
docker exec rabbitmq1 /bin/bash -c 'rabbitmqctl start_app'

#build cluster
echo "Starting to build rabbitmq cluster with two ram nodes."
docker exec rabbitmq2 /bin/bash -c 'rabbitmqctl stop_app'
docker exec rabbitmq2 /bin/bash -c 'rabbitmqctl reset'
docker exec rabbitmq2 /bin/bash -c 'rabbitmqctl join_cluster --ram rabbit@rabbitmq1'
docker exec rabbitmq2 /bin/bash -c 'rabbitmqctl start_app'

docker exec rabbitmq3 /bin/bash -c 'rabbitmqctl stop_app'
docker exec rabbitmq3 /bin/bash -c 'rabbitmqctl reset'
docker exec rabbitmq3 /bin/bash -c 'rabbitmqctl join_cluster --ram rabbit@rabbitmq1'
docker exec rabbitmq3 /bin/bash -c 'rabbitmqctl start_app'

#check cluster status
echo "Check cluster status:"
docker exec rabbitmq1 /bin/bash -c 'rabbitmqctl cluster_status'
docker exec rabbitmq2 /bin/bash -c 'rabbitmqctl cluster_status'
docker exec rabbitmq3 /bin/bash -c 'rabbitmqctl cluster_status'

echo "Starting to create user."
docker exec rabbitmq1 /bin/bash -c 'rabbitmqctl add_user admin admin@123'

echo "Set tags for new user."
docker exec rabbitmq1 /bin/bash -c 'rabbitmqctl set_user_tags admin administrator'

echo "Grant permissions to new user."
docker exec rabbitmq1 /bin/bash -c "rabbitmqctl set_permissions -p '/' admin '.*' '.*' '.*'"