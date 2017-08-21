let ws = new WebSocket('ws://localhost:4242/ws');

let canvas = document.getElementById("canvas");

let count = 0;
let globalArray = [ ];

function draw(graph) {
    $.plot($("#placeholder"), [graph]);
}

ws.addEventListener('message', function(e) {
    let msg = JSON.parse(e.data);

    globalArray.push( [++count, msg["value"]] );

    console.log(globalArray);
    draw(globalArray);

});
