package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/caos/logging"

	caos_errors "github.com/caos/zitadel/internal/errors"

	"github.com/k3a/html2text"

	"github.com/caos/zitadel/internal/notification/channels"
	"github.com/caos/zitadel/internal/notification/messages"
)

func InitFSChannel(path string, config FSConfig) (channels.NotificationChannel, error) {
	if path == "" {
		return nil, nil
	}
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}

	logging.Log("NOTIF-kSvPp").Debug("successfully initialized filesystem email and sms channel")

	return channels.HandleMessageFunc(func(message channels.Message) error {

		fileName := fmt.Sprintf("%d_", time.Now().Unix())
		content := message.GetContent()
		switch msg := message.(type) {
		case *messages.Email:
			recipients := make([]string, len(msg.Recipients))
			copy(recipients, msg.Recipients)
			sort.Strings(recipients)
			fileName = fileName + "mail_to_" + strings.Join(recipients, "_") + ".html"
			if config.Compact {
				content = html2text.HTML2Text(content)
			}
		case *messages.SMS:
			fileName = fileName + "sms_to_" + msg.RecipientPhoneNumber + ".txt"
		default:
			return caos_errors.ThrowUnimplementedf(nil, "NOTIF-6f9a1", "filesystem provider doesn't support message type %T", message)
		}

		return ioutil.WriteFile(filepath.Join(path, fileName), []byte(content), 0666)
	}), nil
}
