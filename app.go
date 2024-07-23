package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/koltyakov/gosip"
	strategy "github.com/koltyakov/gosip-sandbox/strategies/ondemand"
	"github.com/koltyakov/gosip/api"
)

const StampingList = "Lists/Stampable%20Documents"

// App struct
type App struct {
	ctx context.Context
}

type SignInResponse struct {
	CanAccess      bool          `json:"CanAccess"`
	CurrentVersion string        `json:"CurrentVersion"`
	UserInfo       *api.UserInfo `json:"UserInfo"`
}

type StampingItem struct {
	ID            int    `json:"Id"`
	CreationDate  string `json:"Created"`
	Selected      bool   `json:"Selected"`
	StampText     string `json:"StampText"`
	SubmitterName string `json:"Title"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Download the unstamped attachment
func (a *App) DownloadAttachment(id int) string {
	sp := getClient()

	// Load stampable item
	item := sp.Web().
		GetList(StampingList).
		Items().
		GetByID(id)

	// Load attachments
	attachments, _ := item.
		Attachments().
		Get()

	// Get the first attachment (unstamped item)
	pdfFilename := attachments.Data()[0].Data().FileName
	attachment, _ := item.Attachments().GetByName(pdfFilename).Download()

	// Return the Base64 encoded unstamped item
	return base64.StdEncoding.EncodeToString(attachment)
}

func (a *App) LoadUnstamped() []StampingItem {
	sp := getClient()

	// Load unstamped review items
	listItems, _ := sp.Web().
		GetList(StampingList).
		Items().
		Select("Id,Title,StampText,Created").
		Filter("Processed eq null").
		Get()

	// Unmarshal the JSON into a Go struct
	items := []StampingItem{}
	json.Unmarshal(listItems.Normalized(), &items)

	// Return the list of unstamped review items
	return items
}

func (a *App) SignIn() SignInResponse {
	// Set the authentication strategy
	// NOTE:  This is separate from the private getClient method because we need to reuse
	// the client to clear the cookie cache if there is an error
	authCnfg := &strategy.AuthCnfg{
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

	// Get the current version of the application for enforcement
	versionResponse, _ := sp.Web().
		GetList("Lists/StampingApplication").
		Items().
		Filter("Title eq 'CurrentVersion'").
		Top(1).
		Get()

	version := versionResponse.ToMap()[0]["Value"].(string)

	// Check whether the user has the appropriate permissions
	http := api.NewHTTPClient(client)
	endpoint := authCnfg.GetSiteURL() + "/_api/web/lists(guid'53774fcb-80f2-48f3-b2b5-fdaa781c7b75')/effectiveBasePermissions"
	basePermissionsResponse, _ := http.Get(endpoint, nil)
	rawBasePermissions := api.NormalizeODataItem(basePermissionsResponse)

	// Unmarshal the base permissions (this obfuscation here seems to be due to SharePoint returning the integers as strings)
	unmarshaledBasePermissions := &struct {
		EffectiveBasePermissions struct {
			High string `json:"High"`
			Low  string `json:"Low"`
		} `json:"EffectiveBasePermissions"`
	}{}
	json.Unmarshal(rawBasePermissions, &unmarshaledBasePermissions)

	// Convert the base permissions to integers
	high, _ := strconv.ParseInt(unmarshaledBasePermissions.EffectiveBasePermissions.High, 10, 64)
	low, _ := strconv.ParseInt(unmarshaledBasePermissions.EffectiveBasePermissions.Low, 10, 64)

	// Load and check permissions
	basePermissions := api.BasePermissions{
		High: high,
		Low:  low,
	}
	canAccess := api.HasPermissions(basePermissions, api.PermissionKind.EditListItems)

	return SignInResponse{
		CanAccess:      canAccess,
		CurrentVersion: version,
		UserInfo:       response.Data(),
	}
}

// Upload the stamped attachment
func (a *App) UploadStamped(id int, stamped string) error {
	// Decode the Base64 encoded stamped attachment
	pdfArray, _ := base64.StdEncoding.DecodeString(stamped)
	pdf := bytes.NewReader(pdfArray)

	// Get the review item
	sp := getClient()
	item := sp.Web().
		GetList(StampingList).
		Items().
		GetByID(id)

	// Remove unstamped attachment
	attachments, _ := item.Attachments().Get()
	attachment := attachments.Data()[0].Data().FileName
	err := item.Attachments().GetByName(attachment).Delete()
	if err != nil {
		err := item.Attachments().GetByName(attachment).Delete()
		if err != nil {
			return err
		}
	}

	// Add stamped attachment
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
				`{"Processed": "%s"}`,
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
	auth := &strategy.AuthCnfg{
		SiteURL: "https://nysemail.sharepoint.com/sites/DOS/corp/Data",
	}

	client := &gosip.SPClient{AuthCnfg: auth}
	return api.NewSP(client)
}
