package main

import (
		"fmt"
		"labix.org/v2/mgo"
		//"labix.org/v2/mgo/bson"
		"net/http"
		"text/template"
)

type Post struct {
		Content string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//array of Posts
	entries := []Post{}
	err := posts.Find(nil).Limit(12).All(&entries)
	if err != nil {
		fmt.Fprintf(w, "err: %s!", err)
		http.NotFound(w, r)
	}

	ctx := map[string]interface{} {
		"posts":entries,
	}
	indexTemplate.Execute(w, ctx)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	// hardcore security right here ;-)
	if r.Method == "POST"  && (r.FormValue("pw") == "pw") {
		//fmt.Fprintf(w, "%s", r.FormValue("content"))
		content := r.FormValue("content")
		newPost := &Post{Content:content}
		//fmt.Fprintf(w, "%s", newPost)

		//insert new Post
		posts.Insert(newPost)

		http.Redirect(w, r, "/", 302)
	} else {
		//if its a GET, show the input form
		createTemplate.Execute(w, http.StatusFound/*302*/)
	}
}

var posts *mgo.Collection

func main() {
		session, err := mgo.Dial("localhost")
		if err != nil {
				panic(err)
		}
		defer session.Close()
		
		// Optional. Switch the session to a monotonic behavior.
		//session.SetMode(mgo.Monotonic, true)

		posts = session.DB("blog").C("post")

		result := Post{}
		err = posts.Find(nil).One(&result)
		if err != nil {
				panic(err)
		}

		http.HandleFunc("/", indexHandler)
		http.HandleFunc("/create", createHandler)
		http.ListenAndServe(":8081", nil)
}

var indexTemplate = template.Must(template.New("index").Parse(`
<html>
	<body>
		<h2>Stuff... <a href="mailto:stuff@nerdporn.org">stuff@nerdporn.org</a></h2>
		<ul>
			{{range .posts}}
			<li>
			{{.Content}}
			</li>
		{{end}}
		</ul>
	</body>
</html>
`))

var createTemplate = template.Must(template.New("create").Parse(`
<html>
	<body>
		<h2>Stuff... <a href="mailto:stuff@nerdporn.org">stuff@nerdporn.org</a></h2>
		<div>
			<form action="create" method="post">
				<textarea name="content" rows="4" cols="40"><a href=""></a></textarea>
				<input type="input" name="pw" value="pw">
				<input type="submit" value="Submit">
			</form>
		</div>
	</body>
</html>
`))
