package websocket2

//
//import "github.com/pkg/errors"
//
//type Process1 func(state State) (err error)
//
//func NewSyncroniser(state State) (sync *Synchroniser, err error) {
//	return nil, nil
//}
//
//type Synchroniser struct {
//	state State
//
//	processes map[string]Process1
//}
//
//func (s *Synchroniser) DefineProcess(name string, p Process1) (err error) {
//	if _, exists := s.processes[name]; exists {
//		err = errors.New("process '" + name + "' already exists")
//		return err
//	}
//	s.processes[name] = p
//	return nil
//}
//
//func (s *Synchroniser) Execute(name string) (err error) {
//	var p Process1
//	var exists bool
//	if p, exists = s.processes[name]; !exists {
//		err = errors.New("no process found with process name " + name)
//		return err
//	}
//
//	return s.Run(p)
//}
//
//func (s *Synchroniser) Run(p Process1) (err error) {
//	return nil
//}
