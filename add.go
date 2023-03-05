package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
)

type LogInfo struct {
	Timestamp time.Time `json:"timestamp"`
	Log       string    `json:"log"`
}

// конвертация времени в стандартный формат
func convertTime(s string) time.Time {
	theTime, _ := time.Parse("2006-01-02", s)
	return theTime
}

// Функция для добавление логов, принимает ctx, client для отправление запроса и индекс
func add(ctx context.Context, client *elasticsearch.Client, index string) {
	// в слайсе allLogs сохраняем логи и их время, предварительно конвертируя в стандартный формат
	allLogs := []LogInfo{
		{convertTime("2023-02-13"), `localhost - - [13/Feb/2023:19:00:59 +0600] "POST / HTTP/1.1" 401 221 CUPS-Get-Devices successful-ok`},
		{convertTime("2023-02-13"), `[2023-01-26 16:21:07.184] [info]  compliance configDir=/Usjknxa`},
		{convertTime("2023-02-26"), `[2023-01-26 16:21:07.497] [info]  Events - processArgv ["/Applications/FortiClient.app/Contents/MacOS/FortiClient"]`},
		{convertTime("2023-01-16"), `2023-01-16 11:39:59 FortiClient Uninstall: Failed to connect to fctservctl`},
		{convertTime("2023-02-03"), `[2023-02-03 19:04:53.412] [info]  Channel handleMacMsgEvent "update_id"`},
		{convertTime("2023-01-16"), `[2023-01-16 11:40:47.101] [info]  Pipeserver - listening`},
		{convertTime("2023-02-27"), `Mon Feb 27 00:30:53.768 <airport[369]> _configureTCPKeepAlive: Unable to enable TCP keep-alive on en0 (Operation not supported)`},
		{convertTime("2023-02-03"), `[2023-02-03 17:07:11.265] [debug] Pipe server 'connection' [object Object]`},
		{convertTime("2022-11-12"), `/dev/rdisk1s1: fsck_apfs completed at Sat Nov 12 17:22:44 2022`},
		{convertTime("2022-07-20"), `2022-07-20 00:30:48-07 MacBook-Pro softwareupdated[192]: authorizeWithEmptyAuthorizationForRights: Requesting provided rights: 1`},
		{convertTime("2023-02-03"), `[2023-02-03 09:55:49.617] [debug] Pipe server 'connection' [object Object]`},
		{convertTime("2022-07-20"), `2022-07-20 00:30:54-07 MacBook-Pro Language Chooser[294]: LCA: most frequent country code: CN`},
		{convertTime("2023-01-17"), `[2023-01-17 06:19:33.719] [debug] Pipe server 'connection' [object Object]`},
		{convertTime("2022-09-29"), `2022-09-29 14:19:04+06 MacBook-Pro-Asem softwareupdated[222]: SoftwareUpdate: request for status for unknown product MSU_UPDATE_21G115_patch_12.6`},
		{convertTime("2023-02-13"), `localhost - asemalikova [13/Feb/2023:19:00:59 +0600] "POST / HTTP/1.1" 200 478 CUPS-Get-Devices -`},
		{convertTime("2022-11-12"), `/dev/rdisk2s1: fsck_apfs completed at Sat Nov 12 17:22:44 2022`},
		{convertTime("2023-02-27"), `2023-02-27 13:56:30.004427 (user/501/com.apple.mdworker.shared.05000000-0700-0000-0000-000000000000) <Notice>: launching: ipc (mach)`},
		{convertTime("2022-11-12"), `/dev/rdisk1s2: fsck_apfs started at Sat Nov 12 17:22:44 2022`},
		{convertTime("2022-09-29"), `2022-09-29 14:19:04+06 MacBook-Pro-Asem SoftwareUpdateNotificationManager[1229]: AssertionMgr: Could not cancel com.apple.softwareupdate.NotifyAgentAssertion-BadgingCountChanged assertion - no assertion found for pid 1229`},
		{convertTime("2022-09-29"), `2022-09-29 14:19:04+06 MacBook-Pro-Asem softwareupdated[222]: 3 updates found:`},
		{convertTime("2023-02-26"), `2023-02-26 23:47:26.329509 (gui/501/com.apple.amsaccountsd) <Notice>: service state: spawning`},
		{convertTime("2022-10-03"), `2022-10-03 18:25:04+06 MBP-Asem installd[4004]: installd: Starting`},
		{convertTime("2023-02-26"), `2023-02-26 23:14:46.865596 (system) <Notice>: removing child: pid/45393`},
		{convertTime("2022-09-29"), `2022-09-29 14:19:04+06 MacBook-Pro-Asem softwareupdated[222]: Scan (f=1, d=1) completed`},
		{convertTime("2023-02-26"), `2023-02-26 23:14:46.867087 (pid/60300/com.apple.quicklook.satellite) <Notice>: internal event: PETRIFIED, code = 0`},
		{convertTime("2022-07-20"), `2022-07-20 00:30:46-07 localhost Installer Progress[52]: Progress UI App Starting`},
		{convertTime("2023-02-03"), `[2023-02-03 10:11:22.835] [debug] Pipe server 'connection' [object Object]`},
		{convertTime("2022-11-12"), `/dev/rdisk2s1: fsck_apfs completed at Sat Nov 12 17:22:44 2022`},
		{convertTime("2023-02-03"), `[2023-02-03 15:17:28.360] [info]  client connected`},
		{convertTime("2022-09-29"), `2022-09-29 01:14:50 -0700 IMDPersistenceAgent[1007]: Created table (if needed) ok: kvtable`},
		{convertTime("2022-10-03"), `2022-10-03 18:38:07+06 MBP-Asem nbagent[3900]: NBStateController: Connection interrupted`},
		{convertTime("2022-09-30"), `2022-09-30 01:18:57 +0600 IMDPersistenceAgent[892]: Dropped trigger: after_delete_on_message`},
		{convertTime("2022-07-17"), `/dev/rdisk3s3: fsck_apfs completed at Tue Jan 17 15:32:17 2023`},
		{convertTime("2022-09-30"), `2022-09-30 01:18:57 +0600 IMDPersistenceAgent[892]: Dropped table (if needed) ok: sqlite_stat1`},
		{convertTime("2022-12-26"), `2022-12-26 11:32:39+06 MacBook-Pro-Asem softwareupdated[34380]: SUOSUAlarmObserver: Setting alarm event stream handler`},
		{convertTime("2022-12-26"), `2022-12-26 11:32:38+06 MacBook-Pro-Asem suhelperd[34310]: DISPATCH_MACH_DISCONNECTED for client port.`},
		{convertTime("2022-11-12"), `/dev/rdisk1s3: fsck_apfs completed at Sat Nov 12 17:22:44 2022`},
		{convertTime("2022-09-29"), `2022-09-29 14:19:03+06 MacBook-Pro-Asem softwareupdated[222]: Failed to get bridge device`},
		{convertTime("2022-09-30"), `2022-09-30 01:18:57 +0600 IMDPersistenceAgent[892]: Created table (if needed) ok: kvtable`},
		{convertTime("2023-02-25"), `2023-02-25 19:09:55+06 MacBook-Pro-Asem softwareupdated[65477]: SUOSUPowerEventObserver: System will power on`},
		{convertTime("2023-02-19"), `/dev/rdisk3s3: fsck_apfs started at Sun Feb 19 15:55:53 2023`},
		{convertTime("2022-09-30"), `2022-09-30 01:18:57 +0600 IMDPersistenceAgent[892]: Created table (if needed) ok: sync_deleted_chats`},
	}

	//для отправления данных в сервер конвертируем их в json формат
	result := make([]string, len(allLogs))
	for i, w := range allLogs {
		r, err := json.Marshal(&w)
		if err != nil {
			log.Fatalf("json.Marshal(allLogs) Error:", err)
		}
		result[i] = string(r)
	}

	//проходимся по json объектам
	for i, w := range result {
		// Настраиваем  запрос Index API
		req := esapi.IndexRequest{
			Index:      index,
			DocumentID: strconv.Itoa(i + 1),
			Body:       strings.NewReader(w),
			Refresh:    "true",
		}

		// С помощью функции Do() отправляем запрос req, ответ записываем rr
		// проверяем на ошибку, если все в порядке то проверяем ответ res,
		// если res.Status() возвращает ошибку то выводим его на экран
		res, err := req.Do(ctx, client)
		if err != nil {
			log.Fatalf("IndexRequest ERROR: %s", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("%s ERROR indexing document ID=%d", res.Status(), i+1)
		}
	}
}
