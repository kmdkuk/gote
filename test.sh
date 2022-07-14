docker-compose build
# docker-compose run --rm gote /gote --debug --target blog.kmdkuk.com --notification slack --mode ping
docker-compose run --rm gote /gote --debug --target https://hogehoge.hoge --notification slack --mode http
