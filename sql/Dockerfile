FROM postgres:15-alpine

ENV POSTGRES_USER root
ENV POSTGRES_PASSWORD secret
ENV POSTGRES_DB store_management

COPY users.sql /docker-entrypoint-storedb.d/
COPY products.sql /docker-entrypoint-storedb.d/

EXPOSE 5432

CMD ["postgres"]