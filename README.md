A simple go app that allows for you to shorten urls.  This works by generating a shortkey and storing this in the database along side the given url. When the visits the url with the shortkey, it will redirect them to the url.

## Helm installation
### Create postgres secrets
For postgres-database helm chart, execute this command in kubernetes to write the following secrets:
```
kubectl create secret generic postgres-secret \
  --from-literal=POSTGRES_PASSWORD=mypassword \
  --from-literal=POSTGRES_USER=myuser \
  --from-literal=POSTGRES_DB=mydb \
  --namespace=database
```
Change this to something more secure.

### Pass postgres secrets to url app
Execute the same command but change the namespace to be url-app
```
kubectl create secret generic postgres-secret \
  --from-literal=POSTGRES_PASSWORD=mypassword \
  --from-literal=POSTGRES_USER=myuser \
  --from-literal=POSTGRES_DB=mydb \
  --namespace=url-app
```

