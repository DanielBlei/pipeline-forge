# Development Environment

Local development setup for Pipeline Forge with automated database provisioning and testing infrastructure.

## Features

- **Docker Compose Setup** - MySQL and PostgreSQL databases for testing
- **Database Initialization** - Automated schema and sample data setup
- **Environment Configuration** - Development-specific configurations
- **Testing Infrastructure** - Isolated environments for integration testing

## Quick Start

```bash
# Start development databases
docker-compose up -d

# Initialize databases
./init/mysql/mysql-init.sql
./init/postgres/postgres-init.sql
```

## What's Included

### **MySQL Database**

- **Port**: 3306
- **Database**: `pipeline_forge`
- **User**: `pipeline_user`
- **Password**: Managed via Docker secrets
- **Sample Tables**: `account`, `events`, `users`

### **PostgreSQL Database**

- **Port**: 5432
- **Database**: `pipeline_forge`
- **User**: `pipeline_user`
- **Password**: Managed via Docker secrets
- **Sample Tables**: `account`, `events`, `users`

## Configuration

The development environment provides:

- **Isolated Databases** - Separate containers for MySQL and PostgreSQL
- **Sample Data** - Pre-populated tables for testing ingestion workflows
- **Environment Variables** - Development-specific configurations
- **Network Isolation** - Databases accessible only from development containers

## Usage with Workloads

### **Ingest Workload Testing**

```bash
# Test MySQL ingestion
cd ../workloads/ingest
uv run ingest --config tmp/config.yaml --catalog tmp/catalog.yaml --env dev

# Test PostgreSQL ingestion
# Update config.yaml to use PostgreSQL source
uv run ingest --config tmp/config.yaml --catalog tmp/catalog.yaml --env dev
```

### **Integration Testing**

```bash
# Run full integration tests
cd ../workloads/ingest
make test

# Test with specific database
docker-compose -f ../../dev/docker-compose.yml up -d mysql
make test
```

## Development Workflow

1. **Start Databases**

   ```bash
   docker-compose up -d
   ```

2. **Verify Connection**

   ```bash
   docker-compose ps
   ```

3. **Run Tests**

   ```bash
   cd ../workloads/ingest
   make test
   ```

4. **Clean Up**
   ```bash
   docker-compose down -v
   ```

## Troubleshooting

### **Port Conflicts**

If ports 3306 or 5432 are already in use:

```bash
# Check what's using the ports
lsof -i :3306
lsof -i :5432

# Stop conflicting services or change ports in docker-compose.yml
```

### **Database Connection Issues**

```bash
# Check container logs
docker-compose logs mysql
docker-compose logs postgres

# Restart services
docker-compose restart
```

### **Data Persistence**

```bash
# Remove volumes to start fresh
docker-compose down -v
docker-compose up -d
```

## Environment Variables

| Variable              | Default            | Description              |
| --------------------- | ------------------ | ------------------------ |
| `MYSQL_ROOT_PASSWORD` | `rootpassword`     | MySQL root password      |
| `MYSQL_DATABASE`      | `pipeline_forge`   | MySQL database name      |
| `MYSQL_USER`          | `pipeline_user`    | MySQL user               |
| `POSTGRES_DB`         | `pipeline_forge`   | PostgreSQL database name |
| `POSTGRES_USER`       | `pipeline_user`    | PostgreSQL user          |
| `POSTGRES_PASSWORD`   | `postgrespassword` | PostgreSQL password      |
