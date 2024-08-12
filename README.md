## GO CLI to stress test http endpoints

To run this application you can call 

- docker build -t loadtester .

And then just run it passing a url, number of requests and concurrency

- docker run loadtester --url=http://terra.com.br --requests=100 --concurrency=10