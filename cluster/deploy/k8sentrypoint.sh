#!/usr/bin/env bash
mkdir -p /home/logs/xeslog/openapi/$(hostname)

cd  /home/www/openapi.xesv5.com/http
for i in `find ./Config/`; do if [ "${i##*.}" = online ]; then cp $i `echo $i | sed 's/\.online//'`; fi; done

if [ ${KUBERNETES_MODE} = 'doubleAlive' ];then
  for i in `find ./Config/`; do if [ "${i##*.}" = online_ali ]; then cp $i `echo $i | sed 's/\.online_ali//'`; fi; done
  rm -fr  /etc/supervisor/conf.d/confd-tw.conf
else
   rm -fr /etc/supervisor/conf.d/confd-tw_ali.conf
fi
cd  /home/www/openapi.xesv5.com
mv  http/* ./

/usr/bin/supervisord -n -c /etc/supervisord.conf