package cmd

import "testing"

/**
测试匹配service字符
*/

var output = `
● phoenix.service - Phoenix Server
   Loaded: loaded (/usr/lib/systemd/system/phoenix.service; enabled; vendor preset: disabled)
   Active: active (running) since 五 2021-03-26 14:58:36 CST; 32s ago
     Docs: https://www.woqutech.com/
 Main PID: 2850 (phoenix)
   CGroup: /system.slice/phoenix.service
           └─2850 /home/sendoh/qdm_control/packages/holmes/phoenix/bin/phoenix --alert-config=alertmanager.yml --alert-url=http://127.0.0.1:10012/alertmanag...

3月 26 14:58:36 192-168-1-99 phoenix[2850]: [1 rows affected or returned ]
`

func TestService(t *testing.T) {
	service := Service{
		statusOutput:  output,
		autoStartFile: "/etc/systemd/system/multi-user.target.wants/phoenix.service",
	}
	t.Log(service.IsValid())
	t.Log(service.Name())
	t.Log(service.Status())
	t.Log(service.TimeDuration())
	t.Log(service.PID())
	t.Log(service.AutoStart())
}
