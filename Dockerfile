FROM grafana/grafana:latest

COPY --chown=grafana:grafana grinksync .
RUN chmod +x grinksync
#may need to chown grinksync to grafana user

RUN ./grinksync &
# RUN the grafana startup command
