//
//  dsdump.go
//
//  Created by Jens-Uwe Mager <jum@anubis.han.de> on 07.01.14.
//

// Simple dump and restore of all records in the appengine datastore.
package dsdump

import (
	"appengine"
	"appengine/datastore"
	"encoding/gob"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Set DEBUG to true for lots of debug output
const DEBUG = false

func debug(c appengine.Context, format string, a ...interface{}) {
	if DEBUG {
		c.Debugf(format, a...)
	}
}

func init() {
	http.HandleFunc("/admin/dsdump", dsdump)
	gob.Register(&time.Time{})
}

type DSRec struct {
	Key   *datastore.Key
	Props datastore.PropertyList
}

func dsdump(w http.ResponseWriter, r *http.Request) {
	var err error
	c := appengine.NewContext(r)
	debug(c, "req url %#v", r.URL)
	debug(c, "req header %#v", r.Header)
	switch r.Method {
	case "GET":
		it := datastore.NewQuery("").Run(c)
		enc := gob.NewEncoder(w)
		w.Header().Set("Content-Type", "application/x-gob")
		for {
			var rec DSRec
			if rec.Key, err = it.Next(&rec.Props); err != nil {
				if err != datastore.Done {
					c.Errorf("query next: %v", err.Error())
				}
				break
			}
			if rec.Key.Kind()[0] == '_' {
				// skip over internal datastore management records
				continue
			}
			err = enc.Encode(&rec)
			if err != nil {
				c.Errorf("encode rec: %v", err.Error())
				break
			}
		}
	case "POST":
		dec := gob.NewDecoder(r.Body)
		num := 0
		for {
			var rec DSRec
			if err = dec.Decode(&rec); err != nil {
				if err == io.EOF {
					break
				}
				c.Errorf("decode gob post: %v", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			debug(c, "req val: %#v", rec)
			_, err = datastore.Put(c, rebuildKey(c, rec.Key), &rec.Props)
			if err != nil {
				c.Errorf("put rec %+v: %v", rec, err.Error())
			} else {
				num++
			}
		}
		fmt.Fprintf(w, "updated %d records, do not forget to flush memcache\n", num)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func rebuildKey(c appengine.Context, key *datastore.Key) *datastore.Key {
	if key == nil {
		return nil
	}
	return datastore.NewKey(c, key.Kind(), key.StringID(), key.IntID(), rebuildKey(c, key.Parent()))
}
