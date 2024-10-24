FROM grafana/grafana:latest

COPY --chown=grafana:grafana gsync .
RUN chmod +x gsync
#may need to chown gsync to grafana user

RUN ./gsync &
# RUN the grafana startup command
