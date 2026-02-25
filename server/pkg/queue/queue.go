package queue

import (
	"fmt"
	"github.com/RichardKnop/machinery/v2"
	redisbackend "github.com/RichardKnop/machinery/v2/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v2/brokers/redis"
	"github.com/RichardKnop/machinery/v2/config"
	redislock "github.com/RichardKnop/machinery/v2/locks/redis"
	queuelog "github.com/RichardKnop/machinery/v2/log"
)

func NewQueue(username, password, host, queue string,
	port, db, resultExpireSecond int64,
	tasksMap map[string]interface{},
	log *Logger) (*machinery.Server, error) {
	queuelog.Set(log)
	cnf := &config.Config{
		DefaultQueue:    queue,
		ResultsExpireIn: int(resultExpireSecond),
		Redis: &config.RedisConfig{
			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
			DelayedTasksKey:        queue + "_delay",
		},
	}
	//u := &url.URL{
	//	Scheme: "redis",
	//	User:   url.UserPassword(username, password),
	//	Host:   fmt.Sprintf("%s:%d", host, port),
	//}
	//redisURL := u.String()
	//encodedUsername := url.QueryEscape(username)
	//encodedPassword := url.QueryEscape(password)
	redisURL := fmt.Sprintf("%s:%s@%s:%d", username, password, host, port)
	// Create server instance
	broker := redisbroker.NewGR(cnf, []string{redisURL}, int(db))
	backend := redisbackend.NewGR(cnf, []string{redisURL}, int(db))
	l := redislock.New(cnf, []string{redisURL}, int(db), 10)
	server := machinery.NewServer(cnf, broker, backend, l)
	return server, server.RegisterTasks(tasksMap)
}
