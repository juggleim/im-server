package actorsystem

import "google.golang.org/protobuf/proto"

type Processor func(ActorRef, proto.Message)
type NewInput func() proto.Message
