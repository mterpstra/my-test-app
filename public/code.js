$( document ).ready(function() {
    $("body").animate({opacity:'0'}, 2000, function() {
        $(this).css({color:'green'}).animate({opacity:'1'}, 2000, function() {
        })
    });
});

