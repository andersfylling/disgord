package event

import (
	"github.com/andersfylling/disgord/disgordctx"
)

// socket
//

// HelloCallbackStack ***************
type HelloHandler interface {
	Add(cb HelloCallback) error
	Trigger(disgordctx.Context, *HelloBox) error
	ReceiveChan() <-chan *HelloBox
}

func NewHelloCallbackStack() *HelloCallbackStack {
	return &HelloCallbackStack{
		listener: make(chan *HelloBox),
	}
}

type HelloCallbackStack struct {
	sequential     bool
	listeners      []HelloCallback
	listenerExists bool
	listener       chan *HelloBox
}

func (stack *HelloCallbackStack) Add(cb HelloCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []HelloCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *HelloCallbackStack) Trigger(ctx disgordctx.Context, box *HelloBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *HelloCallbackStack) ReceiveChan() <-chan *HelloBox {
	return stack.listener
}

// ReadyCallbackStack ***************
type ReadyHandler interface {
	Add(cb ReadyCallback) error
	Trigger(disgordctx.Context, *ReadyBox) error
	ReceiveChan() <-chan *ReadyBox
}

func NewReadyCallbackStack() *ReadyCallbackStack {
	return &ReadyCallbackStack{
		listener: make(chan *ReadyBox),
	}
}

type ReadyCallbackStack struct {
	sequential     bool
	listeners      []ReadyCallback
	listenerExists bool
	listener       chan *ReadyBox
}

func (stack *ReadyCallbackStack) Add(cb ReadyCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []ReadyCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *ReadyCallbackStack) Trigger(ctx disgordctx.Context, box *ReadyBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}

	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *ReadyCallbackStack) ReceiveChan() <-chan *ReadyBox {
	stack.listenerExists = true
	return stack.listener
}

// ResumedCallbackStack **********
type ResumedHandler interface {
	Add(cb ResumedCallback) error
	Trigger(disgordctx.Context, *ResumedBox) error
	ReceiveChan() <-chan *ResumedBox
}

func NewResumedCallbackStack() *ResumedCallbackStack {
	return &ResumedCallbackStack{
		listener: make(chan *ResumedBox),
	}
}

type ResumedCallbackStack struct {
	sequential     bool
	listeners      []ResumedCallback
	listenerExists bool
	listener       chan *ResumedBox
}

func (stack *ResumedCallbackStack) Add(cb ResumedCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []ResumedCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *ResumedCallbackStack) Trigger(ctx disgordctx.Context, box *ResumedBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *ResumedCallbackStack) ReceiveChan() <-chan *ResumedBox {
	stack.listenerExists = true
	return stack.listener
}

// InvalidSessionCallbackStack ***************
type InvalidSessionHandler interface {
	Add(cb InvalidSessionCallback) error
	Trigger(disgordctx.Context, *InvalidSessionBox) error
	ReceiveChan() <-chan *InvalidSessionBox
}

func NewInvalidSessionCallbackStack() *InvalidSessionCallbackStack {
	return &InvalidSessionCallbackStack{
		listener: make(chan *InvalidSessionBox),
	}
}

type InvalidSessionCallbackStack struct {
	sequential     bool
	listeners      []InvalidSessionCallback
	listenerExists bool
	listener       chan *InvalidSessionBox
}

func (stack *InvalidSessionCallbackStack) Add(cb InvalidSessionCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []InvalidSessionCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *InvalidSessionCallbackStack) Trigger(ctx disgordctx.Context, box *InvalidSessionBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *InvalidSessionCallbackStack) ReceiveChan() <-chan *InvalidSessionBox {
	stack.listenerExists = true
	return stack.listener
}

// channel
//

// ChannelCreateCallbackStack **************
type ChannelCreateHandler interface {
	Add(ChannelCreateCallback) error
	Trigger(disgordctx.Context, *ChannelCreateBox) error
	ReceiveChan() <-chan *ChannelCreateBox
}

func NewChannelCreateCallbackStack() *ChannelCreateCallbackStack {
	return &ChannelCreateCallbackStack{
		listener: make(chan *ChannelCreateBox),
	}
}

type ChannelCreateCallbackStack struct {
	sequential     bool
	listeners      []ChannelCreateCallback
	listenerExists bool
	listener       chan *ChannelCreateBox
}

func (stack *ChannelCreateCallbackStack) Add(cb ChannelCreateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []ChannelCreateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *ChannelCreateCallbackStack) Trigger(ctx disgordctx.Context, box *ChannelCreateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *ChannelCreateCallbackStack) ReceiveChan() <-chan *ChannelCreateBox {
	stack.listenerExists = true
	return stack.listener
}

