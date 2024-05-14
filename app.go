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
	StagedforFiling time.Time `json:"StagedforFiling"`
	SubmitterName string `json:"SubmitterName"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Download the unstamped oath attachment
func (a *App) DownloadAttachment(id int) string {
	sp := getClient()

	// Load oath review item
	item := sp.Web().
		GetList("Lists/OathOfOfficeReviews1").
		Items().
		GetByID(id)

	// Load attachments
	attachments, _ := item.
		Attachments().
		Get()

	// Get the first attachment (unstamped oath)
	pdfFilename := attachments.Data()[0].Data().FileName
	attachment, _ := item.Attachments().GetByName(pdfFilename).Download()

	// Return the Base64 encoded unstamped oath
	return base64.StdEncoding.EncodeToString(attachment)
}

func (a *App) LoadUnstamped() []Oath {
	sp := getClient()

	// Load unstamped oath review items
	listItems, _ := sp.Web().
		GetList("Lists/OathOfOfficeReviews1").
		Items().
		Select("Id,CreationDate,StagedforFiling,SubmitterName").
		Filter("StagedforFiling eq null and Filing/Determination eq 'Accepted'").
		Get()

	// Unmarshal the JSON into a Go struct
	items := []Oath{}
	json.Unmarshal(listItems.Normalized(), &items)

	// Return the list of unstamped oath review items
	return items
}

func (a *App) SignIn() *api.UserInfo {
	// Set the authentication strategy
	// NOTE:  This is separate from the private getClient method because we need to reuse
	// the client to clear the cookie cache if there is an error
	authCnfg := &strategy.AuthCnfg {
		SiteURL: "https://nysemail.sharepoint.com/sites/DOS/corp/Data",
	}

	// Create the SharePoint client
	client := &gosip.SPClient{AuthCnfg: authCnfg}
	sp := api.NewSP(client)

	// Get the current user; if there is an error, clear the cookie cache and try again
	response, err := sp.Web().CurrentUser().Get()
	if err != nil {
		authCnfg.CleanCookieCache()
		sp = api.NewSP(client)
		response, _ = sp.Web().CurrentUser().Get()
	}

	// Return the current user
	return response.Data()
}

// Upload the stamped oath attachment
func (a *App) UploadStamped(id int, stamped string) error {
	// Decode the Base64 encoded stamped oath
	pdfArray, _ := base64.StdEncoding.DecodeString(stamped)
	pdf := bytes.NewReader(pdfArray)

	// Get the oath review item
	sp := getClient()
	item := sp.Web().GetList("Lists/OathOfOfficeReviews1").Items().GetByID(id)

	// Remove unstamped oath attachment
	attachments, _ := item.Attachments().Get()
	attachment := attachments.Data()[0].Data().FileName
	err := item.Attachments().GetByName(attachment).Delete()
	if err != nil {
		err := item.Attachments().GetByName(attachment).Delete()
		if err != nil {
			return err
		}
	}

	// Add stamped oath attachment
	_, err = item.Attachments().Add("stamped.pdf", pdf)
	if err != nil {
		_, err = item.Attachments().Add("stamped.pdf", pdf)
		if err != nil {
			return err
		}
	}

	// Update `StagedforFiling` timestamp
	_, err = item.Update(
		[]byte(
			fmt.Sprintf(
				`{"StagedforFiling": "%s"}`,
				time.Now().Format(time.RFC3339),
			),
		),
	)

	return err
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
