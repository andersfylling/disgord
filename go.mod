module github.com/andersfylling/disgord

go 1.13

require (
	github.com/andersfylling/djp v0.0.0-20190905223822-fbe0bb181ad8
	github.com/andersfylling/snowflake/v4 v4.0.2
	github.com/buger/jsonparser v0.0.0-20181115193947-bf1c66bbce23
	github.com/gorilla/websocket v1.4.1
	github.com/json-iterator/go v1.1.7
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pkg/errors v0.8.1
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.11.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/net v0.0.0-20191021144547-ec77196f6094
	golang.org/x/sys v0.0.0-20191024172528-b4ff53e7a1cb // indirect
	nhooyr.io/websocket v1.7.2
)

replace github.com/andersfylling/djp => ../djp
