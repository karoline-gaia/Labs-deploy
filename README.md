# Weather Service - Sistema de Consulta de Clima por CEP

Sistema em Go que recebe um CEP brasileiro, identifica a cidade e retorna o clima atual em Celsius, Fahrenheit e Kelvin. Desenvolvido para deploy no Google Cloud Run.

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

## 📡 Endpoints

### GET /weather/{cep}

Retorna a temperatura atual para o CEP informado.

**Exemplo de requisição:**

```bash
curl http://localhost:8080/weather/01310100
```

**Respostas:**

#### Sucesso (200)
```json
{
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.5
}
```

#### CEP inválido (422)
```json
{
  "message": "invalid zipcode"
}
```

#### CEP não encontrado (404)
```json
{
  "message": "can not find zipcode"
}
```

### GET /

Health check do serviço.

```bash
curl http://localhost:8080/
```

Resposta:
```json
{
  "status": "ok"
}
```

## 🧪 Exemplos de Teste

```bash
# CEP válido (Av. Paulista, São Paulo)
curl http://localhost:8080/weather/01310100

# CEP válido com hífen
curl http://localhost:8080/weather/01310-100

# CEP inválido (formato incorreto)
curl http://localhost:8080/weather/123

# CEP não encontrado
curl http://localhost:8080/weather/99999999
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

#### 3. Fazer Upload dos Arquivos
- Clique nos **3 pontos `⋮`** → **"Upload"**
- Selecione todos os arquivos do projeto
- Aguarde o upload concluir

#### 4. Executar o Deploy

```bash
# Definir sua API Key do WeatherAPI
export WEATHER_API_KEY="sua_chave_api_aqui"

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

#### 5. Obter a URL do Serviço

Após o deploy, copie a URL exibida no terminal ou execute:

```bash
gcloud run services describe weather-service \
  --region us-central1 \
  --format 'value(status.url)'
```

#### 6. Testar o Serviço

```bash
# Substitua pela sua URL
curl https://weather-service-xxxxx-uc.a.run.app/weather/01310100
```

### Comandos Úteis

```bash
# Ver logs em tempo real
gcloud run services logs tail weather-service --region us-central1

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
- **ViaCEP API**: Consulta de CEPs brasileiros
- **WeatherAPI**: Consulta de dados meteorológicos
- **Docker**: Containerização
- **Google Cloud Run**: Hospedagem serverless

## 📝 Conversões de Temperatura

- **Celsius para Fahrenheit**: F = C × 1.8 + 32
- **Celsius para Kelvin**: K = C + 273

## 🧪 Cobertura de Testes

Os testes automatizados (`main_test.go`) cobrem:

- ✅ Validação de formato de CEP (8 dígitos, com/sem hífen)
- ✅ Conversões de temperatura (Celsius → Fahrenheit, Kelvin)
- ✅ Respostas HTTP corretas para cada cenário (200, 404, 422)
- ✅ Tratamento de erros e edge cases
- ✅ Health check endpoint

**Executar testes:**
```bash
go test -v -cover
```

## 📄 Licença

Este projeto foi desenvolvido para fins educacionais como parte de um desafio técnico.

## 👤 Autor

Desenvolvido em Go com foco em boas práticas, testes automatizados e deploy em cloud.
