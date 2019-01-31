package websocket2

//
//import "fmt"
//
//type State interface {
//	fmt.Stringer
//	Set(state interface{})
//	Equal(state interface{}) bool
//}
//
//const (
//	StateDisconnected DefaultStateCtrl = iota
//	StateConnected
//)
//
//type DefaultStateCtrl int
//
//func (s *DefaultStateCtrl) Set(state interface{}) {
//	var change *DefaultStateCtrl
//	var ok bool
//	if change, ok = state.(*DefaultStateCtrl); !ok {
//		return
//	}
//
//	*s = *change
//}
//func (s *DefaultStateCtrl) Equal(state interface{}) bool {
//	var s2 *DefaultStateCtrl
//	var ok bool
//	if s2, ok = state.(*DefaultStateCtrl); !ok {
//		return false
//	}
//
//	// compare content and not the pointer val
//	return int(*s2) == int(*s)
//}
//
//// TODO: go generate
//func (s *DefaultStateCtrl) String() string {
//	return "TODO"
//}
//
//var _ State = (*DefaultStateCtrl)(nil)
