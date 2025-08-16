FROM python:3.13-slim

WORKDIR /app

# Copy project files
COPY pyproject.toml ./
COPY uv.lock ./
COPY ingest ./ingest

# Install uv and install the package
RUN pip install --upgrade uv && \
    uv install .

CMD ["ingest", "--help"]
