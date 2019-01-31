// +build !disgord_parallelism

package constant

// LockedMethods signifies if the methods of discord objects should handle locking.
// I don't enjoy introducing this, but at the same time, I don't want to completely remove a
// easy way to activate locking for people that require parallel handling of objects.
//
// Yes. They could set mutexes themselves, but if they want performance, they want to only lock when they actually
// deal read and write to the objects, not while doing casting, error checks, and other time consuming operations
// which often takes place in internal methods before actually interaction with the "unsafe" content.
//
// TODO: reconsider - This option affects the behaviour in a breaking way.
const LockedMethods = false
