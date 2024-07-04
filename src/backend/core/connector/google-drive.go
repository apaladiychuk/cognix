package connector

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"strings"

	drive "google.golang.org/api/drive/v2"

	"net/http"

	"google.golang.org/api/option"
	"log"
)

func GetDriver(token *oauth2.Token) {
	ctx := context.Background()

	// If modifying these scopes, delete your previously saved token.json.
	srv, err := drive.NewService(ctx,
		option.WithHTTPClient(&http.Client{Transport: utils.NewTransport(token)}))
	if err != nil {
		log.Fatalf("Unable to retrieve driveactivity Client %v", err)
	}

	//Fields("nextPageToken, files(id, name)").
	r, err := srv.Drives.List().Do()
	if err != nil {
		log.Fatalf("Unable to retrieve list of activities. %v", err)
	}
	fmt.Println("Recent Activity:")
	for _, dr := range r.Items {
		fmt.Printf(" id %s name %s\n", dr.Id, dr.Name)
	}
	fr, err := srv.Files.List().Q("mimeType = 'application/vnd.google-apps.folder'").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve list of activities. %v", err)
	}

	for _, f := range fr.Items {
		fullpath := []string{}
		for _, p := range f.Parents {
			fullpath = append(fullpath, p.Id)
		}

		fmt.Printf("folder %s (%s)\n parent %s  \n", f.Title, f.Id, strings.Join(fullpath, ","))
		// ${dataSource.folder_id}' in parents
		nfr, err := srv.Files.List().Q(fmt.Sprintf("'%s' in parents", f.Id)).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve list of activities. %v", err)
		}
		for _, d := range nfr.Items {
			fmt.Printf("\t---\t\t%s- %s \n", d.Id, d.Title)
			resp, err := srv.Files.Get(d.Id).Download()
			if err != nil {
				log.Fatalf("Unable to retrieve list of activities. %v", err)
			}
			resp.
		}
		fmt.Printf("\t %s- %s \n", f.Id, f.Title)
	}
}