// ChannelUpdateCallbackStack ************
type ChannelUpdateHandler interface {
	Add(ChannelUpdateCallback) error
	Trigger(disgordctx.Context, *ChannelUpdateBox) error
	ReceiveChan() <-chan *ChannelUpdateBox
}

func NewChannelUpdateCallbackStack() *ChannelUpdateCallbackStack {
	return &ChannelUpdateCallbackStack{
		listener: make(chan *ChannelUpdateBox),
	}
}

type ChannelUpdateCallbackStack struct {
	sequential     bool
	listeners      []ChannelUpdateCallback
	listenerExists bool
	listener       chan *ChannelUpdateBox
}

func (stack *ChannelUpdateCallbackStack) Add(cb ChannelUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []ChannelUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *ChannelUpdateCallbackStack) Trigger(ctx disgordctx.Context, box *ChannelUpdateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *ChannelUpdateCallbackStack) ReceiveChan() <-chan *ChannelUpdateBox {
	stack.listenerExists = true
	return stack.listener
}

// ChannelDeleteCallbackStack ***********
type ChannelDeleteHandler interface {
	Add(ChannelDeleteCallback) error
	Trigger(disgordctx.Context, *ChannelDeleteBox) error
	ReceiveChan() <-chan *ChannelDeleteBox
}

func NewChannelDeleteCallbackStack() *ChannelDeleteCallbackStack {
	return &ChannelDeleteCallbackStack{
		listener: make(chan *ChannelDeleteBox),
	}
}

type ChannelDeleteCallbackStack struct {
	sequential     bool
	listeners      []ChannelDeleteCallback
	listenerExists bool
	listener       chan *ChannelDeleteBox
}

func (stack *ChannelDeleteCallbackStack) Add(cb ChannelDeleteCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []ChannelDeleteCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *ChannelDeleteCallbackStack) Trigger(ctx disgordctx.Context, box *ChannelDeleteBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *ChannelDeleteCallbackStack) ReceiveChan() <-chan *ChannelDeleteBox {
	stack.listenerExists = true
	return stack.listener
}

// ChannelPinsUpdateCallbackStack **********
type ChannelPinsUpdateHandler interface {
	Add(ChannelPinsUpdateCallback) error
	Trigger(disgordctx.Context, *ChannelPinsUpdateBox) error
	ReceiveChan() <-chan *ChannelPinsUpdateBox
}

func NewChannelPinsUpdateCallbackStack() *ChannelPinsUpdateCallbackStack {
	return &ChannelPinsUpdateCallbackStack{
		listener: make(chan *ChannelPinsUpdateBox),
	}
}

type ChannelPinsUpdateCallbackStack struct {
	sequential     bool
	listeners      []ChannelPinsUpdateCallback
	listenerExists bool
	listener       chan *ChannelPinsUpdateBox
}

func (stack *ChannelPinsUpdateCallbackStack) Add(cb ChannelPinsUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []ChannelPinsUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *ChannelPinsUpdateCallbackStack) Trigger(ctx disgordctx.Context, box *ChannelPinsUpdateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *ChannelPinsUpdateCallbackStack) ReceiveChan() <-chan *ChannelPinsUpdateBox {
	stack.listenerExists = true
	return stack.listener
}

// Guild in general
//

// GuildCreateCallbackStack **********
type GuildCreateHandler interface {
	Add(GuildCreateCallback) error
	Trigger(disgordctx.Context, *GuildCreateBox) error
	ReceiveChan() <-chan *GuildCreateBox
}

func NewGuildCreateCallbackStack() *GuildCreateCallbackStack {
	return &GuildCreateCallbackStack{
		listener: make(chan *GuildCreateBox),
	}
}

type GuildCreateCallbackStack struct {
	sequential     bool
	listeners      []GuildCreateCallback
	listenerExists bool
	listener       chan *GuildCreateBox
}

func (stack *GuildCreateCallbackStack) Add(cb GuildCreateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildCreateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildCreateCallbackStack) Trigger(ctx disgordctx.Context, box *GuildCreateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *GuildCreateCallbackStack) ReceiveChan() <-chan *GuildCreateBox {
	stack.listenerExists = true
	return stack.listener
}

// GuildUpdateCallbackStack .....
type GuildUpdateHandler interface {
	Add(GuildUpdateCallback) error
	Trigger(disgordctx.Context, *GuildUpdateBox) error
	ReceiveChan() <-chan *GuildUpdateBox
}

func NewGuildUpdateCallbackStack() *GuildUpdateCallbackStack {
	return &GuildUpdateCallbackStack{
		listener: make(chan *GuildUpdateBox),
	}
}

