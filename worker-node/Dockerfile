FROM node:20-alpine

WORKDIR /app

RUN npm install @grpc/grpc-js @grpc/proto-loader

# Copie du proto (contexte racine)
COPY proto/sync.proto .

# --- CORRECTION ICI ---
# On spécifie le dossier source 'node/'
COPY node/index.js .

CMD ["node", "index.js"]