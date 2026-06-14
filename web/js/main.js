const baseUrl = $(location).attr('origin');

let loginPassword = $('#loginPassword')
let registerPassword = $('#registerPassword')

let pathname = $(location).attr('pathname');
let selected

$('.logo').on('click', function () {
    window.location.replace(baseUrl + '/')
})

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
    case "/profile":
        selected = $('.nav-link:contains("Профиль")');
        break;
    case "/admin":
        selected = $('.nav-link:contains("Админ-панель")');
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
            loginPassword.val('')
        } else {
            alert("Все поля должны быть заполнены")
        }
    }
}

$('#fileInput').on('change', function () {
    const fileInput = $(this)[0];
    const file = fileInput.files[0];
    if (!file) {
        console.warn('Выберите файл');
        return;
    }

    const formData = new FormData();
    formData.append('image', file, file.name); // поле 'image' будет ключом на сервере
    // console.log(file)
    $.ajax({
        url: '/upload',
        type: 'POST',
        data: formData,
        processData: false, // критично: не превращать FormData в строку
        contentType: false, // критично: не ставить application/x-www-form-urlencoded
        success: function(res) {
            console.log('Успех:', res);
            window.location.replace(baseUrl + '/profile')
        },
        error: function(xhr, status, err) {
            console.error('Ошибка:', status, err);
        }
    });
});

$('.photo-card').on('click', function () {
    document.getElementById('lightboxImg').src = $(this).children('img').attr('src');
    document.getElementById('lightbox').style.display = 'flex';
    document.body.style.overflow = 'hidden';
})

function closeLightbox() {
    document.getElementById('lightbox').style.display = 'none';
    document.body.style.overflow = '';
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
            if (data.fail){
                alert(data.fail)
            }
            if (data.redirect){
                window.location.replace(baseUrl + data.redirect)
            }
        }
    });
}