type GuildUpdateCallbackStack struct {
	sequential     bool
	listeners      []GuildUpdateCallback
	listenerExists bool
	listener       chan *GuildUpdateBox
}

func (stack *GuildUpdateCallbackStack) Add(cb GuildUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildUpdateCallbackStack) Trigger(ctx disgordctx.Context, box *GuildUpdateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *GuildUpdateCallbackStack) ReceiveChan() <-chan *GuildUpdateBox {
	stack.listenerExists = true
	return stack.listener
}

// GuildDeleteCallbackStack *********
type GuildDeleteHandler interface {
	Add(GuildDeleteCallback) error
	Trigger(disgordctx.Context, *GuildDeleteBox) error
	ReceiveChan() <-chan *GuildDeleteBox
}

func NewGuildDeleteCallbackStack() *GuildDeleteCallbackStack {
	return &GuildDeleteCallbackStack{
		listener: make(chan *GuildDeleteBox),
	}
}

type GuildDeleteCallbackStack struct {
	sequential     bool
	listeners      []GuildDeleteCallback
	listenerExists bool
	listener       chan *GuildDeleteBox
}

func (stack *GuildDeleteCallbackStack) Add(cb GuildDeleteCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildDeleteCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildDeleteCallbackStack) Trigger(ctx disgordctx.Context, box *GuildDeleteBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *GuildDeleteCallbackStack) ReceiveChan() <-chan *GuildDeleteBox {
	stack.listenerExists = true
	return stack.listener
}

// GuildBanAddCallbackStack **************
type GuildBanAddHandler interface {
	Add(GuildBanAddCallback) error
	Trigger(disgordctx.Context, *GuildBanAddBox) error
	ReceiveChan() <-chan *GuildBanAddBox
}

func NewGuildBanAddCallbackStack() *GuildBanAddCallbackStack {
	return &GuildBanAddCallbackStack{
		listener: make(chan *GuildBanAddBox),
	}
}

type GuildBanAddCallbackStack struct {
	sequential     bool
	listeners      []GuildBanAddCallback
	listenerExists bool
	listener       chan *GuildBanAddBox
}

func (stack *GuildBanAddCallbackStack) Add(cb GuildBanAddCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildBanAddCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildBanAddCallbackStack) Trigger(ctx disgordctx.Context, box *GuildBanAddBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *GuildBanAddCallbackStack) ReceiveChan() <-chan *GuildBanAddBox {
	stack.listenerExists = true
	return stack.listener
}

// GuildBanRemoveCallbackStack *********
type GuildBanRemoveHandler interface {
	Add(GuildBanRemoveCallback) error
	Trigger(disgordctx.Context, *GuildBanRemoveBox) error
	ReceiveChan() <-chan *GuildBanRemoveBox
}

func NewGuildBanRemoveCallbackStack() *GuildBanRemoveCallbackStack {
	return &GuildBanRemoveCallbackStack{
		listener: make(chan *GuildBanRemoveBox),
	}
}

type GuildBanRemoveCallbackStack struct {
	sequential     bool
	listeners      []GuildBanRemoveCallback
	listenerExists bool
	listener       chan *GuildBanRemoveBox
}

func (stack *GuildBanRemoveCallbackStack) Add(cb GuildBanRemoveCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildBanRemoveCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildBanRemoveCallbackStack) Trigger(ctx disgordctx.Context, box *GuildBanRemoveBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *GuildBanRemoveCallbackStack) ReceiveChan() <-chan *GuildBanRemoveBox {
	stack.listenerExists = true
	return stack.listener
}

// GuildEmojisUpdateCallbackStack ***********
type GuildEmojisUpdateHandler interface {
	Add(GuildEmojisUpdateCallback) error
	Trigger(disgordctx.Context, *GuildEmojisUpdateBox) error
	ReceiveChan() <-chan *GuildEmojisUpdateBox
}

func NewGuildEmojisUpdateCallbackStack() *GuildEmojisUpdateCallbackStack {
	return &GuildEmojisUpdateCallbackStack{
		listener: make(chan *GuildEmojisUpdateBox),
	}
}

type GuildEmojisUpdateCallbackStack struct {
	sequential     bool
	listeners      []GuildEmojisUpdateCallback
	listenerExists bool
	listener       chan *GuildEmojisUpdateBox
}

func (stack *GuildEmojisUpdateCallbackStack) Add(cb GuildEmojisUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildEmojisUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildEmojisUpdateCallbackStack) Trigger(ctx disgordctx.Context, box *GuildEmojisUpdateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *GuildEmojisUpdateCallbackStack) ReceiveChan() <-chan *GuildEmojisUpdateBox {
	stack.listenerExists = true
	return stack.listener
}

