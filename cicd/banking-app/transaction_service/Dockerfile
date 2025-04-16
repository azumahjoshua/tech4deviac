FROM python:3.13-alpine

WORKDIR /app

# Install system dependencies
RUN apk add --no-cache \
    curl \
    gcc \
    musl-dev \
    postgresql-dev

# Copy requirements first
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Set Python path
ENV PYTHONPATH=/app \
    PYTHONUNBUFFERED=1

EXPOSE 5000
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "5000"]

# FROM python:3.13-alpine

# WORKDIR /app

# # Install system dependencies
# RUN apk add --no-cache \
#     curl \
#     gcc \
#     musl-dev \
#     postgresql-dev

# # Copy requirements first
# COPY transaction_service/requirements.txt .
# RUN pip install --no-cache-dir -r requirements.txt

# # Copy application code
# COPY transaction_service/ .

# # Set Python path
# ENV PYTHONPATH=/app \
#     PYTHONUNBUFFERED=1

# EXPOSE 5000
# CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "5000"]
# # Stage 1: Build stage
# FROM python:3.13-alpine AS builder

# WORKDIR /app

# # Install curl for health checks
# RUN apk add --no-cache curl

# # Copy requirements first for caching
# COPY transaction_service/requirements.txt .
# RUN pip install --user --no-cache-dir -r requirements.txt

# # Copy the entire service code
# COPY transaction_service/ ./transaction_service/

# # Stage 2: Runtime stage
# FROM python:3.13-alpine

# WORKDIR /app

# # Install curl for health checks
# RUN apk add --no-cache curl

# # Copy installed packages and application
# COPY --from=builder /root/.local /root/.local
# COPY --from=builder /app/transaction_service ./transaction_service

# # Environment setup
# ENV PATH=/root/.local/bin:$PATH
# ENV PYTHONPATH=/app

# # Health check with proper timing
# HEALTHCHECK --interval=30s --timeout=5s --retries=3 --start-period=10s \
#     CMD curl -f http://localhost:5000/health || exit 1

# EXPOSE 5000
# CMD ["uvicorn", "transaction_service.main:app", "--host", "0.0.0.0", "--port", "5000"]