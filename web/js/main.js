const baseUrl = $(location).attr('origin');

let loginPassword = $('#loginPassword')
let registerPassword = $('#registerPassword')

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
    login('click')
})

loginPassword.on('keydown', function () {
    if (event.key === 'Enter') {
        login('keydown')
    }
})

$('.register').on('click', function () {
    register('click')
})

registerPassword.on('keydown', function () {
    if (event.key === 'Enter') {
        register('keydown')
    }
})

$('.logout').on('click', function () {
    let post = {
        action: 'logout',
    }
    sendPost(post, '/logout')
})

function register(action) {
    let post = {
        login: $('#registerUsername')[0].value,
        password: registerPassword[0].value
    }
    if (action === 'click' || action === 'keydown') {
        sendPost(post, '/register')
    } else {
        alert("Все поля должны быть заполнены")
    }
}

function login(action){
    let post = {
        login: $('#loginUsername')[0].value,
        password: loginPassword[0].value
    }
    if (action === 'click' || action === 'keydown') {
        if (post.login && post.password) {
            sendPost(post, '/login')
        } else {
            alert("Все поля должны быть заполнены")
        }
    }
}

function sendPost(post, url) {
    $.ajax({
        url: url,
        type: "POST",
        contentType: 'application/json',
        data: JSON.stringify(post),
        dataType: 'json',
        async: false,
        success: function(data) {
            if (data.redirect){
                window.location.replace(baseUrl + data.redirect)
            }
        }
    });
}