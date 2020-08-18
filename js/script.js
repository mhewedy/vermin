var landingPage = window.location.href.split("-")[1] || "create";
getCommandPage(landingPage);

$('#list-tab a').on('click', function (e) {

    e.preventDefault();
    $(this).tab('show');
    var idArr = $(this).attr('id').split("-");

    window.location.href = '#' + idArr[0] + '-' + idArr[1];
    getCommandPage(idArr[1]);
});

function getCommandPage(command) {
    $.get("commands/" + command + ".html", function (data) {
        var ele = $('#list-' + command);
        ele.html(data);
        ele.tab('show');

        $('#list-' + command + '-list').addClass('active');

        processLinks();
    });
}


/**
 * find elements of special links and convert them to actaul links to command page.
 *
 * e.g. find elements of type <link-stop></link-stop> and convert it to a link to stop page
 */
function processLinks() {
    $(document).find('*').each(function () {
        var localName = $(this)[0].localName;
        if (localName.startsWith("link-")) {

            var parts = localName.split("-");
            var txt = $(this).text() || parts[1];

            $(this).empty();
            $(this).append('<a href="#" onclick="$(\'#list-' + parts[1] + '-list\').click()">' + txt + '</a>');
        }
    })
}


