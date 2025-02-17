package board

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/pressly/chi"
	"github.com/tidwall/buntdb"
	"go.fodro/nyx/http/errw"
	"go.fodro/nyx/http/middle"
	"go.fodro/nyx/resources"
)

func serveBoard(w http.ResponseWriter, r *http.Request) {
	dat := bytes.NewBuffer([]byte{})
	db := middle.GetDB(r)
	ctx := middle.GetBaseCtx(r)
	err := db.View(func(tx *buntdb.Tx) error {
		bName := chi.URLParam(r, "board")
		log.Println("Getting board")
		b, err := resources.GetBoard(tx, r.Host, bName)
		if err != nil {
			return err
		}
		ctx["Board"] = b

		log.Println("Listing Threads...")
		threads, err := resources.ListThreads(tx, r.Host, bName)
		if err != nil {
			return err
		}
		log.Println("Number of Thread on board: ", len(threads))

		log.Println("Filling threads")
		for k := range threads {
			err := resources.FillReplies(tx, r.Host, threads[k])
			if err != nil {
				return err
			}
		}
		ctx["Threads"] = threads
		bList, err := resources.ListBoards(tx, r.Host)
		if err != nil {
			return err
		}
		ctx["Boards"] = bList
		return nil
	})
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	err = tmpls.ExecuteTemplate(dat, "board/board", ctx)
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	http.ServeContent(w, r, "board.html", time.Now(), bytes.NewReader(dat.Bytes()))
}
