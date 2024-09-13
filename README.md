Para iniciar o server: 


docker compose up --build


O servidor será iniciado na porta 8080.

---------------------------------------------------------------------


Post - Criar um produto:

http://localhost:8080/products {"name","description","price","stock"}


Exemplo:


curl -X POST http://localhost:8080/products -H "Content-Type: application/json" -d '{"name": "Example Product","description": "This is an example product.","price": 25.00,"stock":10'}


---------------------------------------------------------------------


Put - Atualizar um produto:


http://localhost:8080/products/{id}

exemplo:

curl -X PUT http://localhost:8080/products/1 -H "Content-Type: application/json" -d '{
  "name": "Updated Product",
  "description": "Updated description.",
  "price": 30.00,
  "stock": 5
}'

------------------------------------------------------
Excluir um produto:

http://localhost:8080/products/{id}


Exemplo:


curl -X DELETE http://localhost:8080/products/1


------------------------------------------------------
Get - Listar todos os produtos:
 

http://localhost:8080/products


Exemplo:


curl http://localhost:8080/products  


Pegar item especifico:  
http://localhost:8080/products/{id}

Exemplo:

Get - curl http://localhost:8080/products/1  

