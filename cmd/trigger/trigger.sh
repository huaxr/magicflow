#!/usr/bin/env bash
num=0
while (( $num < 10000 ))
do
    echo "$num"
    num= $num + 1
    curl  -H "Content-Type: application/json"  -X POST -d  '{"data":1000,"app_token":"dXpuYTVjNzZpYWNuaThx","app_name":"first_example","playbook_id":138}' http://127.0.0.1:8080/trigger/execute
done

curl  -H "Content-Type: application/json"  -X POST -d  '{"data":1000,"app_token":"dXpuYTVjNzZpYWNuaThx","app_name":"first_example","playbook_id":138, "sync":true}' http://127.0.0.1:8080/trigger/execute


curl  -H "Content-Type: application/json"  -X POST -d  '{"data":{"Payload":{"info": "{\"description\":\"tmp\",\"isScheme\":\"1\",\"title\":\"hello\",\"userIds\":[\"77155378\"]}","touchBizIds":"xxk-reach-55"}}
,"app_token":"NTRhcWUybGd5bG5rbG4y","app_name":"test_app","playbook_id":2,"sync":false}' http://flow-service.bccv5.vdyoo.com/trigger/execute

curl  -H "Content-Type: application/json"  -X POST -d  '{"data":{"Payload":{"info": "{\"description\":\"tmp\",\"isScheme\":\"1\",\"title\":\"hello\",\"userIds\":[\"77155378\"]}","touchBizIds":"xxk-reach-55"}}
,"app_token":"dXpuYTVjNzZpYWNuaThx","app_name":"first_example","playbook_id":141,"sync":false}' http://127.0.0.1:8080/trigger/execute
