package main

import (
	"context"
	"fmt"

	"github.com/orewaee/nanolink/internal/core/driven/repo"
	"github.com/orewaee/nanolink/internal/datasource/disk"
)

func main() {
	var linkRepo repo.LinkRepo
	linkRepo = &disk.YamlLinkRepo{
		Dir: "links",
	}
	fmt.Println(linkRepo.GetLinkById(context.TODO(), "github"))
}
