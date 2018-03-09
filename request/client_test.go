package request

import "testing"

func missingImplError(t *testing.T, interfaceName string) {
	t.Error("client{} does not implement " + interfaceName + " interface")
}

func TestClientImplementInterfaces(t *testing.T) {
	if _, implemented := interface{}(&Client{}).(DiscordRequester); !implemented {
		missingImplError(t, "DiscordRequester")
	}
	if _, implemented := interface{}(&Client{}).(DiscordGetter); !implemented {
		missingImplError(t, "DiscordGetter")
	}
	if _, implemented := interface{}(&Client{}).(DiscordPoster); !implemented {
		missingImplError(t, "DiscordPoster")
	}
	if _, implemented := interface{}(&Client{}).(DiscordPutter); !implemented {
		missingImplError(t, "DiscordPutter")
	}
	if _, implemented := interface{}(&Client{}).(DiscordPatcher); !implemented {
		missingImplError(t, "DiscordPatcher")
	}
	if _, implemented := interface{}(&Client{}).(DiscordDeleter); !implemented {
		missingImplError(t, "DiscordDeleter")
	}
}

func TestRateLimiter(t *testing.T) {

}
