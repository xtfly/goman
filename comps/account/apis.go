package account

import (
	"fmt"
	"strings"

	"github.com/Unknwon/com"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/xtfly/goman/boot"
	"github.com/xtfly/goman/comps"
	"github.com/xtfly/goman/kits"
	"github.com/xtfly/goman/models"
	"gopkg.in/macaron.v1"
)

// POST /api/account/check/
func ApiCheckUserName(c *macaron.Context) {
	un := strings.TrimSpace(c.QueryEscape("username"))
	if err, ok := CheckUsernameChar(un); !ok {
		c.JSON(200, comps.NewRestErrResp(-1, err))
		return
	}

	if CheckUsernameSensitiveWords(un) || models.UserExistedByName(un) {
		c.JSON(200, comps.NewRestErrResp(-1, "用户名已被注册"))
		return
	}

	c.JSON(200, comps.NewRestErrResp(1, ""))
}

type SignupForm struct {
	Name      string `form:"user_name"`
	Email     string `form:"email"`
	Password  string `form:"password"`
	Gender    int8   `form:"gender"`
	JobId     string `form:"job_id"`
	Province  string `form:"province"`
	City      string `form:"city"`
	Signature string `form:"signature"`
	CsrfToken string `form:"_csrf"`
	ICode     string `form:"icode"`
	IEmail    string `form:"invitation_email"`
	Agree     bool   `form:"agreement_chk"`
}

// POST /api/account/signup/
func ApiUserSignup(c *macaron.Context, f SignupForm, cpt *captcha.Captcha, x csrf.CSRF, s *session.Flash) {
	if x.ValidToken(f.CsrfToken) {
		c.JSON(200, comps.NewRestErrResp(-1, "非法的跨站请求"))
		return
	}

	// if input invalid, will store user_name & email
	defer func() {
		s.Set("email", f.Email)
		s.Set("user_name", f.Name)
	}()

	ra := boot.SysSetting.Ra
	switch ra.RegisterType {
	case models.RegTypeClose:
		c.JSON(200, comps.NewRestErrResp(-1, "本站目前关闭注册"))
		return
	case models.RegTypeInvite:
		if f.ICode == "" {
			c.JSON(200, comps.NewRestErrResp(-1, "本站只能通过邀请注册"))
		}
	default:
		break
	}

	if !cpt.VerifyReq(c.Req) {
		c.JSON(200, comps.NewRestErrResp(-1, "请填写正确的验证码"))
		return
	}

	if err, ok := CheckUsernameChar(f.Name); !ok {
		c.JSON(200, comps.NewRestErrResp(-1, err))
		return
	}

	if CheckUsernameSensitiveWords(f.Name) || models.UserExistedByName(f.Name) {
		c.JSON(200, comps.NewRestErrResp(-1, "用户名已被注册或包含敏感词或系统保留字"))
		return
	}

	if kits.IsEmail(f.Email) || models.UserExistedByEmail(f.Email) {
		c.JSON(200, comps.NewRestErrResp(-1, "EMail 已经被使用, 或格式不正确"))
		return
	}

	if len(f.Password) < 6 {
		c.JSON(200, comps.NewRestErrResp(-1, "密码长度不符合规则"))
		return
	}

	if !f.Agree {
		c.JSON(200, comps.NewRestErrResp(-1, "你必需同意用户协议才能继续"))
		return
	}

	if f.Gender < models.GenderUnknown || f.Gender > models.GenderFemale {
		c.JSON(200, comps.NewRestErrResp(-1, "非法的性别输入"))
		return
	}

	u := &models.Users{}
	u.UserName = f.Name
	u.Email = f.Email
	u.Password = f.Password
	u.Gender = f.Gender
	u.JobId = com.StrTo(f.JobId).MustInt64()
	u.Province = f.Province
	u.City = f.City
	u.Signature = f.Signature
	u.RegIp = c.RemoteAddr()
}