// GuildIntegrationsUpdateCallbackStack *******************
type GuildIntegrationsUpdateHandler interface {
	Add(GuildIntegrationsUpdateCallback) error
	Trigger(disgordctx.Context, *GuildIntegrationsUpdateBox) error
	ReceiveChan() <-chan *GuildIntegrationsUpdateBox
}

func NewGuildIntegrationsUpdateCallbackStack() *GuildIntegrationsUpdateCallbackStack {
	return &GuildIntegrationsUpdateCallbackStack{
		listener: make(chan *GuildIntegrationsUpdateBox),
	}
}

type GuildIntegrationsUpdateCallbackStack struct {
	sequential     bool
	listeners      []GuildIntegrationsUpdateCallback
	listenerExists bool
	listener       chan *GuildIntegrationsUpdateBox
}

func (stack *GuildIntegrationsUpdateCallbackStack) Add(cb GuildIntegrationsUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildIntegrationsUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildIntegrationsUpdateCallbackStack) Trigger(ctx disgordctx.Context, box *GuildIntegrationsUpdateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *GuildIntegrationsUpdateCallbackStack) ReceiveChan() <-chan *GuildIntegrationsUpdateBox {
	stack.listenerExists = true
	return stack.listener
}

// Guild Member
//

// GuildMemberAddCallbackStack ***********************
type GuildMemberAddHandler interface {
	Add(GuildMemberAddCallback) error
	Trigger(disgordctx.Context, *GuildMemberAddBox) error
	ReceiveChan() <-chan *GuildMemberAddBox
}

func NewGuildMemberAddCallbackStack() *GuildMemberAddCallbackStack {
	return &GuildMemberAddCallbackStack{
		listener: make(chan *GuildMemberAddBox),
	}
}

type GuildMemberAddCallbackStack struct {
	sequential     bool
	listeners      []GuildMemberAddCallback
	listenerExists bool
	listener       chan *GuildMemberAddBox
}

func (stack *GuildMemberAddCallbackStack) Add(cb GuildMemberAddCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildMemberAddCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildMemberAddCallbackStack) Trigger(ctx disgordctx.Context, box *GuildMemberAddBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *GuildMemberAddCallbackStack) ReceiveChan() <-chan *GuildMemberAddBox {
	stack.listenerExists = true
	return stack.listener
}

// GuildMemberRemoveCallbackStack *******************
type GuildMemberRemoveHandler interface {
	Add(GuildMemberRemoveCallback) error
	Trigger(disgordctx.Context, *GuildMemberRemoveBox) error
	ReceiveChan() <-chan *GuildMemberRemoveBox
}

func NewGuildMemberRemoveCallbackStack() *GuildMemberRemoveCallbackStack {
	return &GuildMemberRemoveCallbackStack{
		listener: make(chan *GuildMemberRemoveBox),
	}
}

type GuildMemberRemoveCallbackStack struct {
	sequential     bool
	listeners      []GuildMemberRemoveCallback
	listenerExists bool
	listener       chan *GuildMemberRemoveBox
}

func (stack *GuildMemberRemoveCallbackStack) Add(cb GuildMemberRemoveCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildMemberRemoveCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildMemberRemoveCallbackStack) Trigger(ctx disgordctx.Context, box *GuildMemberRemoveBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *GuildMemberRemoveCallbackStack) ReceiveChan() <-chan *GuildMemberRemoveBox {
	stack.listenerExists = true
	return stack.listener
}

// GuildMemberUpdateCallbackStack **************
type GuildMemberUpdateHandler interface {
	Add(GuildMemberUpdateCallback) error
	Trigger(disgordctx.Context, *GuildMemberUpdateBox) error
	ReceiveChan() <-chan *GuildMemberUpdateBox
}

func NewGuildMemberUpdateCallbackStack() *GuildMemberUpdateCallbackStack {
	return &GuildMemberUpdateCallbackStack{
		listener: make(chan *GuildMemberUpdateBox),
	}
}

type GuildMemberUpdateCallbackStack struct {
	sequential     bool
	listeners      []GuildMemberUpdateCallback
	listenerExists bool
	listener       chan *GuildMemberUpdateBox
}

func (stack *GuildMemberUpdateCallbackStack) Add(cb GuildMemberUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildMemberUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildMemberUpdateCallbackStack) Trigger(ctx disgordctx.Context, box *GuildMemberUpdateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *GuildMemberUpdateCallbackStack) ReceiveChan() <-chan *GuildMemberUpdateBox {
	stack.listenerExists = true
	return stack.listener
}

