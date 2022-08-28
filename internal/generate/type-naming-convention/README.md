A little hack to detect and fail when reserved naming convention is abused.

# DirectionalNames: "CRUD" types
to keep it simple a lot has been renamed to only use UPDATE, DELETE and CREATE ("DUC" for the rest of this section). Especially for events and REST communication.
We want to reserve names that starts and ends with "DUC" for updating and receiving information about changes at the Discord state.

In this scope, a type is used for either:
 - incoming changes about the discord state. Where the naming pattern is **DUC***Object*
 - to request changes to the discord state. Where the naming pattern is *Object***DUC**

examples:
 - *Message***Create** (sent by discord as a consequence of change in the discord state)
 - **Create***Message* (sent by client to change discord state)

By preserving we mean that only type definition are whitelisted for usage. consts, vars, etc. are not allowed to use a similar naming convention.