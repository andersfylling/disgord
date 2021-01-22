package disgord

type ApplicationQueryBuilder interface {
	Guild(Snowflake) ApplicationGuildQueryBuilder
	Global() ApplicationGlobalQueryBuilder
}

func (c *Client) Application(id Snowflake) applicationQueryBuilder {
	return applicationQueryBuilder{appID: id}
}

type applicationQueryBuilder struct {
	appID Snowflake
}

type ApplicationGlobalQueryBuilder interface {
	GetCommands()
	CreateCommand()
	Command(Snowflake) ApplicationCommandQueryBuilder
}

type ApplicationGuildQueryBuilder interface {
	GetCommands()
	CreateCommand()
	Command(Snowflake) ApplicationCommandQueryBuilder
}

type ApplicationCommandQueryBuilder interface {
	Create()
	Get()
	Delete()
	UpdateBuilder()
}
