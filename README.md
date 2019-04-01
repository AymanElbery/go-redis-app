<img src="https://www.restapiexample.com/wp-content/uploads/2018/06/golang-redis-databse-example.png">






# go-redis-app
Make web app based on Redis in Go 

use the Redis CLI to add a few additional albums, along with a new likes sorted set by :

`HMSET album:1 title "Book 1" artist "Ayman Elbery" price 4.95 likes 8`
`HMSET album:2 title "Book 2" artist "Ahmed" price 5.95 likes 3`
`HMSET album:3 title "Book 3" artist "Ismail" price 11.95 likes 11`
`HMSET album:4 title "Book 4" artist "Ali" price 2.95 likes 5`
`ZADD likes 8 1 3 2 11 3 5 4`


then run `go run main.go` 

then test on these routes :

GET  Method >> `/album?id=1`
POST Method >> `/like`
GET  Method >> `/popular`
