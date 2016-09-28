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

        ws = new WebSocket("ws://" + location.host + "/ws");
        ws.onopen = function() {
            console.log("ddchat: ws opened");
        }

        ws.onmessage = function(evt) {
            var msg = JSON.parse(evt.data);
            console.log("ddchat: message recieved:");
            //console.log(msg);

            if (msg["type"] == "line") {

                fromX = parseInt(msg.fromX);
                fromY = parseInt(msg.fromY);
                toX = parseInt(msg.toX);
                toY = parseInt(msg.toY);

                drawLine(fromX, fromY, toX, toY);
                return
            }

            t1[3] = msg.username;
            var currentdate = new Date(); 
            var datetime = currentdate.getDate() + "/"
                + (currentdate.getMonth()+1)  + "/" 
                + currentdate.getFullYear() + " @ "  
                + currentdate.getHours() + ":"  
                + currentdate.getMinutes() + ":" 
                + currentdate.getSeconds();
            t1[5] = datetime;
            t1[7] = msg.body;

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

        var msg = {}
        msg["username"] = $("#username").val();
        msg["body"] = $("#body").val();
        msg["type"] = "chat";
        $("#body").val("");
        $("#body").focus();

        ws.send(JSON.stringify(msg));
    };

    $("#send").click(sendFunc);
    $("#msgform").submit(function() {
        sendFunc();
        return false;
    });

    // drawing stuff
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
        mouse.ox = mouse.x;
        mouse.oy = mouse.y;
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
            "type": "line",
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


    /*
       $("#msgform").keydown(function(event) {
       if (event.keyCode == 13) {
       sendFunc();
       return false;
       }
       });
       */

}

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

        }

        ws.onclose = function(e) {
            console.log("point: connection closed");
        }

    } else {
        alert("WebSocket is not supported!");
    }
}
