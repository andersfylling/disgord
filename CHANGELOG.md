# Change Log

## [v0.12.2](https://github.com/andersfylling/disgord/tree/v0.12.2)

[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.12.1...v0.12.2)

**Fixed bugs:**

- Client.CreateBotURL panics in powershell [\#233](https://github.com/andersfylling/disgord/issues/233)

## [v0.12.1](https://github.com/andersfylling/disgord/tree/v0.12.1) (2019-10-24)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.12.0...v0.12.1)

**Fixed bugs:**

- emitter uses continue in select statement [\#230](https://github.com/andersfylling/disgord/issues/230)

**Closed issues:**

- increase timeout for queue checks in gateway emitter [\#231](https://github.com/andersfylling/disgord/issues/231)

## [v0.12.0](https://github.com/andersfylling/disgord/tree/v0.12.0) (2019-10-22)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.12.0-rc8...v0.12.0)

**Fixed bugs:**

- Member is missing internalUpdate [\#193](https://github.com/andersfylling/disgord/issues/193)
- Guild.LoadAllMembers does not check for duplicates [\#190](https://github.com/andersfylling/disgord/issues/190)

**Closed issues:**

- replace shutdown channel with context.Context [\#169](https://github.com/andersfylling/disgord/issues/169)
- "not by bot" filter for reactions [\#157](https://github.com/andersfylling/disgord/issues/157)

**Merged pull requests:**

- add config option LoadMembersQuietly [\#229](https://github.com/andersfylling/disgord/pull/229) ([andersfylling](https://github.com/andersfylling))
- Replace Guild.LoadAllMembers with more suitable Client.LoadMembers [\#221](https://github.com/andersfylling/disgord/pull/221) ([paulhobbel](https://github.com/paulhobbel))

## [v0.12.0-rc8](https://github.com/andersfylling/disgord/tree/v0.12.0-rc8) (2019-10-21)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.12.0-rc7...v0.12.0-rc8)

**Implemented enhancements:**

- Replace uber dep for atomic with sync.atomic [\#213](https://github.com/andersfylling/disgord/issues/213)
- Add a convenience method for reactions [\#211](https://github.com/andersfylling/disgord/issues/211)
- improve docs about snowflake.Snowflake vs disgord.Snowflake [\#161](https://github.com/andersfylling/disgord/issues/161)
- Better support for distributed instances [\#224](https://github.com/andersfylling/disgord/issues/224)
- create workflow that verifies install scipt works on push to develop [\#170](https://github.com/andersfylling/disgord/issues/170)
- dynamic buckets + option to inject custom system [\#173](https://github.com/andersfylling/disgord/pull/173) ([andersfylling](https://github.com/andersfylling))

**Fixed bugs:**

- Remove depalias [\#214](https://github.com/andersfylling/disgord/issues/214)
- Event Guild Members Chunk does not update members count [\#128](https://github.com/andersfylling/disgord/issues/128)

**Closed issues:**

- Sharded caching [\#183](https://github.com/andersfylling/disgord/issues/183)
- replace cache strategy with TLFU \(Time aware Least Frequently Used\) [\#180](https://github.com/andersfylling/disgord/issues/180)
- Standardise error types [\#178](https://github.com/andersfylling/disgord/issues/178)
- add new message fields [\#159](https://github.com/andersfylling/disgord/issues/159)

**Merged pull requests:**

- allow injecting custom identify rate limiter [\#227](https://github.com/andersfylling/disgord/pull/227) ([andersfylling](https://github.com/andersfylling))
- copy only config + doc update [\#226](https://github.com/andersfylling/disgord/pull/226) ([andersfylling](https://github.com/andersfylling))

## [v0.12.0-rc7](https://github.com/andersfylling/disgord/tree/v0.12.0-rc7) (2019-10-13)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.12.0-rc6...v0.12.0-rc7)

**Implemented enhancements:**

- Fix voice options not being used \(self-mute & deafen\) [\#218](https://github.com/andersfylling/disgord/pull/218) ([ikkerens](https://github.com/ikkerens))

**Merged pull requests:**

- Move private pkgs to internal pkg [\#223](https://github.com/andersfylling/disgord/pull/223) ([andersfylling](https://github.com/andersfylling))

## [v0.12.0-rc6](https://github.com/andersfylling/disgord/tree/v0.12.0-rc6) (2019-09-28)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.12.0-rc5...v0.12.0-rc6)

**Implemented enhancements:**

- Use slice in Request Guild Members Command [\#210](https://github.com/andersfylling/disgord/pull/210) ([andersfylling](https://github.com/andersfylling))

**Fixed bugs:**

- UpdateStatus while not connected will silently do nothing [\#209](https://github.com/andersfylling/disgord/issues/209)
- disgord does not complain about unknown handler signatures [\#208](https://github.com/andersfylling/disgord/issues/208)
- Use slice in Request Guild Members Command [\#210](https://github.com/andersfylling/disgord/pull/210) ([andersfylling](https://github.com/andersfylling))

**Merged pull requests:**

- panic when registerring a incorrect handler signature \(fixes \#20… [\#217](https://github.com/andersfylling/disgord/pull/217) ([andersfylling](https://github.com/andersfylling))
- upgrade websocket/nhooyr to fix atomic panic on ARM systems [\#216](https://github.com/andersfylling/disgord/pull/216) ([andersfylling](https://github.com/andersfylling))
- detects premature Emit usage \(fixes \#209\) [\#215](https://github.com/andersfylling/disgord/pull/215) ([andersfylling](https://github.com/andersfylling))
- Removed circle ci [\#212](https://github.com/andersfylling/disgord/pull/212) ([svenwiltink](https://github.com/svenwiltink))
- upgrade deps [\#206](https://github.com/andersfylling/disgord/pull/206) ([andersfylling](https://github.com/andersfylling))
- Some grammar changes/fixes, more to come [\#203](https://github.com/andersfylling/disgord/pull/203) ([GreemDev](https://github.com/GreemDev))

## [v0.12.0-rc5](https://github.com/andersfylling/disgord/tree/v0.12.0-rc5) (2019-09-22)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.12.0-rc4...v0.12.0-rc5)

**Implemented enhancements:**

- Enhanced ready event for guild loading. [\#198](https://github.com/andersfylling/disgord/issues/198)
- refactor websocket logic [\#162](https://github.com/andersfylling/disgord/issues/162)
- Clarify if a message event is from a guild or a direct message [\#153](https://github.com/andersfylling/disgord/issues/153)
- Initiate reconnect instead of stopping when Client.Connect\(\) fails [\#141](https://github.com/andersfylling/disgord/issues/141)
- Allow using channels instead of just handlers in .On method [\#131](https://github.com/andersfylling/disgord/issues/131)
- Internal loop for GetMessages [\#130](https://github.com/andersfylling/disgord/issues/130)
- helper functions for v0.12 [\#126](https://github.com/andersfylling/disgord/issues/126)
- Feature/integration tests [\#205](https://github.com/andersfylling/disgord/pull/205) ([andersfylling](https://github.com/andersfylling))
- rename ShardConfig.TotalNrOfShards to ShardConfig.ShardCount [\#204](https://github.com/andersfylling/disgord/pull/204) ([andersfylling](https://github.com/andersfylling))
- auto release on milestone close [\#192](https://github.com/andersfylling/disgord/pull/192) ([andersfylling](https://github.com/andersfylling))
- Use millisecond precision header [\#165](https://github.com/andersfylling/disgord/pull/165) ([andersfylling](https://github.com/andersfylling))
- specify events to ignore rather than handle [\#149](https://github.com/andersfylling/disgord/pull/149) ([andersfylling](https://github.com/andersfylling))

**Fixed bugs:**

- Update dockerfile [\#191](https://github.com/andersfylling/disgord/issues/191)
- Deadlock during reconnect phase [\#132](https://github.com/andersfylling/disgord/issues/132)
- Deadline for heartbeat ack is too low [\#168](https://github.com/andersfylling/disgord/issues/168)
- refactor websocket logic [\#162](https://github.com/andersfylling/disgord/issues/162)
- fixes issue with identify for distributed bots doing sharding [\#199](https://github.com/andersfylling/disgord/pull/199) ([andersfylling](https://github.com/andersfylling))
- fixes client.Ready for distributed bots [\#197](https://github.com/andersfylling/disgord/pull/197) ([andersfylling](https://github.com/andersfylling))

**Closed issues:**

- Add option to disable listening for presence\_updates and typing events [\#160](https://github.com/andersfylling/disgord/issues/160)

**Merged pull requests:**

- document build tags + introduce legacy build tag for REST method… [\#202](https://github.com/andersfylling/disgord/pull/202) ([andersfylling](https://github.com/andersfylling))
- remove short events pkg [\#201](https://github.com/andersfylling/disgord/pull/201) ([andersfylling](https://github.com/andersfylling))
- add GuildsReady method \(fixes \#198\) [\#200](https://github.com/andersfylling/disgord/pull/200) ([andersfylling](https://github.com/andersfylling))
- Allow registering event channels as if they are handlers [\#147](https://github.com/andersfylling/disgord/pull/147) ([andersfylling](https://github.com/andersfylling))
- Refactor sharding [\#146](https://github.com/andersfylling/disgord/pull/146) ([andersfylling](https://github.com/andersfylling))

## [v0.12.0-rc4](https://github.com/andersfylling/disgord/tree/v0.12.0-rc4) (2019-09-15)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.12.0-rc3...v0.12.0-rc4)

## [v0.12.0-rc3](https://github.com/andersfylling/disgord/tree/v0.12.0-rc3) (2019-09-15)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.12.0-rc2...v0.12.0-rc3)

## [v0.12.0-rc2](https://github.com/andersfylling/disgord/tree/v0.12.0-rc2) (2019-09-15)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.12.0-rc1...v0.12.0-rc2)

**Merged pull requests:**

- User.AvatarURL\(size int\) string [\#189](https://github.com/andersfylling/disgord/pull/189) ([victionn](https://github.com/victionn))

## [v0.12.0-rc1](https://github.com/andersfylling/disgord/tree/v0.12.0-rc1) (2019-09-15)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.11.3...v0.12.0-rc1)

**Implemented enhancements:**

- Implement voice-kicking [\#133](https://github.com/andersfylling/disgord/issues/133)
- use LFU as crs \(breaking\) [\#181](https://github.com/andersfylling/disgord/pull/181) ([andersfylling](https://github.com/andersfylling))

**Fixed bugs:**

- How to handle gateway error 4011? [\#184](https://github.com/andersfylling/disgord/issues/184)

**Closed issues:**

- discord invite links invalid? [\#187](https://github.com/andersfylling/disgord/issues/187)
- Allow setting number of total shards [\#177](https://github.com/andersfylling/disgord/issues/177)
- add internal alias pkg [\#163](https://github.com/andersfylling/disgord/issues/163)
- support ms precision in ratelimit headers [\#158](https://github.com/andersfylling/disgord/issues/158)
- Unknown import path "zeromod" [\#156](https://github.com/andersfylling/disgord/issues/156)

**Merged pull requests:**

- \[br\] auto-scaling on error + re-distributing message queues [\#188](https://github.com/andersfylling/disgord/pull/188) ([andersfylling](https://github.com/andersfylling))
- upgrade time [\#186](https://github.com/andersfylling/disgord/pull/186) ([andersfylling](https://github.com/andersfylling))
- add nhooyr websocket packet [\#176](https://github.com/andersfylling/disgord/pull/176) ([andersfylling](https://github.com/andersfylling))
- set minimum go version to 1.12 [\#175](https://github.com/andersfylling/disgord/pull/175) ([andersfylling](https://github.com/andersfylling))
- Feature/snowflake v4 [\#172](https://github.com/andersfylling/disgord/pull/172) ([andersfylling](https://github.com/andersfylling))
- allow kicking member from voice channel \(fixes \#133\) [\#171](https://github.com/andersfylling/disgord/pull/171) ([andersfylling](https://github.com/andersfylling))
- use callback and not channels for shard syncing [\#167](https://github.com/andersfylling/disgord/pull/167) ([andersfylling](https://github.com/andersfylling))
- add DM check to Message [\#166](https://github.com/andersfylling/disgord/pull/166) ([andersfylling](https://github.com/andersfylling))
- add IsByBot middleware + refactored logic into utils [\#155](https://github.com/andersfylling/disgord/pull/155) ([jfoster](https://github.com/jfoster))
- clarify self-bot support \(resolves \#150\) [\#151](https://github.com/andersfylling/disgord/pull/151) ([nikkelma](https://github.com/nikkelma))

## [v0.11.3](https://github.com/andersfylling/disgord/tree/v0.11.3) (2019-07-10)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.11.2...v0.11.3)

**Closed issues:**

- Panic when guild.Members [\#152](https://github.com/andersfylling/disgord/issues/152)
- \[Docs\] Clarify self-bot support [\#150](https://github.com/andersfylling/disgord/issues/150)

**Merged pull requests:**

- correct the json key name for inline field [\#154](https://github.com/andersfylling/disgord/pull/154) ([victionn](https://github.com/victionn))

## [v0.11.2](https://github.com/andersfylling/disgord/tree/v0.11.2) (2019-06-22)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.10.5...v0.11.2)

**Merged pull requests:**

- GetMessages with internal looping [\#145](https://github.com/andersfylling/disgord/pull/145) ([andersfylling](https://github.com/andersfylling))

## [v0.10.5](https://github.com/andersfylling/disgord/tree/v0.10.5) (2019-06-08)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.11.1...v0.10.5)

## [v0.11.1](https://github.com/andersfylling/disgord/tree/v0.11.1) (2019-06-08)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.10.4...v0.11.1)

**Closed issues:**

- Handlers do not use event from middlewares [\#143](https://github.com/andersfylling/disgord/issues/143)

**Merged pull requests:**

- use event mutated by middleware in event handlers [\#144](https://github.com/andersfylling/disgord/pull/144) ([andersfylling](https://github.com/andersfylling))

## [v0.10.4](https://github.com/andersfylling/disgord/tree/v0.10.4) (2019-05-30)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.11.0...v0.10.4)

## [v0.11.0](https://github.com/andersfylling/disgord/tree/v0.11.0) (2019-05-29)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.11.0-rc3...v0.11.0)

**Closed issues:**

- Query for channel / role ID's? [\#140](https://github.com/andersfylling/disgord/issues/140)

**Merged pull requests:**

- 📝 fix small typo in README.md [\#137](https://github.com/andersfylling/disgord/pull/137) ([BigHeadGeorge](https://github.com/BigHeadGeorge))

## [v0.11.0-rc3](https://github.com/andersfylling/disgord/tree/v0.11.0-rc3) (2019-05-16)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.11.0-rc1...v0.11.0-rc3)

## [v0.11.0-rc1](https://github.com/andersfylling/disgord/tree/v0.11.0-rc1) (2019-05-09)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.11.0-rc2...v0.11.0-rc1)

## [v0.11.0-rc2](https://github.com/andersfylling/disgord/tree/v0.11.0-rc2) (2019-05-09)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.10.3...v0.11.0-rc2)

**Fixed bugs:**

- Unexpected Character when unmarshalling GUILD\_CREATE events after connect [\#135](https://github.com/andersfylling/disgord/issues/135)

**Closed issues:**

- Rewrite socket client\#connect to wait for a given event [\#134](https://github.com/andersfylling/disgord/issues/134)

**Merged pull requests:**

- Fixes references to Evt\* objects [\#129](https://github.com/andersfylling/disgord/pull/129) ([jravesloot](https://github.com/jravesloot))

## [v0.10.3](https://github.com/andersfylling/disgord/tree/v0.10.3) (2019-04-02)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.10.2...v0.10.3)

## [v0.10.2](https://github.com/andersfylling/disgord/tree/v0.10.2) (2019-04-01)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.10.1...v0.10.2)

## [v0.10.1](https://github.com/andersfylling/disgord/tree/v0.10.1) (2019-03-24)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.10.0...v0.10.1)

## [v0.10.0](https://github.com/andersfylling/disgord/tree/v0.10.0) (2019-03-20)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.10.0-rc5...v0.10.0)

## [v0.10.0-rc5](https://github.com/andersfylling/disgord/tree/v0.10.0-rc5) (2019-03-20)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.10.0-rc4...v0.10.0-rc5)

## [v0.10.0-rc4](https://github.com/andersfylling/disgord/tree/v0.10.0-rc4) (2019-03-20)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.10.0-rc3...v0.10.0-rc4)

## [v0.10.0-rc3](https://github.com/andersfylling/disgord/tree/v0.10.0-rc3) (2019-03-18)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.10.0-rc2...v0.10.0-rc3)

**Closed issues:**

- Missing paging support in cache for getting guild members [\#125](https://github.com/andersfylling/disgord/issues/125)

## [v0.10.0-rc2](https://github.com/andersfylling/disgord/tree/v0.10.0-rc2) (2019-03-18)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.10.0-rc1...v0.10.0-rc2)

**Closed issues:**

- Channel caching and copying must be rewritten [\#124](https://github.com/andersfylling/disgord/issues/124)

## [v0.10.0-rc1](https://github.com/andersfylling/disgord/tree/v0.10.0-rc1) (2019-03-18)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.9.8...v0.10.0-rc1)

**Closed issues:**

- Custom marshallers for any discord struct that uses time.Time [\#123](https://github.com/andersfylling/disgord/issues/123)

## [v0.9.8](https://github.com/andersfylling/disgord/tree/v0.9.8) (2019-03-13)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.9.7...v0.9.8)

## [v0.9.7](https://github.com/andersfylling/disgord/tree/v0.9.7) (2019-03-13)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.9.6...v0.9.7)

**Implemented enhancements:**

- git hooks for pre-commit [\#111](https://github.com/andersfylling/disgord/issues/111)

**Fixed bugs:**

- Error when unmarshalling audit log entries [\#121](https://github.com/andersfylling/disgord/issues/121)

**Closed issues:**

- Add a accessible object pool for different data structures [\#89](https://github.com/andersfylling/disgord/issues/89)

**Merged pull requests:**

- Release v0.10 [\#122](https://github.com/andersfylling/disgord/pull/122) ([andersfylling](https://github.com/andersfylling))

## [v0.9.6](https://github.com/andersfylling/disgord/tree/v0.9.6) (2019-02-21)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.9.5...v0.9.6)

## [v0.9.5](https://github.com/andersfylling/disgord/tree/v0.9.5) (2019-02-18)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.9.4...v0.9.5)

**Merged pull requests:**

- feat: update documentation to use New instead of NewClient [\#118](https://github.com/andersfylling/disgord/pull/118) ([CallumDenby](https://github.com/CallumDenby))

## [v0.9.4](https://github.com/andersfylling/disgord/tree/v0.9.4) (2019-02-16)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.9.3...v0.9.4)

## [v0.9.3](https://github.com/andersfylling/disgord/tree/v0.9.3) (2019-02-15)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.9.2...v0.9.3)

## [v0.9.2](https://github.com/andersfylling/disgord/tree/v0.9.2) (2019-02-11)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.9.1...v0.9.2)

**Fixed bugs:**

- New project creation fails \(scripted and manual\) [\#115](https://github.com/andersfylling/disgord/issues/115)

**Closed issues:**

- Adds a handler controller \(builds on middleware\) [\#114](https://github.com/andersfylling/disgord/issues/114)
- Add middleware support for socket events [\#91](https://github.com/andersfylling/disgord/issues/91)

**Merged pull requests:**

- Introduces middleware and handler controllers for events [\#116](https://github.com/andersfylling/disgord/pull/116) ([andersfylling](https://github.com/andersfylling))

## [v0.9.1](https://github.com/andersfylling/disgord/tree/v0.9.1) (2019-02-10)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.9.0...v0.9.1)

## [v0.9.0](https://github.com/andersfylling/disgord/tree/v0.9.0) (2019-02-10)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.9.0-beta...v0.9.0)

## [v0.9.0-beta](https://github.com/andersfylling/disgord/tree/v0.9.0-beta) (2019-02-09)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.9.0-alpha...v0.9.0-beta)

## [v0.9.0-alpha](https://github.com/andersfylling/disgord/tree/v0.9.0-alpha) (2019-02-07)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.8.8...v0.9.0-alpha)

**Implemented enhancements:**

- Rewrite websocket to simplify maintenance [\#102](https://github.com/andersfylling/disgord/issues/102)
- Missing Voice support [\#92](https://github.com/andersfylling/disgord/issues/92)

**Closed issues:**

- add spoiler tag for messages [\#109](https://github.com/andersfylling/disgord/issues/109)

**Merged pull requests:**

- adds script for making a basic bot \(fixes \#83\) [\#112](https://github.com/andersfylling/disgord/pull/112) ([andersfylling](https://github.com/andersfylling))
- refactors websocketing + adds voice support + better shard handling [\#110](https://github.com/andersfylling/disgord/pull/110) ([andersfylling](https://github.com/andersfylling))

## [v0.8.8](https://github.com/andersfylling/disgord/tree/v0.8.8) (2019-01-30)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.8.7...v0.8.8)

## [v0.8.7](https://github.com/andersfylling/disgord/tree/v0.8.7) (2019-01-29)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.8.6...v0.8.7)

**Closed issues:**

- Allow setting hardcoded rate limits [\#107](https://github.com/andersfylling/disgord/issues/107)

## [v0.8.6](https://github.com/andersfylling/disgord/tree/v0.8.6) (2019-01-28)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.8.5...v0.8.6)

## [v0.8.5](https://github.com/andersfylling/disgord/tree/v0.8.5) (2019-01-28)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.8.4...v0.8.5)

**Implemented enhancements:**

- Can't see number of connected guilds [\#79](https://github.com/andersfylling/disgord/issues/79)
- Introduce DisgordErr [\#74](https://github.com/andersfylling/disgord/issues/74)
- Drop requirement that marshalling should equal incoming json data [\#69](https://github.com/andersfylling/disgord/issues/69)
- Add custom error types [\#46](https://github.com/andersfylling/disgord/issues/46)
- output to log when detecting incorrect rate limiters [\#38](https://github.com/andersfylling/disgord/issues/38)
- Consider using fasthttp [\#37](https://github.com/andersfylling/disgord/issues/37)
- allow developers to inject a logger [\#19](https://github.com/andersfylling/disgord/issues/19)

**Fixed bugs:**

- Handle ratelimits for multiple shards to avoid reconnect loops [\#82](https://github.com/andersfylling/disgord/issues/82)
- Drop requirement that marshalling should equal incoming json data [\#69](https://github.com/andersfylling/disgord/issues/69)

**Closed issues:**

- Rework the git log to reduce repo size [\#101](https://github.com/andersfylling/disgord/issues/101)
- Handlers: support simpler functions [\#99](https://github.com/andersfylling/disgord/issues/99)
- Move live tests into subpkg [\#96](https://github.com/andersfylling/disgord/issues/96)

## [v0.8.4](https://github.com/andersfylling/disgord/tree/v0.8.4) (2018-12-12)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.8.3...v0.8.4)

**Closed issues:**

- Panic after connection lost and attempted reconnect [\#98](https://github.com/andersfylling/disgord/issues/98)

## [v0.8.3](https://github.com/andersfylling/disgord/tree/v0.8.3) (2018-12-11)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.8.2...v0.8.3)

**Implemented enhancements:**

- Simplify status updates [\#76](https://github.com/andersfylling/disgord/issues/76)

**Fixed bugs:**

- Crashing after loosing connection and attempting reconnect [\#97](https://github.com/andersfylling/disgord/issues/97)
- Remove pointers from CreateGuildChannelParams [\#93](https://github.com/andersfylling/disgord/issues/93)

**Closed issues:**

- Rate limiting for socket commands [\#84](https://github.com/andersfylling/disgord/issues/84)

**Merged pull requests:**

- feat: Client.SetStatus, Client.SetStatusString [\#85](https://github.com/andersfylling/disgord/pull/85) ([Soumil07](https://github.com/Soumil07))

## [v0.8.2](https://github.com/andersfylling/disgord/tree/v0.8.2) (2018-11-13)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.8.1...v0.8.2)

**Fixed bugs:**

- Panic when losing internet connection\(?\) [\#48](https://github.com/andersfylling/disgord/issues/48)

## [v0.8.1](https://github.com/andersfylling/disgord/tree/v0.8.1) (2018-10-31)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.8.0...v0.8.1)

**Fixed bugs:**

- Status does not update [\#80](https://github.com/andersfylling/disgord/issues/80)

## [v0.8.0](https://github.com/andersfylling/disgord/tree/v0.8.0) (2018-10-30)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.7.0...v0.8.0)

**Implemented enhancements:**

- Evaluate use of pointers and internal interfaces for improved performance [\#66](https://github.com/andersfylling/disgord/issues/66)
- Add configurable lifetime option for items in cache [\#62](https://github.com/andersfylling/disgord/issues/62)
- Refactor event\_dispatcher.go [\#54](https://github.com/andersfylling/disgord/issues/54)
- update myself on socket events [\#53](https://github.com/andersfylling/disgord/issues/53)
- Project install guide [\#75](https://github.com/andersfylling/disgord/issues/75)
- fix godoc [\#57](https://github.com/andersfylling/disgord/issues/57)
- Missing support for sharding [\#45](https://github.com/andersfylling/disgord/issues/45)
- Refactor/websocket [\#72](https://github.com/andersfylling/disgord/pull/72) ([andersfylling](https://github.com/andersfylling))

**Fixed bugs:**

- Add configurable lifetime option for items in cache [\#62](https://github.com/andersfylling/disgord/issues/62)
- Add config param "ActivateEventChannels" [\#78](https://github.com/andersfylling/disgord/issues/78)
- Modify requests should allow for resetting to default values [\#68](https://github.com/andersfylling/disgord/issues/68)

**Closed issues:**

- File uploads [\#55](https://github.com/andersfylling/disgord/issues/55)
- Use json.RawMessage instead [\#47](https://github.com/andersfylling/disgord/issues/47)
- Cannot stop disgord during reconnect [\#11](https://github.com/andersfylling/disgord/issues/11)

**Merged pull requests:**

- Removed println that was flooding the console [\#71](https://github.com/andersfylling/disgord/pull/71) ([pizza61](https://github.com/pizza61))
- Fix/circle ci coverage [\#70](https://github.com/andersfylling/disgord/pull/70) ([andersfylling](https://github.com/andersfylling))
- 🎨📝🔥👕 Add initial implementation of go generate files [\#67](https://github.com/andersfylling/disgord/pull/67) ([ikkerens](https://github.com/ikkerens))
- 🎨📝 Added build constraints to switch between json-iterator and std [\#61](https://github.com/andersfylling/disgord/pull/61) ([ikkerens](https://github.com/ikkerens))
- 🔥🎨📝 Added file upload support \(fixes \#55\) [\#60](https://github.com/andersfylling/disgord/pull/60) ([ikkerens](https://github.com/ikkerens))
- 📝 Updated docs [\#59](https://github.com/andersfylling/disgord/pull/59) ([ikkerens](https://github.com/ikkerens))
- 📝 Adds support for custom httpclient implementations in websocket layer [\#58](https://github.com/andersfylling/disgord/pull/58) ([ikkerens](https://github.com/ikkerens))

## [v0.7.0](https://github.com/andersfylling/disgord/tree/v0.7.0) (2018-09-16)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.6.2...v0.7.0)

**Implemented enhancements:**

- Enforce sequential event handlers [\#50](https://github.com/andersfylling/disgord/issues/50)

**Fixed bugs:**

- Use date from http response in rate limiting [\#51](https://github.com/andersfylling/disgord/issues/51)
- reset is obviusly going to be ahead of time [\#52](https://github.com/andersfylling/disgord/issues/52)

**Closed issues:**

- Missing event handler for `PRESENCES\_REPLACE` [\#49](https://github.com/andersfylling/disgord/issues/49)
- Disconnections [\#10](https://github.com/andersfylling/disgord/issues/10)
- cleanup socket branch [\#9](https://github.com/andersfylling/disgord/issues/9)

## [v0.6.2](https://github.com/andersfylling/disgord/tree/v0.6.2) (2018-09-09)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.6.1...v0.6.2)

## [v0.6.1](https://github.com/andersfylling/disgord/tree/v0.6.1) (2018-09-09)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.6.0...v0.6.1)

## [v0.6.0](https://github.com/andersfylling/disgord/tree/v0.6.0) (2018-09-09)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.5.1...v0.6.0)

## [v0.5.1](https://github.com/andersfylling/disgord/tree/v0.5.1) (2018-09-09)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.5.0...v0.5.1)

## [v0.5.0](https://github.com/andersfylling/disgord/tree/v0.5.0) (2018-09-08)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.4.2...v0.5.0)

**Fixed bugs:**

- Emoji endpoints have incorrect rate limit bucket key [\#43](https://github.com/andersfylling/disgord/issues/43)
- Does not stop goroutine on shutdown [\#42](https://github.com/andersfylling/disgord/issues/42)

## [v0.4.2](https://github.com/andersfylling/disgord/tree/v0.4.2) (2018-09-05)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.4.1...v0.4.2)

## [v0.4.1](https://github.com/andersfylling/disgord/tree/v0.4.1) (2018-09-05)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.4.0...v0.4.1)

## [v0.4.0](https://github.com/andersfylling/disgord/tree/v0.4.0) (2018-09-05)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.3.0...v0.4.0)

**Implemented enhancements:**

- Code duplication in the rest package [\#30](https://github.com/andersfylling/disgord/issues/30)

**Fixed bugs:**

- Socketing/Gateway doesn't recieve updates after READY event [\#41](https://github.com/andersfylling/disgord/issues/41)
- caching\(state\) logic contains a racecondition [\#27](https://github.com/andersfylling/disgord/issues/27)

## [v0.3.0](https://github.com/andersfylling/disgord/tree/v0.3.0) (2018-09-02)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.2.0...v0.3.0)

## [v0.2.0](https://github.com/andersfylling/disgord/tree/v0.2.0) (2018-09-02)
[Full Changelog](https://github.com/andersfylling/disgord/compare/v0.1.0...v0.2.0)

**Implemented enhancements:**

- Missing CONTRIBUTING.md [\#24](https://github.com/andersfylling/disgord/issues/24)
- REST and authentication [\#23](https://github.com/andersfylling/disgord/issues/23)
- missing abstract func for retrieving data [\#20](https://github.com/andersfylling/disgord/issues/20)
- faster json lib [\#1](https://github.com/andersfylling/disgord/issues/1)
- Feature/contributing [\#29](https://github.com/andersfylling/disgord/pull/29) ([andersfylling](https://github.com/andersfylling))
- Test/ratelimit [\#28](https://github.com/andersfylling/disgord/pull/28) ([andersfylling](https://github.com/andersfylling))

**Fixed bugs:**

- Discord uses different epoch on snowflake [\#36](https://github.com/andersfylling/disgord/issues/36)
- REST and authentication [\#23](https://github.com/andersfylling/disgord/issues/23)

**Closed issues:**

- Incorrect rate limiting? [\#21](https://github.com/andersfylling/disgord/issues/21)
- Store event handler pointers in array, and link index to event const for faster lookup [\#2](https://github.com/andersfylling/disgord/issues/2)

**Merged pull requests:**

- Merge pull request \#33 from andersfylling/develop [\#34](https://github.com/andersfylling/disgord/pull/34) ([andersfylling](https://github.com/andersfylling))
- Create LICENSE [\#22](https://github.com/andersfylling/disgord/pull/22) ([andersfylling](https://github.com/andersfylling))
- Feature/http requests WIP [\#16](https://github.com/andersfylling/disgord/pull/16) ([andersfylling](https://github.com/andersfylling))

## [v0.1.0](https://github.com/andersfylling/disgord/tree/v0.1.0) (2018-02-19)
**Closed issues:**

- Sequence number is lost on reconnect [\#12](https://github.com/andersfylling/disgord/issues/12)
- Reconnecting [\#8](https://github.com/andersfylling/disgord/issues/8)
- Correct websocket disconnect [\#7](https://github.com/andersfylling/disgord/issues/7)
- event dispatcher from the socket layer [\#6](https://github.com/andersfylling/disgord/issues/6)
- Handling discord events related to sockets [\#5](https://github.com/andersfylling/disgord/issues/5)
- Heartbeat [\#4](https://github.com/andersfylling/disgord/issues/4)

**Merged pull requests:**

- Feature/websocket [\#13](https://github.com/andersfylling/disgord/pull/13) ([andersfylling](https://github.com/andersfylling))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*