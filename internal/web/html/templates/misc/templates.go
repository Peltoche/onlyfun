package misc

type NotFoundPageTmpl struct{}

func (t *NotFoundPageTmpl) Template() string { return "misc/page_404" }
