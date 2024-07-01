FROM traefik:v2.9

COPY traefik.yml /etc/traefik/traefik.yml
COPY dynamic_config.yml /etc/traefik/dynamic_config.yml

CMD ["traefik", "--configFile=/etc/traefik/traefik.yml"]