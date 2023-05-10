# Init
    -  Executar o docker-compose para subir o mysql e o rabbitmq

# Migrate:

### Requisitos
    - instalar o migrate [https://github.com/golang-migrate/migrate] 
    - instalar o make

### Criar tabelas
    - make migrateup

# Servidor

### Executar o server
    - acessar a pasta cmd/ordersystem
    - rodar o comando: go run main.go wire_gen.go 