
function toggleLed() {
    $.ajax({
        url: "devices/led1"
    }).then(function(data) {
        var json = $.parseJSON(data);
        $("#LED1-state").html("LED turned ON: " + json.State);
    });
}

$(document).ready(function() {
    $("#LED1-button").click(toggleLed); 
});