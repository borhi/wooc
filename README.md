Please run 

`cp config.template.toml config.toml`

Please edit config.toml and run

`docker build -t wooc . && docker run --publish 80:8080 --name wooc --rm wooc`

Swagger url:

`http://localhost/swaggerui/`