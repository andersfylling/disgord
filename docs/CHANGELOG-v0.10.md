This document lists the changes in names, types and behaviors from v0.9 to v0.10.


### Name changes:

 - All the permissions consts ([a-zA-Z]*Permission), have moved the "Permission" suffix to the prefix: eg. SendMessagePermission => PermissionSendMessage.
 This is to group related const together when using an IDE or another auto-completion tool.
 
 - Custom Permissions such as All, AllChannel, etc. Now have the "All" suffix as the first word after Permission: eg. PermissionAll => PermissionAll, PermissionAllChannel => PermissionChannelAll.

#### REST
Discord has a bunch of inconsistencies in their naming scheme for the REST methods. DisGord have always wanted to be as close to the docs as possible to ease bot development while reading the Discord docs. However, it seems more convenient to change the name such that they stay consistent and reflect their purpose that is easier to understand for new comers.

REST methods starting with Modify, Update, Edit, etc. Now all start with "Update". Those endpoints that exists solely to change one attribute, are prefix with "Set" instead of "Update".

Most endpoints starting with "Remove" has been renamed to start with "Delete", as they actually delete the resource instead of removing/revoking it from said resource. eg. RemoveChannel is a deletion, while RemoveMemberRole, removes the role from the member but does not delete it. This distinction matters.

This means that any related config struct has been renamed to match the method name.

But, to make it easier for DisGord users, all the endpoint defined in the Discord Docs are still on the client. They wrap the actual method/methods, and are marked Deprecated. So please update your source code.
 
##### Renamed
 - func CreateChannelMessage => CreateMessage
 - func UpdateChannel was removed, as func ModifyChannel was renamed to UpdateChannel
 - func CreateGuildBan => BanMember
 - func RemoveGuildBan => UnbanMember
 - func BeginGuildPrune => PruneMembers
 - func GetGuildPruneCount => EstimatePruneMembersCount
 - func ModifyCurrentUserNick => SetCurrentUserNick
 - func RemoveGuildMember => KickMember
 - func ModifyGuildMember => UpdateGuildMember
 - func ModifyGuildChannelPositions => UpdateGuildChannelPositions
 - func ModifyGuild => UpdateGuild
 - func ModifyGuildEmbed => UpdateGuildEmbed
 - func ModifyGuildIntegration => UpdateGuildIntegration
 - func ModifyGuildRole => UpdateGuildRole
 - func ModifyGuildRolePositions => UpdateGuildRolePositions
 - func ModifyChannel => UpdateChannel
 - type ModifyChannelParams => UpdateChannelParams
 - func GroupDMRemoveRecipient => KickParticipant
 - func GroupDMAddRecipient => AddDMParticipant
 - func CloseChannel => DeleteChannel
 - func BulkDeleteMessages => DeleteMessages
 - func EditChannelPermissions => UpdateChannelPermissions
 - func AddPinnedChannelMessage => PinMessage + PinMessageID
 - func DeletePinnedChannelMessage => UnpinMessage + UnpinMessageID

#### Removed / unexported
 - type BeginGuildPruneParams
 - type GuildPruneCount
 - type UpdateCurrentUserNickParams
 - .Once(event string, inputs ...interface{}) error
 - NewSession(conf *Config) (Session, error)
 - NewSessionMustCompile(conf *Config) (session Session)
 
 - every REST func. REST functionality now only accessible from the client instance.
 
 #### Changed
  - .On(event string, inputs ...interface{}) error => .On(event string, inputs ...interface{})