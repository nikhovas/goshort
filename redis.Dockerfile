FROM nikhovas/goshort:alpine as goshort


FROM redis:alpine
COPY --from=goshort /bin/goshort /bin/goshort
COPY redis-docker-commands.sh /scripts/redis-docker-commands.sh
RUN ["chmod", "+x", "/scripts/redis-docker-commands.sh"]
ENTRYPOINT ["sh", "/scripts/redis-docker-commands.sh"]