package main

import (
	"awesome/internal/app/domain_role/role"
	"awesome/internal/app/domain_role/role/aggregate"
	"ddd/domain"
	"ddd/infrastructure"
	"ddd/infrastructure/persistence"
	"log"
)

func main() {
	repo, err := persistence.NewRedisRepo("127.0.0.1:6379", "123456", 0)
	if err != nil {
		log.Panic(err)
		return
	}
	err = domain.RunService(domain.Options{
		Repo:           repo,
		Namespace:      "role",
		SessionFactory: aggregate.NewRole,
		EventFactory:   role.Events(),
		CommandFactory: role.Commands(),
		Input:          &infrastructure.DaprProxy{},
	})
	if err != nil {
		log.Panic(err)
	}
}
