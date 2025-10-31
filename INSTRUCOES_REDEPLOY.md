# 🚀 Instruções para Redeploy

## O que foi corrigido?

### ✅ Melhorias Implementadas (v1.1)

1. **Logs Detalhados:**
   - Cada requisição agora é rastreada do início ao fim
   - Logs mostram: CEP recebido, localização encontrada, temperatura obtida
   - Erros incluem detalhes específicos para facilitar debug

2. **Validação Aprimorada:**
   - Verifica se a API Key está configurada antes de fazer requisições
   - Retorna mensagens de erro mais específicas

3. **Tratamento de Erros Robusto:**
   - Captura e loga erros da Weather API com detalhes
   - Diferencia entre erro de conexão e erro de API
   - Facilita identificar a causa raiz do problema

## 📋 Como Fazer o Redeploy

### Passo 1: Acessar Google Cloud Shell

1. Acesse: https://console.cloud.google.com/
2. Clique no ícone **`>_`** (Cloud Shell) no canto superior direito
3. Aguarde o terminal abrir

### Passo 2: Atualizar o Código

```bash
# Entrar no diretório do projeto
cd ~/Labs-deploy

# Baixar as últimas alterações
git pull origin main
```

### Passo 3: Fazer o Redeploy

**Opção A - Usando o script (mais fácil):**

```bash
./redeploy.sh
```

**Opção B - Manualmente:**

```bash
# Configurar projeto
gcloud config set project fullcycle01

# Definir API Key
export WEATHER_API_KEY="e3a7bd5d41e0453391a220046252810"

# Deploy
gcloud run deploy weather-service \
  --source . \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=$WEATHER_API_KEY \
  --memory 256Mi \
  --cpu 1 \
  --max-instances 10
```

### Passo 4: Aguardar o Deploy

O processo leva cerca de 3-5 minutos. Você verá mensagens como:

```
Building using Dockerfile and deploying container...
✓ Creating Container Repository...
✓ Uploading sources...
✓ Building Container...
✓ Creating Revision...
✓ Routing traffic...
Done.
```

### Passo 5: Testar o Serviço

```bash
# Teste básico
curl https://weather-service-175512104676.us-central1.run.app/weather/01310100

# Deve retornar algo como:
# {"temp_C":25.0,"temp_F":77.0,"temp_K":298.0}
```

## 🔍 Verificar Logs

Após fazer uma requisição, você pode ver os logs detalhados:

```bash
# Ver logs em tempo real
gcloud run services logs tail weather-service --region us-central1

# Ver últimas 50 linhas
gcloud run services logs read weather-service --region us-central1 --limit 50
```

**Exemplo de logs esperados (sucesso):**

```
Received request for CEP: 01310100
Found location for CEP 01310100: São Paulo,SP
Fetching weather for location: São Paulo,SP
Successfully fetched temperature for São Paulo,SP: 25.0°C
Successfully processed CEP 01310100: 25.0°C, 77.0°F, 298.0°K
```

**Exemplo de logs esperados (erro):**

```
Received request for CEP: 01310100
Found location for CEP 01310100: São Paulo,SP
Fetching weather for location: São Paulo,SP
ERROR: Weather API returned status 401 for location: São Paulo,SP
Weather API error details: map[error:map[code:1002 message:API key is invalid]]
ERROR: Failed to get temperature for location 'São Paulo,SP': weather API error: status 401
```

## ✅ Checklist de Verificação

Após o redeploy, verifique:

- [ ] Deploy concluído sem erros
- [ ] URL retornada no final do deploy
- [ ] Teste com CEP válido retorna temperaturas
- [ ] Teste com CEP inválido retorna erro 422
- [ ] Teste com CEP não encontrado retorna erro 404
- [ ] Logs mostram informações detalhadas

## 🆘 Problemas Comuns

### Erro: "API key is invalid"

**Solução:** Verificar se a API Key está configurada corretamente:

```bash
gcloud run services describe weather-service \
  --region us-central1 \
  --format 'value(spec.template.spec.containers[0].env[0].value)'
```

Se estiver vazia ou incorreta, atualizar:

```bash
gcloud run services update weather-service \
  --region us-central1 \
  --set-env-vars WEATHER_API_KEY=e3a7bd5d41e0453391a220046252810
```

### Erro: "Permission denied"

**Solução:** Adicionar permissões necessárias:

```bash
PROJECT_NUMBER=$(gcloud projects describe fullcycle01 --format="value(projectNumber)")

gcloud projects add-iam-policy-binding fullcycle01 \
  --member="serviceAccount:${PROJECT_NUMBER}-compute@developer.gserviceaccount.com" \
  --role="roles/editor"
```

## 📞 Suporte

Se continuar com problemas:

1. **Verificar logs:** `gcloud run services logs read weather-service --region us-central1 --limit 100`
2. **Testar localmente:** `docker-compose up` e testar em `http://localhost:8080`
3. **Verificar status do serviço:** `gcloud run services describe weather-service --region us-central1`

---

**Última atualização:** 2025-10-31
**Versão:** 1.1
