let pathname = $(location).attr('pathname');
let selected
$('.nav-link').removeClass('active');
switch (pathname) {
    case "/":
        selected = $('.nav-link:contains("Главная")');
        break;
    case "/login":
        selected = $('.nav-link:contains("Войти")');
        break;
    case "/register":
        selected = $('.nav-link:contains("Регистрация")');
        break;
}
selected.addClass('active');

$('.login').on('click', function () {
    let post = {
        login: $('#loginUsername')[0].value,
        password: $('#loginPassword')[0].value
    }
    console.log(post)
    $.ajax({
        url: '/login',
        type: "POST",
        contentType: 'application/json',
        data: JSON.stringify(post),
        dataType: 'json',
        async: false,
        success: function(data) {
            console.log(data);
        }
    });
})