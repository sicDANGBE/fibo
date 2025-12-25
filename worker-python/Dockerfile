FROM python:3.12-slim

WORKDIR /app

RUN pip install --no-cache-dir grpcio grpcio-tools

# Copie du proto (contexte racine)
COPY proto/sync.proto .

# Génération du code Python
RUN python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. sync.proto

# --- CORRECTION ICI ---
# On spécifie le dossier source 'python/'
COPY python/main.py .

CMD ["python", "main.py"]