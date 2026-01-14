package stat

import (
	"go/adv-demo/pkg/event"
	"log"
)

type StatServiceDeps struct {
	EventBus       *event.EventBus
	StatRepository *StatRepository
}

type StatService struct {
	EventBus       *event.EventBus
	StatRepository *StatRepository
}

func NewStatService(deps *StatServiceDeps) *StatService {
	return &StatService{
		EventBus:       deps.EventBus,
		StatRepository: deps.StatRepository,
	}
}

func (s *StatService) AddClick() {
	for msg := range s.EventBus.Subsribe() {
		if msg.Type == event.LinkVisitedEvent {
			id, ok := msg.Data.(uint)

			if !ok {
				log.Fatalln("Bad LinkVisitedEvent data ", msg.Data)
				continue
			}

			s.StatRepository.AddClick(id)
		}
	}
}
