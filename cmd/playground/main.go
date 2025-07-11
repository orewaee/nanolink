package main

import (
	"context"
	"fmt"
	"time"

	"github.com/orewaee/nanolink/internal/core/domain"
	"github.com/orewaee/nanolink/internal/core/driven/repo"
	"github.com/orewaee/nanolink/internal/datasource/disk"
)

func main() {
	var linkRepo repo.LinkRepo
	linkRepo = &disk.YamlLinkRepo{
		Dir: "links",
	}

	fmt.Println(linkRepo.AddLink(context.TODO(), domain.Link{Id: "github", Location: "https://github.com/orewaee", CreatedAt: time.Now()}))
	fmt.Println(linkRepo.GetLinkById(context.TODO(), "github"))
}
