Para iniciar o server: 


docker compose up --build


O servidor será iniciado na porta 8080.


Criar um produto:


curl -X POST http://localhost:8080/products -H "Content-Type: application/json" -d '{"name": "Example Product","description": "This is an example product.","price": 25.00,"stock":10'}


Atualizar um produto:


curl -X PUT http://localhost:8080/products/1 -H "Content-Type: application/json" -d '{
  "name": "Updated Product",
  "description": "Updated description.",
  "price": 30.00,
  "stock": 5
}'


Excluir um produto:


curl -X DELETE http://localhost:8080/products/1


Listar todos os produtos:
 

curl http://localhost:8080/products  


Pegar item especifico:  


curl http://localhost:8080/products/1  

