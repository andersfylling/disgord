package httd

import "testing"

func missingImplError(t *testing.T, interfaceName string) {
	t.Error("client{} does not implement " + interfaceName + " interface")
}

func TestClientImplementInterfaces(t *testing.T) {
	client := &Client{}
	if _, implemented := interface{}(client).(Requester); !implemented {
		missingImplError(t, "Requester")
	}
	if _, implemented := interface{}(client).(Getter); !implemented {
		missingImplError(t, "Getter")
	}
	if _, implemented := interface{}(client).(Poster); !implemented {
		missingImplError(t, "Poster")
	}
	if _, implemented := interface{}(client).(Puter); !implemented {
		missingImplError(t, "Puter")
	}
	if _, implemented := interface{}(client).(Patcher); !implemented {
		missingImplError(t, "Patcher")
	}
	if _, implemented := interface{}(client).(Deleter); !implemented {
		missingImplError(t, "Deleter")
	}
}

func TestRateLimiter(t *testing.T) {

}
