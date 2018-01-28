package discordws

// type Hub map[string](chan []byte)
//
// func (h Hub) Publish(eventName string, p []byte) {
// 	if val, ok := h[eventName]; ok {
// 		data := make([]byte, len(p))
// 		copy(data, p)
//
// 		h[eventName] <- data
// 	}
// }
//
// func (h Hub) Subscribe(eventName string) chan []byte {
// 	if val, ok := h[eventName]; !ok {
// 		h[eventName] = make(chan []byte)
// 	}
//
// 	return h[eventName]
// }
