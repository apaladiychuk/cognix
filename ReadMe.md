### Configure 
Copy file config/.env.example to config/.env

Declare env variables in config/.env
  

### Run docker container 

```azure

  docker compose up -d cockroach
  docker compose up migration     
  docker compose up -d api 
```

