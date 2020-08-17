getCommandPage("create");

$('#list-tab a').on('click', function (e) {
    e.preventDefault();
    $(this).tab('show');
    var idArr = $(this).attr('id').split("-");

    getCommandPage(idArr[1]);
});

function getCommandPage(command) {
    $.get("commands/" + command + ".html", function (data) {
        $('#list-' + command).html(data);
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


