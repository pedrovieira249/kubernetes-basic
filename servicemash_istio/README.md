# Service Mesh e Istio — Guia Completo

Este guia explica, de ponta a ponta, o que é um Service Mesh, por que utilizar, e como operar o Istio em Kubernetes. Inclui conceitos, arquitetura, recursos principais, exemplos de YAML, guia de instalação em Windows/WSL, boas práticas e troubleshooting.

## Visão Geral: O que é Service Mesh

- **Definição:** Service Mesh é uma camada de infraestrutura que gerencia comunicação serviço‑a‑serviço (east‑west) de forma transparente, adicionando controle de tráfego, segurança, observabilidade e políticas sem alterar o código das aplicações.
- **Como funciona:** Um proxy sidecar (ex.: Envoy) é injetado ao lado de cada pod. O plano de controle (ex.: `istiod`, no Istio) configura dinamicamente esses proxies.
- **Benefícios:**
	- Tráfego avançado: roteamento, balanceamento, retries, timeouts, circuit breaking.
	- Segurança por padrão: identidade de serviço, mTLS, autorização e autenticação.
	- Observabilidade: métricas ricas, logs, tracing distribuído.
	- Políticas consistentes: rate limits, quotas, acesso baseado em identidade.
- **Quando usar:**
	- Muitas comunicações entre microsserviços que precisam de governança.
	- Requisitos de segurança forte (mTLS), auditoria e compliance.
	- Entregas progressivas (canary, A/B), resiliência e controle fino de tráfego.

## Istio: Conceitos e Arquitetura

- **Plano de Dados (Data Plane):**
	- Proxies Envoy rodando como sidecars em cada pod, interceptando tráfego de entrada e saída.
- **Plano de Controle (Control Plane):**
	- `istiod` distribui configuração, gerencia certificados/identidade (SPIFFE), e aplica políticas.
- **Injeção de Sidecar:**
	- Automática via label no namespace (`istio-injection=enabled`) ou manual por anotação no pod.
- **Identidade e Certificados:**
	- Cada workload recebe um certificado X.509 com identidade SPIFFE; mTLS entre serviços é configurado pelo Istio.

## Recursos Principais do Istio

- **Traffic Management:**
	- `Gateway`: entrada norte‑sul (expor HTTP/TCP do cluster).
	- `VirtualService`: regras de roteamento (paths, headers, pesos, retries, timeouts).
	- `DestinationRule`: políticas por destino (circuit breaking, outlier detection, subset/canary).
	- `ServiceEntry`: registrar destinos externos (egress control) para tráfego fora do cluster.
	- `Sidecar`: escopo de visibilidade e configuração por workload.
- **Segurança:**
	- `PeerAuthentication`: modo mTLS (STRICT, PERMISSIVE, DISABLE) por namespace/serviço.
	- `RequestAuthentication`: validação de JWT e identidade do chamador.
	- `AuthorizationPolicy`: autorização baseada em atributos (origem, paths, métodos, principals).
- **Observabilidade:**
	- Telemetry v2 gera métricas (Prometheus), logs e traços (Jaeger/Zipkin).
	- Integrações comuns: Grafana (dashboards), Kiali (visão de mesh), Jaeger (tracing).
- **Resiliência:**
	- Timeouts, retries, circuit breaking, outlier detection e connection pool tuning.

## Instalação (Windows/WSL + Kubernetes)

Pré‑requisitos:

- Windows 10/11 com WSL 2 habilitado e uma distro Linux (Ubuntu).
- Kubernetes disponível (ex.: Kind, Minikube ou um cluster remoto). Este repositório possui um arquivo `kind.yaml` para provisionar um cluster local.

Passos sugeridos (no terminal da sua distro WSL):

```bash
# 1) Preparar um cluster Kind usando o arquivo existente
kind create cluster --config kind.yaml

# 2) Instalar o cliente do Istio (istioctl)
# Baixe a versão estável do Istio no site oficial (istio.io) e adicione o binário istioctl ao PATH
# Exemplo (ajuste a versão conforme necessário):
curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.20.0 sh -
export PATH="$PATH:$HOME/istio-1.20.0/bin"

# 3) Instalar o Istio no cluster
istioctl install --set profile=default -y

# 4) Habilitar injeção automática de sidecar no namespace desejado
kubectl label namespace default istio-injection=enabled --overwrite

# 5) Validar\
istioctl version
istioctl proxy-status
```

Observação: perfis como `default`, `minimal` e `demo` variam em componentes instalados. O `demo` inclui mais recursos e é útil para testes/POCs (consome mais recursos).

Para addons (Kiali, Grafana, Jaeger), consulte os manifests oficiais do Istio (amplamente disponibilizados no diretório `samples/addons` do repositório do Istio). Aplique conforme a versão escolhida.

## Habilitando o Sidecar e Integrando com seus Manifests

- Com a label de namespace aplicada (`istio-injection=enabled`), novos pods receberão o sidecar automaticamente.
- Para pods existentes, faça um rollout/redeploy para reinjeção.
- Seus manifests de `Deployment`/`Service` existentes geralmente funcionam sem mudanças. Caso necessário, ajuste:
	- **Probes:** garanta que liveness/readiness não sejam afetadas por redirecionamentos do sidecar.
	- **Portas e Protocolos:** utilize portas/names adequados e defina protocolos (HTTP, HTTP2, gRPC) quando relevante.

