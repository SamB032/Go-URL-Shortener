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

  namespace = "kube-system"
  values = [
    file("../helm/traefik-reverse-proxy/values.yaml")
  ]
}

## Observability
resource "helm_release" "prometheus" {
  name       = "prometheus"
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "kube-prometheus-stack"

  namespace = kubernetes_namespace.namespace["monitoring"].metadata[0].name
  values    = [file("../helm/prometheus/values.yaml")]
}

resource "helm_release" "loki" {
  name       = "loki"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "loki-stack"
  version    = "2.9.10"

  namespace = kubernetes_namespace.namespace["monitoring"].metadata[0].name
  values = [
    file("../helm/loki/values.yaml")
  ]
}

resource "helm_release" "tempo" {
  name       = "tempo"
  repository = "https://grafana.github.io/helm-charts"
  chart      = "tempo"
  version    = "1.21.1"

  namespace = kubernetes_namespace.namespace["monitoring"].metadata[0].name

  values = [
    file("../helm/tempo/values.yaml")
  ]
}

resource "kubernetes_config_map" "grafana_dashboards" {
  for_each = { for file in fileset("${path.module}/dashboards", "*.json") : file => file }

  metadata {
    name      = replace(each.key, ".json", "")
    namespace = kubernetes_namespace.namespace["monitoring"].metadata[0].name
    labels = {
      grafana_dashboard = "true"
    }
  }

  data = {
    "${each.key}" = file("${path.module}/dashboards/${each.key}")
  }
}
