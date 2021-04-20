package board

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/pressly/chi"
	"github.com/tidwall/buntdb"
	"go.fodro/nyx/http/errw"
	"go.fodro/nyx/http/middle"
	"go.fodro/nyx/resources"
)

func serveThread(w http.ResponseWriter, r *http.Request) {
	dat := bytes.NewBuffer([]byte{})
	db := middle.GetDB(r)
	ctx := middle.GetBaseCtx(r)
	err := db.View(func(tx *buntdb.Tx) error {
		bName := chi.URLParam(r, "board")
		b, err := resources.GetBoard(tx, r.Host, bName)
		if err != nil {
			return err
		}
		ctx["Board"] = b

		id, err := strconv.Atoi(chi.URLParam(r, "thread"))
		if err != nil {
			return err
		}
		thread, err := resources.GetThread(tx, r.Host, bName, id)
		if err != nil {
			return err
		}

		err = resources.FillReplies(tx, r.Host, thread)
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}
		ctx["Thread"] = thread
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
	err = tmpls.ExecuteTemplate(dat, "board/thread", ctx)
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	http.ServeContent(w, r, "board.html", time.Now(), bytes.NewReader(dat.Bytes()))
}
