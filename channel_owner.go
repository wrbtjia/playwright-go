package playwright

type ChannelOwner struct {
	EventEmitter
	objectType  string
	guid        string
	channel     *Channel
	objects     map[string]*ChannelOwner
	connection  *Connection
	initializer map[string]interface{}
	parent      *ChannelOwner
}

func (c *ChannelOwner) Dispose() {
	// Clean up from parent and connection.
	if c.parent != nil {
		delete(c.parent.objects, c.guid)
	}
	delete(c.connection.objects, c.guid)

	// Dispose all children.
	for _, object := range c.objects {
		object.Dispose()
	}
	c.objects = make(map[string]*ChannelOwner)
}

func (c *ChannelOwner) createChannelOwner(self interface{}, parent *ChannelOwner, objectType string, guid string, initializer map[string]interface{}) {
	c.objectType = objectType
	c.guid = guid
	c.parent = parent
	c.objects = make(map[string]*ChannelOwner)
	c.connection = parent.connection
	c.channel = newChannel(c.connection, guid)
	c.channel.object = self
	c.initializer = initializer
	c.connection.objects[guid] = c
	c.parent.objects[guid] = c
	c.initEventEmitter()
}

func newRootChannelOwner(connection *Connection) *ChannelOwner {
	c := &ChannelOwner{
		objectType: "",
		guid:       "",
		connection: connection,
		objects:    make(map[string]*ChannelOwner),
		channel:    newChannel(connection, ""),
	}
	c.channel.object = c
	c.connection.objects[""] = c
	c.initEventEmitter()
	return c
}
