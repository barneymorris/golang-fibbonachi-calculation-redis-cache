package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func main() {
    client := redis.NewClient(&redis.Options{
        Addr:	  "localhost:6379",
        Password: "", // no password set
        DB:		  0,  // use default DB
    })

	fiberApp := fiber.New()
	
	fiberApp.Get("/api/cache/fibbnonachi/:number", func (c *fiber.Ctx) error {
		n := c.Params("number")

		if n == "" {
			return c.Status(400).JSON("please provided a fibbonachi number")
		}

		number, err := strconv.Atoi(n);
		if err != nil {
			return c.Status(400).JSON("please provided a VALID fibbonachi number")
		}

		s := fmt.Sprintf("number:%d", number)

		var fibb int

		cache, err := client.Get(context.Background(), s).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				f := FibonacciRecursion(number)
				client.Set(context.Background(), s, f, time.Second * 10)
				fibb = f
			} else {
				return c.Status(500).JSON(err.Error())
			}
		} else {
			fibb, _ = strconv.Atoi(cache)
		}

		return c.Status(200).JSON(fibb)
	})

	log.Fatal(fiberApp.Listen(":3000"))
}

func FibonacciRecursion(n int) int {
    if n <= 1 {
        return n
    }
    return FibonacciRecursion(n-1) + FibonacciRecursion(n-2)
}