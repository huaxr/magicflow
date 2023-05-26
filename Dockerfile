FROM hub.xesv5.com/standard/debian:buster-slim
# FROM行由前端生成，请不要改动
ADD ./config/xesFlow.conf  /etc/supervisord.conf
ADD  ./git-resource/bin/xesFlow /home/www/xesFlow/bin/xesFlow
ADD ./git-resource/conf /home/www/xesFlow/conf
ADD ./git-resource/static /home/www/xesFlow/static
RUN chmod +x /home/www/xesFlow/bin/xesFlow
WORKDIR /home/www/xesFlow/
CMD ["/usr/local/bin/supervisord", "-c", "/etc/supervisord.conf"]