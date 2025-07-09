package gconf

const (
	// https://github.com/gorilla/securecookie
	HttpCookieProvider_GORILLA = 0

	HttpCookieSameSiteOption_STRICT = 0
	HttpCookieSameSiteOption_LAX    = 1
	HttpCookieSameSiteOption_NONE   = 2
)

type HttpCookie struct {
	Provider int    `json:"provider,omitempty" yaml:"provider,omitempty"`
	Name     string `json:"name,omitempty" yaml:"name,omitempty"`
	Expires  string `json:"expires,omitempty" yaml:"expires,omitempty"`
	MaxAge   int    `json:"max_age,omitempty" yaml:"max_age,omitempty"`
	Domain   string `json:"domain,omitempty" yaml:"domain,omitempty"`
	Path     string `json:"path,omitempty" yaml:"path,omitempty"`
	Secure   bool   `json:"secure,omitempty" yaml:"secure,omitempty"`
	HttpOnly bool   `json:"http_only,omitempty" yaml:"http_only,omitempty"`
	SameSite int    `json:"same_site,omitempty" yaml:"same_site,omitempty"`
	HashKey  string `json:"hash_key,omitempty" yaml:"hash_key,omitempty"`
	BlockKey string `json:"block_key,omitempty" yaml:"block_key,omitempty"`
}

const (
	// https://github.com/gorilla/csrf
	CsrfProvider_GORILLA = 0
)

type Csrf struct {
	Provider int         `json:"provider,omitempty" yaml:"provider,omitempty"`
	Key      string      `json:"key,omitempty" yaml:"key,omitempty"`
	Header   string      `json:"header,omitempty" yaml:"header,omitempty"`
	Cookie   *HttpCookie `json:"cookie,omitempty" yaml:"cookie,omitempty"`
}

type HttpSession struct {
	Handler int         `json:"handler,omitempty" yaml:"handler,omitempty"`
	Header  string      `json:"header,omitempty" yaml:"header,omitempty"`
	Cookie  *HttpCookie `json:"cookie,omitempty" yaml:"cookie,omitempty"`
}

type HttpServer struct {
	Bind      string       `json:"bind,omitempty" yaml:"bind,omitempty"`
	Port      int          `json:"port,omitempty" yaml:"port,omitempty"`
	Hostnames []string     `json:"hostnames,omitempty" yaml:"hostnames,omitempty"`
	Tls       *Tls         `json:"tls,omitempty" yaml:"tls,omitempty"`
	Csrf      *Csrf        `json:"csrf,omitempty" yaml:"csrf,omitempty"`
	Session   *HttpSession `json:"session,omitempty" yaml:"session,omitempty"`
}
