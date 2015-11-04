// category.js

$(function (){
    cate = getURLParameter("v");
    $("ul.sidebar-menu li a").each(function (){
        if ($(this).attr("href").indexOf(decodeURI(cate)) !== -1) {
            $(this).addClass("choosen");
        }
    });
})





