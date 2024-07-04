package connector

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	"fmt"
	"golang.org/x/oauth2"

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

	//q := driveactivity.QueryDriveActivityRequest{PageSize: 10}
	r, err := srv.Drives.List().Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve list of activities. %v", err)
	}
	fmt.Println("Recent Activity:")
	for _, dr := range r.Items {
		fmt.Printf(" id %s name %s\n", dr.Id, dr.Name)
	}
}
