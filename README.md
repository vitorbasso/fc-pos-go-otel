# fc-pos-go-otel
Segundo lab pós go expert (observabilidade e open telemetry)

(Screenshots de exemplo do resultado disponíveis em [exemplo resposta do servidor A](#exemplo-de-resposta-do-servidor-a) e [exemplo trace no zipkin](#exemplo-de-trace-gerado-no-zipkin))

## TL;DR
* Pra rodar localmente <b style='font-size:1.5em'>você vai precisar configurar sua própria chave do weather api</b> ou no app.env ou no docker-compose
* Suba o projeto usando `docker-compose up`
* Realize pelo menos uma requisição para o Serviço A em `POST http://localhost:8080/temperatures` com o cep no corpo da chamada.
* Entre no UI do zipkin disponível em `http://localhost:9411` e clique em `RUN QUERY` para ver os traces gerados.

## Requerimentos
  * golang versão 1.22.3 ou superior
  * ou docker e docker-compose
  * uma api key do serviço [weather api](https://www.weatherapi.com/)
  * Um client http para realizar as chamadas POST

## Subindo o projeto

 * Para subir o projeto é necessário preencher sua [chave da api do weather api](https://www.weatherapi.com/docs/) em um de dois possíveis lugares:
   - no arquivo `app.env` na root do projeto;
   - ou no arquivo `docker-compose.yaml`, também na root do projeto;

    Em qualquer uma das opções que escolher, o campo onde deve-se completar com a chave é o `WEATHER_API_KEY=`

 * Após configurada sua chave, execute o programa com `docker-compose up`
 * Após executar o programa, você poderá acessá-lo por meio de uma ferramenta cliente http como o [postman](https://www.postman.com/) ou o [curl](https://curl.se/), pela url base `http://localhost:8080`
 * Também é fornecido um arquivo `test.http` na raíz do projeto, caso prefire utilizá-lo

## Servidor Zipkin

  * O zipkin fica acessível, após subir o projeto, em `http://localhost:9411`.

## O teste

* Para testar o funcionamento, suba o `docker-compose.yaml` como descrito acima, espere terminar de subir.
* Execute a chamada para o endpoint POST /temperatures com o cep no corpo da chamada para o serviço A pelo menos uma vez.
* Finalmente, entre no UI do zipkin em `http://localhost:9411` e clique no botão `RUN QUERY` para ver os spans.

## Endpoints Serviço A

O serviço A é acessível pelo endpoint `http://localhost:8080`

 * O endpoint que retorna as temperaturas correspondentes a localização do cep informado é o seguinte:
      ```http
          POST /temperatures
          Content-Type: application/json

          {
            "cep": {cep}
          }
      ```
    Devendo-se substituir {cep} pelo cep no qual você deseja obter as temperaturas.

 * Por fim, há um health check caso apenas queira conferir se a aplicação está rodando corretamente, basta realizar uma request para:
      ```http
          GET /health-check
      ```
    Ele deverá devolver http status 200 com apenas um `.` no corpo.

## Endpoints Serviço B

Caso deseje testar o serviço B sem passar pelo serviço A ele é acessível pelo endpoint `http://localhost:8081` com os seguintes endpoints

 * O endpoint que retorna as temperaturas correspondentes a localização do cep informado é o seguinte:
      ```http
          GET /temperatures/{cep}
      ```
    Devendo-se substituir {cep} pelo cep no qual você deseja obter as temperaturas.

 * Por fim, há um health check caso apenas queira conferir se a aplicação está rodando corretamente, basta realizar uma request para:
      ```http
          GET /health-check
      ```
    Ele deverá devolver http status 200 com apenas um `.` no corpo.

## Exemplo de resposta do Servidor A
![image](/docs/images/Screenshot%202024-06-01%20235141.png)

## Exemplo de trace gerado no Zipkin
![image](/docs/images/Screenshot%202024-06-01%20235050.png)
--
![image](/docs/images/Screenshot%202024-06-01%20235537.png)
