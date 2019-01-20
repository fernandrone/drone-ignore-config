package plugin

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"encoding/base64"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
)

// empty context
var noContext = context.Background()

// template for drone ignore file
const droneignoreTemplate = `{
	"type": "file",
	"encoding": "base64",
	"size": 5362,
	"name": ".droneignore",
	"path": ".droneignore",
	"content": "%s",
	"sha": "3d21ec53a331a6f037a91c368710b99387d012c1",
	"url": "https://api.github.com/repos/octocat/hello-world/contents/.droneignore",
	"git_url": "https://api.github.com/repos/octocat/hello-world/git/blobs/3d21ec53a331a6f037a91c368710b99387d012c1",
	"html_url": "https://github.com/octocat/hello-world/blob/master/.droneignore",
	"download_url": "https://raw.githubusercontent.com/octocat/hello-world/master/.droneignore",
	"_links": {
		"git": "https://api.github.com/repos/octocat/hello-world/git/blobs/3d21ec53a331a6f037a91c368710b99387d012c1",
		"self": "https://api.github.com/repos/octocat/hello-world/contents/.droneignore",
		"html": "https://github.com/octocat/hello-world/blob/master/.droneignore"
	}
}`

// template for commit information
const commitInfoTemplate = `{
	"url": "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
	"sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e",
	"node_id": "MDY6Q29tbWl0NmRjYjA5YjViNTc4NzVmMzM0ZjYxYWViZWQ2OTVlMmU0MTkzZGI1ZQ==",
	"html_url": "https://github.com/octocat/Hello-World/commit/6dcb09b5b57875f334f61aebed695e2e4193db5e",
	"comments_url": "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e/comments",
	"commit": {
	  "url": "https://api.github.com/repos/octocat/Hello-World/git/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
	  "author": {
		"name": "Monalisa Octocat",
		"email": "support@github.com",
		"date": "2011-04-14T16:00:49Z"
	  },
	  "committer": {
		"name": "Monalisa Octocat",
		"email": "support@github.com",
		"date": "2011-04-14T16:00:49Z"
	  },
	  "message": "Fix all the bugs",
	  "tree": {
		"url": "https://api.github.com/repos/octocat/Hello-World/tree/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		"sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e"
	  },
	  "comment_count": 0,
	  "verification": {
		"verified": false,
		"reason": "unsigned",
		"signature": null,
		"payload": null
	  }
	},
	"author": {
	  "login": "octocat",
	  "id": 1,
	  "node_id": "MDQ6VXNlcjE=",
	  "avatar_url": "https://github.com/images/error/octocat_happy.gif",
	  "gravatar_id": "",
	  "url": "https://api.github.com/users/octocat",
	  "html_url": "https://github.com/octocat",
	  "followers_url": "https://api.github.com/users/octocat/followers",
	  "following_url": "https://api.github.com/users/octocat/following{/other_user}",
	  "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
	  "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
	  "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
	  "organizations_url": "https://api.github.com/users/octocat/orgs",
	  "repos_url": "https://api.github.com/users/octocat/repos",
	  "events_url": "https://api.github.com/users/octocat/events{/privacy}",
	  "received_events_url": "https://api.github.com/users/octocat/received_events",
	  "type": "User",
	  "site_admin": false
	},
	"committer": {
	  "login": "octocat",
	  "id": 1,
	  "node_id": "MDQ6VXNlcjE=",
	  "avatar_url": "https://github.com/images/error/octocat_happy.gif",
	  "gravatar_id": "",
	  "url": "https://api.github.com/users/octocat",
	  "html_url": "https://github.com/octocat",
	  "followers_url": "https://api.github.com/users/octocat/followers",
	  "following_url": "https://api.github.com/users/octocat/following{/other_user}",
	  "gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
	  "starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
	  "subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
	  "organizations_url": "https://api.github.com/users/octocat/orgs",
	  "repos_url": "https://api.github.com/users/octocat/repos",
	  "events_url": "https://api.github.com/users/octocat/events{/privacy}",
	  "received_events_url": "https://api.github.com/users/octocat/received_events",
	  "type": "User",
	  "site_admin": false
	},
	"parents": [
	  {
		"url": "https://api.github.com/repos/octocat/Hello-World/commits/6dcb09b5b57875f334f61aebed695e2e4193db5e",
		"sha": "6dcb09b5b57875f334f61aebed695e2e4193db5e"
	  }
	],
	"stats": {
	  "additions": 104,
	  "deletions": 4,
	  "total": 108
	},
	"files": [ %s ]
  }`

