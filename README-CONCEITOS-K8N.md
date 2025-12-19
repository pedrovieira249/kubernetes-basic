# Guia Completo: Conceitos Básicos do Kubernetes

## 📚 Índice

1. [Instalação e Configuração](#-instalação-e-configuração)
2. [O que é Kubernetes?](#o-que-é-kubernetes)
3. [Por que usar Kubernetes?](#por-que-usar-kubernetes)
4. [Arquitetura do Kubernetes](#arquitetura-do-kubernetes)
5. [Conceitos Fundamentais](#conceitos-fundamentais)
6. [Objetos do Kubernetes](#objetos-do-kubernetes)
7. [Comandos Básicos](#comandos-básicos)
8. [Arquivos do Projeto](#-arquivos-do-projeto)
9. [Exemplos Práticos](#exemplos-práticos)
10. [Fluxo de Trabalho](#fluxo-de-trabalho)

---

## 🔧 Instalação e Configuração

### Pré-requisitos

Antes de começar com Kubernetes, você precisa ter instalado:

1. **Docker Desktop** (Windows/Mac) ou **Docker Engine** (Linux)
2. **kubectl** - CLI do Kubernetes
3. **Kind** - Kubernetes in Docker (para desenvolvimento local)

---

### 1. Instalando Docker

#### Windows

1. Baixe o [Docker Desktop para Windows](https://www.docker.com/products/docker-desktop)
2. Execute o instalador
3. Reinicie o computador se solicitado
4. Verifique a instalação:

```bash
docker --version
docker ps
```

#### Linux (Ubuntu/Debian)

```bash
# Atualizar repositórios
sudo apt update

# Instalar dependências
sudo apt install -y apt-transport-https ca-certificates curl software-properties-common

# Adicionar chave GPG oficial do Docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Adicionar repositório Docker
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Instalar Docker
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io

# Adicionar usuário ao grupo docker (para não precisar de sudo)
sudo usermod -aG docker $USER

# Verificar instalação
docker --version
```

#### macOS

```bash
# Usando Homebrew
brew install --cask docker

# Ou baixe o Docker Desktop para Mac
# https://www.docker.com/products/docker-desktop
```

---

### 2. Instalando kubectl

O `kubectl` é a ferramenta de linha de comando para interagir com clusters Kubernetes.

#### Windows (PowerShell)

```powershell
# Usando Chocolatey
choco install kubernetes-cli

# Ou usando curl
curl.exe -LO "https://dl.k8s.io/release/v1.28.0/bin/windows/amd64/kubectl.exe"

# Verificar instalação
kubectl version --client
```

#### Linux

```bash
# Baixar a versão mais recente
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"

# Tornar executável
chmod +x kubectl

# Mover para PATH
sudo mv kubectl /usr/local/bin/

# Verificar instalação
kubectl version --client
```

#### macOS

```bash
# Usando Homebrew
brew install kubectl

# Verificar instalação
kubectl version --client
```

---

### 3. Instalando Kind (Kubernetes in Docker)

Kind permite criar clusters Kubernetes locais usando containers Docker.

#### Windows (PowerShell)

```powershell
# Usando Chocolatey
choco install kind

# Ou usando curl
curl.exe -Lo kind-windows-amd64.exe https://kind.sigs.k8s.io/dl/v0.20.0/kind-windows-amd64
Move-Item .\kind-windows-amd64.exe c:\some-dir-in-your-PATH\kind.exe

# Verificar instalação
kind version
```

#### Linux

```bash
# Baixar Kind
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64

# Tornar executável
chmod +x ./kind

# Mover para PATH
sudo mv ./kind /usr/local/bin/kind

# Verificar instalação
kind version
```

#### macOS

```bash
# Usando Homebrew
brew install kind

# Verificar instalação
kind version
```

---

### 4. Criando seu Primeiro Cluster com Kind

#### Cluster Simples (1 node)

```bash
# Criar cluster simples
kind create cluster

# Verificar cluster
kubectl cluster-info --context kind-kind
kubectl get nodes
```

#### Cluster Multi-Node (Recomendado)

Use o arquivo [kind.yaml](kind.yaml) deste projeto:

```yaml
# kind.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4

nodes:
  - role: control-plane
  - role: worker
  - role: worker
  - role: worker
```

```bash
# Criar cluster usando o arquivo de configuração
kind create cluster --name kubkindteste --config=kind.yaml

# Verificar nodes (deve mostrar 1 control-plane + 3 workers)
kubectl get nodes

# Exemplo de saída:
# NAME                         STATUS   ROLES           AGE   VERSION
# kubkindteste-control-plane   Ready    control-plane   1m    v1.27.0
# kubkindteste-worker          Ready    <none>          1m    v1.27.0
# kubkindteste-worker2         Ready    <none>          1m    v1.27.0
# kubkindteste-worker3         Ready    <none>          1m    v1.27.0
```

#### Comandos Úteis do Kind

```bash
# Listar clusters
kind get clusters

# Ver nodes de um cluster específico
kind get nodes --name kubkindteste

# Deletar cluster
kind delete cluster --name kubkindteste

# Carregar imagem Docker no Kind (IMPORTANTE!)
kind load docker-image pedrovieira249/golang-kubernetes-teste:v5.6 --name kubkindteste

# Exportar logs do cluster para debug
kind export logs --name kubkindteste
```

---

### 5. Configurando Metrics Server

O Metrics Server é necessário para o HPA (Horizontal Pod Autoscaler) funcionar.

#### Baixar e Instalar

```bash
# Baixar o arquivo de componentes
curl -LO https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Renomear para metrics-server.yaml
mv components.yaml metrics-server.yaml
```

#### Modificar para Kind (Importante!)

Como estamos usando Kind, precisamos adicionar `--kubelet-insecure-tls` ao Metrics Server.

Edite o arquivo `metrics-server.yaml` e adicione na seção `args` do container:

```yaml
containers:
- args:
  - --cert-dir=/tmp
  - --secure-port=10250
  - --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname
  - --kubelet-use-node-status-port
  - --metric-resolution=15s
  - --kubelet-insecure-tls  # <--- ADICIONAR ESTA LINHA
```

#### Aplicar o Metrics Server

```bash
# Aplicar a configuração
kubectl apply -f metrics-server.yaml

# Verificar se está rodando
kubectl get deployment metrics-server -n kube-system

# Verificar API de métricas (pode levar 1-2 minutos)
kubectl get apiservices | grep metrics

# Deve aparecer:
# v1beta1.metrics.k8s.io   kube-system/metrics-server   True

# Testar métricas
kubectl top nodes
kubectl top pods
```

---

### 6. Instalando Cert-Manager (Gerenciamento de Certificados TLS)

O **cert-manager** é um controlador Kubernetes que automatiza o gerenciamento e emissão de certificados TLS de várias fontes (Let's Encrypt, HashiCorp Vault, Venafi, etc).

#### Por que usar Cert-Manager?

- ✅ Automatiza criação e renovação de certificados TLS/SSL
- ✅ Integração com Let's Encrypt (certificados gratuitos)
- ✅ Suporte a múltiplos issuers (CA providers)
- ✅ Renovação automática antes da expiração
- ✅ Essencial para Ingress com HTTPS

---

#### Instalação do Cert-Manager

```bash
# Aplicar os manifestos oficiais do cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.19.2/cert-manager.yaml

# Exemplo de saída:
# namespace/cert-manager created
# customresourcedefinition.apiextensions.k8s.io/certificaterequests.cert-manager.io created
# customresourcedefinition.apiextensions.k8s.io/certificates.cert-manager.io created
# customresourcedefinition.apiextensions.k8s.io/challenges.acme.cert-manager.io created
# customresourcedefinition.apiextensions.k8s.io/clusterissuers.cert-manager.io created
# customresourcedefinition.apiextensions.k8s.io/issuers.cert-manager.io created
# customresourcedefinition.apiextensions.k8s.io/orders.acme.cert-manager.io created
# serviceaccount/cert-manager created
# ... (muitos recursos criados)
```

**O que foi criado:**
- Namespace `cert-manager`
- CRDs (Custom Resource Definitions) para certificados
- Deployments: cert-manager, cert-manager-cainjector, cert-manager-webhook
- RBAC roles e bindings
- Services e webhooks

---

#### Verificar Instalação

```bash
# Verificar se o namespace foi criado
kubectl get namespaces | grep cert-manager

# Ver pods do cert-manager (aguardar todos ficarem Running)
kubectl get pods --namespace cert-manager

# Exemplo de saída esperada:
# NAME                                       READY   STATUS    RESTARTS   AGE
# cert-manager-7d4b5d746-xxxxx              1/1     Running   0          1m
# cert-manager-cainjector-6d8d7c8f9-xxxxx   1/1     Running   0          1m
# cert-manager-webhook-7b8c8d9f8-xxxxx      1/1     Running   0          1m

# Aguardar até todos os pods estarem prontos
kubectl wait --for=condition=ready pod --all -n cert-manager --timeout=300s

# Verificar deployments
kubectl get deployments -n cert-manager

# Verificar services
kubectl get svc -n cert-manager

# Verificar CRDs criados
kubectl get crd | grep cert-manager
```

---

#### Testar o Cert-Manager

##### Teste 1: Criar um Issuer de Teste (Self-Signed)

```bash
# Criar arquivo test-issuer.yaml
cat <<EOF > test-issuer.yaml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: test-selfsigned
  namespace: default
spec:
  selfSigned: {}
EOF

# Aplicar o issuer
kubectl apply -f test-issuer.yaml

# Verificar
kubectl get issuer
kubectl describe issuer test-selfsigned
```

##### Teste 2: Criar um Certificado de Teste

```bash
# Criar arquivo test-certificate.yaml
cat <<EOF > test-certificate.yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: test-certificate
  namespace: default
spec:
  secretName: test-certificate-tls
  issuerRef:
    name: test-selfsigned
    kind: Issuer
  dnsNames:
    - example.com
    - www.example.com
EOF

# Aplicar o certificado
kubectl apply -f test-certificate.yaml

# Verificar o certificado
kubectl get certificate

# Ver detalhes (deve mostrar Ready=True)
kubectl describe certificate test-certificate

# Ver o secret criado com o certificado
kubectl get secret test-certificate-tls

# Ver conteúdo do certificado
kubectl get secret test-certificate-tls -o yaml
```

##### Teste 3: Verificar Logs

```bash
# Ver logs do cert-manager
kubectl logs -n cert-manager deployment/cert-manager

# Ver eventos relacionados a certificados
kubectl get events --sort-by=.metadata.creationTimestamp | grep -i certificate

# Verificar CertificateRequest criado automaticamente
kubectl get certificaterequest
kubectl describe certificaterequest <request-name>
```

---

#### Exemplo Prático: ClusterIssuer com Let's Encrypt

Para usar em produção com Let's Encrypt:

```yaml
# letsencrypt-issuer.yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    # Email para notificações de expiração
    email: seu-email@exemplo.com
    # Servidor de produção do Let's Encrypt
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      # Secret para armazenar a chave privada da conta ACME
      name: letsencrypt-prod
    # Solver HTTP01 para validação de domínio
    solvers:
    - http01:
        ingress:
          class: nginx
```

```bash
# Aplicar o ClusterIssuer
kubectl apply -f letsencrypt-issuer.yaml

# Verificar
kubectl get clusterissuer
kubectl describe clusterissuer letsencrypt-prod
```

**Usando com Ingress:**

```yaml
# ingress-with-tls.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
  - hosts:
    - exemplo.com
    - www.exemplo.com
    secretName: exemplo-com-tls  # Cert-manager criará este secret
  rules:
  - host: exemplo.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: go-server-service
            port:
              number: 80
```

---

#### Remover Cert-Manager

Quando quiser remover completamente o cert-manager:

```bash
# 1. Remover certificados e issuers criados (opcional, mas recomendado)
kubectl delete certificate --all --all-namespaces
kubectl delete issuer --all --all-namespaces
kubectl delete clusterissuer --all

# 2. Remover os recursos de teste
kubectl delete -f test-certificate.yaml
kubectl delete -f test-issuer.yaml
rm test-certificate.yaml test-issuer.yaml

# 3. Remover o cert-manager
kubectl delete -f https://github.com/cert-manager/cert-manager/releases/download/v1.19.2/cert-manager.yaml

# Aguardar a remoção
kubectl get namespaces | grep cert-manager  # Não deve retornar nada

# 4. Verificar se CRDs foram removidos
kubectl get crd | grep cert-manager  # Não deve retornar nada

# 5. (Se necessário) Forçar remoção de CRDs manualmente
kubectl delete crd certificaterequests.cert-manager.io
kubectl delete crd certificates.cert-manager.io
kubectl delete crd challenges.acme.cert-manager.io
kubectl delete crd clusterissuers.cert-manager.io
kubectl delete crd issuers.cert-manager.io
kubectl delete crd orders.acme.cert-manager.io
```

---

#### Verificar Remoção Completa

```bash
# Verificar namespace
kubectl get namespace cert-manager
# Deve retornar: Error from server (NotFound)

# Verificar pods
kubectl get pods -n cert-manager
# Deve retornar: No resources found

# Verificar CRDs
kubectl get crd | grep cert-manager
# Não deve retornar nada

# Verificar webhooks
kubectl get validatingwebhookconfigurations | grep cert-manager
kubectl get mutatingwebhookconfigurations | grep cert-manager
# Não deve retornar nada
```

---

#### Troubleshooting Cert-Manager

##### Problema 1: Pods não ficam Running

```bash
# Ver status dos pods
kubectl get pods -n cert-manager

# Ver logs de um pod com problema
kubectl logs -n cert-manager <pod-name>

# Descrever pod para ver eventos
kubectl describe pod -n cert-manager <pod-name>

# Verificar se CRDs foram criados
kubectl get crd | grep cert-manager
```

##### Problema 2: Certificado não é emitido (fica Pending)

```bash
# Ver status do certificado
kubectl get certificate
kubectl describe certificate <certificate-name>

# Ver CertificateRequest criado
kubectl get certificaterequest
kubectl describe certificaterequest <request-name>

# Ver Challenge (para ACME/Let's Encrypt)
kubectl get challenge
kubectl describe challenge <challenge-name>

# Ver logs do cert-manager
kubectl logs -n cert-manager deployment/cert-manager -f
```

##### Problema 3: Webhook não funciona

```bash
# Verificar se webhook está rodando
kubectl get pods -n cert-manager | grep webhook

# Ver logs do webhook
kubectl logs -n cert-manager deployment/cert-manager-webhook

# Verificar configuração do webhook
kubectl get validatingwebhookconfigurations cert-manager-webhook -o yaml
```

---

#### Comandos Úteis do Cert-Manager

```bash
# Listar todos os certificados do cluster
kubectl get certificates --all-namespaces

# Listar todos os issuers
kubectl get issuer --all-namespaces
kubectl get clusterissuer

# Ver secrets de certificados
kubectl get secrets | grep tls

# Forçar renovação de um certificado
kubectl delete certificaterequest <request-name>

# Ver status de um certificado
kubectl get certificate <certificate-name> -o wide

# Ver informações detalhadas do certificado
kubectl describe certificate <certificate-name>

# Ver o secret TLS criado
kubectl get secret <secret-name> -o yaml

# Decodificar certificado TLS do secret
kubectl get secret <secret-name> -o jsonpath='{.data.tls\.crt}' | base64 -d | openssl x509 -text -noout
```

---

#### Versões e Atualizações

```bash
# Verificar versão instalada do cert-manager
kubectl get deployment -n cert-manager cert-manager -o jsonpath='{.spec.template.spec.containers[0].image}'

# Listar versões disponíveis
# Acesse: https://github.com/cert-manager/cert-manager/releases

# Atualizar para nova versão
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.20.0/cert-manager.yaml

# Verificar atualização
kubectl rollout status deployment cert-manager -n cert-manager
```

---

### 7. Estrutura do Projeto

Este projeto contém os seguintes arquivos de configuração Kubernetes:

```
c:\Full Cycle\kubernetes\
├── golang/                      # Aplicação Go
│   ├── Dockerfile              # Build da imagem Docker
│   ├── server.go               # Servidor web Go
│   └── go.mod                  # Dependências Go
│
├── kind.yaml                   # Configuração do cluster Kind (3 workers)
├── metrics-server.yaml         # Servidor de métricas para HPA
│
├── pod.yaml                    # Definição básica de Pod
├── replicaset.yaml            # ReplicaSet com 5 réplicas
├── deployment.yaml            # Deployment principal da aplicação
├── service.yaml               # Service LoadBalancer
│
├── configmap-env.yaml         # ConfigMap com variáveis de ambiente
├── configmap-config.yaml      # ConfigMap com arquivo de configuração
├── secret.yaml                # Secrets (USERNAME/PASSWORD)
│
├── pv.yaml                    # PersistentVolume
├── pvc.yaml                   # PersistentVolumeClaim
│
├── hpa.yaml                   # Horizontal Pod Autoscaler
│
├── statefulset.yaml           # StatefulSet do MySQL
├── mysql-service-h.yaml       # Serviço Headless para MySQL
│
└── README-CONCEITOS-K8N.md    # Este arquivo

# Recursos Externos (aplicados via URL):
# - cert-manager.yaml           # Gerenciamento de certificados TLS (v1.19.2)
```

---

### 8. Workflow Completo de Deploy

#### Passo 1: Preparar a Aplicação

```bash
# Navegar para o diretório do projeto
cd "c:\Full Cycle\kubernetes"

# Build da imagem Docker
cd golang
docker build -t pedrovieira249/golang-kubernetes-teste:v5.6 .

# (Opcional) Push para Docker Hub
docker push pedrovieira249/golang-kubernetes-teste:v5.6

# Voltar para o diretório raiz
cd ..
```

#### Passo 2: Criar o Cluster Kind

```bash
# Criar cluster multi-node
kind create cluster --name kubkindteste --config=kind.yaml

# Carregar a imagem no Kind (IMPORTANTE - evita pull do Docker Hub)
kind load docker-image pedrovieira249/golang-kubernetes-teste:v5.6 --name kubkindteste

# Verificar nodes
kubectl get nodes
```

#### Passo 3: Instalar Metrics Server

```bash
# Aplicar metrics server
kubectl apply -f metrics-server.yaml

# Aguardar ficar pronto (1-2 minutos)
kubectl wait --for=condition=ready pod -l k8s-app=metrics-server -n kube-system --timeout=300s

# Verificar
kubectl get apiservices | grep metrics
```

#### Passo 3.5: Instalar Cert-Manager (Opcional)

```bash
# Instalar cert-manager (se precisar gerenciar certificados TLS)
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.19.2/cert-manager.yaml

# Aguardar pods ficarem prontos (1-2 minutos)
kubectl wait --for=condition=ready pod --all -n cert-manager --timeout=300s

# Verificar instalação
kubectl get pods -n cert-manager
kubectl get crd | grep cert-manager
```

#### Passo 4: Aplicar Configurações (ConfigMaps e Secrets)

```bash
# Aplicar ConfigMaps
kubectl apply -f configmap-env.yaml
kubectl apply -f configmap-config.yaml

# Aplicar Secrets
kubectl apply -f secret.yaml

# Verificar
kubectl get configmaps
kubectl get secrets
```

#### Passo 5: Aplicar Storage (PV e PVC)

```bash
# Aplicar PersistentVolume e PersistentVolumeClaim
kubectl apply -f pv.yaml
kubectl apply -f pvc.yaml

# Verificar
kubectl get pv
kubectl get pvc
```

#### Passo 6: Deploy da Aplicação

```bash
# Aplicar Deployment
kubectl apply -f deployment.yaml

# Verificar pods (aguardar todos ficarem Running)
kubectl get pods -w

# Verificar logs de um pod
kubectl logs -l app=go-server -f
```

#### Passo 7: Expor com Service

```bash
# Aplicar Service
kubectl apply -f service.yaml

# Verificar service
kubectl get svc

# No Kind, LoadBalancer não recebe IP externo automaticamente
# Use port-forward para acessar:
kubectl port-forward service/go-server-service 8080:80
```

#### Passo 8: Configurar HPA (Auto-scaling)

```bash
# Aplicar HPA
kubectl apply -f hpa.yaml

# Verificar HPA
kubectl get hpa

# Monitorar em tempo real
kubectl get hpa -w
```

#### Passo 9: Acessar a Aplicação

```bash
# Em um terminal, fazer port-forward
kubectl port-forward service/go-server-service 8080:80

# Em outro terminal ou navegador, acessar:
# http://localhost:8080          → Página principal (com variáveis de ambiente)
# http://localhost:8080/healthz  → Health check
# http://localhost:8080/config   → Conteúdo do ConfigMap
# http://localhost:8080/secret   → Secrets (USERNAME/PASSWORD)
```

---

### 9. Deploy do MySQL (StatefulSet)

```bash
# Aplicar serviço headless
kubectl apply -f mysql-service-h.yaml

# Aplicar StatefulSet
kubectl apply -f statefulset.yaml

# Verificar pods do MySQL (devem ser criados sequencialmente)
kubectl get pods -l app=mysql -w

# Verificar PVCs criados automaticamente pelo StatefulSet
kubectl get pvc

# Conectar ao MySQL para testar
kubectl exec -it mysql-0 -- mysql -uroot -proot

# Dentro do MySQL:
# SHOW DATABASES;
# CREATE DATABASE teste;
# EXIT;
```

---

### 10. Testando Auto-Scaling com Fortio

```bash
# Terminal 1: Monitorar HPA
kubectl get hpa -w

# Terminal 2: Monitorar Pods
kubectl get pods -l app=go-server -w

# Terminal 3: Executar teste de carga
kubectl run fortio --rm -it --image=fortio/fortio -- \
  load -c 300 -qps 0 -t 1m http://go-server-service/

# Observar:
# - HPA aumentando número de réplicas
# - Novos pods sendo criados
# - CPU usage aumentando
```

---

### 11. Comandos de Troubleshooting

```bash
# Ver eventos do cluster
kubectl get events --sort-by=.metadata.creationTimestamp

# Ver logs de um deployment
kubectl logs -l app=go-server --tail=100 -f

# Ver logs de todos os containers de um pod
kubectl logs pod-name --all-containers=true

# Descrever recurso (para debug)
kubectl describe pod <pod-name>
kubectl describe deployment go-server
kubectl describe hpa go-server-hpa

# Ver uso de recursos
kubectl top nodes
kubectl top pods

# Entrar em um pod para debug
kubectl exec -it <pod-name> -- sh

# Ver configuração completa de um recurso
kubectl get deployment go-server -o yaml
kubectl get service go-server-service -o yaml

# Verificar endpoints de um service
kubectl get endpoints go-server-service

# Ver TODAS as resources
kubectl get all
kubectl get all -A  # Todos os namespaces
```

---

### 12. Limpeza e Manutenção

```bash
# Deletar recursos específicos
kubectl delete -f deployment.yaml
kubectl delete -f service.yaml
kubectl delete -f hpa.yaml

# Deletar tudo de uma vez
kubectl delete -f .

# Deletar por label
kubectl delete pods -l app=go-server
kubectl delete all -l app=go-server

# Deletar o cluster Kind
kind delete cluster --name kubkindteste

# Limpar imagens Docker não utilizadas
docker system prune -a
```

---

### 13. Configurações Importantes

#### Arquivo kind.yaml Explicado

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4

nodes:
  - role: control-plane  # Node mestre (gerencia o cluster)
  - role: worker         # Node worker 1 (executa aplicações)
  - role: worker         # Node worker 2
  - role: worker         # Node worker 3
```

**Por que 3 workers?**
- Permite testar distribuição de pods
- Simula ambiente de produção
- Testa alta disponibilidade
- Facilita ver o scheduler em ação

#### Recursos da Aplicação Go

No [deployment.yaml](deployment.yaml), temos:

```yaml
resources:
  requests:      # Mínimo garantido ao pod
    memory: "32Mi"
    cpu: "10m"   # 10 millicores = 0.01 CPU
  limits:        # Máximo que o pod pode usar
    memory: "64Mi"
    cpu: "50m"   # 50 millicores = 0.05 CPU
```

**Unidades:**
- CPU: `m` = millicores (1000m = 1 CPU)
- Memory: `Mi` = Mebibytes, `Gi` = Gibibytes

#### Health Checks Configurados

```yaml
startupProbe:      # Verifica se iniciou corretamente
  httpGet:
    path: /healthz
    port: 8080
  periodSeconds: 3
  failureThreshold: 30  # 30 tentativas = 90 segundos max

readinessProbe:    # Verifica se está pronto para receber tráfego
  httpGet:
    path: /healthz
    port: 8080
  periodSeconds: 2
  failureThreshold: 1

livenessProbe:     # Verifica se está vivo (reinicia se falhar)
  httpGet:
    path: /healthz
    port: 8080
  periodSeconds: 5
  failureThreshold: 3  # 3 falhas = reinicia pod
```

**Lógica do /healthz:**
- Primeiros 10 segundos: retorna erro 500 (simula inicialização)
- Após 10 segundos: retorna 200 OK

Isso demonstra como o `startupProbe` aguarda a aplicação iniciar antes de executar outros probes.

---

## 🎯 O que é Kubernetes?

**Kubernetes (K8s)** é uma plataforma open-source para automação de deploy, escalonamento e gerenciamento de aplicações em contêineres. Foi criado pelo Google e agora é mantido pela Cloud Native Computing Foundation (CNCF).

### Analogia do Mundo Real

Imagine que você gerencia um restaurante (seu sistema):

- **Sem Kubernetes:** Você tem que contratar garçons manualmente, ajustar a quantidade conforme o movimento, substituir quem falta, distribuir as mesas...
- **Com Kubernetes:** Você tem um gerente automático que contrata, demite, redistribui tarefas, substitui funcionários ausentes, e garante que o restaurante sempre funcione perfeitamente.

### O que o Kubernetes faz?

- ✅ **Orquestração de Contêineres:** Gerencia múltiplos containers Docker
- ✅ **Auto-Scaling:** Aumenta ou diminui recursos automaticamente
- ✅ **Auto-Healing:** Reinicia containers que falharam
- ✅ **Load Balancing:** Distribui tráfego entre múltiplas instâncias
- ✅ **Rolling Updates:** Atualiza aplicações sem downtime
- ✅ **Service Discovery:** Conecta serviços automaticamente
- ✅ **Gerenciamento de Configuração:** Centraliza secrets e configs

---

## 🤔 Por que usar Kubernetes?

### Sem Kubernetes (Docker simples)

```bash
# Você precisa manualmente:
docker run -d app1
docker run -d app2
docker run -d app3

# Se um container cair, você precisa reiniciar manualmente
docker restart app1

# Para escalar, você precisa criar mais manualmente
docker run -d app1
docker run -d app1
docker run -d app1

# E configurar load balancer manualmente...
```

### Com Kubernetes

```bash
# Você declara o que quer:
kubectl apply -f deployment.yaml

# Kubernetes garante automaticamente:
# - 3 réplicas sempre rodando
# - Se uma cair, sobe outra automaticamente
# - Load balancer configurado
# - Health checks automáticos
# - Updates sem downtime
```

### Benefícios

| Problema | Solução com Kubernetes |
|----------|------------------------|
| Container caiu | Kubernetes reinicia automaticamente |
| Muito tráfego | Escala automaticamente mais pods |
| Deploy de nova versão | Rolling update sem downtime |
| Configuração espalhada | ConfigMaps e Secrets centralizados |
| Múltiplos servidores | Distribui carga entre nodes automaticamente |
| Rede complexa | Service Discovery automático |

---

## 🏗️ Arquitetura do Kubernetes

O Kubernetes é composto por dois tipos principais de componentes:

```
┌─────────────────────────────────────────────────────────────┐
│                    CLUSTER KUBERNETES                        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │           CONTROL PLANE (Cérebro)                   │   │
│  │  ┌──────────────┐  ┌──────────────┐                │   │
│  │  │ API Server   │  │  Scheduler   │                │   │
│  │  └──────────────┘  └──────────────┘                │   │
│  │  ┌──────────────┐  ┌──────────────┐                │   │
│  │  │   etcd       │  │  Controller  │                │   │
│  │  │  (Database)  │  │   Manager    │                │   │
│  │  └──────────────┘  └──────────────┘                │   │
│  └─────────────────────────────────────────────────────┘   │
│                          │                                   │
│                          │ (Gerencia)                        │
│                          ▼                                   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              WORKER NODES (Trabalho)                │   │
│  │                                                      │   │
│  │  ┌───────────┐  ┌───────────┐  ┌───────────┐      │   │
│  │  │  Node 1   │  │  Node 2   │  │  Node 3   │      │   │
│  │  │           │  │           │  │           │      │   │
│  │  │ ┌───────┐ │  │ ┌───────┐ │  │ ┌───────┐ │      │   │
│  │  │ │ Pod 1 │ │  │ │ Pod 3 │ │  │ │ Pod 5 │ │      │   │
│  │  │ └───────┘ │  │ └───────┘ │  │ └───────┘ │      │   │
│  │  │ ┌───────┐ │  │ ┌───────┐ │  │ ┌───────┐ │      │   │
│  │  │ │ Pod 2 │ │  │ │ Pod 4 │ │  │ │ Pod 6 │ │      │   │
│  │  │ └───────┘ │  │ └───────┘ │  │ └───────┘ │      │   │
│  │  └───────────┘  └───────────┘  └───────────┘      │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## 🧩 Conceitos Fundamentais

### 1. Cluster (Aglomerado)

**O que é:** Um conjunto de máquinas (físicas ou virtuais) trabalhando juntas como um sistema único.

**Analogia:** Um cluster é como uma empresa:
- Tem um **escritório central** (Control Plane) que toma decisões
- Tem **filiais/lojas** (Nodes) que fazem o trabalho real

**Componentes:**

```
Cluster = Control Plane + Worker Nodes
```

**Exemplo:**
```bash
# Ver informações do cluster
kubectl cluster-info

# Ver todos os nodes do cluster
kubectl get nodes
```

---

### 2. Control Plane (Plano de Controle)

**O que é:** O "cérebro" do cluster. Toma todas as decisões sobre onde rodar os pods, monitora o estado, e gerencia o cluster.

**Componentes principais:**

#### a) **API Server** (kube-apiserver)
- **O que faz:** Ponto central de comunicação. Toda interação com o cluster passa por aqui.
- **Analogia:** É a recepcionista da empresa. Tudo passa por ela.
- **Exemplo:** Quando você executa `kubectl get pods`, você está falando com o API Server.

#### b) **etcd**
- **O que faz:** Banco de dados que armazena TODO o estado do cluster.
- **Analogia:** É o arquivo/database da empresa. Guarda tudo.
- **Armazena:** Configurações, secrets, estado atual de todos os recursos.

#### c) **Scheduler** (kube-scheduler)
- **O que faz:** Decide em qual Node um Pod deve rodar.
- **Analogia:** É o RH que decide em qual filial colocar um novo funcionário.
- **Considera:** Recursos disponíveis (CPU, memória), restrições, afinidade.

#### d) **Controller Manager** (kube-controller-manager)
- **O que faz:** Garante que o estado atual seja igual ao estado desejado.
- **Analogia:** É o gerente que verifica se tudo está como deveria estar.
- **Exemplos de Controllers:**
  - **ReplicaSet Controller:** Garante número correto de réplicas
  - **Node Controller:** Monitora nodes com problemas
  - **Endpoint Controller:** Conecta Services e Pods

**Resumo:**
```
┌─────────────────────────────────────────┐
│         CONTROL PLANE                   │
├─────────────────────────────────────────┤
│ API Server    → Interface de comunicação│
│ etcd          → Banco de dados          │
│ Scheduler     → Decide onde rodar Pods  │
│ Controller    → Garante estado desejado │
└─────────────────────────────────────────┘
```

---

### 3. Node (Nó)

**O que é:** Uma máquina (física ou virtual) que executa as aplicações. Cada node pode rodar múltiplos pods.

**Analogia:** É como uma loja/filial da empresa. Lá é onde o trabalho real acontece.

**Componentes de cada Node:**

#### a) **kubelet**
- **O que faz:** Agente que roda em cada node. Garante que os containers estejam rodando.
- **Analogia:** É o gerente da filial que executa as ordens do escritório central.
- **Responsabilidades:**
  - Recebe instruções do Control Plane
  - Inicia/para containers
  - Reporta status dos pods

#### b) **kube-proxy**
- **O que faz:** Gerencia regras de rede no node. Permite comunicação entre pods.
- **Analogia:** É o telefonista da filial que roteia ligações.
- **Função:** Load balancing de rede, encaminha tráfego para pods corretos

#### c) **Container Runtime**
- **O que faz:** Software que roda os containers (Docker, containerd, CRI-O).
- **Analogia:** São as ferramentas que os funcionários usam para trabalhar.

**Exemplo:**
```bash
# Ver todos os nodes
kubectl get nodes

# Ver detalhes de um node específico
kubectl describe node nome-do-node

# Ver uso de recursos dos nodes
kubectl top nodes
```

**Estrutura de um Node:**
```
┌────────────────────────────────────┐
│           NODE                     │
├────────────────────────────────────┤
│  kubelet (gerente)                 │
│  kube-proxy (rede)                 │
│  Container Runtime (Docker)        │
│                                    │
│  ┌──────────┐  ┌──────────┐      │
│  │  Pod 1   │  │  Pod 2   │      │
│  │ ┌──────┐ │  │ ┌──────┐ │      │
│  │ │App A │ │  │ │App B │ │      │
│  │ └──────┘ │  │ └──────┘ │      │
│  └──────────┘  └──────────┘      │
└────────────────────────────────────┘
```

---

### 4. Pod

**O que é:** A menor unidade que você pode criar no Kubernetes. Um pod encapsula um ou mais containers que compartilham recursos.

**Analogia:** Um pod é como uma "sala de trabalho" onde um ou mais funcionários trabalham juntos compartilhando mesa, internet, telefone.

**Características:**

- **IP único:** Cada pod tem seu próprio endereço IP
- **Compartilham rede:** Containers no mesmo pod se comunicam via localhost
- **Compartilham volumes:** Podem acessar os mesmos arquivos
- **Efêmero:** Pods são descartáveis, podem ser substituídos a qualquer momento
- **Co-localização:** Containers no mesmo pod sempre rodam no mesmo node

**Tipos de Pods:**

#### Pod com 1 container (Mais comum - 95% dos casos)
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
spec:
  containers:
  - name: nginx
    image: nginx:latest
    ports:
    - containerPort: 80
```

```
┌─────────────────┐
│      Pod        │
│  ┌───────────┐  │
│  │  Nginx    │  │
│  │ Container │  │
│  └───────────┘  │
└─────────────────┘
```

#### Pod com múltiplos containers (Padrão Sidecar)
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-with-logs
spec:
  containers:
  - name: app
    image: minha-app:latest
  - name: log-collector  # Container auxiliar
    image: fluentd:latest
```

```
┌──────────────────────────────┐
│          Pod                 │
│  ┌───────────┐ ┌──────────┐ │
│  │   App     │ │  Logs    │ │
│  │ Container │ │Collector │ │
│  └───────────┘ └──────────┘ │
│         (compartilham        │
│          rede e volumes)     │
└──────────────────────────────┘
```

**Exemplo prático:**
```bash
# Criar um pod simples
kubectl run nginx --image=nginx

# Ver pods
kubectl get pods

# Ver detalhes do pod
kubectl describe pod nginx

# Ver logs do pod
kubectl logs nginx

# Entrar no pod
kubectl exec -it nginx -- /bin/bash

# Deletar pod
kubectl delete pod nginx
```

**Estados do Pod:**

| Estado | Descrição |
|--------|-----------|
| **Pending** | Pod foi aceito mas containers ainda não foram criados |
| **Running** | Pod foi associado a um node e todos os containers estão rodando |
| **Succeeded** | Todos os containers terminaram com sucesso |
| **Failed** | Todos os containers terminaram e pelo menos um falhou |
| **Unknown** | Estado do pod não pode ser obtido |

---

### 5. Namespace

**O que é:** Uma forma de dividir recursos do cluster em grupos virtuais.

**Analogia:** São como "departamentos" de uma empresa. Marketing, TI, RH - cada um tem seus próprios recursos mas fazem parte da mesma empresa.

**Namespaces padrão:**

- **default:** Namespace padrão para recursos sem namespace especificado
- **kube-system:** Para recursos do próprio Kubernetes
- **kube-public:** Recursos públicos acessíveis por todos
- **kube-node-lease:** Informações de heartbeat dos nodes

**Por que usar Namespaces?**

1. **Organização:** Separar ambientes (dev, staging, prod)
2. **Isolamento:** Recursos de um namespace não afetam outros
3. **Segurança:** Controlar acesso via RBAC
4. **Quotas:** Limitar recursos por namespace

**Exemplo:**
```bash
# Listar namespaces
kubectl get namespaces

# Criar namespace
kubectl create namespace desenvolvimento
kubectl create namespace producao

# Ver pods de um namespace específico
kubectl get pods -n kube-system

# Ver pods de todos os namespaces
kubectl get pods --all-namespaces
# ou
kubectl get pods -A

# Definir namespace padrão para comandos
kubectl config set-context --current --namespace=desenvolvimento
```

**Estrutura:**
```
┌─────────────────────────────────────────┐
│              CLUSTER                    │
├─────────────────────────────────────────┤
│  ┌────────────────────────────────┐    │
│  │  Namespace: desenvolvimento    │    │
│  │  - Pod: app-dev-1              │    │
│  │  - Service: api-dev            │    │
│  └────────────────────────────────┘    │
│  ┌────────────────────────────────┐    │
│  │  Namespace: producao           │    │
│  │  - Pod: app-prod-1             │    │
│  │  - Service: api-prod           │    │
│  └────────────────────────────────┘    │
└─────────────────────────────────────────┘
```

---

## 📦 Objetos do Kubernetes

### 6. Deployment

**O que é:** A forma recomendada de fazer deploy de aplicações. Gerencia a criação e atualização de Pods.

**Analogia:** É como uma "receita" que diz: "Quero sempre 3 cópias da minha aplicação rodando, usando esta imagem Docker".

**O que o Deployment faz:**

- ✅ Cria e gerencia ReplicaSets
- ✅ Garante número de réplicas desejado
- ✅ Rolling updates (atualiza sem downtime)
- ✅ Rollback (volta versão anterior se der problema)
- ✅ Scale up/down (aumenta/diminui réplicas)

**Estrutura:**
```
Deployment
    ↓
ReplicaSet
    ↓
Pods (3 réplicas)
```

**Exemplo:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3  # Quero 3 cópias
  selector:
    matchLabels:
      app: nginx
  template:  # Template do Pod
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.21
        ports:
        - containerPort: 80
```

**Comandos:**
```bash
# Criar deployment
kubectl apply -f deployment.yaml

# Ver deployments
kubectl get deployments

# Ver detalhes
kubectl describe deployment deployment/go-server

# Escalar (mudar número de réplicas)
kubectl scale deployment deployment/go-server --replicas=5

# Atualizar imagem (rolling update)
kubectl set image deployment/nginx-deployment golangserver=pedrovieira249/golang-kubernetes-teste:latest

# Ver histórico de versões
kubectl rollout history deployment/go-server

# Voltar para versão anterior (rollback)
kubectl rollout undo deployment/go-server

# Ver status do rollout
kubectl rollout status deployment/go-server
```

**Fluxo de um Deployment:**
```
1. Você cria um Deployment
   └─> Deployment cria um ReplicaSet
       └─> ReplicaSet cria 3 Pods
           └─> Cada Pod roda 1 container nginx

2. Um Pod morre
   └─> ReplicaSet detecta
       └─> Cria novo Pod automaticamente
           └─> Sempre mantém 3 rodando

3. Você atualiza a imagem
   └─> Deployment cria novo ReplicaSet
       └─> Gradualmente substitui Pods antigos por novos
           └─> Zero downtime!
```

---

### 7. ReplicaSet

**O que é:** Garante que um número específico de réplicas de pods esteja sempre rodando.

**Analogia:** É como um supervisor que garante que sempre tenha exatamente X funcionários trabalhando. Se um sai, ele contrata outro imediatamente.

**Importante:** Você geralmente NÃO cria ReplicaSets diretamente. Eles são criados automaticamente pelos Deployments.

**Exemplo:**
```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: nginx-replicaset
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
```

**Comandos:**
```bash
# Ver replicasets
kubectl get replicasets
kubectl get rs  # forma curta

# Ver detalhes
kubectl describe rs nginx-replicaset

# Deletar (cuidado: deleta os pods também)
kubectl delete rs nginx-replicaset
```

---

### 8. Service (Serviço)

**O que é:** Expõe seus Pods na rede, fornecendo um endereço estável para acessá-los.

**Por que preciso?** Pods são efêmeros (morrem e nascem). Seus IPs mudam. Services fornecem um ponto de acesso estável.

**Analogia:** É como o número de telefone da empresa. Mesmo que funcionários mudem, o número continua o mesmo.

**Tipos de Services:**

#### a) **ClusterIP** (Padrão)
- **Uso:** Comunicação interna no cluster
- **Acesso:** Somente dentro do cluster
- **Quando usar:** Microsserviços conversando entre si

```yaml
apiVersion: v1
kind: Service
metadata:
  name: backend-service
spec:
  type: ClusterIP
  selector:
    app: backend
  ports:
  - port: 80  # Porta do service
    targetPort: 8080  # Porta do container
```

#### b) **NodePort**
- **Uso:** Expõe o service em uma porta de cada Node
- **Acesso:** De fora do cluster via <NodeIP>:<NodePort>
- **Quando usar:** Desenvolvimento, testes

```yaml
apiVersion: v1
kind: Service
metadata:
  name: frontend-service
spec:
  type: NodePort
  selector:
    app: frontend
  ports:
  - port: 80
    targetPort: 8080
    nodePort: 30080  # Porta acessível externamente
```

#### c) **LoadBalancer**
- **Uso:** Cria load balancer externo (na cloud)
- **Acesso:** Via IP público do load balancer
- **Quando usar:** Produção em cloud (AWS, GCP, Azure)

```yaml
apiVersion: v1
kind: Service
metadata:
  name: api-service
spec:
  type: LoadBalancer
  selector:
    app: api
  ports:
  - port: 80
    targetPort: 8080
```

**Como funciona:**
```
┌─────────────────────────────────────────────────┐
│                 SERVICE                         │
│           (IP estável: 10.0.0.100)              │
│                                                  │
│   Distribui tráfego automaticamente para:       │
│                                                  │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│   │  Pod 1   │  │  Pod 2   │  │  Pod 3   │    │
│   │ IP: ...1 │  │ IP: ...2 │  │ IP: ...3 │    │
│   └──────────┘  └──────────┘  └──────────┘    │
│                                                  │
│   (Pods podem morrer e renascer com novos IPs   │
│    mas o Service continua com mesmo IP)         │
└─────────────────────────────────────────────────┘
```

**Comandos:**
```bash
# Ver services
kubectl get services
kubectl get svc  # forma curta

# Ver detalhes
kubectl describe svc frontend-service

# Testar acesso ao service
kubectl port-forward svc/go-server-service 8080:80
kubectl port-forward service/go-server-service 8080:80

# Ver endpoints do service (pods que estão atendendo)
kubectl get endpoints frontend-service
```

---

### 9. ConfigMap

**O que é:** Armazena configurações não-sensíveis em formato chave-valor.

**Analogia:** É como um arquivo de configuração centralizado que todos podem usar.

**Por que usar:** Separar configuração do código. Mesma imagem Docker funciona em dev/staging/prod com ConfigMaps diferentes.

**Exemplo:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  database_url: "mysql://db:3306/myapp"
  cache_enabled: "true"
  log_level: "info"
  app.properties: |
    color=blue
    language=pt-BR
```

**Usando ConfigMap em um Pod:**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-pod
spec:
  containers:
  - name: app
    image: minha-app:latest
    env:
    # Injetar como variável de ambiente
    - name: DATABASE_URL
      valueFrom:
        configMapKeyRef:
          name: app-config
          key: database_url
    volumeMounts:
    # Ou montar como arquivo
    - name: config
      mountPath: /etc/config
  volumes:
  - name: config
    configMap:
      name: app-config
```

**Comandos:**
```bash
# Criar ConfigMap via comando
kubectl create configmap app-config --from-literal=color=blue

# Criar ConfigMap de arquivo
kubectl create configmap app-config --from-file=config.properties

# Ver ConfigMaps
kubectl get configmaps
kubectl get cm  # forma curta

# Ver conteúdo
kubectl describe cm app-config
kubectl get cm app-config -o yaml
```

---

### 10. Secret

**O que é:** Similar ao ConfigMap, mas para dados sensíveis (senhas, tokens, chaves).

**Diferença do ConfigMap:** Dados são codificados em base64 e podem ser criptografados em repouso.

**Exemplo:**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-secret
type: Opaque
data:
  # Valores em base64
  username: YWRtaW4=  # admin
  password: cGFzc3dvcmQxMjM=  # password123
```

**Usando Secret:**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-pod
spec:
  containers:
  - name: app
    image: minha-app:latest
    env:
    - name: DB_USERNAME
      valueFrom:
        secretKeyRef:
          name: db-secret
          key: username
    - name: DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: db-secret
          key: password
```

**Comandos:**
```bash
# Criar Secret via comando
kubectl create secret generic db-secret \
  --from-literal=username=admin \
  --from-literal=password=senha123

# Ver Secrets (valores são ocultos)
kubectl get secrets

# Ver detalhes
kubectl describe secret db-secret

# Ver valores decodificados (CUIDADO!)
kubectl get secret db-secret -o jsonpath='{.data.password}' | base64 -d
```

---

### 11. Volume

**O que é:** Armazenamento que persiste além do ciclo de vida de um Pod.

**Por que preciso:** Containers são efêmeros. Se o pod morre, dados são perdidos. Volumes permitem persistência.

**Tipos principais:**

#### a) **emptyDir** (temporário)
- Criado quando Pod inicia
- Deletado quando Pod é removido
- Compartilhado entre containers do mesmo pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod-com-volume
spec:
  containers:
  - name: app
    image: nginx
    volumeMounts:
    - name: cache
      mountPath: /cache
  volumes:
  - name: cache
    emptyDir: {}
```

#### b) **hostPath** (node local)
- Usa diretório do node
- Dados persistem se pod é recriado no mesmo node

```yaml
volumes:
- name: logs
  hostPath:
    path: /var/log/app
    type: DirectoryOrCreate
```

#### c) **PersistentVolume (PV)** e **PersistentVolumeClaim (PVC)**
- Armazenamento real e persistente
- Independente do ciclo de vida do pod

```yaml
# PersistentVolumeClaim
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pvc
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi

---
# Usando o PVC no Pod
apiVersion: v1
kind: Pod
metadata:
  name: mysql-pod
spec:
  containers:
  - name: mysql
    image: mysql:8
    volumeMounts:
    - name: mysql-storage
      mountPath: /var/lib/mysql
  volumes:
  - name: mysql-storage
    persistentVolumeClaim:
      claimName: mysql-pvc
```

---

### 12. Ingress

**O que é:** Gerencia acesso externo aos services, tipicamente HTTP/HTTPS.

**Analogia:** É como um porteiro/recepcionista que direciona visitantes para o departamento correto baseado no que eles pedem.

**Por que usar:** 
- Roteamento baseado em domínio/path
- SSL/TLS termination
- Load balancing
- Um único ponto de entrada

**Exemplo:**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress
spec:
  rules:
  - host: api.exemplo.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: api-service
            port:
              number: 80
  - host: blog.exemplo.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: blog-service
            port:
              number: 80
```

**Como funciona:**
```
                    ┌─────────────────┐
Internet ─────────> │     Ingress     │
                    │  (nginx/traefik) │
                    └─────────────────┘
                            │
         ┌──────────────────┼──────────────────┐
         ▼                  ▼                  ▼
    api.exemplo.com    blog.exemplo.com   app.exemplo.com
         │                  │                  │
    ┌────▼────┐        ┌────▼────┐       ┌────▼────┐
    │ Service │        │ Service │       │ Service │
    │   API   │        │  Blog   │       │   App   │
    └─────────┘        └─────────┘       └─────────┘
         │                  │                  │
    ┌────┴────┐        ┌────┴────┐       ┌────┴────┐
    │  Pods   │        │  Pods   │       │  Pods   │
    └─────────┘        └─────────┘       └─────────┘
```

---

## 🎯 Comandos Básicos do kubectl

### Informações do Cluster

```bash
# Ver informações do cluster
kubectl cluster-info

# Ver versão do kubectl e do cluster
kubectl version

# Ver nodes
kubectl get nodes

# Detalhes de um node
kubectl describe node <nome-do-node>
```

### Trabalhando com Pods

```bash
# Criar pod simples
kubectl run nginx --image=nginx

# Ver pods
kubectl get pods
kubectl get pods -o wide  # Mais informações (IP, Node)
kubectl get pods --watch  # Atualização em tempo real
kubectl get pods -o wide -w
kubectl get pods -w

# Ver pods de todos os namespaces
kubectl get pods -A

# Ver detalhes de um pod
kubectl describe pod <nome-do-pod>

# Ver logs
kubectl logs <nome-do-pod>
kubectl logs -f <nome-do-pod>  # Seguir logs em tempo real
kubectl logs <nome-do-pod> -c <nome-container>  # Pod com múltiplos containers

# Executar comando no pod
kubectl exec <nome-do-pod> -- ls /app
kubectl exec -it <nome-do-pod> -- /bin/bash  # Entrar no pod

# Deletar pod
kubectl delete pod <nome-do-pod>
```

### Trabalhando com Deployments

```bash
# Criar deployment
kubectl create deployment nginx --image=nginx

# Ver deployments
kubectl get deployments
kubectl get deploy  # forma curta

# Escalar deployment
kubectl scale deployment nginx --replicas=5

# Atualizar imagem (rolling update)
kubectl set image deployment/nginx nginx=nginx:1.22

# Ver histórico de rollouts
kubectl rollout history deployment/nginx

# Fazer rollback
kubectl rollout undo deployment/nginx

# Ver status do rollout
kubectl rollout status deployment/nginx

# Deletar deployment
kubectl delete deployment nginx
```

### Trabalhando com Services

```bash
# Criar service
kubectl expose deployment nginx --port=80 --type=NodePort

# Ver services
kubectl get services
kubectl get svc  # forma curta

# Ver detalhes
kubectl describe svc nginx

# Deletar service
kubectl delete svc nginx
```

### Aplicar Manifestos YAML

```bash
# Aplicar um arquivo
kubectl apply -f deployment.yaml

# Aplicar um arquivo e visualizar os pods (para validar o livenessProbe)
kubectl apply -f deployment.yaml && kubectl get pods --watch

# Aplicar múltiplos arquivos
kubectl apply -f deployment.yaml -f service.yaml

# Aplicar diretório inteiro
kubectl apply -f ./manifests/

# Ver o que seria aplicado (dry-run)
kubectl apply -f deployment.yaml --dry-run=client

# Deletar recursos de um arquivo
kubectl delete -f deployment.yaml
```

### Debugging

```bash
# Ver eventos do cluster
kubectl get events
kubectl get events --sort-by=.metadata.creationTimestamp

# Ver uso de recursos
kubectl top nodes
kubectl top pods

# Port-forward (acessar pod localmente)
kubectl port-forward pod/go-server 8080:80

# Port-forward de service
kubectl port-forward svc/go-server 8080:80

# Copiar arquivos de/para pod
kubectl cp <pod>:/path/to/file ./local-file
kubectl cp ./local-file <pod>:/path/to/file
```

### Namespace

```bash
# Ver namespaces
kubectl get namespaces
kubectl get ns  # forma curta

# Criar namespace
kubectl create namespace dev

# Usar namespace específico
kubectl get pods -n dev

# Definir namespace padrão
kubectl config set-context --current --namespace=dev

# Deletar namespace (CUIDADO: deleta tudo dentro dele)
kubectl delete namespace dev
```

### Outros Comandos Úteis

```bash
# Ver todos os recursos
kubectl get all

# Ver recursos específicos
kubectl get pods,services,deployments

# Editar recurso diretamente
kubectl edit deployment nginx

# Ver definição YAML de um recurso
kubectl get deployment nginx -o yaml

# Ver definição JSON
kubectl get deployment nginx -o json

# Usar JSONPath para extrair dados
kubectl get pods -o jsonpath='{.items[*].metadata.name}'

# Explicar recursos
kubectl explain pods
kubectl explain pods.spec.containers

# Ver API resources disponíveis
kubectl api-resources
```

---

## � Arquivos do Projeto

Esta seção detalha todos os arquivos YAML do projeto e sua finalidade.

---

### 1. Aplicação Go ([golang/](golang/))

#### [server.go](golang/server.go)

Servidor web Go com 4 endpoints:

```go
// Endpoint principal - exibe variáveis de ambiente
http://localhost:8080/

// Health check - usado pelos probes do Kubernetes
http://localhost:8080/healthz

// Exibe conteúdo do ConfigMap montado como arquivo
http://localhost:8080/config

// Exibe secrets (USERNAME e PASSWORD)
http://localhost:8080/secret
```

**Funcionalidades:**
- Lê variáveis de ambiente do ConfigMap
- Lê arquivo de configuração montado via volume
- Exibe secrets para demonstração
- Health check inteligente (falha nos primeiros 10s, depois OK)

#### [Dockerfile](golang/Dockerfile)

```dockerfile
FROM golang:1.25-alpine
WORKDIR /app
COPY go.mod ./
COPY server.go ./
RUN go build -o server server.go
EXPOSE 8080
CMD ["./server"]
```

**Build e Push:**
```bash
cd golang
docker build -t pedrovieira249/golang-kubernetes-teste:v5.6 .
docker push pedrovieira249/golang-kubernetes-teste:v5.6
kind load docker-image pedrovieira249/golang-kubernetes-teste:v5.6 --name kubkindteste
```

---

### 2. Configuração do Cluster

#### [kind.yaml](kind.yaml)

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
  - role: worker
  - role: worker
  - role: worker
```

**Uso:**
```bash
kind create cluster --name kubkindteste --config=kind.yaml
```

**Resultado:** 1 control-plane + 3 workers = 4 nodes totais

---

### 3. Pods e Réplicas

#### [pod.yaml](pod.yaml)

Pod simples para testes:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: "go-server"
  labels:
    app: "go-server"
spec:
  containers:
    - name: "go-server"
      image: "pedrovieira249/golang-kubernetes-teste:latest"
```

**Uso:**
```bash
kubectl apply -f pod.yaml
kubectl get pods
kubectl logs go-server
kubectl delete -f pod.yaml
```

**⚠️ Limitação:** Se o pod morrer, NÃO é recriado automaticamente. Use Deployment em produção.

---

#### [replicaset.yaml](replicaset.yaml)

Garante 5 réplicas do pod:

```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: "go-server"
spec:
  replicas: 5
  selector:
    matchLabels:
      app: "go-server"
  template:
    metadata:
      labels:
        app: "go-server"
    spec:
      containers:
        - name: "go-server"
          image: "pedrovieira249/golang-kubernetes-teste:v2"
```

**Uso:**
```bash
kubectl apply -f replicaset.yaml
kubectl get rs
kubectl get pods  # Verá 5 pods

# Teste: deletar um pod
kubectl delete pod go-server-xxxxx
kubectl get pods  # Nova réplica criada automaticamente

# Escalar
kubectl scale rs go-server --replicas=10
```

**⚠️ Limitação:** Não suporta rolling updates. Use Deployment.

---

#### [deployment.yaml](deployment.yaml) ⭐ PRINCIPAL

Deployment completo com todas as features:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-server
  template:
    spec:
      automountServiceAccountToken: false
      containers:
        - name: go-server
          image: pedrovieira249/golang-kubernetes-teste:v5.6
          
          # PROBES
          startupProbe:
            httpGet:
              path: /healthz
              port: 8080
            periodSeconds: 3
            failureThreshold: 30
          
          readinessProbe:
            httpGet:
              path: /healthz
              port: 8080
            periodSeconds: 2
            failureThreshold: 1
          
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8080
            periodSeconds: 5
            failureThreshold: 3
          
          # VARIÁVEIS DE AMBIENTE
          envFrom:
            - configMapRef:
                name: go-server-config-env
            - secretRef:
                name: goserver-secret
          
          # RECURSOS
          resources:
            requests:
              memory: "32Mi"
              cpu: "10m"
            limits:
              memory: "64Mi"
              cpu: "50m"
          
          # VOLUMES
          volumeMounts:
            - name: config
              mountPath: /app/config
              readOnly: true
            - name: pvc-storage
              mountPath: /app/pvc
      
      volumes:
        - name: config
          configMap:
            name: go-server-configmap
            items:
              - key: valeus
                path: myconfig.txt
        - name: pvc-storage
          persistentVolumeClaim:
            claimName: go-server-pvc
```

**Features implementadas:**
- ✅ 3 réplicas com auto-healing
- ✅ StartupProbe (aguarda 30x3s = 90s para iniciar)
- ✅ ReadinessProbe (controla quando recebe tráfego)
- ✅ LivenessProbe (reinicia se ficar não-saudável)
- ✅ ConfigMap como variáveis de ambiente
- ✅ ConfigMap como arquivo montado
- ✅ Secrets como variáveis de ambiente
- ✅ PersistentVolumeClaim para storage
- ✅ Resource limits (CPU e memória)

**Uso:**
```bash
kubectl apply -f deployment.yaml
kubectl get deployments
kubectl get pods -w

# Ver probes em ação
kubectl describe pod <pod-name>
kubectl logs <pod-name>

# Rolling update
kubectl set image deployment/go-server go-server=pedrovieira249/golang-kubernetes-teste:v6
kubectl rollout status deployment/go-server

# Rollback
kubectl rollout undo deployment/go-server

# Escalar
kubectl scale deployment go-server --replicas=5
```

---

### 4. Services

#### [service.yaml](service.yaml)

Expõe o deployment via LoadBalancer:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: "go-server-service"
spec:
  selector:
    app: "go-server"
  type: LoadBalancer
  ports:
    - protocol: "TCP"
      port: 80
      targetPort: 8080
```

**O que faz:**
- Cria um endpoint estável `go-server-service`
- Distribui tráfego entre todos os pods com label `app: go-server`
- Recebe requisições na porta 80
- Encaminha para porta 8080 dos containers

**Uso:**
```bash
kubectl apply -f service.yaml
kubectl get svc

# No Kind, LoadBalancer não ganha IP externo
# Use port-forward:
kubectl port-forward service/go-server-service 8080:80

# Acessar:
curl http://localhost:8080
```

**Tipos de Service:**
- `ClusterIP`: Apenas interno (padrão)
- `NodePort`: Expõe em porta do node (30000-32767)
- `LoadBalancer`: Pede IP externo (cloud providers)

---

### 5. ConfigMaps

#### [configmap-env.yaml](configmap-env.yaml)

Variáveis de ambiente:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: go-server-config-env
data:
  SAUDACAO: "Olá, Mundo blz!"
  MENSAGEM: "Bem-vindo ao servidor Go!"
```

**Uso no Deployment:**
```yaml
envFrom:
  - configMapRef:
      name: go-server-config-env
```

**Resultado:** Variáveis `SAUDACAO` e `MENSAGEM` disponíveis na aplicação.

---

#### [configmap-config.yaml](configmap-config.yaml)

Arquivo de configuração:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: go-server-configmap
data:
  valeus: "teste, 123141, $H2fer7mreo3$12w1K%#%, ####, @!@!@!, Fim do arquivo"
```

**Uso no Deployment:**
```yaml
volumes:
  - name: config
    configMap:
      name: go-server-configmap
      items:
        - key: valeus
          path: myconfig.txt  # Arquivo criado em /app/config/myconfig.txt

volumeMounts:
  - name: config
    mountPath: /app/config
    readOnly: true
```

**Resultado:** Arquivo `/app/config/myconfig.txt` com o conteúdo do ConfigMap.

**Comandos:**
```bash
kubectl apply -f configmap-env.yaml
kubectl apply -f configmap-config.yaml
kubectl get cm
kubectl describe cm go-server-config-env
kubectl get cm go-server-configmap -o yaml

# Editar ConfigMap (pods precisam ser reiniciados)
kubectl edit cm go-server-config-env
kubectl rollout restart deployment go-server
```

---

### 6. Secrets

#### [secret.yaml](secret.yaml)

Armazena dados sensíveis (codificados em base64):

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: goserver-secret
type: Opaque
data:
  USERNAME: UGVkcm8gQXVnb3N0bwo=  # "Pedro Augusto" em base64
  PASSWORD: MTIzNDU2Cg==          # "123456" em base64
```

**Codificar/Decodificar base64:**
```bash
# Codificar
echo -n "Pedro Augusto" | base64
# UGVkcm8gQXVnb3N0bwo=

# Decodificar
echo "UGVkcm8gQXVnb3N0bwo=" | base64 -d
# Pedro Augusto
```

**Uso no Deployment:**
```yaml
envFrom:
  - secretRef:
      name: goserver-secret
```

**Criar Secret via comando:**
```bash
# Método 1: Literal
kubectl create secret generic goserver-secret \
  --from-literal=USERNAME="Pedro Augusto" \
  --from-literal=PASSWORD="123456"

# Método 2: De arquivo
echo -n "Pedro Augusto" > username.txt
echo -n "123456" > password.txt
kubectl create secret generic goserver-secret \
  --from-file=USERNAME=username.txt \
  --from-file=PASSWORD=password.txt

# Ver secret (valores ocultos)
kubectl get secrets
kubectl describe secret goserver-secret

# Ver valores (CUIDADO!)
kubectl get secret goserver-secret -o yaml
kubectl get secret goserver-secret -o jsonpath='{.data.PASSWORD}' | base64 -d
```

---

### 7. Storage (Volumes)

#### [pv.yaml](pv.yaml)

PersistentVolume - armazenamento disponível no cluster:

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv1
spec:
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce  # Pode ser montado em um único node
```

**Access Modes:**
- `ReadWriteOnce` (RWO): Leitura/escrita em 1 node apenas
- `ReadOnlyMany` (ROX): Somente leitura em múltiplos nodes
- `ReadWriteMany` (RWX): Leitura/escrita em múltiplos nodes
- `ReadWriteOncePod` (RWOP): Leitura/escrita em 1 pod específico

---

#### [pvc.yaml](pvc.yaml)

PersistentVolumeClaim - solicita storage:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: go-server-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
```

**Fluxo:**
1. PVC solicita 5Gi com ReadWriteOnce
2. Kubernetes encontra um PV compatível
3. Faz "bind" entre PVC e PV
4. Deployment monta o PVC no pod

**Comandos:**
```bash
kubectl apply -f pv.yaml
kubectl apply -f pvc.yaml

kubectl get pv
kubectl get pvc

# Ver bind entre PV e PVC
kubectl describe pvc go-server-pvc
```

---

### 8. Auto-Scaling

#### [hpa.yaml](hpa.yaml)

Horizontal Pod Autoscaler - escala baseado em CPU:

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: go-server-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: go-server
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 75
```

**Comportamento:**
- Monitora CPU usage do deployment `go-server`
- Se CPU > 75%: aumenta pods (até max 10)
- Se CPU < 75%: diminui pods (até min 2)
- Aguarda ~5min antes de fazer scale down

**Pré-requisito:** Metrics Server instalado!

```bash
kubectl apply -f hpa.yaml
kubectl get hpa

# Monitorar em tempo real
kubectl get hpa -w

# Ver métricas
kubectl top pods -l app=go-server
```

**Testar com carga:**
```bash
# Terminal 1: Monitorar HPA
kubectl get hpa -w

# Terminal 2: Gerar carga
kubectl run fortio --rm -it --image=fortio/fortio -- \
  load -c 300 -qps 0 -t 1m http://go-server-service/

# Observar pods sendo criados
kubectl get pods -l app=go-server -w
```

---

### 9. StatefulSets (Banco de Dados)

#### [mysql-service-h.yaml](mysql-service-h.yaml)

Serviço Headless para StatefulSet:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql-h
spec:
  clusterIP: None  # Headless = sem IP do service
  selector:
    app: mysql
  ports:
    - port: 3306
```

**Por que Headless?**
- StatefulSets precisam de identidade de rede estável
- Cada pod tem DNS único: `mysql-0.mysql-h`, `mysql-1.mysql-h`, etc
- Permite comunicação direta entre pods

---

#### [statefulset.yaml](statefulset.yaml)

MySQL com storage persistente:

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
spec:
  serviceName: mysql-h  # Referencia o service headless
  replicas: 3
  selector:
    matchLabels:
      app: mysql
  template:
    spec:
      containers:
        - name: mysql-container
          image: mysql:5.7
          env:
            - name: MYSQL_ROOT_PASSWORD
              value: "root"
          resources:
            requests:
              cpu: "500m"
              memory: "512Mi"
          volumeMounts:
            - name: mysql-volume
              mountPath: /var/lib/mysql
  
  volumeClaimTemplates:  # Cria PVC automaticamente para cada pod
    - metadata:
        name: mysql-volume
      spec:
        accessModes:
          - ReadWriteOnce
        resources:
          requests:
            storage: 5Gi
```

**Características do StatefulSet:**
- Pods criados em ordem: mysql-0, mysql-1, mysql-2
- Cada pod tem storage dedicado
- Pods deletados são recriados com mesmo nome
- PVCs persistem mesmo se StatefulSet for deletado

**Uso:**
```bash
# Aplicar service headless primeiro
kubectl apply -f mysql-service-h.yaml

# Depois o StatefulSet
kubectl apply -f statefulset.yaml

# Acompanhar criação sequencial
kubectl get pods -l app=mysql -w

# Ver PVCs criados
kubectl get pvc

# Conectar ao MySQL
kubectl exec -it mysql-0 -- mysql -uroot -proot

# Testar persistência
# No MySQL:
CREATE DATABASE teste;
USE teste;
CREATE TABLE usuarios (id INT, nome VARCHAR(50));
INSERT INTO usuarios VALUES (1, 'Pedro');
EXIT;

# Deletar pod
kubectl delete pod mysql-0

# Pod recriado automaticamente
kubectl get pods -l app=mysql -w

# Conectar novamente e verificar dados
kubectl exec -it mysql-0 -- mysql -uroot -proot -e "SELECT * FROM teste.usuarios;"
```

---

## �💡 Exemplos Práticos

### Exemplo 1: Deploy Completo de uma Aplicação Web

```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: minha-app

---
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: minha-app
data:
  APP_ENV: "production"
  LOG_LEVEL: "info"

---
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secret
  namespace: minha-app
type: Opaque
data:
  db-password: cGFzc3dvcmQxMjM=  # password123 em base64

---
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
  namespace: minha-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web
  template:
    metadata:
      labels:
        app: web
    spec:
      containers:
      - name: web
        image: nginx:latest
        ports:
        - containerPort: 80
        env:
        - name: APP_ENV
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: APP_ENV
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: app-secret
              key: db-password
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 80
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 80
          initialDelaySeconds: 5
          periodSeconds: 5

---
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: web-service
  namespace: minha-app
spec:
  type: LoadBalancer
  selector:
    app: web
  ports:
  - port: 80
    targetPort: 80
```

**Aplicar:**
```bash
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

# Ou tudo de uma vez
kubectl apply -f .
```

---

### Exemplo 2: Aplicação Laravel com MySQL

```yaml
# mysql-deployment.yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pvc
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - name: mysql
        image: mysql:8
        env:
        - name: MYSQL_ROOT_PASSWORD
          value: "rootpassword"
        - name: MYSQL_DATABASE
          value: "laravel"
        - name: MYSQL_USER
          value: "laravel"
        - name: MYSQL_PASSWORD
          value: "laravelpassword"
        ports:
        - containerPort: 3306
        volumeMounts:
        - name: mysql-storage
          mountPath: /var/lib/mysql
      volumes:
      - name: mysql-storage
        persistentVolumeClaim:
          claimName: mysql-pvc

---
apiVersion: v1
kind: Service
metadata:
  name: mysql-service
spec:
  selector:
    app: mysql
  ports:
  - port: 3306
    targetPort: 3306

---
# laravel-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: laravel-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: laravel
  template:
    metadata:
      labels:
        app: laravel
    spec:
      containers:
      - name: laravel
        image: sua-imagem-laravel:latest
        ports:
        - containerPort: 8000
        env:
        - name: DB_HOST
          value: "mysql-service"
        - name: DB_DATABASE
          value: "laravel"
        - name: DB_USERNAME
          value: "laravel"
        - name: DB_PASSWORD
          value: "laravelpassword"
        - name: APP_KEY
          value: "base64:sua-app-key-aqui"

---
apiVersion: v1
kind: Service
metadata:
  name: laravel-service
spec:
  type: NodePort
  selector:
    app: laravel
  ports:
  - port: 80
    targetPort: 8000
    nodePort: 30080
```

**Aplicar:**
```bash
kubectl apply -f mysql-deployment.yaml
kubectl apply -f laravel-deployment.yaml

# Esperar MySQL ficar pronto
kubectl wait --for=condition=ready pod -l app=mysql --timeout=300s

# Verificar
kubectl get pods
kubectl get svc
```

---

### Exemplo 3: Aplicação Go (do seu projeto)

```yaml
# go-app-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-server
  labels:
    app: go-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-server
  template:
    metadata:
      labels:
        app: go-server
    spec:
      containers:
      - name: go-server
        image: pedrovieira249/golang-kubernetes-teste:latest
        ports:
        - containerPort: 8080
        resources:
          requests:
            memory: "32Mi"
            cpu: "100m"
          limits:
            memory: "64Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /
            port: 8080
          initialDelaySeconds: 3
          periodSeconds: 5

---
apiVersion: v1
kind: Service
metadata:
  name: go-server-service
spec:
  type: NodePort
  selector:
    app: go-server
  ports:
  - port: 80
    targetPort: 8080
    nodePort: 30080
```

**Aplicar no Kind:**
```bash
# 1. Carregar imagem no Kind
kind load docker-image pedrovieira249/golang-kubernetes-teste:latest --name kubkindteste

# 2. Aplicar deployment
kubectl apply -f go-app-deployment.yaml

# 3. Verificar
kubectl get deployments
kubectl get pods
kubectl get svc

# 4. Acessar
kubectl port-forward svc/go-server-service 8080:80

# 5. Testar
curl http://localhost:8080
```

---

## 🔄 Fluxo de Trabalho Típico

### Desenvolvimento Local

```bash
# 1. Criar cluster Kind
kind create cluster --name dev

# 2. Escrever código da aplicação
# 3. Criar Dockerfile
# 4. Build da imagem
docker build -t minha-app:v1 .

# 5. Carregar no Kind
kind load docker-image minha-app:v1 --name dev

# 6. Criar manifestos Kubernetes
# 7. Aplicar
kubectl apply -f k8s/

# 8. Verificar
kubectl get all

# 9. Acessar aplicação
kubectl port-forward svc/minha-app 8080:80

# 10. Ver logs
kubectl logs -l app=minha-app -f

# 11. Fazer mudanças
# 12. Rebuild e reload
docker build -t minha-app:v2 .
kind load docker-image minha-app:v2 --name dev
kubectl set image deployment/minha-app app=minha-app:v2

# 13. Verificar rollout
kubectl rollout status deployment/minha-app
```

---

## 🎓 Conceitos Avançados (Resumo)

### Labels e Selectors

**Labels:** Tags chave-valor para organizar recursos

```yaml
metadata:
  labels:
    app: web
    environment: production
    version: v1
```

**Selectors:** Filtram recursos por labels

```bash
kubectl get pods -l app=web
kubectl get pods -l environment=production,version=v1
```

### Health Checks

**Liveness Probe:** Verifica se container está vivo (reinicia se falhar)
**Readiness Probe:** Verifica se container está pronto para receber tráfego

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5
```

### Resource Limits

```yaml
resources:
  requests:  # Mínimo garantido
    memory: "64Mi"
    cpu: "250m"
  limits:    # Máximo permitido
    memory: "128Mi"
    cpu: "500m"
```

### Horizontal Pod Autoscaler (HPA)

Escala automaticamente baseado em métricas:

```bash
kubectl autoscale deployment go-server --cpu-percent=50 --min=2 --max=10
```

**Proxy:** Cria um proxy para acessar a API do Kubernets

```bash
kubectl proxy --port=8080
curl http://localhost:8080/api/v1/namespaces/default/services/go-server-service
```

**Metrics Service:** Cria um mestrics services para a realizar escalonamento automatico do Kubernets
Primeiro baixe esse arquivo: https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
Renomeio para metrics-server.yaml

```bash
kubectl apply -f metrics-server.yaml
kubectl get apiservices
```
Após rodar o segundo comando deve aparecer 
v1beta1.metrics.k8s.io  kube-system/metrics-server True

---

---

## 🔥 Fortio - Testes de Carga

**Fortio** é uma ferramenta de load testing HTTP/gRPC para validar HPA e performance de aplicações Kubernetes.

### Método 1: Usando --rm (Recomendado)

Remove o pod automaticamente após conclusão:

```bash
# Teste básico
kubectl run fortio --rm -i --tty --image=fortio/fortio -- \
  load -c 10 -qps 0 -t 2m http://go-server-service/

# Parâmetros:
# --rm              Remove pod automaticamente ao finalizar
# -i                Modo interativo
# --tty             Aloca terminal
# -c 10             10 conexões simultâneas
# -qps 0            Queries por segundo ilimitadas (máximo)
# -t 2m             Duração: 2 minutos
```

### Método 2: Usando --generator (Versões antigas)

Para kubectl < 1.18 (depreciado):

```bash
kubectl run -it --generator=run-pod/v1 fortio --rm --image=fortio/fortio -- \
  load -c 10 -qps 0 -t 2m http://go-server-service/
```

**⚠️ Nota:** A flag `--generator` foi depreciada no Kubernetes 1.18 e removida na versão 1.25. Use o método sem `--generator` em versões recentes.

### Exemplos de Testes

```bash
# Teste leve (2 conexões, 30 segundos)
kubectl run fortio --rm -i --tty --image=fortio/fortio -- \
  load -c 2 -qps 0 -t 30s http://go-server-service/

# Teste médio (20 conexões, 2 minutos)
kubectl run fortio --rm -i --tty --image=fortio/fortio -- \
  load -c 20 -qps 0 -t 2m http://go-server-service/

# Teste pesado (50 conexões, 5 minutos)
kubectl run fortio --rm -i --tty --image=fortio/fortio -- \
  load -c 50 -qps 0 -t 5m http://go-server-service/

# Teste com QPS limitado (100 queries/segundo)
kubectl run fortio --rm -i --tty --image=fortio/fortio -- \
  load -c 10 -qps 100 -t 2m http://go-server-service/
```

### Monitoramento Durante Teste

Execute em terminais separados:

```bash
# Terminal 1: Executar Fortio
kubectl run fortio --rm -i --tty --image=fortio/fortio -- \
  load -c 20 -qps 0 -t 3m http://go-server-service/

# Terminal 2: Monitorar HPA
kubectl get hpa --watch

# Terminal 3: Monitorar Pods
kubectl get pods -l app=go-server --watch

# ou com output contínuo
kubectl get pods -l app=go-server -w
```

### Interpretando Resultados

Exemplo de saída do Fortio:

```
Fortio 1.60.3 running at 0 queries per second...
Starting at max qps with 10 thread(s)
Code 200 : 245867 (100.0 %)
Response time histogram:
Avg: 0.487 ms, Min: 0.156 ms, Max: 45.2 ms
# target 50% 0.425
# target 75% 0.565
# target 90% 0.712
# target 99% 1.234
# target 99.9% 3.456
All done 245867 calls (plus 10 warmup) 0.487 ms avg, 2049.5 qps
```

**Análise dos Resultados:**
- **QPS:** 2049.5 queries/segundo alcançadas
- **Latência média:** 0.487 ms
- **Taxa de sucesso:** 100% (Code 200)
- **p50:** 50% das requisições em < 0.425 ms
- **p99:** 99% das requisições em < 1.234 ms

### Testando Escalonamento HPA

```bash
# 1. Verificar estado inicial
kubectl get hpa
kubectl get pods -l app=go-server

# 2. Iniciar teste de carga
kubectl run fortio --rm -i --tty --image=fortio/fortio -- \
  load -c 300 -qps 0 -t 1m http://go-server-service/

# 3. Em outro terminal, observar escalonamento
watch -n 1 'kubectl get hpa && echo "---" && kubectl get pods -l app=go-server'

# 4. Após teste, aguardar scale down (5-10 minutos)
```

### Múltiplos Pods de Carga

Para gerar carga mais intensa:

```bash
# Iniciar 3 pods de carga simultaneamente
for i in {1..3}; do
  kubectl run fortio-$i --image=fortio/fortio -- \
    load -c 50 -qps 0 -t 5m http://go-server-service/ &
done

# Limpar depois
kubectl delete pod -l run=fortio
```

### Dicas Importantes

1. **Service Name:** Use o nome do service (`go-server-service`) ao invés de IPs
2. **Namespace:** Se o service está em outro namespace: `http://go-server-service.namespace.svc.cluster.local/`
3. **Duração:** Comece com testes curtos (30s-1m) e aumente gradualmente
4. **Conexões:** Aumente `-c` gradualmente para evitar sobrecarga
5. **Limpeza:** Com `--rm`, o pod é removido automaticamente

---

## 📊 Resumo Visual

### Hierarquia de Objetos

```
Cluster
  └─ Namespaces
      └─ Deployments
          └─ ReplicaSets
              └─ Pods
                  └─ Containers
```

### Fluxo de Tráfego

```
Internet
  ↓
Ingress (nginx/traefik)
  ↓
Service (Load Balancer)
  ↓
Pods (distribuição automática)
  ↓
Containers
```

### Ciclo de Vida de um Deployment

```
1. kubectl apply -f deployment.yaml
   ↓
2. API Server recebe requisição
   ↓
3. Controller Manager cria ReplicaSet
   ↓
4. ReplicaSet cria Pods
   ↓
5. Scheduler atribui Pods aos Nodes
   ↓
6. Kubelet nos Nodes cria containers
   ↓
7. Container Runtime (Docker) roda containers
```

---

## 🎯 Checklist de Aprendizado

### Nível Iniciante
- [ ] Entender o que é Kubernetes e por que usar
- [ ] Conhecer arquitetura básica (Control Plane + Nodes)
- [ ] Criar e gerenciar Pods
- [ ] Usar kubectl para comandos básicos
- [ ] Entender Namespaces
- [ ] Trabalhar com Deployments
- [ ] Criar Services (ClusterIP, NodePort)

### Nível Intermediário
- [ ] Dominar manifestos YAML
- [ ] Usar ConfigMaps e Secrets
- [ ] Entender ReplicaSets
- [ ] Trabalhar com Volumes e PVCs
- [ ] Implementar Health Checks
- [ ] Fazer Rolling Updates e Rollbacks
- [ ] Usar Labels e Selectors efetivamente
- [ ] Configurar Resource Limits

### Nível Avançado
- [ ] Configurar Ingress
- [ ] Implementar HPA (Horizontal Pod Autoscaler)
- [ ] Usar StatefulSets
- [ ] Configurar RBAC (Role-Based Access Control)
- [ ] Trabalhar com DaemonSets
- [ ] Implementar Network Policies
- [ ] Usar Helm para package management
- [ ] Monitoramento e Logging (Prometheus, Grafana)

---

## 📚 Recursos para Estudo

### Documentação Oficial
- **Kubernetes Docs:** https://kubernetes.io/docs/
- **kubectl Cheat Sheet:** https://kubernetes.io/docs/reference/kubectl/cheatsheet/

### Tutoriais Interativos
- **Katacoda:** https://www.katacoda.com/courses/kubernetes
- **Play with Kubernetes:** https://labs.play-with-k8s.com/

### Cursos Recomendados
- Kubernetes para Iniciantes (gratuito na Udemy)
- Certified Kubernetes Administrator (CKA)
- CKAD (Certified Kubernetes Application Developer)

### Ferramentas Úteis
- **k9s:** Interface TUI para Kubernetes
- **kubectx/kubens:** Trocar contexts e namespaces facilmente
- **Helm:** Package manager para Kubernetes
- **Lens:** IDE para Kubernetes

---

## 🤝 Próximos Passos

Agora que você conhece os conceitos básicos:

1. **Pratique com Kind:** Crie clusters locais e experimente
2. **Faça deploys reais:** Use suas aplicações Laravel/Go
3. **Aprenda YAML:** Domine a sintaxe dos manifestos
4. **Estude patterns:** Sidecar, Ambassador, Adapter
5. **Explore ferramentas:** Helm, Kustomize, ArgoCD
6. **Considere certificação:** CKA ou CKAD

**Lembre-se:** Kubernetes tem uma curva de aprendizado, mas com prática constante tudo faz sentido!

---

## 🔧 Troubleshooting e Boas Práticas

### Problemas Comuns e Soluções

#### 1. Pod fica em `ImagePullBackOff` ou `ErrImagePull`

**Problema:** Kubernetes não consegue baixar a imagem Docker.

**Soluções:**
```bash
# Para Kind: carregar a imagem localmente
kind load docker-image pedrovieira249/golang-kubernetes-teste:v5.6 --name kubkindteste

# Verificar se a imagem existe localmente
docker images | grep golang-kubernetes-teste

# Ver exatamente qual o erro
kubectl describe pod <pod-name>

# Verificar nome da imagem no deployment
kubectl get deployment go-server -o yaml | grep image:
```

**Causa comum:** Esqueceu de fazer `kind load` ou nome da imagem está errado.

---

#### 2. Pod fica em `Pending`

**Problema:** Pod não consegue ser agendado em nenhum node.

**Diagnóstico:**
```bash
# Ver detalhes do pod
kubectl describe pod <pod-name>

# Verificar eventos
kubectl get events --sort-by=.metadata.creationTimestamp

# Ver recursos dos nodes
kubectl top nodes
kubectl describe nodes
```

**Causas comuns:**
- **Recursos insuficientes:** CPU/memória solicitada maior que disponível
- **PVC não está bound:** Volume não foi provisionado
- **Node selector incompatível:** Pod requer node específico que não existe

**Soluções:**
```bash
# Reduzir recursos solicitados no deployment
resources:
  requests:
    cpu: "10m"      # Reduzir de 100m para 10m
    memory: "32Mi"  # Reduzir de 128Mi para 32Mi

# Verificar PVCs
kubectl get pvc
kubectl describe pvc go-server-pvc
```

---

#### 3. Pod fica em `CrashLoopBackOff`

**Problema:** Container inicia e morre repetidamente.

**Diagnóstico:**
```bash
# Ver logs do pod
kubectl logs <pod-name>
kubectl logs <pod-name> --previous  # Logs da execução anterior

# Ver eventos
kubectl describe pod <pod-name>
```

**Causas comuns:**
- Aplicação tem erro e termina imediatamente
- Falta variável de ambiente obrigatória
- ConfigMap ou Secret não existe
- Falha no healthcheck logo após iniciar

**Soluções:**
```bash
# Verificar se ConfigMaps e Secrets existem
kubectl get cm
kubectl get secrets

# Ajustar startupProbe para dar mais tempo
startupProbe:
  periodSeconds: 3
  failureThreshold: 30  # 90 segundos total

# Verificar logs da aplicação
kubectl logs -f <pod-name>
```

---

#### 4. HPA não funciona (fica `<unknown>`)

**Problema:** HPA não consegue obter métricas.

**Diagnóstico:**
```bash
# Verificar se metrics server está rodando
kubectl get deployment metrics-server -n kube-system

# Verificar API de métricas
kubectl get apiservices | grep metrics

# Tentar obter métricas manualmente
kubectl top nodes
kubectl top pods
```

**Soluções:**
```bash
# Reinstalar metrics server com flag para Kind
# Editar metrics-server.yaml e adicionar:
- --kubelet-insecure-tls

# Reaplicar
kubectl delete -f metrics-server.yaml
kubectl apply -f metrics-server.yaml

# Aguardar 1-2 minutos
kubectl wait --for=condition=ready pod -l k8s-app=metrics-server -n kube-system --timeout=300s
```

---

#### 5. Service não roteia tráfego

**Problema:** Service criado mas não consegue acessar pods.

**Diagnóstico:**
```bash
# Verificar se service tem endpoints
kubectl get endpoints go-server-service

# Se vazio, problema com selector
kubectl describe svc go-server-service

# Verificar labels dos pods
kubectl get pods --show-labels

# Verificar se pods estão Ready
kubectl get pods -l app=go-server
```

**Soluções:**
```bash
# Garantir que labels do selector correspondem aos labels dos pods
# No service.yaml:
selector:
  app: "go-server"  # Deve ser igual ao label dos pods

# No deployment.yaml:
template:
  metadata:
    labels:
      app: "go-server"  # Deve ser igual ao selector
```

---

### Boas Práticas

#### 1. Ordem de Aplicação dos Recursos

```bash
# 1. Configurações (não dependem de nada)
kubectl apply -f metrics-server.yaml

# 2. Storage
kubectl apply -f pv.yaml
kubectl apply -f pvc.yaml

# 3. Configurações da aplicação
kubectl apply -f configmap-env.yaml
kubectl apply -f configmap-config.yaml
kubectl apply -f secret.yaml

# 4. Aplicação
kubectl apply -f deployment.yaml

# 5. Services
kubectl apply -f service.yaml
kubectl apply -f mysql-service-h.yaml

# 6. StatefulSets
kubectl apply -f statefulset.yaml

# 7. Auto-scaling (por último)
kubectl apply -f hpa.yaml
```

---

#### 2. Labels e Anotações

**Use labels consistentes:**
```yaml
metadata:
  labels:
    app: go-server
    version: v5.6
    environment: production
```

**Filtrar recursos:**
```bash
kubectl get pods -l app=go-server
kubectl get all -l app=go-server
```

---

#### 3. Resource Limits Adequados

```yaml
resources:
  requests:    # Reserva garantida
    cpu: "100m"
    memory: "128Mi"
  limits:      # Máximo permitido
    cpu: "500m"
    memory: "512Mi"
```

**Diretrizes:**
- `requests` < `limits`
- Monitorar uso real: `kubectl top pods`
- Ajustar baseado em métricas reais

---

#### 4. Versionamento de Imagens

**❌ Evite:**
```yaml
image: pedrovieira249/golang-kubernetes-teste:latest
```

**✅ Prefira:**
```yaml
image: pedrovieira249/golang-kubernetes-teste:v5.6
```

**Por quê?**
- `latest` pode mudar sem você saber
- Dificulta rollback
- Impossível saber qual versão está rodando

---

#### 5. Comandos Úteis para Desenvolvimento

```bash
# Criar recursos rapidamente
kubectl run nginx --image=nginx --port=80

# Dry-run (ver YAML sem criar)
kubectl run nginx --image=nginx --dry-run=client -o yaml

# Aplicar e ver resultado
kubectl apply -f deployment.yaml && kubectl get pods -w

# Editar recurso diretamente
kubectl edit deployment go-server

# Ver diferenças antes de aplicar
kubectl diff -f deployment.yaml

# Forçar recriação de pods
kubectl rollout restart deployment go-server

# Ver histórico de rollouts
kubectl rollout history deployment go-server
```

---

**Criado por:** Pedro Vieira  
**Data:** Dezembro de 2025  
**Versão:** 2.0
**Atualização:** Adicionado seção completa de instalação, configuração e arquivos do projeto