FROM python:3.13-slim

WORKDIR /app

# Copy project files
COPY pyproject.toml ./
COPY ingest ./ingest

# Install uv and sync dependencies from pyproject.toml
RUN pip install --upgrade uv && \
    uv sync --system

CMD ["python", "-m", "ingest.main"]
