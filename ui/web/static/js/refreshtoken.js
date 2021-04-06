function refreshTokenIfNeed() {
    var accessToken = $.cookie("access_token");
    var refreshToken = $.cookie("refresh_token");

    var accessTokenSplit = JSON.parse(atob(accessToken.split('.')[1]));
    alert(accessTokenSplit.exp - (new Date().getTime() + 1) / 1000);
    if (accessTokenSplit.exp <= (new Date().getTime() + 1) / 1000) {
        $.ajax({
            'url': '/api/auth/refresh-tokens',
            'method': 'POST',
            'dataType': 'json',
            'contentType': 'application/json',
            'data': JSON.stringify({
                "refresh_token": refreshToken
            }),
            success: function (response) { //Данные отправлены успешно
                $.cookie("access_token", response.access_token);
                $.cookie("refresh_token", response.refresh_token);

            },
            error: function (response) { // Данные не отправлены
                alert("error");
                window.location.href = "/api/login";
            }
        })
    }
}
