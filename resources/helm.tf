resource "helm_release" "database" {
  name       = "postgres-database"
  chart      = "../helm/postgres-database"

  namespace  = kubernetes_namespace.namespace["database"].metadata[0].name
}

resource "helm_release" "url-app" {
  name       = "url-app"
  chart      = "../helm/app"
  
  namespace  = kubernetes_namespace.namespace["url-app"].metadata[0].name
}
