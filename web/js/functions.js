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

function getImages() {
    $.get(baseUrl + '/images', function(data) {
        if (data) {
            refreshImages(data)
        }
    })
}

function viewImage(src) {
    document.getElementById('lightboxImg').src = src;
    document.getElementById('lightbox').style.display = 'flex';
    document.body.style.overflow = 'hidden';
}

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

function refreshImages(table) {
    $('.table-wrapper').html(table);
}

function formatImages(images) {
    console.log("images: ", images)
    let rows = [];
    let newImages = `
                    <table id="adminPhotosTable">
                    <thead>
                        <tr><th>Превью</th><th>Автор</th><th>Дата</th><th>Действия</th></tr>
                    </thead>
                    <tbody id="adminPhotosBody">
                `
    $.each(images, function (index, value) {
        rows.push(`
                        <tr>
                            <td id="${value.id}" title="Нажмите для просмотра">
                                <img class="thumb-mini" src="${value.url}" alt="миниатюра" style="cursor:pointer;">
                            </td>
                            <td>${value.author}</td>
                            <td>${value.createdAt}</td>
                            <td>
                            <button class="btn btn-danger btn-sm">🗑 Удалить</button>
                            </td>
                        </tr>
                    `);
    })
    newImages += rows.join('');
    newImages += '</table>';
    return newImages
}

function deleteImage(imageId){
    $.ajax({
        url: '/images/' + imageId,
        type: 'DELETE',
        success: function(response) {
            console.log('Удалено:', response);
            $('tr#' + imageId).remove()
        },
        error: function(xhr, status, error) {
            console.error('Ошибка:', error);
        }
    });
}