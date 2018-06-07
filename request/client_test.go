package request

import "testing"

func missingImplError(t *testing.T, interfaceName string) {
	t.Error("client{} does not implement " + interfaceName + " interface")
}

func TestClientImplementInterfaces(t *testing.T) {
	client := &Client{}
	if _, implemented := interface{}(client).(DiscordRequester); !implemented {
		missingImplError(t, "DiscordRequester")
	}
	if _, implemented := interface{}(client).(DiscordGetter); !implemented {
		missingImplError(t, "DiscordGetter")
	}
	if _, implemented := interface{}(client).(DiscordPoster); !implemented {
		missingImplError(t, "DiscordPoster")
	}
	if _, implemented := interface{}(client).(DiscordPutter); !implemented {
		missingImplError(t, "DiscordPutter")
	}
	if _, implemented := interface{}(client).(DiscordPatcher); !implemented {
		missingImplError(t, "DiscordPatcher")
	}
	if _, implemented := interface{}(client).(DiscordDeleter); !implemented {
		missingImplError(t, "DiscordDeleter")
	}
}

func TestRateLimiter(t *testing.T) {

}
