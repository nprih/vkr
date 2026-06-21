const baseUrl = $(location).attr('origin');

let loginPassword = $('#loginPassword')
let registerPassword = $('#registerPassword')

let pathname = $(location).attr('pathname');
let selected

/** Переход на главную по нажатию на логотип */
$('.logo').on('click', function () {
    window.location.replace(baseUrl + '/')
})
/** Смена "активности" вкладок меню */
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

/** Отправка логин/пароля по нажатию кнопки */
$('.login').on('click', function () {
    login('click')
})
/** Отправка логин/пароля по нажатию "Enter" находясь в поле "password" */
loginPassword.on('keydown', function () {
    if (event.key === 'Enter') {
        login('keydown')
    }
})

/** Отправка логин/пароля по нажатию кнопки */
$('.register').on('click', function () {
    register('click')
})
/** Отправка логин/пароля по нажатию "Enter" находясь в поле "password" */
registerPassword.on('keydown', function () {
    if (event.key === 'Enter') {
        register('keydown')
    }
})

/** Отправка запроса на прекращение авторизации */
$('.logout').on('click', function () {
    let post = {
        action: 'logout',
    }
    sendPost(post, '/logout')
})

/** Отправка файла для сохранения */
$('#fileInput').on('change', function () {
    const fileInput = $(this)[0];
    const file = fileInput.files[0];
    if (!file) {
        return;
    }

    const formData = new FormData();
    formData.append('image', file, file.name);
    $.ajax({
        url: '/upload',
        type: 'POST',
        data: formData,
        processData: false,
        contentType: false,
        success: function(res) {
            refreshImages()
        },
        error: function(xhr, status, err) {
            console.error('Ошибка:', status, err);
        }
    });
});

/** Вызов модалки для просмотра картинки на главной */
$('.photo-card').on('click',  function () {
    viewImage($(this).children('img').attr('src'))
})

/** Вызов модалки для просмотра картинки в профиле */
$('.table-wrapper').on('click', '.thumb-mini', function() {
    viewImage($(this).attr('src'));
});

/** Удаление в профиле по нажатию кнопки удаления*/
$('.profile-photo').on('click', '.delete_image', function () {
    if (confirm('Вы действительно хотите удалить изображение?')) {
        deleteImage($(this).attr('id'))
    }
})

/** Удаление в админ-панели по нажатию кнопки удаления*/
$('.admin-photo').on('click', '.delete_image', function () {
    if (confirm('Вы действительно хотите удалить изображение?')) {
        deleteImage($(this).attr('id'), true)
    }
})

$('.select').on('click', function () {
    let id = $(this).attr('id')
    let row = $('.row-user#' + id)
    let user = $('td#'+ id +'.select.username').text()

    if (row.hasClass('active')) {
        row.removeClass('active')
        $('.description').text('🖼️ Все фотографии')
        refreshUserImages(0)
    } else {
        $('.row-user').removeClass('active')
        row.addClass('active')
        $('.description').text('🖼️ Все фотографии пользователя ' + user)
        refreshUserImages(id)
    }

})

function refreshUserImages(id = 0) {
    $.get(baseUrl + '/images/admin/' + id, function(data) {
        if (data) {
            $('.admin-photo').html(data);
        }
    })
}

$('.delete_user').on('click', function () {
    console.log('delete: ' + $(this).attr('id'))
})