// GuildMemberChunkCallbackStack **************
type GuildMembersChunkHandler interface {
	Add(GuildMembersChunkCallback) error
	Trigger(disgordctx.Context, *GuildMembersChunkBox) error
	ReceiveChan() <-chan *GuildMembersChunkBox
}

func NewGuildMembersChunkCallbackStack() *GuildMembersChunkCallbackStack {
	return &GuildMembersChunkCallbackStack{
		listener: make(chan *GuildMembersChunkBox),
	}
}

type GuildMembersChunkCallbackStack struct {
	sequential     bool
	listeners      []GuildMembersChunkCallback
	listenerExists bool
	listener       chan *GuildMembersChunkBox
}

func (stack *GuildMembersChunkCallbackStack) Add(cb GuildMembersChunkCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildMembersChunkCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildMembersChunkCallbackStack) Trigger(ctx disgordctx.Context, box *GuildMembersChunkBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *GuildMembersChunkCallbackStack) ReceiveChan() <-chan *GuildMembersChunkBox {
	stack.listenerExists = true
	return stack.listener
}

// Guild role
//

// GuildRoleCreateCallbackStack *************
type GuildRoleCreateHandler interface {
	Add(GuildRoleCreateCallback) error
	Trigger(disgordctx.Context, *GuildRoleCreateBox) error
	ReceiveChan() <-chan *GuildRoleCreateBox
}

func NewGuildRoleCreateCallbackStack() *GuildRoleCreateCallbackStack {
	return &GuildRoleCreateCallbackStack{
		listener: make(chan *GuildRoleCreateBox),
	}
}

type GuildRoleCreateCallbackStack struct {
	sequential     bool
	listeners      []GuildRoleCreateCallback
	listenerExists bool
	listener       chan *GuildRoleCreateBox
}

func (stack *GuildRoleCreateCallbackStack) Add(cb GuildRoleCreateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildRoleCreateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildRoleCreateCallbackStack) Trigger(ctx disgordctx.Context, box *GuildRoleCreateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *GuildRoleCreateCallbackStack) ReceiveChan() <-chan *GuildRoleCreateBox {
	stack.listenerExists = true
	return stack.listener
}

// GuildRoleUpdateCallbackStack ***************
type GuildRoleUpdateHandler interface {
	Add(GuildRoleUpdateCallback) error
	Trigger(disgordctx.Context, *GuildRoleUpdateBox) error
	ReceiveChan() <-chan *GuildRoleUpdateBox
}

func NewGuildRoleUpdateCallbackStack() *GuildRoleUpdateCallbackStack {
	return &GuildRoleUpdateCallbackStack{
		listener: make(chan *GuildRoleUpdateBox),
	}
}

type GuildRoleUpdateCallbackStack struct {
	sequential     bool
	listeners      []GuildRoleUpdateCallback
	listenerExists bool
	listener       chan *GuildRoleUpdateBox
}

func (stack *GuildRoleUpdateCallbackStack) Add(cb GuildRoleUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildRoleUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildRoleUpdateCallbackStack) Trigger(ctx disgordctx.Context, box *GuildRoleUpdateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *GuildRoleUpdateCallbackStack) ReceiveChan() <-chan *GuildRoleUpdateBox {
	stack.listenerExists = true
	return stack.listener
}

// GuildRoleDeleteCallbackStack **************
type GuildRoleDeleteHandler interface {
	Add(GuildRoleDeleteCallback) error
	Trigger(disgordctx.Context, *GuildRoleDeleteBox) error
	ReceiveChan() <-chan *GuildRoleDeleteBox
}

func NewGuildRoleDeleteCallbackStack() *GuildRoleDeleteCallbackStack {
	return &GuildRoleDeleteCallbackStack{
		listener: make(chan *GuildRoleDeleteBox),
	}
}

type GuildRoleDeleteCallbackStack struct {
	sequential     bool
	listeners      []GuildRoleDeleteCallback
	listenerExists bool
	listener       chan *GuildRoleDeleteBox
}

func (stack *GuildRoleDeleteCallbackStack) Add(cb GuildRoleDeleteCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []GuildRoleDeleteCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *GuildRoleDeleteCallbackStack) Trigger(ctx disgordctx.Context, box *GuildRoleDeleteBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *GuildRoleDeleteCallbackStack) ReceiveChan() <-chan *GuildRoleDeleteBox {
	stack.listenerExists = true
	return stack.listener
}

// message
//

// MessageCreateCallbackStack ********************
type MessageCreateHandler interface {
	Add(MessageCreateCallback) error
	Trigger(disgordctx.Context, *MessageCreateBox) error
	ReceiveChan() <-chan *MessageCreateBox
}

