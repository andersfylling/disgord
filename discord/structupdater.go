package discord

type StructUpdater interface {
	// UpdateOrigin compares the current object with the `new` instance, if any changes are detected,
	//              a request is sent to the discord endpoint to update the given object.
	//              Note! that the object ID must be the same.
	UpdateOrigin(new interface{}) (err error)
}

func UpdateOldStruct(old interface{}, new interface{}) (err error) {

	// TODO
	// new can be of a different type, but must be a struct containing similar fields with same type to compare.
	return nil
}
