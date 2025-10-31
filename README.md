# Weather Service - Sistema de Consulta de Clima por CEP

Sistema em Go que recebe um CEP brasileiro, identifica a cidade e retorna o clima atual em Celsius, Fahrenheit e Kelvin.

## 🌐 Serviço em Produção

**URL do serviço deployado no Google Cloud Run:**
```
https://weather-service-175512104676.us-central1.run.app
```

### 🧪 Testes Rápidos

```bash
# Health Check
curl https://weather-service-175512104676.us-central1.run.app/

# CEP válido (São Paulo - Av. Paulista)
curl https://weather-service-175512104676.us-central1.run.app/weather/01310100

# CEP inválido (retorna 422)
curl https://weather-service-175512104676.us-central1.run.app/weather/123

# CEP não encontrado (retorna 404)
curl https://weather-service-175512104676.us-central1.run.app/weather/99999999
```

### 🌍 Testar no Navegador

- **São Paulo:** https://weather-service-175512104676.us-central1.run.app/weather/01310100
- **Rio de Janeiro:** https://weather-service-175512104676.us-central1.run.app/weather/20040020
- **Belo Horizonte:** https://weather-service-175512104676.us-central1.run.app/weather/30130100

## 📋 Requisitos

- Go 1.21 ou superior (para desenvolvimento local)
- Docker e Docker Compose (para testes locais)
- Chave de API do WeatherAPI (gratuita em https://www.weatherapi.com/)
- Conta Google Cloud Platform (para deploy)

## 🚀 Como Executar Localmente

### 1. Configurar a API Key

Crie um arquivo `.env` na raiz do projeto:

```bash
cp .env.example .env
```

Edite o arquivo `.env` e adicione sua chave de API do WeatherAPI:

```
WEATHER_API_KEY=sua_chave_api_aqui
```

### 2. Executar com Docker Compose

```bash
docker-compose up --build
```

O serviço estará disponível em `http://localhost:8080`

### 3. Executar sem Docker

```bash
# Instalar dependências
go mod download

# Configurar a variável de ambiente
export WEATHER_API_KEY=sua_chave_api_aqui

# Executar
go run main.go
```

## 🧪 Executar Testes

```bash
# Executar todos os testes
go test -v

# Executar testes com cobertura
go test -v -cover

# Gerar relatório de cobertura
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 📡 Endpoints da API

### GET /weather/{cep}

Retorna a temperatura atual para o CEP informado.

**Formato do CEP:** 8 dígitos (com ou sem hífen)

**Exemplo de requisição:**

```bash
# Produção
curl https://weather-service-175512104676.us-central1.run.app/weather/01310100

# Local
curl http://localhost:8080/weather/01310100
```

**Respostas:**

#### ✅ Sucesso (200 OK)
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

#### ❌ CEP Inválido (422 Unprocessable Entity)
Quando o CEP não possui 8 dígitos ou contém caracteres inválidos.

```json
{
  "message": "invalid zipcode"
}
```

#### ❌ CEP Não Encontrado (404 Not Found)
Quando o CEP é válido mas não existe na base do ViaCEP.

```json
{
  "message": "can not find zipcode"
}
```

### GET /

Health check do serviço.

```bash
# Produção
curl https://weather-service-175512104676.us-central1.run.app/

# Local
curl http://localhost:8080/
```

**Resposta:**
```json
{
  "status": "ok"
}
```

## 🧪 Exemplos de Teste Completos

### Testando o Serviço em Produção

```bash
# 1. Health Check
curl https://weather-service-175512104676.us-central1.run.app/
# Resposta esperada: {"status":"ok"}

# 2. CEP válido (Av. Paulista, São Paulo)
curl https://weather-service-175512104676.us-central1.run.app/weather/01310100
# Resposta esperada: {"temp_C":XX.X,"temp_F":XX.X,"temp_K":XXX.X}

# 3. CEP válido com hífen
curl https://weather-service-175512104676.us-central1.run.app/weather/01310-100
# Resposta esperada: {"temp_C":XX.X,"temp_F":XX.X,"temp_K":XXX.X}

# 4. CEP inválido (formato incorreto) - Retorna 422
curl https://weather-service-175512104676.us-central1.run.app/weather/123
# Resposta esperada: {"message":"invalid zipcode"}

# 5. CEP não encontrado - Retorna 404
curl https://weather-service-175512104676.us-central1.run.app/weather/99999999
# Resposta esperada: {"message":"can not find zipcode"}
```

### Outros CEPs para Teste

```bash
# Rio de Janeiro - Centro
curl https://weather-service-175512104676.us-central1.run.app/weather/20040020

# Belo Horizonte - Centro
curl https://weather-service-175512104676.us-central1.run.app/weather/30130100

# Curitiba - Centro
curl https://weather-service-175512104676.us-central1.run.app/weather/80010000

# Porto Alegre - Centro
curl https://weather-service-175512104676.us-central1.run.app/weather/90010000
```

## 🏗️ Estrutura do Projeto

```
weather-service/
├── main.go              # Código principal da aplicação
├── main_test.go         # Testes automatizados em Go
├── go.mod               # Dependências do Go
├── go.sum               # Checksums das dependências
├── Dockerfile           # Container Docker (multi-stage build)
├── docker-compose.yml   # Orquestração Docker
├── .env.example         # Exemplo de variáveis de ambiente
├── .dockerignore        # Arquivos ignorados no build Docker
├── .gcloudignore        # Arquivos ignorados no deploy GCP
├── .gitignore           # Arquivos ignorados pelo Git
└── README.md            # Documentação
```

## 🌐 Deploy no Google Cloud Run

### Método Recomendado: Google Cloud Shell

O método mais fácil é usar o **Google Cloud Shell** (terminal integrado no navegador):

#### 1. Acessar o Console
- Acesse: https://console.cloud.google.com/
- Faça login e crie/selecione um projeto
- Habilite o billing (necessário mesmo para free tier)

#### 2. Abrir Cloud Shell
- Clique no ícone **`>_`** no canto superior direito
- Aguarde o terminal abrir

#### 3. Clonar o Repositório

```bash
# Clonar o projeto
git clone https://github.com/karoline-gaia/Labs-deploy.git
cd Labs-deploy
```

#### 4. Executar o Deploy

```bash
# Configurar projeto
gcloud config set project fullcycle01

# Definir API Key do WeatherAPI
export WEATHER_API_KEY="e3a7bd5d41e0453391a220046252810"

# Deploy (aguarde 3-5 minutos)
gcloud run deploy weather-service \
  --source . \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=$WEATHER_API_KEY \
  --memory 256Mi \
  --cpu 1 \
  --max-instances 10
```

#### 5. Redeploy (Atualizar o Serviço)

Para atualizar o serviço após mudanças no código:

```bash
cd ~/Labs-deploy
git pull origin main
./redeploy.sh
```

Ou manualmente:

```bash
cd ~/Labs-deploy
git pull origin main

export WEATHER_API_KEY="e3a7bd5d41e0453391a220046252810"

gcloud run deploy weather-service \
  --source . \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=$WEATHER_API_KEY
```

### Comandos Úteis

```bash
# Ver logs em tempo real (útil para debug)
gcloud run services logs tail weather-service --region us-central1

# Ver logs das últimas 50 linhas
gcloud run services logs read weather-service --region us-central1 --limit 50

# Obter URL do serviço
gcloud run services describe weather-service \
  --region us-central1 \
  --format 'value(status.url)'

# Atualizar variável de ambiente
gcloud run services update weather-service \
  --region us-central1 \
  --set-env-vars WEATHER_API_KEY=nova_chave

# Deletar o serviço
gcloud run services delete weather-service --region us-central1
```

### Custos

- **GRATUITO** no free tier do Google Cloud Run
- 2 milhões de requisições/mês grátis
- 360.000 GB-segundos de memória/mês grátis
- 180.000 vCPU-segundos/mês grátis

## 🔧 Tecnologias Utilizadas

- **Go 1.21**: Linguagem de programação
- **ViaCEP API**: Consulta de CEPs brasileiros (https://viacep.com.br/)
- **WeatherAPI**: Consulta de dados meteorológicos (https://www.weatherapi.com/)
- **Docker**: Containerização com multi-stage build
- **Google Cloud Run**: Hospedagem serverless
- **testify**: Framework de testes para Go

## 📝 Conversões de Temperatura

As conversões são realizadas conforme especificado:

- **Celsius para Fahrenheit**: `F = C × 1.8 + 32`
- **Celsius para Kelvin**: `K = C + 273`

## 🧪 Cobertura de Testes

Os testes automatizados (`main_test.go`) cobrem:

- ✅ **Validação de formato de CEP**: 8 dígitos, com/sem hífen, caracteres inválidos
- ✅ **Conversões de temperatura**: Precisão das fórmulas C→F e C→K
- ✅ **Respostas HTTP corretas**: Status codes 200, 404, 422
- ✅ **Tratamento de erros**: CEP inválido, CEP não encontrado
- ✅ **Health check endpoint**: Verificação de disponibilidade

**Executar testes localmente:**
```bash
go test -v -cover
```

## 🔍 Troubleshooting

### Erro: "error fetching weather data"

Se você receber este erro, verifique:

1. **API Key configurada corretamente:**
```bash
# No Cloud Shell, verificar se a variável está definida
gcloud run services describe weather-service --region us-central1 --format 'value(spec.template.spec.containers[0].env[0].value)'
```

2. **Ver logs detalhados:**
```bash
# Logs em tempo real
gcloud run services logs tail weather-service --region us-central1

# Últimas 100 linhas
gcloud run services logs read weather-service --region us-central1 --limit 100
```

3. **Testar localmente:**
```bash
export WEATHER_API_KEY="e3a7bd5d41e0453391a220046252810"
docker-compose up
curl http://localhost:8080/weather/01310100
```

### Melhorias Implementadas

**v1.1 - Logs e Tratamento de Erros Aprimorados:**
- ✅ Logs detalhados em cada etapa do processamento
- ✅ Mensagens de erro mais específicas
- ✅ Validação da API Key no início
- ✅ Rastreamento completo de requisições
- ✅ Detalhes de erros da Weather API nos logs

Para ver os logs após fazer uma requisição:
```bash
gcloud run services logs read weather-service --region us-central1 --limit 20
```

## 📊 Requisitos Atendidos

- ✅ Sistema recebe CEP válido de 8 dígitos
- ✅ Realiza pesquisa do CEP via ViaCEP
- ✅ Consulta temperatura via WeatherAPI
- ✅ Retorna temperaturas em Celsius, Fahrenheit e Kelvin
- ✅ Responde com código 200 em caso de sucesso
- ✅ Responde com código 422 para CEP inválido
- ✅ Responde com código 404 para CEP não encontrado
- ✅ Testes automatizados implementados
- ✅ Docker e docker-compose configurados
- ✅ Deploy realizado no Google Cloud Run
- ✅ Endereço ativo e acessível
- ✅ Logs detalhados para debugging
- ✅ Tratamento robusto de erros

## 📄 Licença

Este projeto foi desenvolvido para fins educacionais como parte de um desafio técnico.

## 👤 Autor

Desenvolvido em Go com foco em boas práticas, clean code, testes automatizados e deploy em cloud.