func NewMessageCreateCallbackStack() *MessageCreateCallbackStack {
	return &MessageCreateCallbackStack{
		listener: make(chan *MessageCreateBox),
	}
}

type MessageCreateCallbackStack struct {
	sequential     bool
	listeners      []MessageCreateCallback
	listenerExists bool
	listener       chan *MessageCreateBox
}

func (stack *MessageCreateCallbackStack) Add(cb MessageCreateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []MessageCreateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *MessageCreateCallbackStack) Trigger(ctx disgordctx.Context, box *MessageCreateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *MessageCreateCallbackStack) ReceiveChan() <-chan *MessageCreateBox {
	stack.listenerExists = true
	return stack.listener
}

// MessageUpdateCallbackStack ****************
type MessageUpdateHandler interface {
	Add(MessageUpdateCallback) error
	Trigger(disgordctx.Context, *MessageUpdateBox) error
	ReceiveChan() <-chan *MessageUpdateBox
}

func NewMessageUpdateCallbackStack() *MessageUpdateCallbackStack {
	return &MessageUpdateCallbackStack{
		listener: make(chan *MessageUpdateBox),
	}
}

type MessageUpdateCallbackStack struct {
	sequential     bool
	listeners      []MessageUpdateCallback
	listenerExists bool
	listener       chan *MessageUpdateBox
}

func (stack *MessageUpdateCallbackStack) Add(cb MessageUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []MessageUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *MessageUpdateCallbackStack) Trigger(ctx disgordctx.Context, box *MessageUpdateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *MessageUpdateCallbackStack) ReceiveChan() <-chan *MessageUpdateBox {
	stack.listenerExists = true
	return stack.listener
}

// MessageDeleteCallbackStack ***************
type MessageDeleteHandler interface {
	Add(MessageDeleteCallback) error
	Trigger(disgordctx.Context, *MessageDeleteBox) error
	ReceiveChan() <-chan *MessageDeleteBox
}

func NewMessageDeleteCallbackStack() *MessageDeleteCallbackStack {
	return &MessageDeleteCallbackStack{
		listener: make(chan *MessageDeleteBox),
	}
}

type MessageDeleteCallbackStack struct {
	sequential     bool
	listeners      []MessageDeleteCallback
	listenerExists bool
	listener       chan *MessageDeleteBox
}

func (stack *MessageDeleteCallbackStack) Add(cb MessageDeleteCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []MessageDeleteCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *MessageDeleteCallbackStack) Trigger(ctx disgordctx.Context, box *MessageDeleteBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *MessageDeleteCallbackStack) ReceiveChan() <-chan *MessageDeleteBox {
	stack.listenerExists = true
	return stack.listener
}

// MessageDeleteBulkCallbackStack ****************
type MessageDeleteBulkHandler interface {
	Add(MessageDeleteBulkCallback) error
	Trigger(disgordctx.Context, *MessageDeleteBulkBox) error
	ReceiveChan() <-chan *MessageDeleteBulkBox
}

func NewMessageDeleteBulkCallbackStack() *MessageDeleteBulkCallbackStack {
	return &MessageDeleteBulkCallbackStack{
		listener: make(chan *MessageDeleteBulkBox),
	}
}

type MessageDeleteBulkCallbackStack struct {
	sequential     bool
	listeners      []MessageDeleteBulkCallback
	listenerExists bool
	listener       chan *MessageDeleteBulkBox
}

func (stack *MessageDeleteBulkCallbackStack) Add(cb MessageDeleteBulkCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []MessageDeleteBulkCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *MessageDeleteBulkCallbackStack) Trigger(ctx disgordctx.Context, box *MessageDeleteBulkBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *MessageDeleteBulkCallbackStack) ReceiveChan() <-chan *MessageDeleteBulkBox {
	stack.listenerExists = true
	return stack.listener
}

// message reaction
//

// MessageReactionAddCallbackStack ************
type MessageReactionAddHandler interface {
	Add(MessageReactionAddCallback) error
	Trigger(disgordctx.Context, *MessageReactionAddBox) error
	ReceiveChan() <-chan *MessageReactionAddBox
}

func NewMessageReactionAddCallbackStack() *MessageReactionAddCallbackStack {
	return &MessageReactionAddCallbackStack{
		listener: make(chan *MessageReactionAddBox),
	}
}

type MessageReactionAddCallbackStack struct {
	sequential     bool
	listeners      []MessageReactionAddCallback
	listenerExists bool
	listener       chan *MessageReactionAddBox
}

