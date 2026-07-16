package storage

import (
	"strings"
	"testing"
)

func TestCHeckQuotaNotifiesUser(t *testing.T) {
	// Сохранение и восстановление исходного значения notifyUser.
	saved := notifyUser
	defer func() { notifyUser = saved }()

	// Установка поддельной функции для notifyUser.
	var notifiedUser, notifiedMsg string
	notifyUser = func(user, msg string) {
		notifiedUser, notifiedMsg = user, msg
	}
	// ... имитация условия 980-Мбайтной занятости ...
	const user = "joe@example.org"
	CheckQuota(user)
	if notifiedUser == "" && notifiedMsg == "" {
		t.Fatalf("notifyUser not called")
	}
	if notifiedUser != user {
		t.Errorf("Уведомлен (%s) вместо %s", notifiedUser, user)
	}
	const wantSubstring = "98% of your quota"
	if !strings.Contains(notifiedMsg, wantSubstring) {
		t.Errorf("неожиданное уведомление <<%s>>, "+"want substring %q", notifiedMsg, wantSubstring)
	}
}
