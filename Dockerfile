FROM hub.xxx.com/standard/debian:buster-slim
# FROM行由前端生成，请不要改动
ADD ./config/Flow.conf  /etc/supervisord.conf
ADD  ./git-resource/bin/Flow /home/www/Flow/bin/Flow
ADD ./git-resource/conf /home/www/Flow/conf
ADD ./git-resource/static /home/www/Flow/static
RUN chmod +x /home/www/Flow/bin/Flow
WORKDIR /home/www/Flow/
CMD ["/usr/local/bin/supervisord", "-c", "/etc/supervisord.conf"]