func (stack *MessageReactionAddCallbackStack) Add(cb MessageReactionAddCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []MessageReactionAddCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *MessageReactionAddCallbackStack) Trigger(ctx disgordctx.Context, box *MessageReactionAddBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *MessageReactionAddCallbackStack) ReceiveChan() <-chan *MessageReactionAddBox {
	stack.listenerExists = true
	return stack.listener
}

// MessageReactionRemoveCallbackStack *********
type MessageReactionRemoveHandler interface {
	Add(MessageReactionRemoveCallback) error
	Trigger(disgordctx.Context, *MessageReactionRemoveBox) error
	ReceiveChan() <-chan *MessageReactionRemoveBox
}

func NewMessageReactionRemoveCallbackStack() *MessageReactionRemoveCallbackStack {
	return &MessageReactionRemoveCallbackStack{
		listener: make(chan *MessageReactionRemoveBox),
	}
}

type MessageReactionRemoveCallbackStack struct {
	sequential     bool
	listeners      []MessageReactionRemoveCallback
	listenerExists bool
	listener       chan *MessageReactionRemoveBox
}

func (stack *MessageReactionRemoveCallbackStack) Add(cb MessageReactionRemoveCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []MessageReactionRemoveCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *MessageReactionRemoveCallbackStack) Trigger(ctx disgordctx.Context, box *MessageReactionRemoveBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *MessageReactionRemoveCallbackStack) ReceiveChan() <-chan *MessageReactionRemoveBox {
	stack.listenerExists = true
	return stack.listener
}

// MessageReactionRemoveAllCallbackStack *********
type MessageReactionRemoveAllHandler interface {
	Add(MessageReactionRemoveAllCallback) error
	Trigger(disgordctx.Context, *MessageReactionRemoveAllBox) error
	ReceiveChan() <-chan *MessageReactionRemoveAllBox
}

func NewMessageReactionRemoveAllCallbackStack() *MessageReactionRemoveAllCallbackStack {
	return &MessageReactionRemoveAllCallbackStack{
		listener: make(chan *MessageReactionRemoveAllBox),
	}
}

type MessageReactionRemoveAllCallbackStack struct {
	sequential     bool
	listeners      []MessageReactionRemoveAllCallback
	listenerExists bool
	listener       chan *MessageReactionRemoveAllBox
}

func (stack *MessageReactionRemoveAllCallbackStack) Add(cb MessageReactionRemoveAllCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []MessageReactionRemoveAllCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *MessageReactionRemoveAllCallbackStack) Trigger(ctx disgordctx.Context, box *MessageReactionRemoveAllBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *MessageReactionRemoveAllCallbackStack) ReceiveChan() <-chan *MessageReactionRemoveAllBox {
	stack.listenerExists = true
	return stack.listener
}

// presence
//

// PresenceUpdateCallbackStack *************
type PresenceUpdateHandler interface {
	Add(PresenceUpdateCallback) error
	Trigger(disgordctx.Context, *PresenceUpdateBox) error
	ReceiveChan() <-chan *PresenceUpdateBox
}

func NewPresenceUpdateCallbackStack() *PresenceUpdateCallbackStack {
	return &PresenceUpdateCallbackStack{
		listener: make(chan *PresenceUpdateBox),
	}
}

type PresenceUpdateCallbackStack struct {
	sequential     bool
	listeners      []PresenceUpdateCallback
	listenerExists bool
	listener       chan *PresenceUpdateBox
}

func (stack *PresenceUpdateCallbackStack) Add(cb PresenceUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []PresenceUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *PresenceUpdateCallbackStack) Trigger(ctx disgordctx.Context, box *PresenceUpdateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *PresenceUpdateCallbackStack) ReceiveChan() <-chan *PresenceUpdateBox {
	stack.listenerExists = true
	return stack.listener
}

// typing start
//

// TypingStartCallbackStack ******************
type TypingStartHandler interface {
	Add(TypingStartCallback) error
	Trigger(disgordctx.Context, *TypingStartBox) error
	ReceiveChan() <-chan *TypingStartBox
}

func NewTypingStartCallbackStack() *TypingStartCallbackStack {
	return &TypingStartCallbackStack{
		listener: make(chan *TypingStartBox),
	}
}

type TypingStartCallbackStack struct {
	sequential     bool
	listeners      []TypingStartCallback
	listenerExists bool
	listener       chan *TypingStartBox
}

func (stack *TypingStartCallbackStack) Add(cb TypingStartCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []TypingStartCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *TypingStartCallbackStack) Trigger(ctx disgordctx.Context, box *TypingStartBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *TypingStartCallbackStack) ReceiveChan() <-chan *TypingStartBox {
	stack.listenerExists = true
	return stack.listener
}

