package indexer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jaytaylor/html2text"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"time"
	"tryffel.net/go/meilindex/config"
)

type mailSpringPerson struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (p *mailSpringPerson) String() string {
	if p.Name != "" {
		return p.Name
	}
	return p.Email
}

type mailSpringMessageData struct {
	Date   int64 `json:"date"`
	Folder struct {
		Path string `json:"path"`
	} `json:"folder"`

	From      []mailSpringPerson `json:"from"`
	Id        string             `json:"id"`
	PlainText bool               `json:"plaintext"`
	To        []mailSpringPerson `json:"to"`
	Bcc       []mailSpringPerson `json:"bcc"`
	CC        []mailSpringPerson `json:"cc"`
}

func (m *mailSpringMessageData) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return fmt.Errorf("expected source to be string, got %v", src)
	}
	return json.Unmarshal([]byte(str), m)
}

// mailspring mail as joined from db.Message, db.MessageBody.
type mailSpringMail struct {
	Id      string                 `db:"id"`
	MailId  string                 `db:"headerMessageId"`
	Subject string                 `db:"subject"`
	Data    *mailSpringMessageData `db:"data"`
	Body    sql.NullString         `db:"body"`
}

func (msMail *mailSpringMail) ToMail() *Mail {
	mail := &Mail{
		Uid:       msMail.Id,
		Id:        msMail.Id,
		Subject:   msMail.Subject,
		Body:      msMail.Body.String,
		Timestamp: time.Unix(msMail.Data.Date, 0),
		Folder:    msMail.Data.Folder.Path,
	}

	if len(msMail.Data.From) > 0 {
		if msMail.Data.From[0].Name != "" {
			mail.From = msMail.Data.From[0].Name
		} else {
			mail.From = msMail.Data.From[0].Email
		}
	}

	mail.To = make([]string, len(msMail.Data.To))
	for i, v := range msMail.Data.To {
		if v.Name != "" {
			mail.To[i] = v.Name
		} else {
			mail.To[i] = v.Email
		}
	}
	mail.Cc = make([]string, len(msMail.Data.CC))
	for i, v := range msMail.Data.CC {
		if v.Name != "" {
			mail.Cc[i] = v.Name
		} else {
			mail.Cc[i] = v.Email
		}
	}

	if msMail.Body.String != "" {
		plainText, err := html2text.FromString(msMail.Body.String, html2text.Options{
			PrettyTables: false,
		})
		if err != nil {
			logrus.Warningf("format html as plain text (mail %s): %v", msMail.Id, err)
		} else {
			mail.Body = plainText
		}
	}
	return mail
}

func ReadMailspring(file string, recursive bool, flushFunc func(mails []*Mail) error) ([]*Mail, error) {
	logrus.Infof("open mailspring database %s", file)
	db, err := sqlx.Open("sqlite3", fmt.Sprintf("%s?mode=ro", file))
	if err != nil {
		return nil, fmt.Errorf("open database: %v", err)
	}

	defer db.Close()

	batchSize := config.Conf.File.BatchSize
	rawMails := make([]mailSpringMail, 0, batchSize)

	var totalMails int

	row := db.QueryRowx("select count(id) as count from Message;")
	err = row.Scan(&totalMails)
	if err != nil {
		return nil, fmt.Errorf("get total rawMails count: %v", err)
	}

	pages := totalMails / batchSize
	if totalMails%batchSize != 0 {
		pages += 1
	}

	mailSql := `
select
       m.id ,
       m.headerMessageId,
       m.subject,
       m.data,
       body.value as body
from Message as m
left join MessageBody as body on m.id = body.id
order by m.id asc
limit ?
offset ?;

`

	for page := 0; page < pages; page++ {
		err = db.Select(&rawMails, mailSql, batchSize, page*batchSize)
		if err != nil {
			return nil, fmt.Errorf("read rawMails, page: %d: %v", page, err)
		}

		mails := make([]*Mail, len(rawMails))

		for i, v := range rawMails {
			mails[i] = v.ToMail()
		}

		err = flushFunc(mails)
		if err != nil {
			logrus.Errorf("flush mails (page %d): %v", page, err)
		}
	}
	return nil, nil
}
