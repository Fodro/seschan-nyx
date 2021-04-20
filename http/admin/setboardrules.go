package admin

import (
	"net/http"

	"github.com/tidwall/buntdb"
	"go.fodro/nyx/http/errw"
	"go.fodro/nyx/http/middle"
	"go.fodro/nyx/resources"
)

func handleSetRules(w http.ResponseWriter, r *http.Request) {
	sess := middle.GetSession(r)
	if !middle.IsAdminSession(sess) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}
	db := middle.GetDB(r)

	boardName := r.FormValue("shortname")
	rules := r.FormValue("rules")

	if err = db.Update(func(tx *buntdb.Tx) error {
		board, err := resources.GetBoard(tx, r.Host, boardName)
		if err != nil {
			return err
		}
		board.Metadata["rules"] = rules
		return resources.UpdateBoard(tx, r.Host, board)
	}); err != nil {
		errw.ErrorWriter(err, w, r)
		return
	}

	http.Redirect(w, r, "/admin/panel.html", http.StatusSeeOther)
}
