package gmicro

import (
	"time"

	"im-server/commons/gmicro/actorsystem"

	"google.golang.org/protobuf/proto"
)

type Cluster struct {
	currentNode *Node
	actorSystem *actorsystem.ActorSystem
}

type IActorRegister interface {
	RegisterActor(method string, actorCreateFun func() actorsystem.IUntypedActor, concurrentCount int)
	RegisterMultiMethodActor(methods []string, actorCreateFun func() actorsystem.IUntypedActor, concurrentCount int)
}

func NewCluster(nodename string, exts map[string]string) *Cluster {
	actorSystem := actorsystem.NewActorSystemNoRpc(nodename)
	//current Node
	curNode := NewNode(nodename, actorSystem.Host, actorSystem.Prot, exts)
	cluster := &Cluster{
		currentNode: curNode,
		actorSystem: actorSystem,
	}
	return cluster
}

func (cluster *Cluster) GetCurrentNode() *Node {
	return cluster.currentNode
}

func (cluster *Cluster) GetAllNodes() []*Node {
	ret := []*Node{cluster.currentNode}
	return ret
}

func (cluster *Cluster) RegisterActor(method string, actorCreateFun func() actorsystem.IUntypedActor, concurrentCount int) {
	cluster.actorSystem.RegisterActor(method, actorCreateFun, concurrentCount)
	cluster.currentNode.AddMethod(method)
}
func (cluster *Cluster) RegisterMultiMethodActor(methods []string, actorCreateFun func() actorsystem.IUntypedActor, concurrentCount int) {
	cluster.actorSystem.RegisterMultiMethodActor(methods, actorCreateFun, concurrentCount)
	for _, method := range methods {
		cluster.currentNode.AddMethod(method)
	}
}

func (cluster *Cluster) Startup() {

}

func (cluster *Cluster) Shutdown() {

}

func (cluster *Cluster) GetTargetNode(method, targetId string) *Node {
	if _, exist := cluster.currentNode.methodMap[method]; exist {
		return cluster.currentNode
	} else {
		return nil
	}
}

func (cluster *Cluster) GetTargetNodeCount(method string) int {
	return 1
}

func (cluster *Cluster) getNodeList(method string) []*Node {
	return []*Node{
		cluster.currentNode,
	}
}

func (cluster *Cluster) LocalActorOf(method string) actorsystem.ActorRef {
	return cluster.actorSystem.LocalActorOf(method)
}

func (cluster *Cluster) ActorOf(host string, port int, method string) actorsystem.ActorRef {
	return cluster.actorSystem.ActerOf(host, port, method)
}

func (cluster *Cluster) CallbackActorOf(ttl time.Duration, actor actorsystem.ICallbackUntypedActor) actorsystem.ActorRef {
	return cluster.actorSystem.CallbackActerOf(ttl, actor)
}

func (cluster *Cluster) UnicastRouteWithNoSender(method, targetId string, obj proto.Message) bool {
	return cluster.UnicastRoute(method, targetId, obj, actorsystem.NoSender)
}

func (cluster *Cluster) UnicastRoute(method, targetId string, obj proto.Message, sender actorsystem.ActorRef) bool {
	nod := cluster.GetTargetNode(method, targetId)
	if nod != nil {
		cluster.baseRoute(method, nod.Ip, nod.Port, obj, sender)
		return true
	}
	return false
}

func (cluster *Cluster) baseRoute(method string, host string, port int, obj proto.Message, sender actorsystem.ActorRef) {
	actor := cluster.actorSystem.ActerOf(host, port, method)
	actor.Tell(obj, sender)
}

func (cluster *Cluster) BroadcastWithNoSender(method string, obj proto.Message) {
	cluster.BroadcastRoute(method, obj, actorsystem.NoSender, []string{})
}

func (cluster *Cluster) BroadcastRoute(method string, obj proto.Message, sender actorsystem.ActorRef, excludeNotes []string) {
	excludeNode := map[string]bool{}
	for _, nodeName := range excludeNotes {
		excludeNode[nodeName] = true
	}
	nodes := cluster.getNodeList(method)
	for _, node := range nodes {
		if _, exist := excludeNode[node.Name]; exist {
			continue
		}
		cluster.baseRoute(method, node.Ip, node.Port, obj, sender)
	}
}

type Node struct {
	Name      string            `json:"name"`
	Ip        string            `json:"ip"`
	Port      int               `json:"port"`
	Methods   []string          `json:"methods"`
	methodMap map[string]bool   `json:"-"`
	Exts      map[string]string `json:"exts"`
}

func NewNode(name, ip string, port int, exts map[string]string) *Node {
	node := &Node{
		Name:      name,
		Ip:        ip,
		Port:      port,
		Methods:   []string{},
		methodMap: make(map[string]bool),
		Exts:      exts,
	}
	return node
}

func (node *Node) AddMethod(method string) {
	node.methodMap[method] = true
	methodArr := make([]string, 0, len(node.methodMap))
	for method := range node.methodMap {
		methodArr = append(methodArr, method)
	}
	node.Methods = methodArr
}
