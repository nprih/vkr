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

function refreshImages(isadmin = false) {
    $.get(baseUrl + '/images', function(data) {
        if (data) {
            if (isadmin){
                $('.admin-photo').html(data);
            } else {
                $('.table-wrapper').html(data);
            }
        }
    })
}

function refreshUserImages(selectedUserId = 0) {
    $.get(baseUrl + '/images/admin/' + selectedUserId, function(data) {
        if (data) {
            $('.admin-photo').html(data);
        }
    })
}

function viewImage(src) {
    console.log(src)
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

function deleteImage(imageId, isAdmin = false, selectedUserId = 0){
    $.ajax({
        url: '/images/' + imageId,
        type: 'DELETE',
        success: function(response) {

            if (isAdmin){
                refreshUserImages(selectedUserId)
            } else {
                refreshImages()
            }

        },
        error: function(xhr, status, error) {
            console.error('Ошибка:', error);
        }
    });
}

function deleteUser(userId) {
    $.ajax({
        url: '/user/' + userId,
        type: 'DELETE',
        success: function(response) {
            console.log(response)

            // if (isAdmin){
            //     refreshUserImages(selectedUserId)
            // } else {
            //     refreshImages()
            // }

        },
        error: function(xhr, status, error) {
            console.error('Ошибка:', error);
        }
    });
}