FROM golang:1.27-alpine AS build

WORKDIR /src

# Copia e baixa as dependências (aproveita o cache do Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copia o código e compila o binário otimizado
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# Imagem de execução Distroless (Ultra segura e leve)
FROM gcr.io/distroless/static-debian12:nonroot

LABEL org.opencontainers.image.title="Go K8s Observability Platform"
LABEL org.opencontainers.image.description="Production-ready Go API on Kubernetes with native Prometheus metrics instrumentation"
LABEL org.opencontainers.image.source="https://github.com/kauatwn/go-k8s-observability-platform"

# Copia o executável do stage de build
COPY --from=build /out/api /api

USER nonroot

EXPOSE 8080

ENTRYPOINT ["/api"]
