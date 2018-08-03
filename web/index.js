import { h, render } from 'preact';

import { DoorMonitor } from './components'

render((
    <div id="app">
        <span>Hello, world!</span>
        <DoorMonitor />
    </div>
), document.body);


window.onload = function () {
    var conn;
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + process.env.API_BASE + "/sub");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                item.innerText = messages[i];
                appendLog(item);
            }
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};