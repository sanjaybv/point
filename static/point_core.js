function init() {

	var ws;

	if ("WebSocket" in window) {

		ws = new WebSocket("ws://localhost:8082/");
		
		ws.onopen = function() {
			console.log("point: ws opened");
		}

		ws.onmessage = function(evt) {

			var msg = JSON.parse(evt.data);
			console.log("point: message recieved:");

			fromX = parseInt(msg.fromX);
			fromY = parseInt(msg.fromY);
			toX = parseInt(msg.toX);
			toY = parseInt(msg.toY);

			drawLine(fromX, fromY, toX, toY);
		}

		ws.onclose = function(e) {
			console.log("point: connection closed");
		}

	} else {
		alert("WebSocket is not supported!");
	}
	
	var qPt = $("#pointspace");

	qPt.attr("width", $("#sketch").width());

	var mouse = {x: 0, y: 0, ox: 0, oy: 0};
	var pointSpace = document.getElementById("pointspace");
	ctx = pointSpace.getContext('2d');

	ctx.lineWidth = 3;
	ctx.lineJoin = 'round';
	ctx.lineCap = 'round';
	ctx.strokeStyle = 'black';

	qPt.mousemove(function(e) {
		mouse.ox = mouse.x
		mouse.oy = mouse.y
		mouse.x = e.pageX - qPt.offset().left;
		mouse.y = e.pageY - qPt.offset().top;
	});

	qPt.mousedown(function(e) {
		ctx.beginPath();
		ctx.moveTo(mouse.x, mouse.y);
		qPt.mousemove(onPaint);
	});

	qPt.mouseup(function() {
		qPt.unbind('mousemove', onPaint);
	});

	var onPaint = function() {
		drawLine(mouse.ox, mouse.oy, mouse.x, mouse.y);
		ws.send(JSON.stringify({
			fromX: mouse.ox,
			fromY: mouse.oy,
			toX: mouse.x,
			toY: mouse.y
		}));
	};

	function drawLine(fromX, fromY, toX, toY) {
		ctx.beginPath();
		ctx.moveTo(fromX, fromY);
		ctx.lineTo(toX, toY);
		ctx.stroke();
	}

	$("#clear").click(function () {

		pointSpace.width = pointSpace.width;

		ctx.lineWidth = 3;
		ctx.lineJoin = 'round';
		ctx.lineCap = 'round';
		ctx.strokeStyle = 'black';
	});
}