const commitFilesTemplate = `{
	"filename": "%s",
	"additions": 10,
	"deletions": 2,
	"changes": 12,
	"status": "modified",
	"raw_url": "https://github.com/octocat/Hello-World/raw/7ca483543807a51b6079e54ac4cc392bc29ae284/file1.txt",
	"blob_url": "https://github.com/octocat/Hello-World/blob/7ca483543807a51b6079e54ac4cc392bc29ae284/file1.txt",
	"patch": "@@ -29,7 +29,7 @@\n....."
  }`

type TestServer struct {
	*httptest.Server
	getCommitResponse   []byte
	getContentsResponse []byte
}

func NewTestServer() *TestServer {
	ts := &TestServer{}
	ts.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "contents") {
			w.Write(ts.getContentsResponse)
		}
		if strings.Contains(r.URL.Path, "commits") {
			w.Write(ts.getCommitResponse)
		}
	}))
	return ts
}

func (ts *TestServer) SetGetCommitResponse(filenames []string) {
	var contents []string
	for _, filename := range filenames {
		contents = append(contents, fmt.Sprintf(commitFilesTemplate, filename))
	}

	ts.getCommitResponse = []byte(fmt.Sprintf(commitInfoTemplate, strings.Join(contents, ",")))
}

func (ts *TestServer) SetGetContentsResponse(patterns []string) {
	content := strings.Join(patterns, "\n")
	b64Content := base64.URLEncoding.EncodeToString([]byte(content))

	ts.getContentsResponse = []byte(fmt.Sprintf(droneignoreTemplate, b64Content))
}

func TestPlugin(t *testing.T) {
	log.AddHook(filename.NewHook())

	req := &config.Request{
		Build: drone.Build{
			After: "3d21ec53a331a6f037a91c368710b99387d012c1",
		},
		Repo: drone.Repo{
			Slug:      "octocat/hello-world",
			Name:      "hello-world",
			Namespace: "octocat",
		},
	}

	ts := NewTestServer()
	defer ts.Close()

	plugin := New(ts.URL, "d7c559e677ebc489d4e0193c8b97a12e")

	var tests = []struct {
		ignorePatterns []string
		changedFiles   []string
		shouldIgnore   bool
	}{
		{[]string{""}, []string{""}, false},
		{[]string{""}, []string{"src/README.md"}, false},
		{[]string{"*"}, []string{""}, true},
		{[]string{"*"}, []string{"src/README.md"}, true},
		{[]string{"*", "!src/README.md"}, []string{"src/README.md"}, false},
		{[]string{"src/README.md"}, []string{"src/README.md", "src/DCOS.md"}, false},
		{[]string{"*.md"}, []string{"src/README.md", "src/DCOS.md"}, true},
		{[]string{"src"}, []string{"src/README.md", "src/DCOS.md"}, true},
		{[]string{"src/**"}, []string{"src/README.md", "src/DCOS.md"}, true},
		{[]string{"src/**"}, []string{"src/README.md", "src/DCOS.md", "data/file.txt"}, false},
		{[]string{"src/README.md", "src/DCOS.md"}, []string{"src/README.md", "src/DCOS.md"}, true},
		{[]string{"src/README.md", "src/DCOS.md"}, []string{"src/README.md"}, true},
		{[]string{"src/README.md", "src/DCOS.md"}, []string{"src/DCOS.md"}, true},
	}
	for _, test := range tests {
		ts.SetGetContentsResponse(test.ignorePatterns)
		ts.SetGetCommitResponse(test.changedFiles)

		_, err := plugin.Find(noContext, req)

		var got bool

		if err != nil && err.Error() == ".droneignore match, skipping build" {
			got = true
		} else if err != nil {
			t.Errorf("p.Find(ignorePatterns: %q, changedFiles: %q):\n%v", test.ignorePatterns, test.changedFiles, err)
		}

		if got != test.shouldIgnore {
			t.Errorf("p.Find(ignorePatterns: %q, changedFiles: %q):\nExpected \"%t\", got \"%t\"", test.ignorePatterns, test.changedFiles, test.shouldIgnore, got)
		}
	}
}
