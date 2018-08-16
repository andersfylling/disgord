package disgord

import (
	"testing"
	"os"
	"fmt"
	"github.com/andersfylling/disgord/rest"
	"github.com/andersfylling/snowflake"
)

func TestRESTClient_users(t *testing.T) {
	token := os.Getenv("DISGORD_DISCORD_TEST_MAIN")
	if token == "" {
		fmt.Println("not running unit-test: testClient_REST_endpoints; missing token")
		return
	}


	conf := &Config{
		Token: token,
		Debug: true,
	}
	client, err := NewClient(conf)
	if err != nil {
		t.Error(err)
	}
	//client.Connect()

	req := client.Req()

	usr, err := rest.GetUser(req, snowflake.NewID(228846961774559232))
	if err != nil {
		t.Error(err.Error())
	}
	if usr.ID.String() != "228846961774559232" {
		t.Error("did not retrieve correct user")
	}

}