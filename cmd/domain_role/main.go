package main

import (
	"ddd/domain"
	"ddd/infrastructure"
	"ddd/infrastructure/persistence"
	"awesome/internal/app/domain_role/aggregate"
	"awesome/internal/app/domain_role"
	"log"
)


func main(){
	repo, err := persistence.NewRedisRepo("127.0.0.1:6379", "123456", 0)
	if err != nil{
		log.Panic(err)
		return
	}
	err = domain.RunService(domain.Options{
		Repo:           repo,
		Namespace:      "role",
		SessionFactory: aggregate.RoleFactory(),
		EventFactory:   domain_role.EventsFactory(),
		CommandFactory: domain_role.CommandFactory(),
		Input: &infrastructure.DaprProxy{},
	})
	if err != nil{
		log.Panic(err)
	}
}
