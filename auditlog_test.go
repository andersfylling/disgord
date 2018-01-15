package disgord

// import "testing"
//
// func TestConvertAuditLogParamsToStr(t *testing.T) {
// 	params := &AuditLogParams{}
// 	res := convertAuditLogParamsToStr(params)
//
// 	if res != "" {
// 		t.Errorf("Empty AuditLogParams struct did not create an empty string. Got %s, wants \"\"", res)
// 	}
// }

var auditLogSample string = "{\"webhooks\": [], \"users\": [{\"username\": \"Anders\", \"discriminator\": \"7237\", \"id\": \"228846961774559232\", \"avatar\": \"69a7a0e9cb963adfdd69a2224b4ac180\"}, {\"username\": \"alek\", \"discriminator\": \"1049\", \"id\": \"253218433276182528\", \"avatar\": \"38d04eba240fa3cad5816947025644ad\"}, {\"username\": \"CapinoMarket\", \"discrimi nator\": \"2022\", \"id\": \"211951321295749130\", \"avatar\": \"64c9831c2d5ee0e725ee8fb429ab7019\"}, {\"username\": \"JailBotTester\", \"discriminator\": \"6540\", \"bot\": true, \"id\": \"400741409134477323\", \"avatar\": null}], \"audit_log_entries\": [{\"target_id\": \"253218433276182528\", \"reason\": \"testing\", \"user_id\": \"228846961774559232\", \"id\": \"40075566523 8360074\", \"action_type\": 20}, {\"target_id\": \"400741409134477323\", \"changes\": [{\"new_value\": [{\"name\": \"mod\", \"id\": \"399704833185153025\"}], \"key\": \"$add\"}], \"user_id\": \"228846961774559232\", \"id\": \"400755530660053003\", \"action_type\": 25}, {\"target_id\": \"253218433276182528\", \"user_id\": \"228846961774559232\", \"id\": \"400699547866628107\", \" action_type\": 20}, {\"target_id\": \"253218433276182528\", \"user_id\": \"228846961774559232\", \"id\": \"400698947603005440\", \"action_type\": 20}, {\"target_id\": \"211951321295749130\", \"user_id\": \"228846961774559232\", \"id\": \"400698728509341718\", \"action_type\": 20}, {\"target_id\": \"253218433276182528\", \"user_id\": \"228846961774559232\", \"id\": \"400698 703087665153\", \"action_type\": 20}, {\"target_id\": \"253218433276182528\", \"user_id\": \"228846961774559232\", \"id\": \"400697915200503828\", \"action_type\": 20}, {\"target_id\": \"253218433276182528\", \"user_id\": \"228846961774559232\", \"id\": \"400697714590875669\", \"action_type\": 20}, {\"target_id\": \"253218433276182528\", \"user_id\": \"228846961774559232 \", \"id\": \"400697309903585282\", \"action_type\": 20}, {\"target_id\": \"211951321295749130\", \"user_id\": \"228846961774559232\", \"id\": \"400694998104014852\", \"action_type\": 23}]}"