// user update
type UserUpdateHandler interface {
	Add(UserUpdateCallback) error
	Trigger(disgordctx.Context, *UserUpdateBox) error
	ReceiveChan() <-chan *UserUpdateBox
}

func NewUserUpdateCallbackStack() *UserUpdateCallbackStack {
	return &UserUpdateCallbackStack{
		listener: make(chan *UserUpdateBox),
	}
}

type UserUpdateCallbackStack struct {
	sequential     bool
	listeners      []UserUpdateCallback
	listenerExists bool
	listener       chan *UserUpdateBox
}

func (stack *UserUpdateCallbackStack) Add(cb UserUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []UserUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *UserUpdateCallbackStack) Trigger(ctx disgordctx.Context, box *UserUpdateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *UserUpdateCallbackStack) ReceiveChan() <-chan *UserUpdateBox {
	stack.listenerExists = true
	return stack.listener
}

// voice
//

// VoiceStateUpdateCallbackStack *************************
type VoiceStateUpdateHandler interface {
	Add(VoiceStateUpdateCallback) error
	Trigger(disgordctx.Context, *VoiceStateUpdateBox) error
	ReceiveChan() <-chan *VoiceStateUpdateBox
}

func NewVoiceStateUpdateCallbackStack() *VoiceStateUpdateCallbackStack {
	return &VoiceStateUpdateCallbackStack{
		listener: make(chan *VoiceStateUpdateBox),
	}
}

type VoiceStateUpdateCallbackStack struct {
	sequential     bool
	listeners      []VoiceStateUpdateCallback
	listenerExists bool
	listener       chan *VoiceStateUpdateBox
}

func (stack *VoiceStateUpdateCallbackStack) Add(cb VoiceStateUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []VoiceStateUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *VoiceStateUpdateCallbackStack) Trigger(ctx disgordctx.Context, box *VoiceStateUpdateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *VoiceStateUpdateCallbackStack) ReceiveChan() <-chan *VoiceStateUpdateBox {
	stack.listenerExists = true
	return stack.listener
}

// VoiceServerUpdateCallbackStack ***********************
type VoiceServerUpdateHandler interface {
	Add(VoiceServerUpdateCallback) error
	Trigger(disgordctx.Context, *VoiceServerUpdateBox) error
	ReceiveChan() <-chan *VoiceServerUpdateBox
}

func NewVoiceServerUpdateCallbackStack() *VoiceServerUpdateCallbackStack {
	return &VoiceServerUpdateCallbackStack{
		listener: make(chan *VoiceServerUpdateBox),
	}
}

type VoiceServerUpdateCallbackStack struct {
	sequential     bool
	listeners      []VoiceServerUpdateCallback
	listenerExists bool
	listener       chan *VoiceServerUpdateBox
}

func (stack *VoiceServerUpdateCallbackStack) Add(cb VoiceServerUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []VoiceServerUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *VoiceServerUpdateCallbackStack) Trigger(ctx disgordctx.Context, box *VoiceServerUpdateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *VoiceServerUpdateCallbackStack) ReceiveChan() <-chan *VoiceServerUpdateBox {
	stack.listenerExists = true
	return stack.listener
}

// WebhooksUpdateCallbackStack *******************
type WebhooksUpdateHandler interface {
	Add(cb WebhooksUpdateCallback) error
	Trigger(disgordctx.Context, *WebhooksUpdateBox) error
	ReceiveChan() <-chan *WebhooksUpdateBox
}

func NewWebhooksUpdateCallbackStack() *WebhooksUpdateCallbackStack {
	return &WebhooksUpdateCallbackStack{
		listener: make(chan *WebhooksUpdateBox),
	}
}

type WebhooksUpdateCallbackStack struct {
	sequential     bool
	listeners      []WebhooksUpdateCallback
	listenerExists bool
	listener       chan *WebhooksUpdateBox
}

func (stack *WebhooksUpdateCallbackStack) Add(cb WebhooksUpdateCallback) (err error) {
	if stack.listeners == nil {
		stack.listeners = []WebhooksUpdateCallback{}
	}

	stack.listeners = append(stack.listeners, cb)
	return nil
}

func (stack *WebhooksUpdateCallbackStack) Trigger(ctx disgordctx.Context, box *WebhooksUpdateBox) (err error) {
	if stack.listenerExists {
		stack.listener <- box
	}
	for _, listener := range stack.listeners {
		if stack.sequential {
			listener(ctx, box)
		} else {
			go listener(ctx, box)
		}
	}

	return nil
}

func (stack *WebhooksUpdateCallbackStack) ReceiveChan() <-chan *WebhooksUpdateBox {
	stack.listenerExists = true
	return stack.listener
}
