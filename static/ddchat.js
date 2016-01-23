function makeConn() {

		var ws;
		var t1 = [
				'<div class="comment">',
				'<!--a class="avatar"><img src="static/blank.png"></a-->',
				'<div class="content"><a class="author">',
				'',
				'</a><div class="metadata"><span class="date">',
				'time',
				'</span></div><div class="text">',
				'asdfasdfasdf',
				'</div></div></div>'
						];

		if ("WebSocket" in window) {

				ws = new WebSocket("ws://localhost:8081/");
				ws.onopen = function() {
						console.log("ddchat: ws opened");
				}
				ws.onmessage = function(evt) {

						var msg = JSON.parse(evt.data);
						console.log("ddchat: message recieved:");
						//console.log(msg);

						t1[3] = msg.username;
						t1[5] = msg.time;
						t1[7] = emoji(
										msg.body, 
										"static/emoji-images/pngs", 
										23
									 );
						if (msg.body.trim() == "") {
								t1[7] = "!@#$%^&*";
						}
						if (msg.username.trim() == "") {
								t1[3] = "!@#$%^&*";
						}
						$("#conv").append(t1.join(""));
						$("#conv").scrollTop($("#conv")[0].scrollHeight);
				}
				ws.onclose = function() {
						console.log("ddchat: connection closed");
				}
		} else {

				alert("WebSocket is not supported!");
		}

		sendFunc = function() {

				var username = $("#username").val();
				var body = $("#body").val();
				$("#body").val("");
				$("#body").focus();
				var msg = {"body":body, "username":username};

				ws.send(JSON.stringify(msg));
		};

		$("#send").click(sendFunc);
		$("#msgform").submit(function() {
				sendFunc();
				return false;
		});
		/*
		   $("#msgform").keydown(function(event) {
		   if (event.keyCode == 13) {
		   sendFunc();
		   return false;
		   }
		   });
		   */

}
