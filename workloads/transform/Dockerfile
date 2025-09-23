FROM python:3.13-slim

COPY pipeline_forge .

WORKDIR /pipeline_forge

RUN pip install dbt-core==1.10.11 && \
    pip install dbt-bigquery==1.10.2

CMD ["dbt", "-h"]
