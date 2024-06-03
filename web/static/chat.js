// Get the history and populate the "history" div
$.getJSON("/list", function (data, status) {
    $.each(data.Chats, function (i, field) {
        $("#chatList").append(`<a class="w3-bar-item w3-button" href="#${field.Token}" > ${field.Title} </a> ` + "<br>");
    });

});

var converter = new showdown.Converter()
function displayHistory() {
    // Get the anchor and populate the history
    let anchor = window.location.hash.substring(1);
    if (anchor !== "") {
        $.getJSON(`/history?token=${anchor}`, function (data, status) {
            $("#chatTitle").text(`${data.Title}`);
            $("#chatContainer").html("");
            $.each(data.Messages, function (i, field) {
                
                let html = `<div class="w3-row">`
                if (field.Who === "user") {
                    // html +=  <div
                    html += `<div class="w3-col l6 w3-padding-large"></div>
                    <div class="w3-col l6 w3-padding-large w3-right-align">
                        <div class="w3-card w3-round-large w3-padding" style="display: inline-block;">
                            <p><b>${field.Content}</b></p>
                        </div>
                    </div>`
                } else {
                    html += `<div class="w3-padding-large">
                        <div class="w3-card w3-round-large w3-padding" style="display: inline-block;">
                            <p>${converter.makeHtml(field.Content)}</p>
                        </div>
                    </div>
                    `
                }
                html += `</div>`
                $("#chatContainer").append(html)
            });
        });
    }

}


$(document).ready(function () {
    $("#chatList").click(function () {
        // delay 20ms
        setTimeout(function () {
            displayHistory();
        }, 20);
    });
});
displayHistory()