package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/koltyakov/gosip"
	strategy "github.com/koltyakov/gosip-sandbox/strategies/ondemand"
	"github.com/koltyakov/gosip/api"
)

// App struct
type App struct {
	ctx context.Context
}

type Oath struct {
	ID int `json:"Id"`
	CreationDate string `json:"CreationDate"`
	StagedforFiling bool `json:"StagedforFiling"`
	SubmitterName string `json:"SubmitterName"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) DownloadAttachment(id int) string {
	sp := getClient()

	item := sp.Web().
		GetList("Lists/OathOfOfficeReviews1").
		Items().
		GetByID(id)

	attachments, _ := item.
		Attachments().
		Get()

	pdfFilename := attachments.Data()[0].Data().FileName
	attachment, _ := item.Attachments().GetByName(pdfFilename).Download()

	return base64.StdEncoding.EncodeToString(attachment)
}

func (a *App) LoadUnstamped() []Oath {
	sp := getClient()

	listItems, _ := sp.Web().
		GetList("Lists/OathOfOfficeReviews1").
		Items().
		Select("Id,CreationDate,StagedforFiling,SubmitterName").
		Filter("(StagedforFiling eq null) and (Filing_x003a_Determination/Value eq 'Approved')").
		Get()

	items := []Oath{}
	json.Unmarshal(listItems.Normalized(), &items)

	return items
}

func (a *App) SignIn() *api.UserInfo {
	authCnfg := &strategy.AuthCnfg {
		SiteURL: "https://nysemail.sharepoint.com/sites/DOS/corp/Data",
	}

	client := &gosip.SPClient{AuthCnfg: authCnfg}
	sp := api.NewSP(client)

	response, err := sp.Web().CurrentUser().Get()
	if err != nil {
		authCnfg.CleanCookieCache()
		sp = api.NewSP(client)
		response, _ = sp.Web().CurrentUser().Get()
	}

	return response.Data()
}

func (a *App) UploadStamped(id int, stamped string) {
	pdfArray, _ := base64.StdEncoding.DecodeString(stamped)
	pdf := bytes.NewReader(pdfArray)
	sp := getClient()

	item := sp.Web().GetList("Lists/OathOfOfficeReviews1").Items().GetByID(id)

	// Remove old attachment
	attachments, _ := item.Attachments().Get()
	attachment := attachments.Data()[0].Data().FileName
	item.Attachments().GetByName(attachment).Delete()

	// Add new attachment
	item.Attachments().Add("stamped.pdf", pdf)

	// Update timestamp
	_, err := item.Update(
		[]byte(
			fmt.Sprintf(`{"StagedforFiling": "%s"}`, time.Now().Format(time.RFC3339)),
		),
	)

	if err != nil {
		panic(err)
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func getClient() *api.SP {
	auth := &strategy.AuthCnfg {
		SiteURL: "https://nysemail.sharepoint.com/sites/DOS/corp/Data",
	}

	client := &gosip.SPClient{AuthCnfg: auth}
	return api.NewSP(client)
}
