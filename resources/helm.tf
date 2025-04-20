## Core app services

resource "helm_release" "database" {
  name  = "postgres-database"
  chart = "../helm/postgres-database"

  namespace = kubernetes_namespace.namespace["database"].metadata[0].name
  values = [
    file("../helm/postgres-database/values.yaml")
  ]
}

resource "helm_release" "url-app" {
  name  = "url-app"
  chart = "../helm/app"

  namespace = kubernetes_namespace.namespace["url-app"].metadata[0].name
  values = [
    file("../helm/app/values.yaml")
  ]
}

resource "helm_release" "traefik" {
  name       = "traefik"
  repository = "https://helm.traefik.io/traefik"
  chart      = "traefik"
  version    = "35.0.1"
  
  namespace = "default"
  values = [
    file("../helm/traefik-reverse-proxy/values.yaml")
  ]
}

## Observability
resource "helm_release" "prometheus" {
  name       = "prometheus"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-prometheus-stack"

  namespace  = "monitoring"
  values = [file("../helm/prometheus/values.yaml")]
}
