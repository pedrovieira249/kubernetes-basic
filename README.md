# Guia de Instalação do Kind no Windows

## 📚 O que é cada coisa?

### Kubernetes (K8s)
O **Kubernetes** é uma plataforma open-source para automação de deploy, escalonamento e gerenciamento de aplicações em contêineres. Ele foi originalmente desenvolvido pelo Google e agora é mantido pela Cloud Native Computing Foundation (CNCF).

**Para que serve:**
- Orquestrar contêineres Docker em múltiplos servidores
- Escalar aplicações automaticamente conforme a demanda
- Garantir alta disponibilidade e recuperação automática de falhas
- Gerenciar configurações, secrets e volumes de forma centralizada
- Realizar deploys e rollbacks de forma controlada

**Exemplo prático:** Imagine que você tem uma API Laravel rodando em Docker. Com Kubernetes, você pode:
- Rodar várias cópias da sua API simultaneamente
- Distribuir o tráfego entre elas automaticamente
- Se uma cair, outra sobe automaticamente no lugar
- Escalar de 2 para 10 instâncias em segundos quando o tráfego aumentar

### Kind (Kubernetes in Docker)
**Kind** é uma ferramenta para rodar clusters Kubernetes locais usando contêineres Docker como "nós" (nodes). Foi criado principalmente para testar o próprio Kubernetes, mas é perfeito para desenvolvimento local.

**Para que serve:**
- Criar clusters Kubernetes completos em sua máquina local
- Testar aplicações em um ambiente Kubernetes real sem custos de cloud
- Aprender Kubernetes sem precisar de múltiplas máquinas
- Testar configurações e deployments antes de subir para produção

**Vantagem:** É mais leve que outras soluções (como Minikube com VMs) porque usa Docker nativamente.

### Cluster
Um **Cluster** é um conjunto de máquinas (físicas ou virtuais) trabalhando juntas como um sistema único. No Kubernetes, um cluster é composto por:

**Control Plane (Plano de Controle):**
- Gerencia o cluster inteiro
- Toma decisões sobre onde rodar os contêineres
- Monitora o estado do cluster

**Worker Nodes (Nós de Trabalho):**
- Máquinas que executam suas aplicações
- Cada node pode rodar múltiplos contêineres

**No Kind:** Todo o cluster roda dentro de contêineres Docker na sua máquina, simulando um ambiente real.

---

## 🚀 Instalação do Kind no Windows

### Pré-requisitos
- Docker Desktop instalado e rodando
- PowerShell ou Command Prompt
- Permissões de administrador

### Passo 1: Download do Kind

Abra o PowerShell e execute:

```powershell
# Baixe o executável do Kind
curl.exe -Lo kind-windows-amd64.exe https://kind.sigs.k8s.io/dl/v0.30.0/kind-windows-amd64.exe
```

### Passo 2: Mova para um Diretório Permanente

Crie um diretório para o Kind e mova o executável:

```powershell
# Crie o diretório (se não existir)
New-Item -ItemType Directory -Path "C:\Kind" -Force

# Mova o arquivo e renomeie
Move-Item .\kind-windows-amd64.exe C:\Kind\kind.exe
```

### Passo 3: Configurar Variáveis de Ambiente

#### Opção A: Via Interface Gráfica

1. Pressione `Win + Pause/Break` ou clique com botão direito em **"Este Computador"** → **Propriedades**
2. Clique em **"Configurações avançadas do sistema"**
3. Clique em **"Variáveis de Ambiente"**
4. Em **"Variáveis do sistema"**, encontre a variável **Path** e clique em **"Editar"**
5. Clique em **"Novo"** e adicione: `C:\Kind`
6. Clique em **"OK"** em todas as janelas

#### Opção B: Via PowerShell (como Administrador)

```powershell
# Execute o PowerShell como Administrador
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "Machine") + ";C:\Kind",
    "Machine"
)
```

### Passo 4: Verificar a Instalação

**IMPORTANTE:** Após configurar o PATH, feche e abra uma nova janela do PowerShell.

```powershell
# Verifique se o Kind está acessível
kind version
```

**Saída esperada:**
```
kind v0.30.0 go1.21.0 windows/amd64
```

### Passo 5: Instalar kubectl (CLI do Kubernetes)

O `kubectl` é a ferramenta de linha de comando para interagir com clusters Kubernetes.

```powershell
# Opção 1: Via Chocolatey (se você tiver instalado)
choco install kubernetes-cli

# Opção 2: Download manual
curl.exe -LO "https://dl.k8s.io/release/v1.28.0/bin/windows/amd64/kubectl.exe"

# Mova para o mesmo diretório do Kind
Move-Item .\kubectl.exe C:\Kind\kubectl.exe
```

Verifique a instalação:

```powershell
kubectl version --client
```

---

## 🎯 Criando seu Primeiro Cluster

### Passo 1: Verificar o Docker

Certifique-se de que o Docker Desktop está rodando:

```powershell
docker ps
```

### Passo 2: Criar o Cluster

```powershell
# Crie um cluster com o nome padrão "kind"
kind create cluster

# Ou crie com um nome personalizado
kind create cluster --name meu-cluster
```

**O que está acontecendo:**
1. Kind baixa a imagem Docker do Kubernetes
2. Cria um contêiner que funciona como control plane
3. Configura o kubectl para se conectar ao cluster
4. Inicializa todos os componentes do Kubernetes

**Saída esperada:**
```
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.27.3) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind
```

### Passo 3: Verificar o Cluster

