resource "random_id" "username" {
  byte_length = 8
}

resource "random_password" "password" {
  length  = 16
  special = true
  upper   = true
  lower   = true
}

# Secrets for Database
resource "kubernetes_secret" "database" {
  metadata {
    name      = "postgres-secret"
    namespace = kubernetes_namespace.namespace["database"].metadata[0].name
  }

  data = {
    POSTGRES_USER     = random_id.username.hex
    POSTGRES_PASSWORD = random_password.password.result
    POSTGRES_DB       = "shortkey"
  }
}

# Secrets for App
resource "kubernetes_secret" "app" {
  metadata {
    name      = "postgres-secret"
    namespace = kubernetes_namespace.namespace["url-app"].metadata[0].name
  }

  data = {
    POSTGRES_USER     = random_id.username.hex
    POSTGRES_PASSWORD = random_password.password.result
    POSTGRES_DB       = "shortkey"
  }
}
