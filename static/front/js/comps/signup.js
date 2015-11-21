$(document).ready(function() {

  $.get('/static/signup_agreement.txt', function(result) {
    $('#register_agreement').html(result);
  }, 'text');

  $('.aw-agreement-btn').click(function() {
    if ($('.aw-register-agreement').is(':visible')) {
      $('.aw-register-agreement').hide();
    } else {
      $('.aw-register-agreement').show();
    }
  });

  $('.more-information-btn').click(function() {
    $('.more-information').fadeIn();
    $(this).parent().hide();
  });

  verify_register_form('#register_form');

  /* 注册页面验证 */
  function verify_register_form(element) {
    $(element).find('[type=text], [type=password]').on({
      focus: function() {
        if (typeof $(this).attr('tips') != 'undefined' && $(this).attr('tips') != '') {
          $(this).parent().append('<span class="aw-reg-tips">' + $(this).attr('tips') + '</span>');
        }
      },
      blur: function() {
        if ($(this).attr('tips') != '') {
          switch ($(this).attr('name')) {
            case 'user_name':
              var _this = $(this);
              $(this).parent().find('.aw-reg-tips').detach();
              if ($(this).val().length >= 0 && $(this).val().length < 2) {
                $(this).parent().find('.aw-reg-tips').detach();
                $(this).parent().append('<span class="aw-reg-tips aw-reg-err"><i class="aw-icon i-err"></i>' + $(this).attr('errortips') + '</span>');
                return;
              }
              if ($(this).val().length > 17) {
                $(this).parent().find('.aw-reg-tips').detach();
                $(this).parent().append('<span class="aw-reg-tips aw-reg-err"><i class="aw-icon i-err"></i>' + $(this).attr('errortips') + '</span>');
                return;
              } else {
                $.post('/api/account/check/', {
                  username: $(this).val()
                }, function(result) {
                  if (result.errno == -1) {
                    _this.parent().find('.aw-reg-tips').detach();
                    _this.parent().append('<span class="aw-reg-tips aw-reg-err"><i class="aw-icon i-err"></i>' + result.err + '</span>');
                  } else {
                    _this.parent().find('.aw-reg-tips').detach();
                    _this.parent().append('<span class="aw-reg-tips aw-reg-right"><i class="aw-icon i-followed"></i></span>');
                  }
                }, 'json');
              }
              return;

            case 'email':
              $(this).parent().find('.aw-reg-tips').detach();
              var emailreg = /^\w+((-\w+)|(\.\w+))*\@[A-Za-z0-9]+((\.|-)[A-Za-z0-9]+)*\.[A-Za-z0-9]+$/;
              if (!emailreg.test($(this).val())) {
                $(this).parent().find('.aw-reg-tips').detach();
                $(this).parent().append('<span class="aw-reg-tips aw-reg-err"><i class="aw-icon i-err"></i>' + $(this).attr('errortips') + '</span>');
                return;
              } else {
                $(this).parent().find('.aw-reg-tips').detach();
                $(this).parent().append('<span class="aw-reg-tips aw-reg-right"><i class="aw-icon i-followed"></i></span>');
              }
              return;

            case 'password':
              $(this).parent().find('.aw-reg-tips').detach();
              if ($(this).val().length >= 0 && $(this).val().length < 6) {
                $(this).parent().find('.aw-reg-tips').detach();
                $(this).parent().append('<span class="aw-reg-tips aw-reg-err"><i class="aw-icon i-err"></i>' + $(this).attr('errortips') + '</span>');
                return;
              }
              if ($(this).val().length > 17) {
                $(this).parent().find('.aw-reg-tips').detach();
                $(this).parent().append('<span class="aw-reg-tips aw-reg-err"><i class="aw-icon i-err"></i>' + $(this).attr('errortips') + '</span>');
                return;
              } else {
                $(this).parent().find('.aw-reg-tips').detach();
                $(this).parent().append('<span class="aw-reg-tips aw-reg-right"><i class="aw-icon i-followed"></i></span>');
              }
              return;
          }
        }

      }
    });
  }

  $('.select_area').LocationSelect({
    labels: ["请选择省份或直辖市", "请选择城市"],
    elements: document.getElementsByTagName("select"),
    detector: function() {
      //this.select(["{{u.Province}}", "{{u.City}}"]);
    },
    dataUrl: '/static/front/js/areas.js'
  });
});
