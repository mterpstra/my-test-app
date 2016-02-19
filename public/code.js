$( document ).ready(function() {

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
    $(".restaurant-name").parent().removeClass('selected');
    $(this).parent().addClass('selected');
    id = $(this).data("restaurant");
    $.get("/restaurants/"+id+"/items", ShowItems);
}


function ShowItems(data)
{
    var template = $('#items-template').html();
    var info = Mustache.to_html(template, data);
    $('#items-list').html(info);
}
