$(document).ready(() => {
    $('#shorten_actions').hide();
    $('#href').append(location.href);

    $('#button-shorten').on('click', () => {
        let longURL = $('#longurl').val();
        let customcode = $('#shortcode-custom').val();
        if (longURL !== '') {
            if (!longURL.match(/^[a-zA-Z]+:\/\//)) {
                longURL = 'http://' + longURL;
            }
            let obj = { "longurl": longURL };
            if (customcode != "") {
                obj["custom"] = true;
                obj["customcode"] = customcode;
            }
            else
                obj["custom"] = false
            $.ajax({
                type: "POST",
                url: "/url/",
                data: JSON.stringify(obj),
                contentType: "application/json; charset=utf-8",
                success: function (res) {
                    let shortURL = $('#short-link a');
                    $('#long-link').empty();
                    shortURL.empty();
                    if (res.code == 200 || res.code == 201) {
                        $('#long-link').append(longURL);
                        shortURL.append(res.data);
                        shortURL.attr('href', res.data);
                    }
                    if (res.code == 503) {
                        $('#long-link').append("Hệ thống có lỗi, vui lòng thử lại hoặc liên hệ MISA");
                    }
                    if (res.code == 409) {
                        $('#long-link').append("Link rút gọn đã tồn tại vui lòng chọn link khác");
                    }
                    $('#shorten_actions').show();
                }
            });
        }
    })

    $('#button-option').on('click', () => {
        let option = $('#shortcode-option');
        if (option.is(":visible")) {
            option.hide();
            $('#shortcode-custom').val("");
        }
        else
            option.show();
    })
    var clipboard = new ClipboardJS('#copy');
})
