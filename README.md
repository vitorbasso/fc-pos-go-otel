# fc-pos-go-otel
Segundo lab pós go expert (observabilidade e open telemetry)

## TL;DR
* Pra rodar localmente <b style='font-size:1.5em'>você vai precisar configurar sua própria chave do weather api</b> ou no app.env ou no docker-compose

## Requerimentos
  * golang versão 1.22.3 ou superior
  * ou docker e docker-compose
  * uma api key do serviço [weather api](https://www.weatherapi.com/)

## Como executar

### Localmente

 * Para executar o programa localmente é necessário preencher sua [chave da api do weather api](https://www.weatherapi.com/docs/) em um de dois possíveis lugares:
   - no arquivo `app.env` na root do projeto;
   - ou no arquivo `docker-compose.yml`, também na root do projeto;

    Em qualquer uma das opções que escolher, o campo onde deve-se completar com a chave é o `WEATHER_API_KEY=`

 * Após configurada sua chave, execute o programa com `docker-compose up`
 * De forma alternativa, você também pode subir o servidor local com `go run cmd/server/main.go` caso tiver configurado a chave do weather api. Se não tiver configurado, adicione ela no inicio como variável de ambiente da seguinte forma: `WEATHER_API_KEY={chave} go run cmd/server/main.go`.
 * Após executar o programa de uma das formas descritas nessa seção, você poderá acessá-lo pela url base `http://localhost:8080`


## Endpoints

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