/*



  $uid = $this->model('account')->user_register($_POST['user_name'], $_POST['password'], $_POST['email']);


if ($_POST['email'] == $invitation['invitation_email'])
{
  $this->model('active')->set_user_email_valid_by_uid($uid);

  $this->model('active')->active_user_by_uid($uid);
}

if (isset($_POST['sex']))
{
  $update_data['sex'] = intval($_POST['sex']);

  if ($_POST['province'])
  {
    $update_data['province'] = htmlspecialchars($_POST['province']);
    $update_data['city'] = htmlspecialchars($_POST['city']);
  }

  if ($_POST['job_id'])
  {
    $update_data['job_id'] = intval($_POST['job_id']);
  }

  $update_attrib_data['signature'] = htmlspecialchars($_POST['signature']);

  // 更新主表
  $this->model('account')->update_users_fields($update_data, $uid);

  // 更新从表
  $this->model('account')->update_users_attrib_fields($update_attrib_data, $uid);
}

$this->model('account')->setcookie_logout();
$this->model('account')->setsession_logout();

if ($_POST['icode'])
{
  $follow_users = $this->model('invitation')->get_invitation_by_code($_POST['icode']);
}
else if (HTTP::get_cookie('fromuid'))
{
  $follow_users = $this->model('account')->get_user_info_by_uid(HTTP::get_cookie('fromuid'));
}

if ($follow_users['uid'])
{
  $this->model('follow')->user_follow_add($uid, $follow_users['uid']);
  $this->model('follow')->user_follow_add($follow_users['uid'], $uid);

  $this->model('integral')->process($follow_users['uid'], 'INVITE', get_setting('integral_system_config_invite'), '邀请注册: ' . $_POST['user_name'], $follow_users['uid']);
}

if ($_POST['icode'])
{
  $this->model('invitation')->invitation_code_active($_POST['icode'], time(), fetch_ip(), $uid);
}

if (get_setting('register_valid_type') == 'N' OR (get_setting('register_valid_type') == 'email' AND get_setting('register_type') == 'invite'))
{
  $this->model('active')->active_user_by_uid($uid);
}

$user_info = $this->model('account')->get_user_info_by_uid($uid);

if (get_setting('register_valid_type') == 'N' OR $user_info['group_id'] != 3 OR $_POST['email'] == $invitation['invitation_email'])
{
  $this->model('account')->setcookie_login($user_info['uid'], $user_info['user_name'], $_POST['password'], $user_info['salt']);

  if (!$_POST['_is_mobile'])
  {
    H::ajax_json_output(AWS_APP::RSM(array(
      'url' => get_js_url('/home/first_login-TRUE')
    ), 1, null));
  }
}
else
{
  AWS_APP::session()->valid_email = $user_info['email'];

  $this->model('active')->new_valid_email($uid);

  if (!$_POST['_is_mobile'])
  {
    H::ajax_json_output(AWS_APP::RSM(array(
      'url' => get_js_url('/account/valid_email/')
    ), 1, null));
  }
}

if ($_POST['_is_mobile'])
{
  if ($_POST['return_url'])
  {
    $user_info = $this->model('account')->get_user_info_by_uid($uid);

    $this->model('account')->setcookie_login($user_info['uid'], $user_info['user_name'], $_POST['password'], $user_info['salt']);

    $return_url = strip_tags($_POST['return_url']);
  }
  else
  {
    $return_url = get_js_url('/m/');
  }

  H::ajax_json_output(AWS_APP::RSM(array(
    'url' => $return_url
  ), 1, null));
}
*/

func CheckUsernameChar(un string) (string, bool) {
	if kits.IsDigit(un) {
		return "用户名不能为纯数字", false
	}

	if strings.ContainsAny(un, "-./") || strings.Contains(un, "__") {
		return "用户名不能包含 - / . % 与连续的下划线", false
	}

	unlen := len(un)
	min := boot.SysSetting.Ra.UsernameLenMin
	max := boot.SysSetting.Ra.UsernameLenMax
	if unlen < min || unlen > max {
		return fmt.Sprintf("用户名长度只能在[%d,%d]", min, max), false
	}

	switch boot.SysSetting.Ra.UsernameRule {
	case models.UserRuleNotLimit:
		break
	case models.UserRuleChineseLetterNumUnline:
		if !kits.IsChineseLetterNumUnline(un) {
			return fmt.Sprintf("请输入大于[%d,%d]字节的用户名, 允许汉字、字母与数字", min, max), false
		}
	case models.UserRuleLetterNumUnline:
		if !kits.IsLetterNumUnline(un) {
			return fmt.Sprintf("请输入[%d,%d]个字母、数字或下划线", min, max), false
		}
	case models.UserRuleChinese:
		if !kits.IsChinese(un) {
			return fmt.Sprintf("请输入[%d,%d]个汉字", min/2, max/2), false
		}
	default:
		break
	}

	return "", true
}

//检查用户名中是否包含敏感词或用户信息保留字
func CheckUsernameSensitiveWords(un string) bool {
	if kits.SensitiveWordExists(un, boot.SysSetting.Cs.SensitiveWords) {
		return true
	}

	cs := boot.SysSetting.Ra.CensorUser
	if len(cs) == 0 {
		return false
	}

	css := strings.Split(cs, "\n")
	for _, c := range css {
		if strings.ToLower(un) == strings.ToLower(strings.TrimSpace(c)) {
			return true
		}
	}

	return false
}
