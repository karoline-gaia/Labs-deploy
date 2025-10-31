#!/bin/bash

# Script para redeploy no Google Cloud Run
# Execute este script no Google Cloud Shell

echo "🚀 Iniciando redeploy do Weather Service..."

# Configurar projeto
gcloud config set project fullcycle01

# Fazer pull das últimas alterações
echo "📥 Baixando últimas alterações do GitHub..."
cd ~/Labs-deploy
git pull origin main

# Configurar API Key
export WEATHER_API_KEY="e3a7bd5d41e0453391a220046252810"

# Fazer deploy
echo "🔨 Fazendo deploy no Cloud Run..."
gcloud run deploy weather-service \
  --source . \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=$WEATHER_API_KEY \
  --memory 256Mi \
  --cpu 1 \
  --max-instances 10

echo "✅ Deploy concluído!"
echo "🌐 URL: https://weather-service-175512104676.us-central1.run.app"
echo ""
echo "🧪 Teste com:"
echo "curl https://weather-service-175512104676.us-central1.run.app/weather/01310100"
