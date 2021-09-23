
function setAngle() {
    var angle = $("#slider").roundSlider("option", "value");
    $.ajax({
        url: "servo",
        type: "get",
        data: {
            angle: angle
        }
    });
}

$(document).ready(function() {
    $("#slider").roundSlider({
        radius: 360,
        circleShape: "custom-quarter",
        min: -45,
        max: 45,
        value: 0,
        startAngle: "45",
        handleShape: "dot",
        update: "setAngle"
    }); 
});