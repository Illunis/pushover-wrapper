package main

type messages struct {
	ID      int    `json:"id"`
	IDStr   string `json:"id_str"`
	Message string `json:"message"`
	App     string `json:"app"`
	/*
		"message": "This is a test alert",
		"app": "LibreNMS Work",
		"aid": 435443636951759403,
		"aid_str": "435443636951759403",
		"icon": "mnsw2ykc6qa5sbn",
		"date": 1631773887,
		"priority": 0,
		"acked": 0,
		"umid": 502550228871720147,
		"umid_str": "502550228871720147",
		"title": "Testing transport from LibreNMS",
		"url": "mailto:<mail>",
		"url_title": "Reply to <mail>",
		"queued_date": 1631773893,
		"dispatched_date": 1631773893
	*/
}

type device struct {
	Name string `json:"name"`
	/*
		"encryption_enabled": false,
		"default_sound": "po",
		"always_use_default_sound": false,
		"default_high_priority_sound": "po",
		"always_use_default_high_priority_sound": false,
		"dismissal_sync_enabled": false
	*/
}

type respJSON struct {
	Message []messages `json:"messages"`
	Status  int        `json:"status"`
	Request string     `json:"request"`
	Device  device     `json:"device"`
	/*
		"user": {
			"quiet_hours": false,
			"is_android_licensed": true,
			"is_ios_licensed": false,
			"is_desktop_licensed": true,
			"email": "<mail>",
			"show_team_ad": "1"
		},
	*/
}
