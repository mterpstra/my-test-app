$( document ).ready(function() {



    /*
    $("body").animate({opacity:'0'}, 2000, function() {
        $(this).css({color:'green'}).animate({opacity:'1'}, 2000, function() {
        })
    });
    */
    GetRestaurants()

   $(document).on("click", ".restaurant-name", GetItems)


});


function GetRestaurants() 
{
    $.get("/restaurants", ShowRestaurants);
}


function ShowRestaurants(data)
{
    var template = $('#restaurants-template').html();
    var info = Mustache.to_html(template, data);
    $('#restaurant-list').html(info);
}

function GetItems() 
{
    id = $(this).data("restaurant");
    $.get("/restaurants/"+id+"/items", ShowItems);
}


function ShowItems(data)
{
    var template = $('#items-template').html();
    var info = Mustache.to_html(template, data);
    $('#items-list').html(info);
}
