package repo

type Repo struct {
	name   string
	svnUrl string
}

func (r *Repo) GetName() string {
	return r.name
}

func (r *Repo) GetSvnUrl() string {
	return r.svnUrl
}

func CreateRepo(name, snvUrl string) Repo {
	return Repo{name, snvUrl}
}
