docker-compose.yml for Local:
version: '3.8'

services:
  postgres:
    image: postgres:14
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: mysecretpassword
      POSTGRES_DB: Login
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:8
    restart: always
    ports:
      - "6379:6379"
  
  login_backend:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "80:8080"
    depends_on:
      - postgres
      - redis
    env_file:
      - .env

volumes:
  postgres_data:

docker-compose.yml for AWS Deployment:

version: '3.8'

services:
  postgres:
    image: postgres:14
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: mysecretpassword
      POSTGRES_DB: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./db/migrations:/docker-entrypoint-initdb.d

  redis:
    image: redis:8
    restart: always
    ports:
      - "6379:6379"
  
  login_backend:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "80:8080"
    depends_on:
      - postgres
      - redis
    env_file:
      - .env

volumes:
  postgres_data:

1. Change dbname to "postgres"
2. Command to switch to psql console in aws console: docker compose exec postgres psql -U postgres -d postgres