## Exemplos de YAML (Traffic, Segurança)

### VirtualService + DestinationRule (canary e resiliência)

```yaml
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
	name: minha-api
spec:
	hosts:
		- minha-api
	http:
		- route:
				- destination:
						host: minha-api
						subset: v1
					weight: 80
				- destination:
						host: minha-api
						subset: v2
					weight: 20
			retries:
				attempts: 3
				perTryTimeout: 2s
			timeout: 5s
---
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
	name: minha-api
spec:
	host: minha-api
	subsets:
		- name: v1
			labels:
				version: v1
		- name: v2
			labels:
				version: v2
	trafficPolicy:
		outlierDetection:
			consecutive5xxErrors: 5
			interval: 5s
			baseEjectionTime: 30s
		connectionPool:
			http:
				maxRequestsPerConnection: 100
```

### Gateway (expor tráfego externo)

```yaml
apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:
	name: minha-gateway
spec:
	selector:
		istio: ingressgateway
	servers:
		- port:
				number: 80
				name: http
				protocol: HTTP
			hosts:
				- "minha-api.example.com"
```

### PeerAuthentication (mTLS estrito)

```yaml
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
	name: default
	namespace: default
spec:
	mtls:
		mode: STRICT
```

### AuthorizationPolicy (permitir apenas chamadas internas)

```yaml
apiVersion: security.istio.io/v1beta1
kind: AuthorizationPolicy
metadata:
	name: allow-same-namespace
	namespace: default
spec:
	action: ALLOW
	rules:
		- from:
				- source:
						namespaces: ["default"]
```

### RequestAuthentication (JWT)

```yaml
apiVersion: security.istio.io/v1beta1
kind: RequestAuthentication
metadata:
	name: jwt-auth
	namespace: default
spec:
	selector:
		matchLabels:
			app: minha-api
	jwtRules:
		- issuer: "https://issuer.example.com"
			jwksUri: "https://issuer.example.com/.well-known/jwks.json"
```

## Observabilidade

- **Métricas:** Exportadas via Prometheus; dashboards no Grafana.
- **Tracing:** Jaeger/Zipkin para spans distribuídos.
- **Kiali:** Visualização de serviços, topologia, status de mTLS, rotas e políticas.
- **Dicas:**
	- Garanta headers de tracing (ex.: `x-request-id`, `x-b3-*`, `traceparent`) propagados pela aplicação.
	- Habilite logs de acesso no Envoy quando necessário para investigar rotas.

## Boas Práticas

- **Perfis e Recursos:** Evite perfis muito pesados em produção. Configure requests/limits para sidecars.
- **mTLS:** Use `STRICT` entre serviços internos; trate integrações externas com `ServiceEntry` e TLS apropriado.
- **Protocolos:** Declare protocolos corretos nas portas para funcionalidades avançadas (ex.: gRPC).
- **Canary/Release:** Utilize subsets por label (`version`) e roteamento ponderado.
- **Sidecar Scope:** Restrinja a visibilidade com `Sidecar` para reduzir domínio de falhas e consumo.
- **Upgrades:** Prefira upgrades por **revisões** do Istio (multi‑revision) e migrações graduais.

## Troubleshooting

- **Sidecar não injeta:**
	- Verifique a label do namespace: `istio-injection=enabled`.
	- Cheque admission webhook e logs do `istiod`.
- **Falhas de tráfego (503/404):**
	- Use `istioctl analyze` para validar configuração.
	- Confirme hosts em `VirtualService` e subsets em `DestinationRule`.
- **mTLS quebrando chamadas:**
	- Ajuste `PeerAuthentication` (STRICT vs PERMISSIVE) e políticas de destino.
- **Probes e saúde:**
	- Certifique que as probes não passam pelo proxy ou estão corretamente configuradas.
- **Egress externo bloqueado:**
	- Adicione `ServiceEntry` e, se necessário, `EgressGateway`.

Comandos úteis:

```bash
istioctl analyze
istioctl proxy-status
istioctl proxy-config routes <pod>.<ns>
istioctl proxy-config clusters <pod>.<ns>
istioctl proxy-config listeners <pod>.<ns>
```

## Como aplicar neste repositório

- Crie/valide o cluster com o `kind.yaml`.
- Instale o Istio e habilite a injeção no namespace onde seus manifests (`deployment.yaml`, `service.yaml`, etc.) estão sendo aplicados.
- Para expor serviços, adicione um `Gateway` e `VirtualService` apontando para o seu `Service` Kubernetes existente.
- Ajuste segurança com `PeerAuthentication` e `AuthorizationPolicy` conforme necessidade.

## Referências

- Documentação oficial do Istio: conceitos, instalação e exemplos.
- Envoy Proxy: fundamentos do plano de dados.
- Addons (Kiali, Grafana, Jaeger): disponíveis no repositório oficial do Istio por versão.

---

Se quiser, posso criar arquivos de exemplo (`Gateway`, `VirtualService`, `DestinationRule`) no diretório do projeto para um serviço específico seu e integrar com os manifests já presentes.