```powershell
# Veja informações do cluster
kubectl cluster-info

# Liste os nós do cluster
kubectl get nodes

# Veja todos os pods do sistema
kubectl get pods -A
```

### Passo 4: Teste com uma Aplicação Simples

Vamos fazer o deploy de um nginx para testar:

```powershell
# Crie um deployment do nginx
kubectl create deployment nginx --image=nginx

# Exponha o deployment como um serviço
kubectl expose deployment nginx --port=80 --type=NodePort

# Veja o status
kubectl get deployments
kubectl get pods
kubectl get services
```

### Passo 5: Acessar a Aplicação

```powershell
# Faça port-forward para acessar localmente
kubectl port-forward service/nginx 8080:80
kubectl port-forward pod/go-server 8080:8080
```

Abra o navegador em `http://localhost:8080` - você verá a página do nginx!

---

## 📋 Comandos Úteis do Kind

### Gerenciamento de Clusters

```powershell
# Listar todos os clusters
kind get clusters

# Ver informações de um cluster específico
kind get nodes --name meu-cluster

# Deletar um cluster
kind delete cluster --name meu-cluster

# Deletar todos os clusters
kind delete clusters --all
```

### Criar Cluster com Configuração Personalizada

Crie um arquivo `kind-config.yaml`:

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
  - role: worker
  - role: worker
```

```powershell
# Crie o cluster usando a configuração
kind create cluster --config kind-config.yaml --name cluster-multi-node
```

### Carregar Imagens Docker Locais

```powershell
# Útil para testar suas próprias imagens Docker
docker build -t minha-app:latest .
kind load docker-image minha-app:latest --name meu-cluster
```

---

## 🔧 Comandos Úteis do kubectl

### Visualização

```powershell
# Ver todos os recursos
kubectl get all

# Ver pods em tempo real
kubectl get pods --watch

# Detalhes de um pod específico
kubectl describe pod <nome-do-pod>

# Ver logs de um pod
kubectl logs <nome-do-pod>

# Ver logs em tempo real
kubectl logs -f <nome-do-pod>
```

### Gerenciamento

```powershell
# Aplicar configuração de um arquivo YAML
kubectl apply -f deployment.yaml

# Deletar recursos
kubectl delete deployment nginx
kubectl delete service nginx

# Escalar um deployment
kubectl scale deployment nginx --replicas=3

# Ver recursos em um namespace específico
kubectl get pods -n kube-system
```

### Troubleshooting

```powershell
# Entrar em um pod (similar ao docker exec)
kubectl exec -it <nome-do-pod> -- /bin/bash

# Ver eventos do cluster
kubectl get events

# Ver uso de recursos
kubectl top nodes
kubectl top pods
```

---

## 🎓 Próximos Passos

Agora que você tem um cluster Kubernetes funcionando, pode:

1. **Fazer deploy de uma aplicação Laravel:**
   - Criar Dockerfile para sua aplicação
   - Criar manifestos Kubernetes (Deployment, Service, ConfigMap)
   - Fazer deploy no cluster Kind

2. **Aprender sobre recursos do Kubernetes:**
   - Deployments (gerenciar réplicas da aplicação)
   - Services (expor aplicações na rede)
   - ConfigMaps e Secrets (gerenciar configurações)
   - Persistent Volumes (armazenamento persistente)
   - Ingress (roteamento HTTP)

3. **Praticar conceitos:**
   - Rolling updates e rollbacks
   - Auto-scaling horizontal
   - Health checks (liveness e readiness probes)
   - Namespaces para organização

---

## 🐛 Problemas Comuns

### "kind: command not found"
- Verifique se adicionou `C:\Kind` no PATH
- Feche e abra uma nova janela do PowerShell
- Execute: `$env:PATH -split ';'` para verificar o PATH

### "Cannot connect to the Docker daemon"
- Verifique se o Docker Desktop está rodando
- Execute: `docker ps` para testar

### "Failed to create cluster"
- Verifique se não há clusters com o mesmo nome: `kind get clusters`
- Delete o cluster antigo se necessário: `kind delete cluster --name <nome>`
- Verifique os logs: `kind create cluster --verbosity=3`

### Cluster muito lento
- Kind usa recursos da sua máquina local
- Feche outros aplicativos pesados
- Configure limites de recursos no Docker Desktop

---

## 📚 Recursos Adicionais

- **Documentação Oficial do Kind:** https://kind.sigs.k8s.io/
- **Documentação do Kubernetes:** https://kubernetes.io/docs/
- **kubectl Cheat Sheet:** https://kubernetes.io/docs/reference/kubectl/cheatsheet/
- **Tutoriais Interativos:** https://www.katacoda.com/courses/kubernetes

---

## 💡 Dicas Finais

1. **Use aliases para comandos comuns:**
   ```powershell
   # Adicione ao seu perfil do PowerShell
   Set-Alias -Name k -Value kubectl
   ```

2. **Crie um cluster dedicado para cada projeto:**
   ```powershell
   kind create cluster --name projeto-laravel
   kind create cluster --name projeto-nodejs
   ```

3. **Sempre delete clusters não utilizados para liberar recursos:**
   ```powershell
   kind delete cluster --name <nome>
   ```

4. **Use context do kubectl para alternar entre clusters:**
   ```powershell
   kubectl config get-contexts
   kubectl config use-context kind-meu-cluster
   ```

---

**Criado por:** Pedro Vieira  
**Data:** Novembro de 2025  
**Versão:** 1.0