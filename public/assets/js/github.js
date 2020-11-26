$(function() {
    $('.datetime').each(function()
    {
        $(this).html(moment($(this).text()).format('YYYY-MM-DD HH:mm'));
    });

    $('#releases').DataTable({
        'pageLength': 25,
        'order':      [[ 3, 'desc' ]],
    });
});
